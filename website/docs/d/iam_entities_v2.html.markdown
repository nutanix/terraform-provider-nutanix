---
layout: "nutanix"
page_title: "NUTANIX: nutanix_iam_entities_v2"
sidebar_current: "docs-nutanix-datasource-entities-v2"
description: |-
  Provides a datasource to list IAM Entities with optional filtering and pagination.
---

# nutanix_iam_entities_v2

Provides a datasource to list IAM Entities. Entities are used in authorization policies (e.g. user, role, cluster). Supports pagination and OData `filter`, `order_by`, and `select`.

## Example Usage

```hcl
# List all entities (default page/limit)
data "nutanix_iam_entities_v2" "all" {}

# List entities with filter and pagination
data "nutanix_iam_entities_v2" "filtered" {
  filter = "name eq 'user'"
  limit  = 20
  page   = 0
}

# List with order_by
data "nutanix_iam_entities_v2" "ordered" {
  order_by = "name asc"
  limit    = 50
}

# List with select to specify returned fields
data "nutanix_iam_entities_v2" "selected" {
  select = "name,displayName,extId"
  limit  = 10
}
```

## Argument Reference

* `page` - (Optional) Page number of the result set (0-based). Must be between 0 and the maximum number of pages.
* `limit` - (Optional) Number of records to return. Must be between 1 and 100. Default is 50.

* `filter` - (Optional) OData filter expression. The filter can be applied to the following fields:
  * `createdBy` - Filter by creator (user or service ext_id).
  * `createdTime` - Filter by creation time (ISO 8601 format).
  * `displayName` - Filter by display name.
  * `extId` - Filter by entity external identifier.
  * `lastUpdatedTime` - Filter by last updated time (ISO 8601 format).
  * `name` - Filter by entity name.

  **Filter examples:**

  ```hcl
  filter = "createdBy eq '390b7801-7a80-5c94-8a07-8de63651b27b'"
  filter = "createdTime eq '2009-09-23T14:30:00-07:00'"
  filter = "displayName eq 'Role'"
  filter = "extId eq '1e1a9608-79f0-415b-88ad-69d23e325de9'"
  filter = "lastUpdatedTime eq '2009-09-23T14:30:00-07:00'"
  filter = "name eq 'role'"
  ```

* `order_by` - (Optional) OData orderby expression. The orderby can be applied to the following fields:

  * `createdTime` - Sort by creation time (ISO 8601 format).
  * `displayName` - Sort by display name.
  * `extId` - Sort by entity external identifier.
  * `lastUpdatedTime` - Sort by last updated time (ISO 8601 format).
  * `name` - Sort by entity name.

  **Orderby examples:**
  ```hcl
  order_by = "name asc"
  order_by = "createdTime desc"
  ```

* `select` - (Optional) OData select expression. The select can be applied to the following fields:
  * `clientName` - Select by client name.
  * `createdBy` - Select by creator (user or service ext_id).
  * `createdTime` - Select by creation time (ISO 8601 format).
  * `description` - Select by description.
  * `displayName` - Select by display name.
  * `extId` - Select by entity external identifier.
  * `isLogicalAndSupportedForAttributes` - Select by whether logical AND is supported for attributes.
  * `lastUpdatedTime` - Select by last updated time (ISO 8601 format).
  * `name` - Select by entity name.

  **Select examples:**
  ```hcl
  select = "name,displayName,extId"
  ```

---

## Attributes Reference

* `entities` - List of IAM entities.


### Entities
The `entities` attribute supports the following:

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

See [Nutanix List Entities v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Entities/operation/listEntities) for API details.
