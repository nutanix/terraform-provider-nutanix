package selfservice

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNutanixCalmAppPatchSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, ResourceNutanixCalmAppPatch().InternalValidate(ResourceNutanixCalmAppPatch().Schema, true))
}

func TestUnitExpandPatchSpec(t *testing.T) {
	t.Parallel()

	patch := map[string]interface{}{
		"name": "scale_out",
		"uuid": "patch-uuid-1",
	}

	tests := []struct {
		name      string
		pr        map[string]interface{}
		patchName string
		wantSpec  map[string]interface{}
		wantUUID  string
	}{
		{name: "nil", pr: nil, patchName: "scale_out", wantSpec: nil, wantUUID: ""},
		{name: "empty", pr: map[string]interface{}{}, patchName: "scale_out", wantSpec: nil, wantUUID: ""},
		{
			name: "populated",
			pr: map[string]interface{}{
				"resources": map[string]interface{}{
					"patch_list": []interface{}{patch},
				},
			},
			patchName: "scale_out",
			wantSpec:  patch,
			wantUUID:  "patch-uuid-1",
		},
		{
			name: "not found",
			pr: map[string]interface{}{
				"resources": map[string]interface{}{
					"patch_list": []interface{}{patch},
				},
			},
			patchName: "missing",
			wantSpec:  nil,
			wantUUID:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotSpec, gotUUID := expandPatchSpec(tt.pr, tt.patchName)
			require.Equal(t, tt.wantSpec, gotSpec)
			require.Equal(t, tt.wantUUID, gotUUID)
		})
	}
}

func TestUnitGetAttrsListFromPatchList(t *testing.T) {
	t.Parallel()

	attrs := map[string]interface{}{"foo": "bar"}

	tests := []struct {
		name      string
		pr        map[string]interface{}
		patchName string
		want      map[string]interface{}
	}{
		{name: "nil", pr: nil, patchName: "p1", want: nil},
		{name: "empty", pr: map[string]interface{}{}, patchName: "p1", want: nil},
		{
			name: "populated",
			pr: map[string]interface{}{
				"resources": map[string]interface{}{
					"patch_list": []interface{}{
						map[string]interface{}{
							"name":       "p1",
							"attrs_list": []interface{}{attrs},
						},
					},
				},
			},
			patchName: "p1",
			want:      attrs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, getAttrsListFromPatchList(tt.pr, tt.patchName))
		})
	}
}

func TestUnitResourceNutanixCalmAppPatchReadDelete(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixCalmAppPatch().Schema, map[string]interface{}{
		"app_uuid":    "app-uuid-1",
		"patch_name":  "scale_out",
		"config_name": "config-1",
	})

	require.Nil(t, resourceNutanixCalmAppPatchRead(t.Context(), d, nil))
	require.Nil(t, resourceNutanixCalmAppPatchDelete(t.Context(), d, nil))
}
