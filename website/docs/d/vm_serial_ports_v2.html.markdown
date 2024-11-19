---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_serial_ports_v4"
sidebar_current: "docs-nutanix-datasource-vm-serial-ports-v4"
description: |-
   Lists the serial ports attached to a Virtual Machine. 
---

# nutanix_vm_serial_ports_v4
Lists the serial ports attached to a Virtual Machine.


## Argument Reference

The following arguments are supported:

* `vm_ext_id`: The globally unique identifier of a VM. It should be of type UUID
* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results. Default is 0.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `serial_ports`: List of all serial ports attached to vms

## Attribute Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `is_connected`: Indicates whether the serial port is connected or not.
* `index`: Index of the serial port.

See detailed information in [Nutanix VMs Serial Ports](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).