package sqlstore

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/gchaincl/sqlhooks"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/services/sqlstore/migrator"
	"github.com/prometheus/client_golang/prometheus"
	"xorm.io/core"
)

var (
	databaseQueryHistogram *prometheus.HistogramVec
)

func init() {
	databaseQueryHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "grafana",
		Name:      "database_queries_duration_seconds",
		Help:      "Database query histogram",
		Buckets:   prometheus.ExponentialBuckets(0.00001, 4, 10),
	}, []string{"status"})

	prometheus.MustRegister(databaseQueryHistogram)
}

// WrapDatabaseDriverWithHooks creates a fake database driver that
// executes pre and post functions which we use to gather metrics about
// database queries. It also registers the metrics.
func WrapDatabaseDriverWithHooks(dbType string) string {
	drivers := map[string]driver.Driver{
		migrator.SQLite:   &sqlite3.SQLiteDriver{},
		migrator.MySQL:    &mysql.MySQLDriver{},
		migrator.Postgres: &pq.Driver{},
	}

	d, exist := drivers[dbType]
	if !exist {
		return dbType
	}

	driverWithHooks := dbType + "WithHooks"
	sql.Register(driverWithHooks, sqlhooks.Wrap(d, &databaseQueryWrapper{log: log.New("sqlstore.metrics")}))
	core.RegisterDriver(driverWithHooks, &databaseQueryWrapperDriver{dbType: dbType})
	return driverWithHooks
}

// databaseQueryWrapper satisfies the sqlhook.databaseQueryWrapper interface
// which allow us to wrap all SQL queries with a `Before` & `After` hook.
type databaseQueryWrapper struct {
	log log.Logger
}

// databaseQueryWrapperKey is used as key to save values in `context.Context`
type databaseQueryWrapperKey struct{}

// Before hook will print the query with its args and return the context with the timestamp
func (h *databaseQueryWrapper) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, databaseQueryWrapperKey{}, time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *databaseQueryWrapper) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value(databaseQueryWrapperKey{}).(time.Time)
	elapsed := time.Since(begin)
	databaseQueryHistogram.WithLabelValues("success").Observe(elapsed.Seconds())
	h.log.Debug("query finished", "status", "success", "elapsed time", elapsed, "sql", query)
	return ctx, nil
}

// OnError will be called if any error happens
func (h *databaseQueryWrapper) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	status := "error"
	// https://golang.org/pkg/database/sql/driver/#ErrSkip
	if err == nil || errors.Is(err, driver.ErrSkip) {
		status = "success"
	}

	begin := ctx.Value(databaseQueryWrapperKey{}).(time.Time)
	elapsed := time.Since(begin)
	databaseQueryHistogram.WithLabelValues(status).Observe(elapsed.Seconds())
	h.log.Debug("query finished", "status", status, "elapsed time", elapsed, "sql", query, "error", err)
	return err
}

// databaseQueryWrapperDriver satisfies the xorm.io/core.Driver interface
type databaseQueryWrapperDriver struct {
	dbType string
}

func (hp *databaseQueryWrapperDriver) Parse(driverName, dataSourceName string) (*core.Uri, error) {
	driver := core.QueryDriver(hp.dbType)
	if driver == nil {
		return nil, fmt.Errorf("could not find driver with name %s", hp.dbType)
	}
	return driver.Parse(driverName, dataSourceName)
}
