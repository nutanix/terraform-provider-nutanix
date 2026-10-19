package objectstoresv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	objectsCommon "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/common/v1/config"
	objectsResponse "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/common/v1/response"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/objects/v4/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitObjectStoreSchemas(t *testing.T) {
	t.Parallel()

	r := ResourceNutanixObjectStoresV2()
	require.NoError(t, r.InternalValidate(r.Schema, true))
	require.NoError(t, linksSchema().Elem.(*schema.Resource).InternalValidate(
		linksSchema().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, SchemaForValuePrefixLength(32).Elem.(*schema.Resource).InternalValidate(
		SchemaForValuePrefixLength(32).Elem.(*schema.Resource).Schema, true))
}

func TestUnitExpandMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *objectsCommon.Metadata
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"owner_reference_id":   "owner-1",
					"owner_user_name":      "admin",
					"project_reference_id": "proj-1",
					"project_name":         "demo",
					"category_ids":         []interface{}{"cat-1", "cat-2"},
				},
			},
			want: &objectsCommon.Metadata{
				OwnerReferenceId:   utils.StringPtr("owner-1"),
				OwnerUserName:      utils.StringPtr("admin"),
				ProjectReferenceId: utils.StringPtr("proj-1"),
				ProjectName:        utils.StringPtr("demo"),
				CategoryIds:        []string{"cat-1", "cat-2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandMetadata(tt.in))
		})
	}
}

func TestUnitExpandIPAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
	}{
		{name: "empty", in: []interface{}{}},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"ipv4": []interface{}{
						map[string]interface{}{
							"value":         "10.0.0.1",
							"prefix_length": 24,
						},
					},
					"ipv6": []interface{}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandIPAddress(tt.in)
			if tt.name == "empty" {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.NotNil(t, got[0].Ipv4)
			require.Equal(t, "10.0.0.1", utils.StringValue(got[0].Ipv4.Value))
			require.Equal(t, 24, utils.IntValue(got[0].Ipv4.PrefixLength))
			require.Nil(t, got[0].Ipv6)
		})
	}
}

func TestUnitExpandIPv4Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
	}{
		{name: "empty", in: []interface{}{}},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"value":         "10.0.0.8",
					"prefix_length": 32,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandIPv4Address(tt.in)
			if tt.name == "empty" {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.Equal(t, "10.0.0.8", utils.StringValue(got.Value))
			require.Equal(t, 32, utils.IntValue(got.PrefixLength))
		})
	}
}

func TestUnitExpandIPv6Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
	}{
		{name: "empty", in: []interface{}{}},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"value":         "fe80::1",
					"prefix_length": 128,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandIPv6Address(tt.in)
			if tt.name == "empty" {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.Equal(t, "fe80::1", utils.StringValue(got.Value))
			require.Equal(t, 128, utils.IntValue(got.PrefixLength))
		})
	}
}

func TestUnitExpandState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want config.State
	}{
		{name: "empty", in: "", want: config.STATE_UNKNOWN},
		{name: "available", in: "OBJECT_STORE_AVAILABLE", want: config.STATE_OBJECT_STORE_AVAILABLE},
		{name: "deployment failed", in: "OBJECT_STORE_DEPLOYMENT_FAILED", want: config.STATE_OBJECT_STORE_DEPLOYMENT_FAILED},
		{name: "operation failed", in: "OBJECT_STORE_OPERATION_FAILED", want: config.STATE_OBJECT_STORE_OPERATION_FAILED},
		{name: "operation pending", in: "OBJECT_STORE_OPERATION_PENDING", want: config.STATE_OBJECT_STORE_OPERATION_PENDING},
		{name: "cert creation failed", in: "OBJECT_STORE_CERT_CREATION_FAILED", want: config.STATE_OBJECT_STORE_CERT_CREATION_FAILED},
		{name: "creating cert", in: "CREATING_OBJECT_STORE_CERT", want: config.STATE_CREATING_OBJECT_STORE_CERT},
		{name: "deletion failed", in: "OBJECT_STORE_DELETION_FAILED", want: config.STATE_OBJECT_STORE_DELETION_FAILED},
		{name: "deploying", in: "DEPLOYING_OBJECT_STORE", want: config.STATE_DEPLOYING_OBJECT_STORE},
		{name: "deleting", in: "DELETING_OBJECT_STORE", want: config.STATE_DELETING_OBJECT_STORE},
		{name: "undeployed", in: "UNDEPLOYED_OBJECT_STORE", want: config.STATE_UNDEPLOYED_OBJECT_STORE},
		{name: "unknown default", in: "NOPE", want: config.STATE_UNKNOWN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandState(tt.in)
			require.NotNil(t, got)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestUnitFlattenMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *objectsCommon.Metadata
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &objectsCommon.Metadata{},
			want: []map[string]interface{}{{
				"owner_reference_id":   (*string)(nil),
				"owner_user_name":      (*string)(nil),
				"project_reference_id": (*string)(nil),
				"project_name":         (*string)(nil),
				"category_ids":         []string(nil),
			}},
		},
		{
			name: "populated",
			in: &objectsCommon.Metadata{
				OwnerReferenceId:   utils.StringPtr("owner-1"),
				OwnerUserName:      utils.StringPtr("admin"),
				ProjectReferenceId: utils.StringPtr("proj-1"),
				ProjectName:        utils.StringPtr("demo"),
				CategoryIds:        []string{"cat-1"},
			},
			want: []map[string]interface{}{{
				"owner_reference_id":   utils.StringPtr("owner-1"),
				"owner_user_name":      utils.StringPtr("admin"),
				"project_reference_id": utils.StringPtr("proj-1"),
				"project_name":         utils.StringPtr("demo"),
				"category_ids":         []string{"cat-1"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenMetadata(tt.in))
		})
	}
}

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []objectsResponse.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []objectsResponse.ApiLink{}, want: nil},
		{
			name: "populated",
			in: []objectsResponse.ApiLink{{
				Href: utils.StringPtr("https://example.test/obj/1"),
				Rel:  utils.StringPtr("self"),
			}},
			want: []map[string]interface{}{{
				"href": utils.StringPtr("https://example.test/obj/1"),
				"rel":  utils.StringPtr("self"),
			}},
		},
		{
			name: "missing fields",
			in:   []objectsResponse.ApiLink{{}},
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

func TestUnitFlattenFloatingIPv4Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *objectsCommon.IPv4Address
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &objectsCommon.IPv4Address{
				Value:        utils.StringPtr("10.0.0.8"),
				PrefixLength: utils.IntPtr(24),
			},
			want: []map[string]interface{}{{
				"value":         utils.StringPtr("10.0.0.8"),
				"prefix_length": utils.IntPtr(24),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFloatingIPv4Address(tt.in))
		})
	}
}

func TestUnitFlattenFloatingIPv6Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *objectsCommon.IPv6Address
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &objectsCommon.IPv6Address{
				Value:        utils.StringPtr("fe80::1"),
				PrefixLength: utils.IntPtr(128),
			},
			want: []map[string]interface{}{{
				"value":         utils.StringPtr("fe80::1"),
				"prefix_length": utils.IntPtr(128),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFloatingIPv6Address(tt.in))
		})
	}
}

func TestUnitFlattenIPAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []objectsCommon.IPAddress
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []objectsCommon.IPAddress{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []objectsCommon.IPAddress{{
				Ipv4: &objectsCommon.IPv4Address{
					Value:        utils.StringPtr("10.0.0.8"),
					PrefixLength: utils.IntPtr(24),
				},
			}},
			want: []map[string]interface{}{{
				"ipv4": []map[string]interface{}{{
					"value":         utils.StringPtr("10.0.0.8"),
					"prefix_length": utils.IntPtr(24),
				}},
				"ipv6": []map[string]interface{}(nil),
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
