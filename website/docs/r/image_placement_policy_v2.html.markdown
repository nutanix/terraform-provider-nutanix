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
resource "nutanix_image_placement_policy_v2" "example"{
	name           = "image_placement_policy"
	description    = "%[2]s"
	placement_type = "SOFT"
	cluster_entity_filter {
		category_ext_ids = [
			"ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
		]
		type = "CATEGORIES_MATCH_ALL"
	}
	image_entity_filter {
		category_ext_ids = [
			"ab520e1d-4950-1db1-917f-a9e2ea35b8e3",
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

## Import

This helps to manage existing entities which are not created through terraform. Image Placement Policies can be imported using the `UUID`. (ext_id in v4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_image_placement_policy_v2" "import_ipp"{}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_image_placement_policies_v2" "list_ipps"{}
terraform import nutanix_image_placement_policy_v2.import_ipp <UUID>
```

See detailed information in [Nutanix Create Image Placement Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/ImagePlacementPolicies/operation/createPlacementPolicy)
