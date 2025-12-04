---
layout: "nutanix"
page_title: "NUTANIX: ngt_installation_v2"
sidebar_current: "docs-nutanix-resource-ngt-installation-v2"
description: |-
  Installs Nutanix Guest Tools in a Virtual Machine by using the provided credentials.

---

# nutanix_ngt_installation_v2

Provides Nutanix resource to Installs Nutanix Guest Tools in a Virtual Machine by using the provided credentials.


## Example

```hcl
resource "nutanix_ngt_installation_v2" "example"{
    ext_id = "ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
    credential {
        username = "username"
        password = "pass.1234567890"
    }
    reboot_preference {
        schedule_type = "IMMEDIATE"
    }
    capablities = ["VSS_SNAPSHOT"]
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) uuid of the Virtual Machine.
* `credential`:(Optional) Sign in credentials for the server.
* `capablities`:(Optional) The list of the application names that are enabled on the guest VM. [`SELF_SERVICE_RESTORE`, `VSS_SNAPSHOT`]
* `reboot_preference`:(Optional) The restart schedule after installing or upgrading Nutanix Guest Tools.
* `is_enabled`:(Optional) Indicates whether Nutanix Guest Tools is enabled or not.


### Credential

The credential attribute supports the following:

* `username`: - (Required) username to sign in to server
* `password`: - (Required) password to sign in to server

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





See detailed information in [Nutanix Install VM Guest Tools V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/installVmGuestTools).
