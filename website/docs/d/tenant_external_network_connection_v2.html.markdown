---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_external_network_connection_v2"
sidebar_current: "docs-nutanix-datasource-tenant-external-network-connection-v2"
description: |-
  Retrieves external network connection settings for a tenant domain allocation.
---

# nutanix_tenant_external_network_connection_v2

Retrieves external network connection settings for a tenant domain allocation.

## Example

```hcl
data "nutanix_tenant_external_network_connection_v2" "example" {
  tenant_ext_id            = "00000000-0000-0000-0000-000000000000"
  domain_allocation_ext_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

- `tenant_ext_id`: (Required) External identifier of the tenant.
- `domain_allocation_ext_id`: (Required) External identifier of the domain allocation.

## Attributes Reference

The following attributes are exported:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `config`: The NAT or no-NAT external network connection configuration.
- `state_info`: The state details of the external network connection.
- `state_info.state`: The state of the external network connection.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `links`: A HATEOAS style link for the response.

See detailed information in [Nutanix Get Tenant External Network Connection V4](https://developers.nutanix.com/).
