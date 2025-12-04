---
layout: "nutanix"
page_title: "NUTANIX: nutanix_users"
sidebar_current: "docs-nutanix-datasource-users"
description: |-
  This operation retrieves a list of all the users.
---

# nutanix_users

Provides a datasource to retrieve all the users.

## Example Usage

```hcl
resource "nutanix_user" "user" {
	directory_service_user {
		user_principal_name = "test-user@ntnxlab.local"
		directory_service_reference {
		uuid = "<directory-service-uuid>"
		}
	}
}

data "nutanix_user" "users" {}
```


## Attribute Reference

The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of Users

# Entities

The entities attribute element contains the following attributes:

* `metadata`: - The user kind metadata.
* `api_version` - The version of the API.
* `state`: - The state of the entity.
* `name`: - The name of the user.
* `user_type`: - The name of the user.
* `display_name`: - The display name of the user (common name) provided by the directory service.
* `project_reference_list`: - A list of projects the user is part of. See #reference for more details.
* `access_control_policy_reference_list`: - List of ACP references. See #reference for more details.
* `directory_service_user`: - (Optional) The directory service user configuration. See below for more information.
* `identity_provider_user`: - (Optional) (Optional) The identity provider user configuration. See below for more information.
* `categories`: - (Optional) Categories for the user.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when user was last updated.
* `uuid`: - User UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when user was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - User name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Directory Service User

The directory_service_user argument supports the following:

* `user_principal_name`: - (Optional) The UserPrincipalName of the user from the directory service.
* `directory_service_reference`: - (Optional) The reference to a directory service. See #reference for to look the supported attributes. 

### Identity Provider User

The identity_provider_user argument supports the following:

* `username`: - (Optional) The username from identity provider. Name ID for SAML Identity Provider. 
* `identity_provider_reference`: - (Optional) The reference to a identity provider. See #reference for to look the supported attributes. 

### Reference

The `project_reference`, `owner_reference`, `role_reference` `directory_service_reference` attributes supports the following:

* `kind`: - The kind name. (Default depends on the resource you are referencing)
* `name`: - the name(Optional).
* `uuid`: - the UUID(Required).

For `access_control_policy_reference_list` and `project_reference_list` are the same as reference but used as list.

See detailed information in [Nutanix Users](https://www.nutanix.dev/api_references/prism-central-v3/#/7ddb9e0b5e5df-get-a-list-of-existing-users).
