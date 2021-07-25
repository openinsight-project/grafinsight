package cloudwatch

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
	"github.com/openinsight-project/grafinsight/pkg/components/simplejson"
	"github.com/openinsight-project/grafinsight/pkg/setting"
	"github.com/openinsight-project/grafinsight/pkg/tsdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery_Metrics(t *testing.T) {
	origNewCWClient := NewCWClient
	t.Cleanup(func() {
		NewCWClient = origNewCWClient
	})

	var client FakeCWClient

	NewCWClient = func(sess *session.Session) cloudwatchiface.CloudWatchAPI {
		return client
	}

	t.Run("Custom metrics", func(t *testing.T) {
		client = FakeCWClient{
			Metrics: []*cloudwatch.Metric{
				{
					MetricName: aws.String("Test_MetricName"),
					Dimensions: []*cloudwatch.Dimension{
						{
							Name: aws.String("Test_DimensionName"),
						},
					},
				},
			},
		}
		executor := newExecutor(nil, newTestConfig(), fakeSessionCache{})
		resp, err := executor.Query(context.Background(), fakeDataSource(), &tsdb.TsdbQuery{
			Queries: []*tsdb.Query{
				{
					Model: simplejson.NewFromAny(map[string]interface{}{
						"type":      "metricFindQuery",
						"subtype":   "metrics",
						"region":    "us-east-1",
						"namespace": "custom",
					}),
				},
			},
		})
		require.NoError(t, err)

		assert.Equal(t, &tsdb.Response{
			Results: map[string]*tsdb.QueryResult{
				"": {
					Meta: simplejson.NewFromAny(map[string]interface{}{
						"rowCount": 1,
					}),
					Tables: []*tsdb.Table{
						{
							Columns: []tsdb.TableColumn{
								{
									Text: "text",
								},
								{
									Text: "value",
								},
							},
							Rows: []tsdb.RowValues{
								{
									"Test_MetricName",
									"Test_MetricName",
								},
							},
						},
					},
				},
			},
		}, resp)
	})

	t.Run("Dimension keys for custom metrics", func(t *testing.T) {
		client = FakeCWClient{
			Metrics: []*cloudwatch.Metric{
				{
					MetricName: aws.String("Test_MetricName"),
					Dimensions: []*cloudwatch.Dimension{
						{
							Name: aws.String("Test_DimensionName"),
						},
					},
				},
			},
		}
		executor := newExecutor(nil, newTestConfig(), fakeSessionCache{})
		resp, err := executor.Query(context.Background(), fakeDataSource(), &tsdb.TsdbQuery{
			Queries: []*tsdb.Query{
				{
					Model: simplejson.NewFromAny(map[string]interface{}{
						"type":      "metricFindQuery",
						"subtype":   "dimension_keys",
						"region":    "us-east-1",
						"namespace": "custom",
					}),
				},
			},
		})
		require.NoError(t, err)

		assert.Equal(t, &tsdb.Response{
			Results: map[string]*tsdb.QueryResult{
				"": {
					Meta: simplejson.NewFromAny(map[string]interface{}{
						"rowCount": 1,
					}),
					Tables: []*tsdb.Table{
						{
							Columns: []tsdb.TableColumn{
								{
									Text: "text",
								},
								{
									Text: "value",
								},
							},
							Rows: []tsdb.RowValues{
								{
									"Test_DimensionName",
									"Test_DimensionName",
								},
							},
						},
					},
				},
			},
		}, resp)
	})
}

func TestQuery_Regions(t *testing.T) {
	origNewEC2Client := newEC2Client
	t.Cleanup(func() {
		newEC2Client = origNewEC2Client
	})

	var cli fakeEC2Client

	newEC2Client = func(client.ConfigProvider) ec2iface.EC2API {
		return cli
	}

	t.Run("An extra region", func(t *testing.T) {
		const regionName = "xtra-region"
		cli = fakeEC2Client{
			regions: []string{regionName},
		}
		executor := newExecutor(nil, newTestConfig(), fakeSessionCache{})
		resp, err := executor.Query(context.Background(), fakeDataSource(), &tsdb.TsdbQuery{
			Queries: []*tsdb.Query{
				{
					Model: simplejson.NewFromAny(map[string]interface{}{
						"type":      "metricFindQuery",
						"subtype":   "regions",
						"region":    "us-east-1",
						"namespace": "custom",
					}),
				},
			},
		})
		require.NoError(t, err)

		rows := []tsdb.RowValues{}
		for _, region := range knownRegions {
			rows = append(rows, []interface{}{
				region,
				region,
			})
		}
		rows = append(rows, []interface{}{
			regionName,
			regionName,
		})
		assert.Equal(t, &tsdb.Response{
			Results: map[string]*tsdb.QueryResult{
				"": {
					Meta: simplejson.NewFromAny(map[string]interface{}{
						"rowCount": len(knownRegions) + 1,
					}),
					Tables: []*tsdb.Table{
						{
							Columns: []tsdb.TableColumn{
								{
									Text: "text",
								},
								{
									Text: "value",
								},
							},
							Rows: rows,
						},
					},
				},
			},
		}, resp)
	})
}

func TestQuery_InstanceAttributes(t *testing.T) {
	origNewEC2Client := newEC2Client
	t.Cleanup(func() {
		newEC2Client = origNewEC2Client
	})

	var cli fakeEC2Client

	newEC2Client = func(client.ConfigProvider) ec2iface.EC2API {
		return cli
	}

	t.Run("Get instance ID", func(t *testing.T) {
		const instanceID = "i-12345678"
		cli = fakeEC2Client{
			reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							InstanceId: aws.String(instanceID),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Environment"),
									Value: aws.String("production"),
								},
							},
						},
					},
				},
			},
		}
		executor := newExecutor(nil, newTestConfig(), fakeSessionCache{})
		resp, err := executor.Query(context.Background(), fakeDataSource(), &tsdb.TsdbQuery{
			Queries: []*tsdb.Query{
				{
					Model: simplejson.NewFromAny(map[string]interface{}{
						"type":          "metricFindQuery",
						"subtype":       "ec2_instance_attribute",
						"region":        "us-east-1",
						"attributeName": "InstanceId",
						"filters": map[string]interface{}{
							"tag:Environment": []string{"production"},
						},
					}),
				},
			},
		})
		require.NoError(t, err)

		assert.Equal(t, &tsdb.Response{
			Results: map[string]*tsdb.QueryResult{
				"": {
					Meta: simplejson.NewFromAny(map[string]interface{}{
						"rowCount": 1,
					}),
					Tables: []*tsdb.Table{
						{
							Columns: []tsdb.TableColumn{
								{
									Text: "text",
								},
								{
									Text: "value",
								},
							},
							Rows: []tsdb.RowValues{
								{
									instanceID,
									instanceID,
								},
							},
						},
					},
				},
			},
		}, resp)
	})
}

func TestQuery_EBSVolumeIDs(t *testing.T) {
	origNewEC2Client := newEC2Client
	t.Cleanup(func() {
		newEC2Client = origNewEC2Client
	})

	var cli fakeEC2Client

	newEC2Client = func(client.ConfigProvider) ec2iface.EC2API {
		return cli
	}

	t.Run("", func(t *testing.T) {
		const instanceIDs = "{i-1, i-2, i-3}"

		cli = fakeEC2Client{
			reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							InstanceId: aws.String("i-1"),
							BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
								{Ebs: &ec2.EbsInstanceBlockDevice{VolumeId: aws.String("vol-1-1")}},
								{Ebs: &ec2.EbsInstanceBlockDevice{VolumeId: aws.String("vol-1-2")}},
							},
						},
						{
							InstanceId: aws.String("i-2"),
							BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
								{Ebs: &ec2.EbsInstanceBlockDevice{VolumeId: aws.String("vol-2-1")}},
								{Ebs: &ec2.EbsInstanceBlockDevice{VolumeId: aws.String("vol-2-2")}},
							},
						},
					},
				},
				{
					Instances: []*ec2.Instance{
						{
							InstanceId: aws.String("i-3"),
							BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
								{Ebs: &ec2.EbsInstanceBlockDevice{VolumeId: aws.String("vol-3-1")}},
								{Ebs: &ec2.EbsInstanceBlockDevice{VolumeId: aws.String("vol-3-2")}},
							},
						},
						{
							InstanceId: aws.String("i-4"),
							BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
								{Ebs: &ec2.EbsInstanceBlockDevice{VolumeId: aws.String("vol-4-1")}},
								{Ebs: &ec2.EbsInstanceBlockDevice{VolumeId: aws.String("vol-4-2")}},
							},
						},
					},
				},
			},
		}
		executor := newExecutor(nil, newTestConfig(), fakeSessionCache{})
		resp, err := executor.Query(context.Background(), fakeDataSource(), &tsdb.TsdbQuery{
			Queries: []*tsdb.Query{
				{
					Model: simplejson.NewFromAny(map[string]interface{}{
						"type":       "metricFindQuery",
						"subtype":    "ebs_volume_ids",
						"region":     "us-east-1",
						"instanceId": instanceIDs,
					}),
				},
			},
		})
		require.NoError(t, err)

		assert.Equal(t, &tsdb.Response{
			Results: map[string]*tsdb.QueryResult{
				"": {
					Meta: simplejson.NewFromAny(map[string]interface{}{
						"rowCount": 6,
					}),
					Tables: []*tsdb.Table{
						{
							Columns: []tsdb.TableColumn{
								{
									Text: "text",
								},
								{
									Text: "value",
								},
							},
							Rows: []tsdb.RowValues{
								{
									"vol-1-1",
									"vol-1-1",
								},
								{
									"vol-1-2",
									"vol-1-2",
								},
								{
									"vol-2-1",
									"vol-2-1",
								},
								{
									"vol-2-2",
									"vol-2-2",
								},
								{
									"vol-3-1",
									"vol-3-1",
								},
								{
									"vol-3-2",
									"vol-3-2",
								},
							},
						},
					},
				},
			},
		}, resp)
	})
}

func TestQuery_ResourceARNs(t *testing.T) {
	origNewRGTAClient := newRGTAClient
	t.Cleanup(func() {
		newRGTAClient = origNewRGTAClient
	})

	var cli fakeRGTAClient

	newRGTAClient = func(client.ConfigProvider) resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI {
		return cli
	}

	t.Run("", func(t *testing.T) {
		cli = fakeRGTAClient{
			tagMapping: []*resourcegroupstaggingapi.ResourceTagMapping{
				{
					ResourceARN: aws.String("arn:aws:ec2:us-east-1:123456789012:instance/i-12345678901234567"),
					Tags: []*resourcegroupstaggingapi.Tag{
						{
							Key:   aws.String("Environment"),
							Value: aws.String("production"),
						},
					},
				},
				{
					ResourceARN: aws.String("arn:aws:ec2:us-east-1:123456789012:instance/i-76543210987654321"),
					Tags: []*resourcegroupstaggingapi.Tag{
						{
							Key:   aws.String("Environment"),
							Value: aws.String("production"),
						},
					},
				},
			},
		}
		executor := newExecutor(nil, newTestConfig(), fakeSessionCache{})
		resp, err := executor.Query(context.Background(), fakeDataSource(), &tsdb.TsdbQuery{
			Queries: []*tsdb.Query{
				{
					Model: simplejson.NewFromAny(map[string]interface{}{
						"type":         "metricFindQuery",
						"subtype":      "resource_arns",
						"region":       "us-east-1",
						"resourceType": "ec2:instance",
						"tags": map[string]interface{}{
							"Environment": []string{"production"},
						},
					}),
				},
			},
		})
		require.NoError(t, err)

		assert.Equal(t, &tsdb.Response{
			Results: map[string]*tsdb.QueryResult{
				"": {
					Meta: simplejson.NewFromAny(map[string]interface{}{
						"rowCount": 2,
					}),
					Tables: []*tsdb.Table{
						{
							Columns: []tsdb.TableColumn{
								{
									Text: "text",
								},
								{
									Text: "value",
								},
							},
							Rows: []tsdb.RowValues{
								{
									"arn:aws:ec2:us-east-1:123456789012:instance/i-12345678901234567",
									"arn:aws:ec2:us-east-1:123456789012:instance/i-12345678901234567",
								},
								{
									"arn:aws:ec2:us-east-1:123456789012:instance/i-76543210987654321",
									"arn:aws:ec2:us-east-1:123456789012:instance/i-76543210987654321",
								},
							},
						},
					},
				},
			},
		}, resp)
	})
}

func TestQuery_ListMetricsPagination(t *testing.T) {
	origNewCWClient := NewCWClient
	t.Cleanup(func() {
		NewCWClient = origNewCWClient
	})

	var client FakeCWClient

	NewCWClient = func(sess *session.Session) cloudwatchiface.CloudWatchAPI {
		return client
	}

	metrics := []*cloudwatch.Metric{
		{MetricName: aws.String("Test_MetricName1")},
		{MetricName: aws.String("Test_MetricName2")},
		{MetricName: aws.String("Test_MetricName3")},
		{MetricName: aws.String("Test_MetricName4")},
		{MetricName: aws.String("Test_MetricName5")},
		{MetricName: aws.String("Test_MetricName6")},
		{MetricName: aws.String("Test_MetricName7")},
		{MetricName: aws.String("Test_MetricName8")},
		{MetricName: aws.String("Test_MetricName9")},
		{MetricName: aws.String("Test_MetricName10")},
	}

	t.Run("List Metrics and page limit is reached", func(t *testing.T) {
		client = FakeCWClient{Metrics: metrics, MetricsPerPage: 2}
		executor := newExecutor(nil, &setting.Cfg{AWSListMetricsPageLimit: 3, AWSAllowedAuthProviders: []string{"default"}, AWSAssumeRoleEnabled: true}, fakeSessionCache{})
		executor.DataSource = fakeDataSource()
		response, err := executor.listMetrics("default", &cloudwatch.ListMetricsInput{})
		require.NoError(t, err)

		expectedMetrics := client.MetricsPerPage * executor.cfg.AWSListMetricsPageLimit
		assert.Equal(t, expectedMetrics, len(response))
	})

	t.Run("List Metrics and page limit is not reached", func(t *testing.T) {
		client = FakeCWClient{Metrics: metrics, MetricsPerPage: 2}
		executor := newExecutor(nil, &setting.Cfg{AWSListMetricsPageLimit: 1000, AWSAllowedAuthProviders: []string{"default"}, AWSAssumeRoleEnabled: true}, fakeSessionCache{})
		executor.DataSource = fakeDataSource()
		response, err := executor.listMetrics("default", &cloudwatch.ListMetricsInput{})
		require.NoError(t, err)

		assert.Equal(t, len(metrics), len(response))
	})
}
