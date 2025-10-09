---
layout: "nutanix"
page_title: "NUTANIX: nutanix_accept_eula_v2"
sidebar_current: "docs-nutanix-resource-accept-eula-v2"
description: |-
  API for allowing users to accept End User License Agreement.
---

# nutanix_accept_eula_v2

API for allowing users to accept End User License Agreement.

## Example Usage

```hcl
resource "nutanix_accept_eula_v2" "accept_eula" {
  user_name   = "admin"
  job_title   = "Nutanix Administrator"
  login_id    = "12345"
  company_name = "Nutanix"
}
```

## Argument Reference

The following arguments are supported:

* `user_name`: - (Required) User name of the user accepting the EULA.
* `login_id`: - (Required) Login ID of the user accepting the EULA.
* `job_title`: - (Required) Job title of the user accepting the EULA.
* `company_name`: - (Required) Company name of the user accepting the EULA.

## Attributes Reference
The following attributes are exported:

* `message`: - The message string.
* `severity`: - The message severity.
* `code`: - The code associated with this message. This string is typically prefixed with the namespace to which the endpoint belongs. For example: VMM-40000
* `locale`: - Locale for this message. The default locale would be 'en-US'.
* `error_group`: - The error group associated with this message of severity ERROR.
* `arguments_map`: - The map of argument name to value.



See detailed information in [Accept End User License Agreement](https://developers.nutanix.com/api-reference?namespace=licensing&version=v4.1#tag/EndUserLicenseAgreement/operation/addUser).
