---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_cdrom_insert_eject_v2"
sidebar_current: "docs-nutanix-resource-vm-cdroms-insert-eject-v2"
description: |-
   Inserts the defined ISO into a CD-ROM device attached to a Virtual Machine.
   Ejects the ISO currently inserted into a CD-ROM device on a Virtual Machine.
---

# nutanix_vm_cdrom_insert_eject_v2

Inserts the defined ISO into a CD-ROM device attached to a Virtual Machine.
Ejects the ISO currently inserted into a CD-ROM device on a Virtual Machine.


## Example

```hcl
##############################################
# ------------------------------------------------
# This resource allows inserting a custom ISO into
# a VM’s CD-ROM device.
#
# You can manage both:
# 1. **Insertion** — via `apply`
# 2. **Ejection** — automatically on `delete`
#  You can also eject the ISO by setting `action = "eject"` → triggers eject operation explicitly.
##############################################
resource "nutanix_vm_cdrom_insert_eject_v2" "insert-cdrom"{
  vm_ext_id = "8a938cc5-282b-48c4-81be-de22de145d07"
  ext_id    = "c2c249b0-98a0-43fa-9ff6-dcde578d3936"
  backing_info {
    data_source {
      reference {
        image_reference {
          image_ext_id = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `vm_ext_id`: (Required) The globally unique identifier of a VM. It should be of type UUID
* `ext_id`: (Required) The globally unique identifier of a CD-ROM. It should be of type UUID.
* `backing_info`: (Required) Storage provided by Nutanix ADSF
* `action`: (Optional) Default value: "insert". Accepted values: "insert" → Mounts the specified ISO image to the VM’s CD-ROM, "eject" → Unmounts (ejects) the ISO image from the VM’s CD-ROM.

### backing_info
* `disk_size_bytes`: (Required) Size of the disk in Bytes
* `storage_container`: (Optional) This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: (Optional) Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: (Required) Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: (Optional) A reference to a disk or image that contains the contents of a disk.


### backing_info.data_source
* `reference`: (Required) Reference to image or vm disk. Either `image_reference` or `vm_disk_reference`.
* `image_reference`: (Optional) Image Reference
* `image_reference.image_ext_id`: (Required) The globally unique identifier of an image. It should be of type UUID.

* `vm_disk_reference`: (Optional) Vm Disk Reference
* `vm_disk_reference.disk_address`: (Required) Disk address.
* `vm_disk_reference.vm_reference`: (Required) This is a reference to a VM.

## Import

This functionality allows you to eject CD-ROMs that were not mounted through Terraform.
To do so, you can import the inserted CD-ROMs (that need to be ejected) using their respective vm_ext_id and cdrom_ext_id (entity UUIDs).
```hcl
// Step 1: Create a placeholder resource in your root module. For example:
resource "nutanix_vm_cdrom_insert_eject_v2" "import_cdrom_inserted" {}

// Step 2: execute this command in cli
terraform import nutanix_vm_cdrom_insert_eject_v2.import_cdrom_inserted vm_ext_id/cdrom_ext_id

// Step 3: Once imported, update the resource configuration(resource placeholder added in Step 1) to perform the eject operation
resource "nutanix_vm_cdrom_insert_eject_v2" "import_cdrom_inserted" {
  vm_ext_id = <Virtual_Machine_UUID>
  ext_id    = <CD_ROM_UUID>
  action    = "eject"
}
```

See detailed information in [Nutanix VMs CDROM Insert V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/insertCdRomById).
