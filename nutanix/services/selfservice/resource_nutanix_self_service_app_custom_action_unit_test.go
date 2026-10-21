package selfservice

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNutanixCalmAppCustomActionSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, ResourceNutanixCalmAppCustomAction().InternalValidate(ResourceNutanixCalmAppCustomAction().Schema, true))
}

func TestUnitExpandCustomActionSpec(t *testing.T) {
	t.Parallel()

	action := map[string]interface{}{
		"name": "action_restart",
		"uuid": "act-uuid-1",
	}

	tests := []struct {
		name       string
		pr         map[string]interface{}
		actionName string
		wantSpec   map[string]interface{}
		wantUUID   string
	}{
		{name: "nil", pr: nil, actionName: "restart", wantSpec: nil, wantUUID: ""},
		{name: "empty", pr: map[string]interface{}{}, actionName: "restart", wantSpec: nil, wantUUID: ""},
		{
			name: "match calm name",
			pr: map[string]interface{}{
				"resources": map[string]interface{}{
					"action_list": []interface{}{action},
				},
			},
			actionName: "restart",
			wantSpec:   action,
			wantUUID:   "act-uuid-1",
		},
		{
			name: "match exact name",
			pr: map[string]interface{}{
				"resources": map[string]interface{}{
					"action_list": []interface{}{
						map[string]interface{}{
							"name": "MyAction",
							"uuid": "act-uuid-2",
						},
					},
				},
			},
			actionName: "MyAction",
			wantSpec: map[string]interface{}{
				"name": "MyAction",
				"uuid": "act-uuid-2",
			},
			wantUUID: "act-uuid-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotSpec, gotUUID := expandCustomActionSpec(tt.pr, tt.actionName)
			require.Equal(t, tt.wantSpec, gotSpec)
			require.Equal(t, tt.wantUUID, gotUUID)
		})
	}
}

func TestUnitResourceNutanixCalmAppCustomActionNoops(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixCalmAppCustomAction().Schema, map[string]interface{}{
		"app_name":    "app-1",
		"action_name": "restart",
	})

	require.Nil(t, ResourceNutanixCalmAppCustomActionRead(t.Context(), d, nil))
	require.Nil(t, ResourceNutanixCalmAppCustomActionUpdate(t.Context(), d, nil))
	require.Nil(t, ResourceNutanixCalmAppCustomActionDelete(t.Context(), d, nil))
}
