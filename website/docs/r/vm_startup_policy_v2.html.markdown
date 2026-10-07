---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_startup_policy_v2"
sidebar_current: "docs-nutanix-resource-vm-startup-policy-v2"
description: |-
  Creates a VM startup policy.
---

# nutanix_vm_startup_policy_v2

Creates a VM startup policy. VM startup policies define the order in which VMs should be started during an HA event or Cluster restart event.

## Example Usage

```hcl
resource "nutanix_vm_startup_policy_v2" "example" {
  name           = "example-vm-startup-policy"
  description    = "Example VM startup policy"
  project_ext_id = "<project-uuid>"

  groups {
    categories {
      ext_id = nutanix_category_v2.cat1.id
    }
  }

  groups {
    categories {
      ext_id = nutanix_category_v2.cat2.id
    }
  }

  start_conditions {
    delay_duration_secs = 30
    power_state_criteria {
      power_on {}
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the VM startup policy.
* `description` - (Optional) Description of the VM startup policy.
* `groups` - (Required) Ordered list of groups configured for the VM startup policy. Minimum 2, maximum 6 groups. Each group is represented by one or more Categories which VMs are expected to be associated with. The list should be ordered in the order in which VMs should be started in an HA event or Cluster restart event.
  * `categories` - (Optional) Categories configured for the group. Minimum 1, maximum 5 categories per group.
    * `ext_id` - (Optional) The external ID (UUID) of the category.
* `project_ext_id` - (Optional) The external ID (UUID) of the project. Once set, this cannot be updated.
* `start_conditions` - (Required) Ordered list of start conditions for the VM startup policy. Minimum 1, maximum 5. The number of start conditions must be exactly one less than the number of groups.
  * `delay_duration_secs` - (Optional) The delay in seconds after the power state criteria is met before the dependent VMs are started. Valid range: 0–600. Default: 0.
  * `power_state_criteria` - (Optional) The power state criteria that the VM must attain before the dependent VMs are started. Exactly one of `power_on` or `guest_bootup` must be specified.
    * `power_on` - (Optional) Power on criteria. Empty block.
    * `guest_bootup` - (Optional) Guest bootup criteria.
      * `timeout_duration_secs` - (Optional) The timeout in seconds in which the VM's Guest OS boot up should be detected successfully.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
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

See detailed information in [Nutanix VM Startup Policy V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/VmStartupPolicies/operation/createVmStartupPolicy).
