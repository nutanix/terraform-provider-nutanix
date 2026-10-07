---
layout: "nutanix"
page_title: "NUTANIX: nutanix_template_placement_policy_v2"
sidebar_current: "docs-nutanix-datasource-template-placement-policy-v2"
description: |-
  Retrieves the details of the template placement policy for the provided external identifier.
---

# nutanix_template_placement_policy_v2

Retrieves the details of the template placement policy for the provided external identifier.

## Example

```hcl
data "nutanix_template_placement_policy_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) The external identifier of the template placement policy.

## Attribute Reference

The following attributes are exported:

* `name` - Name of the placement policy.
* `description` - Description of the placement policy.
* `placement_type` - The placement type. Valid values: `SOFT`.
* `cluster_filter` - Category filter for clusters.
* `content_filter` - Category filter for content.
* `create_time` - The time when the placement policy was created.
* `created_by` - External identifier of the user who created the placement policy.
* `update_time` - The time when the placement policy was last updated.
* `updated_by` - External identifier of the user who updated the placement policy.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

### cluster_filter, content_filter

* `category_ext_ids` - Filter to match entities based on the provided categories.
* `type` - The match type for categories. Valid values: `CATEGORIES_MATCH_ALL`, `CATEGORIES_MATCH_ANY`.

See detailed information in [Nutanix Get Template Placement Policy V2](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/TemplatePlacementPolicies/operation/getTemplatePlacementPolicyById).