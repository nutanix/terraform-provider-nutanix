---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_startup_policy_start_condition_conflict_dependee_vms_v2"
sidebar_current: "docs-nutanix-datasource-vm-startup-policy-start-condition-conflict-dependee-vms-v2"
description: |-
  List dependee VMs of a start condition conflict for the provided start condition conflict external identifier and VM startup policy external identifier.
---

# nutanix_vm_startup_policy_start_condition_conflict_dependee_vms_v2

List dependee VMs of a start condition conflict for the provided start condition conflict external identifier and VM startup policy external identifier.

## Example Usage

```hcl
data "nutanix_vm_startup_policy_start_condition_conflict_dependee_vms_v2" "example" {
  vm_startup_policy_ext_id       = "00000000-0000-0000-0000-000000000000"
  start_condition_conflict_ext_id = "00000000-0000-0000-0000-000000000001"
}
```

## Argument Reference

* `vm_startup_policy_ext_id` - (Required) The external ID of the VM startup policy.
* `start_condition_conflict_ext_id` - (Required) The external ID of the start condition conflict of a VM startup policy.

## Attributes Reference

* `vms` - List of VM references.
  * `ext_id` - The external ID (UUID) of the VM.

See detailed information in [Nutanix VM Startup Policy V2](https://developers.nutanix.com/).
