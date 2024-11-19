---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_serial_port_v4"
sidebar_current: "docs-nutanix-datasource-vm-serial-port-v4"
description: |-
   Retrieves configuration details for the provided serial port attached to a Virtual Machine.
---

# nutanix_vm_serial_port_v4
Retrieves configuration details for the provided serial port attached to a Virtual Machine.


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: The globally unique identifier of a VM. It should be of type UUID
* `ext_id`: The globally unique identifier of a serial port. It should be of type UUID.

## Attribute Reference

The following attributes are exported:

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `is_connected`: Indicates whether the serial port is connected or not.
* `index`: Index of the serial port.

See detailed information in [Nutanix VMs Serial Port](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).