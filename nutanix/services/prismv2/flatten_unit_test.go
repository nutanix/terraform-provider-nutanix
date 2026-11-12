package prismv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clustermgmtConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/clustermgmt/v4/config"
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/config"
	commonResponse "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFlattenTime(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		name string
		in   *time.Time
		want *string
	}{
		{name: "nil", in: nil, want: nil},
		{name: "populated", in: &when, want: utils.StringPtr(when.String())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTime(tt.in))
		})
	}
}

func TestUnitFlattenCategoryType(t *testing.T) {
	t.Parallel()

	user := config.CategoryType(2)
	system := config.CategoryType(3)
	internal := config.CategoryType(4)
	other := config.CategoryType(99)

	tests := []struct {
		name string
		in   *config.CategoryType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "user", in: &user, want: "USER"},
		{name: "system", in: &system, want: "SYSTEM"},
		{name: "internal", in: &internal, want: "INTERNAL"},
		{name: "unrecognized", in: &other, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCategoryType(tt.in))
		})
	}
}

func TestUnitFlattenResourceGroup(t *testing.T) {
	t.Parallel()

	entity := config.ResourceGroup(2)
	policy := config.ResourceGroup(3)
	other := config.ResourceGroup(99)

	tests := []struct {
		name string
		in   *config.ResourceGroup
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "entity", in: &entity, want: "ENTITY"},
		{name: "policy", in: &policy, want: "POLICY"},
		{name: "unrecognized", in: &other, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceGroup(tt.in))
		})
	}
}

func TestUnitFlattenResourceType(t *testing.T) {
	t.Parallel()

	vm := config.ResourceType(2)
	mhVM := config.ResourceType(3)
	policySchema := config.ResourceType(26)
	other := config.ResourceType(99)

	tests := []struct {
		name string
		in   *config.ResourceType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "vm", in: &vm, want: "VM"},
		{name: "mh vm", in: &mhVM, want: "MH_VM"},
		{name: "policy schema", in: &policySchema, want: "POLICY_SCHEMA"},
		{name: "unrecognized", in: &other, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceType(tt.in))
		})
	}
}

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	href := utils.StringPtr("https://example.test/categories/1")
	rel := utils.StringPtr("self")

	tests := []struct {
		name string
		in   []commonResponse.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []commonResponse.ApiLink{}, want: nil},
		{
			name: "populated",
			in:   []commonResponse.ApiLink{{Href: href, Rel: rel}},
			want: []map[string]interface{}{{"href": href, "rel": rel}},
		},
		{
			name: "missing fields",
			in:   []commonResponse.ApiLink{{}},
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

func TestUnitFlattenAssociationSummary(t *testing.T) {
	t.Parallel()

	rg := config.ResourceGroup(2)
	rt := config.ResourceType(2)
	catID := utils.StringPtr("cat-1")
	count := utils.IntPtr(3)

	tests := []struct {
		name string
		in   []config.AssociationSummary
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{name: "empty", in: []config.AssociationSummary{}, want: []interface{}{}},
		{
			name: "populated",
			in: []config.AssociationSummary{{
				CategoryId:    catID,
				Count:         count,
				ResourceGroup: &rg,
				ResourceType:  &rt,
			}},
			want: []interface{}{map[string]interface{}{
				"category_id":    catID,
				"count":          count,
				"resource_group": "ENTITY",
				"resource_type":  "VM",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAssociationSummary(tt.in))
		})
	}
}

func TestUnitFlattenAssociationDetail(t *testing.T) {
	t.Parallel()

	rg := config.ResourceGroup(3)
	rt := config.ResourceType(6)
	catID := utils.StringPtr("cat-2")
	resID := utils.StringPtr("res-1")

	tests := []struct {
		name string
		in   []config.AssociationDetail
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{name: "empty", in: []config.AssociationDetail{}, want: []interface{}{}},
		{
			name: "populated",
			in: []config.AssociationDetail{{
				CategoryId:    catID,
				ResourceGroup: &rg,
				ResourceType:  &rt,
				ResourceId:    resID,
			}},
			want: []interface{}{map[string]interface{}{
				"category_id":    catID,
				"resource_group": "POLICY",
				"resource_type":  "CLUSTER",
				"resource_id":    resID,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAssociationDetail(tt.in))
		})
	}
}

func TestUnitFlattenCategoriesEntities(t *testing.T) {
	t.Parallel()

	catType := config.CategoryType(2)
	extID := utils.StringPtr("ext-1")
	key := utils.StringPtr("Environment")
	value := utils.StringPtr("Prod")
	desc := utils.StringPtr("desc")
	owner := utils.StringPtr("owner-1")
	project := utils.StringPtr("proj-1")
	shared := []string{"proj-a"}
	sharedAll := utils.BoolPtr(true)
	tenant := utils.StringPtr("tenant-1")
	href := utils.StringPtr("https://example.test/categories/ext-1")
	rel := utils.StringPtr("self")

	tests := []struct {
		name string
		in   []config.Category
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []config.Category{}, want: nil},
		{
			name: "populated",
			in: []config.Category{{
				ExtId:                   extID,
				Key:                     key,
				Value:                   value,
				Type:                    &catType,
				Description:             desc,
				OwnerUuid:               owner,
				ProjectExtId:            project,
				SharedWithProjects:      shared,
				IsSharedWithAllProjects: sharedAll,
				TenantId:                tenant,
				Links:                   []commonResponse.ApiLink{{Href: href, Rel: rel}},
			}},
			want: []interface{}{map[string]interface{}{
				"ext_id":                   extID,
				"key":                      key,
				"value":                    value,
				"type":                     "USER",
				"description":              desc,
				"owner_uuid":               owner,
				"project_ext_id":           project,
				"shared_with_projects":     shared,
				"shared_with_all_projects": sharedAll,
				"associations":             []interface{}{},
				"detailed_associations":    []interface{}{},
				"tenant_id":                tenant,
				"links":                    []map[string]interface{}{{"href": href, "rel": rel}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCategoriesEntities(tt.in))
		})
	}
}

func TestUnitFlattenClusterSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   config.Size
		want string
	}{
		{name: "starter", in: config.Size(2), want: "STARTER"},
		{name: "small", in: config.Size(3), want: "SMALL"},
		{name: "large", in: config.Size(4), want: "LARGE"},
		{name: "extralarge", in: config.Size(5), want: "EXTRALARGE"},
		{name: "unknown", in: config.Size(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenClusterSize(tt.in))
		})
	}
}

func TestUnitFlattenEnvironmentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   config.EnvironmentType
		want string
	}{
		{name: "onprem", in: config.EnvironmentType(2), want: "ONPREM"},
		{name: "ntnx cloud", in: config.EnvironmentType(3), want: "NTNX_CLOUD"},
		{name: "unknown", in: config.EnvironmentType(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEnvironmentType(tt.in))
		})
	}
}

func TestUnitFlattenEnvironmentProvisioningType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   config.ProvisioningType
		want string
	}{
		{name: "ntnx", in: config.ProvisioningType(2), want: "NTNX"},
		{name: "native", in: config.ProvisioningType(3), want: "NATIVE"},
		{name: "unknown", in: config.ProvisioningType(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEnvironmentProvisioningType(tt.in))
		})
	}
}

func TestUnitFlattenEnvironmentProviderType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   config.ProviderType
		want string
	}{
		{name: "ntnx", in: config.ProviderType(2), want: "NTNX"},
		{name: "azure", in: config.ProviderType(3), want: "AZURE"},
		{name: "aws", in: config.ProviderType(4), want: "AWS"},
		{name: "gcp", in: config.ProviderType(5), want: "GCP"},
		{name: "vsphere", in: config.ProviderType(6), want: "VSPHERE"},
		{name: "unknown", in: config.ProviderType(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEnvironmentProviderType(tt.in))
		})
	}
}

func TestUnitFlattenPCConfig(t *testing.T) {
	t.Parallel()

	size := config.Size(3)
	tests := []struct {
		name string
		in   *config.DomainManagerClusterConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &config.DomainManagerClusterConfig{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &config.DomainManagerClusterConfig{
				ShouldEnableLockdownMode: utils.BoolPtr(true),
				Name:                     utils.StringPtr("pc-1"),
				Size:                     &size,
				BuildInfo: &clustermgmtConfig.BuildInfo{
					Version: utils.StringPtr("2024.1"),
				},
			},
			want: []map[string]interface{}{{
				"should_enable_lockdown_mode": true,
				"name":                        "pc-1",
				"size":                        "SMALL",
				"build_info":                  []map[string]interface{}{{"version": "2024.1"}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPCConfig(tt.in))
		})
	}
}

func TestUnitFlattenBootstrapConfig(t *testing.T) {
	t.Parallel()

	envType := config.EnvironmentType(2)
	tests := []struct {
		name string
		in   *config.BootstrapConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &config.BootstrapConfig{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &config.BootstrapConfig{
				EnvironmentInfo: &config.EnvironmentInfo{Type: &envType},
			},
			want: []map[string]interface{}{{
				"environment_info": []map[string]interface{}{{"type": "ONPREM"}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenBootstrapConfig(tt.in))
		})
	}
}

func TestUnitFlattenResourceConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.DomainManagerResourceConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &config.DomainManagerResourceConfig{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &config.DomainManagerResourceConfig{
				NumVcpus:          utils.IntPtr(8),
				MemorySizeBytes:   utils.Int64Ptr(1024),
				DataDiskSizeBytes: utils.Int64Ptr(2048),
				ContainerExtIds:   []string{"ctr-1"},
			},
			want: []map[string]interface{}{{
				"num_vcpus":            8,
				"memory_size_bytes":    int64(1024),
				"data_disk_size_bytes": int64(2048),
				"container_ext_ids":    []string{"ctr-1"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceConfig(tt.in))
		})
	}
}

func TestUnitFlattenEnvironmentInfo(t *testing.T) {
	t.Parallel()

	envType := config.EnvironmentType(2)
	provider := config.ProviderType(2)
	provisioning := config.ProvisioningType(2)

	tests := []struct {
		name string
		in   *config.EnvironmentInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &config.EnvironmentInfo{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &config.EnvironmentInfo{
				Type:             &envType,
				ProviderType:     &provider,
				ProvisioningType: &provisioning,
			},
			want: []map[string]interface{}{{
				"type":              "ONPREM",
				"provider_type":     "NTNX",
				"provisioning_type": "NTNX",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEnvironmentInfo(tt.in))
		})
	}
}

func TestUnitFlattenPCNetwork(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.DomainManagerNetwork
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &config.DomainManagerNetwork{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &config.DomainManagerNetwork{
				Fqdn: utils.StringPtr("pc.example.test"),
				ExternalAddress: &commonConfig.IPAddress{
					Ipv4: &commonConfig.IPv4Address{
						Value:        utils.StringPtr("10.0.0.1"),
						PrefixLength: utils.IntPtr(32),
					},
				},
			},
			want: []map[string]interface{}{{
				"fqdn": "pc.example.test",
				"external_address": []map[string]interface{}{{
					"ipv4": []interface{}{map[string]interface{}{
						"value":         utils.StringPtr("10.0.0.1"),
						"prefix_length": utils.IntPtr(32),
					}},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPCNetwork(tt.in))
		})
	}
}

func TestUnitFlattenIPAddress(t *testing.T) {
	t.Parallel()

	val := utils.StringPtr("10.0.0.2")
	prefix := utils.IntPtr(24)

	tests := []struct {
		name string
		in   *commonConfig.IPAddress
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &commonConfig.IPAddress{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &commonConfig.IPAddress{
				Ipv4: &commonConfig.IPv4Address{Value: val, PrefixLength: prefix},
			},
			want: []map[string]interface{}{{
				"ipv4": []interface{}{map[string]interface{}{
					"value":         val,
					"prefix_length": prefix,
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPAddress(tt.in))
		})
	}
}

func TestUnitFlattenIPv4Address(t *testing.T) {
	t.Parallel()

	val := utils.StringPtr("192.168.1.1")
	prefix := utils.IntPtr(32)

	tests := []struct {
		name string
		in   *commonConfig.IPv4Address
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &commonConfig.IPv4Address{Value: val, PrefixLength: prefix},
			want: []interface{}{map[string]interface{}{"value": val, "prefix_length": prefix}},
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

	val := utils.StringPtr("2001:db8::1")
	prefix := utils.IntPtr(128)

	tests := []struct {
		name string
		in   *commonConfig.IPv6Address
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &commonConfig.IPv6Address{Value: val, PrefixLength: prefix},
			want: []interface{}{map[string]interface{}{"value": val, "prefix_length": prefix}},
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

	val := utils.StringPtr("ntp.example.test")
	tests := []struct {
		name string
		in   *commonConfig.FQDN
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &commonConfig.FQDN{Value: val},
			want: []interface{}{map[string]interface{}{"value": val}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFQDN(tt.in))
		})
	}
}

func TestUnitFlattenIPAddressOrFQDN(t *testing.T) {
	t.Parallel()

	val := utils.StringPtr("8.8.8.8")
	prefix := utils.IntPtr(32)

	tests := []struct {
		name string
		in   []commonConfig.IPAddressOrFQDN
		want interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []commonConfig.IPAddressOrFQDN{}, want: nil},
		{
			name: "populated",
			in: []commonConfig.IPAddressOrFQDN{{
				Ipv4: &commonConfig.IPv4Address{Value: val, PrefixLength: prefix},
			}},
			want: []map[string]interface{}{{
				"ipv4": []interface{}{map[string]interface{}{"value": val, "prefix_length": prefix}},
				"ipv6": []interface{}(nil),
				"fqdn": []interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPAddressOrFQDN(tt.in))
		})
	}
}

func TestUnitFlattenIPRanges(t *testing.T) {
	t.Parallel()

	beginVal := utils.StringPtr("10.0.0.10")
	endVal := utils.StringPtr("10.0.0.20")
	prefix := utils.IntPtr(32)

	tests := []struct {
		name string
		in   []commonConfig.IpRange
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []commonConfig.IpRange{}, want: nil},
		{
			name: "populated",
			in: []commonConfig.IpRange{{
				Begin: &commonConfig.IPAddress{Ipv4: &commonConfig.IPv4Address{Value: beginVal, PrefixLength: prefix}},
				End:   &commonConfig.IPAddress{Ipv4: &commonConfig.IPv4Address{Value: endVal, PrefixLength: prefix}},
			}},
			want: []map[string]interface{}{{
				"begin": []map[string]interface{}{{
					"ipv4": []interface{}{map[string]interface{}{"value": beginVal, "prefix_length": prefix}},
				}},
				"end": []map[string]interface{}{{
					"ipv4": []interface{}{map[string]interface{}{"value": endVal, "prefix_length": prefix}},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPRanges(tt.in))
		})
	}
}

func TestUnitFlattenExternalNetworks(t *testing.T) {
	t.Parallel()

	gwVal := utils.StringPtr("10.0.0.1")
	maskVal := utils.StringPtr("255.255.255.0")
	prefix := utils.IntPtr(32)
	netExtID := utils.StringPtr("net-1")

	gw := commonConfig.IPAddressOrFQDN{
		Ipv4: &commonConfig.IPv4Address{Value: gwVal, PrefixLength: prefix},
	}
	mask := commonConfig.IPAddressOrFQDN{
		Ipv4: &commonConfig.IPv4Address{Value: maskVal, PrefixLength: prefix},
	}

	tests := []struct {
		name string
		in   []config.ExternalNetwork
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []config.ExternalNetwork{}, want: nil},
		{
			name: "populated",
			in: []config.ExternalNetwork{{
				DefaultGateway: &gw,
				SubnetMask:     &mask,
				NetworkExtId:   netExtID,
			}},
			want: []map[string]interface{}{{
				"default_gateway": []map[string]interface{}{{
					"ipv4": []interface{}{map[string]interface{}{"value": gwVal, "prefix_length": prefix}},
					"ipv6": []interface{}(nil),
					"fqdn": []interface{}(nil),
				}},
				"subnet_mask": []map[string]interface{}{{
					"ipv4": []interface{}{map[string]interface{}{"value": maskVal, "prefix_length": prefix}},
					"ipv6": []interface{}(nil),
					"fqdn": []interface{}(nil),
				}},
				"ip_ranges":      []map[string]interface{}(nil),
				"network_ext_id": netExtID,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExternalNetworks(tt.in))
		})
	}
}

func TestUnitFlattenPcs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.DomainManager
		want []map[string]interface{}
	}{
		{name: "empty", in: []config.DomainManager{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.DomainManager{{
				ExtId:                          utils.StringPtr("pc-1"),
				TenantId:                       utils.StringPtr("tenant-1"),
				IsRegisteredWithHostingCluster: utils.BoolPtr(true),
				HostingClusterExtId:            utils.StringPtr("cluster-1"),
				ShouldEnableHighAvailability:   utils.BoolPtr(false),
				NodeExtIds:                     []string{"node-1"},
			}},
			want: []map[string]interface{}{{
				"ext_id":                             "pc-1",
				"tenant_id":                          "tenant-1",
				"links":                              []map[string]interface{}(nil),
				"config":                             []map[string]interface{}(nil),
				"is_registered_with_hosting_cluster": true,
				"network":                            []map[string]interface{}(nil),
				"hosting_cluster_ext_id":             "cluster-1",
				"should_enable_high_availability":    false,
				"node_ext_ids":                       []string{"node-1"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPcs(tt.in))
		})
	}
}

func TestUnitFlattenRestorablePcs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []management.RestorableDomainManager
		want []map[string]interface{}
	}{
		{name: "empty", in: []management.RestorableDomainManager{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []management.RestorableDomainManager{{
				ExtId:                          utils.StringPtr("rpc-1"),
				TenantId:                       utils.StringPtr("tenant-1"),
				IsRegisteredWithHostingCluster: utils.BoolPtr(false),
				HostingClusterExtId:            utils.StringPtr("cluster-2"),
				ShouldEnableHighAvailability:   utils.BoolPtr(true),
				NodeExtIds:                     []string{"node-2"},
			}},
			want: []map[string]interface{}{{
				"ext_id":                             "rpc-1",
				"tenant_id":                          "tenant-1",
				"links":                              []map[string]interface{}(nil),
				"config":                             []map[string]interface{}(nil),
				"is_registered_with_hosting_cluster": false,
				"network":                            []map[string]interface{}(nil),
				"hosting_cluster_ext_id":             "cluster-2",
				"should_enable_high_availability":    true,
				"node_ext_ids":                       []string{"node-2"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRestorablePcs(tt.in))
		})
	}
}

func TestUnitFlattenDomainManager(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.DomainManager
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &config.DomainManager{},
			want: []map[string]interface{}{{
				"tenant_id":                          "",
				"ext_id":                             "",
				"config":                             []map[string]interface{}(nil),
				"is_registered_with_hosting_cluster": false,
				"network":                            []map[string]interface{}(nil),
				"hosting_cluster_ext_id":             "",
				"should_enable_high_availability":    false,
				"node_ext_ids":                       []string(nil),
			}},
		},
		{
			name: "populated",
			in: &config.DomainManager{
				TenantId:                       utils.StringPtr("tenant-1"),
				ExtId:                          utils.StringPtr("dm-1"),
				IsRegisteredWithHostingCluster: utils.BoolPtr(true),
				HostingClusterExtId:            utils.StringPtr("hc-1"),
				ShouldEnableHighAvailability:   utils.BoolPtr(true),
				NodeExtIds:                     []string{"n-1"},
			},
			want: []map[string]interface{}{{
				"tenant_id":                          "tenant-1",
				"ext_id":                             "dm-1",
				"config":                             []map[string]interface{}(nil),
				"is_registered_with_hosting_cluster": true,
				"network":                            []map[string]interface{}(nil),
				"hosting_cluster_ext_id":             "hc-1",
				"should_enable_high_availability":    true,
				"node_ext_ids":                       []string{"n-1"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDomainManager(tt.in))
		})
	}
}

func TestUnitFlattenRestorePoints(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		in   []management.RestorePoint
		want []map[string]interface{}
	}{
		{name: "empty", in: []management.RestorePoint{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []management.RestorePoint{{
				TenantId:     utils.StringPtr("tenant-1"),
				ExtId:        utils.StringPtr("rp-1"),
				CreationTime: &when,
			}},
			want: []map[string]interface{}{{
				"tenant_id":      "tenant-1",
				"ext_id":         "rp-1",
				"links":          []map[string]interface{}(nil),
				"creation_time":  utils.StringPtr(when.String()),
				"domain_manager": []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRestorePoints(tt.in))
		})
	}
}

func TestUnitFlattenClusterReference(t *testing.T) {
	t.Parallel()

	extID := utils.StringPtr("cluster-1")
	name := utils.StringPtr("cluster-name")

	tests := []struct {
		name string
		in   *management.ClusterReference
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &management.ClusterReference{ExtId: extID, Name: name},
			want: []map[string]interface{}{{"ext_id": extID, "name": name}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenClusterReference(tt.in))
		})
	}
}

func TestUnitFlattenClusterLocation(t *testing.T) {
	t.Parallel()

	extID := utils.StringPtr("cluster-1")
	name := utils.StringPtr("cluster-name")
	ref := management.NewClusterReference()
	ref.ExtId = extID
	ref.Name = name
	loc := management.NewClusterLocation()
	cfg := management.NewOneOfClusterLocationConfig()
	require.NoError(t, cfg.SetValue(*ref))
	loc.Config = cfg

	tests := []struct {
		name string
		in   management.ClusterLocation
		want []map[string]interface{}
	}{
		{name: "empty", in: management.ClusterLocation{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in:   *loc,
			want: []map[string]interface{}{{
				"config": []map[string]interface{}{{"ext_id": extID, "name": name}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenClusterLocation(tt.in))
		})
	}
}

func TestUnitFlattenAccessKeyCredentials(t *testing.T) {
	t.Parallel()

	accessKey := utils.StringPtr("AKIA...")
	secret := utils.StringPtr("secret")

	tests := []struct {
		name string
		in   *management.AccessKeyCredentials
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &management.AccessKeyCredentials{
				AccessKeyId:     accessKey,
				SecretAccessKey: secret,
			},
			want: []map[string]interface{}{{
				"access_key_id":     accessKey,
				"secret_access_key": secret,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAccessKeyCredentials(tt.in))
		})
	}
}

func TestUnitFlattenBackupPolicy(t *testing.T) {
	t.Parallel()

	rpo := utils.IntPtr(60)
	tests := []struct {
		name string
		in   *management.BackupPolicy
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &management.BackupPolicy{RpoInMinutes: rpo},
			want: []map[string]interface{}{{"rpo_in_minutes": rpo}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenBackupPolicy(tt.in))
		})
	}
}

func TestUnitFlattenProviderConfig(t *testing.T) {
	t.Parallel()

	bucket := utils.StringPtr("bucket-1")
	region := utils.StringPtr("us-west-2")
	accessKey := utils.StringPtr("AKIA")
	secret := utils.StringPtr("secret")

	aws := management.NewAWSS3Config()
	aws.BucketName = bucket
	aws.Region = region
	aws.Credentials = &management.AccessKeyCredentials{
		AccessKeyId:     accessKey,
		SecretAccessKey: secret,
	}
	pc := management.NewOneOfObjectStoreLocationProviderConfig()
	require.NoError(t, pc.SetValue(*aws))

	tests := []struct {
		name string
		in   *management.OneOfObjectStoreLocationProviderConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   pc,
			want: []map[string]interface{}{{
				"bucket_name": bucket,
				"region":      region,
				"credentials": []map[string]interface{}{{
					"access_key_id":     accessKey,
					"secret_access_key": secret,
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenProviderConfig(tt.in))
		})
	}
}

func TestUnitFlattenObjectStoreLocation(t *testing.T) {
	t.Parallel()

	rpo := utils.IntPtr(30)
	bucket := utils.StringPtr("bucket-1")
	region := utils.StringPtr("us-east-1")

	aws := management.NewAWSS3Config()
	aws.BucketName = bucket
	aws.Region = region
	pc := management.NewOneOfObjectStoreLocationProviderConfig()
	require.NoError(t, pc.SetValue(*aws))

	obj := management.ObjectStoreLocation{
		ProviderConfig: pc,
		BackupPolicy:   &management.BackupPolicy{RpoInMinutes: rpo},
	}

	got := flattenObjectStoreLocation(obj)
	require.Len(t, got, 1)
	require.Equal(t, []map[string]interface{}{{
		"bucket_name": bucket,
		"region":      region,
		"credentials": []map[string]interface{}(nil),
	}}, got[0]["provider_config"])
	require.Equal(t, []map[string]interface{}{{"rpo_in_minutes": rpo}}, got[0]["backup_policy"])
}

func TestUnitFlattenBackupTargetLocation(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenBackupTargetLocation(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenBackupTargetLocation(management.NewOneOfBackupTargetLocation()))
	})

	t.Run("cluster location", func(t *testing.T) {
		t.Parallel()
		extID := utils.StringPtr("cluster-1")
		ref := management.NewClusterReference()
		ref.ExtId = extID
		loc := management.NewClusterLocation()
		cfg := management.NewOneOfClusterLocationConfig()
		require.NoError(t, cfg.SetValue(*ref))
		loc.Config = cfg
		oneOf := management.NewOneOfBackupTargetLocation()
		require.NoError(t, oneOf.SetValue(*loc))

		got := flattenBackupTargetLocation(oneOf)
		require.Len(t, got, 1)
		require.NotNil(t, got[0]["cluster_location"])
	})

	t.Run("object store location", func(t *testing.T) {
		t.Parallel()
		bucket := utils.StringPtr("bucket-1")
		aws := management.NewAWSS3Config()
		aws.BucketName = bucket
		pc := management.NewOneOfObjectStoreLocationProviderConfig()
		require.NoError(t, pc.SetValue(*aws))
		obj := management.NewObjectStoreLocation()
		obj.ProviderConfig = pc
		oneOf := management.NewOneOfBackupTargetLocation()
		require.NoError(t, oneOf.SetValue(*obj))

		got := flattenBackupTargetLocation(oneOf)
		require.Len(t, got, 1)
		require.NotNil(t, got[0]["object_store_location"])
	})
}

func TestUnitFlattenBackupTargets(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)
	extID := utils.StringPtr("bt-1")
	tenant := utils.StringPtr("tenant-1")
	paused := utils.BoolPtr(true)
	reason := utils.StringPtr("maintenance")

	tests := []struct {
		name string
		in   []management.BackupTarget
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []management.BackupTarget{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []management.BackupTarget{{
				ExtId:             extID,
				TenantId:          tenant,
				LastSyncTime:      &when,
				IsBackupPaused:    paused,
				BackupPauseReason: reason,
			}},
			want: []map[string]interface{}{{
				"ext_id":              extID,
				"tenant_id":           tenant,
				"links":               []map[string]interface{}(nil),
				"location":            []map[string]interface{}(nil),
				"last_sync_time":      utils.StringPtr(when.String()),
				"is_backup_paused":    paused,
				"backup_pause_reason": reason,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenBackupTargets(tt.in))
		})
	}
}

func TestUnitFlattenRestoreSourceClusterReference(t *testing.T) {
	t.Parallel()

	extID := utils.StringPtr("cluster-1")
	tests := []struct {
		name string
		in   *management.ClusterReference
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &management.ClusterReference{ExtId: extID, Name: utils.StringPtr("ignored")},
			want: []map[string]interface{}{{"ext_id": extID}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRestoreSourceClusterReference(tt.in))
		})
	}
}

func TestUnitFlattenRestoreSourceClusterLocation(t *testing.T) {
	t.Parallel()

	t.Run("nil config", func(t *testing.T) {
		t.Parallel()
		got := flattenRestoreSourceClusterLocation(management.ClusterLocation{})
		require.Equal(t, []map[string]interface{}{{"config": nil}}, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		extID := utils.StringPtr("cluster-1")
		ref := management.NewClusterReference()
		ref.ExtId = extID
		loc := management.NewClusterLocation()
		cfg := management.NewOneOfClusterLocationConfig()
		require.NoError(t, cfg.SetValue(*ref))
		loc.Config = cfg

		got := flattenRestoreSourceClusterLocation(*loc)
		require.Equal(t, []map[string]interface{}{{
			"config": []map[string]interface{}{{"ext_id": extID}},
		}}, got)
	})
}

func TestUnitFlattenRestoreSourceObjectStoreLocation(t *testing.T) {
	t.Parallel()

	bucket := utils.StringPtr("bucket-1")
	aws := management.NewAWSS3Config()
	aws.BucketName = bucket
	pc := management.NewOneOfObjectStoreLocationProviderConfig()
	require.NoError(t, pc.SetValue(*aws))

	got := flattenRestoreSourceObjectStoreLocation(management.ObjectStoreLocation{ProviderConfig: pc})
	require.Len(t, got, 1)
	// credentials are WriteOnly and must not be rehydrated into state.
	_, hasCreds := got[0]["provider_config"].([]map[string]interface{})[0]["credentials"]
	require.False(t, hasCreds)
}

func TestUnitFlattenRestoreSourceLocation(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenRestoreSourceLocation(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenRestoreSourceLocation(management.NewOneOfRestoreSourceLocation()))
	})

	t.Run("cluster location", func(t *testing.T) {
		t.Parallel()
		extID := utils.StringPtr("cluster-1")
		ref := management.NewClusterReference()
		ref.ExtId = extID
		loc := management.NewClusterLocation()
		cfg := management.NewOneOfClusterLocationConfig()
		require.NoError(t, cfg.SetValue(*ref))
		loc.Config = cfg
		oneOf := management.NewOneOfRestoreSourceLocation()
		require.NoError(t, oneOf.SetValue(*loc))

		got := flattenRestoreSourceLocation(oneOf)
		require.Len(t, got, 1)
		require.NotNil(t, got[0]["cluster_location"])
	})
}

func TestUnitFlattenResourceBackupTargetLocation(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenResourceBackupTargetLocation(nil, nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenResourceBackupTargetLocation(management.NewOneOfBackupTargetLocation(), nil))
	})

	t.Run("cluster location", func(t *testing.T) {
		t.Parallel()
		extID := utils.StringPtr("cluster-1")
		ref := management.NewClusterReference()
		ref.ExtId = extID
		loc := management.NewClusterLocation()
		cfg := management.NewOneOfClusterLocationConfig()
		require.NoError(t, cfg.SetValue(*ref))
		loc.Config = cfg
		oneOf := management.NewOneOfBackupTargetLocation()
		require.NoError(t, oneOf.SetValue(*loc))

		got := flattenResourceBackupTargetLocation(oneOf, nil)
		require.Len(t, got, 1)
		require.NotNil(t, got[0]["cluster_location"])
	})
}

func TestUnitFlattenResourceObjectStoreLocation(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenResourceObjectStoreLocation(nil, nil))
	})

	t.Run("populated without state", func(t *testing.T) {
		t.Parallel()
		bucket := utils.StringPtr("bucket-1")
		aws := management.NewAWSS3Config()
		aws.BucketName = bucket
		pc := management.NewOneOfObjectStoreLocationProviderConfig()
		require.NoError(t, pc.SetValue(*aws))
		obj := &management.ObjectStoreLocation{ProviderConfig: pc}

		got := flattenResourceObjectStoreLocation(obj, nil)
		require.Len(t, got, 1)
		require.NotNil(t, got[0]["provider_config"])
	})

	t.Run("does not rehydrate write-only credentials from state", func(t *testing.T) {
		t.Parallel()
		bucket := utils.StringPtr("bucket-1")
		aws := management.NewAWSS3Config()
		aws.BucketName = bucket
		pc := management.NewOneOfObjectStoreLocationProviderConfig()
		require.NoError(t, pc.SetValue(*aws))
		obj := &management.ObjectStoreLocation{ProviderConfig: pc}

		stateCreds := []interface{}{map[string]interface{}{
			"access_key_id":     "from-state",
			"secret_access_key": "state-secret",
		}}
		state := []interface{}{map[string]interface{}{
			"object_store_location": []interface{}{map[string]interface{}{
				"provider_config": []interface{}{map[string]interface{}{
					"credentials": stateCreds,
				}},
			}},
		}}

		got := flattenResourceObjectStoreLocation(obj, state)
		_, hasCreds := got[0]["provider_config"].([]map[string]interface{})[0]["credentials"]
		require.False(t, hasCreds)
	})
}

func TestUnitUnregisterDomainManagerUUIDFromTask(t *testing.T) {
	t.Parallel()

	extID := "pc-ext-1"
	other := "other-pc"

	tests := []struct {
		name    string
		task    config.Task
		want    *string
		wantErr bool
	}{
		{
			name:    "empty",
			task:    config.Task{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "match config rel",
			task: config.Task{EntitiesAffected: []config.EntityReference{{
				ExtId: utils.StringPtr(extID),
				Rel:   utils.StringPtr(utils.RelEntityTypeDomainManager),
			}}},
			want: utils.StringPtr(extID),
		},
		{
			name: "match management rel",
			task: config.Task{EntitiesAffected: []config.EntityReference{{
				ExtId: utils.StringPtr(extID),
				Rel:   utils.StringPtr(utils.RelEntityTypeDomainManagerManagement),
			}}},
			want: utils.StringPtr(extID),
		},
		{
			name: "wrong ext id",
			task: config.Task{EntitiesAffected: []config.EntityReference{{
				ExtId: utils.StringPtr(other),
				Rel:   utils.StringPtr(utils.RelEntityTypeDomainManager),
			}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := unregisterDomainManagerUUIDFromTask(tt.task, extID)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitNoopResourceHandlers(t *testing.T) {
	t.Parallel()

	restorePc := schema.TestResourceDataRaw(t, ResourceNutanixRestorePcV2().Schema, map[string]interface{}{})
	require.Nil(t, ResourceNutanixRestorePcRead(t.Context(), restorePc, nil))
	require.Nil(t, ResourceNutanixRestorePcUpdate(t.Context(), restorePc, nil))
	require.Nil(t, ResourceNutanixRestorePcDelete(t.Context(), restorePc, nil))

	unregister := schema.TestResourceDataRaw(t, ResourceNutanixUnregisterClusterV2().Schema, map[string]interface{}{
		"ext_id": "pc-1",
	})
	require.Nil(t, ResourceNutanixUnregisterClusterV2Read(t.Context(), unregister, nil))
	require.Nil(t, ResourceNutanixUnregisterClusterV2Update(t.Context(), unregister, nil))
	require.Nil(t, ResourceNutanixUnregisterClusterV2Delete(t.Context(), unregister, nil))

	deploy := schema.TestResourceDataRaw(t, ResourceNutanixDeployPcV2().Schema, map[string]interface{}{})
	require.Nil(t, ResourceNutanixDeployPcV2Read(t.Context(), deploy, nil))
	require.Nil(t, ResourceNutanixDeployPcV2Delete(t.Context(), deploy, nil))

	reg := schema.TestResourceDataRaw(t, ResourceNutanixClusterPCRegistrationV2().Schema, map[string]interface{}{
		"pc_ext_id": "pc-1",
	})
	require.Nil(t, ResourceNutanixClusterPCRegistrationV2Update(t.Context(), reg, nil))
	require.Nil(t, ResourceNutanixClusterPCRegistrationV2Delete(t.Context(), reg, nil))
}
