---
layout: "nutanix"
page_title: "NUTANIX: nutanix_associate_category_to_volume_group_v2"
sidebar_current: "docs-nutanix-resource-associate-category-to-volume-group-v2"
description: |-
  This operation submits a request to Creates a new Volume Disk.
---

# nutanix_volume_group_v2

Provides a resource to Creates a new Volume Disk.

## Example Usage

```hcl

resource "nutanix_volume_group_v2" "example"{
  name                               = "test_volume_group"
  description                        = "Test Volume group with min spec and no Auth"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  target_name                        = "volumegroup-test-0"
  created_by                         = "Test"
  cluster_reference                  = "<Cluster uuid>"
  iscsi_features {
    enabled_authentications = "CHAP"
    target_secret           = "1234567891011"
  }

  storage_features {
    flash_mode {
      is_enabled = true
    }
  }
  usage_type = "USER"
  is_hidden  = false

  lifecycle {
    ignore_changes = [
      iscsi_features[0].target_secret
    ]
  }
}


# List categories
data "nutanix_categories_v2" "categories"{}

# Associate categories to volume group
resource "nutanix_associate_category_to_volume_group_v2" "example"{
  ext_id = nutanix_volume_group_v2.example.id
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.0.ext_id
  }
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.1.ext_id
  }
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.2.ext_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `ext_id `: -(Required) The external identifier of the Volume Group.
* `categories`: -(Required) The category to be associated/disassociated with the Volume Group. This is a mandatory field.


### categories

The categories attribute supports the following:

* `ext_id`: -(Required) The external identifier of the category.
* `name`: -(Optional) Name of entity that's represented by this reference
* `uris`: -(Optional) URI of entities that's represented by this reference.
* `entity_type`: -(Optional) Type of entity that's represented by this reference. Default value is "CATEGORY". Valid values are:
  * "CATEGORY".

See detailed information in [Nutanix Associate/Disassociate category to/from a Volume Group V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.0#tag/VolumeGroups/operation/associateCategory).
