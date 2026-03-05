---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_anti_affinity_policy_v2"
sidebar_current: "docs-nutanix-resource-vm-anti-affinity-policy-v2"
description: |-
  Provides a Nutanix VM-VM Anti-Affinity Policy resource to create and manage anti-affinity policies for virtual machines.
---

# nutanix_vm_anti_affinity_policy_v2

Provides a resource to create, read, update, and delete VM-VM Anti-Affinity policies. VM-VM Anti-Affinity policies ensure that VMs in specified categories are spread across different hosts for high availability and fault tolerance. This helps prevent single points of failure by distributing VMs across the cluster. For more information on VM-VM Anti-Affinity policies, see the [AHV Administration Guide](https://portal.nutanix.com/page/documents/details?targetId=AHV-Admin-Guide-v11_0:ahv-anti-affinity-policies-c.html).

## Example Usage

```hcl
# Create VM category
resource "nutanix_category_v2" "anti_affinity_category" {
  key   = "vm-anti-affinity"
  value = "anti-affinity-group-1"
}

# Create VM-VM Anti-Affinity policy
resource "nutanix_vm_anti_affinity_policy_v2" "example" {
  name        = "vm-anti-affinity-policy"
  description = "Policy to spread VMs across different hosts for high availability"
  categories  = [nutanix_category_v2.anti_affinity_category.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VM-VM Anti-Affinity policy.
* `description` - (Optional) A description of the VM-VM Anti-Affinity policy.
* `categories` - (Required) List of VM category external IDs (`ext_id`) that this policy applies to. VMs with these categories will be spread across different hosts according to the anti-affinity rules.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id` - The external identifier of the policy.
* `create_time` - The timestamp when the policy was created.
* `update_time` - The timestamp when the policy was last updated.
* `created_by` - Information about the entity that created the policy.
* `updated_by` - Information about the entity that last updated the policy.

## Import

VM-VM Anti-Affinity policies can be imported using the `ext_id`. You can fetch the external ID using the datasource `nutanix_vm_anti_affinity_policies_v2`.

```hcl
terraform import nutanix_vm_anti_affinity_policy_v2.example <ext_id>
```

## Example Import

```hcl
# First, get the ext_id of the policy
data "nutanix_vm_anti_affinity_policies_v2" "policies" {}

# Then import using the ext_id
# terraform import nutanix_vm_anti_affinity_policy_v2.imported <ext_id>

# Create the configuration
resource "nutanix_vm_anti_affinity_policy_v2" "imported" {}
```

See detailed information in [Nutanix VM-VM Anti-Affinity Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmAntiAffinityPolicies)
