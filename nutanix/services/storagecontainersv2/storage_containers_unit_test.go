package storagecontainersv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clustermgmt "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	clsConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/config"
	clsResponse "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/response"
	clsstats "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/stats"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitStorageContainerSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixStorageContainersV2().InternalValidate(ResourceNutanixStorageContainersV2().Schema, true))
	require.NoError(t, DatasourceNutanixStorageContainerV2().InternalValidate(DatasourceNutanixStorageContainerV2().Schema, false))
	require.NoError(t, DatasourceNutanixStorageContainersV2().InternalValidate(DatasourceNutanixStorageContainersV2().Schema, false))
	require.NoError(t, DatasourceNutanixStorageStatsInfoV2().InternalValidate(DatasourceNutanixStorageStatsInfoV2().Schema, false))
	require.NoError(t, resourceSchemaForValuePrefixLength().Elem.(*schema.Resource).InternalValidate(resourceSchemaForValuePrefixLength().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, resourceSchemaForFqdnValue().Elem.(*schema.Resource).InternalValidate(resourceSchemaForFqdnValue().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, SchemaForValueTimestamp().Elem.(*schema.Resource).InternalValidate(SchemaForValueTimestamp().Elem.(*schema.Resource).Schema, false))
	require.NoError(t, SchemaForValuePrefixLength().Elem.(*schema.Resource).InternalValidate(SchemaForValuePrefixLength().Elem.(*schema.Resource).Schema, false))
	require.NoError(t, SchemaForFqdnValue().Elem.(*schema.Resource).InternalValidate(SchemaForFqdnValue().Elem.(*schema.Resource).Schema, false))
}

func TestUnitFlattenCacheDeduplication(t *testing.T) {
	t.Parallel()

	redacted := clustermgmt.CACHEDEDUPLICATION_REDACTED
	none := clustermgmt.CACHEDEDUPLICATION_NONE
	off := clustermgmt.CACHEDEDUPLICATION_OFF
	on := clustermgmt.CACHEDEDUPLICATION_ON
	unknown := clustermgmt.CACHEDEDUPLICATION_UNKNOWN

	tests := []struct {
		name string
		in   *clustermgmt.CacheDeduplication
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "unknown", in: &unknown, want: "UNKNOWN"},
		{name: "redacted", in: &redacted, want: "REDACTED"},
		{name: "none", in: &none, want: "NONE"},
		{name: "off", in: &off, want: "OFF"},
		{name: "on", in: &on, want: "ON"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCacheDeduplication(tt.in))
		})
	}
}

func TestUnitFlattenErasureCodeStatus(t *testing.T) {
	t.Parallel()

	redacted := clustermgmt.ERASURECODESTATUS_REDACTED
	none := clustermgmt.ERASURECODESTATUS_NONE
	off := clustermgmt.ERASURECODESTATUS_OFF
	on := clustermgmt.ERASURECODESTATUS_ON

	tests := []struct {
		name string
		in   *clustermgmt.ErasureCodeStatus
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "redacted", in: &redacted, want: "REDACTED"},
		{name: "none", in: &none, want: "NONE"},
		{name: "off", in: &off, want: "OFF"},
		{name: "on", in: &on, want: "ON"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenErasureCodeStatus(tt.in))
		})
	}
}

func TestUnitFlattenOnDiskDedup(t *testing.T) {
	t.Parallel()

	redacted := clustermgmt.ONDISKDEDUP_REDACTED
	none := clustermgmt.ONDISKDEDUP_NONE
	off := clustermgmt.ONDISKDEDUP_OFF
	post := clustermgmt.ONDISKDEDUP_POST_PROCESS

	tests := []struct {
		name string
		in   *clustermgmt.OnDiskDedup
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "redacted", in: &redacted, want: "REDACTED"},
		{name: "none", in: &none, want: "NONE"},
		{name: "off", in: &off, want: "OFF"},
		{name: "post process", in: &post, want: "POST_PROCESS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenOnDiskDedup(tt.in))
		})
	}
}

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	href := utils.StringPtr("https://example.test/sc/1")
	rel := utils.StringPtr("self")

	tests := []struct {
		name string
		in   []clsResponse.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []clsResponse.ApiLink{}, want: nil},
		{
			name: "populated",
			in:   []clsResponse.ApiLink{{Href: href, Rel: rel}},
			want: []map[string]interface{}{{"href": href, "rel": rel}},
		},
		{
			name: "missing fields",
			in:   []clsResponse.ApiLink{{}},
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

func TestUnitFlattenIPv4Address(t *testing.T) {
	t.Parallel()

	value := utils.StringPtr("10.0.0.1")
	prefix := utils.IntPtr(24)

	tests := []struct {
		name string
		in   *clsConfig.IPv4Address
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &clsConfig.IPv4Address{},
			want: []interface{}{map[string]interface{}{"value": (*string)(nil), "prefix_length": (*int)(nil)}},
		},
		{
			name: "populated",
			in:   &clsConfig.IPv4Address{Value: value, PrefixLength: prefix},
			want: []interface{}{map[string]interface{}{"value": value, "prefix_length": prefix}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPv4Address(tt.in))
		})
	}
}

func TestUnitFlattenIPv6Address(t *testing.T) {
	t.Parallel()

	value := utils.StringPtr("2001:db8::1")
	prefix := utils.IntPtr(64)

	tests := []struct {
		name string
		in   *clsConfig.IPv6Address
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &clsConfig.IPv6Address{},
			want: []interface{}{map[string]interface{}{"value": (*string)(nil), "prefix_length": (*int)(nil)}},
		},
		{
			name: "populated",
			in:   &clsConfig.IPv6Address{Value: value, PrefixLength: prefix},
			want: []interface{}{map[string]interface{}{"value": value, "prefix_length": prefix}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPv6Address(tt.in))
		})
	}
}

func TestUnitFlattenFQDN(t *testing.T) {
	t.Parallel()

	value := utils.StringPtr("storage.example.test")

	tests := []struct {
		name string
		in   *clsConfig.FQDN
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &clsConfig.FQDN{},
			want: []interface{}{map[string]interface{}{"value": (*string)(nil)}},
		},
		{
			name: "populated",
			in:   &clsConfig.FQDN{Value: value},
			want: []interface{}{map[string]interface{}{"value": value}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFQDN(tt.in))
		})
	}
}

func TestUnitFlattenNfsWhitelistAddresses(t *testing.T) {
	t.Parallel()

	value := utils.StringPtr("10.0.0.8")
	prefix := utils.IntPtr(24)

	tests := []struct {
		name string
		in   []clsConfig.IPAddressOrFQDN
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []clsConfig.IPAddressOrFQDN{}, want: nil},
		{
			name: "populated",
			in: []clsConfig.IPAddressOrFQDN{{
				Ipv4: &clsConfig.IPv4Address{Value: value, PrefixLength: prefix},
			}},
			want: []map[string]interface{}{{
				"ipv4": []interface{}{map[string]interface{}{"value": value, "prefix_length": prefix}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNfsWhitelistAddresses(tt.in))
		})
	}
}

func TestUnitFlattenStorageContainers(t *testing.T) {
	t.Parallel()

	name := utils.StringPtr("sc-1")
	extID := utils.StringPtr("ext-1")
	erasure := clustermgmt.ERASURECODESTATUS_ON
	cache := clustermgmt.CACHEDEDUPLICATION_OFF
	dedup := clustermgmt.ONDISKDEDUP_NONE

	tests := []struct {
		name string
		in   []clustermgmt.StorageContainer
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []clustermgmt.StorageContainer{}, want: nil},
		{
			name: "populated",
			in: []clustermgmt.StorageContainer{{
				ContainerExtId:     extID,
				Name:               name,
				ErasureCode:        &erasure,
				CacheDeduplication: &cache,
				OnDiskDedup:        &dedup,
			}},
			want: nil, // asserted field-by-field below
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenStorageContainers(tt.in)
			if tt.name != "populated" {
				require.Equal(t, tt.want, got)
				return
			}
			require.Len(t, got, 1)
			m := got[0].(map[string]interface{})
			require.Equal(t, extID, m["container_ext_id"])
			require.Equal(t, name, m["name"])
			require.Equal(t, "ON", m["erasure_code"])
			require.Equal(t, "OFF", m["cache_deduplication"])
			require.Equal(t, "NONE", m["on_disk_dedup"])
		})
	}
}

func TestUnitFlattenValueTimestamp(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	val := utils.Int64Ptr(42)

	tests := []struct {
		name string
		in   []clsstats.TimeIntValuePair
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []clsstats.TimeIntValuePair{}, want: nil},
		{
			name: "populated",
			in:   []clsstats.TimeIntValuePair{{Value: val, Timestamp: &when}},
			want: []map[string]interface{}{{
				"value":     val,
				"timestamp": "2024-01-02T03:04:05Z",
			}},
		},
		{
			name: "missing fields",
			in:   []clsstats.TimeIntValuePair{{}},
			want: []map[string]interface{}{{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenValueTimestamp(tt.in))
		})
	}
}

func TestUnitExpandIPv4Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv4Address(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Address([]interface{}{map[string]interface{}{
			"value":         "10.0.0.1",
			"prefix_length": 24,
		}})
		require.NotNil(t, got)
		require.Equal(t, "10.0.0.1", utils.StringValue(got.Value))
		require.Equal(t, 24, utils.IntValue(got.PrefixLength))
	})
}

func TestUnitExpandIPv6Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv6Address(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6Address([]interface{}{map[string]interface{}{
			"value":         "2001:db8::1",
			"prefix_length": 64,
		}})
		require.NotNil(t, got)
		require.Equal(t, "2001:db8::1", utils.StringValue(got.Value))
		require.Equal(t, 64, utils.IntValue(got.PrefixLength))
	})
}

func TestUnitExpandFQDN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *string
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{name: "populated", in: []interface{}{map[string]interface{}{"value": "storage.example.test"}}, want: utils.StringPtr("storage.example.test")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandFQDN(tt.in)
			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.Equal(t, utils.StringValue(tt.want), utils.StringValue(got.Value))
		})
	}
}

func TestUnitExpandNfsWhitelistAddresses(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandNfsWhitelistAddresses(nil))
	})

	t.Run("populated from schema", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixStorageContainersV2().Schema, map[string]interface{}{
			"name":           "sc-1",
			"cluster_ext_id": "cluster-1",
			"nfs_whitelist_addresses": []interface{}{map[string]interface{}{
				"ipv4": []interface{}{map[string]interface{}{
					"value":         "10.0.0.8",
					"prefix_length": 24,
				}},
			}},
		})
		got := expandNfsWhitelistAddresses(d.Get("nfs_whitelist_addresses"))
		require.Len(t, got, 1)
		require.NotNil(t, got[0].Ipv4)
		require.Equal(t, "10.0.0.8", utils.StringValue(got[0].Ipv4.Value))
		require.Equal(t, 24, utils.IntValue(got[0].Ipv4.PrefixLength))
	})
}
