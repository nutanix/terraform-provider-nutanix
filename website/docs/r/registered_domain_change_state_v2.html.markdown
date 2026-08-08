---
layout: "nutanix"
page_title: "NUTANIX: nutanix_registered_domain_change_state_v2"
sidebar_current: "docs-nutanix-resource-registered-domain-change-state-v2"
description: |-
  Update state of the registered domain.
---

# nutanix_registered_domain_change_state_v2

Update state of the registered domain (sync-registration action).

## Example Usage

```hcl
resource "nutanix_registered_domain_change_state_v2" "sync" {
  ext_id             = "00000000-0000-0000-0000-000000000000"
  registration_state = "REGISTERED"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`:- (Required) External identifier of the registered domain.
* `registration_state`:- (Required) The new registration state to set. One of `UNREGISTERED`, `REGISTERING`, `REGISTERED`, `UNREGISTERING`, `REGISTRATION_ERROR`, `UNREGISTRATION_ERROR`, `DRAFT`.

See detailed information in [Nutanix Registered Domain V2](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.0).
