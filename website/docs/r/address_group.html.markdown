---
layout: "nutanix"
page_title: "NUTANIX: nutanix_address_group"
sidebar_current: "docs-nutanix-resource-address-group"
description: |-
  This operation submits a request to create a address group based on the input parameters.
---

# nutanix_address_group

Provides a resource to create a address group based on the input parameters.

## Example Usage

```hcl
resource "nutanix_address_group" "test_address" {
	name = "test"
	description = "test address groups resource"

	ip_address_block_list {
		ip = "10.0.0.0"
		prefix_length = 24
	}
}
```


## Argument Reference

The following arguments are supported:

* `name`: - (Required) Name of the service group
* `description`: - (Optional) Description of the service group
* `ip_address_block_list`: - (Required) list of IP address blocks with their prefix length
* `address_group_string`: - (ReadOnly) Address Group string

### IP Address List

The ip_address_block_list argument supports the following:

* `ip`: - (Required) IP of the address block
* `prefix_length`: - (Required) Prefix length of address block in int

See detailed information in [Nutanix Address Groups](https://www.nutanix.dev/api_references/prism-central-v3/#/5ccef53a546a4-create-a-new-address-group).
