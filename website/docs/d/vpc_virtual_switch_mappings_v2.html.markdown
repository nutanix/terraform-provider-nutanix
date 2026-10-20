---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vpc_virtual_switch_mappings_v2"
sidebar_current: "docs-nutanix-datasource-vpc-virtual-switch-mappings-v2"
description: |-
  List the VPC to Virtual Switch mappings and traffic configurations from the Nutanix v4 API.
---

# nutanix_vpc_virtual_switch_mappings_v2

The `nutanix_vpc_virtual_switch_mappings_v2` data source retrieves the current set of mappings between Virtual Private Clouds (VPCs) and Virtual Switches, along with their East-West traffic configurations.

Use this data source to inspect which Virtual Switches have traffic policies applied, the cluster boundaries they are effective in, and the organizational metadata (project ownership and categories) associated with each mapping.

## Example Usage

```hcl
data "nutanix_vpc_virtual_switch_mappings_v2" "example" {}
```

## Argument Reference

No arguments are required.

## Attribute Reference

The following attributes are exported:

* `vpc_virtual_switch_mappings` - The list of VPC to Virtual Switch mapping entries. See the [`vpc_virtual_switch_mappings`](#vpc_virtual_switch_mappings) block below for details.

### vpc_virtual_switch_mappings

Each entry in the `vpc_virtual_switch_mappings` list exports the following attributes:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `virtual_switch_uuid` - The Universally Unique Identifier (UUID) of the Virtual Switch the traffic policies are applied to.
* `cluster_uuids` - A list of cluster UUIDs where this Virtual Switch mapping is effective.
* `is_all_traffic_permitted` - A boolean flag indicating the traffic permission policy. When `true`, all network traffic is permitted through the virtual switch. When `false`, the switch only permits basic ICMP and essential statistics collection requests.
* `project_ext_id` - The UUID of the project that owns this entity.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A list of HATEOAS style links for the response. See the [`links`](#links) block below for details.
* `metadata` - Administrative and organizational metadata associated with the mapping. See the [`metadata`](#metadata) block below for details.

### links

The `links` block exports the following attributes:

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object that is returned by the URL.

### metadata

The `metadata` block exports the following attributes:

* `owner_reference_id` - A globally unique identifier (UUID) that represents the user, group, or service account that owns this mapping resource.
* `owner_user_name` - The human-readable username of the resource owner.
* `project_reference_id` - A globally unique identifier (UUID) of the Nutanix Project this resource belongs to.
* `project_name` - The human-readable name of the Nutanix Project this resource is assigned to.
* `category_ids` - A list of globally unique identifiers (UUIDs) representing the categories (tags) associated with the mapping.

See detailed information in [Nutanix List VPC Virtual Switch Mappings V2](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.4#tag/VpcVirtualSwitchMappings/operation/listVpcVirtualSwitchMappings).