---
layout: "nutanix"
page_title: "NUTANIX: nutanix_iam_entity_v2"
sidebar_current: "docs-nutanix-datasource-entity-v2"
description: |-
  Provides a datasource to retrieve an IAM Entity by its external identifier.
---

# nutanix_iam_entity_v2

Provides a datasource to retrieve an IAM Entity by its external identifier. Entities are used in authorization policies (e.g. user, role, cluster).

## Example Usage

```hcl
# Get entity by ext_id
data "nutanix_iam_entity_v2" "example" {
  ext_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "entity_name" {
  value = data.nutanix_iam_entity_v2.example.name
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) External identifier of the IAM Entity.

## Attributes Reference

The following attributes are exported:

* `id` - The external identifier of the entity (same as `ext_id`).
* `tenant_id` - Tenant ID for the Entity. A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name` - Name of the Entity. Unique name of the entity.
* `description` - Description of the Entity.
* `display_name` - Display name for the Entity. UI display name of the entity.
* `client_name` - Client that created the entity.
* `search_url` - Search URL for the Entity. URL provided by the client to search the entities.
* `created_time` - Creation time of the Entity.
* `last_updated_time` - Last updated time of the Entity.
* `created_by` - User or Service that created the Entity.
* `attribute_list` - List of attributes for the Entity (used in authorization policy filters).
* `is_logical_and_supported_for_attributes` - Whether logical AND is supported for attributes. Indicates whether the entity supports scoping using multiple attributes which will result in a logical AND.

### links

Each link in the `links` list supports the following:

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object.

### attribute_list

Each element in `attribute_list` supports the following:

* `tenant_id` - Tenant identifier for the attribute.
* `ext_id` - External identifier of the attribute.
* `links` - HATEOAS links for the attribute (each with `href` and `rel`).
* `display_name` - Display name of the entity's attribute.
* `name` - Name of the entity's attribute used in Authorization Policy filters.
* `supported_operator` - List of supported operators for this entity attribute.
* `attribute_values` - List of attribute values supported for access control.

See [Nutanix Get Entity v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Entities/operation/getEntityById) for API details.
