package selfservice

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNutanixCalmAppRestoreSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, ResourceNutanixCalmAppRestore().InternalValidate(ResourceNutanixCalmAppRestore().Schema, true))
}

func TestUnitFetchRestoreActionUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		in         map[string]interface{}
		actionName string
		wantUUID   string
		wantTask   string
	}{
		{name: "nil", in: nil, actionName: "action_restore", wantUUID: "", wantTask: ""},
		{name: "empty", in: map[string]interface{}{}, actionName: "action_restore", wantUUID: "", wantTask: ""},
		{
			name: "populated",
			in: map[string]interface{}{
				"resources": map[string]interface{}{
					"action_list": []interface{}{
						map[string]interface{}{
							"name": "action_restore",
							"uuid": "restore-uuid",
							"runbook": map[string]interface{}{
								"task_definition_list": []interface{}{
									map[string]interface{}{
										"type": "CALL_CONFIG",
										"uuid": "task-uuid",
									},
								},
							},
						},
					},
				},
			},
			actionName: "action_restore",
			wantUUID:   "restore-uuid",
			wantTask:   "task-uuid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotUUID, gotTask := fetchRestoreActionUUID(tt.in, tt.actionName)
			require.Equal(t, tt.wantUUID, gotUUID)
			require.Equal(t, tt.wantTask, gotTask)
		})
	}
}

func TestUnitResourceNutanixCalmAppRestoreNoops(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixCalmAppRestore().Schema, map[string]interface{}{
		"app_name":            "app-1",
		"snapshot_uuid":       "snap-1",
		"restore_action_name": "action_restore",
	})

	require.Nil(t, resourceNutanixCalmAppRestoreRead(t.Context(), d, nil))
	require.Nil(t, resourceNutanixCalmAppRestoreUpdate(t.Context(), d, nil))
	require.Nil(t, resourceNutanixCalmAppRestoreDelete(t.Context(), d, nil))
}
