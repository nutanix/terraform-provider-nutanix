---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ngt_upgrade_v2"
sidebar_current: "docs-nutanix-resource-ngt-upgrade-v2"
description: |-
  Installs Nutanix Guest Tools in a Virtual Machine by using the provided credentials.

---

# nutanix_ngt_upgrade_v2

Provides Nutanix resource to Trigger an in-guest upgrade of Nutanix Guest Tools.


## Example

```hcl
resource "nutanix_ngt_upgrade_v2" "example"{
    ext_id = "ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
    reboot_preference {
        schedule_type = "IMMEDIATE"
    }
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) uuid of the Virtual Machine.
* `reboot_preference`:(Optional) The restart schedule after installing or upgrading Nutanix Guest Tools.


### Reboot Preference

The reboot_preference attribute supports the following:

* `schedule_type`: - Schedule type for restart.
    * `LATER` : Schedule a restart for a specific time.
    * `SKIP` : Do not schedule a restart.
    * `IMMEDIATE` : Schedule an immediate restart.
* `schedule`: - Restart schedule.

#### schedule

The schedule attribute supports the following:

* `start_time`: - The start time for a scheduled restart.

## Attribute Reference

The following attributes are exported:
* `ext_id`: uuid of the Virtual Machine.
* `capablities`: The list of the application names that are enabled on the guest VM. [`SELF_SERVICE_RESTORE`, `VSS_SNAPSHOT`]
* `is_enabled`: The entities being qualified by the Authorization Policy.
* `version`: Version of Nutanix Guest Tools installed on the VM.
* `is_installed`: Indicates whether Nutanix Guest Tools is installed on the VM or not.
* `is_enabled`: Indicates whether Nutanix Guest Tools is enabled or not.
* `is_iso_inserted`: Indicates whether Nutanix Guest Tools ISO is inserted or not.
* `available_version`: Version of Nutanix Guest Tools available on the cluster.
* `guest_os_version`: Version of the operating system on the VM.
* `is_reachable`: Indicates whether the communication from VM to CVM is active or not.
* `is_vss_snapshot_capable`: Indicates whether the VM is configured to take VSS snapshots through NGT or not.
* `is_vm_mobility_drivers_installed`: Indicates whether the VM mobility drivers are installed on the VM or not.





See detailed information in [Nutanix Upgrade VM Guest Tools V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/upgradeVmGuestTools).
