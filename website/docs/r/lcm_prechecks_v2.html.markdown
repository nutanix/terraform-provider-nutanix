---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_prechecks_v2"
sidebar_current: "docs-nutanix_lcm_prechecks_v2"
description: |-
  Perform LCM prechecks for the intended update operation.
---

# nutanix_lcm_prechecks_v2

Perform LCM prechecks for the intended update operation.

## Example

```hcl
resource "nutanix_lcm_prechecks_v2" "pre-checks" {
  x_cluster_id = "0005a104-0b0b-4b0-8005-0b0b0b0b0b0b"
  entity_update_specs {
    entity_uuid = "0c5c9e53-3551-4c5d-b13c-e41c04cbfaf7"
    to_version  = "4.0.0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `x_cluster_id`: (Optional) Cluster uuid on which the resource is present or operation is being performed.
* `management_server`: (Optional) Cluster management server configuration used while updating clusters with ESX or Hyper-V.
* `entity_update_specs`: (Required) List of entity update objects for getting recommendations.
* `skipped_precheck_flags`: (Optional) List of prechecks to skip. The allowed value is 'powerOffUvms' that skips the pinned VM prechecks. Items Enum: `POWER_OFF_UVMS`

### Management Server
The `management_server` attribute supports the following:

* `hypervisor_type`: (Required) Type of Hypervisor present in the cluster. Enum Values:
    * "HYPERV" : Hyper-V Hypervisor.
    * "ESX" : ESX Hypervisor.
    * "AHV" : Nutanix AHV Hypervisor.
* `ip`: (Required) IP address of the management server.
* `username`: (Required) Username to login to the management server.
* `password`: (Required) Password to login to the management server.

### Entity Update Specs
The `entity_update_specs` attribute supports the following:

* `entity_uuid`: (Required) UUID of the LCM entity.
* `to_version`: (Required) Version to upgrade to.

See detailed information in [Nutanix LCM Prechecks v4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Prechecks/operation/performPrechecks)
