---
layout: "nutanix"
page_title: "NUTANIX: nutanix_scope_templates_v2"
sidebar_current: "docs-nutanix-datasource-scope-templates-v2"
description: |-
  Provides a datasource to list all scope templates.
---
# nutanix_scope_templates_v2

Provides a datasource to list all scope templates.

## Example Usage

```hcl
# List all scope templates
data "nutanix_scope_templates_v2" "all" {}

# List scope templates with filter
data "nutanix_scope_templates_v2" "filtered" {
  filter = "displayName eq 'my_scope_template'"
}
```

## Argument Reference

The following arguments are supported:

* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attribute Reference

The following attributes are exported:

* `scope_templates`: List of scope templates.

### Scope Templates

Each scope template in the list contains the following attributes:

* `ext_id`: External identifier of the scope template.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.
* `display_name`: The display name for the scope template.
* `description`: Description of the scope template.
* `entities`: List of entities being scoped for the template.
* `created_by`: Service name that created the scope template.
* `created_time`: The creation time of the scope template.

### Links

The links attribute supports the following:

* `href`: The URL at which the entity described by the link can be accessed.
* `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.

### Entities

The entities attribute supports the following:

* `entity_filter`: Information of the entity filter present in the EntityFilter object.

See detailed information in [Nutanix List Scope Templates v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/ScopeTemplates/operation/listScopeTemplates).
