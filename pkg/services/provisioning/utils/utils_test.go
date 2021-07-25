package utils

import (
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/services/sqlstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckOrgExists(t *testing.T) {
	Convey("with default org in database", t, func() {
		sqlstore.InitTestDB(t)

		defaultOrg := models.CreateOrgCommand{Name: "Main Org."}
		err := sqlstore.CreateOrg(&defaultOrg)
		So(err, ShouldBeNil)

		Convey("default org exists", func() {
			err := CheckOrgExists(defaultOrg.Result.Id)
			So(err, ShouldBeNil)
		})

		Convey("other org doesn't exist", func() {
			err := CheckOrgExists(defaultOrg.Result.Id + 1)
			So(err, ShouldEqual, models.ErrOrgNotFound)
		})
	})
}
