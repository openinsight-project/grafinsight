package dtos

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	"github.com/openinsight-project/grafinsight/pkg/components/simplejson"
	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/setting"
)

var regNonAlphaNumeric = regexp.MustCompile("[^a-zA-Z0-9]+")

type AnyId struct {
	Id int64 `json:"id"`
}

type LoginCommand struct {
	User     string `json:"user" binding:"Required"`
	Password string `json:"password" binding:"Required"`
	Remember bool   `json:"remember"`
}

type CurrentUser struct {
	IsSignedIn                 bool              `json:"isSignedIn"`
	Id                         int64             `json:"id"`
	Login                      string            `json:"login"`
	Email                      string            `json:"email"`
	Name                       string            `json:"name"`
	LightTheme                 bool              `json:"lightTheme"`
	OrgCount                   int               `json:"orgCount"`
	OrgId                      int64             `json:"orgId"`
	OrgName                    string            `json:"orgName"`
	OrgRole                    models.RoleType   `json:"orgRole"`
	IsGrafanaAdmin             bool              `json:"isGrafanaAdmin"`
	GravatarUrl                string            `json:"gravatarUrl"`
	Timezone                   string            `json:"timezone"`
	Locale                     string            `json:"locale"`
	HelpFlags1                 models.HelpFlags1 `json:"helpFlags1"`
	HasEditPermissionInFolders bool              `json:"hasEditPermissionInFolders"`
}

type MetricRequest struct {
	From    string             `json:"from"`
	To      string             `json:"to"`
	Queries []*simplejson.Json `json:"queries"`
	Debug   bool               `json:"debug"`
}

func GetGravatarUrl(text string) string {
	if setting.DisableGravatar {
		return setting.AppSubUrl + "/public/img/user_profile.png"
	}

	if text == "" {
		return ""
	}

	hasher := md5.New()
	if _, err := hasher.Write([]byte(strings.ToLower(text))); err != nil {
		log.Warnf("Failed to hash text: %s", err)
	}
	return fmt.Sprintf(setting.AppSubUrl+"/avatar/%x", hasher.Sum(nil))
}

func GetGravatarUrlWithDefault(text string, defaultText string) string {
	if text != "" {
		return GetGravatarUrl(text)
	}

	text = regNonAlphaNumeric.ReplaceAllString(defaultText, "") + "@localhost"

	return GetGravatarUrl(text)
}

func IsHiddenUser(userLogin string, signedInUser *models.SignedInUser, cfg *setting.Cfg) bool {
	if userLogin == "" || signedInUser.IsGrafanaAdmin || userLogin == signedInUser.Login {
		return false
	}

	if _, hidden := cfg.HiddenUsers[userLogin]; hidden {
		return true
	}

	return false
}
