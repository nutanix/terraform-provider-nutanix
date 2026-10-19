package vmm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitVmmSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixVirtualMachine().InternalValidate(ResourceNutanixVirtualMachine().Schema, true))
	require.NoError(t, ResourceNutanixImage().InternalValidate(ResourceNutanixImage().Schema, true))
	require.NoError(t, DataSourceNutanixVirtualMachine().InternalValidate(DataSourceNutanixVirtualMachine().Schema, false))
	require.NoError(t, DataSourceNutanixImage().InternalValidate(DataSourceNutanixImage().Schema, false))
	require.NoError(t, categoriesSchema().Elem.(*schema.Resource).InternalValidate(categoriesSchema().Elem.(*schema.Resource).Schema, true))
}

func TestUnitExpandStringList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []*string
	}{
		{name: "nil", in: nil, want: []*string{}},
		{name: "empty", in: []interface{}{}, want: []*string{}},
		{name: "skips empty strings", in: []interface{}{"", "a", 1}, want: []*string{utils.StringPtr("a")}},
		{name: "populated", in: []interface{}{"x", "y"}, want: []*string{utils.StringPtr("x"), utils.StringPtr("y")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandStringList(tt.in)
			require.Equal(t, len(tt.want), len(got))
			for i := range tt.want {
				require.Equal(t, *tt.want[i], *got[i])
			}
		})
	}
}

func TestUnitExpandCategories(t *testing.T) {
	t.Parallel()

	set := schema.NewSet(categoriesSchema().Set, []interface{}{
		map[string]interface{}{"name": "env", "value": "prod"},
	})
	got := expandCategories(set)
	require.Equal(t, map[string]string{"env": "prod"}, got)
}

func TestUnitFlattenCategories(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenCategories(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenCategories(map[string]string{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCategories(map[string]string{"env": "prod"})
		require.Equal(t, []interface{}{
			map[string]interface{}{"name": "env", "value": "prod"},
		}, got)
	})
}

func TestUnitFlattenReferenceValues(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, map[string]interface{}{}, flattenReferenceValues(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenReferenceValues(&v3.Reference{
			Kind: utils.StringPtr("vm"),
			UUID: utils.StringPtr("uuid-1"),
			Name: utils.StringPtr("vm-1"),
		})
		require.Equal(t, map[string]interface{}{
			"kind": "vm",
			"uuid": "uuid-1",
			"name": "vm-1",
		}, got)
	})
}

func TestUnitFlattenArrayOfReferenceValues(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenArrayOfReferenceValues(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenArrayOfReferenceValues([]*v3.ReferenceValues{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenArrayOfReferenceValues([]*v3.ReferenceValues{
			{Kind: "image", UUID: "img-1", Name: "img"},
			nil,
		})
		require.Equal(t, []map[string]interface{}{{
			"kind": "image",
			"uuid": "img-1",
			"name": "img",
		}}, got)
	})
}

func TestUnitFlattenNicListStatus(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenNicListStatus(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenNicListStatus([]*v3.VMNicOutputStatus{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicListStatus([]*v3.VMNicOutputStatus{{
			NicType:     utils.StringPtr("NORMAL_NIC"),
			UUID:        utils.StringPtr("nic-1"),
			MacAddress:  utils.StringPtr("00:11:22:33:44:55"),
			IsConnected: utils.BoolPtr(true),
			SubnetReference: &v3.Reference{
				UUID: utils.StringPtr("subnet-1"),
				Name: utils.StringPtr("subnet"),
			},
			IPEndpointList: []*v3.IPAddress{{
				IP:   utils.StringPtr("10.0.0.5"),
				Type: utils.StringPtr("ASSIGNED"),
			}},
		}})
		require.Len(t, got, 1)
		require.Equal(t, "NORMAL_NIC", got[0]["nic_type"])
		require.Equal(t, "true", got[0]["is_connected"])
		require.Equal(t, "subnet-1", got[0]["subnet_uuid"])
	})
}

func TestUnitFlattenNicList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenNicList(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenNicList([]*v3.VMNic{}))
	})

	t.Run("populated skips learned", func(t *testing.T) {
		t.Parallel()
		got := flattenNicList([]*v3.VMNic{{
			NicType: utils.StringPtr("NORMAL_NIC"),
			UUID:    utils.StringPtr("nic-1"),
			IPEndpointList: []*v3.IPAddress{
				{IP: utils.StringPtr("10.0.0.5"), Type: utils.StringPtr("ASSIGNED")},
				{IP: utils.StringPtr("10.0.0.6"), Type: utils.StringPtr("LEARNED")},
			},
		}})
		require.Len(t, got, 1)
		ips := got[0]["ip_endpoint_list"].([]map[string]interface{})
		require.Len(t, ips, 1)
		require.Equal(t, "10.0.0.5", ips[0]["ip"])
	})
}

func TestUnitFlattenDisk(t *testing.T) {
	t.Parallel()

	t.Run("minimal", func(t *testing.T) {
		t.Parallel()
		got := flattenDisk(&v3.VMDisk{
			UUID:          utils.StringPtr("disk-1"),
			DiskSizeBytes: utils.Int64Ptr(1024),
			DiskSizeMib:   utils.Int64Ptr(1),
		})
		require.Equal(t, "disk-1", got["uuid"])
		require.Equal(t, int64(1024), got["disk_size_bytes"])
		require.Equal(t, []map[string]interface{}(nil), got["device_properties"])
	})

	t.Run("with device properties", func(t *testing.T) {
		t.Parallel()
		got := flattenDisk(&v3.VMDisk{
			UUID: utils.StringPtr("disk-2"),
			DeviceProperties: &v3.VMDiskDeviceProperties{
				DeviceType: utils.StringPtr("DISK"),
				DiskAddress: &v3.DiskAddress{
					DeviceIndex: utils.Int64Ptr(0),
					AdapterType: utils.StringPtr("SCSI"),
				},
			},
		})
		props := got["device_properties"].([]map[string]interface{})
		require.Len(t, props, 1)
		addr := props[0]["disk_address"].(map[string]interface{})
		require.Equal(t, "0", addr["device_index"])
		require.Equal(t, utils.StringPtr("SCSI"), addr["adapter_type"])
	})
}

func TestUnitFlattenDiskList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenDiskList(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenDiskList([]*v3.VMDisk{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDiskList([]*v3.VMDisk{{UUID: utils.StringPtr("d1")}})
		require.Len(t, got, 1)
		require.Equal(t, "d1", got[0]["uuid"])
	})
}

func TestUnitFlattenSerialPortList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenSerialPortList(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenSerialPortList([]*v3.VMSerialPort{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSerialPortList([]*v3.VMSerialPort{{
			Index:       utils.Int64Ptr(0),
			IsConnected: utils.BoolPtr(true),
		}})
		require.Equal(t, []map[string]interface{}{{
			"index":        int64(0),
			"is_connected": true,
		}}, got)
	})
}

func TestUnitFlattenGPUList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenGPUList(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenGPUList([]*v3.VMGpuOutputStatus{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGPUList([]*v3.VMGpuOutputStatus{{
			Vendor:   utils.StringPtr("NVIDIA"),
			UUID:     utils.StringPtr("gpu-1"),
			DeviceID: utils.Int64Ptr(1),
			Mode:     utils.StringPtr("PASSTHROUGH_GRAPHICS"),
		}})
		require.Len(t, got, 1)
		require.Equal(t, "NVIDIA", got[0]["vendor"])
		require.Equal(t, int64(1), got[0]["device_id"])
	})
}

func TestUnitExpandIPAddressList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPAddressList(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPAddressList([]interface{}{}))
	})

	t.Run("populated defaults type", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddressList([]interface{}{
			map[string]interface{}{"ip": "10.0.0.8"},
			map[string]interface{}{"ip": "10.0.0.9", "type": "ASSIGNED"},
		})
		require.Len(t, got, 2)
		require.Equal(t, "10.0.0.8", *got[0].IP)
		require.Equal(t, "ASSIGNED", *got[0].Type)
		require.Equal(t, "ASSIGNED", *got[1].Type)
	})
}

func TestUnitExpandDiskListRaw(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandDiskListRaw([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandDiskListRaw([]interface{}{
			map[string]interface{}{
				"uuid":            "disk-1",
				"disk_size_mib":   1024,
				"disk_size_bytes": 0,
				"device_properties": []interface{}{
					map[string]interface{}{
						"device_type": "DISK",
						"disk_address": map[string]interface{}{
							"device_index": "0",
							"adapter_type": "SCSI",
						},
					},
				},
				"data_source_reference":  map[string]interface{}{},
				"volume_group_reference": map[string]interface{}{},
				"storage_config":         []interface{}{},
			},
		})
		require.Len(t, got, 1)
		require.Equal(t, "disk-1", *got[0].UUID)
		require.Equal(t, int64(1024), *got[0].DiskSizeMib)
		require.Equal(t, "DISK", *got[0].DeviceProperties.DeviceType)
	})
}

func TestUnitExpandStorageConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandStorageConfig(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandStorageConfig([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandStorageConfig([]interface{}{
			map[string]interface{}{
				"flash_mode": "ENABLED",
				"storage_container_reference": []interface{}{
					map[string]interface{}{
						"url":  "http://example",
						"kind": "storage_container",
						"uuid": "sc-1",
						"name": "container",
					},
				},
			},
		})
		require.Equal(t, "ENABLED", got.FlashMode)
		require.Equal(t, "sc-1", got.StorageContainerReference.UUID)
	})
}

func TestUnitExpandDeviceProperties(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandDeviceProperties(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandDeviceProperties([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandDeviceProperties([]interface{}{
			map[string]interface{}{
				"device_type": "CDROM",
				"disk_address": map[string]interface{}{
					"device_index": "1",
					"adapter_type": "IDE",
				},
			},
		})
		require.Equal(t, "CDROM", *got.DeviceType)
		require.Equal(t, int64(1), *got.DiskAddress.DeviceIndex)
		require.Equal(t, "IDE", *got.DiskAddress.AdapterType)
	})
}

func TestUnitExpandNicList(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.Nil(t, expandNicList(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{
			"nic_list": []interface{}{
				map[string]interface{}{
					"nic_type":     "NORMAL_NIC",
					"subnet_uuid":  "subnet-1",
					"is_connected": "true",
					"ip_endpoint_list": []interface{}{
						map[string]interface{}{"ip": "10.0.0.10", "type": "ASSIGNED"},
					},
				},
			},
		})
		got := expandNicList(d)
		require.Len(t, got, 1)
		require.Equal(t, "NORMAL_NIC", *got[0].NicType)
		require.Equal(t, "subnet-1", *got[0].SubnetReference.UUID)
		require.True(t, *got[0].IsConnected)
	})
}

func TestUnitExpandSerialPortList(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.Nil(t, expandSerialPortList(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{
			"serial_port_list": []interface{}{
				map[string]interface{}{"index": 0, "is_connected": true},
			},
		})
		got := expandSerialPortList(d)
		require.Len(t, got, 1)
		require.Equal(t, int64(0), *got[0].Index)
		require.True(t, *got[0].IsConnected)
	})
}

func TestUnitExpandGPUList(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.Nil(t, expandGPUList(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{
			"gpu_list": []interface{}{
				map[string]interface{}{
					"vendor":    "NVIDIA",
					"device_id": 1,
					"mode":      "PASSTHROUGH_GRAPHICS",
				},
			},
		})
		got := expandGPUList(d)
		require.Len(t, got, 1)
		require.Equal(t, "NVIDIA", *got[0].Vendor)
		require.Equal(t, int64(1), *got[0].DeviceID)
	})
}

func TestUnitExpandDiskList(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.Nil(t, expandDiskList(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{
			"disk_list": []interface{}{
				map[string]interface{}{
					"uuid":          "disk-1",
					"disk_size_mib": 512,
				},
			},
		})
		got := expandDiskList(d)
		require.Len(t, got, 1)
		require.Equal(t, "disk-1", *got[0].UUID)
	})
}

func TestUnitExpandNGT(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.Nil(t, expandNGT(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{
			"nutanix_guest_tools": map[string]interface{}{
				"state": "ENABLED",
			},
			"ngt_enabled_capability_list": []interface{}{"SELF_SERVICE_RESTORE"},
		})
		got := expandNGT(d)
		require.NotNil(t, got)
		require.Equal(t, "ENABLED", *got.NutanixGuestTools.State)
		require.Len(t, got.NutanixGuestTools.EnabledCapabilityList, 1)
	})
}

func TestUnitGetCdromDiskList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []*v3.VMDisk{}, GetCdromDiskList(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []*v3.VMDisk{}, GetCdromDiskList([]*v3.VMDisk{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		disks := []*v3.VMDisk{
			{UUID: utils.StringPtr("disk"), DeviceProperties: &v3.VMDiskDeviceProperties{DeviceType: utils.StringPtr("DISK")}},
			{UUID: utils.StringPtr("cdrom"), DeviceProperties: &v3.VMDiskDeviceProperties{DeviceType: utils.StringPtr("CDROM")}},
		}
		got := GetCdromDiskList(disks)
		require.Len(t, got, 1)
		require.Equal(t, "cdrom", *got[0].UUID)
	})
}

func TestUnitCountDiskListCdrom(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		n, err := CountDiskListCdrom(nil)
		require.NoError(t, err)
		require.Equal(t, 0, n)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		n, err := CountDiskListCdrom([]*v3.VMDisk{
			{DeviceProperties: &v3.VMDiskDeviceProperties{DeviceType: utils.StringPtr("CDROM")}},
			{DeviceProperties: &v3.VMDiskDeviceProperties{DeviceType: utils.StringPtr("DISK")}},
		})
		require.NoError(t, err)
		require.Equal(t, 1, n)
	})
}

func TestUnitValidateRefHelpers(t *testing.T) {
	t.Parallel()

	t.Run("validateRef empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, validateRef(map[string]interface{}{}))
	})

	t.Run("validateRef populated", func(t *testing.T) {
		t.Parallel()
		got := validateRef(map[string]interface{}{
			"kind": "vm",
			"uuid": "u-1",
			"name": "n-1",
		})
		require.Equal(t, "vm", *got.Kind)
		require.Equal(t, "u-1", *got.UUID)
		require.Equal(t, "n-1", *got.Name)
	})

	t.Run("validateShortRef", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, validateShortRef(map[string]interface{}{}))
		got := validateShortRef(map[string]interface{}{"kind": "image", "uuid": "i-1"})
		require.Equal(t, "image", *got.Kind)
		require.Equal(t, "i-1", *got.UUID)
	})

	t.Run("buildReference", func(t *testing.T) {
		t.Parallel()
		got := buildReference("uuid-1", "subnet")
		require.Equal(t, "uuid-1", *got.UUID)
		require.Equal(t, "subnet", *got.Kind)
	})
}

func TestUnitConvertMapInterfaceToMapString(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, map[string]string{}, convertMapInterfaceToMapString(map[string]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := convertMapInterfaceToMapString(map[string]interface{}{"a": 1, "b": "x"})
		require.Equal(t, map[string]string{"a": "1", "b": "x"}, got)
	})
}

func TestUnitFlattenClusterReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.NoError(t, flattenClusterReference(nil, d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.NoError(t, flattenClusterReference(&v3.Reference{
			UUID: utils.StringPtr("cluster-1"),
			Name: utils.StringPtr("cluster"),
		}, d))
		require.Equal(t, "cluster-1", d.Get("cluster_uuid"))
		require.Equal(t, "cluster", d.Get("cluster_name"))
	})
}

func TestUnitFlattenNutanixGuestTools(t *testing.T) {
	t.Parallel()

	t.Run("nil guest", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.NoError(t, flattenNutanixGuestTools(d, nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVirtualMachine().Schema, map[string]interface{}{})
		require.NoError(t, flattenNutanixGuestTools(d, &v3.GuestToolsStatus{
			NutanixGuestTools: &v3.NutanixGuestToolsStatus{
				State:                 utils.StringPtr("ENABLED"),
				Version:               utils.StringPtr("1.0"),
				EnabledCapabilityList: []*string{utils.StringPtr("VSS_SNAPSHOT")},
				IsReachable:           utils.BoolPtr(true),
			},
		}))
		ngt := d.Get("nutanix_guest_tools").(map[string]interface{})
		require.Equal(t, "ENABLED", ngt["state"])
		require.Equal(t, "true", ngt["is_reachable"])
	})
}
