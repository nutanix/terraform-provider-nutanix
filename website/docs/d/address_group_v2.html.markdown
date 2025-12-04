---
layout: "nutanix"
page_title: "NUTANIX: nutanix_address_group_v2"
sidebar_current: "docs-nutanix-datasource-address-group-v2"
description: |-
  This operation retrieves an address_group.
---

# nutanix_address_group_v2

Get an Address Group by ExtID

## Example Usage

```hcl
data "nutanix_address_group_v2" "get-addr-group"{
  ext_id = "0005b3b0-0b3b-4b3b-8b3b-0b3b3b3b3b3b"
}
```


## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) Address group UUID.

## Attribute Reference

The following attributes are exported:

* `name`: A short identifier for an Address Group.
* `description`: A user defined annotation for an Address Group.
* `ipv4_addresses`: List of CIDR blocks in the Address Group.
* `ip_ranges`: List of IP range containing start and end IP
* `policy_references`: Reference to policy associated with Address Group.
* `created_by`: created by.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.


### ipv4_addresses
* `value`: ip of address
* `prefix_length`: The prefix length of the network to which this host IPv4 address belongs.


### ip_ranges
* `start_ip`: start ip
* `end_ip`: end ip




See detailed information in [Nutanix Address Group v4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.0#tag/AddressGroups/operation/getAddressGroupById).
