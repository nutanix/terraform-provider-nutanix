---
layout: "nutanix"
page_title: "NUTANIX: nutanix_sda_policies_v2"
sidebar_current: "docs-nutanix-datasource-sda-policies-v2"
description: |-
  Retrieves a list of all System-Defined Alert Policies.
---

# nutanix_sda_policies_v2

Retrieves a list of all System-Defined Alert Policies.

## Example

```hcl
data "nutanix_sda_policies_v2" "example" {}
```

## Argument Reference

No arguments are required.

## Attribute Reference

The following attributes are exported:

- `sda_policies`: List of System-Defined Alert Policies. Each entry has the same attributes as `nutanix_sda_policy_v2`.

See `nutanix_sda_policy_v2` datasource documentation for the full attribute reference of each policy.
