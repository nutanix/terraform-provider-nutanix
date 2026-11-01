package selfservice

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitDatsourceNutanixCalmRuntimeEditablesSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DatsourceNutanixCalmRuntimeEditables().InternalValidate(DatsourceNutanixCalmRuntimeEditables().Schema, false))
}

func TestUnitFlattenRuntimeSpec(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{"key":"value"}`)

	tests := []struct {
		name string
		in   []*selfservice.RuntimeSpec
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*selfservice.RuntimeSpec{}, want: nil},
		{
			name: "populated",
			in: []*selfservice.RuntimeSpec{{
				UUID:        utils.StringPtr("uuid-1"),
				Name:        utils.StringPtr("var-1"),
				Description: utils.StringPtr("desc"),
				Type:        utils.StringPtr("LOCAL"),
				Context:     utils.StringPtr("ctx"),
				Value:       &raw,
			}},
			want: []interface{}{map[string]interface{}{
				"uuid":        utils.StringPtr("uuid-1"),
				"name":        utils.StringPtr("var-1"),
				"description": utils.StringPtr("desc"),
				"type":        utils.StringPtr("LOCAL"),
				"context":     utils.StringPtr("ctx"),
				"value":       `{"key":"value"}`,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRuntimeSpec(tt.in))
		})
	}
}

func TestUnitFlattenRuntimeEditables(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`"v"`)
	spec := &selfservice.RuntimeSpec{
		UUID:  utils.StringPtr("u1"),
		Name:  utils.StringPtr("n1"),
		Value: &raw,
	}

	tests := []struct {
		name string
		in   *selfservice.RuntimeEditables
		want []interface{}
	}{
		{
			name: "empty",
			in:   &selfservice.RuntimeEditables{},
			want: []interface{}{map[string]interface{}{
				"action_list":         []interface{}(nil),
				"service_list":        []interface{}(nil),
				"credential_list":     []interface{}(nil),
				"substrate_list":      []interface{}(nil),
				"package_list":        []interface{}(nil),
				"task_list":           []interface{}(nil),
				"restore_config_list": []interface{}(nil),
				"variable_list":       []interface{}(nil),
				"deployment_list":     []interface{}(nil),
			}},
		},
		{
			name: "populated",
			in: &selfservice.RuntimeEditables{
				VariableList: []*selfservice.RuntimeSpec{spec},
			},
			want: []interface{}{map[string]interface{}{
				"action_list":         []interface{}(nil),
				"service_list":        []interface{}(nil),
				"credential_list":     []interface{}(nil),
				"substrate_list":      []interface{}(nil),
				"package_list":        []interface{}(nil),
				"task_list":           []interface{}(nil),
				"restore_config_list": []interface{}(nil),
				"variable_list": []interface{}{map[string]interface{}{
					"uuid":        utils.StringPtr("u1"),
					"name":        utils.StringPtr("n1"),
					"description": (*string)(nil),
					"type":        (*string)(nil),
					"context":     (*string)(nil),
					"value":       `"v"`,
				}},
				"deployment_list": []interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRuntimeEditables(tt.in))
		})
	}
}
