package elasticsearch

import (
	"context"
	"fmt"

	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/tsdb"
	es "github.com/openinsight-project/grafinsight/pkg/tsdb/elasticsearch/client"
)

// ElasticsearchExecutor represents a handler for handling elasticsearch datasource request
type ElasticsearchExecutor struct{}

var (
	intervalCalculator tsdb.IntervalCalculator
)

// NewElasticsearchExecutor creates a new elasticsearch executor
func NewElasticsearchExecutor(dsInfo *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	return &ElasticsearchExecutor{}, nil
}

func init() {
	intervalCalculator = tsdb.NewIntervalCalculator(nil)
	tsdb.RegisterTsdbQueryEndpoint("elasticsearch", NewElasticsearchExecutor)
}

// Query handles an elasticsearch datasource request
func (e *ElasticsearchExecutor) Query(ctx context.Context, dsInfo *models.DataSource, tsdbQuery *tsdb.TsdbQuery) (*tsdb.Response, error) {
	if len(tsdbQuery.Queries) == 0 {
		return nil, fmt.Errorf("query contains no queries")
	}

	client, err := es.NewClient(ctx, dsInfo, tsdbQuery.TimeRange)
	if err != nil {
		return nil, err
	}

	if tsdbQuery.Debug {
		client.EnableDebug()
	}

	query := newTimeSeriesQuery(client, tsdbQuery, intervalCalculator)
	return query.execute()
}
