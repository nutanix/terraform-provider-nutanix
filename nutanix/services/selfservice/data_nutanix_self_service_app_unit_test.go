package selfservice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitDatsourceNutanixCalmAppSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DatsourceNutanixCalmApp().InternalValidate(DatsourceNutanixCalmApp().Schema, false))
}

func TestUnitFlattenAppSummary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pr   map[string]interface{}
		meta map[string]interface{}
		want []interface{}
	}{
		{
			name: "nil inputs",
			pr:   nil,
			meta: nil,
			want: []interface{}{map[string]interface{}{}},
		},
		{
			name: "empty maps",
			pr:   map[string]interface{}{},
			meta: map[string]interface{}{},
			want: []interface{}{map[string]interface{}{}},
		},
		{
			name: "populated",
			pr: map[string]interface{}{
				"resources": map[string]interface{}{
					"app_profile_config_reference": map[string]interface{}{"name": "Default"},
					"app_blueprint_reference":      map[string]interface{}{"name": "bp-1"},
				},
			},
			meta: map[string]interface{}{
				"project_reference": map[string]interface{}{"name": "proj"},
				"owner_reference":   map[string]interface{}{"name": "owner"},
				"creation_time":     "2024-01-02T03:04:05Z",
				"last_update_time":  "2024-01-03T03:04:05Z",
				"uuid":              "app-uuid-1",
			},
			want: []interface{}{map[string]interface{}{
				"application_profile": "Default",
				"blueprint":           "bp-1",
				"project":             "proj",
				"owner":               "owner",
				"created_on":          "2024-01-02T03:04:05Z",
				"last_updated_on":     "2024-01-03T03:04:05Z",
				"application_uuid":    "app-uuid-1",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAppSummary(tt.pr, tt.meta))
		})
	}
}

func TestUnitFlattenActions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pr   map[string]interface{}
		want []interface{}
	}{
		{name: "nil", pr: nil, want: nil},
		{name: "empty", pr: map[string]interface{}{}, want: nil},
		{
			name: "populated",
			pr: map[string]interface{}{
				"resources": map[string]interface{}{
					"action_list": []interface{}{
						map[string]interface{}{
							"name":        "action_start",
							"uuid":        "act-1",
							"description": "start app",
						},
						map[string]interface{}{
							"name":        "custom_restart",
							"uuid":        "act-2",
							"description": "restart",
						},
					},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"name":        "start",
					"uuid":        "act-1",
					"description": "start app",
				},
				map[string]interface{}{
					"name":        "custom restart",
					"uuid":        "act-2",
					"description": "restart",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenActions(tt.pr))
		})
	}
}
