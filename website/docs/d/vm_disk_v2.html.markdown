---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_disk_v2"
sidebar_current: "docs-nutanix-datasource-vm-disk-v2"
description: |-
   This operation retrieves configuration details for the provided disk device attached to a Virtual Machine. 
---

# nutanix_vm_disk_v2

Retrieves configuration details for the provided disk device attached to a Virtual Machine.

## Example

```hcl

    data "nutanix_vm_disk_v4" "test"{
        vm_ext_id = {{ vm_ext_id }}
        ext_id = {{ disk uuid }}
    }
```


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: The globally unique identifier of a VM. It should be of type UUID
* `ext_id`: The globally unique identifier of a VM disk. It should be of type UUID.

## Attribute Reference

The following attributes are exported:
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `disk_address`: Disk address.
* `backing_info`: Supporting storage to create virtual disk on.

### disk_address
* `bus_type`: Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
* `index`: Device index on the bus. This field is ignored unless the bus details are specified.


### backing_info
* `vm_disk`: VM Disk info
* `adfs_volume_group_reference`:  Volume group reference


### backing_info.vm_disk
* `disk_ext_id`: The globally unique identifier of a VM disk. It should be of type UUID.
* `disk_size_bytes`: Size of the disk in Bytes
* `storage_container`: This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: A reference to a disk or image that contains the contents of a disk.
* `is_migration_in_progress`: Indicates if the disk is undergoing migration to another container.

### backing_info.vm_disk.data_source
* `reference`: Reference to image or vm disk
* `image_reference`: Image Reference
* `image_reference.image_ext_id`: The globally unique identifier of an image. It should be of type UUID.
* `vm_disk_reference`: Vm Disk Reference
* `vm_disk_reference.disk_ext_id`:  The globally unique identifier of a VM disk. It should be of type UUID.
* `vm_disk_reference.disk_address`: Disk address.
* `vm_disk_reference.vm_reference`: This is a reference to a VM.


See detailed information in [Nutanix VMs Disk](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).