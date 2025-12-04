---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_revert_v2"
sidebar_current: "docs-nutanix-resource-vm-revert-v2"
description: |-
  This operation Revert VM identified by {extId}.
---

# nutanix_vm_revert_v2

This operation Revert VM identified by {extId}. This does an in-place VM restore from a specified VM Recovery Point.

## Example Usage

```hcl
# revert Vm
resource "nutanix_vm_revert_v2" "example"{
  ext_id = "8a938cc5-282b-48c4-81be-de22de145d07"
  vm_recovery_point_ext_id = "c2c249b0-98a0-43fa-9ff6-dcde578d3936"
}

```


## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The globally unique identifier of a VM. It should be of type UUID.
* `vm_recovery_point_ext_id`: -(Required) The external identifier of the VM Recovery Point.


## Attribute Reference

The following attributes are exported:

* `ext_id`: - The globally unique identifier of a VM. It should be of type UUID.
* `vm_recovery_point_ext_id`: -The external identifier of the VM Recovery Point.
* `status`: - The status of the Revert operation.


See detailed information in [Nutanix Revert VM V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/revertVm).
