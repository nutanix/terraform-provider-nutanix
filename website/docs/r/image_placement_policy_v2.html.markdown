---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_placement_v2"
sidebar_current: "docs-nutanix-resource-image-placement-v2"
description: |-
  Provides a Nutanix Image resource to Create a Image.
---

# nutanix_image_placement_v2

Create an image placement policy using the provided request body. Name, placement_type, image_entity_filter and source are mandatory fields to create an policy.


```hcl

data "nutanix_categories_v2" "categories"{}

locals {
	category0 = data.nutanix_categories_v4.categories.categories.0.ext_id
}
resource "nutanix_image_placement_policy_v2" "example"{
	name           = "image_placement_policy"
	description    = "%[2]s"
	placement_type = "SOFT"
	cluster_entity_filter {
		category_ext_ids = [
			local.category0,
		]
		type = "CATEGORIES_MATCH_ALL"
	}
	image_entity_filter {
		category_ext_ids = [
			local.category0,
		]
		type = "CATEGORIES_MATCH_ALL"
	}
}
```

## Argument Reference

The following arguments are supported:
* `name`: (Required) Name of the image placement policy.
* `description`: (Optional) Description of the image placement policy.
* `placement_type`: (Required) Type of the image placement policy. Valid values "HARD", "SOFT"
* `image_entity_filter`: (Required) Category-based entity filter.
* `cluster_entity_filter`: (Required) Category-based entity filter.
* `enforcement_state`: (Optional) Enforcement status of the image placement policy. Valid values "ACTIVE", "SUSPENDED"
* `action` : (Optional) Action to be performed on the image placement policy. Valid values "RESUME", "SUSPEND"

### image_entity_filter
* `type`: (Required) Filter matching type. Valid values "CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"
* `category_ext_ids`: Array of strings

### cluster_entity_filter
* `type`: (Required) Filter matching type. Valid values "CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"
* `category_ext_ids`: Array of strings


See detailed information in [Nutanix Image Placement Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1)