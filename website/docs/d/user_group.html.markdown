---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_group"
sidebar_current: "docs-nutanix-datasource-user-group"
description: |-
  This operation retrieves a user based on the input parameters.
---

# nutanix_user_group

Provides a datasource to retrieve a user group based on the input parameters.

## Example Usage

```hcl

//Retrieve by UUID
data "nutanix_user_group" "usergroup" {
	user_group_id = "dd30a856-8e72-4158-b716-98455ceda220"
}

//Retrieve by Name
data "nutanix_user_group" "usergroupbyname" {
	user_group_name = "example-group-1"
}

//Retrieve by Distinguished Name
data "nutanix_user_group" "test" {
	user_group_distinguished_name = "cn=example-group-1,cn=users,dc=ntnxlab,dc=local"
}
```




## Argument Reference

The following arguments are supported:

* `user_group_id`: - (Optional) The UUID for the user group
* `user_group_name`: - (Optional) The name for the user group
* `user_group_distinguished_name` - (Optional) The distinguished name for the user group

## Attributes Reference

The following attributes are exported:

* `api_version` - The version of the API.
* `metadata`: - The user group kind metadata.
* `categories`: - The Distinguished Categories for the user group.
* `owner_reference`: - The reference to a user.
* `project_reference`: - The Distinguished The reference to a project.
* `user_group_type`: - The type of the user group.
* `display_name`: - The display name of the user group.
* `directory_service_user_group`: - A Directory Service User Group.
* `project_reference_list`: - A list of projects the user is part of. See #reference for more details.
* `access_control_policy_reference_list`: - List of ACP references. See #reference for more details.
* `state`: - The state of the entity.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when the user group was last updated.
* `uuid`: - User group UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when the user group was created.
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

See detailed information in [Nutanix Users](https://www.nutanix.dev/api_references/prism-central-v3/#/ec9f993c00b11-get-a-existing-user-group).
