---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_external_network_connection_v2"
sidebar_current: "docs-nutanix-resource-tenant-external-network-connection-v2"
description: |-
  Update external network connection settings for a tenant domain allocation.
---

# nutanix_tenant_external_network_connection_v2

Update external network connection settings for a tenant domain allocation. This
resource manages the external network connection that lives with a domain
allocation; it supports create/update (both map to the update API) and read only.

## Example

```hcl
resource "nutanix_tenant_external_network_connection_v2" "example" {
  tenant_ext_id            = "00000000-0000-0000-0000-000000000000"
  domain_allocation_ext_id = "11111111-1111-1111-1111-111111111111"

  config {
    nat {
      transit_vpc_ext_id    = "22222222-2222-2222-2222-222222222222"
      overlay_subnet_ext_id = "33333333-3333-3333-3333-333333333333"
      floating_ip_count     = 2
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `tenant_ext_id`: (Required) External identifier of the tenant.
- `domain_allocation_ext_id`: (Required) External identifier of the domain allocation.
- `config`: (Optional) The NAT or no-NAT external network connection configuration.

### config

- `nat`: (Optional) NAT configuration.
- `no_nat`: (Optional) No-NAT configuration.

#### nat

- `transit_vpc_ext_id`: (Required) The external identifier of the transit VPC for NAT configuration.
- `overlay_subnet_ext_id`: (Required) The external identifier of the overlay subnet for NAT configuration.
- `floating_ip_count`: (Optional) Total floating IPs to allocate. Existing assignments are preserved, and only additional IPs are created. Default is 0.

#### no_nat

- `overlay_subnet_ext_id`: (Optional) The external identifier of the overlay subnet for no-NAT configuration.
- `transit_vpc_ext_id`: (Optional) The external identifier of the transit VPC for no-NAT configuration.
- `externally_routable_prefixes`: (Optional) Externally Routable Prefix (ERP) CIDRs for no-NAT mode. The list replaces existing values on each request, and CIDRs are normalized.

## Attributes Reference

The following attributes are exported:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `state_info`: The state details of the external network connection.
- `state_info.state`: The state of the external network connection.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

See detailed information in [Nutanix Update Tenant External Network Connection V4](https://developers.nutanix.com/).
