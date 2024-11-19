---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_power_action_v4"
sidebar_current: "docs-nutanix-resource-vm-power-action-v4"
description: |-
  Perform power actions on VM. 
---

# nutanix_vm_power_action_v4

Perform power actions on VM. It supports 'Power Off', 'Power On', 'Power cycle', 'Reset'.


## Example

```hcl

  resource "nutanix_vm_power_action_v4" "test" {
    ext_id= resource.nutanix_virtual_machine_v2.rtest.id
    action = "power_on"
  }

  resource "nutanix_vm_power_action_v4" "test" {
    ext_id= resource.nutanix_virtual_machine_v2.rtest.id
    action = "vm_reset"
  }
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) The globally unique identifier of a VM. It should be of type UUID.
* `action`: (Required) It supports "power_on", "power_off", "power_cycle", "vm_reset" . 



See detailed information in [Nutanix VMs Power Action](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).