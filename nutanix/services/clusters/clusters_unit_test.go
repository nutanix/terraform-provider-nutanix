package clusters

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitClusterSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, DataSourceNutanixCluster().InternalValidate(DataSourceNutanixCluster().Schema, false))
	require.NoError(t, DataSourceNutanixClusters().InternalValidate(DataSourceNutanixClusters().Schema, false))
	require.NoError(t, DataSourceNutanixHost().InternalValidate(DataSourceNutanixHost().Schema, false))
	require.NoError(t, DataSourceNutanixHosts().InternalValidate(DataSourceNutanixHosts().Schema, false))
	require.NoError(t, categoriesSchema().Elem.(*schema.Resource).InternalValidate(categoriesSchema().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, resourceNutanixDatasourceClusterResourceV0().InternalValidate(resourceNutanixDatasourceClusterResourceV0().Schema, false))
	require.NoError(t, resourceNutanixDatasourceClustersResourceV0().InternalValidate(resourceNutanixDatasourceClustersResourceV0().Schema, false))
}

func TestUnitFlattenCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]string
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{name: "empty", in: map[string]string{}, want: []interface{}{}},
		{
			name: "populated",
			in:   map[string]string{"env": "prod"},
			want: []interface{}{map[string]interface{}{"name": "env", "value": "prod"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenCategories(tt.in)
			if len(tt.in) > 1 {
				require.ElementsMatch(t, tt.want, got)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitFlattenReferenceValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.Reference
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{
			name: "empty",
			in:   &v3.Reference{},
			want: map[string]interface{}{"kind": "", "uuid": ""},
		},
		{
			name: "populated",
			in: &v3.Reference{
				Kind: utils.StringPtr("cluster"),
				UUID: utils.StringPtr("uuid-1"),
				Name: utils.StringPtr("pc"),
			},
			want: map[string]interface{}{
				"kind": "cluster",
				"uuid": "uuid-1",
				"name": "pc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReferenceValues(tt.in))
		})
	}
}

func TestUnitSetRSEntityMetadata(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	specVersion := int64(3)

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		meta, cats := setRSEntityMetadata(&v3.Metadata{})
		require.Equal(t, map[string]interface{}{
			"last_update_time": time.Time{}.String(),
			"uuid":             "",
			"creation_time":    time.Time{}.String(),
			"spec_version":     "0",
			"spec_hash":        "",
			"name":             "",
		}, meta)
		require.Equal(t, []interface{}{}, cats)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		meta, cats := setRSEntityMetadata(&v3.Metadata{
			LastUpdateTime: &when,
			CreationTime:   &when,
			UUID:           utils.StringPtr("uuid-1"),
			SpecVersion:    &specVersion,
			SpecHash:       utils.StringPtr("hash-1"),
			Name:           utils.StringPtr("host-1"),
			Categories:     map[string]string{"tier": "app"},
		})
		require.Equal(t, when.String(), meta["last_update_time"])
		require.Equal(t, when.String(), meta["creation_time"])
		require.Equal(t, "uuid-1", meta["uuid"])
		require.Equal(t, "3", meta["spec_version"])
		require.Equal(t, "hash-1", meta["spec_hash"])
		require.Equal(t, "host-1", meta["name"])
		require.Equal(t, []interface{}{map[string]interface{}{"name": "tier", "value": "app"}}, cats)
	})
}

func TestUnitFlattenFailOverCluster(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.FailoverCluster
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{name: "empty", in: &v3.FailoverCluster{}, want: map[string]interface{}{"ip": "", "name": ""}},
		{name: "populated", in: &v3.FailoverCluster{IP: "10.0.0.1", Name: "fo"}, want: map[string]interface{}{"ip": "10.0.0.1", "name": "fo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFailOverCluster(tt.in))
		})
	}
}

func TestUnitFlattenDomainCredencials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.DomainCredencial
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{name: "empty", in: &v3.DomainCredencial{}, want: map[string]interface{}{"username": "", "password": ""}},
		{name: "populated", in: &v3.DomainCredencial{Username: "admin", Password: "secret"}, want: map[string]interface{}{"username": "admin", "password": "secret"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDomainCredencials(tt.in))
		})
	}
}

func TestUnitFlattenIMPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.IPMI
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{name: "empty", in: &v3.IPMI{}, want: map[string]interface{}{"ip": ""}},
		{name: "populated", in: &v3.IPMI{IP: "10.0.0.2"}, want: map[string]interface{}{"ip": "10.0.0.2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIMPI(tt.in))
		})
	}
}

func TestUnitFlattenWindowsDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.WindowsDomain
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{
			name: "empty",
			in:   &v3.WindowsDomain{},
			want: map[string]interface{}{
				"name":                   "",
				"name_server_ip":         "",
				"organization_unit_path": "",
				"name_prefix":            "",
			},
		},
		{
			name: "populated",
			in: &v3.WindowsDomain{
				Name:                 "corp",
				NameServerIP:         "10.0.0.3",
				OrganizationUnitPath: "OU=Hosts",
				NamePrefix:           "hv-",
			},
			want: map[string]interface{}{
				"name":                   "corp",
				"name_server_ip":         "10.0.0.3",
				"organization_unit_path": "OU=Hosts",
				"name_prefix":            "hv-",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenWindowsDomain(tt.in))
		})
	}
}

func TestUnitFlattenReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.ReferenceValues
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{name: "empty", in: &v3.ReferenceValues{}, want: map[string]interface{}{"kind": "", "uuid": "", "name": ""}},
		{name: "populated", in: &v3.ReferenceValues{Kind: "disk", UUID: "u-1", Name: "d1"}, want: map[string]interface{}{"kind": "disk", "uuid": "u-1", "name": "d1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReference(tt.in))
		})
	}
}

func TestUnitFlattenReferenceList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.ReferenceValues
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*v3.ReferenceValues{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in:   []*v3.ReferenceValues{{Kind: "disk", UUID: "u-1", Name: "d1"}},
			want: []map[string]interface{}{{"kind": "disk", "uuid": "u-1", "name": "d1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReferenceList(tt.in))
		})
	}
}

func TestUnitFlattenExternalNetworkListReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.ReferenceValues
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{name: "empty", in: &v3.ReferenceValues{}, want: map[string]interface{}{"uuid": "", "name": ""}},
		{name: "populated", in: &v3.ReferenceValues{UUID: "net-1", Name: "vlan0"}, want: map[string]interface{}{"uuid": "net-1", "name": "vlan0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExternalNetworkListReference(tt.in))
		})
	}
}

func TestUnitFlattenExternalNetworkListReferenceList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.ReferenceValues
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*v3.ReferenceValues{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in:   []*v3.ReferenceValues{{UUID: "net-1", Name: "vlan0"}},
			want: []map[string]interface{}{{"uuid": "net-1", "name": "vlan0"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenExternalNetworkListReferenceList(tt.in))
		})
	}
}

func TestUnitFlattenHypervisor(t *testing.T) {
	t.Parallel()

	numVMs := int64(4)

	tests := []struct {
		name string
		in   *v3.Hypervisor
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{
			name: "empty",
			in:   &v3.Hypervisor{},
			want: map[string]interface{}{"num_vms": "", "ip": "", "hypervisor_full_name": ""},
		},
		{
			name: "populated",
			in:   &v3.Hypervisor{NumVms: &numVMs, IP: "10.0.0.4", HypervisorFullName: "AHV"},
			want: map[string]interface{}{"num_vms": "4", "ip": "10.0.0.4", "hypervisor_full_name": "AHV"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenHypervisor(tt.in))
		})
	}
}

func TestUnitFlattenControllerVM(t *testing.T) {
	t.Parallel()

	pct := 12.5
	size := int64(1024)
	port := int64(9440)

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, map[string]interface{}{}, flattenControllerVM(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenControllerVM(&v3.ControllerVM{
			IP:      "10.0.0.5",
			NatIP:   "203.0.113.5",
			NatPort: &port,
			OplogUsage: &v3.OplogUsage{
				OplogDiskPct:  &pct,
				OplogDiskSize: &size,
			},
		})
		require.Equal(t, map[string]interface{}{
			"ip":                           "10.0.0.5",
			"nat_ip":                       "203.0.113.5",
			"oplog_usage.opolog_disk_pct":  "12.5",
			"oplog_usage.opolog_disk_size": "1024",
			"nat_port":                     &port,
		}, got)
	})
}

func TestUnitFlattenBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *v3.Block
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{name: "empty", in: &v3.Block{}, want: map[string]interface{}{"block_serial_number": "", "block_model": ""}},
		{name: "populated", in: &v3.Block{BlockSerialNumber: "SN-1", BlockModel: "NX"}, want: map[string]interface{}{"block_serial_number": "SN-1", "block_model": "NX"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenBlock(tt.in))
		})
	}
}

func TestUnitFlattenGpuList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*v3.GPU
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*v3.GPU{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []*v3.GPU{{
				Status: "USED",
				Vendor: "NVIDIA",
				Name:   "A100",
				UUID:   "gpu-1",
				ConsumerReference: &v3.ReferenceValues{
					Kind: "vm",
					UUID: "vm-1",
					Name: "vm",
				},
			}},
			want: []map[string]interface{}{{
				"status":                    "USED",
				"vendor":                    "NVIDIA",
				"num_virtual_display_heads": (*int64)(nil),
				"assignable":                false,
				"license_list":              []*string(nil),
				"num_vgpus_allocated":       (*int64)(nil),
				"pci_address":               "",
				"name":                      "A100",
				"frame_buffer_size_mib":     (*int64)(nil),
				"index":                     (*int64)(nil),
				"uuid":                      "gpu-1",
				"numa_node":                 (*int64)(nil),
				"max_resolution":            "",
				"consumer_reference":        map[string]interface{}{"kind": "vm", "uuid": "vm-1", "name": "vm"},
				"mode":                      "",
				"fraction":                  (*int64)(nil),
				"guest_driver_version":      "",
				"device_id":                 (*int64)(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenGpuList(tt.in))
		})
	}
}

func TestUnitFlattenHostEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenHostEntities(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenHostEntities([]*v3.HostResponse{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenHostEntities([]*v3.HostResponse{{
			APIVersion: "3.1",
			Metadata: &v3.Metadata{
				UUID: utils.StringPtr("host-uuid"),
				Name: utils.StringPtr("host-1"),
			},
			Status: &v3.HostStatus{
				Name: "host-1",
				Resources: &v3.HostResources{
					GPUDriverVersion: "1.0",
					CPUModel:         "Intel",
					HostType:         "HYPER_CONVERGED",
					IPMI:             &v3.IPMI{IP: "10.0.0.9"},
					Block:            &v3.Block{BlockSerialNumber: "SN", BlockModel: "NX"},
				},
				ClusterReference: &v3.ReferenceValues{Kind: "cluster", UUID: "c-1", Name: "cluster"},
			},
		}})
		require.Len(t, got, 1)
		require.Equal(t, "host-1", got[0]["name"])
		require.Equal(t, "1.0", got[0]["gpu_driver_version"])
		require.Equal(t, "Intel", got[0]["cpu_model"])
		require.Equal(t, "3.1", got[0]["api_version"])
		require.Equal(t, map[string]interface{}{"ip": "10.0.0.9"}, got[0]["ipmi"])
		require.Equal(t, map[string]interface{}{"block_serial_number": "SN", "block_model": "NX"}, got[0]["block"])
		require.Equal(t, map[string]interface{}{"kind": "cluster", "uuid": "c-1", "name": "cluster"}, got[0]["cluster_reference"])
	})
}

func TestUnitResourceNutanixCategoriesMigrateState(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got, err := resourceNutanixCategoriesMigrateState(nil, nil)
		require.NoError(t, err)
		require.Nil(t, got)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got, err := resourceNutanixCategoriesMigrateState(map[string]interface{}{}, nil)
		require.NoError(t, err)
		require.Equal(t, map[string]interface{}{}, got)
	})

	t.Run("populated map categories", func(t *testing.T) {
		t.Parallel()
		got, err := resourceNutanixCategoriesMigrateState(map[string]interface{}{
			"categories": map[string]interface{}{
				"os_type": "ubuntu",
			},
		}, nil)
		require.NoError(t, err)
		require.Equal(t, []interface{}{
			map[string]interface{}{"name": "os_type", "value": "ubuntu"},
		}, got["categories"])
	})
}

func TestUnitResourceDatasourceStateUpgradeV0(t *testing.T) {
	t.Parallel()

	t.Run("cluster", func(t *testing.T) {
		t.Parallel()
		got, err := resourceDatasourceClusterStateUpgradeV0(context.Background(), map[string]interface{}{
			"categories": map[string]interface{}{"tier": "app"},
		}, nil)
		require.NoError(t, err)
		require.Equal(t, []interface{}{
			map[string]interface{}{"name": "tier", "value": "app"},
		}, got["categories"])
	})

	t.Run("clusters", func(t *testing.T) {
		t.Parallel()
		got, err := resourceDatasourceClustersStateUpgradeV0(context.Background(), map[string]interface{}{
			"categories": map[string]interface{}{"tier": "app"},
		}, nil)
		require.NoError(t, err)
		require.Equal(t, []interface{}{
			map[string]interface{}{"name": "tier", "value": "app"},
		}, got["categories"])
	})
}
