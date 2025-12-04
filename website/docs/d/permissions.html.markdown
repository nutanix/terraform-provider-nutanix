---
layout: "nutanix"
page_title: "NUTANIX: nutanix_permissions"
sidebar_current: "docs-nutanix-datasource-permissions"
description: |-
  Provides a datasource to retrieve all the permissions.
---

# nutanix_permission

Provides a datasource to retrieve all the permissions.

## Example Usage

```hcl
data "nutanix_permission" "permissions" {}
```

## Attributes Reference

The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of Permissions

# Entities

The entities attribute element contains the following attributes:

* `metadata`: The permission kind metadata.
* `categories`: The categories for this resource.
* `owner_reference`: The reference to a user.
* `project_reference`: The reference to a project.
* `name` The name for the permission.
* `state`: The state of the permission.
* `description` A description for the permission.
* `operation` The operation that is being performed on a given kind.
* `kind` The kind on which the operation is being performed.
* `fields` . The fields that can/cannot be accessed during the specified operation. field_name_list will be a list of fields. e.g. if field_mode = disallowed, field_name_list = [“xyz”] then the list of allowed fields is ALL fields minus xyz. Seee [Field](#field) for more info.

### Field

The field attribute exports the following:

* `field_mode` Allow or disallow the fields mentioned.
* `field_name_list` The list of fields.

### Metadata

The metadata attribute exports the following:

* `last_update_time` - UTC date and time in RFC-3339 format when the permission was last updated.
* `uuid` - permission UUID.
* `creation_time` - UTC date and time in RFC-3339 format when the permission was created.
* `spec_version` - Version number of the latest spec.
* `spec_hash` - Hash of the spec. This will be returned from server.
* `name` - permission name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories
The categories attribute supports the following:

* `name` - the key name.
* `value` - value of the key.

### Reference
The `project_reference`, `owner_reference` attributes supports the following:

* `kind` - (Required) The kind name (Default value: `project`).
* `name` - the name.
* `uuid` - (Required) the UUID.

See detailed information in [Nutanix Permissions](https://www.nutanix.dev/api_references/prism-central-v3/#/0fabd54b59459-list-the-permissions).
