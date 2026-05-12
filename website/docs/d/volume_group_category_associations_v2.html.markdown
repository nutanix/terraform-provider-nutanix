---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_category_associations_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-category-associations-v2"
description: |-
  Query the category details associated with a Volume Group.
---

# nutanix_volume_group_category_associations_v2

Query the category details that are associated with the Volume Group identified by {volumeGroupExtId}. Deprecated: This API has been deprecated.

## Example Usage

```hcl
data "nutanix_volume_group_category_associations_v2" "example" {
  volume_group_ext_id = "d09aeec9-5bb7-4bfd-9717-a051178f6e7c"
}
```

## Argument Reference

The following arguments are supported:

* `volume_group_ext_id`: -(Required) The external identifier of a Volume Group.

## Attributes Reference

The following attributes are exported:

* `category_associations`: - List of category details associated with the Volume Group.

### Category Associations

Each element in `category_associations` has the following fields:

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `name`: - Name of the entity represented by this reference.
* `entity_type`: - The entity type.
* `uris`: - URI of entity represented by this reference.
