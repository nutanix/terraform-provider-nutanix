package iamv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitExpandKVPair(t *testing.T) {
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
			in:       []interface{}{map[string]interface{}{"name": "attr"}},
			wantName: "attr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandKVPair(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.Equal(t, tt.wantName, utils.StringValue(got[0].Name))
		})
	}
}

func TestUnitExpandDsServiceAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      interface{}
		wantNil bool
		user    string
	}{
		{name: "nil", in: nil, wantNil: true},
		{name: "empty", in: []interface{}{}, wantNil: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"username": "admin",
				"password": "secret",
			}},
			user: "admin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandDsServiceAccount(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.user, utils.StringValue(got.Username))
			require.Equal(t, "secret", utils.StringValue(got.Password))
		})
	}
}

func TestUnitExpandUserConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      interface{}
		wantNil bool
	}{
		{name: "nil", in: nil, wantNil: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"user_object_class":  "person",
				"user_search_base":   "ou=users",
				"username_attribute": "uid",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandUserConfiguration(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, "person", utils.StringValue(got.UserObjectClass))
			require.Equal(t, "ou=users", utils.StringValue(got.UserSearchBase))
			require.Equal(t, "uid", utils.StringValue(got.UsernameAttribute))
		})
	}
}

func TestUnitExpandUserGroupConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      interface{}
		wantNil bool
	}{
		{name: "nil", in: nil, wantNil: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"group_object_class":           "group",
				"group_search_base":            "ou=groups",
				"group_member_attribute":       "member",
				"group_member_attribute_value": "dn",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandUserGroupConfiguration(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, "group", utils.StringValue(got.GroupObjectClass))
			require.Equal(t, "ou=groups", utils.StringValue(got.GroupSearchBase))
			require.Equal(t, "member", utils.StringValue(got.GroupMemberAttribute))
			require.Equal(t, "dn", utils.StringValue(got.GroupMemberAttributeValue))
		})
	}
}

func TestUnitExpandOpenLdapConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      interface{}
		wantNil bool
	}{
		{name: "nil", in: nil, wantNil: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"user_configuration": []interface{}{map[string]interface{}{
					"username_attribute": "uid",
				}},
				"user_group_configuration": []interface{}{map[string]interface{}{
					"group_object_class": "group",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandOpenLdapConfig(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got.UserConfiguration)
			require.Equal(t, "uid", utils.StringValue(got.UserConfiguration.UsernameAttribute))
			require.NotNil(t, got.UserGroupConfiguration)
			require.Equal(t, "group", utils.StringValue(got.UserGroupConfiguration.GroupObjectClass))
		})
	}
}

func TestUnitExpandIdpMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      interface{}
		wantNil bool
	}{
		{name: "nil", in: nil, wantNil: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"entity_id":             "entity",
				"login_url":             "https://login",
				"logout_url":            "https://logout",
				"error_url":             "https://error",
				"certificate":           "cert",
				"name_id_policy_format": "emailAddress",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandIdpMetadata(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, "entity", utils.StringValue(got.EntityId))
			require.Equal(t, "https://login", utils.StringValue(got.LoginUrl))
			require.Equal(t, "https://logout", utils.StringValue(got.LogoutUrl))
			require.Equal(t, "https://error", utils.StringValue(got.ErrorUrl))
			require.Equal(t, "cert", utils.StringValue(got.Certificate))
			require.NotNil(t, got.NameIdPolicyFormat)
			require.Equal(t, import1.NameIdPolicyFormat(2), *got.NameIdPolicyFormat)
		})
	}
}

func TestUnitExpandIdentityFilter(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got, err := expandIdentityFilter(nil)
		require.NoError(t, err)
		require.Nil(t, got)
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got, err := expandIdentityFilter([]interface{}{})
		require.NoError(t, err)
		require.Nil(t, got)
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got, err := expandIdentityFilter([]interface{}{
			map[string]interface{}{"reserved": `{"anyof":["user"]}`},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, []interface{}{"user"}, got[0].Reserved_["anyof"])
	})
	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		_, err := expandIdentityFilter([]interface{}{
			map[string]interface{}{"reserved": `{`},
		})
		require.Error(t, err)
	})
}

func TestUnitExpandEntityFilter(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got, err := expandEntityFilter(nil)
		require.NoError(t, err)
		require.Nil(t, got)
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got, err := expandEntityFilter([]interface{}{})
		require.NoError(t, err)
		require.Nil(t, got)
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got, err := expandEntityFilter([]interface{}{
			map[string]interface{}{"reserved": `{"eq":"vm"}`},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "vm", got[0].Reserved_["eq"])
	})
}

func TestUnitAuthPolicyStringsEquivalent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		{name: "both empty", s1: "", s2: "", want: true},
		{name: "empty and braces", s1: "", s2: "{}", want: true},
		{name: "braces and empty", s1: "{}", s2: "", want: true},
		{name: "both braces", s1: "{}", s2: "{}", want: true},
		{name: "different", s1: `{"a":1}`, s2: `{"a":2}`, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, AuthPolicyStringsEquivalent(tt.s1, tt.s2))
		})
	}
}

func TestUnitJSONStringsEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		{name: "equivalent", s1: `{"a":1,"b":2}`, s2: `{ "b": 2, "a": 1 }`, want: true},
		{name: "different", s1: `{"a":1}`, s2: `{"a":2}`, want: false},
		{name: "invalid", s1: `{`, s2: `{}`, want: false},
		{name: "empty invalid", s1: "", s2: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, JSONStringsEqual(tt.s1, tt.s2))
		})
	}
}

func TestUnitJSONBytesEqual(t *testing.T) {
	t.Parallel()

	require.True(t, JSONBytesEqual([]byte(`{"a":1}`), []byte(`{"a":1}`)))
	require.False(t, JSONBytesEqual([]byte(`{"a":1}`), []byte(`{"a":2}`)))
	require.False(t, JSONBytesEqual([]byte(`{`), []byte(`{}`)))
}

func TestUnitSuppressEquivalentDiffs(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixAuthPoliciesV2().Schema, map[string]interface{}{})
	require.True(t, SuppressEquivalentAuthPolicyDiffs("identities", "", "{}", d))
	require.True(t, SuppressEquivalentJSONDiffs("entities", `{"a":1}`, `{ "a": 1 }`, d))
	require.False(t, SuppressEquivalentJSONDiffs("entities", `{"a":1}`, `{"a":2}`, d))
}

func TestUnitDeserializeJSONStringToMap(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got, err := deserializeJSONStringToMap(`{"k":"v"}`)
		require.NoError(t, err)
		require.Equal(t, "v", got["k"])
	})
	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		_, err := deserializeJSONStringToMap(`{`)
		require.Error(t, err)
	})
}
