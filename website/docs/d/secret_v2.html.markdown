---
layout: "nutanix"
page_title: "NUTANIX: nutanix_secret_v2"
sidebar_current: "docs-nutanix-datasource-secret-v2"
description: |-
  Returns the secret value of a claim token identified by its external ID.
---

# nutanix_secret_v2

Returns the secret value of a claim token identified by its external ID.

## Example

```hcl
data "nutanix_secret_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id`: (Required) External ID of the claim token.

## Attributes Reference

* `secret`: Claim token secret data. This value is marked sensitive.

See detailed information in [Nutanix Get Secret By Claim Token V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
