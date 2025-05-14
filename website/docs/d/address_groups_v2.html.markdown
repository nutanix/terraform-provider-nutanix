---
layout: "nutanix"
page_title: "NUTANIX: nutanix_address_groups_v2"
sidebar_current: "docs-nutanix-datasource-address-groups-v2"
description: |-
  This operation retrieves the list of address_groups.
---

# nutanix_address_groups_v2

List all the Address Groups.

## Example Usage

```hcl
# list all address groups
data "nutanix_address_groups_v2" "list-addr-groups" {
}

# filtered the address groups
data "nutanix_address_groups_v2" "list-addr-group-filtered"{
  filter = "name eq 'td-addr-group'"
}

# filtered and limit the number of address groups
data "nutanix_address_groups_v2" "list-addr-groups-filter-limit" {
  filter = "name eq 'td-addr-group'"
  limit  = 1
}

```


## Argument Reference

The following arguments are supported:

* `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources. The filter can be applied to the following fields:
  - createdBy
  - description
  - extId
  - name

* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
  - description
  - extId
  - name
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields:
  - createdBy
  - description
  - extId
  - ipRanges
  - links
  - name
  - policyReferences
  - tenantId



## Attribute Reference
The following attributes are exported as a list:

* `address_groups`: List of address groups

### Address Group
The `address_groups` object contains the following attributes:

* `ext_id`: Address group UUID.
* `name`: A short identifier for an Address Group.
* `description`: A user defined annotation for an Address Group.
* `ipv4_addresses`: List of CIDR blocks in the Address Group.
* `ip_ranges`: List of IP range containing start and end IP
* `policy_references`: Reference to policy associated with Address Group.
* `created_by`: created by.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.


#### ipv4_addresses
* `value`: ip of address
* `prefix_length`: The prefix length of the network to which this host IPv4 address belongs.


#### ip_ranges
* `start_ip`: start ip
* `end_ip`: end ip




See detailed information in [Nutanix List Address Groups v4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.0#tag/AddressGroups/operation/listAddressGroups).
