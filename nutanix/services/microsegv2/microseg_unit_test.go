package microsegv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonCfg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/response"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitMicrosegSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixEntityGroupV2().InternalValidate(ResourceNutanixEntityGroupV2().Schema, true))
	require.NoError(t, ResourceNutanixCategoryMappingV2().InternalValidate(ResourceNutanixCategoryMappingV2().Schema, true))
	require.NoError(t, ResourceNutanixDirectoryServerConfigV2().InternalValidate(ResourceNutanixDirectoryServerConfigV2().Schema, true))
	require.NoError(t, ResourceNutanixNetworkSecurityPolicyExportV2().InternalValidate(ResourceNutanixNetworkSecurityPolicyExportV2().Schema, true))
	require.NoError(t, ResourceNutanixNetworkSecurityPolicyImportV2().InternalValidate(ResourceNutanixNetworkSecurityPolicyImportV2().Schema, true))

	require.NoError(t, DatasourceNutanixEntityGroupV2().InternalValidate(DatasourceNutanixEntityGroupV2().Schema, false))
	require.NoError(t, DatasourceNutanixEntityGroupsV2().InternalValidate(DatasourceNutanixEntityGroupsV2().Schema, false))
	require.NoError(t, DatasourceNutanixCategoryMappingV2().InternalValidate(DatasourceNutanixCategoryMappingV2().Schema, false))
	require.NoError(t, DatasourceNutanixCategoryMappingsV2().InternalValidate(DatasourceNutanixCategoryMappingsV2().Schema, false))
	require.NoError(t, DatasourceNutanixDirectoryServerConfigV2().InternalValidate(DatasourceNutanixDirectoryServerConfigV2().Schema, false))
	require.NoError(t, DatasourceNutanixDirectoryServerConfigsV2().InternalValidate(DatasourceNutanixDirectoryServerConfigsV2().Schema, false))
	require.NoError(t, DataSourceNutanixNetworkSecurityPolicyRulesV2().InternalValidate(DataSourceNutanixNetworkSecurityPolicyRulesV2().Schema, false))
}

func TestUnitExpandAllowedConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{name: "no entities", in: []interface{}{map[string]interface{}{}}, nilOut: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"entities": []interface{}{map[string]interface{}{
					"selected_by": "EXT_ID",
					"type":        "VM",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandAllowedConfig(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.Len(t, got.Entities, 1)
			require.Equal(t, "EXT_ID", got.Entities[0].SelectBy.GetName())
			require.Equal(t, "VM", got.Entities[0].Type.GetName())
		})
	}
}

func TestUnitExpandAllowedEntities(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, expandAllowedEntities([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandAllowedEntities([]interface{}{map[string]interface{}{
			"selected_by":       "IP_VALUES",
			"type":              "ADDRESS_GROUP",
			"reference_ext_ids": []interface{}{"ref-1"},
			"fqdns":             []interface{}{"example.test"},
			"kube_entities":     []interface{}{"kube-1"},
			"match_criteria":    "EQUALS",
			"reference_string":  "ref-str",
			"addresses": []interface{}{map[string]interface{}{
				"ipv4_addresses": []interface{}{map[string]interface{}{
					"value":         "10.0.0.1",
					"prefix_length": 32,
				}},
			}},
			"ip_ranges": []interface{}{map[string]interface{}{
				"ipv4_ranges": []interface{}{map[string]interface{}{
					"start_ip": "10.0.0.1",
					"end_ip":   "10.0.0.10",
				}},
			}},
		}})
		require.Len(t, got, 1)
		require.Equal(t, "IP_VALUES", got[0].SelectBy.GetName())
		require.Equal(t, "ADDRESS_GROUP", got[0].Type.GetName())
		require.Equal(t, []string{"ref-1"}, got[0].ReferenceExtIds)
		require.Equal(t, []string{"example.test"}, got[0].Fqdns)
		require.Equal(t, []string{"kube-1"}, got[0].KubeEntities)
		require.Equal(t, "EQUALS", got[0].MatchCriteria.GetName())
		require.Equal(t, "ref-str", utils.StringValue(got[0].ReferenceString))
		require.NotNil(t, got[0].Addresses)
		require.NotNil(t, got[0].IpRanges)
	})
}

func TestUnitExpandExceptConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"entities": []interface{}{map[string]interface{}{
					"selected_by": "EXT_ID",
					"type":        "ADDRESS_GROUP",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandExceptConfig(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.Len(t, got.Entities, 1)
		})
	}
}

func TestUnitExpandExceptEntities(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, expandExceptEntities([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandExceptEntities([]interface{}{map[string]interface{}{
			"selected_by":       "EXT_ID",
			"type":              "ADDRESS_GROUP",
			"reference_ext_ids": []interface{}{"ex-1"},
		}})
		require.Len(t, got, 1)
		require.Equal(t, "EXT_ID", got[0].SelectBy.GetName())
		require.Equal(t, []string{"ex-1"}, got[0].ReferenceExtIds)
	})
}

func TestUnitExpandAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{name: "no ipv4", in: []interface{}{map[string]interface{}{}}, nilOut: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"ipv4_addresses": []interface{}{map[string]interface{}{
					"value":         "192.168.1.1",
					"prefix_length": 24,
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandAddresses(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.Len(t, got.Ipv4Addresses, 1)
			require.Equal(t, "192.168.1.1", utils.StringValue(got.Ipv4Addresses[0].Value))
			require.Equal(t, 24, utils.IntValue(got.Ipv4Addresses[0].PrefixLength))
		})
	}
}

func TestUnitExpandIpRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"ipv4_ranges": []interface{}{map[string]interface{}{
					"start_ip": "10.0.0.1",
					"end_ip":   "10.0.0.5",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandIpRange(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.Len(t, got.Ipv4Ranges, 1)
			require.Equal(t, "10.0.0.1", utils.StringValue(got.Ipv4Ranges[0].StartIp))
			require.Equal(t, "10.0.0.5", utils.StringValue(got.Ipv4Ranges[0].EndIp))
		})
	}
}

func TestUnitFlattenLinksEntityGroup(t *testing.T) {
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
				Href: utils.StringPtr("https://example.test/eg/1"),
				Rel:  utils.StringPtr("self"),
			}},
			want: []map[string]interface{}{{
				"href": "https://example.test/eg/1",
				"rel":  "self",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLinksEntityGroup(tt.in))
		})
	}
}

func TestUnitFlattenAllowedConfig(t *testing.T) {
	t.Parallel()

	selectBy := import2.ALLOWEDSELECTBY_EXT_ID
	typ := import2.ALLOWEDTYPE_VM
	tests := []struct {
		name string
		in   *import2.AllowedConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import2.AllowedConfig{
				Entities: []import2.AllowedEntity{{
					SelectBy: &selectBy,
					Type:     &typ,
				}},
			},
			want: []map[string]interface{}{{
				"entities": []map[string]interface{}{{
					"selected_by":       "EXT_ID",
					"type":              "VM",
					"addresses":         []map[string]interface{}(nil),
					"ip_ranges":         []map[string]interface{}(nil),
					"kube_entities":     []string(nil),
					"reference_ext_ids": []string(nil),
					"fqdns":             []string(nil),
					"reference_string":  "",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAllowedConfig(tt.in))
		})
	}
}

func TestUnitFlattenAllowedEntities(t *testing.T) {
	t.Parallel()

	selectBy := import2.ALLOWEDSELECTBY_NAME
	typ := import2.ALLOWEDTYPE_SUBNET
	criteria := import2.MATCHCRITERIA_CONTAINS
	tests := []struct {
		name string
		in   []import2.AllowedEntity
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.AllowedEntity{}, want: nil},
		{
			name: "populated",
			in: []import2.AllowedEntity{{
				SelectBy:        &selectBy,
				Type:            &typ,
				MatchCriteria:   &criteria,
				ReferenceString: utils.StringPtr("web"),
				Fqdns:           []string{"a.test"},
			}},
			want: []map[string]interface{}{{
				"selected_by":       "NAME",
				"type":              "SUBNET",
				"addresses":         []map[string]interface{}(nil),
				"ip_ranges":         []map[string]interface{}(nil),
				"kube_entities":     []string(nil),
				"reference_ext_ids": []string(nil),
				"fqdns":             []string{"a.test"},
				"reference_string":  "web",
				"match_criteria":    "CONTAINS",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAllowedEntities(tt.in))
		})
	}
}

func TestUnitFlattenExceptConfig(t *testing.T) {
	t.Parallel()

	selectBy := import2.EXCEPTSELECTBY_EXT_ID
	typ := import2.EXCEPTTYPE_ADDRESS_GROUP
	tests := []struct {
		name string
		in   *import2.ExceptConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import2.ExceptConfig{
				Entities: []import2.ExceptEntity{{
					SelectBy: &selectBy,
					Type:     &typ,
				}},
			},
			want: []map[string]interface{}{{
				"entities": []map[string]interface{}{{
					"selected_by":       "EXT_ID",
					"type":              "ADDRESS_GROUP",
					"addresses":         []map[string]interface{}(nil),
					"ip_ranges":         []map[string]interface{}(nil),
					"reference_ext_ids": []string(nil),
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExceptConfig(tt.in))
		})
	}
}

func TestUnitFlattenExceptEntities(t *testing.T) {
	t.Parallel()

	selectBy := import2.EXCEPTSELECTBY_IP_VALUES
	typ := import2.EXCEPTTYPE_ADDRESS_GROUP
	tests := []struct {
		name string
		in   []import2.ExceptEntity
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.ExceptEntity{}, want: nil},
		{
			name: "populated",
			in: []import2.ExceptEntity{{
				SelectBy:        &selectBy,
				Type:            &typ,
				ReferenceExtIds: []string{"ex-1"},
			}},
			want: []map[string]interface{}{{
				"selected_by":       "IP_VALUES",
				"type":              "ADDRESS_GROUP",
				"addresses":         []map[string]interface{}(nil),
				"ip_ranges":         []map[string]interface{}(nil),
				"reference_ext_ids": []string{"ex-1"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExceptEntities(tt.in))
		})
	}
}

func TestUnitFlattenAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import2.Addresses
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: &import2.Addresses{}, want: nil},
		{
			name: "populated",
			in: &import2.Addresses{
				Ipv4Addresses: []commonCfg.IPv4Address{{
					Value:        utils.StringPtr("10.0.0.8"),
					PrefixLength: utils.IntPtr(32),
				}},
			},
			want: []map[string]interface{}{{
				"ipv4_addresses": []map[string]interface{}{{
					"value":         "10.0.0.8",
					"prefix_length": 32,
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAddresses(tt.in))
		})
	}
}

func TestUnitFlattenIpRanges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import2.IpRange
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: &import2.IpRange{}, want: nil},
		{
			name: "populated",
			in: &import2.IpRange{
				Ipv4Ranges: []import2.IPv4Range{{
					StartIp: utils.StringPtr("10.0.0.1"),
					EndIp:   utils.StringPtr("10.0.0.2"),
				}},
			},
			want: []map[string]interface{}{{
				"ipv4_ranges": []map[string]interface{}{{
					"start_ip": "10.0.0.1",
					"end_ip":   "10.0.0.2",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIpRanges(tt.in))
		})
	}
}

func TestUnitExpandDirectoryServerConfigMatchingCriterias(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want int
	}{
		{name: "nil", in: nil, want: 0},
		{name: "empty", in: []interface{}{}, want: 0},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"criteria":     "prod",
				"match_entity": "VM",
				"match_field":  "NAME",
				"match_type":   "CONTAINS",
			}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandDirectoryServerConfigMatchingCriterias(tt.in)
			if tt.want == 0 {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.Equal(t, "prod", utils.StringValue(got[0].Criteria))
			require.Equal(t, "VM", got[0].MatchEntity.GetName())
			require.Equal(t, "NAME", got[0].MatchField.GetName())
			require.Equal(t, "CONTAINS", got[0].MatchType.GetName())
		})
	}
}

func TestUnitExpandIPAddressOrFQDNList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want int
	}{
		{name: "nil", in: nil, want: 0},
		{name: "empty", in: []interface{}{}, want: 0},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"ipv4": []interface{}{map[string]interface{}{"value": "1.2.3.4", "prefix_length": 32}},
				"ipv6": []interface{}{map[string]interface{}{"value": "::1", "prefix_length": 128}},
				"fqdn": []interface{}{map[string]interface{}{"value": "dc.example.test"}},
			}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandIPAddressOrFQDNList(tt.in)
			if tt.want == 0 {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.Equal(t, "1.2.3.4", utils.StringValue(got[0].Ipv4.Value))
			require.Equal(t, "::1", utils.StringValue(got[0].Ipv6.Value))
			require.Equal(t, "dc.example.test", utils.StringValue(got[0].Fqdn.Value))
		})
	}
}

func TestUnitExpandIPv4AddressValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
		value  string
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{name: "populated", in: []interface{}{map[string]interface{}{"value": "10.0.0.1", "prefix_length": 24}}, value: "10.0.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandIPv4AddressValue(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.value, utils.StringValue(got.Value))
			require.Equal(t, 24, utils.IntValue(got.PrefixLength))
		})
	}
}

func TestUnitExpandIPv6AddressValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{name: "populated", in: []interface{}{map[string]interface{}{"value": "fe80::1", "prefix_length": 64}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandIPv6AddressValue(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.Equal(t, "fe80::1", utils.StringValue(got.Value))
		})
	}
}

func TestUnitExpandFQDNValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{name: "populated", in: []interface{}{map[string]interface{}{"value": "host.example.test"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandFQDNValue(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.Equal(t, "host.example.test", utils.StringValue(got.Value))
		})
	}
}

func TestUnitExpandAdInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"directory_service_reference": "ds-1",
				"object_identifier":           "obj-1",
				"status":                      "USABLE",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandAdInfo(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.Equal(t, "ds-1", utils.StringValue(got.DirectoryServiceReference))
			require.Equal(t, "obj-1", utils.StringValue(got.ObjectIdentifier))
			require.Equal(t, "USABLE", got.Status.GetName())
		})
	}
}

func TestUnitFlattenMatchingCriterias(t *testing.T) {
	t.Parallel()

	entity := import2.MATCHENTITY_VM
	field := import2.MATCHFIELD_NAME
	matchType := import2.MATCHTYPE_CONTAINS
	tests := []struct {
		name string
		in   []import2.MatchingCriteria
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.MatchingCriteria{}, want: nil},
		{
			name: "populated",
			in: []import2.MatchingCriteria{{
				Criteria:    utils.StringPtr("prod"),
				MatchEntity: &entity,
				MatchField:  &field,
				MatchType:   &matchType,
			}},
			want: []map[string]interface{}{{
				"criteria":     "prod",
				"match_entity": "VM",
				"match_field":  "NAME",
				"match_type":   "CONTAINS",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenMatchingCriterias(tt.in))
		})
	}
}

func TestUnitFlattenIPAddressOrFQDNList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []commonCfg.IPAddressOrFQDN
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []commonCfg.IPAddressOrFQDN{}, want: nil},
		{
			name: "populated",
			in: []commonCfg.IPAddressOrFQDN{{
				Ipv4: &commonCfg.IPv4Address{Value: utils.StringPtr("1.1.1.1"), PrefixLength: utils.IntPtr(32)},
				Fqdn: &commonCfg.FQDN{Value: utils.StringPtr("dc.example.test")},
			}},
			want: []map[string]interface{}{{
				"ipv4": []map[string]interface{}{{"value": "1.1.1.1", "prefix_length": 32}},
				"fqdn": []map[string]interface{}{{"value": "dc.example.test"}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPAddressOrFQDNList(tt.in))
		})
	}
}

func TestUnitFlattenAdInfo(t *testing.T) {
	t.Parallel()

	status := import2.ADSTATUS_USABLE
	tests := []struct {
		name string
		in   *import2.AdInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import2.AdInfo{
				DirectoryServiceReference: utils.StringPtr("ds-1"),
				ObjectIdentifier:          utils.StringPtr("obj-1"),
				ObjectPath:                utils.StringPtr("/OU=Users"),
				Status:                    &status,
			},
			want: []map[string]interface{}{{
				"directory_service_reference": "ds-1",
				"object_identifier":           "obj-1",
				"object_path":                 "/OU=Users",
				"status":                      "USABLE",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAdInfo(tt.in))
		})
	}
}

func TestUnitFlattenCategoryMappings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import2.CategoryMapping
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []import2.CategoryMapping{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []import2.CategoryMapping{{
				ExtId:         utils.StringPtr("cm-1"),
				Name:          utils.StringPtr("map-1"),
				CategoryName:  utils.StringPtr("AppType"),
				CategoryValue: utils.StringPtr("Web"),
				ProjectExtId:  utils.StringPtr("p-1"),
				TenantId:      utils.StringPtr("t-1"),
			}},
			want: []map[string]interface{}{{
				"ext_id":         "cm-1",
				"name":           "map-1",
				"category_name":  "AppType",
				"category_value": "Web",
				"ad_info":        []map[string]interface{}(nil),
				"project_ext_id": "p-1",
				"tenant_id":      "t-1",
				"links":          []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCategoryMappings(tt.in))
		})
	}
}

func TestUnitFlattenEntityGroups(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)
	tests := []struct {
		name string
		in   []import2.EntityGroup
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []import2.EntityGroup{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []import2.EntityGroup{{
				ExtId:                   utils.StringPtr("eg-1"),
				Name:                    utils.StringPtr("group-1"),
				Description:             utils.StringPtr("desc"),
				PolicyExtIds:            []string{"pol-1"},
				OwnerExtId:              utils.StringPtr("owner-1"),
				TenantId:                utils.StringPtr("t-1"),
				ProjectExtId:            utils.StringPtr("p-1"),
				CreationTime:            &when,
				LastUpdateTime:          &when,
				IsSharedWithAllProjects: utils.BoolPtr(true),
				IsSystemDefined:         utils.BoolPtr(false),
			}},
			want: []map[string]interface{}{{
				"ext_id":                      "eg-1",
				"name":                        "group-1",
				"description":                 "desc",
				"allowed_config":              []map[string]interface{}(nil),
				"except_config":               []map[string]interface{}(nil),
				"policy_ext_ids":              []string{"pol-1"},
				"owner_ext_id":                "owner-1",
				"tenant_id":                   "t-1",
				"links":                       []map[string]interface{}(nil),
				"project_ext_id":              "p-1",
				"creation_time":               "2024-05-06T07:08:09.000Z",
				"last_update_time":            "2024-05-06T07:08:09.000Z",
				"is_shared_with_all_projects": true,
				"is_system_defined":           false,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEntityGroups(tt.in))
		})
	}
}

func TestUnitFlattenDirectoryServerConfigs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import2.DirectoryServerConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []import2.DirectoryServerConfig{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []import2.DirectoryServerConfig{{
				ExtId:                            utils.StringPtr("dsc-1"),
				DirectoryServiceReference:        utils.StringPtr("ds-1"),
				IsDefaultCategoryEnabled:         utils.BoolPtr(true),
				ShouldKeepDefaultCategoryOnLogin: utils.BoolPtr(false),
				ProjectExtId:                     utils.StringPtr("p-1"),
				TenantId:                         utils.StringPtr("t-1"),
			}},
			want: []map[string]interface{}{{
				"ext_id":                                "dsc-1",
				"directory_service_reference":           "ds-1",
				"domain_controllers":                    []map[string]interface{}(nil),
				"matching_criterias":                    []map[string]interface{}(nil),
				"is_default_category_enabled":           true,
				"should_keep_default_category_on_login": false,
				"project_ext_id":                        "p-1",
				"tenant_id":                             "t-1",
				"links":                                 []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDirectoryServerConfigs(tt.in))
		})
	}
}

func TestUnitExpandCategoryMappingBody(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixCategoryMappingV2().Schema, map[string]interface{}{
		"name":           "map-1",
		"category_name":  "AppType",
		"category_value": "Web",
		"project_ext_id": "p-1",
		"ad_info": []interface{}{map[string]interface{}{
			"directory_service_reference": "ds-1",
			"object_identifier":           "obj-1",
		}},
	})
	got := expandCategoryMappingBody(d)
	require.Equal(t, "map-1", utils.StringValue(got.Name))
	require.Equal(t, "AppType", utils.StringValue(got.CategoryName))
	require.Equal(t, "Web", utils.StringValue(got.CategoryValue))
	require.Equal(t, "p-1", utils.StringValue(got.ProjectExtId))
	require.NotNil(t, got.AdInfo)
}

func TestUnitExpandDirectoryServerConfigBody(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixDirectoryServerConfigV2().Schema, map[string]interface{}{
		"directory_service_reference":           "ds-1",
		"is_default_category_enabled":           true,
		"should_keep_default_category_on_login": false,
		"project_ext_id":                        "p-1",
		"domain_controllers": []interface{}{map[string]interface{}{
			"fqdn": []interface{}{map[string]interface{}{"value": "dc.example.test"}},
		}},
		"matching_criterias": []interface{}{map[string]interface{}{
			"criteria":     "prod",
			"match_entity": "VM",
			"match_field":  "NAME",
			"match_type":   "CONTAINS",
		}},
	})
	got := expandDirectoryServerConfigBody(d)
	require.Equal(t, "ds-1", utils.StringValue(got.DirectoryServiceReference))
	require.Equal(t, true, utils.BoolValue(got.IsDefaultCategoryEnabled))
	require.Equal(t, false, utils.BoolValue(got.ShouldKeepDefaultCategoryOnLogin))
	require.Equal(t, "p-1", utils.StringValue(got.ProjectExtId))
	require.Len(t, got.DomainControllers, 1)
	require.Len(t, got.MatchingCriterias, 1)
}

func TestUnitFlattenNetworkSecurityPolicyRules(t *testing.T) {
	t.Parallel()

	ruleType := import2.RULETYPE_APPLICATION
	tests := []struct {
		name string
		in   []import2.NetworkSecurityPolicyRule
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.NetworkSecurityPolicyRule{}, want: nil},
		{
			name: "populated",
			in: []import2.NetworkSecurityPolicyRule{{
				ExtId:       utils.StringPtr("rule-1"),
				Description: utils.StringPtr("desc"),
				TenantId:    utils.StringPtr("t-1"),
				Type:        &ruleType,
			}},
			want: []interface{}{map[string]interface{}{
				"ext_id":      "rule-1",
				"description": "desc",
				"tenant_id":   "t-1",
				"type":        "APPLICATION",
				"links":       []map[string]interface{}(nil),
				"spec":        []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNetworkSecurityPolicyRules(tt.in))
		})
	}
}

func TestUnitFlattenLinksMicroseg(t *testing.T) {
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
				Href: utils.StringPtr("https://example.test/r/1"),
				Rel:  utils.StringPtr("self"),
			}},
			want: []map[string]interface{}{{
				"href": utils.StringPtr("https://example.test/r/1"),
				"rel":  utils.StringPtr("self"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLinksMicroseg(tt.in))
		})
	}
}

func TestUnitFlattenRuleTypeMicroseg(t *testing.T) {
	t.Parallel()

	app := import2.RULETYPE_APPLICATION
	tests := []struct {
		name string
		in   *import2.RuleType
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "populated", in: &app, want: "APPLICATION"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRuleTypeMicroseg(tt.in))
		})
	}
}

func TestUnitFlattenOneOfNetworkSecurityPolicyRuleSpecMicroseg(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfNetworkSecurityPolicyRuleSpecMicroseg(nil))
	})
	t.Run("empty object type", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfNetworkSecurityPolicyRuleSpecMicroseg(&import2.OneOfNetworkSecurityPolicyRuleSpec{}))
	})
	t.Run("populated two env", func(t *testing.T) {
		t.Parallel()
		spec := import2.NewOneOfNetworkSecurityPolicyRuleSpec()
		v := import2.NewTwoEnvIsolationRuleSpec()
		v.FirstIsolationGroup = []string{"cat-1"}
		v.SecondIsolationGroup = []string{"cat-2"}
		require.NoError(t, spec.SetValue(*v))
		got := flattenOneOfNetworkSecurityPolicyRuleSpecMicroseg(spec)
		require.Equal(t, []map[string]interface{}{{
			"two_env_isolation_rule_spec": []map[string]interface{}{{
				"first_isolation_group":  []string{"cat-1"},
				"second_isolation_group": []string{"cat-2"},
			}},
		}}, got)
	})
}

func TestUnitFlattenIPv4Microseg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *commonCfg.IPv4Address
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &commonCfg.IPv4Address{Value: utils.StringPtr("10.0.0.0"), PrefixLength: utils.IntPtr(24)},
			want: []interface{}{map[string]interface{}{"value": "10.0.0.0", "prefix_length": 24}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPv4Microseg(tt.in))
		})
	}
}

func TestUnitFlattenTCPPortRangeMicroseg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import2.TcpPortRangeSpec
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.TcpPortRangeSpec{}, want: nil},
		{
			name: "populated",
			in:   []import2.TcpPortRangeSpec{{StartPort: utils.IntPtr(80), EndPort: utils.IntPtr(80)}},
			want: []interface{}{map[string]interface{}{"start_port": 80, "end_port": 80}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTCPPortRangeMicroseg(tt.in))
		})
	}
}

func TestUnitFlattenUDPPortRangeMicroseg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import2.UdpPortRangeSpec
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.UdpPortRangeSpec{}, want: nil},
		{
			name: "populated",
			in:   []import2.UdpPortRangeSpec{{StartPort: utils.IntPtr(53), EndPort: utils.IntPtr(53)}},
			want: []interface{}{map[string]interface{}{"start_port": 53, "end_port": 53}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenUDPPortRangeMicroseg(tt.in))
		})
	}
}

func TestUnitFlattenIcmpTypeCodeMicroseg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import2.IcmpTypeCodeSpec
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.IcmpTypeCodeSpec{}, want: nil},
		{
			name: "populated",
			in: []import2.IcmpTypeCodeSpec{{
				IsAllAllowed: utils.BoolPtr(false),
				Type:         utils.IntPtr(8),
				Code:         utils.IntPtr(0),
			}},
			want: []interface{}{map[string]interface{}{
				"is_all_allowed": false,
				"type":           8,
				"code":           0,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIcmpTypeCodeMicroseg(tt.in))
		})
	}
}
