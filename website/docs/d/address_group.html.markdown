---
layout: "nutanix"
page_title: "NUTANIX: nutanix_address_group"
sidebar_current: "docs-nutanix-datasource-address-group"
description: |-
  This operation retrieves an address_group.
---

# nutanix_address_group

Provides a datasource to retrieve a address group.

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

  data "nutanix_address_group" "addr_group" {
    uuid = nutanix_address_group.test_address.id
  }
```


## Attribute Reference

The following attributes are exported:

* `uuid`:- (Required) UUID of the address group
* `name`:- (ReadOnly) Name of the address group
* `description`:- (ReadOnly) Description of the address group
* `ip_address_block_list`: - (ReadOnly) list of IP address blocks with their prefix length
* `address_group_string`: - (ReadOnly) Address Group string


### IP Address Block List

The ip_address_block_list argument supports the following:

* `ip`: - (ReadOnly) IP of the address block
* `prefix_length`: - (ReadOnly) Prefix length of address block in int


See detailed information in [Nutanix Address Group](https://www.nutanix.dev/api_references/prism-central-v3/#/7921eaae69b35-get-a-existing-address-group).
