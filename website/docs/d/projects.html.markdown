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
* `account_reference_list` - List of accounts associated with the project.
* `account_reference_list.#.kind` - The kind name. Default value is `account`
* `account_reference_list.#.uuid` - The UUID of an account.
* `account_reference_list.#.name` - The name of an account.

### Environment Reference List
* `environment_reference_list` - List of environments associated with the project.
* `environment_reference_list.#.kind` - The kind name. Default value is `environment`
* `environment_reference_list.#.uuid` - The UUID of an environment.
* `environment_reference_list.#.name` - The name of an environment.

### Default Subnet Reference Map
* `default_subnet_reference` - Reference to a subnet.
* `default_subnet_reference.kind` - The kind name. Default value is `subnet`
* `default_subnet_reference.uuid` - The UUID of a subnet.
* `default_subnet_reference.name` - The name of a subnet.

### user_reference_list
* `user_reference_list` - List of users in the project.
* `user_reference_list.#.kind` - The kind name. Default value is `user`
* `user_reference_list.#.uuid` - The UUID of a user
* `user_reference_list.#.name` - The name of a user.

### External User Group Reference List
* `external_user_group_reference_list` - List of directory service user groups. These groups are not managed by Nutanix.
* `external_user_group_reference_list.#.kind` - The kind name. Default value is `user_group`
* `external_user_group_reference_list.#.uuid` - The UUID of a user_group
* `external_user_group_reference_list.#.name` - The name of a user_group

### Subnet Reference List
* `subnet_reference_list` - List of subnets for the project.
* `subnet_reference_list.#.kind` - The kind name. Default value is `subnet`
* `subnet_reference_list.#.uuid` - The UUID of a subnet
* `subnet_reference_list.#.name` - The name of a subnet.

### External Network List
* `external_network_list` - List of external networks associated with the project.
* `external_network_list.#.uuid` - The UUID of a network.
* `external_network_list.#.name` - The name of a network.

### Resource Domain
* `resource_domain.resources.#.units` - The units of the resource type
* `resource_domain.resources.#.value` - The amount of resource consumed

### Metadata
The metadata attribute exports the following:

* `last_update_time` - UTC date and time in RFC-3339 format when the project was last updated.
* `uuid` - Project UUID.
* `creation_time` - UTC date and time in RFC-3339 format when the project was created.
* `spec_version` - Version number of the latest spec.
* `spec_hash` - Hash of the spec. This will be returned from server.
* `name` - Project name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories
The categories attribute supports the following:

* `name` - the key name.
* `value` - value of the key.

### Reference
The `project_reference`, `owner_reference` attributes supports the following:

* `kind` - (Required) The kind name (Default value: `project`).
* `name` - (Optional) the name.
* `uuid` - (Required) the UUID.


See detailed information in [Nutanix Projects](https://www.nutanix.dev/api_references/prism-central-v3/#/226263506f77a-get-a-list-of-existing-projects).
