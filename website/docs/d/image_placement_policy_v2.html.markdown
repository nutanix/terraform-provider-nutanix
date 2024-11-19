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
        ext_id = {{ ext_id of image placement policy }}
    }

```


## Argument Reference

The following arguments are supported:

* `ext_id`: The external identifier of an image placement policy.

## Attribute Reference

The following arguments are supported:
* `name`: (Required) Name of the image placement policy.
* `description`: (Optional) Description of the image placement policy.
* `placement_type`: (Required) Type of the image placement policy. Valid values "HARD", "SOFT"
* `image_entity_filter`: (Required) Category-based entity filter.
* `cluster_entity_filter`: (Required) Category-based entity filter.
* `enforcement_state`: (Optional) Enforcement status of the image placement policy. Valid values "ACTIVE", "SUSPENDED"

### image_entity_filter
* `type`: (Required) Filter matching type. Valid values "CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"
* `category_ext_ids`: Array of strings

### cluster_entity_filter
* `type`: (Required) Filter matching type. Valid values "CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"
* `category_ext_ids`: Array of strings

See detailed information in [Nutanix Image placement policy](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1)