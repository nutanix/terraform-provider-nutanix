package objectstoresv2

import (
	"testing"

	objectsCommon "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/common/v1/config"
	objectsResponse "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/common/v1/response"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/objects/v4/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitObjectStoreCertificatesDataSourceSchema(t *testing.T) {
	t.Parallel()

	r := DatasourceNutanixObjectStoreCertificatesV2()
	require.NoError(t, r.InternalValidate(r.Schema, false))
}

func TestUnitFlattenCertificates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.Certificate
		want interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []config.Certificate{}, want: nil},
		{
			name: "populated",
			in: []config.Certificate{{
				TenantId: utils.StringPtr("tenant-1"),
				ExtId:    utils.StringPtr("cert-1"),
				Links: []objectsResponse.ApiLink{{
					Href: utils.StringPtr("https://example.test/cert/1"),
					Rel:  utils.StringPtr("self"),
				}},
				Metadata: &objectsCommon.Metadata{
					OwnerReferenceId: utils.StringPtr("owner-1"),
				},
				AlternateFqdns: []objectsCommon.FQDN{{
					Value: utils.StringPtr("alt.example.test"),
				}},
				AlternateIps: []objectsCommon.IPAddress{{
					Ipv4: &objectsCommon.IPv4Address{
						Value:        utils.StringPtr("10.0.0.30"),
						PrefixLength: utils.IntPtr(32),
					},
				}},
			}},
			want: []map[string]interface{}{{
				"tenant_id": utils.StringPtr("tenant-1"),
				"ext_id":    utils.StringPtr("cert-1"),
				"links": []map[string]interface{}{{
					"href": utils.StringPtr("https://example.test/cert/1"),
					"rel":  utils.StringPtr("self"),
				}},
				"metadata": []map[string]interface{}{{
					"owner_reference_id":   utils.StringPtr("owner-1"),
					"owner_user_name":      (*string)(nil),
					"project_reference_id": (*string)(nil),
					"project_name":         (*string)(nil),
					"category_ids":         []string(nil),
				}},
				"alternate_fqdns": []map[string]interface{}{{
					"value": "alt.example.test",
				}},
				"alternate_ips": []map[string]interface{}{{
					"ipv4": []map[string]interface{}{{
						"value":         utils.StringPtr("10.0.0.30"),
						"prefix_length": utils.IntPtr(32),
					}},
					"ipv6": []map[string]interface{}(nil),
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCertificates(tt.in))
		})
	}
}
