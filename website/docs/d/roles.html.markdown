---
layout: "nutanix"
page_title: "NUTANIX: nutanix_roles"
sidebar_current: "docs-nutanix-resource-roles"
description: |-
  Describes a list of roles
---

# nutanix_roles

Describes a list of roles.

## Example Usage

```hcl
data "nutanix_roles" "test" {}
```

## Attribute Reference

The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of Roles

# Entities

The entities attribute element contains the followings attributes:

* `name`: - Name of the role.
* `description`: - The description of the role.
* `categories`: - Categories for the role.
* `project_reference`: - The reference to a project.
* `owner_reference`: - The reference to a user.
* `project_reference`: - The reference to a project.
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

See detailed information in [Nutanix Roles](https://www.nutanix.dev/api_references/prism-central-v3/#/3de7424ca8221-list-the-roles).
