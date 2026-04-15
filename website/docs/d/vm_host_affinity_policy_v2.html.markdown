---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_host_affinity_policy_v2"
sidebar_current: "docs-nutanix-datasource-vm-host-affinity-policy-v2"
description: |-
  Describes a VM-Host Affinity policy
---

# nutanix_vm_host_affinity_policy_v2

Retrieve the VM-Host Affinity policy details for the provided external identifier (ext_id).

## Example

```hcl
data "nutanix_vm_host_affinity_policy_v2" "policy" {
  ext_id = "12345678-1234-1234-1234-123456789012"
}

output "policy_name" {
  value = data.nutanix_vm_host_affinity_policy_v2.policy.name
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) The external identifier of a VM-Host Affinity policy.

## Attribute Reference

The following attributes are exported:

* `name` - The name of the VM-Host Affinity policy.
* `description` - A description of the VM-Host Affinity policy.
* `vm_categories` - List of VM category external IDs that this policy applies to.
* `host_categories` - List of host category external IDs that define where the VMs can be placed.
* `create_time` - The timestamp when the policy was created.
* `update_time` - The timestamp when the policy was last updated.
* `created_by` - Information about the entity that created the policy.
* `last_updated_by` - Information about the entity that last updated the policy.
* `num_vms` - Number of VMs associated with the VM-host affinity policy.
* `num_hosts` - Number of hosts associated with the VM-host affinity policy.
* `num_compliant_vms` - Number of VMs which are compliant with the VM-host affinity policy.
* `num_non_compliant_vms` - Number of VMs which are not compliant with the VM-host affinity policy.

See detailed information in [Nutanix Get VM Host Affinity Policy V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmHostAffinityPolicies/operation/getVmHostAffinityPolicyById)
