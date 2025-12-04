---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_shutdown_action_v2"
sidebar_current: "docs-nutanix-resource-vm-shutdown-action-v2"
description: |-
  Perform shutdown actions on VM.
---

# nutanix_vm_shutdown_action_v2

Collaborative reboot or shutdown of a Virtual Machine through the ACPI support in the operating system. Also, Collaborative reboot or shutdown of a Virtual Machine, requesting Nutanix Guest Tools to trigger a reboot or shutdown from within the VM.

## Example

```hcl
resource "nutanix_vm_shutdown_action_v2" "vmShuts"{
    ext_id= {{ vm uuid }}
    action = "shutdown"
}

resource "nutanix_vm_shutdown_action_v2" "vmShuts"{
    ext_id= {{ vm uuid }}
    action = "reboot"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) The globally unique identifier of a VM. It should be of type UUID.
* `action`: (Required) It supports "shutdown", "guest_shutdown", "reboot", "guest_reboot".
* `guest_power_state_transition_config`: (Optional) Additional configuration for Nutanix Gust Tools power state transition. It should be only used with `guest_shutdown` or `guest_reboot`.

### guest_power_state_transition_config
* `should_enable_script_exec`: (Optional) Indicates whether to run the set script before the VM shutdowns/restarts.
* `should_fail_on_script_failure`: (Optional) Indicates whether to abort VM shutdown/restart if the script fails.


See detailed information in [Nutanix VMs Power Action Shutdown V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/shutdownVm).
See detailed information in [Nutanix VMs Power Action Shutdown Guest Vm V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/shutdownGuestVm).
See detailed information in [Nutanix VMs Power Action Reboot V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/rebootVm).
See detailed information in [Nutanix VMs Power Action Reboot Guest Vm V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/rebootGuestVm).
