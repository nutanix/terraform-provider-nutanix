---
layout: "nutanix"
page_title: "NUTANIX: nutanix_address_groups_v2"
sidebar_current: "docs-nutanix-resource-address-groups-v2"
description: |-
  This operation submits a request to create a address group based on the input parameters.
---

# nutanix_address_groups_v2

Create an Address Group

## Example Usage

```hcl
# Create Address group with ipv4 addresses
resource "nutanix_address_groups_v2" "ipv4-address" {
  name        = "address_group_ipv4_address"
  description = "address group description"
  ipv4_addresses {
    value         = "10.0.0.0"
    prefix_length = 24
  }
  ipv4_addresses {
    value         = "172.0.0.0"
    prefix_length = 24
  }
}

# Create Address group. with ip range
resource "nutanix_address_groups_v2" "ip-ranges" {
  name        = "address_group_ip_ranges"
  description = "address group description"
  ip_ranges {
    start_ip = "10.0.0.1"
    end_ip   = "10.0.0.10"
  }
}
```


## Argument Reference

The following arguments are supported:

* `name`: - (Required) Name of the Address group
* `description`: - (Optional) Description of the Address group
* `ipv4_addresses`: - (Optional) List of CIDR blocks in the Address Group.
* `ip_ranges`: - (Optional) List of IP range containing start and end IP.


### ipv4_addresses
* `value`: (Optional) ip of address
* `prefix_length`: (Optional) The prefix length of the network to which this host IPv4 address belongs.


### ip_ranges
* `start_ip`: (Required) start ip
* `end_ip`: (Required) end ip


## Attributes Reference

The following attributes are exported:

* `ext_id`: address group uuid.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `policy_references`: Reference to policy associated with Address Group.
* `created_by`: created by.


See detailed information in [Nutanix Address Group V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.0#tag/AddressGroups/operation/createAddressGroup).
