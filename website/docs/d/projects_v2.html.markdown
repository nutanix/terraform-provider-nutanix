---
layout: "nutanix"
page_title: "NUTANIX: nutanix_projects_v2"
sidebar_current: "docs-nutanix-datasource-projects-v2"
description: |-
  List the projects defined on the system.
---

# nutanix_projects_v2

List the projects defined on the system.

## Example Usage

```hcl
data "nutanix_projects_v2" "example" {}
```

## Argument Reference
The following arguments are supported:

- `page`:- A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`:- A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`:- A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - createdBy
  - createdTimestamp
  - description
  - extId
  - id
  - isDefault
  - isSystemDefined
  - modifiedTimestamp
  - name
  - state
  - updatedBy
- `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - createdBy
  - createdTimestamp
  - description
  - extId
  - id
  - modifiedTimestamp
  - name
  - state
  - updatedBy
- `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields:
  - createdBy
  - createdTimestamp
  - description
  - id
  - isDefault
  - isSystemDefined
  - modifiedTimestamp
  - name
  - state
  - updatedBy


## Attributes Reference

The following attributes are exported:

* `projects`:- List of projects.

## Projects

The `projects` attribute is a list of project objects. Each project supports the following attributes:

* `ext_id`:- A globally unique identifier of the project.
* `name`:- Name of the project.
* `id`:- ID of the project
* `description`:- Description of the project.
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity.
* `state`:- State of the project.
* `is_default`:- Indicates if this is the default project.
* `is_system_defined`:- Indicates if this project is system defined.
* `created_by`:- User who created the project.
* `updated_by`:- User who last updated the project.
* `created_timestamp`:- Creation timestamp in microseconds.
* `modified_timestamp`:- Last modified timestamp in microseconds.
* `links`:- A HATEOAS style link for the response.

### Links

The `links` attribute supports the following:

* `href`:- The URL at which the entity described by the link can be accessed.
* `rel`:- A name that identifies the relationship of the link to the object.
