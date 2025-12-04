---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_key_v2"
sidebar_current: "docs-nutanix-datasource-user-key-v2"
description: |-
  Fetches the requested key through the provided external identifier for the user and the key.
---

# nutanix_user_key_v2
Fetches the requested key through the provided external identifier for the user and the key.

## Example Usage

```hcl
# Get key
data "nutanix_user_key_v2" "get_key"{
  user_ext_id = "<SERVICE_ACCOUNT_UUID>"
  ext_id = "<USER_KEY_UUID>"
}
```

##  Argument Reference
The following arguments are supported:

* `user_ext_id`: - ( Required ) External Identifier of the User.
* `ext_id`: - ( Required ) External identifier of the key.


## Attributes Reference

The following attributes are exported:

* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id` - The External Identifier of the User Group.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`: - Identifier for the key in the form of a name.
* `description`: - Brief description of the key.
* `key_type`: - The type of key.
* `created_time`: - The creation time of the key.
* `last_updated_by`: - User who updated the key.
* `creation_type`: - The creation mechanism of this entity.
* `expiry_time`: - The time when the key will expire.
* `status`: - The status of the key.
* `created_by`: - User or service who created the key.
* `last_updated_time`: - The time when the key was updated.
* `assigned_to`: - External client to whom the given key is allocated.
* `last_used_time`: - The time when the key was last used.
* `key_details`: - Details specific to type of the key.



See detailed information in [Nutanix Get the Requested User Key V4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Users/operation/getUserKeyById)