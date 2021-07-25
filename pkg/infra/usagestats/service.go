package usagestats

import (
	"context"
	"fmt"
	"time"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/login/social"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/services/alerting"
	"github.com/openinsight-project/grafinsight/pkg/services/sqlstore"

	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/registry"
	"github.com/openinsight-project/grafinsight/pkg/setting"
)

var metricsLogger log.Logger = log.New("metrics")

func init() {
	registry.RegisterService(&UsageStatsService{})
}

type UsageStats interface {
	GetUsageReport(ctx context.Context) (UsageReport, error)

	RegisterMetric(name string, fn MetricFunc)
}

type MetricFunc func() (interface{}, error)

type UsageStatsService struct {
	Cfg                *setting.Cfg               `inject:""`
	Bus                bus.Bus                    `inject:""`
	SQLStore           *sqlstore.SQLStore         `inject:""`
	AlertingUsageStats alerting.UsageStatsQuerier `inject:""`
	License            models.Licensing           `inject:""`

	log log.Logger

	oauthProviders           map[string]bool
	externalMetrics          map[string]MetricFunc
	concurrentUserStatsCache memoConcurrentUserStats
}

func (uss *UsageStatsService) Init() error {
	uss.log = log.New("infra.usagestats")
	uss.oauthProviders = social.GetOAuthProviders(uss.Cfg)
	uss.externalMetrics = make(map[string]MetricFunc)
	return nil
}

func (uss *UsageStatsService) Run(ctx context.Context) error {
	uss.updateTotalStats()

	sendReportTicker := time.NewTicker(time.Hour * 24)
	updateStatsTicker := time.NewTicker(time.Minute * 30)
	defer sendReportTicker.Stop()
	defer updateStatsTicker.Stop()

	for {
		select {
		case <-sendReportTicker.C:
			if err := uss.sendUsageStats(ctx); err != nil {
				metricsLogger.Warn("Failed to send usage stats", "err", err)
			}
		case <-updateStatsTicker.C:
			uss.updateTotalStats()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type memoConcurrentUserStats struct {
	stats *concurrentUsersStats

	memoized time.Time
}

const concurrentUserStatsCacheLifetime = time.Hour

func (uss *UsageStatsService) GetConcurrentUsersStats(ctx context.Context) (*concurrentUsersStats, error) {
	memoizationPeriod := time.Now().Add(-concurrentUserStatsCacheLifetime)
	if !uss.concurrentUserStatsCache.memoized.Before(memoizationPeriod) {
		return uss.concurrentUserStatsCache.stats, nil
	}

	uss.concurrentUserStatsCache.stats = &concurrentUsersStats{}
	err := uss.SQLStore.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		// Retrieves concurrent users stats as a histogram. Buckets are accumulative and upper bound is inclusive.
		rawSQL := `
SELECT
    COUNT(CASE WHEN tokens <= 3 THEN 1 END) AS bucket_le_3,
    COUNT(CASE WHEN tokens <= 6 THEN 1 END) AS bucket_le_6,
    COUNT(CASE WHEN tokens <= 9 THEN 1 END) AS bucket_le_9,
    COUNT(CASE WHEN tokens <= 12 THEN 1 END) AS bucket_le_12,
    COUNT(CASE WHEN tokens <= 15 THEN 1 END) AS bucket_le_15,
    COUNT(1) AS bucket_le_inf
FROM (select count(1) as tokens from user_auth_token group by user_id) uat;`
		_, err := sess.SQL(rawSQL).Get(uss.concurrentUserStatsCache.stats)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get concurrent users stats from database: %w", err)
	}

	uss.concurrentUserStatsCache.memoized = time.Now()
	return uss.concurrentUserStatsCache.stats, nil
}
