---
layout: "nutanix"
page_title: "NUTANIX: nutanix_project_v2"
sidebar_current: "docs-nutanix-resource-project-v2"
description: |-
  Creates and manages a multidomain project.
---

# nutanix_project_v2

Creates and manages a multidomain project.

## Example Usage

```hcl
resource "nutanix_project_v2" "example" {
  name        = "my-multidomain-project"
  description = "Project for multidomain namespace"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the project.
* `description` - (Optional) Description of the project.

## Attributes Reference

The following attributes are exported:

* `ext_id` - A globally unique identifier of the project.
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
