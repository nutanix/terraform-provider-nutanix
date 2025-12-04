---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ngt_insert_iso_v2"
sidebar_current: "docs-nutanix-resource-ngt-insert_iso-v2"
description: |-
  Installs Nutanix Guest Tools in a Virtual Machine by using the provided credentials.

---

# nutanix_ngt_insert_iso_v2

Provides Nutanix resource toInserts the Nutanix Guest Tools installation and configuration ISO into a virtual machine.


## Example

```hcl
##############################################
# ------------------------------------------------
# This resource allows inserting a NGT ISO into
# a VM’s CD-ROM device.
#
# You can manage both:
# 1. **Insertion** — via `apply`
# 2. **Ejection** — automatically on `delete`
#  You can also eject the NGT ISO by setting `action = "eject"` → triggers eject operation explicitly.
##############################################
resource "nutanix_ngt_insert_iso_v2" "example"{
    ext_id = "ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
    capablities = ["VSS_SNAPSHOT"]
    is_config_only = false
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) uuid of the Virtual Machine.
* `capablities`:(Optional) The list of the application names that are enabled on the guest VM. [`SELF_SERVICE_RESTORE`, `VSS_SNAPSHOT`]
* `is_config_only`:(Optional) Indicates that the Nutanix Guest Tools are already installed on the guest VM, and the ISO is being inserted to update the configuration of these tools.
* `action`: (Optional) Default value: "insert". Accepted values: "insert" → Mounts the specified ISO image to the VM’s CD-ROM, "eject" → Unmounts (ejects) the ISO image from the VM’s CD-ROM.


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





See detailed information in [Nutanix Insert VM Guest Tools V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/insertVmGuestTools).
