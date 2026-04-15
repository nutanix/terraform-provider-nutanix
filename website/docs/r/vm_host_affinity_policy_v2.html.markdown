---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_host_affinity_policy_v2"
sidebar_current: "docs-nutanix-resource-vm-host-affinity-policy-v2"
description: |-
  Provides a Nutanix VM-Host Affinity Policy resource to create and manage host affinity policies for virtual machines.
---

# nutanix_vm_host_affinity_policy_v2

Provides a resource to create, read, update, and delete VM-Host Affinity policies. VM-Host Affinity policies ensure that VMs in specific categories are placed on hosts in specified categories. This enables better control over VM placement for compliance, performance, or licensing requirements. For more information on VM-Host Affinity Policies, see the [AHV Administration Guide](https://portal.nutanix.com/page/documents/details?targetId=AHV-Admin-Guide-v11_0:ahv-affinity-policies.html).

## Example Usage

```hcl
# Create VM categories
resource "nutanix_category_v2" "vm_affinity_category" {
  key   = "vm-host-affinity"
  value = "vm-affinity-group-1"
}

# Create host categories
resource "nutanix_category_v2" "host_affinity_category" {
  key   = "vm-host-affinity"
  value = "host-affinity-group-1"
}

# Create VM-Host Affinity policy
resource "nutanix_vm_host_affinity_policy_v2" "example" {
  name              = "vm-host-affinity-policy"
  description       = "Policy to place VMs on specific hosts"
  vm_categories     = [nutanix_category_v2.vm_affinity_category.id]
  host_categories   = [nutanix_category_v2.host_affinity_category.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VM-Host Affinity policy.
* `description` - (Optional) A description of the VM-Host Affinity policy.
* `vm_categories` - (Required) List of VM category external IDs that this policy applies to. VMs with these categories will be subject to the affinity placement rules.
* `host_categories` - (Required) List of host category external IDs that define where the VMs can be placed. Hosts with these categories will be used for VM placement.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id` - The external identifier of the policy.
* `create_time` - The timestamp when the policy was created.
* `update_time` - The timestamp when the policy was last updated.
* `created_by` - Information about the entity that created the policy.
* `last_updated_by` - Information about the entity that last updated the policy.
* `num_vms` - Number of VMs associated with the VM-host affinity policy.
* `num_hosts` - Number of hosts associated with the VM-host affinity policy.
* `num_compliant_vms` - Number of VMs which are compliant with the VM-host affinity policy.
* `num_non_compliant_vms` - Number of VMs which are not compliant with the VM-host affinity policy.

## Import

VM-Host Affinity policies can be imported using the `ext_id`. You can fetch the external ID using the datasource `nutanix_vm_host_affinity_policies_v2`.

```hcl
terraform import nutanix_vm_host_affinity_policy_v2.example <ext_id>
```

## Example Import

```hcl
# First, get the ext_id of the policy
data "nutanix_vm_host_affinity_policies_v2" "policies" {}

# Create the configuration
resource "nutanix_vm_host_affinity_policy_v2" "imported" {}

# Then import using the ext_id
# terraform import nutanix_vm_host_affinity_policy_v2.imported <ext_id>

```

See detailed information in [Nutanix VM-Host Affinity Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmHostAffinityPolicies)
