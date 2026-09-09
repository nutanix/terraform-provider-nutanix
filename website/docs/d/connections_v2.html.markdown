---
layout: "nutanix"
page_title: "NUTANIX: nutanix_connections_v2"
sidebar_current: "docs-nutanix-datasource-connections-v2"
description: |-
  Returns a paginated list of connections for a specific hardware provider.
---

# nutanix_connections_v2

Returns a paginated list of connections for a specific hardware provider.

## Example

```hcl
data "nutanix_connections_v2" "example" {
  hardware_provider_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `hardware_provider_ext_id`: (Required) External ID of the hardware provider.
* `page`: (Optional) A URL query parameter specifying the page number of the result set.
* `limit`: (Optional) A URL query parameter specifying the total number of records returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties.

## Attributes Reference

The following attributes are exported:

* `connections`: List of connections for the hardware provider. Each element has the same attributes as the [`nutanix_connection_v2`](connection_v2.html) datasource.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
