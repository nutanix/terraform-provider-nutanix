---
layout: "nutanix"
page_title: "NUTANIX: nutanix_sda_cluster_configs_v2"
sidebar_current: "docs-nutanix-datasource-sda-cluster-configs-v2"
description: |-
  Retrieves cluster specific configurations associated with a System-Defined Alert Policy identified by the external identifier of the System-Defined Alert Policy.
---

# nutanix_sda_cluster_configs_v2

Retrieves cluster specific configurations associated with a System-Defined Alert Policy identified by the external identifier of the System-Defined Alert Policy.

## Example Usage

```hcl
data "nutanix_sda_cluster_configs_v2" "example" {
  system_defined_policy_ext_id = "<system-defined-policy-ext-id>"
}
```

## Argument Reference

The following arguments are supported:

* `system_defined_policy_ext_id`: (Required) Unique ID of the System-Defined Alert Policy.
* `page`: (Optional) Page number for pagination.
* `limit`: (Optional) Number of results per page.
* `filter`: (Optional) OData filter expression.
* `order_by`: (Optional) OData order-by expression.
* `select`: (Optional) OData select expression.

## Attribute Reference

The following attributes are exported:

* `cluster_configs`: List of cluster configurations.

Each entry in `cluster_configs` contains the same attributes as the `nutanix_sda_cluster_config_v2` datasource. Please refer to the [nutanix_sda_cluster_config_v2](sda_cluster_config_v2) documentation for attribute details.

See detailed information in [Nutanix Monitoring v4 API Reference](https://developers.nutanix.com/).
