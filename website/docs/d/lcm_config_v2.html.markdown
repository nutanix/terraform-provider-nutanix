---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_config_v2"
sidebar_current: "docs-nutanix-datasource-lcm-config-v2"
description: |-
  Get LCM configuration.
---

# nutanix_lcm_config_v2
Get LCM configuration.

## Example

```hcl
data "nutanix_lcm_config_v2" "lcm-configuration"{}

# Get LCM configuration for a specific cluster
data "nutanix_lcm_config_v2" "lcm-configuration-cluster" {
  x_cluster_id = "0005a104-0b0b-4b0b-8005-0b0b0b0b0b0b"
}
```

## Argument Reference
The following arguments are supported:

* `x_cluster_id`: (Optional) Cluster uuid on which the resource is present or operation is being performed.


## Attributes Reference
The following attributes are exported:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `url`: URL of the LCM repository.
* `is_auto_inventory_enabled`: Indicates if the auto inventory operation is enabled. The default value is set to False.
* `auto_inventory_schedule`: The scheduled time in "%H:%M" 24-hour format of the next inventory execution. Used when auto_inventory_enabled is set to True. The default schedule time is 03:00(AM).
* `version`: LCM version installed on the cluster.
* `display_version`: User friendly display version of LCM installed on the cluster.
* `connectivity_type`: This field indicates whether LCM framework on the cluster is running in connected-site mode or darksite mode. Values are :
  - `CONNECTED_SITE`: In connected-site, LCM on the cluster has internet connectivity to reach configured portal for downloading LCM modules/bundles etc.
  - `DARKSITE_DIRECT_UPLOAD`: LCM on the cluster does not have external connectivity and will have a facility to upload darksite bundles through LCM.
  - `DARKSITE_WEB_SERVER`: LCM on the cluster does not have external connectivity and will have a connection to darksite webserver maintained by the customer.
* `is_https_enabled`: Indicates if the LCM URL has HTTPS enabled. The default value is True.
* `supported_software_entities`: List of entities for which One-Click upgrades are supported.
* `deprecated_software_entities`: List of entities for which One-Click upgrades are not available.
* `is_framework_bundle_uploaded`: Indicates if the bundle is uploaded or not.
* `has_module_auto_upgrade_enabled`: Indicates if LCM is enabled to auto-upgrade products. The default value is False.

See detailed information in [Nutanix Get LCM Config V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Config/operation/getConfig)
