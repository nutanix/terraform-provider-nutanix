---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_startup_policy_dependency_conflicts_v2"
sidebar_current: "docs-nutanix-datasource-vm-startup-policy-dependency-conflicts-v2"
description: |-
  List Dependency Conflicts of the provided VM startup policy external identifier.
---

# nutanix_vm_startup_policy_dependency_conflicts_v2

List Dependency Conflicts of the provided VM startup policy external identifier.

## Example Usage

```hcl
data "nutanix_vm_startup_policy_dependency_conflicts_v2" "example" {
  vm_startup_policy_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `vm_startup_policy_ext_id` - (Required) The external ID of the VM startup policy.

## Attributes Reference

* `dependency_conflicts` - List of dependency conflicts. Each entry has the same attributes as the `nutanix_vm_startup_policy_dependency_conflict_v2` datasource, except `dependee_vms` and `dependent_vms` which are only populated when using the single-item datasource.

See detailed information in [Nutanix List VM Startup Policy Dependency Conflicts V2](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/VmStartupPolicies/operation/listVmStartupPolicyDependencyConflicts).
