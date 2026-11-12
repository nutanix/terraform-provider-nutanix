package prism

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitPrismSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, DataSourceNutanixCategoryKey().InternalValidate(DataSourceNutanixCategoryKey().Schema, false))
	require.NoError(t, DataSourceNutanixProject().InternalValidate(DataSourceNutanixProject().Schema, false))
	require.NoError(t, DataSourceNutanixProjects().InternalValidate(DataSourceNutanixProjects().Schema, false))
	require.NoError(t, DataSourceNutanixProtectionRule().InternalValidate(DataSourceNutanixProtectionRule().Schema, false))
	require.NoError(t, DataSourceNutanixProtectionRules().InternalValidate(DataSourceNutanixProtectionRules().Schema, false))
	require.NoError(t, DataSourceNutanixRecoveryPlan().InternalValidate(DataSourceNutanixRecoveryPlan().Schema, false))
	require.NoError(t, DataSourceNutanixRecoveryPlans().InternalValidate(DataSourceNutanixRecoveryPlans().Schema, false))
	require.NoError(t, ResourceNutanixAccessControlPolicy().InternalValidate(ResourceNutanixAccessControlPolicy().Schema, true))
	require.NoError(t, ResourceNutanixCategoryKey().InternalValidate(ResourceNutanixCategoryKey().Schema, true))
	require.NoError(t, ResourceNutanixCategoryValue().InternalValidate(ResourceNutanixCategoryValue().Schema, true))
	require.NoError(t, ResourceNutanixProject().InternalValidate(ResourceNutanixProject().Schema, true))
	require.NoError(t, ResourceNutanixProtectionRule().InternalValidate(ResourceNutanixProtectionRule().Schema, true))
	require.NoError(t, ResourceNutanixRecoveryPlan().InternalValidate(ResourceNutanixRecoveryPlan().Schema, true))
	require.NoError(t, categoriesSchema().Elem.(*schema.Resource).InternalValidate(categoriesSchema().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, categoriesSchemaOptional().Elem.(*schema.Resource).InternalValidate(categoriesSchemaOptional().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, DataSourceFiltersSchema().Elem.(*schema.Resource).InternalValidate(DataSourceFiltersSchema().Elem.(*schema.Resource).Schema, true))
}

func TestUnitExpandStringList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []interface{}
		want []*string
	}{
		{name: "nil", in: nil, want: []*string{}},
		{name: "empty", in: []interface{}{}, want: []*string{}},
		{name: "populated", in: []interface{}{"a", "", "b"}, want: []*string{utils.StringPtr("a"), utils.StringPtr("b")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandStringList(tt.in))
		})
	}
}

func TestUnitExpandFlattenCategories(t *testing.T) {
	t.Parallel()
	t.Run("flatten nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenCategories(nil))
	})
	t.Run("flatten empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenCategories(map[string]string{}))
	})
	t.Run("flatten populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCategories(map[string]string{"env": "prod"})
		require.Equal(t, []interface{}{map[string]interface{}{"name": "env", "value": "prod"}}, got)
	})
	t.Run("expand empty set", func(t *testing.T) {
		t.Parallel()
		s := schema.NewSet(func(v interface{}) int {
			c := v.(map[string]interface{})
			return utils.HashcodeString(c["name"].(string) + c["value"].(string))
		}, nil)
		require.Equal(t, map[string]string{}, expandCategories(s))
	})
	t.Run("expand populated", func(t *testing.T) {
		t.Parallel()
		s := schema.NewSet(func(v interface{}) int {
			c := v.(map[string]interface{})
			return utils.HashcodeString(c["name"].(string) + c["value"].(string))
		}, []interface{}{map[string]interface{}{"name": "env", "value": "prod"}})
		require.Equal(t, map[string]string{"env": "prod"}, expandCategories(s))
	})
}

func TestUnitExpandFilterParams(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   map[string][]string
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: map[string][]string{}, want: []map[string]interface{}{}},
		{name: "populated", in: map[string][]string{"name": {"a"}}, want: []map[string]interface{}{{"name": "name", "values": []string{"a"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandFilterParams(tt.in))
		})
	}
}

func TestUnitBuildAndReplaceFilters(t *testing.T) {
	t.Parallel()
	t.Run("build empty", func(t *testing.T) {
		t.Parallel()
		s := schema.NewSet(schema.HashResource(DataSourceFiltersSchema().Elem.(*schema.Resource)), nil)
		require.Empty(t, BuildFiltersDataSource(s))
	})
	t.Run("build populated", func(t *testing.T) {
		t.Parallel()
		s := schema.NewSet(schema.HashResource(DataSourceFiltersSchema().Elem.(*schema.Resource)), []interface{}{
			map[string]interface{}{"name": "vm_name", "values": []interface{}{"vm1"}},
		})
		got := BuildFiltersDataSource(s)
		require.Len(t, got, 1)
		require.Equal(t, "vm_name", got[0].Name)
		require.Equal(t, []string{"vm1"}, got[0].Values)
	})
	t.Run("replace nil mappings", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, ReplaceFilterPrefixes(nil, nil))
	})
	t.Run("replace empty", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, ReplaceFilterPrefixes([]*client.AdditionalFilter{}, map[string]string{}))
	})
	t.Run("replace populated", func(t *testing.T) {
		t.Parallel()
		in := []*client.AdditionalFilter{{Name: "vm.name", Values: []string{"x"}}}
		got := ReplaceFilterPrefixes(in, map[string]string{"vm": "virtual_machine"})
		require.Equal(t, "virtual_machine.name", got[0].Name)
	})
	t.Run("filterParamsHash", func(t *testing.T) {
		t.Parallel()
		require.NotEqual(t, 0, filterParamsHash(map[string]interface{}{"name": "a"}))
	})
}

func TestUnitGetACPRoleUUIDAndRemoval(t *testing.T) {
	t.Parallel()
	t.Run("nil-ish", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", getACPRoleUUID(nil))
		require.Equal(t, "", getACPRoleUUID("x"))
		require.False(t, acpRemovalFromMiddle(nil, nil))
		require.False(t, acpRemovalFromMiddleRoles(nil, nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", getACPRoleUUID(map[string]interface{}{}))
		require.False(t, acpRemovalFromMiddle([]interface{}{}, []interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		item := map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "r1"}}}
		require.Equal(t, "r1", getACPRoleUUID(item))
		old := []interface{}{
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "a"}}},
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "b"}}},
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "c"}}},
		}
		newL := []interface{}{
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "a"}}},
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "c"}}},
		}
		require.True(t, acpRemovalFromMiddle(old, newL))
		require.Equal(t, 0, acpHash("bad"))
		require.NotEqual(t, 0, acpHash(item))
	})
}

func TestUnitValidateAndExpandReference(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, map[string]interface{}{}, flattenReferenceValues(nil))
		require.Nil(t, expandReference(map[string]interface{}{}))
		require.Nil(t, validateRef(map[string]interface{}{}))
		require.Equal(t, map[string]interface{}{}, flattenReference(nil))
		require.Equal(t, map[string]interface{}{}, flattenExternalNetworkListReference(nil))
		require.Equal(t, []interface{}{}, flattenReferenceValuesList(nil))
	})
	t.Run("empty slices", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenArrayReferenceValues(nil))
		require.Equal(t, []map[string]interface{}{}, flattenArrayReferenceValues([]*v3.Reference{}))
		require.Equal(t, []map[string]interface{}{}, flattenReferenceList(nil))
		require.Equal(t, []map[string]interface{}{}, flattenReferenceList([]*v3.ReferenceValues{}))
		require.Equal(t, []map[string]interface{}{}, flattenExternalNetworkListReferenceList(nil))
		require.Nil(t, expandDirectoryUserGroup(nil))
		require.Nil(t, expandDirectoryUserGroup([]interface{}{}))
		require.Nil(t, expandSamlUserGroup(nil))
		require.Nil(t, expandSamlUserGroup([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		ref := &v3.Reference{Kind: utils.StringPtr("vm"), UUID: utils.StringPtr("u1"), Name: utils.StringPtr("n1")}
		require.Equal(t, map[string]interface{}{"kind": "vm", "uuid": "u1", "name": "n1"}, flattenReferenceValues(ref))
		require.Equal(t, []interface{}{map[string]interface{}{"kind": "vm", "uuid": "u1", "name": "n1"}}, flattenReferenceValuesList(ref))
		exp := expandReference(map[string]interface{}{"kind": "vm", "uuid": "u1", "name": "n1"})
		require.Equal(t, "vm", *exp.Kind)
		rv := &v3.ReferenceValues{Kind: "vm", UUID: "u1", Name: "n1"}
		require.Equal(t, map[string]interface{}{"kind": "vm", "uuid": "u1", "name": "n1"}, flattenReference(rv))
		require.Equal(t, map[string]interface{}{"uuid": "u1", "name": "n1"}, flattenExternalNetworkListReference(rv))
		dug := expandDirectoryUserGroup([]interface{}{map[string]interface{}{"distinguished_name": "cn=a"}})
		require.Equal(t, "cn=a", *dug.DistinguishedName)
		sug := expandSamlUserGroup([]interface{}{map[string]interface{}{"idp_uuid": "i1", "name": "g1"}})
		require.Equal(t, "i1", *sug.IdpUUID)
		require.NotNil(t, validateRef(map[string]interface{}{"kind": "vm", "uuid": "u1"}))
		require.Nil(t, validateRefList(nil, nil))
		require.NotNil(t, validateRefList([]interface{}{map[string]interface{}{"uuid": "u1"}}, utils.StringPtr("vm")))
	})
}

func TestUnitPrismFlattenExpandHelpers(t *testing.T) {
	t.Parallel()

	t.Run("expandAccessControlPolicyResources", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixAccessControlPolicy().Schema, map[string]interface{}{})
			access := &v3.AccessControlPolicyResources{}
			expandAccessControlPolicyResources(d, access)
		})
	})

	t.Run("expandScopeExpressionList", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = expandScopeExpressionList(map[string]interface{}{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandScopeExpressionList(map[string]interface{}{"operator": "IN", "right_hand_side": []interface{}{}})
		})
	})

	t.Run("expandEntityExpressionList", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = expandEntityExpressionList(map[string]interface{}{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandEntityExpressionList(map[string]interface{}{"operator": "IN", "right_hand_side": []interface{}{}})
		})
	})

	t.Run("expandRightHandSide", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = expandRightHandSide(map[string]interface{}{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandRightHandSide(map[string]interface{}{"operator": "IN", "right_hand_side": []interface{}{}})
		})
	})

	t.Run("flattenContextList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenContextList([]*v3.ContextList{&v3.ContextList{}})
		})
	})

	t.Run("flattenScopeExpressionList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenScopeExpressionList([]*v3.ScopeFilterExpressionList{&v3.ScopeFilterExpressionList{}})
		})
	})

	t.Run("flattenEntityExpressionList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenEntityExpressionList([]v3.EntityFilterExpressionList{{}})
		})
	})

	t.Run("flattenRightHandSide", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenRightHandSide(v3.RightHandSide{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenRightHandSide(v3.RightHandSide{Collection: utils.StringPtr("x")})
		})
	})

	t.Run("flattenTightHandsideCategories", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			_ = flattenTightHandsideCategories(nil)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenTightHandsideCategories(map[string][]string{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenTightHandsideCategories(map[string][]string{"a": {"b"}})
		})
	})

	t.Run("flattenResourceDomain", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenResourceDomain(&v3.ResourceDomain{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenResourceDomain(&v3.ResourceDomain{})
		})
	})

	t.Run("expandReferenceByMap", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = expandReferenceByMap(map[string]interface{}{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandReferenceByMap(map[string]interface{}{"operator": "IN", "right_hand_side": []interface{}{}})
		})
	})

	t.Run("expandReferenceList", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixProject().Schema, map[string]interface{}{})
			_ = expandReferenceList(d, "account_reference_list")
		})
	})

	t.Run("expandReferenceSet", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixProject().Schema, map[string]interface{}{})
			_ = expandReferenceSet(d, "user_reference_list")
		})
	})

	t.Run("expandDirectoryReferenceList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandDirectoryReferenceList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("flattenDirectoryReferenceList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDirectoryReferenceList([]*v3.ReferenceValues{&v3.ReferenceValues{}})
		})
	})

	t.Run("expandAcp", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixProject().Schema, map[string]interface{}{})
			_ = expandAcp(nil, d)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixProject().Schema, map[string]interface{}{})
			_ = expandAcp([]interface{}{}, d)
		})
	})

	t.Run("expandDirectoryServiceUserPI", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandDirectoryServiceUserPI(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandDirectoryServiceUserPI([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandProjectContextFilterList", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandProjectContextFilterList(nil))
		})
	})

	t.Run("expandProjectInternalResourceDomain", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandProjectInternalResourceDomain(nil))
		})
	})

	t.Run("expandOptionalReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixProject().Schema, map[string]interface{}{})
			_ = expandOptionalReference(d, "subnet_reference_list", "subnet")
		})
	})

	t.Run("expandSnapshotScheduleList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandSnapshotScheduleList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandRetentionPolicy", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = expandRetentionPolicy(map[string]interface{}{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandRetentionPolicy(map[string]interface{}{"num_snapshots": 1})
		})
	})

	t.Run("flattenCategoriesFilter", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenCategoriesFilter(&v3.CategoryFilter{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenCategoriesFilter(&v3.CategoryFilter{})
		})
	})

	t.Run("flattenAvailabilityZoneConnectivityList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenAvailabilityZoneConnectivityList([]*v3.AvailabilityZoneConnectivityList{&v3.AvailabilityZoneConnectivityList{}})
		})
	})

	t.Run("flattenSnapshotScheduleList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnapshotScheduleList([]*v3.SnapshotScheduleList{&v3.SnapshotScheduleList{}})
		})
	})

	t.Run("flattenRetentionPolicy", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenRetentionPolicy(&v3.SnapshotRetentionPolicy{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenRetentionPolicy(&v3.SnapshotRetentionPolicy{})
		})
	})

	t.Run("flattenOrderAvailibilityList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenOrderAvailibilityList([]*v3.OrderedAvailabilityZoneList{&v3.OrderedAvailabilityZoneList{}})
		})
	})

	t.Run("expandStageWork", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = expandStageWork(map[string]interface{}{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandStageWork(map[string]interface{}{"operator": "IN", "right_hand_side": []interface{}{}})
		})
	})

	t.Run("expandEntityInfoList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandEntityInfoList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandFloatingAssignmentList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandFloatingAssignmentList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandVMIPAssignmentList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandVMIPAssignmentList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandZoneNetworkMappingList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandZoneNetworkMappingList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandRecoveryNetwork", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandRecoveryNetwork([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandIPAssignmentList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandIPAssignmentList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandSubnetList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandSubnetList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("flattenStageWork", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenStageWork(&v3.StageWork{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenStageWork(&v3.StageWork{})
		})
	})

	t.Run("flattenParameters", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenParameters(&v3.Parameters{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenParameters(&v3.Parameters{})
		})
	})

	t.Run("flattenRecoveryNetwork", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenRecoveryNetwork(&v3.Network{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenRecoveryNetwork(&v3.Network{})
		})
	})

	t.Run("flattenSubnetList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSubnetList([]*v3.SubnetList{&v3.SubnetList{}})
		})
	})

}
