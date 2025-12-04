---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_upgrade_v2"
sidebar_current: "docs-nutanix-lcm-upgrade-v2"
description: |-
  Perform upgrade operation to a specific target version for discovered LCM entity/entities.
---

# nutanix_lcm_upgrade_v2

Perform upgrade operation to a specific target version for discovered LCM entity/entities.


## Example

```hcl

# upgrade the entity
resource "nutanix_lcm_upgrade_v2" "upgrade" {
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
* `auto_handle_flags`: (Optional) List of automated system operations to perform, to avoid precheck failure and let the system restore state after an update is complete. The allowed flag is: - 'powerOffUvms': This allows the system to automatically power off user VMs which cannot be migrated to other hosts and power them on when the update is done. This option can avoid pinned VM precheck failure on the host which needs to enter maintenance mode during the update and allow the update to go through. Items Enum: `POWER_OFF_UVMS`
* `max_wait_time_in_secs`: (Optional) Number of seconds LCM waits for the VMs to come up after exiting host maintenance mode. Value in Range [ 60 .. 86400]

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


See detailed information in [Nutanix LCM Upgrade v4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Upgrades/operation/performUpgrade).

