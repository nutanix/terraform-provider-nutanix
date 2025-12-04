---
layout: "nutanix"
page_title: "NUTANIX: nutanix_address_groups"
sidebar_current: "docs-nutanix-datasource-address-groups"
description: |-
  This operation retrieves list of address_groups.
---

# nutanix_address_group

Provides a datasource to retrieve list of address groups.

## Example Usage

```hcl
  data "nutanix_address_groups" "addr_groups" {}
```


## Attribute Reference

The following attributes are exported:

* `entities`:- (ReadOnly) List of address groups
* `metadata`:- (Optional) Use metadata to specify filters

### Metadata

The following arguments are supported:
* `filter`: (Optional) Filter in FIQL Syntax
* `sort_order`:  (Optional) order of sorting
* `offset`:  (Optional) Integer 
* `length`:  (Optional) Integer
* `sort_attribute`:  (Optional) attribute to sort

### Entities

The following attributes are exported as list:

* `address_group`: Information about address_group
* `associated_policies_list`: List of associated policies to address group

#### Address Group

The following attributes are exported:

* `name`:- (ReadOnly) Name of the address group
* `description`:- (ReadOnly) Description of the address group
* `ip_address_block_list`: - (ReadOnly) list of IP address blocks with their prefix length
* `address_group_string`: - (ReadOnly) Address Group string

##### IP Address Block List

The ip_address_block_list argument supports the following:

* `ip`: - (ReadOnly) IP of the address block
* `prefix_length`: - (ReadOnly) Prefix length of address block in int

#### Associated policies

The following attributes are exported as list:
* `name`: - (ReadOnly) Name of associated policy
* `uuid`: - (ReadOnly) UUID of associated policy


See detailed information in [Nutanix Address Group List](https://www.nutanix.dev/api_references/prism-central-v3/#/7504287ad168d-address-groups-lists).
