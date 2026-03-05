---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_anti_affinity_policy_v2"
sidebar_current: "docs-nutanix-datasource-vm-anti-affinity-policy-v2"
description: |-
  Describes a VM-VM Anti-Affinity policy
---

# nutanix_vm_anti_affinity_policy_v2

Retrieve the VM-VM Anti-Affinity policy details for the provided external identifier (ext_id).

## Example

```hcl
data "nutanix_vm_anti_affinity_policy_v2" "policy" {
  ext_id = "12345678-1234-1234-1234-123456789012"
}

output "policy_name" {
  value = data.nutanix_vm_anti_affinity_policy_v2.policy.name
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) The external identifier of a VM-VM Anti-Affinity policy.

## Attribute Reference

The following attributes are exported:

* `name` - The name of the VM-VM Anti-Affinity policy.
* `description` - A description of the VM-VM Anti-Affinity policy.
* `categories` - List of VM category external IDs that this policy applies to.
* `create_time` - The timestamp when the policy was created.
* `update_time` - The timestamp when the policy was last updated.
* `created_by` - Information about the entity that created the policy.
* `updated_by` - Information about the entity that last updated the policy.


See detailed information in [Nutanix Get VM Anti-Affinity Policy V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmAntiAffinityPolicies/operation/getVmAntiAffinityPolicyById)
