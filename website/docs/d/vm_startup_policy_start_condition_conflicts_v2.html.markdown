---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_startup_policy_start_condition_conflicts_v2"
sidebar_current: "docs-nutanix-datasource-vm-startup-policy-start-condition-conflicts-v2"
description: |-
  List Start condition Conflicts of the provided VM startup policy external identifier.
---

# nutanix_vm_startup_policy_start_condition_conflicts_v2

List Start condition Conflicts of the provided VM startup policy external identifier.

## Example Usage

```hcl
data "nutanix_vm_startup_policy_start_condition_conflicts_v2" "example" {
  vm_startup_policy_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `vm_startup_policy_ext_id` - (Required) The external ID of the VM startup policy.

## Attributes Reference

* `start_condition_conflicts` - List of start condition conflicts. Each entry has the same attributes as the `nutanix_vm_startup_policy_start_condition_conflict_v2` datasource, except `dependee_vms` and `dependent_vms` which are only populated when using the single-item datasource.

See detailed information in [Nutanix VM Startup Policy V2](https://developers.nutanix.com/).
