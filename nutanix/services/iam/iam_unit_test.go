package iam

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitIAMSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixUser().InternalValidate(ResourceNutanixUser().Schema, true))
	require.NoError(t, ResourceNutanixUserGroups().InternalValidate(ResourceNutanixUserGroups().Schema, true))
	require.NoError(t, ResourceNutanixRole().InternalValidate(ResourceNutanixRole().Schema, true))
	require.NoError(t, DataSourceNutanixUser().InternalValidate(DataSourceNutanixUser().Schema, false))
	require.NoError(t, DataSourceNutanixUsers().InternalValidate(DataSourceNutanixUsers().Schema, false))
	require.NoError(t, DataSourceNutanixUserGroup().InternalValidate(DataSourceNutanixUserGroup().Schema, false))
	require.NoError(t, DataSourceNutanixUserGroups().InternalValidate(DataSourceNutanixUserGroups().Schema, false))
	require.NoError(t, DataSourceNutanixRole().InternalValidate(DataSourceNutanixRole().Schema, false))
	require.NoError(t, DataSourceNutanixRoles().InternalValidate(DataSourceNutanixRoles().Schema, false))
	require.NoError(t, DataSourceNutanixPermission().InternalValidate(DataSourceNutanixPermission().Schema, false))
	require.NoError(t, DataSourceNutanixPermissions().InternalValidate(DataSourceNutanixPermissions().Schema, false))
	require.NoError(t, DataSourceNutanixAccessControlPolicy().InternalValidate(DataSourceNutanixAccessControlPolicy().Schema, false))
	require.NoError(t, DataSourceNutanixAccessControlPolicies().InternalValidate(DataSourceNutanixAccessControlPolicies().Schema, false))
	require.NoError(t, categoriesSchema().Elem.(*schema.Resource).InternalValidate(categoriesSchema().Elem.(*schema.Resource).Schema, true))
}

func TestUnitExpandCategories(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(categoriesSchema().Set, []interface{}{})
		require.Equal(t, map[string]string{}, expandCategories(set))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(categoriesSchema().Set, []interface{}{
			map[string]interface{}{"name": "Environment", "value": "Dev"},
		})
		require.Equal(t, map[string]string{"Environment": "Dev"}, expandCategories(set))
	})
}

func TestUnitFlattenCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]string
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{name: "empty", in: map[string]string{}, want: []interface{}{}},
		{
			name: "populated",
			in:   map[string]string{"Environment": "Prod"},
			want: []interface{}{map[string]interface{}{"name": "Environment", "value": "Prod"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCategories(tt.in))
		})
	}
}

func TestUnitFlattenReferenceValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.Reference
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{
			name: "populated",
			in: &v3.Reference{
				Kind: utils.StringPtr("user"),
				UUID: utils.StringPtr("uuid-1"),
				Name: utils.StringPtr("alice"),
			},
			want: map[string]interface{}{
				"kind": "user",
				"uuid": "uuid-1",
				"name": "alice",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReferenceValues(tt.in))
		})
	}
}

func TestUnitValidateRef(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]interface{}
		want *v3.Reference
	}{
		{name: "empty", in: map[string]interface{}{}, want: nil},
		{
			name: "populated",
			in:   map[string]interface{}{"kind": "role", "uuid": "r-1", "name": "admin"},
			want: &v3.Reference{
				Kind: utils.StringPtr("role"),
				UUID: utils.StringPtr("r-1"),
				Name: utils.StringPtr("admin"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, validateRef(tt.in))
		})
	}
}

func TestUnitExpandReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]interface{}
		want *v3.Reference
	}{
		{name: "empty", in: map[string]interface{}{}, want: nil},
		{
			name: "populated",
			in:   map[string]interface{}{"kind": "project", "uuid": "p-1", "name": "demo"},
			want: &v3.Reference{
				Kind: utils.StringPtr("project"),
				UUID: utils.StringPtr("p-1"),
				Name: utils.StringPtr("demo"),
			},
		},
		{
			name: "empty name skipped",
			in:   map[string]interface{}{"kind": "project", "uuid": "p-1", "name": ""},
			want: &v3.Reference{
				Kind: utils.StringPtr("project"),
				UUID: utils.StringPtr("p-1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandReference(tt.in))
		})
	}
}

func TestUnitFlattenArrayReferenceValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.Reference
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*v3.Reference{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []*v3.Reference{{
				Kind: utils.StringPtr("permission"),
				UUID: utils.StringPtr("perm-1"),
				Name: utils.StringPtr("view"),
			}},
			want: []map[string]interface{}{{
				"kind": "permission",
				"uuid": "perm-1",
				"name": "view",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenArrayReferenceValues(tt.in))
		})
	}
}

func TestUnitValidateArrayRef(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(schema.HashResource(&schema.Resource{
			Schema: map[string]*schema.Schema{
				"kind": {Type: schema.TypeString, Optional: true},
				"uuid": {Type: schema.TypeString, Optional: true},
				"name": {Type: schema.TypeString, Optional: true},
			},
		}), []interface{}{})
		require.Nil(t, validateArrayRef(set, nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(schema.HashResource(&schema.Resource{
			Schema: map[string]*schema.Schema{
				"kind": {Type: schema.TypeString, Optional: true},
				"uuid": {Type: schema.TypeString, Optional: true},
				"name": {Type: schema.TypeString, Optional: true},
			},
		}), []interface{}{
			map[string]interface{}{"kind": "ignored", "uuid": "u-1", "name": "n1"},
		})
		got := validateArrayRef(set, utils.StringPtr("permission"))
		require.Len(t, got, 1)
		require.Equal(t, utils.StringPtr("permission"), got[0].Kind)
		require.Equal(t, utils.StringPtr("u-1"), got[0].UUID)
		require.Equal(t, utils.StringPtr("n1"), got[0].Name)
	})
}

func TestUnitValidateRefList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		in        []interface{}
		kindValue *string
		want      *v3.Reference
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name:      "populated with override kind",
			in:        []interface{}{map[string]interface{}{"kind": "other", "uuid": "p-1", "name": "proj"}},
			kindValue: utils.StringPtr("project"),
			want: &v3.Reference{
				Kind: utils.StringPtr("project"),
				UUID: utils.StringPtr("p-1"),
				Name: utils.StringPtr("proj"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, validateRefList(tt.in, tt.kindValue))
		})
	}
}

func TestUnitFlattenReferenceValuesList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.Reference
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.Reference{
				Kind: utils.StringPtr("owner"),
				UUID: utils.StringPtr("o-1"),
				Name: utils.StringPtr("bob"),
			},
			want: []interface{}{map[string]interface{}{
				"kind": "owner",
				"uuid": "o-1",
				"name": "bob",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReferenceValuesList(tt.in))
		})
	}
}

func TestUnitSetRSEntityMetadata(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	meta := &v3.Metadata{
		LastUpdateTime: &when,
		UUID:           utils.StringPtr("meta-1"),
		CreationTime:   &when,
		SpecVersion:    utils.Int64Ptr(3),
		SpecHash:       utils.StringPtr("hash"),
		Name:           utils.StringPtr("role-a"),
		Categories:     map[string]string{"App": "db"},
	}
	gotMeta, gotCats := setRSEntityMetadata(meta)
	require.Equal(t, "meta-1", gotMeta["uuid"])
	require.Equal(t, "3", gotMeta["spec_version"])
	require.Equal(t, "hash", gotMeta["spec_hash"])
	require.Equal(t, "role-a", gotMeta["name"])
	require.Equal(t, []interface{}{map[string]interface{}{"name": "App", "value": "db"}}, gotCats)
}

func TestUnitBuildDataSourceListMetadata(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(schema.HashResource(DataSourceNutanixRoles().Schema["metadata"].Elem.(*schema.Resource)), []interface{}{})
		got := buildDataSourceListMetadata(set)
		require.NotNil(t, got)
		require.Nil(t, got.Filter)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		set := schema.NewSet(schema.HashResource(DataSourceNutanixRoles().Schema["metadata"].Elem.(*schema.Resource)), []interface{}{
			map[string]interface{}{
				"filter":         "name==admin",
				"kind":           "role",
				"sort_order":     "ASCENDING",
				"offset":         5,
				"length":         10,
				"sort_attribute": "name",
			},
		})
		got := buildDataSourceListMetadata(set)
		require.Equal(t, utils.StringPtr("name==admin"), got.Filter)
		require.Equal(t, utils.StringPtr("role"), got.Kind)
		require.Equal(t, utils.StringPtr("ASCENDING"), got.SortOrder)
		require.Equal(t, utils.Int64Ptr(5), got.Offset)
		require.Equal(t, utils.Int64Ptr(10), got.Length)
		require.Equal(t, utils.StringPtr("name"), got.SortAttribute)
	})
}

func TestUnitFlattenRightHandSide(t *testing.T) {
	t.Parallel()

	got := flattenRightHandSide(v3.RightHandSide{
		Collection: utils.StringPtr("ALL"),
		UUIDList:   []string{"u1"},
		Categories: map[string][]string{"Env": {"Dev"}},
	})
	require.Len(t, got, 1)
	m := got[0].(map[string]interface{})
	require.Equal(t, "ALL", m["collection"])
	require.Equal(t, []string{"u1"}, m["uuid_list"])
	require.Equal(t, []interface{}{map[string]interface{}{"name": "Env", "value": []string{"Dev"}}}, m["categories"])
}

func TestUnitFlattenTightHandsideCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string][]string
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{name: "empty", in: map[string][]string{}, want: []interface{}{}},
		{
			name: "populated",
			in:   map[string][]string{"Dept": {"Eng"}},
			want: []interface{}{map[string]interface{}{"name": "Dept", "value": []string{"Eng"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTightHandsideCategories(tt.in))
		})
	}
}

func TestUnitFlattenScopeExpressionList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []*v3.ScopeFilterExpressionList
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []*v3.ScopeFilterExpressionList{}, wantLen: 0},
		{
			name: "populated",
			in: []*v3.ScopeFilterExpressionList{{
				LeftHandSide:  "CATEGORY",
				Operator:      "IN",
				RightHandSide: v3.RightHandSide{Collection: utils.StringPtr("ALL")},
			}},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenScopeExpressionList(tt.in)
			require.Len(t, got, tt.wantLen)
			if tt.wantLen > 0 {
				m := got[0].(map[string]interface{})
				require.Equal(t, "CATEGORY", m["left_hand_side"])
				require.Equal(t, "IN", m["operator"])
			}
		})
	}
}

func TestUnitFlattenEntityExpressionList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []v3.EntityFilterExpressionList
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []v3.EntityFilterExpressionList{}, wantLen: 0},
		{
			name: "populated",
			in: []v3.EntityFilterExpressionList{{
				LeftHandSide:  v3.LeftHandSide{EntityType: utils.StringPtr("vm")},
				Operator:      "IN",
				RightHandSide: v3.RightHandSide{UUIDList: []string{"vm-1"}},
			}},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenEntityExpressionList(tt.in)
			require.Len(t, got, tt.wantLen)
			if tt.wantLen > 0 {
				m := got[0].(map[string]interface{})
				require.Equal(t, "vm", m["left_hand_side_entity_type"])
				require.Equal(t, "IN", m["operator"])
			}
		})
	}
}

func TestUnitFlattenContextList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []*v3.ContextList
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []*v3.ContextList{}, wantLen: 0},
		{
			name: "populated",
			in: []*v3.ContextList{{
				ScopeFilterExpressionList: []*v3.ScopeFilterExpressionList{{
					LeftHandSide: "CATEGORY",
					Operator:     "IN",
				}},
				EntityFilterExpressionList: []v3.EntityFilterExpressionList{{
					LeftHandSide: v3.LeftHandSide{EntityType: utils.StringPtr("vm")},
					Operator:     "IN",
				}},
			}},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenContextList(tt.in)
			require.Len(t, got, tt.wantLen)
		})
	}
}

func TestUnitFlattenDirectoryServiceUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.DirectoryServiceUser
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.DirectoryServiceUser{
				UserPrincipalName:         utils.StringPtr("user@example.com"),
				DefaultUserPrincipalName:  utils.StringPtr("default@example.com"),
				DirectoryServiceReference: &v3.Reference{Kind: utils.StringPtr("directory_service"), UUID: utils.StringPtr("ds-1")},
			},
			want: []interface{}{map[string]interface{}{
				"user_principal_name":         utils.StringPtr("user@example.com"),
				"default_user_principal_name": utils.StringPtr("default@example.com"),
				"directory_service_reference": []interface{}{map[string]interface{}{
					"kind": "directory_service",
					"uuid": "ds-1",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDirectoryServiceUser(tt.in))
		})
	}
}

func TestUnitFlattenIdentityProviderUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.IdentityProvider
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.IdentityProvider{
				Username:                  utils.StringPtr("idp-user"),
				IdentityProviderReference: &v3.Reference{Kind: utils.StringPtr("identity_provider"), UUID: utils.StringPtr("idp-1")},
			},
			want: []interface{}{map[string]interface{}{
				"username": utils.StringPtr("idp-user"),
				"identity_provider_reference": []interface{}{map[string]interface{}{
					"kind": "identity_provider",
					"uuid": "idp-1",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIdentityProviderUser(tt.in))
		})
	}
}

func TestUnitExpandDirectoryServiceUser(t *testing.T) {
	t.Parallel()

	t.Run("missing", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixUser().Schema, map[string]interface{}{})
		require.Nil(t, expandDirectoryServiceUser(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixUser().Schema, map[string]interface{}{
			"directory_service_user": []interface{}{map[string]interface{}{
				"user_principal_name": "user@example.com",
				"directory_service_reference": []interface{}{map[string]interface{}{
					"kind": "directory_service",
					"uuid": "ds-1",
					"name": "ds",
				}},
			}},
		})
		got := expandDirectoryServiceUser(d)
		require.NotNil(t, got)
		require.Equal(t, utils.StringPtr("user@example.com"), got.UserPrincipalName)
		require.NotNil(t, got.DirectoryServiceReference)
		require.Equal(t, utils.StringPtr("ds-1"), got.DirectoryServiceReference.UUID)
	})
}

func TestUnitExpandIdentityProviderUser(t *testing.T) {
	t.Parallel()

	t.Run("missing", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixUser().Schema, map[string]interface{}{})
		require.Nil(t, expandIdentityProviderUser(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixUser().Schema, map[string]interface{}{
			"identity_provider_user": []interface{}{map[string]interface{}{
				"username": "idp-user",
				"identity_provider_reference": []interface{}{map[string]interface{}{
					"kind": "identity_provider",
					"uuid": "idp-1",
					"name": "idp",
				}},
			}},
		})
		got := expandIdentityProviderUser(d)
		require.NotNil(t, got)
		require.Equal(t, utils.StringPtr("idp-user"), got.Username)
		require.NotNil(t, got.IdentityProviderReference)
		require.Equal(t, utils.StringPtr("idp-1"), got.IdentityProviderReference.UUID)
	})
}

func TestUnitExpandDirectoryUserGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *v3.DirectoryServiceUserGroup
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in:   []interface{}{map[string]interface{}{"distinguished_name": "CN=group,DC=example"}},
			want: &v3.DirectoryServiceUserGroup{DistinguishedName: utils.StringPtr("CN=group,DC=example")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandDirectoryUserGroup(tt.in))
		})
	}
}

func TestUnitExpandSamlUserGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *v3.SamlUserGroup
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in:   []interface{}{map[string]interface{}{"idp_uuid": "idp-1", "name": "saml-group"}},
			want: &v3.SamlUserGroup{IdpUUID: utils.StringPtr("idp-1"), Name: utils.StringPtr("saml-group")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandSamlUserGroup(tt.in))
		})
	}
}

func TestUnitFlattenSamlUserGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.SamlUserGroup
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &v3.SamlUserGroup{IdpUUID: utils.StringPtr("idp-1"), Name: utils.StringPtr("g1")},
			want: []interface{}{map[string]interface{}{
				"idp_uuid": utils.StringPtr("idp-1"),
				"name":     utils.StringPtr("g1"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSamlUserGroup(tt.in))
		})
	}
}

func TestUnitFlattenDirectoryServiceUserGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.DirectoryServiceUserGroup
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.DirectoryServiceUserGroup{
				DistinguishedName:         utils.StringPtr("CN=group"),
				DirectoryServiceReference: &v3.Reference{Kind: utils.StringPtr("directory_service"), UUID: utils.StringPtr("ds-1")},
			},
			want: []interface{}{map[string]interface{}{
				"distinguished_name": utils.StringPtr("CN=group"),
				"directory_service_reference": []interface{}{map[string]interface{}{
					"kind": "directory_service",
					"uuid": "ds-1",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDirectoryServiceUserGroup(tt.in))
		})
	}
}

func TestUnitFlattenFieldsPermission(t *testing.T) {
	t.Parallel()

	mode := "ALLOWED"
	got := flattenFieldsPermission(&v3.FieldsPermission{
		FieldMode:     &mode,
		FieldNameList: []*string{utils.StringPtr("name")},
	})
	require.Equal(t, []map[string]interface{}{{
		"field_mode":      &mode,
		"field_name_list": []*string{utils.StringPtr("name")},
	}}, got)
}

func TestUnitMatchUserGroupHelpers(t *testing.T) {
	t.Parallel()

	t.Run("match by name", func(t *testing.T) {
		t.Parallel()
		require.False(t, matchUserGroupByName(nil, "x"))
		require.True(t, matchUserGroupByName(&v3.UserGroupIntentResponse{
			Status: &v3.UserGroupStatus{
				Resources: &v3.UserGroupResources{DisplayName: utils.StringPtr("admins")},
			},
		}, "admins"))
		require.False(t, matchUserGroupByName(&v3.UserGroupIntentResponse{
			Status: &v3.UserGroupStatus{
				Resources: &v3.UserGroupResources{DisplayName: utils.StringPtr("admins")},
			},
		}, "other"))
	})

	t.Run("match by distinguished name", func(t *testing.T) {
		t.Parallel()
		require.False(t, matchUserGroupByDistinguishedName(nil, "x"))
		require.True(t, matchUserGroupByDistinguishedName(&v3.UserGroupIntentResponse{
			Status: &v3.UserGroupStatus{
				Resources: &v3.UserGroupResources{
					DirectoryServiceUserGroup: &v3.DirectoryServiceUserGroup{
						DistinguishedName: utils.StringPtr("CN=admins"),
					},
				},
			},
		}, "CN=admins"))
	})
}

func TestUnitGetMetadataAttributes(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixUser().Schema, map[string]interface{}{
		"categories": []interface{}{
			map[string]interface{}{"name": "Env", "value": "Dev"},
		},
		"project_reference": map[string]interface{}{
			"kind": "project",
			"uuid": "p-1",
			"name": "proj",
		},
		"owner_reference": map[string]interface{}{
			"kind": "user",
			"uuid": "o-1",
			"name": "owner",
		},
	})
	meta := &v3.Metadata{}
	require.NoError(t, getMetadataAttributes(d, meta, "user"))
	require.Equal(t, utils.StringPtr("user"), meta.Kind)
	require.Equal(t, map[string]string{"Env": "Dev"}, meta.Categories)
	require.NotNil(t, meta.ProjectReference)
	require.Equal(t, utils.StringPtr("p-1"), meta.ProjectReference.UUID)
	require.NotNil(t, meta.OwnerReference)
	require.Equal(t, utils.StringPtr("o-1"), meta.OwnerReference.UUID)
}

func TestUnitGetMetadataAttributesV2(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixRole().Schema, map[string]interface{}{
		"name": "role-a",
		"categories": []interface{}{
			map[string]interface{}{"name": "Env", "value": "Prod"},
		},
		"project_reference": []interface{}{map[string]interface{}{
			"kind": "ignored",
			"uuid": "p-2",
			"name": "proj2",
		}},
		"owner_reference": []interface{}{map[string]interface{}{
			"kind": "user",
			"uuid": "o-2",
			"name": "owner2",
		}},
	})
	meta := &v3.Metadata{}
	require.NoError(t, getMetadataAttributesV2(d, meta, "role"))
	require.Equal(t, utils.StringPtr("role"), meta.Kind)
	require.Equal(t, map[string]string{"Env": "Prod"}, meta.Categories)
	require.Equal(t, utils.StringPtr("project"), meta.ProjectReference.Kind)
	require.Equal(t, utils.StringPtr("p-2"), meta.ProjectReference.UUID)
	require.Equal(t, utils.StringPtr("user"), meta.OwnerReference.Kind)
}
