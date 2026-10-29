package networking

import (
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFlattenStaticRouteSpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.StaticRouteSpec
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.StaticRouteSpec{
				Name: utils.StringPtr("sr-1"),
			},
			want: []interface{}{
				map[string]interface{}{
					"name":      utils.StringPtr("sr-1"),
					"resources": []interface{}(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStaticRouteSpec(tt.in))
		})
	}
}

func TestUnitFlattenStaticRouteStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.StaticRouteDefStatus
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.StaticRouteDefStatus{
				State: utils.StringPtr("COMPLETE"),
			},
			want: []interface{}{
				map[string]interface{}{
					"state":             utils.StringPtr("COMPLETE"),
					"resources":         []interface{}(nil),
					"execution_context": []interface{}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStaticRouteStatus(tt.in))
		})
	}
}

func TestUnitFlattenStaticRouteSpecResources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.StaticRouteResources
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.StaticRouteResources{
				StaticRoutesList: []*v3.StaticRoutesList{{
					Destination: utils.StringPtr("10.0.0.0/24"),
					NextHop: &v3.NextHop{
						ExternalSubnetReference: &v3.Reference{
							Kind: utils.StringPtr("subnet"),
							UUID: utils.StringPtr("sub-1"),
						},
					},
				}},
				DefaultRouteNexthop: &v3.NextHop{
					NexthopIPAddress: utils.StringPtr("10.0.0.1"),
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"static_routes_list": []map[string]interface{}{{
						"destination": utils.StringPtr("10.0.0.0/24"),
						"nexthop": []interface{}{
							map[string]interface{}{
								"external_subnet_reference": map[string]interface{}{
									"kind": "subnet",
									"uuid": "sub-1",
								},
							},
						},
					}},
					"default_route_nexthop": []interface{}{
						map[string]interface{}{
							"nexthop_ip_address": utils.StringPtr("10.0.0.1"),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStaticRouteSpecResources(tt.in))
		})
	}
}

func TestUnitFlattenStaticRouteStatusResources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.StaticRouteResources
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.StaticRouteResources{
				DefaultRoute: &v3.StaticRoutesList{
					Destination: utils.StringPtr("0.0.0.0/0"),
					IsActive:    utils.BoolPtr(true),
					Priority:    utils.Int64Ptr(1),
					NextHop: &v3.NextHop{
						NexthopIPAddress: utils.StringPtr("10.0.0.1"),
					},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"static_routes_list":  []map[string]interface{}(nil),
					"local_routes_list":   []map[string]interface{}(nil),
					"dynamic_routes_list": []map[string]interface{}(nil),
					"default_route": []map[string]interface{}{{
						"destination": utils.StringPtr("0.0.0.0/0"),
						"is_active":   utils.BoolPtr(true),
						"priority":    utils.Int64Ptr(1),
						"nexthop": []interface{}{
							map[string]interface{}{
								"nexthop_ip_address": utils.StringPtr("10.0.0.1"),
							},
						},
					}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStaticRouteStatusResources(tt.in))
		})
	}
}

func TestUnitFlattenSpecStaticRouteList(t *testing.T) {
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
				Destination: utils.StringPtr("10.0.0.0/8"),
				NextHop:     &v3.NextHop{NexthopIPAddress: utils.StringPtr("1.1.1.1")},
			}},
			want: []map[string]interface{}{{
				"destination": utils.StringPtr("10.0.0.0/8"),
				"nexthop": []interface{}{
					map[string]interface{}{
						"nexthop_ip_address": utils.StringPtr("1.1.1.1"),
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSpecStaticRouteList(tt.in))
		})
	}
}

func TestUnitFlattenStatusDefaultRouteList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.StaticRoutesList
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.StaticRoutesList{
				Destination: utils.StringPtr("0.0.0.0/0"),
				IsActive:    utils.BoolPtr(true),
				Priority:    utils.Int64Ptr(10),
				NextHop:     nil,
			},
			want: []map[string]interface{}{{
				"destination": utils.StringPtr("0.0.0.0/0"),
				"nexthop":     []interface{}(nil),
				"is_active":   utils.BoolPtr(true),
				"priority":    utils.Int64Ptr(10),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStatusDefaultRouteList(tt.in))
		})
	}
}

func TestUnitFlattenStatusStaticRouteList(t *testing.T) {
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
				Destination: utils.StringPtr("192.168.0.0/16"),
				IsActive:    utils.BoolPtr(false),
				Priority:    utils.Int64Ptr(5),
				NextHop: &v3.NextHop{
					LocalSubnetReference: &v3.Reference{
						Kind: utils.StringPtr("subnet"),
						UUID: utils.StringPtr("local-1"),
					},
				},
			}},
			want: []map[string]interface{}{{
				"destination": utils.StringPtr("192.168.0.0/16"),
				"is_active":   utils.BoolPtr(false),
				"priority":    utils.Int64Ptr(5),
				"nexthop": []interface{}{
					map[string]interface{}{
						"local_subnet_reference": map[string]interface{}{
							"kind": "subnet",
							"uuid": "local-1",
						},
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStatusStaticRouteList(tt.in))
		})
	}
}

func TestUnitFlattenNextHop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.NextHop
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.NextHop{
				ExternalSubnetReference: &v3.Reference{
					Kind: utils.StringPtr("subnet"),
					UUID: utils.StringPtr("ext-1"),
				},
				NexthopIPAddress: utils.StringPtr("10.0.0.1"),
			},
			want: []interface{}{
				map[string]interface{}{
					"external_subnet_reference": map[string]interface{}{
						"kind": "subnet",
						"uuid": "ext-1",
					},
					"nexthop_ip_address": utils.StringPtr("10.0.0.1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNextHop(tt.in))
		})
	}
}
