---
layout: "nutanix"
page_title: "NUTANIX: nutanix_roles_v2"
sidebar_current: "docs-nutanix-datasource-roles-v2"
description: |-
  Describes a List Role(s).
---

# nutanix_roles_v2

Describes a List all the Role(s).

## Example Usage

```hcl
# List all Roles
data "nutanix_roles_v2" "roles"{}

# List Roles with filter
data "nutanix_roles_v2" "filtered-roles"{
  filter = "displayName eq 'example_role'"
}

# List Roles with filter and orderby
data "nutanix_roles_v2" "filtered-ordered-roles"{
  filter = "displayName eq 'example_role'"
  order_by = "createdTime desc"
}

```

##  Argument Reference

The following arguments are supported:

* `page`: - A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` :A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - clientName
  - createdBy
  - extId
  - createdTime
  - displayName
  - extId
  - isSystemDefined
  - lastUpdatedTime
* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - createdTime
  - distinguishedName
  - displayName
  - extId-
  - lastUpdatedTime
* `select` : A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. following fields:
  - accessibleClients
  - accessibleEntityTypes
  - assignedUserGroupsCount
  - assignedUsersCount
  - clientName
  - createdBy
  - createdTime
  - description
  - displayName
  - extId
  - isSystemDefined
  - lastUpdatedTime
  - links
  - operations
  - tenantId

## Attributes Reference
The following attributes are exported:

* `roles`: - List of Roles.

### roles
The `roles` attribute contains list of Role objects. Each Role object contains the following attributes:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `display_name`: - The display name for the Role.
* `description`: - Description of the Role.
* `client_name`: - Client that created the entity.
* `operations`: - List of operations for the role.
* `accessible_clients`: - List of Accessible Clients for the Role.
* `accessible_entity_types`: - List of Accessible Entity Types for the Role.
* `assigned_users_count`: - Number of Users assigned to given Role.
* `assigned_users_groups_count`: - Number of User Groups assigned to given Role.
* `created_time`: - The creation time of the Role.
* `last_updated_time`: - The time when the Role was last updated.
* `created_by`: - User or Service Name that created the Role.
* `is_system_defined`: - Flag identifying if the Role is system defined or not.

### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix List Roles v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Roles/operation/listRoles).
