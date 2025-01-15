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
resource "nutanix_vm_cdrom_insert_eject_v2" "insert-cdrom"{
  vm_ext_id = "<vm uuid>"
  ext_id    = "<vm cdrom uuid>"
  backing_info {
    data_source {
      reference {
        image_reference {
          image_ext_id = "<image_uuid>"
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


See detailed information in [Nutanix VMs CDROM Insert Eject V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0).
