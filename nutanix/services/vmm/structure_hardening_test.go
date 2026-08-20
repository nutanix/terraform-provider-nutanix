package vmm

import (
	"testing"

	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// TestFlattenDisk_NilNestedPointers is the regression test for GH #1099
// ("Segfault on flattenDisk"). device_properties.disk_address and
// storage_config.storage_container_reference are optional in the v3 API — CD-ROM
// devices and volume-group-backed disks come back without them — and the flattener
// dereferenced both unconditionally.
func TestFlattenDisk_NilNestedPointers(t *testing.T) {
	for _, tc := range []struct {
		name string
		disk *v3.VMDisk
	}{
		{
			name: "device_properties present but disk_address nil",
			disk: &v3.VMDisk{
				UUID: utils.StringPtr("disk-1"),
				DeviceProperties: &v3.VMDiskDeviceProperties{
					DeviceType: utils.StringPtr("CDROM"),
				},
			},
		},
		{
			name: "storage_config present but container reference nil",
			disk: &v3.VMDisk{
				UUID:          utils.StringPtr("disk-2"),
				StorageConfig: &v3.VMStorageConfig{FlashMode: "ON"},
			},
		},
		{
			name: "both nested pointers nil",
			disk: &v3.VMDisk{
				UUID:             utils.StringPtr("disk-3"),
				DeviceProperties: &v3.VMDiskDeviceProperties{},
				StorageConfig:    &v3.VMStorageConfig{},
			},
		},
		{
			name: "bare disk",
			disk: &v3.VMDisk{UUID: utils.StringPtr("disk-4")},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := flattenDisk(tc.disk)
			if got == nil {
				t.Fatal("flattenDisk returned nil")
			}
			if got["uuid"] != utils.StringValue(tc.disk.UUID) {
				t.Errorf("uuid = %v, want %v", got["uuid"], utils.StringValue(tc.disk.UUID))
			}
		})
	}
}

// TestFlattenDisk_PopulatedDiskUnchanged pins the happy path so the nil-guards did not
// alter the shape produced for a fully-populated disk.
func TestFlattenDisk_PopulatedDiskUnchanged(t *testing.T) {
	disk := &v3.VMDisk{
		UUID:          utils.StringPtr("disk-1"),
		DiskSizeBytes: utils.Int64Ptr(1024),
		DeviceProperties: &v3.VMDiskDeviceProperties{
			DeviceType: utils.StringPtr("DISK"),
			DiskAddress: &v3.DiskAddress{
				DeviceIndex: utils.Int64Ptr(3),
				AdapterType: utils.StringPtr("SCSI"),
			},
		},
		StorageConfig: &v3.VMStorageConfig{
			FlashMode: "ON",
			StorageContainerReference: &v3.StorageContainerReference{
				UUID: "container-uuid",
				Name: "container-name",
			},
		},
	}

	got := flattenDisk(disk)

	deviceProps, ok := got["device_properties"].([]map[string]interface{})
	if !ok || len(deviceProps) != 1 {
		t.Fatalf("device_properties = %#v", got["device_properties"])
	}
	addr, ok := deviceProps[0]["disk_address"].(map[string]interface{})
	if !ok {
		t.Fatalf("disk_address = %#v", deviceProps[0]["disk_address"])
	}
	if addr["device_index"] != "3" {
		t.Errorf("device_index = %v, want \"3\"", addr["device_index"])
	}

	storageConfig, ok := got["storage_config"].([]map[string]interface{})
	if !ok || len(storageConfig) != 1 {
		t.Fatalf("storage_config = %#v", got["storage_config"])
	}
	refs, ok := storageConfig[0]["storage_container_reference"].([]map[string]interface{})
	if !ok || len(refs) != 1 {
		t.Fatalf("storage_container_reference = %#v", storageConfig[0]["storage_container_reference"])
	}
	if refs[0]["uuid"] != "container-uuid" {
		t.Errorf("container uuid = %v, want container-uuid", refs[0]["uuid"])
	}
}

// TestFlattenDiskList_MixedNilPointers walks the list path used by the VM read, since
// that is where the reported crash actually surfaced.
func TestFlattenDiskList_MixedNilPointers(t *testing.T) {
	disks := []*v3.VMDisk{
		{UUID: utils.StringPtr("ok"), DeviceProperties: &v3.VMDiskDeviceProperties{
			DiskAddress: &v3.DiskAddress{DeviceIndex: utils.Int64Ptr(0), AdapterType: utils.StringPtr("SCSI")},
		}},
		{UUID: utils.StringPtr("no-address"), DeviceProperties: &v3.VMDiskDeviceProperties{}},
		{UUID: nil},
	}

	got := flattenDiskList(disks)
	if len(got) != len(disks) {
		t.Errorf("flattenDiskList returned %d disks, want %d", len(got), len(disks))
	}
}
