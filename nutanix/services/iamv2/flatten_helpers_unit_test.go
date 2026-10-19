package iamv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamCommonConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/common/v1/config"
	iamResponse "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/common/v1/response"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authn"
	iamConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authz"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamResponse.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamResponse.ApiLink{}, want: nil},
		{
			name: "populated",
			in: []iamResponse.ApiLink{{
				Href: utils.StringPtr("https://example.test/roles/1"),
				Rel:  utils.StringPtr("self"),
			}},
			want: []map[string]interface{}{{
				"href": utils.StringPtr("https://example.test/roles/1"),
				"rel":  utils.StringPtr("self"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLinks(tt.in))
		})
	}
}

func TestUnitFlattenEntityLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamResponse.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamResponse.ApiLink{}, want: nil},
		{
			name: "populated",
			in: []iamResponse.ApiLink{{
				Href: utils.StringPtr("https://example.test/entities/1"),
				Rel:  utils.StringPtr("self"),
			}},
			want: []map[string]interface{}{{
				"href": "https://example.test/entities/1",
				"rel":  "self",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEntityLinks(tt.in))
		})
	}
}

func TestUnitFlattenRmIdentityType(t *testing.T) {
	t.Parallel()

	user := iamConfig.RMIDENTITYTYPE_USER
	group := iamConfig.RMIDENTITYTYPE_GROUP
	other := iamConfig.RmIdentityType(99)

	tests := []struct {
		name string
		in   *iamConfig.RmIdentityType
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "user", in: &user, want: "USER"},
		{name: "group", in: &group, want: "GROUP"},
		{name: "unknown", in: &other, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRmIdentityType(tt.in))
		})
	}
}

func TestUnitExpandRmIdentityType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want iamConfig.RmIdentityType
	}{
		{name: "empty", in: "", want: iamConfig.RMIDENTITYTYPE_UNKNOWN},
		{name: "user", in: "USER", want: iamConfig.RMIDENTITYTYPE_USER},
		{name: "group", in: "GROUP", want: iamConfig.RMIDENTITYTYPE_GROUP},
		{name: "other", in: "OTHER", want: iamConfig.RMIDENTITYTYPE_UNKNOWN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandRmIdentityType(tt.in)
			require.NotNil(t, got)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestUnitFlattenScopeTemplateNameValues(t *testing.T) {
	t.Parallel()

	oneOf := iamCommonConfig.NewOneOfKVPairValue()
	require.NoError(t, oneOf.SetValue("proj-1"))

	tests := []struct {
		name string
		in   []iamCommonConfig.KVPair
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamCommonConfig.KVPair{}, want: nil},
		{
			name: "populated",
			in: []iamCommonConfig.KVPair{{
				Name:  utils.StringPtr("project"),
				Value: oneOf,
			}},
			want: []map[string]interface{}{{
				"name":  "project",
				"value": "proj-1",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenScopeTemplateNameValues(tt.in))
		})
	}
}

func TestUnitExpandScopeTemplateNameValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       []interface{}
		wantNil  bool
		wantName string
	}{
		{name: "nil", in: nil, wantNil: true},
		{name: "empty", in: []interface{}{}, wantNil: true},
		{
			name:     "populated",
			in:       []interface{}{map[string]interface{}{"name": "project", "value": "p-1"}},
			wantName: "project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandScopeTemplateNameValues(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.Equal(t, tt.wantName, utils.StringValue(got[0].Name))
			require.NotNil(t, got[0].Value)
		})
	}
}

func TestUnitFlattenUserType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.UserType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "local", in: userTypePtr(2), want: "LOCAL"},
		{name: "saml", in: userTypePtr(3), want: "SAML"},
		{name: "ldap", in: userTypePtr(4), want: "LDAP"},
		{name: "external", in: userTypePtr(5), want: "EXTERNAL"},
		{name: "service account", in: userTypePtr(6), want: "SERVICE_ACCOUNT"},
		{name: "unknown", in: userTypePtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenUserType(tt.in))
		})
	}
}

func TestUnitFlattenUserStatusType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.UserStatusType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "active", in: userStatusPtr(2), want: "ACTIVE"},
		{name: "inactive", in: userStatusPtr(3), want: "INACTIVE"},
		{name: "unknown", in: userStatusPtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenUserStatusType(tt.in))
		})
	}
}

func TestUnitFlattenGroupType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.GroupType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "saml", in: groupTypePtr(2), want: "SAML"},
		{name: "ldap", in: groupTypePtr(3), want: "LDAP"},
		{name: "unknown", in: groupTypePtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenGroupType(tt.in))
		})
	}
}

func TestUnitFlattenDirectoryType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.DirectoryType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "ad", in: directoryTypePtr(2), want: "ACTIVE_DIRECTORY"},
		{name: "ldap", in: directoryTypePtr(3), want: "OPEN_LDAP"},
		{name: "unknown", in: directoryTypePtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDirectoryType(tt.in))
		})
	}
}

func TestUnitFlattenGroupSearchType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.GroupSearchType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "non recursive", in: groupSearchTypePtr(2), want: "NON_RECURSIVE"},
		{name: "recursive", in: groupSearchTypePtr(3), want: "RECURSIVE"},
		{name: "unknown", in: groupSearchTypePtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenGroupSearchType(tt.in))
		})
	}
}

func TestUnitFlattenAuthorizationPolicyType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *iamConfig.AuthorizationPolicyType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "user defined", in: authPolicyTypePtr(2), want: "USER_DEFINED"},
		{name: "service defined", in: authPolicyTypePtr(3), want: "SERVICE_DEFINED"},
		{name: "predefined read only", in: authPolicyTypePtr(4), want: "PREDEFINED_READ_ONLY"},
		{name: "predefined update identity", in: authPolicyTypePtr(5), want: "PREDEFINED_UPDATE_IDENTITY_ONLY"},
		{name: "service defined read only", in: authPolicyTypePtr(6), want: "SERVICE_DEFINED_READ_ONLY"},
		{name: "unknown", in: authPolicyTypePtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAuthorizationPolicyType(tt.in))
		})
	}
}

func TestUnitFlattenOperationType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *iamConfig.OperationType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "internal", in: operationTypePtr(2), want: "INTERNAL"},
		{name: "system defined only", in: operationTypePtr(3), want: "SYSTEM_DEFINED_ONLY"},
		{name: "external", in: operationTypePtr(4), want: "EXTERNAL"},
		{name: "unknown", in: operationTypePtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenOperationType(tt.in))
		})
	}
}

func TestUnitFlattenHTTPMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *iamConfig.HttpMethod
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "post", in: httpMethodPtr(2), want: "HTTPMETHOD_POST"},
		{name: "get", in: httpMethodPtr(3), want: "HTTPMETHOD_GET"},
		{name: "put", in: httpMethodPtr(4), want: "HTTPMETHOD_PUT"},
		{name: "patch", in: httpMethodPtr(5), want: "HTTPMETHOD_PATCH"},
		{name: "delete", in: httpMethodPtr(6), want: "HTTPMETHOD_DELETE"},
		{name: "unknown", in: httpMethodPtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenHTTPMethod(tt.in))
		})
	}
}

func TestUnitFlattenNameIDPolicyFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.NameIdPolicyFormat
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "email", in: nameIDPtr(2), want: "emailAddress"},
		{name: "unspecified", in: nameIDPtr(3), want: "unspecified"},
		{name: "x509", in: nameIDPtr(4), want: "X509SubjectName"},
		{name: "windows", in: nameIDPtr(5), want: "WindowsDomainQualifiedName"},
		{name: "encrypted", in: nameIDPtr(6), want: "encrypted"},
		{name: "entity", in: nameIDPtr(7), want: "entity"},
		{name: "kerberos", in: nameIDPtr(8), want: "kerberos"},
		{name: "persistent", in: nameIDPtr(9), want: "persistent"},
		{name: "transient", in: nameIDPtr(10), want: "transient"},
		{name: "unknown", in: nameIDPtr(99), want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNameIDPolicyFormat(tt.in))
		})
	}
}

func TestUnitFlattenTime(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		name string
		in   *time.Time
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "populated", in: &when, want: "2024-01-02T03:04:05Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTime(tt.in))
		})
	}
}

func TestUnitFlattenDsServiceAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.DsServiceAccount
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &import1.DsServiceAccount{},
			want: []map[string]interface{}{{
				"username": (*string)(nil),
				"password": (*string)(nil),
			}},
		},
		{
			name: "populated",
			in: &import1.DsServiceAccount{
				Username: utils.StringPtr("admin"),
				Password: utils.StringPtr("secret"),
			},
			want: []map[string]interface{}{{
				"username": utils.StringPtr("admin"),
				"password": utils.StringPtr("secret"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDsServiceAccount(tt.in))
		})
	}
}

func TestUnitFlattenDsServiceAccountForResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.DsServiceAccount
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated omits write-only password",
			in: &import1.DsServiceAccount{
				Username: utils.StringPtr("svc"),
				Password: utils.StringPtr("pw"),
			},
			want: []map[string]interface{}{{
				"username": utils.StringPtr("svc"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDsServiceAccountForResource(tt.in))
		})
	}
}

func TestUnitFlattenUserConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.UserConfiguration
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &import1.UserConfiguration{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &import1.UserConfiguration{
				UsernameAttribute: utils.StringPtr("uid"),
				UserSearchBase:    utils.StringPtr("ou=users"),
				UserObjectClass:   utils.StringPtr("person"),
			},
			want: []map[string]interface{}{{
				"username_attribute": utils.StringPtr("uid"),
				"user_search_base":   utils.StringPtr("ou=users"),
				"user_object_class":  utils.StringPtr("person"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenUserConfiguration(tt.in))
		})
	}
}

func TestUnitFlattenUserGroupConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.UserGroupConfiguration
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &import1.UserGroupConfiguration{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &import1.UserGroupConfiguration{
				GroupMemberAttribute:      utils.StringPtr("member"),
				GroupMemberAttributeValue: utils.StringPtr("dn"),
				GroupObjectClass:          utils.StringPtr("group"),
				GroupSearchBase:           utils.StringPtr("ou=groups"),
			},
			want: []map[string]interface{}{{
				"group_object_class":           utils.StringPtr("member"),
				"group_search_base":            utils.StringPtr("dn"),
				"group_member_attribute":       utils.StringPtr("group"),
				"group_member_attribute_value": utils.StringPtr("ou=groups"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenUserGroupConfiguration(tt.in))
		})
	}
}

func TestUnitFlattenOpenLdapConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *import1.OpenLdapConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &import1.OpenLdapConfig{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &import1.OpenLdapConfig{
				UserConfiguration: &import1.UserConfiguration{
					UsernameAttribute: utils.StringPtr("uid"),
				},
			},
			want: []map[string]interface{}{{
				"user_configuration": []map[string]interface{}{{
					"username_attribute": utils.StringPtr("uid"),
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenOpenLdapConfig(tt.in))
		})
	}
}

func TestUnitFlattenIdpMetadata(t *testing.T) {
	t.Parallel()

	fmtVal := import1.NameIdPolicyFormat(2)
	tests := []struct {
		name string
		in   *import1.IdpMetadata
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &import1.IdpMetadata{},
			want: []map[string]interface{}{{
				"entity_id":   (*string)(nil),
				"login_url":   (*string)(nil),
				"logout_url":  (*string)(nil),
				"error_url":   (*string)(nil),
				"certificate": (*string)(nil),
			}},
		},
		{
			name: "populated",
			in: &import1.IdpMetadata{
				EntityId:           utils.StringPtr("entity"),
				LoginUrl:           utils.StringPtr("https://login"),
				LogoutUrl:          utils.StringPtr("https://logout"),
				ErrorUrl:           utils.StringPtr("https://error"),
				Certificate:        utils.StringPtr("cert"),
				NameIdPolicyFormat: &fmtVal,
			},
			want: []map[string]interface{}{{
				"entity_id":             utils.StringPtr("entity"),
				"login_url":             utils.StringPtr("https://login"),
				"logout_url":            utils.StringPtr("https://logout"),
				"error_url":             utils.StringPtr("https://error"),
				"certificate":           utils.StringPtr("cert"),
				"name_id_policy_format": "emailAddress",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIdpMetadata(tt.in))
		})
	}
}

func TestUnitFlattenIdentityFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamConfig.IdentityFilter
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamConfig.IdentityFilter{}, want: nil},
		{
			name: "populated",
			in: []iamConfig.IdentityFilter{{
				Reserved_: map[string]interface{}{"anyof": []interface{}{"user"}},
			}},
			want: []interface{}{
				map[string]interface{}{"reserved": `{"anyof":["user"]}`},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIdentityFilters(tt.in))
		})
	}
}

func TestUnitFlattenEntityFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamConfig.EntityFilter
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamConfig.EntityFilter{}, want: nil},
		{
			name: "populated",
			in: []iamConfig.EntityFilter{{
				Reserved_: map[string]interface{}{"eq": "cluster"},
			}},
			want: []interface{}{
				map[string]interface{}{"reserved": `{"eq":"cluster"}`},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEntityFilters(tt.in))
		})
	}
}

func TestUnitFlattenArgumentsMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]string
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: map[string]string{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in:   map[string]string{"reason": "revoke"},
			want: []map[string]interface{}{{
				"property_name": "reason",
				"value":         "revoke",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenArgumentsMap(tt.in))
		})
	}
}

func TestUnitFlattenUsersEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.User
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.User{}, want: nil},
		{
			name: "populated",
			in: []import1.User{{
				ExtId:    utils.StringPtr("user-1"),
				Username: utils.StringPtr("alice"),
				UserType: userTypePtr(2),
			}},
			want: []interface{}{
				map[string]interface{}{
					"ext_id":    utils.StringPtr("user-1"),
					"username":  utils.StringPtr("alice"),
					"user_type": "LOCAL",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenUsersEntities(tt.in))
		})
	}
}

func TestUnitFlattenUserGroupEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.UserGroup
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.UserGroup{}, want: nil},
		{
			name: "populated",
			in: []import1.UserGroup{{
				ExtId:     utils.StringPtr("ug-1"),
				Name:      utils.StringPtr("admins"),
				GroupType: groupTypePtr(3),
			}},
			want: []interface{}{
				map[string]interface{}{
					"ext_id":     utils.StringPtr("ug-1"),
					"name":       utils.StringPtr("admins"),
					"group_type": "LDAP",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenUserGroupEntities(tt.in))
		})
	}
}

func TestUnitFlattenEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamConfig.Entity
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []iamConfig.Entity{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []iamConfig.Entity{{
				ExtId:       utils.StringPtr("ent-1"),
				Name:        utils.StringPtr("cluster"),
				DisplayName: utils.StringPtr("Cluster"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenEntities(tt.in)
			if tt.name == "populated" {
				require.Len(t, got, 1)
				require.Equal(t, "ent-1", got[0]["ext_id"])
				require.Equal(t, "cluster", got[0]["name"])
				require.Equal(t, "Cluster", got[0]["display_name"])
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitFlattenAttributeList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamConfig.AttributeEntity
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamConfig.AttributeEntity{}, want: nil},
		{
			name: "populated",
			in: []iamConfig.AttributeEntity{{
				DisplayName: utils.StringPtr("Name"),
			}},
			want: []map[string]interface{}{{
				"display_name":        "Name",
				"supported_operators": []string{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenAttributeList(tt.in)
			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.Equal(t, "Name", got[0]["display_name"])
		})
	}
}

func TestUnitFlattenDirectoryServiceSearchAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.DirectoryServiceSearchAttribute
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []import1.DirectoryServiceSearchAttribute{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []import1.DirectoryServiceSearchAttribute{{
				Name:   utils.StringPtr("mail"),
				Values: []string{"a@example.test"},
			}},
			want: []map[string]interface{}{{
				"name":   utils.StringPtr("mail"),
				"values": []string{"a@example.test"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDirectoryServiceSearchAttributes(tt.in))
		})
	}
}

func TestUnitFlattenDirectoryServiceSearchResults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.DirectoryServiceSearchEntity
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []import1.DirectoryServiceSearchEntity{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []import1.DirectoryServiceSearchEntity{{
				Name:          utils.StringPtr("alice"),
				IdentityExtId: utils.StringPtr("id-1"),
			}},
			want: []map[string]interface{}{{
				"entity_type":     (*string)(nil),
				"name":            utils.StringPtr("alice"),
				"identity_ext_id": utils.StringPtr("id-1"),
				"attributes":      []map[string]interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDirectoryServiceSearchResults(tt.in))
		})
	}
}

func TestUnitFlattenRoleMembershipSummaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamConfig.RoleMembershipSummary
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamConfig.RoleMembershipSummary{}, want: nil},
		{
			name: "populated",
			in: []iamConfig.RoleMembershipSummary{{
				ExtId:      utils.StringPtr("sum-1"),
				UsersCount: utils.Int64Ptr(2),
			}},
			want: []interface{}{
				map[string]interface{}{
					"ext_id":      "sum-1",
					"users_count": int64(2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRoleMembershipSummaries(tt.in))
		})
	}
}

func TestUnitFlattenRoleMembershipEntities(t *testing.T) {
	t.Parallel()

	idType := iamConfig.RMIDENTITYTYPE_USER
	tests := []struct {
		name string
		in   []iamConfig.RoleMembershipProjection
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamConfig.RoleMembershipProjection{}, want: nil},
		{
			name: "populated",
			in: []iamConfig.RoleMembershipProjection{{
				ExtId:        utils.StringPtr("rm-1"),
				RoleExtId:    utils.StringPtr("role-1"),
				IdentityType: &idType,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenRoleMembershipEntities(tt.in)
			if tt.name == "populated" {
				require.Len(t, got, 1)
				m := got[0].(map[string]interface{})
				require.Equal(t, "rm-1", m["ext_id"])
				require.Equal(t, "role-1", m["role_ext_id"])
				require.Equal(t, "USER", m["identity_type"])
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitFlattenKeyDetails(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenKeyDetails(nil))
	})

	t.Run("api key details", func(t *testing.T) {
		t.Parallel()
		details := import1.NewOneOfKeyKeyDetails()
		api := import1.NewApiKeyDetails()
		api.ApiKey = utils.StringPtr("api-secret")
		require.NoError(t, details.SetValue(*api))
		got := flattenKeyDetails(details)
		require.Equal(t, []map[string]interface{}{{
			"api_key_details": []map[string]interface{}{{
				"api_key": utils.StringPtr("api-secret"),
			}},
		}}, got)
	})

	t.Run("object key details", func(t *testing.T) {
		t.Parallel()
		details := import1.NewOneOfKeyKeyDetails()
		obj := import1.NewObjectKeyDetails()
		obj.SecretKey = utils.StringPtr("secret")
		obj.AccessKey = utils.StringPtr("access")
		require.NoError(t, details.SetValue(*obj))
		got := flattenKeyDetails(details)
		require.Equal(t, []map[string]interface{}{{
			"object_key_details": []map[string]interface{}{{
				"secret_key": utils.StringPtr("secret"),
				"access_key": utils.StringPtr("access"),
			}},
		}}, got)
	})
}

func TestUnitFlattenAssociatedEndpointList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []iamConfig.AssociatedEndpoint
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []iamConfig.AssociatedEndpoint{}, wantLen: 0},
		{name: "populated", in: []iamConfig.AssociatedEndpoint{{
			EndpointUrl: utils.StringPtr("/v4/ops"),
			HttpMethod:  httpMethodPtr(3),
		}}, wantLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenAssociatedEndpointList(tt.in)
			if tt.wantLen == 0 {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, tt.wantLen)
			require.Equal(t, "/v4/ops", got[0]["endpoint_url"])
			require.Equal(t, "HTTPMETHOD_GET", got[0]["http_method"])
		})
	}
}

func TestUnitUserKeyUpdate(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixUserKeyV2().Schema, map[string]interface{}{
		"user_ext_id": "user-1",
		"name":        "key",
	})
	require.Nil(t, resourceNutanixUserKeyV2Update(t.Context(), d, nil))
}

func userTypePtr(v int) *import1.UserType {
	x := import1.UserType(v)
	return &x
}

func userStatusPtr(v int) *import1.UserStatusType {
	x := import1.UserStatusType(v)
	return &x
}

func groupTypePtr(v int) *import1.GroupType {
	x := import1.GroupType(v)
	return &x
}

func directoryTypePtr(v int) *import1.DirectoryType {
	x := import1.DirectoryType(v)
	return &x
}

func groupSearchTypePtr(v int) *import1.GroupSearchType {
	x := import1.GroupSearchType(v)
	return &x
}

func authPolicyTypePtr(v int) *iamConfig.AuthorizationPolicyType {
	x := iamConfig.AuthorizationPolicyType(v)
	return &x
}

func operationTypePtr(v int) *iamConfig.OperationType {
	x := iamConfig.OperationType(v)
	return &x
}

func httpMethodPtr(v int) *iamConfig.HttpMethod {
	x := iamConfig.HttpMethod(v)
	return &x
}

func nameIDPtr(v int) *import1.NameIdPolicyFormat {
	x := import1.NameIdPolicyFormat(v)
	return &x
}

func TestUnitFlattenRolesEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamConfig.Role
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamConfig.Role{}, want: nil},
		{
			name: "populated",
			in: []iamConfig.Role{{
				ExtId:       utils.StringPtr("role-1"),
				DisplayName: utils.StringPtr("Admin"),
			}},
			want: []interface{}{
				map[string]interface{}{
					"ext_id":       utils.StringPtr("role-1"),
					"display_name": utils.StringPtr("Admin"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRolesEntities(tt.in))
		})
	}
}

func TestUnitFlattenPermissionEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamConfig.Operation
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamConfig.Operation{}, want: nil},
		{
			name: "populated",
			in: []iamConfig.Operation{{
				ExtId:         utils.StringPtr("op-1"),
				DisplayName:   utils.StringPtr("View"),
				OperationType: operationTypePtr(2),
			}},
			want: []interface{}{
				map[string]interface{}{
					"ext_id":         utils.StringPtr("op-1"),
					"display_name":   utils.StringPtr("View"),
					"operation_type": "INTERNAL",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPermissionEntities(tt.in))
		})
	}
}

func TestUnitFlattenAuthorizationPolicyEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []iamConfig.AuthorizationPolicyProjection
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []iamConfig.AuthorizationPolicyProjection{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAuthorizationPolicyEntities(tt.in))
		})
	}

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenAuthorizationPolicyEntities([]iamConfig.AuthorizationPolicyProjection{{
			ExtId:                   utils.StringPtr("acp-1"),
			DisplayName:             utils.StringPtr("policy"),
			AuthorizationPolicyType: authPolicyTypePtr(2),
		}})
		require.Len(t, got, 1)
		m := got[0].(map[string]interface{})
		require.Equal(t, utils.StringPtr("acp-1"), m["ext_id"])
		require.Equal(t, "USER_DEFINED", m["authorization_policy_type"])
	})
}

func TestUnitFlattenDirectoryServicesEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.DirectoryService
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.DirectoryService{}, want: nil},
		{
			name: "populated",
			in: []import1.DirectoryService{{
				ExtId:         utils.StringPtr("ds-1"),
				Name:          utils.StringPtr("ad"),
				DirectoryType: directoryTypePtr(2),
			}},
			want: []interface{}{
				map[string]interface{}{
					"ext_id":         utils.StringPtr("ds-1"),
					"name":           utils.StringPtr("ad"),
					"domain_name":    (*string)(nil),
					"directory_type": "ACTIVE_DIRECTORY",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenDirectoryServicesEntities(tt.in)
			if tt.name == "populated" {
				require.Len(t, got, 1)
				m := got[0].(map[string]interface{})
				require.Equal(t, utils.StringPtr("ds-1"), m["ext_id"])
				require.Equal(t, "ACTIVE_DIRECTORY", m["directory_type"])
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitFlattenIdentityProvidersEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import1.SamlIdentityProvider
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import1.SamlIdentityProvider{}, want: nil},
		{
			name: "populated",
			in: []import1.SamlIdentityProvider{{
				ExtId: utils.StringPtr("idp-1"),
				Name:  utils.StringPtr("saml"),
			}},
			want: []interface{}{
				map[string]interface{}{
					"ext_id": utils.StringPtr("idp-1"),
					"name":   utils.StringPtr("saml"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenIdentityProvidersEntities(tt.in)
			if tt.name == "populated" {
				require.Len(t, got, 1)
				m := got[0].(map[string]interface{})
				require.Equal(t, utils.StringPtr("idp-1"), m["ext_id"])
				require.Equal(t, utils.StringPtr("saml"), m["name"])
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitFlattenBucketsAccessKeys(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenBucketsAccessKeys(import1.User{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenBucketsAccessKeys(import1.User{
			BucketsAccessKeys: []import1.BucketsAccessKey{{
				ExtId:         utils.StringPtr("bak-1"),
				AccessKeyName: utils.StringPtr("key"),
			}},
		})
		require.Equal(t, []map[string]interface{}{{
			"ext_id":          "bak-1",
			"access_key_name": "key",
		}}, got)
	})
}

func TestUnitFlattenAdditionalAttributes(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenAdditionalAttributes(import1.User{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		oneOf := iamCommonConfig.NewOneOfKVPairValue()
		require.NoError(t, oneOf.SetValue("eng"))
		got := flattenAdditionalAttributes(import1.User{
			AdditionalAttributes: []iamCommonConfig.KVPair{{
				Name:  utils.StringPtr("dept"),
				Value: oneOf,
			}},
		})
		require.Len(t, got, 1)
		m := got[0].(map[string]interface{})
		require.Equal(t, "dept", m["name"])
		require.NotNil(t, m["value"])
	})
}

func TestUnitFlattenKeysEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, flattenKeysEntities(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, flattenKeysEntities([]import1.Key{}))
	})
}
