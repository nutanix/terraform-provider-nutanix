---
layout: "nutanix"
page_title: "NUTANIX: nutanix_server_identity_pool_v2"
sidebar_current: "docs-nutanix-datasource-server-identity-pool-v2"
description: |-
  Returns details of a server identity pool from a hardware provider connection.
---

# nutanix_server_identity_pool_v2

Returns details of a server identity pool from a hardware provider connection.

## Example

```hcl
data "nutanix_server_identity_pool_v2" "example" {
  hardware_provider_ext_id = "00000000-0000-0000-0000-000000000000"
  connection_ext_id        = "11111111-1111-1111-1111-111111111111"
  ext_id                   = "22222222-2222-2222-2222-222222222222"
}
```

## Argument Reference

The following arguments are supported:

* `hardware_provider_ext_id`: (Required) External ID of the hardware provider.
* `connection_ext_id`: (Required) External ID of the hardware provider connection.
* `ext_id`: (Required) External ID of the server identity pool.

## Attributes Reference

The following attributes are exported:

* `name`: Name of the server identity pool.
* `reference_id`: Reference identifier of the pool from the provider.
* `available_count`: Number of available server identities.
* `groups`: Groups associated with the pool.
* `display_name_mapping`: Display name mapping for the pool.
* `last_updated_time`: Last time the pool was updated.
* `links`: A HATEOAS style list of links for the response.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
