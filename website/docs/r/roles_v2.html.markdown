---
layout: "nutanix"
page_title: "NUTANIX: nutanix_roles_v2"
sidebar_current: "docs-nutanix-resource-roles-v2"
description: |-
  This operation submits a request to add a Role..
---

# nutanix_roles_v2

Provides a resource to add a Role.

## Example Usage

```hcl
# filtered list operation
data "nutanix_operations_v2" "operations-filtered-list" {
  filter = "startswith(displayName, 'Create_')"
}

# Create role
resource "nutanix_roles_v2" "example-role"{
  display_name = "example_role"
  description  = "create example role"
  operations = [
    data.nutanix_operations_v2.operations-filtered-list.operations[0].ext_id,
    data.nutanix_operations_v2.operations-filtered-list.operations[1].ext_id,
    data.nutanix_operations_v2.operations-filtered-list.operations[2].ext_id,
    data.nutanix_operations_v2.operations-filtered-list.operations[3].ext_id
  ]
}
```

## Argument Reference

The following arguments are supported:

- `display_name`: -(Required) The display name for the Role.
- `description`: - Description of the Role.
- `client_name`: - Client that created the entity.
- `operations`: -(Required) List of operations for the role.

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

### Links

The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

## Import

This helps to manage existing entities which are not created through terraform. Role can be imported using the `UUID`. (ext_id in v4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_roles_v2" "import_role" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_roles_v2" "fetch_roles"{}
terraform import nutanix_roles_v2.import_role <UUID>
```

See detailed information in [Nutanix Create Role ](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Roles/operation/createRole).
