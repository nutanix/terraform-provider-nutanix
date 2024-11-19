---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_cdroms_v4"
sidebar_current: "docs-nutanix-resource-vm-cdroms-v4"
description: |-
   Creates and attaches a CD-ROM device to a Virtual Machine.
---

# nutanix_vm_cdroms_v4

Creates and attaches a CD-ROM device to a Virtual Machine.

## Example

```hcl

   resource "nutanix_vm_cdroms_v4" "test" {
        vm_ext_id = {{ vm uuid }}
        disk_address{
            bus_type = "SATA"
            index= 1
        }
    }

    resource "nutanix_vm_cdroms_v4" "test" {
        vm_ext_id = {{ vm uuid }}
        disk_address{
            bus_type = "IDE"
            index= 1
        }
        backing_info{
            disk_size_bytes = 1073741824
            storage_container {
            ext_id = "{{ container uuid }}"
            }
        }
    }

    resource "nutanix_vm_cdroms_v4" "test" {
			vm_ext_id = {{ vm uuid }}
			disk_address{
			  bus_type = "IDE"
			  index= 1
			}
			backing_info{
				disk_size_bytes = 21474836480
				data_source {
					reference{
						image_reference{
							image_ext_id = "{{ image uuid}} "
						}
					}
				}
			}
		}

```


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: (Required) The globally unique identifier of a VM. It should be of type UUID
* `disk_address`: (Required) Disk address.
* `backing_info`: (Optional) Supporting storage to create virtual disk on.
* `iso_type`: (Optional) Type of ISO image inserted in CD-ROM. 


### disk_address
* `bus_type`: (Required) Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
* `index`: (Required) Device index on the bus. This field is ignored unless the bus details are specified.


### backing_info
* `disk_size_bytes`: (Required) Size of the disk in Bytes
* `storage_container`: (Optional) This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: (Optional) Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: (Required) Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: (Optional) A reference to a disk or image that contains the contents of a disk.


## Attribute Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.


See detailed information in [Nutanix VMs CDROM](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).