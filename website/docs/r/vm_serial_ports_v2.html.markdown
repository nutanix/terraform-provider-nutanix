---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_serial_ports_v4"
sidebar_current: "docs-nutanix-resource-vm-serial-ports-v2"
description: |-
  Creates and attaches a disk device to a Virtual Machine.
---

# nutanix_vm_serial_ports_v4

Creates and attaches a disk device to a Virtual Machine.


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: (Required) The globally unique identifier of a VM. It should be of type UUID.
* `index`: (Required) Index of the serial port.
* `is_connected`: (Required) Indicates whether the serial port is connected or not. 


## Attribute Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.


See detailed information in [Nutanix VMs Serial Port](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).