---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_startup_policy_start_condition_conflict_v2"
sidebar_current: "docs-nutanix-datasource-vm-startup-policy-start-condition-conflict-v2"
description: |-
  Get Start condition conflict for the provided start condition conflict external identifier and VM startup policy external identifier.
---

# nutanix_vm_startup_policy_start_condition_conflict_v2

Get Start condition conflict for the provided start condition conflict external identifier and VM startup policy external identifier.

## Example Usage

```hcl
data "nutanix_vm_startup_policy_start_condition_conflict_v2" "example" {
  vm_startup_policy_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id                   = "00000000-0000-0000-0000-000000000001"
}
```

## Argument Reference

* `vm_startup_policy_ext_id` - (Required) The external ID of the VM startup policy.
* `ext_id` - (Required) The external ID of the start condition conflict of a VM startup policy.

## Attributes Reference

* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `conflicting_policy` - The conflicting VM startup policy reference.
  * `ext_id` - The external ID (UUID) of the conflicting VM startup policy.
* `conflicting_start_condition` - The conflicting start condition.
  * `delay_duration_secs` - The delay in seconds after the power state criteria is met before the dependent VMs are started.
  * `power_state_criteria` - The power state criteria.
    * `power_on` - Power on criteria.
    * `guest_bootup` - Guest bootup criteria.
      * `timeout_duration_secs` - The timeout in seconds.
* `start_condition` - The current start condition.
  * `delay_duration_secs` - The delay in seconds.
  * `power_state_criteria` - The power state criteria.
    * `power_on` - Power on criteria.
    * `guest_bootup` - Guest bootup criteria.
      * `timeout_duration_secs` - The timeout in seconds.
* `dependee_category` - The dependee category reference.
  * `ext_id` - The globally unique identifier of an instance of type UUID.
* `dependent_category` - The dependent category reference.
  * `ext_id` - The globally unique identifier of an instance of type UUID.
* `dependee_vms_associated_categories` - The categories through which the dependee VMs are associated with the policies.
  * `ext_id` - The globally unique identifier of an instance of type UUID.
* `dependent_vms_associated_categories` - The categories through which the dependent VMs are associated with the policies.
  * `ext_id` - The globally unique identifier of an instance of type UUID.
* `dependee_vms` - List of dependee VMs involved in the start condition conflict. This is fetched automatically and avoids the need for a separate `nutanix_vm_startup_policy_start_condition_conflict_dependee_vms_v2` datasource call.
  * `ext_id` - The external ID (UUID) of the VM.
* `dependent_vms` - List of dependent VMs involved in the start condition conflict. This is fetched automatically and avoids the need for a separate `nutanix_vm_startup_policy_start_condition_conflict_dependent_vms_v2` datasource call.
  * `ext_id` - The external ID (UUID) of the VM.
* `links` - A HATEOAS style link for the response.
  * `href` - The URL at which the entity described by the link can be accessed.
  * `rel` - A name that identifies the relationship of the link to the object that is returned by the URL.

See detailed information in [Nutanix Get VM Startup Policy Start Condition Conflict V2](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/VmStartupPolicies/operation/getVmStartupPolicyStartConditionConflictById).
