---
layout: "nutanix"
page_title: "NUTANIX: nutanix_registered_domain_establish_trust_v2"
sidebar_current: "docs-nutanix-resource-registered-domain-establish-trust-v2"
description: |-
  Verifies the registered domain by exchanging the client signed request with Nutanix Central.
---

# nutanix_registered_domain_establish_trust_v2

Verifies the registered domain by exchanging the client signed request with Nutanix Central (establish-trust action).

## Example Usage

```hcl
resource "nutanix_registered_domain_establish_trust_v2" "trust" {
  ext_id                 = "00000000-0000-0000-0000-000000000000"
  client_signing_request = "client-signed-request-value"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`:- (Required) External identifier of the registered domain.
* `client_signing_request`:- (Required, Sensitive) The request value signed by the client for domain registration.

See detailed information in [Nutanix Registered Domain V2](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.0).
