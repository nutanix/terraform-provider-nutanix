---
layout: "nutanix"
page_title: "NUTANIX: nutanix_resource_group_v2"
sidebar_current: "docs-nutanix-datasource-resource-group-v2"
description: |-
  Fetches the resource group identified by an external identifier.
---

# nutanix_resource_group_v2

Fetches the resource group identified by an external identifier.

## Example Usage

```hcl
data "nutanix_resource_group_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) The external identifier of the resource group.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the resource group.
* `project_ext_id` - External identifier of the project this resource group belongs to.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `created_by` - User who created the resource group.
* `last_updated_by` - User who last updated the resource group.
* `create_time` - Creation time (RFC3339).
* `last_update_time` - Last update time (RFC3339).
* `placement_targets` - List of placement targets (cluster and storage containers).
* `links` - A HATEOAS style link for the response.

## Placement Targets

The `placement_targets` attribute supports the following:

* `cluster_ext_id`:- UUID of the AOS cluster.
* `storage_containers`:- List of storage containers available for this cluster target.

## Storage Containers

The `storage_containers` attribute supports the following:
* `ext_id`:- UUID of the storage container.

## Links

The `links` attribute supports the following:

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object.
