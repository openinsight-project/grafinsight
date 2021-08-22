package models

import (
	"errors"
	"time"

	"github.com/openinsight-project/grafinsight/pkg/components/securejsondata"
	"github.com/openinsight-project/grafinsight/pkg/components/simplejson"
)

const (
	DS_GRAPHITE      = "graphite"
	DS_INFLUXDB      = "influxdb"
	DS_INFLUXDB_08   = "influxdb_08"
	DS_ES            = "elasticsearch"
	DS_OPENTSDB      = "opentsdb"
	DS_CLOUDWATCH    = "cloudwatch"
	DS_KAIROSDB      = "kairosdb"
	DS_PROMETHEUS    = "prometheus"
	DS_POSTGRES      = "postgres"
	DS_MYSQL         = "mysql"
	DS_MSSQL         = "mssql"
	DS_ACCESS_DIRECT = "direct"
	DS_ACCESS_PROXY  = "proxy"
	// Stackdriver was renamed Google Cloud monitoring 2020-05 but we keep
	// "stackdriver" to avoid breaking changes in reporting.
	DS_CLOUD_MONITORING = "stackdriver"
	DS_AZURE_MONITOR    = "grafinsight-azure-monitor-datasource"
	DS_LOKI             = "loki"
	DS_ES_OPEN_DISTRO   = "grafinsight-es-open-distro-datasource"
)

var (
	ErrDataSourceNotFound                = errors.New("data source not found")
	ErrDataSourceNameExists              = errors.New("data source with the same name already exists")
	ErrDataSourceUidExists               = errors.New("data source with the same uid already exists")
	ErrDataSourceUpdatingOldVersion      = errors.New("trying to update old version of datasource")
	ErrDatasourceIsReadOnly              = errors.New("data source is readonly, can only be updated from configuration")
	ErrDataSourceAccessDenied            = errors.New("data source access denied")
	ErrDataSourceFailedGenerateUniqueUid = errors.New("failed to generate unique datasource ID")
	ErrDataSourceIdentifierNotSet        = errors.New("unique identifier and org id are needed to be able to get or delete a datasource")
)

type DsAccess string

type DataSource struct {
	Id      int64 `json:"id"`
	OrgId   int64 `json:"orgId"`
	Version int   `json:"version"`

	Name              string                        `json:"name"`
	Type              string                        `json:"type"`
	Access            DsAccess                      `json:"access"`
	Url               string                        `json:"url"`
	Password          string                        `json:"password"`
	User              string                        `json:"user"`
	Database          string                        `json:"database"`
	BasicAuth         bool                          `json:"basicAuth"`
	BasicAuthUser     string                        `json:"basicAuthUser"`
	BasicAuthPassword string                        `json:"basicAuthPassword"`
	WithCredentials   bool                          `json:"withCredentials"`
	IsDefault         bool                          `json:"isDefault"`
	JsonData          *simplejson.Json              `json:"jsonData"`
	SecureJsonData    securejsondata.SecureJsonData `json:"secureJsonData"`
	ReadOnly          bool                          `json:"readOnly"`
	Uid               string                        `json:"uid"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// DecryptedBasicAuthPassword returns data source basic auth password in plain text. It uses either deprecated
// basic_auth_password field or encrypted secure_json_data[basicAuthPassword] variable.
func (ds *DataSource) DecryptedBasicAuthPassword() string {
	return ds.decryptedValue("basicAuthPassword", ds.BasicAuthPassword)
}

// DecryptedPassword returns data source password in plain text. It uses either deprecated password field
// or encrypted secure_json_data[password] variable.
func (ds *DataSource) DecryptedPassword() string {
	return ds.decryptedValue("password", ds.Password)
}

// decryptedValue returns decrypted value from secureJsonData
func (ds *DataSource) decryptedValue(field string, fallback string) string {
	if value, ok := ds.DecryptedValue(field); ok {
		return value
	}
	return fallback
}

var knownDatasourcePlugins = map[string]bool{
	DS_ES:                                  true,
	DS_GRAPHITE:                            true,
	DS_INFLUXDB:                            true,
	DS_INFLUXDB_08:                         true,
	DS_KAIROSDB:                            true,
	DS_CLOUDWATCH:                          true,
	DS_PROMETHEUS:                          true,
	DS_OPENTSDB:                            true,
	DS_POSTGRES:                            true,
	DS_MYSQL:                               true,
	DS_MSSQL:                               true,
	DS_CLOUD_MONITORING:                    true,
	DS_AZURE_MONITOR:                       true,
	DS_LOKI:                                true,
	"opennms":                              true,
	"abhisant-druid-datasource":            true,
	"dalmatinerdb-datasource":              true,
	"gnocci":                               true,
	"zabbix":                               true,
	"newrelic-app":                         true,
	"grafinsight-datadog-datasource":       true,
	"grafinsight-simple-json":              true,
	"grafinsight-splunk-datasource":        true,
	"udoprog-heroic-datasource":            true,
	"grafinsight-openfalcon-datasource":    true,
	"opennms-datasource":                   true,
	"rackerlabs-blueflood-datasource":      true,
	"crate-datasource":                     true,
	"ayoungprogrammer-finance-datasource":  true,
	"monasca-datasource":                   true,
	"vertamedia-clickhouse-datasource":     true,
	"alexanderzobnin-zabbix-datasource":    true,
	"grafinsight-influxdb-flux-datasource": true,
	"doitintl-bigquery-datasource":         true,
	"grafinsight-azure-data-explorer-datasource": true,
	"tempo": true,
}

func IsKnownDataSourcePlugin(dsType string) bool {
	_, exists := knownDatasourcePlugins[dsType]
	return exists
}

// ----------------------
// COMMANDS

// Also acts as api DTO
type AddDataSourceCommand struct {
	Name              string            `json:"name" binding:"Required"`
	Type              string            `json:"type" binding:"Required"`
	Access            DsAccess          `json:"access" binding:"Required"`
	Url               string            `json:"url"`
	Password          string            `json:"password"`
	Database          string            `json:"database"`
	User              string            `json:"user"`
	BasicAuth         bool              `json:"basicAuth"`
	BasicAuthUser     string            `json:"basicAuthUser"`
	BasicAuthPassword string            `json:"basicAuthPassword"`
	WithCredentials   bool              `json:"withCredentials"`
	IsDefault         bool              `json:"isDefault"`
	JsonData          *simplejson.Json  `json:"jsonData"`
	SecureJsonData    map[string]string `json:"secureJsonData"`
	Uid               string            `json:"uid"`

	OrgId    int64 `json:"-"`
	ReadOnly bool  `json:"-"`

	Result *DataSource
}

// Also acts as api DTO
type UpdateDataSourceCommand struct {
	Name              string            `json:"name" binding:"Required"`
	Type              string            `json:"type" binding:"Required"`
	Access            DsAccess          `json:"access" binding:"Required"`
	Url               string            `json:"url"`
	Password          string            `json:"password"`
	User              string            `json:"user"`
	Database          string            `json:"database"`
	BasicAuth         bool              `json:"basicAuth"`
	BasicAuthUser     string            `json:"basicAuthUser"`
	BasicAuthPassword string            `json:"basicAuthPassword"`
	WithCredentials   bool              `json:"withCredentials"`
	IsDefault         bool              `json:"isDefault"`
	JsonData          *simplejson.Json  `json:"jsonData"`
	SecureJsonData    map[string]string `json:"secureJsonData"`
	Version           int               `json:"version"`
	Uid               string            `json:"uid"`

	OrgId    int64 `json:"-"`
	Id       int64 `json:"-"`
	ReadOnly bool  `json:"-"`

	Result *DataSource
}

// DeleteDataSourceCommand will delete a DataSource based on OrgID as well as the UID (preferred), ID, or Name.
// At least one of the UID, ID, or Name properties must be set in addition to OrgID.
type DeleteDataSourceCommand struct {
	ID   int64
	UID  string
	Name string

	OrgID int64

	DeletedDatasourcesCount int64
}

// ---------------------
// QUERIES

type GetDataSourcesQuery struct {
	OrgId           int64
	DataSourceLimit int
	User            *SignedInUser
	Result          []*DataSource
}

type GetDataSourcesByTypeQuery struct {
	Type   string
	Result []*DataSource
}

type GetDefaultDataSourceQuery struct {
	OrgId  int64
	User   *SignedInUser
	Result *DataSource
}

// GetDataSourceQuery will get a DataSource based on OrgID as well as the UID (preferred), ID, or Name.
// At least one of the UID, ID, or Name properties must be set in addition to OrgID.
type GetDataSourceQuery struct {
	Id   int64
	Uid  string
	Name string

	OrgId int64

	Result *DataSource
}

// ---------------------
//  Permissions
// ---------------------

type DsPermissionType int

const (
	DsPermissionNoAccess DsPermissionType = iota
	DsPermissionQuery
)

func (p DsPermissionType) String() string {
	names := map[int]string{
		int(DsPermissionQuery):    "Query",
		int(DsPermissionNoAccess): "No Access",
	}
	return names[int(p)]
}

type DatasourcesPermissionFilterQuery struct {
	User        *SignedInUser
	Datasources []*DataSource
	Result      []*DataSource
}
