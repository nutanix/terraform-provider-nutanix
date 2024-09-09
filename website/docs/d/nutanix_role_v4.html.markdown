---
layout: "nutanix"
page_title: "NUTANIX: nutanix_roles_v2"
sidebar_current: "docs-nutanix-datasource-roles-v4"
description: |-
  Describes a List Role(s).
---

# nutanix_volume_groups_v4

Describes a List all the Role(s).

## Example Usage

```hcl
data "nutanix_roles_v2" "roles"{
  ext_id = var.role_ext_id
}

```

##  Argument Reference

The following arguments are supported:

* `ext_id`: - (Required) ExtId for the Role.

## Attributes Reference
The following attributes are exported:


* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `display_name`: - The display name for the Role.
* `description`: - Description of the Role.
* `client_name`: - Client that created the entity.
* `operations`: - Indicates whether to enable Volume Group load balancing for VM attachments. This cannot be enabled if there are iSCSI client attachments already associated with the Volume Group, and vice-versa. This is an optional field.
* `accessible_clients`: - List of Accessible Clients for the Role.
* `accessible_entity_types`: - List of Accessible Entity Types for the Role.
* `assigned_users_count`: - Number of Users assigned to given Role.
* `assigned_users_groups_count`: - Number of User Groups assigned to given Role.
* `created_time`: - The creation time of the Role.
* `last_updated_time`: - The time when the Role was last updated.
* `created_by`: - User or Service Name that created the Role.
* `is_system_defined`: - Flag identifying if the Role is system defined or not.

#### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix Roles](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0.b1).
