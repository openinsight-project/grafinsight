// +build integration

package sqlstore

import (
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestIntegration_GetAdminStats(t *testing.T) {
	InitTestDB(t)

	query := models.GetAdminStatsQuery{}
	err := GetAdminStats(&query)
	require.NoError(t, err)
}
