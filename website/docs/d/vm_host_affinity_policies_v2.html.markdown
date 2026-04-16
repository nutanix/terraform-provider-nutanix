---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_host_affinity_policies_v2"
sidebar_current: "docs-nutanix-datasource-vm-host-affinity-policies-v2"
description: |-
  Describes VM-Host Affinity policies
---

# nutanix_vm_host_affinity_policies_v2

List VM-Host Affinity policies details with support for pagination, filtering, and sorting.

## Example

```hcl
# List all VM-Host Affinity policies
data "nutanix_vm_host_affinity_policies_v2" "all_policies" {}

# List with pagination
data "nutanix_vm_host_affinity_policies_v2" "paginated_policies" {
  page  = 0
  limit = 10
}

# List with filtering
data "nutanix_vm_host_affinity_policies_v2" "filtered_policies" {
  filter = "name eq 'my-policy'"
}

# List with ordering
data "nutanix_vm_host_affinity_policies_v2" "sorted_policies" {
  order_by = "name asc"
}

# List with multiple filters
data "nutanix_vm_host_affinity_policies_v2" "complex_policies" {
  filter   = "startswith(name, 'prod-')"
  order_by = "create_time desc"
  page     = 0
  limit    = 20
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions. For example, `$filter=name eq 'my-policy'` would filter the result on policy name 'my-policy', `$filter=startswith(name, 'prod-')` would filter on policy names starting with 'prod-'. The filter can be applied to the following fields:
  - `name`
  - `ext_id`
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using `asc` or descending order using `desc`. If `asc` or `desc` are not specified, the resources will be sorted in ascending order by default. For example, `name asc` would sort policies by name in ascending order, `update_time desc` would sort by update time in descending order. The order_by can be applied to the following fields:
  - `name`
  - `update_time`

## Attribute Reference

* `policies` - List of all VM-Host Affinity policies

### Policies

The `policies` object is a list of VM-Host Affinity policies. Each VM-Host Affinity policy object contains the following attributes:

* `ext_id` - The external identifier of the VM-Host Affinity policy.
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

See detailed information in [Nutanix List VM Host Affinity Policies V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmHostAffinityPolicies/operation/listVmHostAffinityPolicies)
