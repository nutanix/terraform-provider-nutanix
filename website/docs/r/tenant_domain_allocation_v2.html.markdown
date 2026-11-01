---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_domain_allocation_v2"
sidebar_current: "docs-nutanix-resource-tenant-domain-allocation-v2"
description: |-
  Creates a new domain allocation for a tenant.
---

# nutanix_tenant_domain_allocation_v2

Creates a new domain allocation for a tenant.

## Example

```hcl
resource "nutanix_tenant_domain_allocation_v2" "example" {
  tenant_ext_id = "00000000-0000-0000-0000-000000000000"
  domain_ext_id = "11111111-1111-1111-1111-111111111111"
  projects {
    quota {
      vcpu    = 4
      memory  = 8589934592
      storage = 53687091200
    }
    resource_group {
      placement_targets {
        cluster_uuid            = "22222222-2222-2222-2222-222222222222"
        storage_container_uuids = ["33333333-3333-3333-3333-333333333333"]
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `tenant_ext_id`: (Required) External identifier of the tenant.
- `domain_ext_id`: (Optional) External identifier of the domain.
- `external_network_config`: (Optional) External network connection settings for the domain allocation.
- `projects`: (Optional) A project owns the resource allocation and placement settings within a domain allocation.

### external_network_config

- `config`: (Optional) The NAT or no-NAT external network connection configuration.

#### config

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

### projects

- `ext_id`: (Optional) External identifier of the project.
- `name`: (Optional) The name of the project.
- `quota`: (Optional) The resource quota of the project.
- `resource_group`: (Optional) The resource group of the project.

#### quota

- `ext_id`: (Optional) The external identifier of the quota for the placement target.
- `memory`: (Optional) The amount of memory allocated as quota, in bytes.
- `storage`: (Optional) The amount of storage allocated as quota, in bytes.
- `vcpu`: (Optional) The number of vCPUs allocated as quota.

#### resource_group

- `uuid`: (Optional) The UUID of the resource group for the placement target.
- `placement_targets`: (Optional) The placement targets of the resource group.
- `placement_targets.cluster_uuid`: (Optional) The UUID of the cluster for resource placement.
- `placement_targets.storage_container_uuids`: (Optional) List of storage container UUIDs for the placement target.

## Attributes Reference

The following attributes are exported:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `state_info`: The state details of the domain allocation.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

### state_info

- `state`: The state of the domain allocation.

See detailed information in [Nutanix Create Tenant Domain Allocation V4](http://developers.nutanix.com/api-reference?namespace=tenancy&version=v4.0.a1#tag/Tenants/operation/createTenantDomainAllocation).
