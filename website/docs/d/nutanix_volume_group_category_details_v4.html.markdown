---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_category_details_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-category-details-v4"
description: |-
  Describes a List all the category details that are associated with the Volume Group.
---

# nutanix_volume_group_disks_v2

Query the category details that are associated with the Volume Group identified by {volumeGroupExtId}.
## Example Usage

```hcl

# List of all category details that are associated with the Volume Group.
data "nutanix_volume_group_category_details_v2" "vg_cat_example"{
  ext_id = var.volume_group_ext_id
  limit  = 6
}

```

##  Argument Reference

The following arguments are supported:

* `volume_group_ext_id`: -(Required) The external identifier of the Volume Group.
* `page`: - A query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.


## Attributes Reference
The following attributes are exported:

* `category_details`: - List of all category details that are associated with the Volume Group.

### Category Details

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `name`: - The name of the category detail.
* `uris`: - The uri list of the category detail.
* `entity_type`: - SThe Entity Type of the category detail.



See detailed information in [Nutanix Volumes](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0.b1).
