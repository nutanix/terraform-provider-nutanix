---
layout: "nutanix"
page_title: "NUTANIX: nutanix_sda_cluster_configs_v2"
sidebar_current: "docs-nutanix-datasource-sda-cluster-configs-v2"
description: |-
  Retrieves cluster specific configurations associated with a System-Defined Alert Policy identified by the external identifier of the System-Defined Alert Policy.
---

# nutanix_sda_cluster_configs_v2

Retrieves cluster specific configurations associated with a System-Defined Alert Policy identified by the external identifier of the System-Defined Alert Policy.

## Example

```hcl
data "nutanix_sda_cluster_configs_v2" "example" {
  system_defined_policy_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

- `system_defined_policy_ext_id`: (Required) Unique ID of the System-Defined Alert Policy.

## Attribute Reference

The following attributes are exported:

- `cluster_configs`: List of cluster-specific configurations for the SDA policy. Each entry has the same attributes as `nutanix_sda_cluster_config_v2`.

See `nutanix_sda_cluster_config_v2` datasource documentation for the full attribute reference of each cluster config.
