package selfservice

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNutanixCalmAppRecoveryPointSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, ResourceNutanixCalmAppRecoveryPoint().InternalValidate(ResourceNutanixCalmAppRecoveryPoint().Schema, true))
}

func TestUnitFetchSnapshotActionUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		in         map[string]interface{}
		actionName string
		wantUUID   string
		wantTask   string
	}{
		{name: "nil", in: nil, actionName: "action_snapshot", wantUUID: "", wantTask: ""},
		{name: "empty", in: map[string]interface{}{}, actionName: "action_snapshot", wantUUID: "", wantTask: ""},
		{
			name: "populated",
			in: map[string]interface{}{
				"resources": map[string]interface{}{
					"action_list": []interface{}{
						map[string]interface{}{
							"name": "action_snapshot",
							"uuid": "snap-uuid",
							"runbook": map[string]interface{}{
								"task_definition_list": []interface{}{
									map[string]interface{}{
										"type": "CALL_CONFIG",
										"uuid": "snap-task-uuid",
									},
								},
							},
						},
					},
				},
			},
			actionName: "action_snapshot",
			wantUUID:   "snap-uuid",
			wantTask:   "snap-task-uuid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotUUID, gotTask := fetchSnapshotActionUUID(tt.in, tt.actionName)
			require.Equal(t, tt.wantUUID, gotUUID)
			require.Equal(t, tt.wantTask, gotTask)
		})
	}
}

func TestUnitResourceNutanixCalmAppRecoveryPointNoops(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixCalmAppRecoveryPoint().Schema, map[string]interface{}{
		"app_name":    "app-1",
		"action_name": "action_snapshot",
	})

	require.Nil(t, resourceNutanixCalmAppRecoveryPointRead(t.Context(), d, nil))
	require.Nil(t, resourceNutanixCalmAppRecoveryPointUpdate(t.Context(), d, nil))
}
