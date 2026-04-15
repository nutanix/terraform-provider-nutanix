---
layout: "nutanix"
page_title: "NUTANIX: nutanix_entity_group_v2"
sidebar_current: "docs-nutanix-datasource-entity-group-v2"
description: |-
  Fetches the entity group identified by an external identifier (ext_id).
---

# nutanix_entity_group_v2

Fetches a single entity group by its external identifier (ext_id). Use this data source when you know the entity group UUID and need its attributes.

## Example Usage

```hcl
# Fetch by known ext_id
data "nutanix_entity_group_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}

# Fetch entity group created by a resource
data "nutanix_entity_group_v2" "by_id" {
  ext_id = nutanix_entity_group_v2.my_group.id
}
```

## Argument Reference

* `ext_id` - (Required) The external identifier (UUID) of the entity group.

## Attributes Reference

The following attributes are exported:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `name` - A short identifier of an Entity Group.
* `description` - A user defined annotation for an Entity Group.
* `allowed_config` - Configuration of the allowed entities in the Entity Group.
* `except_config` - Configuration of except entities in the Entity Group.
* `policy_ext_ids` - Mapping of entity group to the list of policy external identifiers.
* `creation_time` - The timestamp when the Entity Group was created.
* `last_update_time` - The timestamp when the Entity Group was last updated.
* `links` - A HATEOAS style link for the response.
* `owner_ext_id` - The external identifier of the user who created the Entity Group.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

### allowed_config

* `entities` - List of allowed entities. Each entity may contain:
  * `type` - The type of entity (e.g. VM, ADDRESS_GROUP).
  * `selected_by` - The selection method (e.g. CATEGORY_EXT_ID, IP_VALUES).
  * `addresses` - With `ipv4_addresses` (value, prefix_length).
  * `ip_ranges` - With `ipv4_ranges` (start_ip, end_ip).
  * `kube_entities` - List of kube entities.
  * `reference_ext_ids` - List of reference external identifiers.

### except_config

* `entities` - List of except entities. Each entity may contain:
  * `addresses` - With `ipv4_addresses` (value, prefix_length).
  * `ip_ranges` - With `ipv4_ranges` (start_ip, end_ip).
  * `reference_ext_ids` - List of reference external identifiers.

See detailed information in [Nutanix Get Entity Group V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.2#tag/EntityGroups/operation/getEntityGroupById).
