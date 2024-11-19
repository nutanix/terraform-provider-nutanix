---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_gpus_v2"
sidebar_current: "docs-nutanix-datasource-vm-gpus-v4"
description: |-
   Lists the GPU devices attached to a Virtual Machine.
---

# nutanix_vm_gpus_v2
Lists the GPU devices attached to a Virtual Machine.


## Exanple

```hcl

    data "nutanix_vm_gpus_v2" "test" {
        vm_ext_id = {{ vm uuid }}
    }

    data "nutanix_vm_disks_v4" "test"{
        page = 0
        limit = 1
        vm_ext_id = {{ vm uuid }}
	}

```


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: The globally unique identifier of a VM. It should be of type UUID
* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions
* `gpus`: List of all GPUs. 


### gpus

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `mode`: The mode of this GPU.
* `device_id`: The device Id of the GPU.
* `vendor`: The vendor of the GPU.
* `pci_address`: The (S)egment:(B)us:(D)evice.(F)unction hardware address. 
* `guest_driver_version`: Last determined guest driver version.
* `name`: Name of the GPU resource.
* `frame_buffer_size_bytes`: GPU frame buffer size in bytes.
* `num_virtual_display_heads`: Number of supported virtual display heads.
* `fraction`: Fraction of the physical GPU assigned.

### pci_address
* `segment`
* `bus`
* `device`
* `func`

See detailed information in [Nutanix VM GPUs](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).
