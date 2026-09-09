---
layout: "nutanix"
page_title: "NUTANIX: nutanix_hardware_providers_v2"
sidebar_current: "docs-nutanix-datasource-hardware-providers-v2"
description: |-
  Returns a paginated list of all available hardware providers.
---

# nutanix_hardware_providers_v2

Returns a paginated list of all available hardware providers.

## Example

```hcl
data "nutanix_hardware_providers_v2" "example" {}
```

## Argument Reference

The following arguments are supported:

* `page`: (Optional) A URL query parameter specifying the page number of the result set.
* `limit`: (Optional) A URL query parameter specifying the total number of records returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria.
* `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties.

## Attributes Reference

The following attributes are exported:

* `hardware_providers`: List of hardware providers. Each element has the same attributes as the [`nutanix_hardware_provider_v2`](hardware_provider_v2.html) datasource.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
