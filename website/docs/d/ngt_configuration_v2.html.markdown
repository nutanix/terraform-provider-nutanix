---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ngt_configuration_v2"
sidebar_current: "docs-nutanix-datasource-ngt-configuration-v2"
description: |-
  Retrieves the Nutanix Guest Tools configuration for a Virtual Machine.



---

# nutanix_ngt_configuration_v2

Provides Nutanix datasource to Retrieves the Nutanix Guest Tools configuration for a Virtual Machine.


## Example

```hcl
data "nutanix_ngt_configuration_v2" "example" {
  ext_id  = "f29535e2-6bd8-4782-b879-409f17217b31"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) uuid of the Virtual Machine.


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





See detailed information in [Nutanix Get VM NGT configuration V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/getGuestToolsById).
