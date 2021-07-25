package cloudwatch

import (
	"context"
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/tsdb"
	"github.com/stretchr/testify/assert"
)

func TestTimeSeriesQuery(t *testing.T) {
	executor := newExecutor(nil, newTestConfig(), fakeSessionCache{})

	t.Run("End time before start time should result in error", func(t *testing.T) {
		_, err := executor.executeTimeSeriesQuery(context.TODO(), &tsdb.TsdbQuery{TimeRange: tsdb.NewTimeRange("now-1h", "now-2h")})
		assert.EqualError(t, err, "invalid time range: start time must be before end time")
	})

	t.Run("End time equals start time should result in error", func(t *testing.T) {
		_, err := executor.executeTimeSeriesQuery(context.TODO(), &tsdb.TsdbQuery{TimeRange: tsdb.NewTimeRange("now-1h", "now-1h")})
		assert.EqualError(t, err, "invalid time range: start time must be before end time")
	})
}
