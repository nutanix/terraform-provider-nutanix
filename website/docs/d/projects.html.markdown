---
layout: "nutanix"
page_title: "NUTANIX: nutanix_projects"
sidebar_current: "docs-nutanix-datasource-projects"
description: |-
 Describes a Projects
---

# nutanix_projects

Describes Projects

## Example Usage

```hcl
data "nutanix_projects" "projects" {}
```


## Attribute Reference

The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of Projects

# Entities

The entities attribute element contains the followings attributes:

* `name` The name for the project.
* `description` A description for project.

### Resource Domain
* `resource_domain` The status for a resource domain (limits and values)
* `resource_domain.resources` Array of the utilization/limit for resource types
* `resource_domain.resources.#.limit` The resource consumption limit (unspecified is unlimited)
* `resource_domain.resources.#.resource_type` The type of resource (for example storage, CPUs)
* `resource_domain.resources.#.units` - The units of the resource type
* `resource_domain.resources.#.value` - The amount of resource consumed

### Account Reference List
* `account_reference_list`
* `account_reference_list.#.kind`
* `account_reference_list.#.uuid`
* `account_reference_list.#.name`

### Environment Reference List
* `environment_reference_list`
* `environment_reference_list.#.kind`
* `environment_reference_list.#.uuid`
* `environment_reference_list.#.name`

### Default Subnet Reference Map
* `default_subnet_reference`
* `default_subnet_reference.kind`
* `default_subnet_reference.uuid`
* `default_subnet_reference.name`

### user_reference_list
* `user_reference_list`
* `user_reference_list.#.kind`
* `user_reference_list.#.uuid`
* `user_reference_list.#.name`

### External User Group Reference List
* `external_user_group_reference_list`
* `external_user_group_reference_list.#.kind`
* `external_user_group_reference_list.#.uuid`
* `external_user_group_reference_list.#.name`

### Subnet Reference List
* `subnet_reference_list`
* `subnet_reference_list.#.kind`
* `subnet_reference_list.#.uuid`
* `subnet_reference_list.#.name`

### External Network List
* `subnet_reference_list`
* `subnet_reference_list.#.uuid`
* `subnet_reference_list.#.name`

### Resource Domain
* `resource_domain.resources.#.units` - The units of the resource type
* `resource_domain.resources.#.value` - The amount of resource consumed

### Metadata
The metadata attribute exports the following:

* `last_update_time` - UTC date and time in RFC-3339 format when vm was last updated.
* `uuid` - vm UUID.
* `creation_time` - UTC date and time in RFC-3339 format when vm was created.
* `spec_version` - Version number of the latest spec.
* `spec_hash` - Hash of the spec. This will be returned from server.
* `name` - vm name.

### Categories
The categories attribute supports the following:

* `name` - the key name.
* `value` - value of the key.

### Reference
The `project_reference`, `owner_reference` attributes supports the following:

* `kind` - (Required) The kind name (Default value: `project`).
* `name` - (Optional) the name.
* `uuid` - (Required) the UUID.


See detailed information in [Nutanix Projects](https://www.nutanix.dev/reference/prism_central/v3/api/projects/postprojectslist).
