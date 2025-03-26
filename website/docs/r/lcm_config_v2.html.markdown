---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_config_v2"
sidebar_current: "docs-nutanix-lcm-config-v2"
description: |-
  Update LCM configuration.
---

# nutanix_lcm_entity_v2
Update LCM configuration.

## Example

```hcl
# Enable Auto Inventory, Add Auto Inventory Schedule and enable auto upgrade
resource "nutanix_lcm_config_v2" "lcm-configuration-update" {
    x_cluster_id = "0005a104-0b0b-4b0b-8005-0b0b0b0b0b0b"
    is_auto_inventory_enabled = true
	auto_inventory_schedule = "16:30"
    has_module_auto_upgrade_enabled = true
}

# Update the LCM url to darksite server
resource "nutanix_lcm_config_v2" "lcm-configuration-update-connectivity-type" {
    x_cluster_id = "0005a104-0b0b-4b0b-8005-0b0b0b0b0b0b"
    url = "https://x.x.x.x:8000/builds"
	connectivity_type = "DARKSITE_WEB_SERVER"
}

```
## Argument Reference
The following arguments are supported:

* `x_cluster_id`: (Optional) Cluster uuid on which the resource is present or operation is being performed.
* `url`: (Optional) URL of the LCM repository.
* `is_auto_inventory_enabled`: (Optional) Indicates if the auto inventory operation is enabled. The default value is set to False.
* `auto_inventory_schedule`: (Optional) The scheduled time in "%H:%M" 24-hour format of the next inventory execution. Used when auto_inventory_enabled is set to True. The default schedule time is 03:00(AM).
* `connectivity_type`: (Optional)This field indicates whether LCM framework on the cluster is running in connected-site mode or darksite mode.
* `is_https_enabled`: (Optional) Indicates if the LCM URL has HTTPS enabled. The default value is True.
* `has_module_auto_upgrade_enabled`: (Optional) Indicates if LCM is enabled to auto-upgrade products. The default value is False.

See detailed information in [Nutanix Update LCM Config V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Config/operation/updateConfig)
