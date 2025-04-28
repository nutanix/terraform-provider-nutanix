---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_key_v2"
sidebar_current: "docs-nutanix-resource-user-key-v2"
description: |-
  Create key of a requested type for a user.
---


# nutanix_user_key_v2

Provides Nutanix resource to Create key of a requested type for a user.

## Example Usage

``` hcl
# Create Service Account
resource "nutanix_users_v2" "service_account" {
      username = "service_account_terraform_example"
      description = "service account tf"
      email_id = "terraform_plugin@domain.com"
      user_type = "SERVICE_ACCOUNT"
}

# Create key under service account, never expires
resource "nutanix_user_key_v2" "create_key" {
   user_ext_id = nutanix_users_v2.service_account.ext_id
   name = "api_key_developers"
   key_type = "API_KEY"
   expiry_time = "2125-01-01T00:00:00Z"
   assigned_to = "developer_user_1"
}
```
##  Argument Reference

The following arguments are supported:

* `user_ext_id`: - ( Required ) External Identifier of the User.
* `name`: - ( Required ) Identifier for the key in the form of a name.
* `description`: - ( Optional ) Brief description of the key.
* `key_type`: - ( Required ) The type of key. Enum Values:
      * "API_KEY":	A key type that is used to identify a service.
      * "OBJECT_KEY":	A combination of access key and secret key to sign an API request.
* `creation_type`: - ( Optional ) The creation mechanism of this entity. Enum Values:
      * "PREDEFINED":	Predefined creator workflow type is for entity created by the system.
      * "SERVICEDEFINED":	Servicedefined creator workflow type is for entity created by the service.
      * "USERDEFINED":	Userdefined creator workflow type is for entity created by the users.
* `expiry_time`: - ( Optional ) The time when the key will expire.
* `status`: - ( Optional ) The status of the key. Enum Values:
      * "REVOKED":	Key is revoked.
      * "VALID":	Key is valid.
      * "EXPIRED":	Key is expired.
* `assigned_to`: - ( Optional ) External client to whom the given key is allocated.

## Attributes Reference
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


### Links

The links attribute supports the following:

See detailed information in [Nutanix Create User Key V4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Users/operation/createUserKey)