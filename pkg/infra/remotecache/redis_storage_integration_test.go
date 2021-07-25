// +build redis

package remotecache

import (
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/setting"
)

func TestRedisCacheStorage(t *testing.T) {

	opts := &setting.RemoteCacheOptions{Name: redisCacheType, ConnStr: "addr=localhost:6379"}
	client := createTestClient(t, opts, nil)
	runTestsForClient(t, client)
}
