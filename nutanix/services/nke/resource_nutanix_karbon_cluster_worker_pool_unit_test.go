package nke

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/karbon"
)

func TestUnitResourceNutanixKarbonWorkerNodePoolSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, ResourceNutanixKarbonWorkerNodePool().InternalValidate(ResourceNutanixKarbonWorkerNodePool().Schema, true))
}

func TestUnitFlattenAHVNodePoolConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *karbon.ClusterNodePoolAHVConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &karbon.ClusterNodePoolAHVConfig{},
			want: []map[string]interface{}{{
				"cpu":                        int64(0),
				"disk_mib":                   int64(0),
				"memory_mib":                 int64(0),
				"network_uuid":               "",
				"prism_element_cluster_uuid": "",
			}},
		},
		{
			name: "populated",
			in: &karbon.ClusterNodePoolAHVConfig{
				CPU:                     8,
				DiskMib:                 122880,
				MemoryMib:               8192,
				NetworkUUID:             "net-uuid",
				PrismElementClusterUUID: "pe-uuid",
				IscsiNetworkUUID:        "iscsi-uuid",
			},
			want: []map[string]interface{}{{
				"cpu":                        int64(8),
				"disk_mib":                   int64(122880),
				"memory_mib":                 int64(8192),
				"network_uuid":               "net-uuid",
				"prism_element_cluster_uuid": "pe-uuid",
				"iscsi_network_uuid":         "iscsi-uuid",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAHVNodePoolConfig(tt.in))
		})
	}
}
