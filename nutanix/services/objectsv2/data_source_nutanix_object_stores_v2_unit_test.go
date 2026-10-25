package objectstoresv2

import (
	"testing"
	"time"

	objectsCommon "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/config"
	objectsResponse "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/objects/v4/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitObjectStoresDataSourceSchema(t *testing.T) {
	t.Parallel()

	r := DatasourceNutanixObjectStoresV2()
	require.NoError(t, r.InternalValidate(r.Schema, false))
}

func TestUnitFlattenObjectStoreEntities(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	state := config.STATE_OBJECT_STORE_AVAILABLE
	vip := objectsCommon.IPAddress{
		Ipv4: &objectsCommon.IPv4Address{
			Value:        utils.StringPtr("10.0.0.20"),
			PrefixLength: utils.IntPtr(24),
		},
	}
	dnsIP := objectsCommon.IPAddress{
		Ipv4: &objectsCommon.IPv4Address{
			Value:        utils.StringPtr("10.0.0.21"),
			PrefixLength: utils.IntPtr(24),
		},
	}

	tests := []struct {
		name string
		in   []config.ObjectStore
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.ObjectStore{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.ObjectStore{{
				ExtId:    utils.StringPtr("os-1"),
				TenantId: utils.StringPtr("tenant-1"),
				Links: []objectsResponse.ApiLink{{
					Href: utils.StringPtr("https://example.test/os/1"),
					Rel:  utils.StringPtr("self"),
				}},
				Metadata: &objectsCommon.Metadata{
					OwnerReferenceId: utils.StringPtr("owner-1"),
				},
				Name:                    utils.StringPtr("objstore"),
				CreationTime:            &when,
				LastUpdateTime:          &when,
				Description:             utils.StringPtr("desc"),
				DeploymentVersion:       utils.StringPtr("1.0"),
				Domain:                  utils.StringPtr("objects.example.test"),
				Region:                  utils.StringPtr("us"),
				NumWorkerNodes:          utils.Int64Ptr(3),
				ClusterExtId:            utils.StringPtr("cluster-1"),
				StorageNetworkReference: utils.StringPtr("net-1"),
				StorageNetworkVip:       &vip,
				StorageNetworkDnsIp:     &dnsIP,
				PublicNetworkReference:  utils.StringPtr("pub-1"),
				PublicNetworkIps:        []objectsCommon.IPAddress{vip},
				TotalCapacityGiB:        utils.Int64Ptr(100),
				State:                   &state,
				CertificateExtIds:       []string{"cert-1"},
			}},
			want: []map[string]interface{}{{
				"ext_id":    utils.StringPtr("os-1"),
				"tenant_id": utils.StringPtr("tenant-1"),
				"links": []map[string]interface{}{{
					"href": utils.StringPtr("https://example.test/os/1"),
					"rel":  utils.StringPtr("self"),
				}},
				"metadata": []map[string]interface{}{{
					"owner_reference_id":   utils.StringPtr("owner-1"),
					"owner_user_name":      (*string)(nil),
					"project_reference_id": (*string)(nil),
					"project_name":         (*string)(nil),
					"category_ids":         []string(nil),
				}},
				"name":                      "objstore",
				"creation_time":             "2024-01-02T03:04:05Z",
				"last_update_time":          "2024-01-02T03:04:05Z",
				"description":               "desc",
				"deployment_version":        "1.0",
				"domain":                    "objects.example.test",
				"region":                    "us",
				"num_worker_nodes":          int64(3),
				"cluster_ext_id":            "cluster-1",
				"storage_network_reference": "net-1",
				"storage_network_vip": []map[string]interface{}{{
					"ipv4": []map[string]interface{}{{
						"value":         utils.StringPtr("10.0.0.20"),
						"prefix_length": utils.IntPtr(24),
					}},
					"ipv6": []map[string]interface{}(nil),
				}},
				"storage_network_dns_ip": []map[string]interface{}{{
					"ipv4": []map[string]interface{}{{
						"value":         utils.StringPtr("10.0.0.21"),
						"prefix_length": utils.IntPtr(24),
					}},
					"ipv6": []map[string]interface{}(nil),
				}},
				"public_network_reference": "pub-1",
				"public_network_ips": []map[string]interface{}{{
					"ipv4": []map[string]interface{}{{
						"value":         utils.StringPtr("10.0.0.20"),
						"prefix_length": utils.IntPtr(24),
					}},
					"ipv6": []map[string]interface{}(nil),
				}},
				"total_capacity_gib":  int64(100),
				"state":               "OBJECT_STORE_AVAILABLE",
				"certificate_ext_ids": []string{"cert-1"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenObjectStoreEntities(tt.in))
		})
	}
}
