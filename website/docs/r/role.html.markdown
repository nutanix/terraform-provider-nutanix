---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role"
sidebar_current: "docs-nutanix-resource-role"
description: |-
  This operation submits a request to create a role based on the input parameters.
---

# nutanix_role

Provides a resource to create a role based on the input parameters.

## Example Usage

```hcl
resource "nutanix_role" "test" {
	name        = "NAME"
	description = "DESCRIPTION"
	permission_reference_list {
		kind = "permission"
		uuid = "ID OF PERMISSION"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "ID OF PERMISSION"
	}
	permission_reference_list {
		kind = "permission"
		uuid = "ID OF PERMISSION"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Optional) Name of the role.
* `description`: - (Optional) The description of the role.
* `categories`: - (Optional) Categories for the role.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.
* `permission_reference_list`: - (Required) List of permission references.

## Attributes Reference

The following attributes are exported:

* `metadata`: - The role kind metadata.
* `api_version` - The version of the API.
* `state`: - The state of the role.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when role was last updated.
* `uuid`: - Role UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when role was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - Role name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Reference

The `project_reference`, `owner_reference` attributes supports the following:

* `kind`: - The kind name (Default value: project)(Required).
* `name`: - the name(Optional).
* `uuid`: - the UUID(Required).

For `permission_reference_list` are the same as reference but used as array.

See detailed information in [Nutanix Roles](https://www.nutanix.dev/api_references/prism-central-v3/#/acde86df0c657-create-a-role).
