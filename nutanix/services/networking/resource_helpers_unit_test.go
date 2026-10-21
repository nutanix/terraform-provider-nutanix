package networking

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFlattenAddressEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.IPAddressBlock
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*v3.IPAddressBlock{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []*v3.IPAddressBlock{{
				IPAddress:    utils.StringPtr("10.0.0.0"),
				PrefixLength: utils.Int64Ptr(24),
			}},
			want: []map[string]interface{}{{
				"ip":            utils.StringPtr("10.0.0.0"),
				"prefix_length": utils.Int64Ptr(24),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAddressEntry(tt.in))
		})
	}
}

func TestUnitExpandAddressEntry(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixAddressGroup().Schema, map[string]interface{}{
			"name": "ag",
		})
		got, err := expandAddressEntry(d)
		require.NoError(t, err)
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixAddressGroup().Schema, map[string]interface{}{
			"name": "ag",
			"ip_address_block_list": []interface{}{
				map[string]interface{}{
					"ip":            "10.0.0.0",
					"prefix_length": 24,
				},
			},
		})
		got, err := expandAddressEntry(d)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "10.0.0.0", *got[0].IPAddress)
		require.Equal(t, int64(24), *got[0].PrefixLength)
	})
}

func TestUnitExpandUpdateReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		uuid string
		kind string
		want *v3.Reference
	}{
		{name: "empty uuid", uuid: "", kind: "subnet", want: nil},
		{
			name: "populated",
			uuid: "sub-1",
			kind: "subnet",
			want: &v3.Reference{
				Kind: utils.StringPtr("subnet"),
				UUID: utils.StringPtr("sub-1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandUpdateReference(tt.uuid, tt.kind))
		})
	}
}

func TestUnitFlattenExtSubnetListStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.ExternalSubnetList
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.ExternalSubnetList{}, want: nil},
		{
			name: "populated",
			in: []*v3.ExternalSubnetList{{
				ExternalSubnetReference: &v3.Reference{
					Kind: utils.StringPtr("subnet"),
					UUID: utils.StringPtr("ext-sub"),
				},
				ExternalIPList: []*string{utils.StringPtr("1.2.3.4")},
				ActiveGatewayNode: &v3.ActiveGatewayNode{
					IPAddress: utils.StringPtr("10.0.0.1"),
					HostReference: &v3.Reference{
						Kind: utils.StringPtr("host"),
						UUID: utils.StringPtr("host-1"),
					},
				},
			}},
			want: []map[string]interface{}{{
				"external_subnet_reference": map[string]interface{}{
					"kind": "subnet",
					"uuid": "ext-sub",
				},
				"active_gateway_node": []interface{}{
					map[string]interface{}{
						"host_reference": map[string]interface{}{
							"kind": "host",
							"uuid": "host-1",
						},
						"ip_address": utils.StringPtr("10.0.0.1"),
					},
				},
				"external_ip_list": []*string{utils.StringPtr("1.2.3.4")},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExtSubnetListStatus(tt.in))
		})
	}
}

func TestUnitFlattenActiveGatewayNode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.ActiveGatewayNode
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.ActiveGatewayNode{
				IPAddress: utils.StringPtr("10.0.0.1"),
				HostReference: &v3.Reference{
					Kind: utils.StringPtr("host"),
					UUID: utils.StringPtr("host-1"),
					Name: utils.StringPtr("h1"),
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"host_reference": map[string]interface{}{
						"kind": "host",
						"uuid": "host-1",
						"name": "h1",
					},
					"ip_address": utils.StringPtr("10.0.0.1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenActiveGatewayNode(tt.in))
		})
	}
}

func TestUnitExpandExternallyRoutablePL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []*v3.ExternallyRoutablePrefixList
	}{
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{"ip": "10.1.0.0", "prefix_length": 16},
			},
			want: []*v3.ExternallyRoutablePrefixList{{
				IP:           utils.StringPtr("10.1.0.0"),
				PrefixLength: utils.IntPtr(16),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandExternallyRoutablePL(tt.in))
		})
	}
}

func TestUnitExpandCommonDNSIPList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []*v3.CommonDomainNameServerIPList
	}{
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{"ip": "8.8.8.8"},
			},
			want: []*v3.CommonDomainNameServerIPList{{
				IP: utils.StringPtr("8.8.8.8"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandCommonDNSIPList(tt.in))
		})
	}
}

func TestUnitFlattenExtSubnetList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.ExternalSubnetList
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.ExternalSubnetList{}, want: nil},
		{
			name: "populated",
			in: []*v3.ExternalSubnetList{{
				ExternalSubnetReference: &v3.Reference{
					Kind: utils.StringPtr("subnet"),
					UUID: utils.StringPtr("ext-1"),
				},
			}},
			want: []map[string]interface{}{{
				"external_subnet_reference": map[string]interface{}{
					"kind": "subnet",
					"uuid": "ext-1",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExtSubnetList(tt.in))
		})
	}
}

func TestUnitFlattenExtRoutableList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.ExternallyRoutablePrefixList
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.ExternallyRoutablePrefixList{}, want: nil},
		{
			name: "populated",
			in: []*v3.ExternallyRoutablePrefixList{{
				IP:           utils.StringPtr("10.1.0.0"),
				PrefixLength: utils.IntPtr(16),
			}},
			want: []map[string]interface{}{{
				"ip":            utils.StringPtr("10.1.0.0"),
				"prefix_length": utils.IntPtr(16),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExtRoutableList(tt.in))
		})
	}
}

func TestUnitFlattenCommonDNSIPList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.CommonDomainNameServerIPList
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.CommonDomainNameServerIPList{}, want: nil},
		{
			name: "populated",
			in:   []*v3.CommonDomainNameServerIPList{{IP: utils.StringPtr("8.8.8.8")}},
			want: []map[string]interface{}{{"ip": utils.StringPtr("8.8.8.8")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCommonDNSIPList(tt.in))
		})
	}
}

func TestUnitExpandStaticRouteList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []*v3.StaticRoutesList
	}{
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"destination":                    "10.0.0.0/24",
					"external_subnet_reference_uuid": "sub-1",
					"vpn_connection_reference_uuid":  "",
				},
			},
			want: []*v3.StaticRoutesList{{
				Destination: utils.StringPtr("10.0.0.0/24"),
				NextHop: &v3.NextHop{
					ExternalSubnetReference: &v3.Reference{
						Kind: utils.StringPtr("subnet"),
						UUID: utils.StringPtr("sub-1"),
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandStaticRouteList(tt.in))
		})
	}
}

func TestUnitFlattenStaticRouteList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.StaticRoutesList
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.StaticRoutesList{}, want: nil},
		{
			name: "populated",
			in: []*v3.StaticRoutesList{{
				Destination: utils.StringPtr("10.0.0.0/24"),
				NextHop: &v3.NextHop{
					ExternalSubnetReference: &v3.Reference{
						UUID: utils.StringPtr("sub-1"),
					},
				},
			}},
			want: []map[string]interface{}{{
				"destination":                    utils.StringPtr("10.0.0.0/24"),
				"external_subnet_reference_uuid": utils.StringPtr("sub-1"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStaticRouteList(tt.in))
		})
	}
}

func TestUnitExpandDefaultRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want *v3.NextHop
	}{
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{
				map[string]interface{}{
					"external_subnet_reference_uuid": "sub-def",
				},
			},
			want: &v3.NextHop{
				ExternalSubnetReference: &v3.Reference{
					Kind: utils.StringPtr("subnet"),
					UUID: utils.StringPtr("sub-def"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandDefaultRoute(tt.in))
		})
	}
}

func TestUnitFlattendefaultRouteList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.NextHop
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: &v3.NextHop{
				ExternalSubnetReference: &v3.Reference{UUID: utils.StringPtr("sub-def")},
				VpnConnectionReference:  &v3.Reference{UUID: utils.StringPtr("vpn-1")},
			},
			want: []map[string]interface{}{{
				"external_subnet_reference_uuid": utils.StringPtr("sub-def"),
				"vpn_connection_reference_uuid":  utils.StringPtr("vpn-1"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattendefaultRouteList(tt.in))
		})
	}
}
