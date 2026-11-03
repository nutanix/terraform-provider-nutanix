package vmmv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const testDiskSerialID = "NFS_6_0_382_a36dcd9e_d389_4394_6557_be092fda879d"

func newTestVMDisk() *config.VmDisk {
	vmDisk := config.NewVmDisk()
	vmDisk.DiskExtId = utils.StringPtr("a36dcd9e-d389-4394-6557-be092fda879d")
	vmDisk.DiskSizeBytes = utils.Int64Ptr(10737418240)
	vmDisk.SerialId = utils.StringPtr(testDiskSerialID)
	return vmDisk
}

func newTestDisk(t *testing.T) config.Disk {
	t.Helper()
	backingInfo := config.NewOneOfDiskBackingInfo()
	if err := backingInfo.SetValue(*newTestVMDisk()); err != nil {
		t.Fatalf("failed to build OneOfDiskBackingInfo: %v", err)
	}
	disk := config.NewDisk()
	disk.BackingInfo = backingInfo
	return *disk
}

func TestFlattenVMDisk_SerialID(t *testing.T) {
	t.Parallel()

	got := flattenVMDisk(newTestVMDisk())

	if len(got) != 1 {
		t.Fatalf("expected exactly one element, got %d", len(got))
	}
	serialID, ok := got[0]["serial_id"].(*string)
	if !ok || serialID == nil || *serialID != testDiskSerialID {
		t.Fatalf("unexpected serial_id: %v", got[0]["serial_id"])
	}
}

func TestFlattenVMDisk_NoSerialID(t *testing.T) {
	t.Parallel()

	vmDisk := newTestVMDisk()
	vmDisk.SerialId = nil
	got := flattenVMDisk(vmDisk)

	if _, ok := got[0]["serial_id"]; ok {
		t.Fatalf("expected serial_id to be omitted when the API returns none, got %v", got[0]["serial_id"])
	}
}

func TestFlattenOneOfDiskBackingInfo_SerialID(t *testing.T) {
	t.Parallel()

	disk := newTestDisk(t)
	got := flattenOneOfDiskBackingInfo(disk.BackingInfo)

	if len(got) != 1 {
		t.Fatalf("expected exactly one element, got %d", len(got))
	}
	vmDisks, ok := got[0]["vm_disk"].([]map[string]interface{})
	if !ok || len(vmDisks) != 1 {
		t.Fatalf("expected one vm_disk, got %v", got[0]["vm_disk"])
	}
	serialID, ok := vmDisks[0]["serial_id"].(*string)
	if !ok || serialID == nil || *serialID != testDiskSerialID {
		t.Fatalf("unexpected serial_id: %v", vmDisks[0]["serial_id"])
	}
}

// Every vm_disk schema that receives flattenVMDisk or flattenOneOfDiskBackingInfo
// output must declare serial_id, otherwise d.Set fails for that resource.
func TestVMDiskSchemas_HaveSerialID(t *testing.T) {
	t.Parallel()

	resources := map[string]*schema.Resource{
		"nutanix_virtual_machine_v2":       ResourceNutanixVirtualMachineV2(),
		"data.nutanix_virtual_machine_v2":  DatasourceNutanixVirtualMachineV4(),
		"data.nutanix_virtual_machines_v2": DatasourceNutanixVirtualMachinesV4(),
		"nutanix_vm_clone_v2":              ResourceNutanixVMCloneV2(),
		"nutanix_template_v2":              ResourceNutanixTemplatesV2(),
		"nutanix_ova_vm_deploy_v2":         ResourceNutanixOvaVMDeploymentV2(),
		"nutanix_vm_cdrom_insert_eject_v2": ResourceNutanixVmsCdRomsInsertEjectV2(),
		"data.nutanix_template_v2":         DatasourceNutanixTemplateV2(),
		"data.nutanix_templates_v2":        DatasourceNutanixTemplatesV2(),
	}

	for name, r := range resources {
		found := 0
		walkVMDiskSchemas(r.Schema, func(path string, s map[string]*schema.Schema) {
			found++
			serialID, ok := s["serial_id"]
			if !ok {
				t.Errorf("%s: %s has no serial_id", name, path)
				return
			}
			if serialID.Type != schema.TypeString || !serialID.Computed || serialID.Optional || serialID.Required {
				t.Errorf("%s: %s.serial_id must be a computed-only string", name, path)
			}
		})
		if found == 0 {
			t.Errorf("%s: no vm_disk schema found", name)
		}
	}
}

// walkVMDiskSchemas calls fn for every nested schema shaped like a VmDisk,
// identified by its is_migration_in_progress attribute.
func walkVMDiskSchemas(s map[string]*schema.Schema, fn func(string, map[string]*schema.Schema)) {
	var walk func(string, map[string]*schema.Schema)
	walk = func(path string, s map[string]*schema.Schema) {
		if _, ok := s["is_migration_in_progress"]; ok {
			fn(path, s)
		}
		for k, v := range s {
			if r, ok := v.Elem.(*schema.Resource); ok {
				walk(path+"."+k, r.Schema)
			}
		}
	}
	walk("", s)
}

func TestSetDisksAndCdRoms_SerialID(t *testing.T) {
	t.Parallel()

	cdRom := config.NewCdRom()
	cdRom.BackingInfo = newTestVMDisk()
	disks := flattenDisk([]config.Disk{newTestDisk(t)})
	cdRoms := flattenCdRom([]config.CdRom{*cdRom})

	resources := map[string]*schema.Resource{
		"nutanix_virtual_machine_v2":      ResourceNutanixVirtualMachineV2(),
		"data.nutanix_virtual_machine_v2": DatasourceNutanixVirtualMachineV4(),
		"nutanix_vm_clone_v2":             ResourceNutanixVMCloneV2(),
	}

	for name, r := range resources {
		d := r.TestResourceData()
		if err := d.Set("disks", disks); err != nil {
			t.Fatalf("%s: setting disks: %v", name, err)
		}
		if err := d.Set("cd_roms", cdRoms); err != nil {
			t.Fatalf("%s: setting cd_roms: %v", name, err)
		}
		if got := d.Get("disks.0.backing_info.0.vm_disk.0.serial_id"); got != testDiskSerialID {
			t.Errorf("%s: disks serial_id = %q", name, got)
		}
		if got := d.Get("cd_roms.0.backing_info.0.serial_id"); got != testDiskSerialID {
			t.Errorf("%s: cd_roms serial_id = %q", name, got)
		}
	}

	d := ResourceNutanixVmsCdRomsInsertEjectV2().TestResourceData()
	if err := d.Set("backing_info", flattenVMDisk(newTestVMDisk())); err != nil {
		t.Fatalf("nutanix_vm_cdrom_insert_eject_v2: setting backing_info: %v", err)
	}
	if got := d.Get("backing_info.0.serial_id"); got != testDiskSerialID {
		t.Errorf("nutanix_vm_cdrom_insert_eject_v2: backing_info serial_id = %q", got)
	}
}
