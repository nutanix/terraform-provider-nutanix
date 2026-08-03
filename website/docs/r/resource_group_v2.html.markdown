---
layout: "nutanix"
page_title: "NUTANIX: nutanix_resource_group_v2"
sidebar_current: "docs-nutanix-resource-resource-group-v2"
description: |-
  Creates and manages a resource group.
---

# nutanix_resource_group_v2

Creates and manages a resource group.

## Example Usage

```hcl
resource "nutanix_resource_group_v2" "resource_group" {
  name           = "resource-group-example"
  project_ext_id = "project_ext_id"
  placement_targets {
    cluster_ext_id = "cluster_ext_id"
    storage_containers {
      ext_id = "storage_container_ext_id"
    }
    storage_containers {
      ext_id = "storage_container_ext_id"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name`:- (Required) Name of the resource group.
* `project_ext_id`:- (Required) External identifier of the project this resource group belongs to.
* `placement_targets`:- (Optional) List of placement targets.

## Placement Targets

The `placement_targets` attribute supports the following:

* `cluster_ext_id`:- (Required) UUID of the AOS cluster.
* `storage_containers`:- (Optional) List of storage containers available for this cluster target.

## Storage Containers

The `storage_containers` attribute supports the following:
* `ext_id`:- (Required) UUUID of the storage container.


## Attributes Reference

The following attributes are exported:
* `name`:- Resource Group name 
* `project_ext_id`:- External identifier of the project this resource group belongs to.
* `placement_targets`:- List of placement targets.
* `ext_id`:- A globally unique identifier of the resource group.
* `created_by`:- User who created the resource group.
* `last_updated_by`:- User who last updated the resource group.
* `create_time`:- Creation time (RFC3339).
* `last_update_time`:- Last update time (RFC3339).
* `links` - A HATEOAS style link for the response.

## Links

The `links` attribute supports the following:

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object.
