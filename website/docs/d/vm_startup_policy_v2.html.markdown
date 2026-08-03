---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_startup_policy_v2"
sidebar_current: "docs-nutanix-datasource-vm-startup-policy-v2"
description: |-
  Fetches the VM startup policy of the provided VM startup policy external identifier.
---

# nutanix_vm_startup_policy_v2

Fetches the VM startup policy of the provided VM startup policy external identifier.

## Example Usage

```hcl
data "nutanix_vm_startup_policy_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id` - (Required) The external ID of the VM startup policy.

## Attributes Reference

* `name` - Name of the VM startup policy.
* `description` - Description of the VM startup policy.
* `groups` - Ordered list of groups configured for the VM startup policy.
  * `categories` - Categories configured for the group.
    * `ext_id` - The globally unique identifier of an instance of type UUID.
* `start_conditions` - Ordered list of start conditions for the VM startup policy.
  * `delay_duration_secs` - The delay in seconds after the power state criteria is met before the dependent VMs are started.
  * `power_state_criteria` - The power state criteria that the VM must attain before the dependent VMs are started.
    * `power_on` - Power on criteria.
    * `guest_bootup` - Guest bootup criteria.
      * `timeout_duration_secs` - The timeout in seconds in which the VM's Guest OS boot up should be detected successfully.
* `project_ext_id` - The external ID (UUID) of the project.
* `create_time` - VM startup policy creation time.
* `update_time` - VM startup policy last updated time.
* `created_by` - The user who created the policy.
  * `ext_id` - The external ID (UUID) of the user.
* `updated_by` - The user who last updated the policy.
  * `ext_id` - The external ID (UUID) of the user.
* `num_compliant_vms` - Number of compliant VMs in the VM startup policy.
* `num_non_compliant_vms` - Number of non-compliant VMs in the VM startup policy.
* `num_pending_vms` - Number of pending VMs in the VM startup policy.
* `num_dependency_conflicts` - Number of dependency conflicts of the VM startup policy.
* `num_start_condition_conflicts` - Number of start condition conflicts of the VM startup policy.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.
  * `href` - The URL at which the entity described by the link can be accessed.
  * `rel` - A name that identifies the relationship of the link to the object that is returned by the URL.

See detailed information in [Nutanix VM Startup Policy V2](https://developers.nutanix.com/).
