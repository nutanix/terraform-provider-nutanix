---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role_v2"
sidebar_current: "docs-nutanix-datasource-roles-v2"
description: |-
  Fetches a role based on the provided external identifier.
---

# nutanix_role_v2

Fetches a role based on the provided external identifier.

## Example Usage

```hcl
data "nutanix_role_v2" "role"{
  ext_id = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
}

```

## Argument Reference

The following arguments are supported:

- `ext_id`: - (Required) ExtId for the Role.

## Attributes Reference

The following attributes are exported:

- `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `display_name`: - The display name for the Role.
- `description`: - Description of the Role.
- `client_name`: - Client that created the entity.
- `operations`: - List of operations for the role.
- `accessible_clients`: - List of Accessible Clients for the Role.
- `accessible_entity_types`: - List of Accessible Entity Types for the Role.
- `assigned_users_count`: - Number of Users assigned to given Role.
- `assigned_users_groups_count`: - Number of User Groups assigned to given Role.
- `created_time`: - The creation time of the Role.
- `last_updated_time`: - The time when the Role was last updated.
- `created_by`: - User or Service Name that created the Role.
- `is_system_defined`: - Flag identifying if the Role is system defined or not.

#### Links

The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix Get Role v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Roles/operation/getRoleById).
