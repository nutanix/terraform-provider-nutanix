---
layout: "nutanix"
page_title: "NUTANIX: nutanix_hardware_provider_v2"
sidebar_current: "docs-nutanix-datasource-hardware-provider-v2"
description: |-
  Returns details of a hardware provider identified by its external ID.
---

# nutanix_hardware_provider_v2

Returns details of a hardware provider identified by its external ID.

## Example

```hcl
data "nutanix_hardware_provider_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) External ID of the hardware provider.

## Attributes Reference

The following attributes are exported:

* `name`: Name of the hardware provider.
* `vendor`: Vendor of the hardware provider.
* `type`: Type of the hardware provider.
* `auth_types`: List of supported authentication types.
* `available_connection_count`: Number of available connections.
* `links`: A HATEOAS style list of links for the response.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
