// +build integration

package sqlstore

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/openinsight-project/grafinsight/pkg/models"
)

func TestSavingTags(t *testing.T) {
	Convey("Testing tags saving", t, func() {
		InitTestDB(t)

		tagPairs := []*models.Tag{
			{Key: "outage"},
			{Key: "type", Value: "outage"},
			{Key: "server", Value: "server-1"},
			{Key: "error"},
		}
		tags, err := EnsureTagsExist(newSession(), tagPairs)

		So(err, ShouldBeNil)
		So(len(tags), ShouldEqual, 4)
	})
}
