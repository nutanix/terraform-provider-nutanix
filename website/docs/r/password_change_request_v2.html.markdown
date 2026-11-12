---
layout: "nutanix"
page_title: "NUTANIX: nutanix_password_change_request_v2"
sidebar_current: "docs-nutanix-resource-password-change-request-v2"
description: |-
  Initiate change password request for a system user on a supported product.
---

# nutanix_password_change_request_v2


Initiate change password request for a system user on a supported product.

> **WriteOnly credentials:** `current_password`, `new_password` are write-only. They are never stored in Terraform state. Re-supply them in configuration on every apply. Do not use `lifecycle.ignore_changes` on these attributes. See the [WriteOnly Credentials](/docs/providers/nutanix/guides/writeonly_credentials.html) guide for upgrade caveats.

## Example Usage

```hcl
resource "nutanix_password_change_request_v2" "change_admin_aos_password" {
	ext_id = "557329fc-6f28-44e4-8905-12c54d704ff9"
	current_password = "M4zGWn^Haxs0za~"
	new_password = "*yk@1+U0syIr"
}

```

## Argument Reference

The following arguments are supported:

- `ext_id`: -(Required) External identifier of the system user password.
- `current_password`: -(Optional, WriteOnly) Existing password of a user account.
- `new_password`: -(Required, WriteOnly) New password for a user account.


See detailed information in [Nutanix Initiate password update for a system user V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.3#tag/PasswordManager/operation/changeSystemUserPasswordById).
