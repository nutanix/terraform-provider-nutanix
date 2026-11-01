package selfservice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitResourceNutanixCalmAppProvisionSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, ResourceNutanixCalmAppProvision().InternalValidate(ResourceNutanixCalmAppProvision().Schema, true))
}

func TestUnitFlattenImageConfigName(t *testing.T) {
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
				"create_spec": map[string]interface{}{
					"resources": map[string]interface{}{
						"disk_list": []interface{}{
							map[string]interface{}{
								"data_source_reference": map[string]interface{}{
									"name": "centos-image",
								},
							},
						},
					},
				},
			},
			want: "centos-image",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenImageConfigName(tt.in))
		})
	}
}

func TestUnitFlattenNicMacAddress(t *testing.T) {
	t.Parallel()

	platformJSON := `[{"status":{"resources":{"nic_list":[{"mac_address":"aa:bb:cc:dd:ee:ff","nic_type":"NORMAL_NIC","subnet_reference":{"name":"vlan0"}}]}}}]`

	tests := []struct {
		name string
		in   map[string]interface{}
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: map[string]interface{}{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: map[string]interface{}{
				"variable_list": []interface{}{
					map[string]interface{}{
						"name":  "platform_data",
						"value": platformJSON,
					},
				},
			},
			want: []map[string]interface{}{{
				"mac_address": "aa:bb:cc:dd:ee:ff",
				"type":        "NORMAL_NIC",
				"subnet":      "vlan0",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNicMacAddress(tt.in))
		})
	}
}

func TestUnitFlattenVM(t *testing.T) {
	t.Parallel()

	platformJSON := `[{"status":{"resources":{"nic_list":[{"mac_address":"11:22:33:44:55:66","nic_type":"NORMAL_NIC","subnet_reference":{"name":"net1"}}]}}}]`

	tests := []struct {
		name string
		in   map[string]interface{}
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: map[string]interface{}{}, want: nil},
		{
			name: "populated",
			in: map[string]interface{}{
				"resources": map[string]interface{}{
					"deployment_list": []interface{}{
						map[string]interface{}{
							"substrate_configuration": map[string]interface{}{
								"create_spec": map[string]interface{}{
									"resources": map[string]interface{}{
										"disk_list": []interface{}{
											map[string]interface{}{
												"data_source_reference": map[string]interface{}{"name": "img-1"},
											},
										},
									},
								},
								"variable_list": []interface{}{
									map[string]interface{}{
										"name":  "platform_data",
										"value": platformJSON,
									},
								},
								"element_list": []interface{}{
									map[string]interface{}{
										"instance_name": "vm-1",
										"address":       "10.0.0.5",
										"instance_id":   "vm-uuid-1",
										"create_spec": map[string]interface{}{
											"resources": map[string]interface{}{
												"num_sockets":          2,
												"num_vcpus_per_socket": 4,
												"memory_size_mib":      8192,
											},
											"categories": map[string]interface{}{"AppType": "Demo"},
											"cluster_reference": map[string]interface{}{
												"name": "cluster-a",
												"uuid": "cluster-uuid",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"nics": []map[string]interface{}{{
						"mac_address": "11:22:33:44:55:66",
						"type":        "NORMAL_NIC",
						"subnet":      "net1",
					}},
					"categories": map[string]interface{}{"AppType": "Demo"},
					"cluster_info": []map[string]interface{}{{
						"cluster_name": "cluster-a",
						"cluster_uuid": "cluster-uuid",
					}},
					"configuration": []map[string]interface{}{{
						"name":       "vm-1",
						"ip_address": "10.0.0.5",
						"vm_uuid":    "vm-uuid-1",
						"image":      "img-1",
						"vcpus":      2,
						"cores":      4,
						"memory":     8192,
					}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenVM(tt.in))
		})
	}
}
