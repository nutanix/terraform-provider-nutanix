---
layout: "nutanix"
page_title: "NUTANIX: nutanix_entity_group_v2"
sidebar_current: "docs-nutanix-resource-entity-group-v2"
description: |-
  Creates and manages an Entity Group for microsegmentation.
---

# nutanix_entity_group_v2

Create and manage an Entity Group for microsegmentation.

## Example Usage

```hcl
# Simple entity group with name and description
resource "nutanix_entity_group_v2" "simple" {
  name        = "my-entity-group"
  description = "Entity group for microsegmentation"
}

# Entity group with allowed_config (VM by category + address group by IP)
resource "nutanix_entity_group_v2" "with_allowed" {
  name        = "entity_group_with_allowed"
  description = "Entity group with allowed entities"

  allowed_config {
    # VM entities selected by category (valid pair: CATEGORY_EXT_ID, VM)
    entities {
      type             = "VM"
      selected_by      = "CATEGORY_EXT_ID"
      reference_ext_ids = ["category-uuid-1", "category-uuid-2"]
    }
    # Address group by IP (valid pair: IP_VALUES, ADDRESS_GROUP). Only one entity per (selected_by, type)—combine addresses and ip_ranges in one block.
    entities {
      type        = "ADDRESS_GROUP"
      selected_by = "IP_VALUES"
      addresses {
        ipv4_addresses {
          value         = "10.0.0.0"
          prefix_length = 24
        }
      }
      ip_ranges {
        ipv4_ranges {
          start_ip = "192.168.1.1"
          end_ip   = "192.168.1.10"
        }
      }
    }
  }
}
```

## Argument Reference

* `name` - (Required) A short identifier of an Entity Group.
* `description` - (Optional) A user defined annotation for an Entity Group.
* `allowed_config` - (Optional) Configuration of the allowed entities in the Entity Group.
* `except_config` - (Optional) Configuration of except entities in the Entity Group.
* `policy_ext_ids` - (Optional) List of policy external identifiers.

### allowed_config

* `entities` - (Optional) List of allowed entities. Each entity may contain:
  * `type` - (Required) The type of entity. Valid values: `KUBE_NAMESPACE`, `SUBNET`, `VM`, `VPC`, `KUBE_SERVICE`, `KUBE_CLUSTER`, `KUBE_PODS`, `ADDRESS_GROUP`.
  * `selected_by` - (Optional) The selection method for the entity. Valid values: `IP_VALUES`, `EXT_ID`, `CATEGORY_EXT_ID`, `LABELS`, `NAME`.
  * `addresses` - (Optional) With `ipv4_addresses` block(s):
    * `value` - (Required) IPv4 address value.
    * `prefix_length` - (Optional) Prefix length.
  * `ip_ranges` - (Optional) With `ipv4_ranges` block(s):
    * `start_ip` - (Required) Start IP of the range.
    * `end_ip` - (Required) End IP of the range.
  * `kube_entities` - (Optional) List of kube entity identifiers. Required when `type` is a kube type (`KUBE_NAMESPACE`, `KUBE_SERVICE`, `KUBE_CLUSTER`, or `KUBE_PODS`).
  * `reference_ext_ids` - (Optional) List of reference external identifiers. Required when `selected_by` is `EXT_ID`.

### except_config

* `entities` - (Optional) List of except entities. Each entity may contain:
  * `addresses` - (Optional) With `ipv4_addresses` block(s).
  * `ip_ranges` - (Optional) With `ipv4_ranges` block(s).
  * `reference_ext_ids` - (Optional) List of reference external identifiers. Required when `selected_by` is `EXT_ID`.

## Validation Requirements

The following validation rules apply to `allowed_config` entities:

### Required Fields

* `type` - (Required) Must be specified for all entities in `allowed_config`.

### Conditional Requirements

* `kube_entities` - Required when `type` is one of: `KUBE_NAMESPACE`, `KUBE_SERVICE`, `KUBE_CLUSTER`, or `KUBE_PODS`. Must not be empty.
* `reference_ext_ids` - Required when `selected_by` is `EXT_ID`. Must not be empty.

### Valid Combinations

The combination of `selected_by` and `type` must be one of the following valid pairs:

* `(CATEGORY_EXT_ID, VM)`
* `(CATEGORY_EXT_ID, SUBNET)`
* `(CATEGORY_EXT_ID, VPC)`
* `(EXT_ID, KUBE_CLUSTER)`
* `(EXT_ID, ADDRESS_GROUP)`
* `(LABELS, KUBE_PODS)`
* `(NAME, KUBE_NAMESPACE)`
* `(NAME, KUBE_SERVICE)`
* `(IP_VALUES, ADDRESS_GROUP)`

Any other combination will result in a validation error.

### Duplicate (selected_by, type) Not Allowed

Within one entity group, you cannot have two entities with the same `(selected_by, type)` pair. For example, two entities both using `IP_VALUES` and `ADDRESS_GROUP` are invalid. Combine all addresses and ip_ranges into a single entity block when using `(IP_VALUES, ADDRESS_GROUP)`.

## Attributes Reference

* `ext_id` - Entity group UUID.
* `creation_time` - The timestamp when the Entity Group was created.
* `last_update_time` - The timestamp when the Entity Group was last updated.
* `links` - A HATEOAS style link for the response.
* `owner_ext_id` - The external identifier of the user who created the Entity Group.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

## Import

Entity Group can be imported using the entity group uuid `entityGroupUUID` (ext_id in v4 terms). eg,

```bash
// create its configuration in the root module. For example:
resource "nutanix_entity_group_v2" "import_entity_group"{}

// execute the below command.
terraform import nutanix_entity_group_v2.import_entity_group <entityGroupUUID>
```

See detailed information in [Nutanix Entity Groups V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.2#tag/EntityGroups/operation/createEntityGroup).
