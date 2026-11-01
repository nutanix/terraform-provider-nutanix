package selfservice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitDataSourceNutanixCalmSnapshotsSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, DataSourceNutanixCalmSnapshots().InternalValidate(DataSourceNutanixCalmSnapshots().Schema, false))
}

func TestUnitFetchSubstrateReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]interface{}
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "empty", in: map[string]interface{}{}, want: ""},
		{
			name: "populated",
			in: map[string]interface{}{
				"resources": map[string]interface{}{
					"deployment_list": []interface{}{
						map[string]interface{}{
							"substrate_configuration": map[string]interface{}{
								"uuid": "sub-uuid-1",
							},
						},
					},
				},
			},
			want: "sub-uuid-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, fetchSubstrateReference(tt.in))
		})
	}
}

func TestUnitFlattenSnapshotEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []map[string]interface{}
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []map[string]interface{}{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []map[string]interface{}{{
				"api_version": "3.0",
				"spec":        map[string]interface{}{"name": "snap-spec"},
				"status": map[string]interface{}{
					"type":                     "AHV_VM",
					"name":                     "snap-1",
					"uuid":                     "snap-uuid",
					"description":              "desc",
					"action_name":              "action_snapshot",
					"recovery_point_info_list": []interface{}{map[string]interface{}{"id": "rp-1"}},
				},
				"metadata": map[string]interface{}{
					"creation_time":    "2024-01-02T03:04:05Z",
					"last_update_time": "2024-01-03T03:04:05Z",
					"spec_version":     2,
					"kind":             "app_protection",
				},
			}},
			want: []map[string]interface{}{{
				"api_version":              "3.0",
				"spec":                     map[string]interface{}{"name": "snap-spec"},
				"type":                     "AHV_VM",
				"name":                     "snap-1",
				"uuid":                     "snap-uuid",
				"description":              "desc",
				"action_name":              "action_snapshot",
				"recovery_point_info_list": []interface{}{map[string]interface{}{"id": "rp-1"}},
				"creation_time":            "2024-01-02T03:04:05Z",
				"last_update_time":         "2024-01-03T03:04:05Z",
				"spec_version":             2,
				"kind":                     "app_protection",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSnapshotEntities(tt.in))
		})
	}
}
