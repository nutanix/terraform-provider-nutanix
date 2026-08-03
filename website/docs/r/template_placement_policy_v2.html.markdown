---
layout: "nutanix"
page_title: "NUTANIX: nutanix_template_placement_policy_v2"
sidebar_current: "docs-nutanix-resource-template-placement-policy-v2"
description: |-
  Creates a template placement policy based on the provided request body.
---

# nutanix_template_placement_policy_v2

Creates a template placement policy based on the provided request body. Supports create, read, update, and delete operations.

## Example

```hcl
resource "nutanix_template_placement_policy_v2" "example" {
  name           = "example-policy"
  description    = "Example template placement policy"
  placement_type = "SOFT"
  content_filter {
    type             = "CATEGORIES_MATCH_ANY"
    category_ext_ids = ["00000000-0000-0000-0000-000000000001"]
  }

  cluster_filter {
    type             = "CATEGORIES_MATCH_ALL"
    category_ext_ids = ["00000000-0000-0000-0000-000000000000"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the placement policy.
* `description` - (Optional) Description of the placement policy.
* `placement_type` - (Optional) The placement type. Valid values: `SOFT`.
* `cluster_filter` - (Optional) Category filter for clusters. See below for nested schema.
* `content_filter` - (Optional) Category filter for content. See below for nested schema.

### cluster_filter, content_filter

* `category_ext_ids` - (Required) Filter to match entities based on the provided categories.
* `type` - (Required) The match type for categories. Valid values: `CATEGORIES_MATCH_ALL`, `CATEGORIES_MATCH_ANY`.

## Attribute Reference

The following attributes are exported:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `create_time` - The time when the placement policy was created.
* `created_by` - External identifier of the user who created the placement policy.
* `update_time` - The time when the placement policy was last updated.
* `updated_by` - External identifier of the user who updated the placement policy.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

## Import

Template placement policies can be imported using the `ext_id`:

```shell
terraform import nutanix_template_placement_policy_v2.example <ext_id>
```
