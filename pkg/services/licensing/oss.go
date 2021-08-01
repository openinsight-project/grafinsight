package licensing

import (
	"github.com/openinsight-project/grafinsight/pkg/api/dtos"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/services/hooks"
	"github.com/openinsight-project/grafinsight/pkg/setting"
)

const (
	openSource = "Open Source"
)

type OSSLicensingService struct {
	Cfg          *setting.Cfg        `inject:""`
	HooksService *hooks.HooksService `inject:""`
}

func (*OSSLicensingService) HasLicense() bool {
	return false
}

func (*OSSLicensingService) Expiry() int64 {
	return 0
}

func (*OSSLicensingService) Edition() string {
	return openSource
}

func (*OSSLicensingService) StateInfo() string {
	return ""
}

func (*OSSLicensingService) ContentDeliveryPrefix() string {
	return "grafana-oss"
}

func (l *OSSLicensingService) LicenseURL(user *models.SignedInUser) string {
	if user.IsGrafanaAdmin {
		return l.Cfg.AppSubURL + "/admin/upgrading"
	}

	return "https://grafana.com/products/enterprise/?utm_source=grafana_footer"
}

func (l *OSSLicensingService) Init() error {
	l.HooksService.AddIndexDataHook(func(indexData *dtos.IndexViewData, req *models.ReqContext) {})

	return nil
}

func (*OSSLicensingService) HasValidLicense() bool {
	return false
}
