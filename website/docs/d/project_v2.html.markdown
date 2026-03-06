---
layout: "nutanix"
page_title: "NUTANIX: nutanix_project_v2"
sidebar_current: "docs-nutanix-datasource-project-v2"
description: |-
  Fetches the multidomain project identified by an external identifier.
---

# nutanix_project_v2

Fetches the multidomain project identified by an external identifier.

## Example Usage

```hcl
data "nutanix_project_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) The external identifier of the project.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the project.
* `description` - Description of the project.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `state` - State of the project.
* `is_default` - Indicates if this is the default project.
* `is_system_defined` - Indicates if this project is system defined.
* `created_by` - User who created the project.
* `updated_by` - User who last updated the project.
* `created_timestamp` - Creation timestamp in microseconds.
* `modified_timestamp` - Last modified timestamp in microseconds.
* `links` - A HATEOAS style link for the response.

### Links

The `links` attribute supports the following:

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object.
