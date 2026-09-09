---
layout: "nutanix"
page_title: "NUTANIX: nutanix_nodes_v2"
sidebar_current: "docs-nutanix-datasource-nodes-v2"
description: |-
  Returns a list of nodes discovered through a specific hardware provider connection.
---

# nutanix_nodes_v2

Returns a list of nodes discovered through a specific hardware provider connection.

## Example

```hcl
data "nutanix_nodes_v2" "example" {
  hardware_provider_ext_id = "00000000-0000-0000-0000-000000000000"
  connection_ext_id        = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

* `hardware_provider_ext_id`: (Required) External ID of the hardware provider.
* `connection_ext_id`: (Required) External ID of the hardware provider connection.
* `page`: (Optional) A URL query parameter specifying the page number of the result set.
* `limit`: (Optional) A URL query parameter specifying the total number of records returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties.

## Attributes Reference

The following attributes are exported:

* `nodes`: List of discovered nodes. Each element has the same attributes as the [`nutanix_connection_node_v2`](connection_node_v2.html) datasource.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
