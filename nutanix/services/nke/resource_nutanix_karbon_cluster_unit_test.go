package nke

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/karbon"
)

func TestUnitResourceNutanixKarbonClusterSchema(t *testing.T) {
	t.Parallel()
	require.NoError(t, ResourceNutanixKarbonCluster().InternalValidate(ResourceNutanixKarbonCluster().Schema, true))
}

func TestUnitExpandPrivateRegistry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      map[string]interface{}
		want    string
		wantErr bool
	}{
		{name: "nil", in: nil, wantErr: true},
		{name: "empty", in: map[string]interface{}{}, wantErr: true},
		{name: "populated", in: map[string]interface{}{"registry_name": "reg-1"}, want: "reg-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := expandPrivateRegistry(tt.in)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, *got.RegistryName)
		})
	}
}

func TestUnitExpandPrivateRegistries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []interface{}
		wantLen int
		wantErr bool
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []interface{}{}, wantLen: 0},
		{
			name:    "populated",
			in:      []interface{}{map[string]interface{}{"registry_name": "reg-a"}},
			wantLen: 1,
		},
		{
			name:    "missing name",
			in:      []interface{}{map[string]interface{}{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := expandPrivateRegistries(tt.in)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Len(t, *got, tt.wantLen)
		})
	}
}

func TestUnitFlattenCalicoConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *karbon.ClusterCalicoConfig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: &karbon.ClusterCalicoConfig{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: &karbon.ClusterCalicoConfig{
				IPPoolConfigs: []karbon.ClusterCalicoConfigIPPoolConfig{{CIDR: "172.20.0.0/16"}},
			},
			want: []map[string]interface{}{{
				"ip_pool_config": []map[string]interface{}{{"cidr": "172.20.0.0/16"}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCalicoConfig(tt.in))
		})
	}
}

func TestUnitFlattenCNIConfig(t *testing.T) {
	t.Parallel()

	flannel := &karbon.ClusterFlannelConfig{}
	calico := &karbon.ClusterCalicoConfig{
		IPPoolConfigs: []karbon.ClusterCalicoConfigIPPoolConfig{{CIDR: "10.0.0.0/16"}},
	}

	tests := []struct {
		name string
		in   karbon.ClusterCNIConfig
		want []map[string]interface{}
	}{
		{
			name: "empty",
			in:   karbon.ClusterCNIConfig{},
			want: []map[string]interface{}{{
				"flannel_config":      (*karbon.ClusterFlannelConfig)(nil),
				"node_cidr_mask_size": int64(0),
				"pod_ipv4_cidr":       "",
				"service_ipv4_cidr":   "",
			}},
		},
		{
			name: "populated flannel",
			in: karbon.ClusterCNIConfig{
				NodeCIDRMaskSize: 24,
				PodIPv4CIDR:      "172.20.0.0/16",
				ServiceIPv4CIDR:  "172.19.0.0/16",
				FlannelConfig:    flannel,
			},
			want: []map[string]interface{}{{
				"flannel_config":      flannel,
				"node_cidr_mask_size": int64(24),
				"pod_ipv4_cidr":       "172.20.0.0/16",
				"service_ipv4_cidr":   "172.19.0.0/16",
			}},
		},
		{
			name: "populated calico",
			in: karbon.ClusterCNIConfig{
				NodeCIDRMaskSize: 24,
				PodIPv4CIDR:      "10.0.0.0/16",
				ServiceIPv4CIDR:  "10.1.0.0/16",
				CalicoConfig:     calico,
			},
			want: []map[string]interface{}{{
				"flannel_config":      (*karbon.ClusterFlannelConfig)(nil),
				"node_cidr_mask_size": int64(24),
				"pod_ipv4_cidr":       "10.0.0.0/16",
				"service_ipv4_cidr":   "10.1.0.0/16",
				"calico_config": []map[string]interface{}{{
					"ip_pool_config": []map[string]interface{}{{"cidr": "10.0.0.0/16"}},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenCNIConfig(tt.in))
		})
	}
}

func TestUnitFlattenNodePool(t *testing.T) {
	t.Parallel()

	hostname := "node-1"
	ip := "10.0.0.10"
	name := "worker-node-pool"
	osVersion := "ntnx-1.0"
	num := int64(1)
	nodes := []karbon.ClusterNodeIntentResponse{{Hostname: &hostname, IPv4Address: &ip}}

	nodepool := &karbon.ClusterNodePool{
		Name:          &name,
		NodeOSVersion: &osVersion,
		NumInstances:  &num,
		Nodes:         &nodes,
		AHVConfig: &karbon.ClusterNodePoolAHVConfig{
			CPU:                     4,
			DiskMib:                 10240,
			MemoryMib:               8192,
			NetworkUUID:             "net-uuid",
			PrismElementClusterUUID: "pe-uuid",
		},
	}

	userDefined := &karbon.ClusterNodePool{
		AHVConfig: &karbon.ClusterNodePoolAHVConfig{
			DiskMib:     20480,
			NetworkUUID: "user-net-uuid",
		},
	}

	t.Run("without user defined", func(t *testing.T) {
		t.Parallel()
		got := flattenNodePool(nil, nodepool)
		require.Equal(t, &name, got["name"])
		require.Equal(t, &num, got["num_instances"])
		require.Equal(t, &osVersion, got["node_os_version"])
		require.Equal(t, []map[string]interface{}{{
			"hostname":     &hostname,
			"ipv4_address": &ip,
		}}, got["nodes"])
		require.Equal(t, []map[string]interface{}{{
			"cpu":                        int64(4),
			"disk_mib":                   int64(10240),
			"memory_mib":                 int64(8192),
			"network_uuid":               "net-uuid",
			"prism_element_cluster_uuid": "pe-uuid",
		}}, got["ahv_config"])
	})

	t.Run("with user defined disk network", func(t *testing.T) {
		t.Parallel()
		got := flattenNodePool(userDefined, nodepool)
		ahv := got["ahv_config"].([]map[string]interface{})
		require.Equal(t, int64(20480), ahv[0]["disk_mib"])
		require.Equal(t, "user-net-uuid", ahv[0]["network_uuid"])
	})
}

func TestUnitExpandStorageClassConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []interface{}
		wantErr bool
		check   func(t *testing.T, got *karbon.ClusterStorageClassConfigIntentInput)
	}{
		{name: "nil", in: nil, wantErr: true},
		{name: "empty", in: []interface{}{}, wantErr: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"name":           "default-storageclass",
				"reclaim_policy": "Delete",
				"volumes_config": []interface{}{map[string]interface{}{
					"file_system":                "ext4",
					"flash_mode":                 true,
					"password":                   "secret",
					"prism_element_cluster_uuid": "pe-1",
					"storage_container":          "ctr-1",
					"username":                   "admin",
				}},
			}},
			check: func(t *testing.T, got *karbon.ClusterStorageClassConfigIntentInput) {
				t.Helper()
				require.Equal(t, "default-storageclass", got.Name)
				require.Equal(t, "Delete", got.ReclaimPolicy)
				require.True(t, got.DefaultStorageClass)
				require.Equal(t, "ext4", got.VolumesConfig.FileSystem)
				require.True(t, got.VolumesConfig.FlashMode)
				require.Equal(t, "secret", got.VolumesConfig.Password)
				require.Equal(t, "pe-1", got.VolumesConfig.PrismElementClusterUUID)
				require.Equal(t, "ctr-1", got.VolumesConfig.StorageContainer)
				require.Equal(t, "admin", got.VolumesConfig.Username)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := expandStorageClassConfig(tt.in)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			tt.check(t, got)
		})
	}
}

func TestUnitExpandCNI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []interface{}
		wantErr bool
		check   func(t *testing.T, got *karbon.ClusterCNIConfigIntentInput)
	}{
		{name: "nil", in: nil, wantErr: true},
		{name: "empty", in: []interface{}{}, wantErr: true},
		{
			name: "flannel",
			in: []interface{}{map[string]interface{}{
				"node_cidr_mask_size": 24,
				"pod_ipv4_cidr":       "172.20.0.0/16",
				"service_ipv4_cidr":   "172.19.0.0/16",
			}},
			check: func(t *testing.T, got *karbon.ClusterCNIConfigIntentInput) {
				t.Helper()
				require.Equal(t, int64(24), got.NodeCIDRMaskSize)
				require.Equal(t, "172.20.0.0/16", got.PodIPv4CIDR)
				require.Equal(t, "172.19.0.0/16", got.ServiceIPv4CIDR)
				require.NotNil(t, got.FlannelConfig)
				require.Nil(t, got.CalicoConfig)
			},
		},
		{
			name: "calico",
			in: []interface{}{map[string]interface{}{
				"node_cidr_mask_size": 24,
				"pod_ipv4_cidr":       "10.0.0.0/16",
				"service_ipv4_cidr":   "10.1.0.0/16",
				"calico_config": []interface{}{map[string]interface{}{
					"ip_pool_config": []interface{}{map[string]interface{}{"cidr": "10.0.0.0/16"}},
				}},
			}},
			check: func(t *testing.T, got *karbon.ClusterCNIConfigIntentInput) {
				t.Helper()
				require.Nil(t, got.FlannelConfig)
				require.NotNil(t, got.CalicoConfig)
				require.Equal(t, "10.0.0.0/16", got.CalicoConfig.IPPoolConfigs[0].CIDR)
			},
		},
		{
			name: "both calico and flannel",
			in: []interface{}{map[string]interface{}{
				"calico_config":  []interface{}{map[string]interface{}{}},
				"flannel_config": []interface{}{map[string]interface{}{}},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := expandCNI(tt.in)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			tt.check(t, got)
		})
	}
}

func TestUnitExpandNodePool(t *testing.T) {
	t.Parallel()

	karbonVersion := &karbon.MetaSemanticVersionResponse{
		MajorVersion:    2,
		MinorVersion:    2,
		RevisionVersion: 3,
	}

	tests := []struct {
		name    string
		in      []interface{}
		wantErr bool
		check   func(t *testing.T, got []karbon.ClusterNodePool)
	}{
		{name: "nil", in: nil, check: func(t *testing.T, got []karbon.ClusterNodePool) {
			t.Helper()
			require.Empty(t, got)
		}},
		{name: "empty", in: []interface{}{}, check: func(t *testing.T, got []karbon.ClusterNodePool) {
			t.Helper()
			require.Empty(t, got)
		}},
		{
			name: "missing name",
			in: []interface{}{map[string]interface{}{
				"name": "",
			}},
			wantErr: true,
		},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"name":            "worker-node-pool",
				"node_os_version": "ntnx-1.0",
				"num_instances":   2,
				"ahv_config": []interface{}{map[string]interface{}{
					"cpu":                        8,
					"disk_mib":                   122880,
					"memory_mib":                 8192,
					"network_uuid":               "net-1",
					"prism_element_cluster_uuid": "pe-1",
				}},
				"nodes": []interface{}{map[string]interface{}{
					"hostname":     "w1",
					"ipv4_address": "10.0.0.20",
				}},
			}},
			check: func(t *testing.T, got []karbon.ClusterNodePool) {
				t.Helper()
				require.Len(t, got, 1)
				require.Equal(t, "worker-node-pool", *got[0].Name)
				require.Equal(t, "ntnx-1.0", *got[0].NodeOSVersion)
				require.Equal(t, int64(2), *got[0].NumInstances)
				require.Equal(t, int64(8), got[0].AHVConfig.CPU)
				require.Equal(t, int64(122880), got[0].AHVConfig.DiskMib)
				require.Equal(t, "net-1", got[0].AHVConfig.NetworkUUID)
				require.Equal(t, "w1", *(*got[0].Nodes)[0].Hostname)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := expandNodePool(tt.in, karbonVersion)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.check(t, got)
		})
	}
}

func TestUnitGetSupportedFileSystems(t *testing.T) {
	t.Parallel()
	require.Equal(t, []string{"ext4", "xfs"}, getSupportedFileSystems())
}

func TestUnitGetSupportedReclaimPolicies(t *testing.T) {
	t.Parallel()
	require.Equal(t, []string{"Delete", "Retain"}, getSupportedReclaimPolicies())
}

func TestUnitCalculateCPURequirement(t *testing.T) {
	t.Parallel()

	oldVersion := &karbon.MetaSemanticVersionResponse{MajorVersion: 2, MinorVersion: 2, RevisionVersion: 1}
	newVersion := &karbon.MetaSemanticVersionResponse{MajorVersion: 2, MinorVersion: 2, RevisionVersion: 2}

	tests := []struct {
		name    string
		ver     *karbon.MetaSemanticVersionResponse
		cpu     int64
		want    int64
		wantErr bool
	}{
		{name: "old even", ver: oldVersion, cpu: 4, want: 2},
		{name: "old odd", ver: oldVersion, cpu: 3, wantErr: true},
		{name: "new no divide", ver: newVersion, cpu: 8, want: 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := calculateCPURequirement(tt.ver, tt.cpu)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitDiffFlatPrivateRegistrySlices(t *testing.T) {
	t.Parallel()

	a := "reg-a"
	b := "reg-b"
	c := "reg-c"

	tests := []struct {
		name string
		s1   []karbon.PrivateRegistryOperationIntentInput
		s2   []karbon.PrivateRegistryOperationIntentInput
		want []string
	}{
		{name: "nil", s1: nil, s2: nil, want: []string{}},
		{name: "empty", s1: []karbon.PrivateRegistryOperationIntentInput{}, s2: []karbon.PrivateRegistryOperationIntentInput{}, want: []string{}},
		{
			name: "populated",
			s1: []karbon.PrivateRegistryOperationIntentInput{
				{RegistryName: &a},
				{RegistryName: &b},
			},
			s2: []karbon.PrivateRegistryOperationIntentInput{
				{RegistryName: &b},
				{RegistryName: &c},
			},
			want: []string{"reg-a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := diffFlatPrivateRegistrySlices(tt.s1, tt.s2)
			names := make([]string, 0, len(got))
			for _, g := range got {
				names = append(names, *g.RegistryName)
			}
			require.Equal(t, tt.want, names)
		})
	}
}

func TestUnitConvertKarbonPrivateRegistryIntentInputToOperation(t *testing.T) {
	t.Parallel()

	name := "reg-1"
	got := convertKarbonPrivateRegistryIntentInputToOperation(karbon.PrivateRegistryResponse{Name: &name})
	require.Equal(t, &name, got.RegistryName)
}

func TestUnitConvertKarbonPrivateRegistriesIntentInputToOperations(t *testing.T) {
	t.Parallel()

	a := "a"
	b := "b"

	tests := []struct {
		name string
		in   karbon.PrivateRegistryListResponse
		want []string
	}{
		{name: "nil", in: nil, want: []string{}},
		{name: "empty", in: karbon.PrivateRegistryListResponse{}, want: []string{}},
		{
			name: "populated",
			in: karbon.PrivateRegistryListResponse{
				{Name: &a},
				{Name: &b},
			},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := convertKarbonPrivateRegistriesIntentInputToOperations(tt.in)
			names := make([]string, 0, len(got))
			for _, g := range got {
				names = append(names, *g.RegistryName)
			}
			require.Equal(t, tt.want, names)
		})
	}
}

func TestUnitAddActivePassiveConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      interface{}
		wantErr bool
		wantIP  string
	}{
		{name: "empty list", in: []interface{}{}, wantErr: true},
		{
			name:    "missing address",
			in:      []interface{}{map[string]interface{}{}},
			wantErr: true,
		},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"external_ipv4_address": "10.0.0.100",
			}},
			wantIP: "10.0.0.100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cluster := &karbon.ClusterIntentInput{}
			err := addActivePassiveConfig(tt.in, cluster)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantIP, cluster.MastersConfig.ActivePassiveConfig.ExternalIPv4Address)
		})
	}
}

func TestUnitAddExternalLBConfig(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixKarbonCluster().Schema, map[string]interface{}{
		"name":    "kc-1",
		"version": "1.25.0",
		"external_lb_config": []interface{}{map[string]interface{}{
			"external_ipv4_address": "10.0.0.50",
			"master_nodes_config": []interface{}{
				map[string]interface{}{"ipv4_address": "10.0.0.51", "node_pool_name": "master-node-pool"},
				map[string]interface{}{"ipv4_address": "10.0.0.52", "node_pool_name": "master-node-pool"},
			},
		}},
	})

	cluster := &karbon.ClusterIntentInput{}
	require.NoError(t, addExternalLBConfig(d.Get("external_lb_config"), cluster))
	require.Equal(t, "10.0.0.50", cluster.MastersConfig.ExternalLBConfig.ExternalIPv4Address)
	require.Len(t, cluster.MastersConfig.ExternalLBConfig.MasterNodesConfig, 2)

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		require.Error(t, addExternalLBConfig([]interface{}{}, &karbon.ClusterIntentInput{}))
	})
}

func TestUnitGetTimeout(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixKarbonCluster().Schema, map[string]interface{}{
			"name":                 "kc-1",
			"version":              "1.25.0",
			"wait_timeout_minutes": 30,
		})
		got, err := getTimeout(d)
		require.NoError(t, err)
		require.Equal(t, int64(30), got)
	})
}

func TestUnitCheckKarbonVersion(t *testing.T) {
	t.Parallel()

	respVersion := "1.28.0"
	resp := &karbon.ClusterIntentResponse{Version: &respVersion}

	t.Run("old version uses state", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixKarbonCluster().Schema, map[string]interface{}{
			"name":    "kc-1",
			"version": "1.25.0",
		})
		old := &karbon.MetaSemanticVersionResponse{MajorVersion: 2, MinorVersion: 2, RevisionVersion: 2}
		require.Equal(t, "1.25.0", checkKarbonVersion(d, old, resp))
	})

	t.Run("new version uses response", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixKarbonCluster().Schema, map[string]interface{}{
			"name":    "kc-1",
			"version": "1.25.0",
		})
		newVer := &karbon.MetaSemanticVersionResponse{MajorVersion: 2, MinorVersion: 2, RevisionVersion: 3}
		require.Equal(t, "1.28.0", checkKarbonVersion(d, newVer, resp))
	})

	t.Run("old version falls back to response", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"wait_timeout_minutes": ResourceNutanixKarbonCluster().Schema["wait_timeout_minutes"],
			"version":              ResourceNutanixKarbonCluster().Schema["version"],
		}, map[string]interface{}{})
		old := &karbon.MetaSemanticVersionResponse{MajorVersion: 2, MinorVersion: 2, RevisionVersion: 2}
		require.Equal(t, "1.28.0", checkKarbonVersion(d, old, resp))
	})
}
