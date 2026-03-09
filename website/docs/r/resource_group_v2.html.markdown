---
layout: "nutanix"
page_title: "NUTANIX: nutanix_resource_group_v2"
sidebar_current: "docs-nutanix-resource-resource-group-v2"
description: |-
  Creates and manages a multidomain resource group.
---

# nutanix_resource_group_v2

Creates and manages a multidomain resource group.

## Example Usage

```hcl
resource "nutanix_resource_group_v2" "example" {
  name           = "my-resource-group"
  project_ext_id = nutanix_project_v2.my_project.ext_id
}
```

With placement targets:

```hcl
resource "nutanix_resource_group_v2" "example" {
  name           = "my-resource-group"
  project_ext_id = nutanix_project_v2.my_project.ext_id

  placement_targets {
    cluster_ext_id = "cluster-uuid"

    storage_containers {
      ext_id = "storage-container-uuid"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource group.
* `project_ext_id` - (Optional) External identifier of the project this resource group belongs to.
* `tenant_id` - (Optional) Tenant identifier.
* `placement_targets` - (Optional) List of placement targets. Each block supports:
  * `cluster_ext_id` - (Optional) External identifier of the cluster.
  * `storage_containers` - (Optional) List of storage container details. Each block supports:
    * `ext_id` - (Required) External identifier of the storage container.

## Attributes Reference

The following attributes are exported:

* `ext_id` - A globally unique identifier of the resource group.
* `created_by` - User who created the resource group.
* `last_updated_by` - User who last updated the resource group.
* `create_time` - Creation time (RFC3339).
* `last_update_time` - Last update time (RFC3339).
* `links` - A HATEOAS style link for the response.

### Links

The `links` attribute supports the following:

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object.
