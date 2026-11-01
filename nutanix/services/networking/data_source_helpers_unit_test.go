package networking

import (
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFlattenVPCEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.VPCIntentResponse
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.VPCIntentResponse{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenVPCEntities(tt.in))
		})
	}

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		in := []*v3.VPCIntentResponse{{
			Metadata: &v3.Metadata{
				UUID: utils.StringPtr("meta-1"),
				Name: utils.StringPtr("vpc-meta"),
			},
			Status: &v3.VPCDefStatus{
				State: utils.StringPtr("COMPLETE"),
				Name:  utils.StringPtr("vpc-1"),
			},
			Spec: &v3.VPC{
				Name: utils.StringPtr("vpc-1"),
			},
		}}
		got := flattenVPCEntities(in)
		require.Len(t, got, 1)
		require.NotNil(t, got[0]["status"])
		require.NotNil(t, got[0]["spec"])
		require.NotNil(t, got[0]["metadata"])
	})
}

func TestUnitFlattenStatusVPC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.VPCDefStatus
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.VPCDefStatus{
				State: utils.StringPtr("COMPLETE"),
				Name:  utils.StringPtr("vpc-1"),
				ExecutionContext: &v3.ExecutionContext{
					TaskUUID: "task-1",
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"state":     utils.StringPtr("COMPLETE"),
					"name":      utils.StringPtr("vpc-1"),
					"resources": []interface{}{},
					"execution_context": []interface{}{
						map[string]interface{}{"task_uuid": "task-1"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStatusVPC(tt.in))
		})
	}
}

func TestUnitFlattenSpecVPC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.VPC
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in:   &v3.VPC{Name: utils.StringPtr("vpc-1")},
			want: []interface{}{
				map[string]interface{}{
					"name":      utils.StringPtr("vpc-1"),
					"resources": []interface{}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSpecVPC(tt.in))
		})
	}
}

func TestUnitFlattenResourcesVPC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.VpcResources
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.VpcResources{
				CommonDomainNameServerIPList: []*v3.CommonDomainNameServerIPList{
					{IP: utils.StringPtr("8.8.8.8")},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"common_domain_name_server_ip_list": []map[string]interface{}{
						{"ip": utils.StringPtr("8.8.8.8")},
					},
					"external_subnet_list":            []map[string]interface{}(nil),
					"externally_routable_prefix_list": []map[string]interface{}(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourcesVPC(tt.in))
		})
	}
}

func TestUnitFlattenExecutionContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.ExecutionContext
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in:   &v3.ExecutionContext{TaskUUID: "task-abc"},
			want: []interface{}{
				map[string]interface{}{"task_uuid": "task-abc"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExecutionContext(tt.in))
		})
	}
}

func TestUnitFlattenVPCResources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.VpcResources
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.VpcResources{
				ExternallyRoutablePrefixList: []*v3.ExternallyRoutablePrefixList{
					{IP: utils.StringPtr("10.1.0.0"), PrefixLength: utils.IntPtr(16)},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"common_domain_name_server_ip_list": []map[string]interface{}(nil),
					"external_subnet_list":              []map[string]interface{}(nil),
					"externally_routable_prefix_list": []map[string]interface{}{{
						"ip":            utils.StringPtr("10.1.0.0"),
						"prefix_length": utils.IntPtr(16),
					}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenVPCResources(tt.in))
		})
	}
}

func TestUnitFlattenVPCMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.ListMetadataOutput
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.ListMetadataOutput{
				TotalMatches: utils.Int64Ptr(2),
				Kind:         utils.StringPtr("vpc"),
				Length:       utils.Int64Ptr(2),
				Offset:       utils.Int64Ptr(0),
				Filter:       utils.StringPtr("name==vpc"),
			},
			want: []interface{}{
				map[string]interface{}{
					"total_matches": utils.Int64Ptr(2),
					"kind":          utils.StringPtr("vpc"),
					"length":        utils.Int64Ptr(2),
					"offset":        utils.Int64Ptr(0),
					"filter":        utils.StringPtr("name==vpc"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenVPCMetadata(tt.in))
		})
	}
}

func TestUnitFlattenStatusFIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.FloatingIPDefStatus
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.FloatingIPDefStatus{
				State: utils.StringPtr("COMPLETE"),
				Resource: &v3.FIPResource{
					FloatingIP: utils.StringPtr("1.2.3.4"),
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"state": utils.StringPtr("COMPLETE"),
					"resources": []interface{}{
						map[string]interface{}{
							"external_subnet_reference": map[string]interface{}{},
							"floating_ip":               utils.StringPtr("1.2.3.4"),
							"vm_nic_reference":          map[string]interface{}{},
							"vpc_reference":             map[string]interface{}{},
						},
					},
					"execution_context": []interface{}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStatusFIP(tt.in))
		})
	}
}

func TestUnitFlattenSpecFIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.FloatingIPSpec
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.FloatingIPSpec{
				Resource: &v3.FIPResource{
					ExternalSubnetReference: &v3.Reference{
						Kind: utils.StringPtr("subnet"),
						UUID: utils.StringPtr("sub-1"),
					},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"resources": []interface{}{
						map[string]interface{}{
							"external_subnet_reference": map[string]interface{}{
								"kind": "subnet",
								"uuid": "sub-1",
							},
							"vm_nic_reference": map[string]interface{}{},
							"vpc_reference":    map[string]interface{}{},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSpecFIP(tt.in))
		})
	}
}

func TestUnitFlattenFIPResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.FIPResource
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.FIPResource{
				FloatingIP: utils.StringPtr("9.9.9.9"),
				VPCReference: &v3.Reference{
					Kind: utils.StringPtr("vpc"),
					UUID: utils.StringPtr("vpc-1"),
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"external_subnet_reference": map[string]interface{}{},
					"floating_ip":               utils.StringPtr("9.9.9.9"),
					"vm_nic_reference":          map[string]interface{}{},
					"vpc_reference": map[string]interface{}{
						"kind": "vpc",
						"uuid": "vpc-1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFIPResource(tt.in))
		})
	}
}

func TestUnitFlattenFIPsEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.FloatingIPsIntentResponse
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.FloatingIPsIntentResponse{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFIPsEntities(tt.in))
		})
	}

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		in := []*v3.FloatingIPsIntentResponse{{
			Metadata: &v3.Metadata{UUID: utils.StringPtr("fip-1")},
			Status:   &v3.FloatingIPDefStatus{State: utils.StringPtr("COMPLETE")},
			Spec:     &v3.FloatingIPSpec{},
		}}
		got := flattenFIPsEntities(in)
		require.Len(t, got, 1)
		require.NotNil(t, got[0]["status"])
		require.NotNil(t, got[0]["spec"])
		require.NotNil(t, got[0]["metadata"])
	})
}

func TestUnitFlattenFIPsMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.ListMetadataOutput
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.ListMetadataOutput{
				TotalMatches: utils.Int64Ptr(1),
				Kind:         utils.StringPtr("floating_ip"),
			},
			want: []interface{}{
				map[string]interface{}{
					"total_matches": utils.Int64Ptr(1),
					"kind":          utils.StringPtr("floating_ip"),
					"length":        (*int64)(nil),
					"offset":        (*int64)(nil),
					"filter":        (*string)(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFIPsMetadata(tt.in))
		})
	}
}

func TestUnitFlattenServiceGroups(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenServiceGroups(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenServiceGroups([]*v3.ServiceGroupListEntry{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		in := []*v3.ServiceGroupListEntry{{
			UUID: utils.StringPtr("sg-1"),
			ServiceGroup: &v3.ServiceGroupInput{
				Name:          utils.StringPtr("svc"),
				Description:   utils.StringPtr("desc"),
				SystemDefined: utils.BoolPtr(false),
				ServiceList: []*v3.ServiceListEntry{{
					Protocol: utils.StringPtr("ALL"),
				}},
			},
			AssociatedPoliciesList: []*v3.ReferenceValues{{
				Kind: "policy", UUID: "p-1", Name: "pol",
			}},
		}}
		got := flattenServiceGroups(in).([]map[string]interface{})
		require.Len(t, got, 1)
		require.Equal(t, "sg-1", got[0]["uuid"])
	})
}

func TestUnitFlattenPbrStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.PbrDefStatus
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.PbrDefStatus{
				State: utils.StringPtr("COMPLETE"),
			},
			want: []interface{}{
				map[string]interface{}{
					"state":             utils.StringPtr("COMPLETE"),
					"execution_context": []interface{}{},
					"resources":         []interface{}(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPbrStatus(tt.in))
		})
	}
}

func TestUnitFlattenPbrSpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.PbrSpec
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &v3.PbrSpec{Name: utils.StringPtr("pbr-1")},
			want: []interface{}{
				map[string]interface{}{
					"name":      utils.StringPtr("pbr-1"),
					"resources": []interface{}(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPbrSpec(tt.in))
		})
	}
}

func TestUnitFlattenPbrResources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.PbrResources
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &v3.PbrResources{
				IsBidirectional: utils.BoolPtr(true),
				Priority:        utils.IntPtr(10),
				ProtocolType:    utils.StringPtr("TCP"),
				VpcReference: &v3.Reference{
					Kind: utils.StringPtr("vpc"),
					UUID: utils.StringPtr("vpc-1"),
				},
				Source: &v3.PbrSourDest{
					AddressType: utils.StringPtr("ANY"),
				},
				Action: &v3.PbrAction{
					Action:        utils.StringPtr("PERMIT"),
					ServiceIPList: []string{"1.1.1.1"},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"is_bidirectional": utils.BoolPtr(true),
					"priority":         utils.IntPtr(10),
					"protocol_type":    utils.StringPtr("TCP"),
					"vpc_reference": map[string]interface{}{
						"kind": "vpc",
						"uuid": "vpc-1",
					},
					"source": []interface{}{
						map[string]interface{}{
							"address_type": utils.StringPtr("ANY"),
						},
					},
					"destination":         []interface{}{},
					"protocol_parameters": []map[string]interface{}(nil),
					"action": []interface{}{
						map[string]interface{}{
							"action":          utils.StringPtr("PERMIT"),
							"service_ip_list": []*string{utils.StringPtr("1.1.1.1")},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPbrResources(tt.in))
		})
	}
}

func TestUnitFlattenPbrEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.PbrIntentResponse
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*v3.PbrIntentResponse{}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPbrEntities(tt.in))
		})
	}

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		in := []*v3.PbrIntentResponse{{
			Metadata: &v3.Metadata{UUID: utils.StringPtr("pbr-1")},
			Spec:     &v3.PbrSpec{Name: utils.StringPtr("pbr")},
			Status:   &v3.PbrDefStatus{State: utils.StringPtr("COMPLETE")},
		}}
		got := flattenPbrEntities(in)
		require.Len(t, got, 1)
		require.NotNil(t, got[0]["spec"])
		require.NotNil(t, got[0]["status"])
		require.NotNil(t, got[0]["metadata"])
	})
}

func TestUnitFlattenPbrMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.ListMetadataOutput
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &v3.ListMetadataOutput{
				TotalMatches: utils.Int64Ptr(3),
				Kind:         utils.StringPtr("routing_policy"),
			},
			want: []interface{}{
				map[string]interface{}{
					"total_matches": utils.Int64Ptr(3),
					"kind":          utils.StringPtr("routing_policy"),
					"length":        (*int64)(nil),
					"offset":        (*int64)(nil),
					"filter":        (*string)(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenPbrMetadata(tt.in))
		})
	}
}

func TestUnitFlattenAddressGroup(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenAddressGroup(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenAddressGroup([]*v3.AddressGroupListEntry{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		in := []*v3.AddressGroupListEntry{{
			UUID: utils.StringPtr("ag-1"),
			AddressGroup: &v3.AddressGroupInput{
				Name:               utils.StringPtr("ag"),
				Description:        utils.StringPtr("desc"),
				AddressGroupString: utils.StringPtr("10.0.0.0/24"),
				BlockList: []*v3.IPAddressBlock{{
					IPAddress:    utils.StringPtr("10.0.0.0"),
					PrefixLength: utils.Int64Ptr(24),
				}},
			},
			AssociatedPoliciesList: []*v3.ReferenceValues{{
				Kind: "policy", UUID: "p-1", Name: "pol",
			}},
		}}
		got := flattenAddressGroup(in).([]map[string]interface{})
		require.Len(t, got, 1)
		require.Equal(t, utils.StringPtr("ag-1"), got[0]["uuid"])
	})
}
