---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_gpu_v4"
sidebar_current: "docs-nutanix-datasource-vm-gpu-v4"
description: |-
   Retrieves configuration details for the provided GPU device attached to a Virtual Machine.
---

# nutanix_vm_gpu_v4
Retrieves configuration details for the provided GPU device attached to a Virtual Machine.


## Exanple

```hcl

    data "nutanix_vm_gpu_v4" "test" {
        vm_ext_id = {{ vm uuid }}
        ext_id = {{ gpu uuid }}
    }

```


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: The globally unique identifier of a VM. It should be of type UUID
* `ext_id`: The globally unique identifier of a VM GPU. It should be of type UUID.

## Attribute Reference

The following attributes are exported:

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

See detailed information in [Nutanix VM GPU](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).
