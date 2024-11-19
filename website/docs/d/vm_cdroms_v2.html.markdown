---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_cdroms_v4"
sidebar_current: "docs-nutanix-datasource-vm-cdroms-v4"
description: |-
   Lists the CD-ROM devices attached to a Virtual Machine.
---

# nutanix_vm_cdroms_v4

Lists the CD-ROM devices attached to a Virtual Machine.

## Example

```hcl

    data "nutanix_vm_cdroms_v4" "test"{
        vm_ext_id = {{ vm uuid }}
    }

    data "nutanix_vm_cdroms_v4" "test"{
        page = 0
        limit = 1
        vm_ext_id = {{ vm uuid }}
    }
```


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: The globally unique identifier of a VM. It should be of type UUID
* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results. Default is 0.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. Default is 50.
* `cdroms`: List of all cdroms attached to VM.

## Attribute Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `disk_address`: Disk address.
* `backing_info`: Supporting storage to create virtual disk on.
* `iso_type`: Type of ISO image inserted in CD-ROM. 

### disk_address
* `bus_type`: Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
* `index`: Device index on the bus. This field is ignored unless the bus details are specified.


### backing_info
* `disk_ext_id`: The globally unique identifier of a VM disk. It should be of type UUID.
* `disk_size_bytes`: Size of the disk in Bytes
* `storage_container`: This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: A reference to a disk or image that contains the contents of a disk.
* `is_migration_in_progress`: Indicates if the disk is undergoing migration to another container.


See detailed information in [Nutanix VMs CDROMs](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).