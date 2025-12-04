---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_placement_v2"
sidebar_current: "docs-nutanix-datasource-image-placement-v2"
description: |-
 Describes a Image placement policy
---

# nutanix_image_v2

Retrieve the image placement policy details for the provided external identifier.

## Example

```hcl
data "nutanix_image_placement_policy_v2" "ipp"{
  ext_id = "cf96e27a-4e52-4cec-b563-d0b25413cc4a"
}
```


## Argument Reference

The following arguments are supported:

* `ext_id`: The external identifier of an image placement policy.

## Attribute Reference

The following arguments are supported:
* `name`: (Required) Name of the image placement policy.
* `description`: (Optional) Description of the image placement policy.
* `placement_type`: (Required) Type of the image placement policy. Valid values:
    - HARD: Hard placement policy. Images can only be placed on clusters enforced by the image placement policy.
    - SOFT: Soft placement policy. Images can be placed on clusters apart from those enforced by the image placement policy.
* `image_entity_filter`: (Required) Category-based entity filter.
* `cluster_entity_filter`: (Required) Category-based entity filter.
* `enforcement_state`: (Optional) Enforcement status of the image placement policy. Valid values:
    - ACTIVE: The image placement policy is being actively enforced.
    - SUSPENDED: The policy enforcement for image placement is suspended.

### image_entity_filter
* `type`: (Required) Filter matching type. Valid values:
    - CATEGORIES_MATCH_ALL: Image policy only applies to the entities that are matched to all the corresponding entity categories attached to the image policy.
    - CATEGORIES_MATCH_ANY: Image policy applies to the entities that match any subset of the entity categories attached to the image policy.
* `category_ext_ids`: Array of strings

### cluster_entity_filter
* `type`: (Required) Filter matching type. Valid values:
    - CATEGORIES_MATCH_ALL: Image policy only applies to the entities that are matched to all the corresponding entity categories attached to the image policy.
    - CATEGORIES_MATCH_ANY: Image policy applies to the entities that match any subset of the entity categories attached to the image policy.
* `category_ext_ids`: Array of strings

See detailed information in [Nutanix Get Image placement policy V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/ImagePlacementPolicies/operation/getPlacementPolicyById)
