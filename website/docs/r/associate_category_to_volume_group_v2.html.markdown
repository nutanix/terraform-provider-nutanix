---
layout: "nutanix"
page_title: "NUTANIX: nutanix_associate_category_to_volume_group_v2"
sidebar_current: "docs-nutanix-resource-associate-category-to-volume-group-v2"
description: |-
  This is a action module. It will associate categories to Volume groups in every apply and not maintain state. Terraform destroy will dissociate the categories from given volume group.

---

# nutanix_associate_category_to_volume_group_v2

Provides a resource to Creates a new Volume Disk.

## Example Usage

```hcl

# Associate categories to volume group
resource "nutanix_associate_category_to_volume_group_v2" "example"{
  ext_id = "f0c0a4ac-c734-4770-b5d7-eca6793eeeb7" # Volume Group extId
  categories{
    ext_id = "85e68112-5b2b-4220-bc8d-e529e4bf420e" # Category extId
  }
  categories{
    ext_id = "45588de3-7c18-4230-a147-7e26ad92d8a6" # Category extId
  }
  categories{
    ext_id = "1c6638f2-5215-4086-8f21-a30e75cb8068" # Category extId
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
