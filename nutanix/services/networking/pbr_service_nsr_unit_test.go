package networking

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitExpandSourDest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *v3.PbrSourDest
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"address_type":  "SUBNET",
					"subnet_ip":     "10.0.0.0",
					"prefix_length": 24,
				},
			},
			want: &v3.PbrSourDest{
				AddressType: utils.StringPtr("SUBNET"),
				IPSubnet: &v3.PbrIPSubnet{
					IP:           utils.StringPtr("10.0.0.0"),
					PrefixLength: utils.IntPtr(24),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandSourDest(tt.in))
		})
	}
}

func TestUnitExpandProtocolParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *v3.PbrProtocolParams
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated icmp",
			in: []interface{}{
				map[string]interface{}{
					"protocol_number": "1",
					"icmp": []interface{}{
						map[string]interface{}{"icmp_code": 0, "icmp_type": 8},
					},
					"udp": []interface{}{},
					"tcp": []interface{}{},
				},
			},
			want: &v3.PbrProtocolParams{
				ProtocolNumber: utils.IntPtr(1),
				Icmp: &v3.PbrIcmp{
					IcmpCode: utils.IntPtr(0),
					IcmpType: utils.IntPtr(8),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandProtocolParameters(tt.in))
		})
	}
}

func TestUnitExpandRouteProtocol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want *v3.PortRangeList
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"source_port_range_list": []interface{}{
						map[string]interface{}{"start_port": 80, "end_port": 80},
					},
					"destination_port_range_list": []interface{}{
						map[string]interface{}{"start_port": 443, "end_port": 443},
					},
				},
			},
			want: &v3.PortRangeList{
				SourcePortRangeList: []*v3.PortRange{{
					StartPort: utils.Int64Ptr(80),
					EndPort:   utils.Int64Ptr(80),
				}},
				DestinationPortRangeList: []*v3.PortRange{{
					StartPort: utils.Int64Ptr(443),
					EndPort:   utils.Int64Ptr(443),
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandRouteProtocol(tt.in))
		})
	}
}

func TestUnitFlattenSourceDest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.PbrSourDest
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.PbrSourDest{
				AddressType: utils.StringPtr("SUBNET"),
				IPSubnet: &v3.PbrIPSubnet{
					IP:           utils.StringPtr("10.0.0.0"),
					PrefixLength: utils.IntPtr(24),
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"address_type":  utils.StringPtr("SUBNET"),
					"subnet_ip":     utils.StringPtr("10.0.0.0"),
					"prefix_length": utils.IntPtr(24),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSourceDest(tt.in))
		})
	}
}

func TestUnitFlattenProtocolParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.PbrProtocolParams
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.PbrProtocolParams{
				ProtocolNumber: utils.IntPtr(6),
				Icmp: &v3.PbrIcmp{
					IcmpCode: utils.IntPtr(0),
					IcmpType: utils.IntPtr(8),
				},
			},
			want: []map[string]interface{}{{
				"protocol_number": "6",
				"icmp": []interface{}{
					map[string]interface{}{
						"icmp_code": utils.IntPtr(0),
						"icmp_type": utils.IntPtr(8),
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenProtocolParams(tt.in))
		})
	}
}

func TestUnitFlattenRouteSpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.PortRangeList
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.PortRangeList{
				SourcePortRangeList: []*v3.PortRange{{
					StartPort: utils.Int64Ptr(80),
					EndPort:   utils.Int64Ptr(80),
				}},
			},
			want: []interface{}{
				map[string]interface{}{
					"source_port_range_list": []map[string]interface{}{{
						"end_port":   utils.Int64Ptr(80),
						"start_port": utils.Int64Ptr(80),
					}},
					"destination_port_range_list": []map[string]interface{}(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRouteSpec(tt.in))
		})
	}
}

func TestUnitFlattenIcmpPbr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.PbrIcmp
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &v3.PbrIcmp{IcmpCode: utils.IntPtr(0), IcmpType: utils.IntPtr(8)},
			want: []interface{}{
				map[string]interface{}{
					"icmp_code": utils.IntPtr(0),
					"icmp_type": utils.IntPtr(8),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIcmpPbr(tt.in))
		})
	}
}

func TestUnitFlattenPortRangeList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.PortRange
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.PortRange{}, want: nil},
		{
			name: "populated",
			in: []*v3.PortRange{{
				StartPort: utils.Int64Ptr(80),
				EndPort:   utils.Int64Ptr(90),
			}},
			want: []map[string]interface{}{{
				"end_port":   utils.Int64Ptr(90),
				"start_port": utils.Int64Ptr(80),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPortRangeList(tt.in))
		})
	}
}

func TestUnitIsValidProtocol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "empty", in: "", want: false},
		{name: "ALL", in: "ALL", want: true},
		{name: "ICMP", in: "ICMP", want: true},
		{name: "TCP", in: "TCP", want: true},
		{name: "UDP", in: "UDP", want: true},
		{name: "invalid", in: "HTTP", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, IsValidProtocol(tt.in))
		})
	}
}

func TestUnitFlattenServiceEntry(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenServiceEntry(&v3.ServiceGroupInput{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		in := &v3.ServiceGroupInput{
			ServiceList: []*v3.ServiceListEntry{{
				Protocol: utils.StringPtr("TCP"),
				TCPPortRangeList: []*v3.PortRange{{
					StartPort: utils.Int64Ptr(80),
					EndPort:   utils.Int64Ptr(80),
				}},
			}},
		}
		want := []map[string]interface{}{{
			"protocol": "TCP",
			"tcp_port_range_list": []map[string]interface{}{{
				"end_port":   int64(80),
				"start_port": int64(80),
			}},
		}}
		require.Equal(t, want, flattenServiceEntry(in))
	})
}

func TestUnitExpandServiceEntry(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixServiceGroup().Schema, map[string]interface{}{
			"name": "sg",
		})
		got, err := expandServiceEntry(d)
		require.NoError(t, err)
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixServiceGroup().Schema, map[string]interface{}{
			"name": "sg",
			"service_list": []interface{}{
				map[string]interface{}{
					"protocol": "TCP",
					"tcp_port_range_list": []interface{}{
						map[string]interface{}{"start_port": 80, "end_port": 80},
					},
				},
			},
		})
		got, err := expandServiceEntry(d)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "TCP", *got[0].Protocol)
		require.Len(t, got[0].TCPPortRangeList, 1)
		require.Equal(t, int64(80), *got[0].TCPPortRangeList[0].StartPort)
	})

	t.Run("invalid protocol", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixServiceGroup().Schema, map[string]interface{}{
			"name": "sg",
			"service_list": []interface{}{
				map[string]interface{}{"protocol": "HTTP"},
			},
		})
		got, err := expandServiceEntry(d)
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestUnitExpandReferencesList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []*v3.Reference
	}{
		{name: "empty", in: []interface{}{}, want: []*v3.Reference{}},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{"kind": "address_group", "uuid": "ag-1", "name": "ag"},
			},
			want: []*v3.Reference{{
				Kind: utils.StringPtr("address_group"),
				UUID: utils.StringPtr("ag-1"),
				Name: utils.StringPtr("ag"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandReferencesList(tt.in))
		})
	}
}

func TestUnitExpandFilterParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string][]string
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: map[string][]string{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in:   map[string][]string{"Environment": {"Prod"}},
			want: []map[string]interface{}{{
				"name":   "Environment",
				"values": []string{"Prod"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandFilterParams(tt.in))
		})
	}
}

func TestUnitExpandPortRangeList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []*v3.PortRange
	}{
		{name: "empty", in: []interface{}{}, want: []*v3.PortRange{}},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{"start_port": 22, "end_port": 22},
			},
			want: []*v3.PortRange{{
				StartPort: utils.Int64Ptr(22),
				EndPort:   utils.Int64Ptr(22),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandPortRangeList(tt.in))
		})
	}
}

func TestUnitExpandIcmpTypeCodeList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []*v3.NetworkRuleIcmpTypeCodeList
	}{
		{name: "empty", in: []interface{}{}, want: []*v3.NetworkRuleIcmpTypeCodeList{}},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{"code": "0", "type": "8"},
			},
			want: []*v3.NetworkRuleIcmpTypeCodeList{{
				Code: utils.Int64Ptr(0),
				Type: utils.Int64Ptr(8),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandIcmpTypeCodeList(tt.in))
		})
	}
}

func TestUnitFlattenNetworkRuleList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.NetworkRule
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*v3.NetworkRule{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []*v3.NetworkRule{{
				Protocol:              utils.StringPtr("TCP"),
				PeerSpecificationType: utils.StringPtr("IP_SUBNET"),
				ExpirationTime:        utils.StringPtr("2025-01-01T00:00:00Z"),
				IPSubnet: &v3.IPSubnet{
					IP:           utils.StringPtr("10.0.0.0"),
					PrefixLength: utils.Int64Ptr(24),
				},
				TCPPortRangeList: []*v3.PortRange{{
					StartPort: utils.Int64Ptr(443),
					EndPort:   utils.Int64Ptr(443),
				}},
			}},
			want: []map[string]interface{}{{
				"protocol":                "TCP",
				"ip_subnet":               "10.0.0.0",
				"ip_subnet_prefix_length": "24",
				"tcp_port_range_list": []map[string]interface{}{{
					"end_port":   int64(443),
					"start_port": int64(443),
				}},
				"peer_specification_type": "IP_SUBNET",
				"expiration_time":         "2025-01-01T00:00:00Z",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNetworkRuleList(tt.in))
		})
	}
}
