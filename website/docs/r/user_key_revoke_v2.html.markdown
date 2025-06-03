---
layout: "nutanix"
page_title: "NUTANIX: nutanix_user_key_revoke_v2"
sidebar_current: "docs-nutanix-resource-user-key-revoke-v2"
description: |-
  Revoke the requested key for a user
---

# nutanix_user_key_revoke_v2

Provides Nutanix resource to Revoke the requested key for a user.

## Example Usage

``` hcl
# Create Service Account
resource "nutanix_users_v2" "service_account" {
      username = "service_account_terraform_example 2"
      description = "service account tf"
      email_id = "terraform_plugin@domain.com"
      user_type = "SERVICE_ACCOUNT"
}

# Create key
resource "nutanix_user_key_v2" "create_key" {
  user_ext_id = nutanix_users_v2.service_account.ext_id
  name = "tf_user_api_key"
  key_type = "API_KEY"
  expiry_time = "2026-01-01T00:00:00Z"
  assigned_to = "user1"
}

# revoke key
resource "nutanix_user_key_revoke_v2" "revoke-key"{
  user_ext_id = nutanix_users_v2.service_account.ext_id
  ext_id = nutanix_user_key_v2.create_key.ext_id
}
```

##  Argument Reference

The following arguments are supported:

* `user_ext_id`: - ( Required ) External Identifier of the User.
* `ext_id`: - ( Required ) External identifier of the key.


## Attributes Reference
* `message`: - The message string.
* `severity`: - The message severity.
* `code`: - The code associated with this message.This string is typically prefixed by the namespace the endpoint belongs to. For example: VMM-40000.
* `locale`: - Locale for this message. The default locale would be 'en-US'.
* `error_group`: - The error group associated with this message of severity ERROR.
* `arguments_map`: - The map of argument name to value.

See detailed information in [Nutanix Revoke the requested key V4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/Users/operation/revokeUserKey)
