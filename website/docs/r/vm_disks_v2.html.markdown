---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_disks_v4"
sidebar_current: "docs-nutanix-resource-vm-disks-v4"
description: |-
  Creates and attaches a disk device to a Virtual Machine.
---

# nutanix_vm_disks_v4

Creates and attaches a disk device to a Virtual Machine.

## Example

```hcl

    resource "nutanix_vm_disks_v4" "test" {
        vm_ext_id = {{ vm_ext_id }}
        disk_address{
            bus_type = "IDE"
            index= 0
        }
        backing_info{
            vm_disk{
            disk_size_bytes = 1073741824*2
            storage_container {
                ext_id = "{{ storage_ext_id }}"
            }
            }
        }
    }

    resource "nutanix_vm_disks_v4" "test" {
        vm_ext_id = {{ vm_ext_id }}
        disk_address{
            bus_type = "SCSI"
            index= 0
        }
        backing_info{
            vm_disk{
            disk_size_bytes = 1073741824
            storage_container {
                ext_id = "{{ storage_ext_id }}"
            }
            data_source{
                reference{
                    vm_disk_reference {
                        disk_address{
                        bus_type="SCSI"
                        index = 0
                        }
                        vm_reference{
                            ext_id= {{ vm disk reference ext_id }}
                        }
                    }
                }
            }
            }
        }
    }
```

## Argument Reference

The following arguments are supported:
* `vm_ext_id`: (Required) The globally unique identifier of a VM. It should be of type UUID.

* `disk_address`: (Required) Disk address.
* `backing_info`: (Required) Supporting storage to create virtual disk on. 


### disk_address
* `bus_type`: (Required) Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
* `index`: (Required) Device index on the bus. This field is ignored unless the bus details are specified.

### backing_info
* `vm_disk`: (Optional) VM Disk info
* `adfs_volume_group_reference`: (Optional) Volume group reference


### backing_info.vm_disk
* `disk_size_bytes`: (Required) Size of the disk in Bytes
* `storage_container`: (Optional) This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: (Optional) Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: (Optional) Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: (Optional) A reference to a disk or image that contains the contents of a disk.

### backing_info.vm_disk.data_source
* `reference`: Reference to image or vm disk. Either `image_reference` or `vm_disk_reference`.
* `image_reference`: (Optional) Image Reference
* `image_reference.image_ext_id`: (Required) The globally unique identifier of an image. It should be of type UUID.
* `vm_disk_reference`: (Optional) Vm Disk Reference
* `vm_disk_reference.disk_address`: (Required) Disk address.
* `vm_disk_reference.vm_reference`: (Required) This is a reference to a VM.


See detailed information in [Nutanix VMs Disk](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).
