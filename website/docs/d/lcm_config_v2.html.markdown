---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_config_v2"
sidebar_current: "docs-nutanix-datasource-lcm-config-v2"
description: |-
  Get LCM configuration.
---

# nutanix_lcm_entity_v2
Get LCM configuration.

## Example

```hcl
data "nutanix_lcm_config_v2" "lcm-configuration" {}
```

## Argument Reference
The following arguments are supported:

* `x_cluster_id`: (Optional) Cluster uuid on which the resource is present or operation is being performed.

See detailed information in [Nutanix LCM Config V4] https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Config/operation/getConfig