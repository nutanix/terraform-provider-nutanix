package passwordmanagerv2

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clusterConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/response"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitPasswordManagerSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixPasswordManagerV2().InternalValidate(ResourceNutanixPasswordManagerV2().Schema, true))
	require.NoError(t, DataSourceNutanixPasswordManagersV2().InternalValidate(DataSourceNutanixPasswordManagersV2().Schema, false))
	require.NoError(t, schemaForLinks().Elem.(*schema.Resource).InternalValidate(schemaForLinks().Elem.(*schema.Resource).Schema, true))
}

func TestUnitFlattenPasswordStatus(t *testing.T) {
	t.Parallel()

	unknown := clusterConfig.PASSWORDSTATUS_UNKNOWN
	secure := clusterConfig.PASSWORDSTATUS_SECURE
	other := clusterConfig.PasswordStatus(99)

	tests := []struct {
		name string
		in   *clusterConfig.PasswordStatus
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "unknown", in: &unknown, want: "UNKNOWN"},
		{name: "redacted", in: statusPtr(clusterConfig.PASSWORDSTATUS_REDACTED), want: "REDACTED"},
		{name: "default", in: statusPtr(clusterConfig.PASSWORDSTATUS_DEFAULT), want: "DEFAULT"},
		{name: "secure", in: &secure, want: "SECURE"},
		{name: "nopasswd", in: statusPtr(clusterConfig.PASSWORDSTATUS_NOPASSWD), want: "NOPASSWD"},
		{name: "multiple issues", in: statusPtr(clusterConfig.PASSWORDSTATUS_MULTIPLE_ISSUES), want: "MULTIPLE_ISSUES"},
		{name: "unrecognized", in: &other, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPasswordStatus(tt.in))
		})
	}
}

func TestUnitFlattenPasswordSystemType(t *testing.T) {
	t.Parallel()

	other := clusterConfig.SystemType(99)
	tests := []struct {
		name string
		in   *clusterConfig.SystemType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "unknown", in: systemTypePtr(clusterConfig.SYSTEMTYPE_UNKNOWN), want: "UNKNOWN"},
		{name: "redacted", in: systemTypePtr(clusterConfig.SYSTEMTYPE_REDACTED), want: "REDACTED"},
		{name: "pc", in: systemTypePtr(clusterConfig.SYSTEMTYPE_PC), want: "PC"},
		{name: "aos", in: systemTypePtr(clusterConfig.SYSTEMTYPE_AOS), want: "AOS"},
		{name: "ahv", in: systemTypePtr(clusterConfig.SYSTEMTYPE_AHV), want: "AHV"},
		{name: "ipmi", in: systemTypePtr(clusterConfig.SYSTEMTYPE_IPMI), want: "IPMI"},
		{name: "unrecognized", in: &other, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPasswordSystemType(tt.in))
		})
	}
}

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.ApiLink{}, want: nil},
		{
			name: "populated",
			in: []import1.ApiLink{{
				Rel:  utils.StringPtr("self"),
				Href: utils.StringPtr("https://example.test/passwords/1"),
			}},
			want: []map[string]interface{}{{
				"rel":  "self",
				"href": "https://example.test/passwords/1",
			}},
		},
		{
			name: "missing fields",
			in:   []import1.ApiLink{{}},
			want: []map[string]interface{}{{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLinks(tt.in))
		})
	}
}

func TestUnitFlattenPasswordEntities(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	status := clusterConfig.PASSWORDSTATUS_SECURE
	systemType := clusterConfig.SYSTEMTYPE_AHV

	tests := []struct {
		name string
		in   []clusterConfig.SystemUserPassword
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []clusterConfig.SystemUserPassword{}, want: []map[string]interface{}{}},
		{
			name: "empty entity",
			in:   []clusterConfig.SystemUserPassword{{}},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: []clusterConfig.SystemUserPassword{{
				ExtId:          utils.StringPtr("ext-1"),
				TenantId:       utils.StringPtr("tenant-1"),
				Username:       utils.StringPtr("admin"),
				ClusterExtId:   utils.StringPtr("cluster-1"),
				LastUpdateTime: &when,
				ExpiryTime:     &when,
				Status:         &status,
				SystemType:     &systemType,
				HasHspInUse:    utils.BoolPtr(true),
				HostIp: &import4.IPv4Address{
					Value:        utils.StringPtr("10.0.0.8"),
					PrefixLength: utils.IntPtr(24),
				},
				Links: []import1.ApiLink{{
					Rel:  utils.StringPtr("self"),
					Href: utils.StringPtr("https://example.test/passwords/ext-1"),
				}},
			}},
			want: []map[string]interface{}{{
				"ext_id":           "ext-1",
				"tenant_id":        "tenant-1",
				"username":         "admin",
				"cluster_ext_id":   "cluster-1",
				"last_update_time": "2024-01-02T03:04:05Z",
				"expiry_time":      "2024-01-02T03:04:05Z",
				"status":           "SECURE",
				"system_type":      "AHV",
				"has_hsp_in_use":   true,
				"host_ip": []map[string]interface{}{{
					"value":         "10.0.0.8",
					"prefix_length": 24,
				}},
				"links": []map[string]interface{}{{
					"rel":  "self",
					"href": "https://example.test/passwords/ext-1",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPasswordEntities(tt.in))
		})
	}
}

func TestUnitIsUnauthorizedErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "other", err: errors.New("connection reset"), want: false},
		{name: "401", err: errors.New("status 401"), want: true},
		{name: "unauthorized", err: errors.New("Unauthorized"), want: true},
		{name: "invalid auth credentials", err: errors.New("invalid auth credentials"), want: true},
		{name: "invalid credentials", err: errors.New("Invalid Credentials"), want: true},
		{name: "empty scheme", err: errors.New(`unsupported protocol scheme ""`), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, isUnauthorizedErr(tt.err))
		})
	}
}

func TestUnitPasswordManagerReadAndDelete(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixPasswordManagerV2().Schema, map[string]interface{}{
		"ext_id":       "ext-1",
		"new_password": "secret",
	})

	require.Nil(t, resourceNutanixPasswordManagerV2Read(t.Context(), d, nil))
	require.Nil(t, resourceNutanixPasswordManagerV2Delete(t.Context(), d, nil))
}

func statusPtr(v clusterConfig.PasswordStatus) *clusterConfig.PasswordStatus { return &v }

func systemTypePtr(v clusterConfig.SystemType) *clusterConfig.SystemType { return &v }
