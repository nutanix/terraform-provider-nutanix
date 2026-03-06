---
layout: "nutanix"
page_title: "NUTANIX: nutanix_projects_v2"
sidebar_current: "docs-nutanix-datasource-projects-v2"
description: |-
  List the multidomain projects defined on the system.
---

# nutanix_projects_v2

List the multidomain projects defined on the system.

## Example Usage

```hcl
data "nutanix_projects_v2" "example" {}
```

## Attributes Reference

The following attributes are exported:

* `projects` - List of projects.

## Projects

The `projects` attribute is a list of project objects. Each project supports the following attributes:

* `ext_id` - A globally unique identifier of the project.
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
