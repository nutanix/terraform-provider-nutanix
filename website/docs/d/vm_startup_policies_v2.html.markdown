---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_startup_policies_v2"
sidebar_current: "docs-nutanix-datasource-vm-startup-policies-v2"
description: |-
  List VM startup policies.
---

# nutanix_vm_startup_policies_v2

List VM startup policies.

## Example Usage

```hcl
data "nutanix_vm_startup_policies_v2" "example" {}

data "nutanix_vm_startup_policies_v2" "filtered" {
  limit  = 10
  filter = "name eq 'my-policy'"
}
```

## Argument Reference

* `page` - (Optional) Page number for pagination.
* `limit` - (Optional) Number of results per page.
* `filter` - (Optional) OData filter expression.
* `order_by` - (Optional) OData order by expression.

## Attributes Reference

* `policies` - List of VM startup policies. Each entry has the same attributes as the `nutanix_vm_startup_policy_v2` datasource, including `project_ext_id`.

See detailed information in [Nutanix VM Startup Policies V2](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.3#tag/VmStartupPolicies/operation/listVmStartupPolicies).
