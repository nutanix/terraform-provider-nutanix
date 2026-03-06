---
layout: "nutanix"
page_title: "NUTANIX: nutanix_resource_groups_v2"
sidebar_current: "docs-nutanix-datasource-resource-groups-v2"
description: |-
  List the multidomain resource groups defined on the system.
---

# nutanix_resource_groups_v2

List the multidomain resource groups defined on the system.

## Example Usage

```hcl
data "nutanix_resource_groups_v2" "example" {}
```

## Attributes Reference

The following attributes are exported:

* `resource_groups` - List of resource groups.

## Resource Groups

The `resource_groups` attribute is a list of resource group objects. Each resource group supports the following attributes:

* `ext_id` - A globally unique identifier of the resource group.
* `name` - Name of the resource group.
* `project_ext_id` - External identifier of the project.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `created_by` - User who created the resource group.
* `last_updated_by` - User who last updated the resource group.
* `create_time` - Creation time (RFC3339).
* `last_update_time` - Last update time (RFC3339).
* `placement_targets` - List of placement targets.
* `links` - A HATEOAS style link for the response.

### Placement Targets

Each placement target supports:

* `cluster_ext_id` - External identifier of the cluster.
* `storage_containers` - List of storage container details, each with `ext_id`.

### Links

The `links` attribute supports the following:

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object.
