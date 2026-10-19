package foundation

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFoundationDiscoverNodesSchema(t *testing.T) {
	t.Parallel()

	r := DataSourceFoundationDiscoverNodes()
	require.NoError(t, r.InternalValidate(r.Schema, false))
}

func TestUnitFlattenDiscoveredNodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []foundation.DiscoveredNode
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []foundation.DiscoveredNode{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []foundation.DiscoveredNode{{
				FoundationVersion:       "5.0",
				Ipv6Address:             "fe80::1",
				NodeUUID:                "uuid-1",
				CurrentNetworkInterface: "eth0",
				NodePosition:            "A",
				Hypervisor:              "kvm",
				Configured:              utils.BoolPtr(true),
				NosVersion:              "6.8",
				ClusterID:               float64(42),
				CurrentCvmVlanTag:       float64(10),
				HypervisorVersion:       "8.2",
				SvmIP:                   "10.0.0.2",
				Model:                   "NX",
				NodeSerial:              "SN1",
			}},
			want: []map[string]interface{}{{
				"foundation_version":        "5.0",
				"ipv6_address":              "fe80::1",
				"node_uuid":                 "uuid-1",
				"current_network_interface": "eth0",
				"node_position":             "A",
				"hypervisor":                "kvm",
				"configured":                utils.BoolPtr(true),
				"nos_version":               "6.8",
				"node_serial":               "SN1",
				"cluster_id":                "42",
				"current_cvm_vlan_tag":      "10",
				"hypervisor_version":        "8.2",
				"svm_ip":                    "10.0.0.2",
				"model":                     "NX",
			}},
		},
		{
			name: "non-float cluster fields",
			in: []foundation.DiscoveredNode{{
				ClusterID:         "not-a-float",
				CurrentCvmVlanTag: true,
			}},
			want: []map[string]interface{}{{
				"foundation_version":        "",
				"ipv6_address":              "",
				"node_uuid":                 "",
				"current_network_interface": "",
				"node_position":             "",
				"hypervisor":                "",
				"configured":                (*bool)(nil),
				"nos_version":               "",
				"node_serial":               "",
				"cluster_id":                "",
				"current_cvm_vlan_tag":      "",
				"hypervisor_version":        "",
				"svm_ip":                    "",
				"model":                     "",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDiscoveredNodes(tt.in))
		})
	}
}
