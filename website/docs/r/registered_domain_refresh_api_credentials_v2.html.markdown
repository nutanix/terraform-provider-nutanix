---
layout: "nutanix"
page_title: "NUTANIX: nutanix_registered_domain_refresh_api_credentials_v2"
sidebar_current: "docs-nutanix-resource-registered-domain-refresh-api-credentials-v2"
description: |-
  Refreshes the API credentials for the registered domain.
---

# nutanix_registered_domain_refresh_api_credentials_v2

Refreshes the API credentials for the registered domain (refresh-api-credentials action).

## Example Usage

```hcl
resource "nutanix_registered_domain_refresh_api_credentials_v2" "refresh" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`:- (Required) External identifier of the registered domain.

See detailed information in [Nutanix Registered Domain V2](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.4.b1#tag/RegisteredDomains/operation/refreshApiCredentials).
