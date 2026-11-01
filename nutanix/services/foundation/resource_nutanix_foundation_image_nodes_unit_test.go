package foundation

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFoundationImageNodesSchema(t *testing.T) {
	t.Parallel()

	r := ResourceFoundationImageNodes()
	require.NoError(t, r.InternalValidate(r.Schema, true))
}

func TestUnitFoundationImageNodesReadAndDelete(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceFoundationImageNodes().Schema, map[string]interface{}{
		"cvm_gateway":        "10.0.0.1",
		"hypervisor_netmask": "255.255.255.0",
		"cvm_netmask":        "255.255.255.0",
		"blocks": []interface{}{
			map[string]interface{}{
				"block_id": "block-1",
				"nodes": []interface{}{
					map[string]interface{}{
						"node_position":       "A",
						"hypervisor_hostname": "host1",
						"hypervisor_ip":       "10.0.0.10",
						"cvm_ip":              "10.0.0.11",
						"ipmi_ip":             "10.0.0.12",
					},
				},
			},
		},
	})
	require.Nil(t, resourceFoundationImageNodesRead(t.Context(), d, nil))
	require.Nil(t, resourceFoundationImageNodesDelete(t.Context(), d, nil))
}

func TestUnitExpandTests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     map[string]interface{}
		want    *foundation.Tests
		wantErr bool
	}{
		{name: "empty", raw: map[string]interface{}{}, want: nil},
		{
			name: "populated",
			raw: map[string]interface{}{
				"tests": []interface{}{
					map[string]interface{}{
						"run_ncc":      true,
						"run_syscheck": false,
					},
				},
			},
			// Production mapping swaps field names into the SDK struct.
			want: &foundation.Tests{
				RunSyscheck: utils.BoolPtr(true),
				RunNcc:      utils.BoolPtr(false),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceFoundationImageNodes().Schema, tt.raw)
			got, err := expandTests(d)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitExpandEosMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  map[string]interface{}
		want *foundation.EosMetadata
	}{
		{name: "empty", raw: map[string]interface{}{}, want: nil},
		{
			name: "populated",
			raw: map[string]interface{}{
				"eos_metadata": []interface{}{
					map[string]interface{}{
						"config_id":    "cfg-1",
						"account_name": []interface{}{"acct-a", "acct-b"},
						"email":        "ops@example.test",
					},
				},
			},
			want: &foundation.EosMetadata{
				ConfigID:    "cfg-1",
				AccountName: []string{"acct-a", "acct-b"},
				Email:       "ops@example.test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceFoundationImageNodes().Schema, tt.raw)
			got, err := expandEosMetadata(d)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitExpandFcSetting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     map[string]interface{}
		want    *foundation.FcSettings
		wantErr bool
	}{
		{name: "empty", raw: map[string]interface{}{}, want: nil},
		{
			name: "missing fc_metadata",
			raw: map[string]interface{}{
				"fc_settings": []interface{}{
					map[string]interface{}{
						"foundation_central": true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "populated",
			raw: map[string]interface{}{
				"fc_settings": []interface{}{
					map[string]interface{}{
						"foundation_central": true,
						"fc_metadata": []interface{}{
							map[string]interface{}{
								"fc_ip":   "10.1.1.1",
								"api_key": "key-1",
							},
						},
					},
				},
			},
			want: &foundation.FcSettings{
				FoundationCentral: utils.BoolPtr(true),
				FcMetadata: &foundation.FcMetadata{
					FcIP:   "10.1.1.1",
					APIKey: "key-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceFoundationImageNodes().Schema, tt.raw)
			got, err := expandFcSetting(d)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitExpandHypervisor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *foundation.Hypervisor
	}{
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"checksum": "abc",
					"filename": "kvm.iso",
				},
			},
			want: &foundation.Hypervisor{
				Checksum: "abc",
				Filename: "kvm.iso",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandHypervisor(tt.in))
		})
	}
}

func TestUnitExpandHypervisorIso(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want foundation.HypervisorIso
	}{
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"kvm": []interface{}{
						map[string]interface{}{
							"checksum": "abc",
							"filename": "kvm.iso",
						},
					},
				},
			},
			want: foundation.HypervisorIso{
				Kvm: &foundation.Hypervisor{
					Checksum: "abc",
					Filename: "kvm.iso",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandHypervisorIso(tt.in))
		})
	}
}

func TestUnitExpandVswitches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []*foundation.Vswitches
	}{
		{name: "empty", in: []interface{}{}, want: []*foundation.Vswitches{}},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"lacp":         "active",
					"bond_mode":    "balance-slb",
					"mtu":          1500,
					"name":         "br0",
					"other_config": []interface{}{"opt1"},
					"uplinks":      []interface{}{"eth0", "eth1"},
				},
			},
			want: []*foundation.Vswitches{{
				Lacp:        "active",
				BondMode:    "balance-slb",
				Mtu:         utils.Int64Ptr(1500),
				Name:        "br0",
				OtherConfig: []string{"opt1"},
				Uplinks:     []string{"eth0", "eth1"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandVswitches(tt.in))
		})
	}
}

func TestUnitExpandUcsmParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want *foundation.UcsmParams
	}{
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"native_vlan":        true,
					"keep_ucsm_settings": false,
					"mac_pool":           "pool-1",
					"vlan_name":          "vlan-a",
				},
			},
			want: &foundation.UcsmParams{
				NativeVlan:       utils.BoolPtr(true),
				KeepUcsmSettings: utils.BoolPtr(false),
				MacPool:          "pool-1",
				VlanName:         "vlan-a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandUcsmParams(tt.in))
		})
	}
}

func TestUnitExpandCluster(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  map[string]interface{}
		want []*foundation.Clusters
	}{
		{name: "empty", raw: map[string]interface{}{}, want: []*foundation.Clusters{}},
		{
			name: "populated",
			raw: map[string]interface{}{
				"clusters": []interface{}{
					map[string]interface{}{
						"enable_ns":               true,
						"backplane_subnet":        "10.2.0.0",
						"cluster_init_successful": false,
						"backplane_netmask":       "255.255.255.0",
						"redundancy_factor":       2,
						"backplane_vlan":          "100",
						"cluster_name":            "cls-1",
						"cluster_external_ip":     "10.3.0.1",
						"cvm_ntp_servers":         "ntp.example",
						"single_node_cluster":     false,
						"cluster_members":         []interface{}{"10.0.0.11"},
						"cvm_dns_servers":         "8.8.8.8",
						"cluster_init_now":        true,
						"hypervisor_ntp_servers":  "ntp.example",
						"timezone":                "UTC",
					},
				},
			},
			want: []*foundation.Clusters{{
				EnableNs:              utils.BoolPtr(true),
				BackplaneSubnet:       "10.2.0.0",
				ClusterInitSuccessful: utils.BoolPtr(false),
				BackplaneNetmask:      "255.255.255.0",
				RedundancyFactor:      utils.Int64Ptr(2),
				BackplaneVlan:         "100",
				ClusterName:           "cls-1",
				ClusterExternalIP:     utils.StringPtr("10.3.0.1"),
				CvmNtpServers:         "ntp.example",
				SingleNodeCluster:     utils.BoolPtr(false),
				ClusterMembers:        []string{"10.0.0.11"},
				CvmDNSServers:         "8.8.8.8",
				ClusterInitNow:        utils.BoolPtr(true),
				HypervisorNtpServers:  "ntp.example",
				Timezone:              "UTC",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceFoundationImageNodes().Schema, tt.raw)
			got, err := expandCluster(d)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitExpandNodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []*foundation.Node
	}{
		{name: "empty", in: []interface{}{}, want: []*foundation.Node{}},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"ipmi_netmask":              "255.255.255.0",
					"ipmi_gateway":              "10.0.0.1",
					"ipv6_address":              "fe80::1",
					"node_position":             "A",
					"image_delay":               5,
					"hypervisor_hostname":       "host1",
					"cvm_gb_ram":                32,
					"device_hint":               "hint",
					"bond_mode":                 "active-backup",
					"bond_uplinks":              []interface{}{"eth0"},
					"rdma_passthrough":          true,
					"cluster_id":                "c1",
					"ucsm_node_serial":          "ucs-1",
					"hypervisor_ip":             "10.0.0.10",
					"node_serial":               "sn-1",
					"ipmi_configure_now":        true,
					"image_successful":          true,
					"ipv6_interface":            "eth0",
					"cvm_num_vcpus":             8,
					"ipmi_mac":                  "aa:bb",
					"rdma_mac_addr":             "cc:dd",
					"current_network_interface": "eth0",
					"hypervisor":                "kvm",
					"bond_lacp_rate":            "fast",
					"image_now":                 true,
					"ucsm_managed_mode":         "mode",
					"ipmi_ip":                   "10.0.0.12",
					"current_cvm_vlan_tag":      10,
					"cvm_ip":                    "10.0.0.11",
					"exlude_boot_serial":        "boot",
					"mitigate_low_boot_space":   true,
					"ipmi_user":                 "admin",
					"ipmi_password":             "secret",
					"ucsm_params": []interface{}{
						map[string]interface{}{
							"mac_pool":  "pool-1",
							"vlan_name": "vlan-a",
						},
					},
					"vswitches": []interface{}{
						map[string]interface{}{
							"name": "br0",
						},
					},
				},
			},
			want: []*foundation.Node{{
				IPMINetmask:             "255.255.255.0",
				IPMIGateway:             "10.0.0.1",
				Ipv6Address:             "fe80::1",
				NodePosition:            "A",
				ImageDelay:              utils.Int64Ptr(5),
				HypervisorHostname:      "host1",
				CvmGbRAM:                utils.Int64Ptr(32),
				DeviceHint:              "hint",
				BondMode:                "active-backup",
				BondUplinks:             []string{"eth0"},
				RdmaPassthrough:         utils.BoolPtr(true),
				ClusterID:               "c1",
				UcsmNodeSerial:          "ucs-1",
				HypervisorIP:            "10.0.0.10",
				NodeSerial:              "sn-1",
				IpmiConfigureNow:        utils.BoolPtr(true),
				ImageSuccessful:         utils.BoolPtr(true),
				Ipv6Interface:           "eth0",
				CvmNumVcpus:             utils.Int64Ptr(8),
				IpmiMac:                 "aa:bb",
				RdmaMacAddr:             "cc:dd",
				CurrentNetworkInterface: "eth0",
				Hypervisor:              "kvm",
				BondLacpRate:            "fast",
				ImageNow:                utils.BoolPtr(true),
				UcsmManagedMode:         "mode",
				IpmiIP:                  "10.0.0.12",
				CurrentCvmVlanTag:       utils.Int64Ptr(10),
				CvmIP:                   "10.0.0.11",
				ExludeBootSerial:        "boot",
				MitigateLowBootSpace:    utils.BoolPtr(true),
				IpmiUser:                "admin",
				IpmiPassword:            "secret",
				UcsmParams: &foundation.UcsmParams{
					MacPool:  "pool-1",
					VlanName: "vlan-a",
				},
				Vswitches: []*foundation.Vswitches{{
					Name: "br0",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandNodes(tt.in))
		})
	}
}

func TestUnitExpandBlocks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  map[string]interface{}
		want []*foundation.Block
	}{
		{name: "empty", raw: map[string]interface{}{}, want: nil},
		{
			name: "populated",
			raw: map[string]interface{}{
				"blocks": []interface{}{
					map[string]interface{}{
						"block_id": "block-1",
						"nodes": []interface{}{
							map[string]interface{}{
								"node_position":       "A",
								"hypervisor_hostname": "host1",
								"hypervisor_ip":       "10.0.0.10",
								"cvm_ip":              "10.0.0.11",
								"ipmi_ip":             "10.0.0.12",
								"image_now":           true,
							},
						},
					},
				},
			},
			want: []*foundation.Block{{
				BlockID: "block-1",
				Nodes: []*foundation.Node{{
					NodePosition:       "A",
					HypervisorHostname: "host1",
					HypervisorIP:       "10.0.0.10",
					CvmIP:              "10.0.0.11",
					IpmiIP:             "10.0.0.12",
					ImageNow:           utils.BoolPtr(true),
					UcsmParams:         nil,
					Vswitches:          []*foundation.Vswitches{},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceFoundationImageNodes().Schema, tt.raw)
			got, err := expandBlocks(d)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitCollectIndividualErrorDiagnostics(t *testing.T) {
	t.Parallel()

	okPct := float64(AggregatePercentComplete)
	failPct := float64(50)

	tests := []struct {
		name string
		in   *foundation.ImageNodesProgressResponse
		want int
	}{
		{
			name: "empty",
			in:   &foundation.ImageNodesProgressResponse{},
			want: 0,
		},
		{
			name: "success",
			in: &foundation.ImageNodesProgressResponse{
				Nodes: []*foundation.NodeProgress{{
					CvmIP:           "10.0.0.11",
					Status:          "COMPLETED",
					PercentComplete: &okPct,
				}},
				Clusters: []*foundation.ClusterProgress{{
					ClusterName:     "cls-1",
					Status:          "COMPLETED",
					PercentComplete: &okPct,
				}},
			},
			want: 0,
		},
		{
			name: "populated failures",
			in: &foundation.ImageNodesProgressResponse{
				Nodes: []*foundation.NodeProgress{{
					CvmIP:           "10.0.0.11",
					Status:          "FAILED",
					Messages:        []string{"node boom"},
					PercentComplete: &failPct,
				}},
				Clusters: []*foundation.ClusterProgress{{
					ClusterName:     "cls-1",
					Status:          "FAILED",
					Messages:        []string{"cluster boom"},
					PercentComplete: &failPct,
				}},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := collectIndividualErrorDiagnostics(tt.in)
			require.Len(t, got, tt.want)
			for _, d := range got {
				require.Equal(t, diag.Error, d.Severity)
			}
		})
	}
}
