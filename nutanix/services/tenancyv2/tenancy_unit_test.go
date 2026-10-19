package tenancyv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/common/v1/response"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/iam/v4/authz"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func testNatConfig(transitVPC string, floatingIPCount int, overlaySubnet string) *import1.NatConfig {
	nat := &import1.NatConfig{}
	if transitVPC != "" {
		nat.TransitVpcExtId = utils.StringPtr(transitVPC)
	}
	if floatingIPCount != 0 {
		nat.FloatingIpCount = utils.IntPtr(floatingIPCount)
	}
	if overlaySubnet != "" {
		nat.OverlaySubnetExtId = utils.StringPtr(overlaySubnet)
	}
	return nat
}

func TestUnitTenancySchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixTenantV2().InternalValidate(ResourceNutanixTenantV2().Schema, true))
	require.NoError(t, ResourceNutanixAddDomainV2().InternalValidate(ResourceNutanixAddDomainV2().Schema, true))
	require.NoError(t, ResourceNutanixTenantDomainAllocationV2().InternalValidate(ResourceNutanixTenantDomainAllocationV2().Schema, true))
	require.NoError(t, ResourceNutanixTenantMembershipV2().InternalValidate(ResourceNutanixTenantMembershipV2().Schema, true))

	require.NoError(t, DatasourceNutanixTenantV2().InternalValidate(DatasourceNutanixTenantV2().Schema, false))
	require.NoError(t, DatasourceNutanixTenantsV2().InternalValidate(DatasourceNutanixTenantsV2().Schema, false))
	require.NoError(t, DatasourceNutanixDomainsV2().InternalValidate(DatasourceNutanixDomainsV2().Schema, false))
	require.NoError(t, DatasourceNutanixTenantDomainAllocationV2().InternalValidate(DatasourceNutanixTenantDomainAllocationV2().Schema, false))
	require.NoError(t, DatasourceNutanixTenantDomainAllocationsV2().InternalValidate(DatasourceNutanixTenantDomainAllocationsV2().Schema, false))
	require.NoError(t, DataSourceNutanixTenantMembershipV2().InternalValidate(DataSourceNutanixTenantMembershipV2().Schema, false))
	require.NoError(t, DataSourceNutanixTenantMembershipsV2().InternalValidate(DataSourceNutanixTenantMembershipsV2().Schema, false))
}

func TestUnitAddDomainDeleteNoop(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixAddDomainV2().Schema, map[string]interface{}{
		"domain_ext_id": "domain-1",
	})
	require.Nil(t, ResourceNutanixAddDomainV2Delete(t.Context(), d, nil))
}

func TestUnitExpandLocalUserInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *import1.LocalUserInfo
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"email":      "a@example.com",
				"first_name": "Ada",
				"last_name":  "Lovelace",
				"password":   "secret",
				"username":   "ada",
			}},
			want: &import1.LocalUserInfo{
				Email:     utils.StringPtr("a@example.com"),
				FirstName: utils.StringPtr("Ada"),
				LastName:  utils.StringPtr("Lovelace"),
				Password:  utils.StringPtr("secret"),
				Username:  utils.StringPtr("ada"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandLocalUserInfo(tt.in))
		})
	}
}

func TestUnitFlattenLocalUserInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.LocalUserInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import1.LocalUserInfo{
				Email:     utils.StringPtr("a@example.com"),
				FirstName: utils.StringPtr("Ada"),
				LastName:  utils.StringPtr("Lovelace"),
				Username:  utils.StringPtr("ada"),
				Password:  utils.StringPtr("secret"),
			},
			want: []map[string]interface{}{{
				"email":      "a@example.com",
				"first_name": "Ada",
				"last_name":  "Lovelace",
				"username":   "ada",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocalUserInfo(tt.in))
		})
	}
}

func TestUnitExpandQuota(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *import1.Quota
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"memory":  1024,
				"storage": 2048,
				"vcpu":    4,
			}},
			want: &import1.Quota{
				Memory:  utils.Int64Ptr(1024),
				Storage: utils.Int64Ptr(2048),
				Vcpu:    utils.IntPtr(4),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandQuota(tt.in))
		})
	}
}

func TestUnitFlattenQuota(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.Quota
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import1.Quota{
				ExtId:   utils.StringPtr("q-1"),
				Memory:  utils.Int64Ptr(1024),
				Storage: utils.Int64Ptr(2048),
				Vcpu:    utils.IntPtr(4),
			},
			want: []map[string]interface{}{{
				"ext_id":  "q-1",
				"memory":  int64(1024),
				"storage": int64(2048),
				"vcpu":    4,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenQuota(tt.in))
		})
	}
}

func TestUnitFlattenTenantStateInfo(t *testing.T) {
	t.Parallel()

	state := import1.TENANTSTATE_ACTIVE
	tests := []struct {
		name string
		in   *import1.TenantStateInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: &import1.TenantStateInfo{}, want: []map[string]interface{}{{}}},
		{
			name: "populated",
			in:   &import1.TenantStateInfo{State: &state},
			want: []map[string]interface{}{{"state": "ACTIVE"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTenantStateInfo(tt.in))
		})
	}
}

func TestUnitFlattenDomainAllocationStateInfo(t *testing.T) {
	t.Parallel()

	state := import1.DOMAINENTITYSTATE_ACTIVE
	tests := []struct {
		name string
		in   *import1.DomainAllocationStateInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: &import1.DomainAllocationStateInfo{}, want: []map[string]interface{}{{}}},
		{
			name: "populated",
			in:   &import1.DomainAllocationStateInfo{State: &state},
			want: []map[string]interface{}{{"state": "ACTIVE"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDomainAllocationStateInfo(tt.in))
		})
	}
}

func TestUnitFlattenExternalNetworkConfigStateInfo(t *testing.T) {
	t.Parallel()

	state := import1.EXTERNALNETWORKCONFIGSTATE_CONFIGURED
	tests := []struct {
		name string
		in   *import1.ExternalNetworkConfigStateInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: &import1.ExternalNetworkConfigStateInfo{}, want: []map[string]interface{}{{}}},
		{
			name: "populated",
			in:   &import1.ExternalNetworkConfigStateInfo{State: &state},
			want: []map[string]interface{}{{"state": "CONFIGURED"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExternalNetworkConfigStateInfo(tt.in))
		})
	}
}

func TestUnitExpandProjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []import1.Project
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"name": "proj-1",
				"quota": []interface{}{map[string]interface{}{
					"memory":  10,
					"storage": 20,
					"vcpu":    2,
				}},
			}},
			want: []import1.Project{{
				Name: utils.StringPtr("proj-1"),
				Quota: &import1.Quota{
					Memory:  utils.Int64Ptr(10),
					Storage: utils.Int64Ptr(20),
					Vcpu:    utils.IntPtr(2),
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandProjects(tt.in))
		})
	}
}

func TestUnitFlattenProjects(t *testing.T) {
	t.Parallel()

	state := import1.DOMAINENTITYSTATE_PENDING
	tests := []struct {
		name string
		in   []import1.Project
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.Project{}, want: nil},
		{
			name: "populated",
			in: []import1.Project{{
				ExtId: utils.StringPtr("p-1"),
				Name:  utils.StringPtr("proj-1"),
				State: &state,
				Quota: &import1.Quota{Vcpu: utils.IntPtr(2)},
			}},
			want: []map[string]interface{}{{
				"ext_id":         "p-1",
				"name":           "proj-1",
				"state":          "PENDING",
				"quota":          []map[string]interface{}{{"ext_id": "", "memory": int64(0), "storage": int64(0), "vcpu": 2}},
				"resource_group": []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenProjects(tt.in))
		})
	}
}

func TestUnitExpandResourceGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *import1.ResourceGroup
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"placement_targets": []interface{}{map[string]interface{}{
					"cluster_uuid":            "c-1",
					"storage_container_uuids": schema.NewSet(schema.HashString, []interface{}{"sc-1"}),
				}},
			}},
			want: &import1.ResourceGroup{
				PlacementTargets: []import1.PlacementTarget{{
					ClusterUuid:           utils.StringPtr("c-1"),
					StorageContainerUuids: []string{"sc-1"},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandResourceGroup(tt.in))
		})
	}
}

func TestUnitFlattenResourceGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.ResourceGroup
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import1.ResourceGroup{
				Uuid: utils.StringPtr("rg-1"),
				PlacementTargets: []import1.PlacementTarget{{
					ClusterUuid:           utils.StringPtr("c-1"),
					StorageContainerUuids: []string{"sc-1"},
				}},
			},
			want: []map[string]interface{}{{
				"uuid": "rg-1",
				"placement_targets": []map[string]interface{}{{
					"cluster_uuid":            "c-1",
					"storage_container_uuids": []string{"sc-1"},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceGroup(tt.in))
		})
	}
}

func TestUnitExpandPlacementTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []import1.PlacementTarget
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"cluster_uuid":            "c-1",
				"storage_container_uuids": schema.NewSet(schema.HashString, []interface{}{"sc-1", "sc-2"}),
			}},
			want: []import1.PlacementTarget{{
				ClusterUuid:           utils.StringPtr("c-1"),
				StorageContainerUuids: []string{"sc-1", "sc-2"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandPlacementTargets(tt.in)
			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.Equal(t, tt.want[0].ClusterUuid, got[0].ClusterUuid)
			require.ElementsMatch(t, tt.want[0].StorageContainerUuids, got[0].StorageContainerUuids)
		})
	}
}

func TestUnitFlattenPlacementTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.PlacementTarget
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.PlacementTarget{}, want: nil},
		{
			name: "populated",
			in: []import1.PlacementTarget{{
				ClusterUuid:           utils.StringPtr("c-1"),
				StorageContainerUuids: []string{"sc-1"},
			}},
			want: []map[string]interface{}{{
				"cluster_uuid":            "c-1",
				"storage_container_uuids": []string{"sc-1"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPlacementTargets(tt.in))
		})
	}
}

func TestUnitExpandExternalNetworkConnection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *import1.ExternalNetworkConnection
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"config": []interface{}{map[string]interface{}{
					"nat": []interface{}{map[string]interface{}{
						"transit_vpc_ext_id":    "vpc-1",
						"floating_ip_count":     2,
						"overlay_subnet_ext_id": "subnet-1",
					}},
				}},
			}},
			want: &import1.ExternalNetworkConnection{
				Config: &import1.ExternalNetworkConnConfig{
					Nat: testNatConfig("vpc-1", 2, "subnet-1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandExternalNetworkConnection(tt.in))
		})
	}
}

func TestUnitFlattenExternalNetworkConnection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.ExternalNetworkConnection
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import1.ExternalNetworkConnection{
				ExtId:    utils.StringPtr("enc-1"),
				TenantId: utils.StringPtr("t-1"),
				Config: &import1.ExternalNetworkConnConfig{
					Nat: testNatConfig("vpc-1", 1, ""),
				},
			},
			want: []map[string]interface{}{{
				"config": []map[string]interface{}{{
					"nat": []map[string]interface{}{{
						"transit_vpc_ext_id":    "vpc-1",
						"floating_ip_count":     1,
						"overlay_subnet_ext_id": "",
					}},
					"no_nat": []map[string]interface{}(nil),
				}},
				"ext_id":     "enc-1",
				"tenant_id":  "t-1",
				"state_info": []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExternalNetworkConnection(tt.in))
		})
	}
}

func TestUnitExpandExternalNetworkConnConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *import1.ExternalNetworkConnConfig
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"no_nat": []interface{}{map[string]interface{}{
					"overlay_subnet_ext_id":        "subnet-1",
					"transit_vpc_ext_id":           "vpc-1",
					"externally_routable_prefixes": []interface{}{"10.0.0.0/24"},
				}},
			}},
			want: &import1.ExternalNetworkConnConfig{
				NoNat: &import1.NoNatConfig{
					OverlaySubnetExtId:         utils.StringPtr("subnet-1"),
					TransitVpcExtId:            utils.StringPtr("vpc-1"),
					ExternallyRoutablePrefixes: []string{"10.0.0.0/24"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandExternalNetworkConnConfig(tt.in))
		})
	}
}

func TestUnitFlattenExternalNetworkConnConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.ExternalNetworkConnConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import1.ExternalNetworkConnConfig{
				Nat: testNatConfig("", 3, ""),
			},
			want: []map[string]interface{}{{
				"nat": []map[string]interface{}{{
					"transit_vpc_ext_id":    "",
					"floating_ip_count":     3,
					"overlay_subnet_ext_id": "",
				}},
				"no_nat": []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExternalNetworkConnConfig(tt.in))
		})
	}
}

func TestUnitExpandNatConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *import1.NatConfig
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"transit_vpc_ext_id":    "vpc-1",
				"floating_ip_count":     5,
				"overlay_subnet_ext_id": "subnet-1",
			}},
			want: testNatConfig("vpc-1", 5, "subnet-1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandNatConfig(tt.in))
		})
	}
}

func TestUnitFlattenNatConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.NatConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   testNatConfig("vpc-1", 5, "subnet-1"),
			want: []map[string]interface{}{{
				"transit_vpc_ext_id":    "vpc-1",
				"floating_ip_count":     5,
				"overlay_subnet_ext_id": "subnet-1",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNatConfig(tt.in))
		})
	}
}

func TestUnitExpandNoNatConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *import1.NoNatConfig
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"overlay_subnet_ext_id":        "subnet-1",
				"transit_vpc_ext_id":           "vpc-1",
				"externally_routable_prefixes": []interface{}{"10.1.0.0/16"},
			}},
			want: &import1.NoNatConfig{
				OverlaySubnetExtId:         utils.StringPtr("subnet-1"),
				TransitVpcExtId:            utils.StringPtr("vpc-1"),
				ExternallyRoutablePrefixes: []string{"10.1.0.0/16"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandNoNatConfig(tt.in))
		})
	}
}

func TestUnitFlattenNoNatConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.NoNatConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import1.NoNatConfig{
				OverlaySubnetExtId:         utils.StringPtr("subnet-1"),
				TransitVpcExtId:            utils.StringPtr("vpc-1"),
				ExternallyRoutablePrefixes: []string{"10.1.0.0/16"},
			},
			want: []map[string]interface{}{{
				"overlay_subnet_ext_id":        "subnet-1",
				"transit_vpc_ext_id":           "vpc-1",
				"externally_routable_prefixes": []string{"10.1.0.0/16"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNoNatConfig(tt.in))
		})
	}
}

func TestUnitExpandExternalUserDetails(t *testing.T) {
	t.Parallel()

	userType := import5.RMIDENTITYTYPE_USER
	tests := []struct {
		name string
		in   []interface{}
		want *import1.ExternalUserDetails
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"identity_provider_ext_id": "idp-1",
				"principal_value":          "user@example.com",
				"type":                     "USER",
			}},
			want: &import1.ExternalUserDetails{
				IdentityProviderExtId: utils.StringPtr("idp-1"),
				PrincipalValue:        utils.StringPtr("user@example.com"),
				Type:                  &userType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandExternalUserDetails(tt.in))
		})
	}
}

func TestUnitExpandRmIdentityType(t *testing.T) {
	t.Parallel()

	user := import5.RMIDENTITYTYPE_USER
	group := import5.RMIDENTITYTYPE_GROUP
	unknown := import5.RMIDENTITYTYPE_UNKNOWN
	tests := []struct {
		name string
		in   string
		want *import5.RmIdentityType
	}{
		{name: "empty", in: "", want: nil},
		{name: "unknown", in: "SOMETHING", want: &unknown},
		{name: "user", in: "USER", want: &user},
		{name: "group", in: "GROUP", want: &group},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, common.ExpandEnum[import5.RmIdentityType](tt.in))
		})
	}
}

func TestUnitFlattenExternalUserDetails(t *testing.T) {
	t.Parallel()

	userType := import5.RMIDENTITYTYPE_GROUP
	tests := []struct {
		name string
		in   *import1.ExternalUserDetails
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &import1.ExternalUserDetails{
				IdentityProviderExtId: utils.StringPtr("idp-1"),
				PrincipalValue:        utils.StringPtr("group@example.com"),
				Type:                  &userType,
			},
			want: []map[string]interface{}{{
				"identity_provider_ext_id": "idp-1",
				"principal_value":          "group@example.com",
				"type":                     "GROUP",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExternalUserDetails(tt.in))
		})
	}
}

func TestUnitFlattenMembershipStateInfo(t *testing.T) {
	t.Parallel()

	state := import1.MEMBERSHIPSTATE_ACTIVE
	tests := []struct {
		name string
		in   *import1.MembershipStateInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: &import1.MembershipStateInfo{}, want: []map[string]interface{}{{"state": ""}}},
		{
			name: "populated",
			in:   &import1.MembershipStateInfo{State: &state},
			want: []map[string]interface{}{{"state": "ACTIVE"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenMembershipStateInfo(tt.in))
		})
	}
}

func TestUnitFlattenAPILinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import2.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.ApiLink{}, want: nil},
		{
			name: "populated",
			in: []import2.ApiLink{{
				Href: utils.StringPtr("https://example.test/1"),
				Rel:  utils.StringPtr("self"),
			}},
			want: []map[string]interface{}{{
				"href": "https://example.test/1",
				"rel":  "self",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAPILinks(tt.in))
		})
	}
}

func TestUnitFlattenDomains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.Domain
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.Domain{}, want: nil},
		{
			name: "populated",
			in: []import1.Domain{{
				DomainExtId: utils.StringPtr("d-1"),
				TenantId:    utils.StringPtr("t-1"),
			}},
			want: []map[string]interface{}{{
				"domain_ext_id": "d-1",
				"ext_id":        "d-1",
				"tenant_id":     "t-1",
				"links":         []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDomains(tt.in))
		})
	}
}

func TestUnitFlattenDomainAllocationsEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.DomainAllocation
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.DomainAllocation{}, want: nil},
		{
			name: "populated",
			in: []import1.DomainAllocation{{
				ExtId:       utils.StringPtr("da-1"),
				DomainExtId: utils.StringPtr("d-1"),
				TenantId:    utils.StringPtr("t-1"),
			}},
			want: []map[string]interface{}{{
				"ext_id":                  "da-1",
				"domain_ext_id":           "d-1",
				"external_network_config": []map[string]interface{}(nil),
				"projects":                []map[string]interface{}(nil),
				"state_info":              []map[string]interface{}(nil),
				"tenant_id":               "t-1",
				"links":                   []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDomainAllocationsEntities(tt.in))
		})
	}
}

func TestUnitFlattenTenantsEntities(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name string
		in   []import1.Tenant
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.Tenant{}, want: nil},
		{
			name: "populated",
			in: []import1.Tenant{{
				ExtId:       utils.StringPtr("tenant-1"),
				Name:        utils.StringPtr("Acme"),
				Description: utils.StringPtr("desc"),
				Url:         utils.StringPtr("https://acme.test"),
				AddedBy:     utils.StringPtr("admin"),
				Domains:     []string{"d-1"},
				CreatedTime: &when,
				TenantId:    utils.StringPtr("t-1"),
			}},
			want: []map[string]interface{}{{
				"ext_id":           "tenant-1",
				"name":             "Acme",
				"description":      "desc",
				"url":              "https://acme.test",
				"added_by":         "admin",
				"domains":          []string{"d-1"},
				"admin_info":       []map[string]interface{}(nil),
				"aggregated_quota": []map[string]interface{}(nil),
				"created_time":     when.String(),
				"state_info":       []map[string]interface{}(nil),
				"tenant_id":        "t-1",
				"links":            []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTenantsEntities(tt.in))
		})
	}
}

func TestUnitFlattenDomainMemberships(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC)
	tests := []struct {
		name string
		in   []import1.Membership
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []import1.Membership{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []import1.Membership{{
				ExtId:         utils.StringPtr("m-1"),
				DomainExtId:   utils.StringPtr("d-1"),
				RoleExtId:     utils.StringPtr("r-1"),
				IdentityExtId: utils.StringPtr("i-1"),
				CreatedBy:     utils.StringPtr("admin"),
				Principal:     utils.StringPtr("user@example.com"),
				TenantId:      utils.StringPtr("t-1"),
				CreatedTime:   &when,
			}},
			want: []map[string]interface{}{{
				"ext_id":           "m-1",
				"domain_ext_id":    "d-1",
				"role_ext_id":      "r-1",
				"identity_ext_id":  "i-1",
				"identity_details": []map[string]interface{}(nil),
				"created_by":       "admin",
				"principal":        "user@example.com",
				"tenant_id":        "t-1",
				"links":            []map[string]interface{}(nil),
				"state_info":       []map[string]interface{}(nil),
				"created_time":     when.Format(time.RFC3339),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTenantMemberships(tt.in))
		})
	}
}

func TestUnitFlattenTenantMembership(t *testing.T) {
	t.Parallel()

	got := flattenTenantMembership(import1.Membership{
		ExtId: utils.StringPtr("m-1"),
	})
	require.Equal(t, "m-1", got["ext_id"])
	require.Nil(t, got["identity_details"])
}
