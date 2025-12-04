---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user"
sidebar_current: "docs-nutanix-resource-user"
description: |-
  This operation submits a request to create a user based on the input parameters.
---

# nutanix_user

Provides a resource to create a user based on the input parameters.

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
```


```hcl
resource "nutanix_user" "user" {
	identity_provider_user {
		username = "username"
		identity_provider_reference {
		uuid = "<identity-provider-uuid>"
		}
	}
}
```

## Argument Reference

The following arguments are supported:

* `directory_service_user`: - (Optional) The directory service user configuration. See below for more information.
* `identity_provider_user`: - (Optional) (Optional) The identity provider user configuration. See below for more information.
* `categories`: - (Optional) Categories for the Access Control Policy.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.

### Directory Service User

The directory_service_user argument supports the following:

* `user_principal_name`: - (Optional) The UserPrincipalName of the user from the directory service.
* `directory_service_reference`: - (Optional) The reference to a directory service. See #reference for to look the supported attributes. 

### Identity Provider User

The identity_provider_user argument supports the following:

* `username`: - (Optional) The username from identity provider. Name ID for SAML Identity Provider. 
* `identity_provider_reference`: - (Optional) The reference to a identity provider. See #reference for to look the supported attributes. 

### Context List

The context_list attribute supports the following:

* `scope_filter_expression_list`: - (Optional) The device ID which is used to uniquely identify this particular disk.
* `entity_filter_expression_list` - (Required) A list of Entity filter expressions.

## Attributes Reference

The following attributes are exported:

* `metadata`: - The vm kind metadata.
* `api_version` - The version of the API.
* `state`: - The state of the entity.
* `name`: - The name of the user.
* `user_type`: - The name of the user.
* `display_name`: - The display name of the user (common name) provided by the directory service.
* `project_reference_list`: - A list of projects the user is part of. See #reference for more details.
* `access_control_policy_reference_list`: - List of ACP references. See #reference for more details.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when vm was last updated.
* `uuid`: - vm UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when vm was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - vm name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Reference

The `project_reference`, `owner_reference`, `role_reference` `directory_service_reference` attributes supports the following:

* `kind`: - The kind name. (Default depends on the resource you are referencing)
* `name`: - the name(Optional).
* `uuid`: - the UUID(Required).

For `access_control_policy_reference_list` and `project_reference_list` are the same as reference but used as list.

See detailed information in [Nutanix Users](https://www.nutanix.dev/api_references/prism-central-v3/#/e7c2691629db9-create-a-new-user).
