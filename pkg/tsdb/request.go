package tsdb

import (
	"context"
	"fmt"

	"github.com/openinsight-project/grafinsight/pkg/models"
)

type HandleRequestFunc func(ctx context.Context, dsInfo *models.DataSource, req *TsdbQuery) (*Response, error)

func HandleRequest(ctx context.Context, dsInfo *models.DataSource, req *TsdbQuery) (*Response, error) {
	var endpoint TsdbQueryEndpoint
	fn, exists := registry[dsInfo.Type]
	if !exists {
		return nil, fmt.Errorf("could not find executor for data source type: %s", dsInfo.Type)
	}

	var err error
	endpoint, err = fn(dsInfo)
	if err != nil {
		return nil, err
	}

	return endpoint.Query(ctx, dsInfo, req)
}
