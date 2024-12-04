---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_gpus_v2"
sidebar_current: "docs-nutanix-resource-vm-gpus-v2"
description: |-
  Attaches a GPU device to a Virtual Machine.
---

# nutanix_vm_gpus_v2

Attaches a GPU device to a Virtual Machine.

## Example 

```hcl
    resource "nutanix_vm_gpus_v2" "test" {
        vm_ext_id = {{ vm uuid }}
        mode = "PASSTHROUGH_COMPUTE"
        vendor= "NVIDIA"
    }
```


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: (Required) The globally unique identifier of a VM. It should be of type UUID.
* `mode`: (Required) The mode of this GPU. The acceptable values are "PASSTHROUGH_GRAPHICS", "PASSTHROUGH_COMPUTE", "VIRTUAL" . 
* `device_id`: (Optional) The device Id of the GPU.
* `vendor`: (Optional) The vendor of the GPU. The acceptables values are "NVIDIA", "AMD", "INTEL" .
* `pci_address`: (Optional) The (S)egment:(B)us:(D)evice.(F)unction hardware address. 


### pci_address
* `segment`
* `bus`
* `device`
* `func`


## Attribute Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `guest_driver_version`: Last determined guest driver version.
* `name`: Name of the GPU resource.
* `frame_buffer_size_bytes`: GPU frame buffer size in bytes.
* `num_virtual_display_heads`: Number of supported virtual display heads.
* `fraction`: Fraction of the physical GPU assigned.


See detailed information in [Nutanix VM GPU](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).

