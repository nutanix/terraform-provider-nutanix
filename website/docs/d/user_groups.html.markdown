---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_groups"
sidebar_current: "docs-nutanix-datasource-user-groups"
description: |-
  Provides a datasource to retrieve all the user groups.
---

# nutanix_user_groups

Provides a datasource to retrieve all the user groups.

## Example Usage

```hcl
data "nutanix_user_groups" "usergroups" {}
```

## Argument Reference

The following attributes are exported:

# Entities

The entities attribute element contains the following attributes:

The following attributes are exported:

* `api_version` - The version of the API.
* `metadata`: - The user group kind metadata.
* `categories`: - The Categories for the user group.
* `owner_reference`: - The reference to a user.
* `project_reference`: - The reference to a project.
* `user_group_type`: - The type of the user group.
* `display_name`: - The display name of the user group.
* `directory_service_user_group`: - A Directory Service User Group.
* `project_reference_list`: - A list of projects the user is part of. See #reference for more details.
* `access_control_policy_reference_list`: - List of ACP references. See #reference for more details.
* `state`: - The state of the entity.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when user group was last updated.
* `uuid`: - User group UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when user group was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - User group name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Directory Service User Group

The directory_service_user argument supports the following:

* `distinguished_name`: - The Distinguished name for the user group
* `directory_service_reference`: - The reference to a directory service. See #reference for to look the supported attributes. 


### Reference

The `project_reference`, `owner_reference`, `role_reference` `directory_service_reference` attributes supports the following:

* `kind`: - The kind name. (Default depends on the resource you are referencing)
* `name`: - the name(Optional).
* `uuid`: - the UUID(Required).

For `access_control_policy_reference_list` and `project_reference_list` are the same as reference but used as list.

See detailed information in [Nutanix Users](https://www.nutanix.dev/api_references/prism-central-v3/#/6016c890e9122-get-a-list-of-existing-user-groups).
