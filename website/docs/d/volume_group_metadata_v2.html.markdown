---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_metadata_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-metadata-v2"
description: |-
  Describes Volume Group metadata.
---

# nutanix_volume_group_metadata_v2

Query for metadata information which is associated with the Volume Group identified by {extId}. Deprecated: This API has been deprecated.

## Example Usage

```hcl
data "nutanix_volume_group_metadata_v2" "example" {
  volume_group_ext_id = "d09aeec9-5bb7-4bfd-9717-a051178f6e7c"
}
```

## Argument Reference

The following arguments are supported:

* `volume_group_ext_id`: -(Required) The external identifier of a Volume Group.

## Attributes Reference

The following attributes are exported:

* `category_ids`: - A list of globally unique identifiers that represent all the categories the resource is associated with.
* `owner_reference_id`: - A globally unique identifier that represents the owner of this resource.
* `owner_user_name`: - The userName of the owner of this resource.
* `project_name`: - The name of the project this resource belongs to.
* `project_reference_id`: - A globally unique identifier that represents the project this resource belongs to.
