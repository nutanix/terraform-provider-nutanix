---
layout: "nutanix"
page_title: "NUTANIX: nutanix_scope_template_v2"
sidebar_current: "docs-nutanix-datasource-scope-template-v2"
description: |-
  Provides a datasource to retrieve a scope template by its external identifier.
---
# nutanix_scope_template_v2

Provides a datasource to retrieve a scope template by its external identifier.

## Example Usage

```hcl
data "nutanix_scope_template_v2" "example" {
  ext_id = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) External identifier of the scope template.

## Attribute Reference

The following attributes are exported:

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

See detailed information in [Nutanix Get Scope Template v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/ScopeTemplates/operation/getScopeTemplateById).
