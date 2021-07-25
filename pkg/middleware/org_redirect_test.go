package middleware

import (
	"context"
	"fmt"
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestOrgRedirectMiddleware(t *testing.T) {
	middlewareScenario(t, "when setting a correct org for the user", func(t *testing.T, sc *scenarioContext) {
		sc.withTokenSessionCookie("token")
		bus.AddHandler("test", func(query *models.SetUsingOrgCommand) error {
			return nil
		})

		bus.AddHandler("test", func(query *models.GetSignedInUserQuery) error {
			query.Result = &models.SignedInUser{OrgId: 1, UserId: 12}
			return nil
		})

		sc.userAuthTokenService.LookupTokenProvider = func(ctx context.Context, unhashedToken string) (*models.UserToken, error) {
			return &models.UserToken{
				UserId:        0,
				UnhashedToken: "",
			}, nil
		}

		sc.m.Get("/", sc.defaultHandler)
		sc.fakeReq("GET", "/?orgId=3").exec()

		assert.Equal(t, 302, sc.resp.Code)
	})

	middlewareScenario(t, "when setting an invalid org for user", func(t *testing.T, sc *scenarioContext) {
		sc.withTokenSessionCookie("token")
		bus.AddHandler("test", func(query *models.SetUsingOrgCommand) error {
			return fmt.Errorf("")
		})

		bus.AddHandler("test", func(query *models.GetSignedInUserQuery) error {
			query.Result = &models.SignedInUser{OrgId: 1, UserId: 12}
			return nil
		})

		sc.userAuthTokenService.LookupTokenProvider = func(ctx context.Context, unhashedToken string) (*models.UserToken, error) {
			return &models.UserToken{
				UserId:        12,
				UnhashedToken: "",
			}, nil
		}

		sc.m.Get("/", sc.defaultHandler)
		sc.fakeReq("GET", "/?orgId=3").exec()

		assert.Equal(t, 404, sc.resp.Code)
	})
}
