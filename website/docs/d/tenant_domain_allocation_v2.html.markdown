---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_domain_allocation_v2"
sidebar_current: "docs-nutanix-datasource-tenant-domain-allocation-v2"
description: |-
  Retrieves a specific domain allocation for a tenant by its external identifier.
---

# nutanix_tenant_domain_allocation_v2

Retrieves a specific domain allocation for a tenant by its external identifier.

## Example

```hcl
data "nutanix_tenant_domain_allocation_v2" "example" {
  tenant_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id        = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

- `tenant_ext_id`: (Required) External identifier of the tenant.
- `ext_id`: (Required) External identifier of the tenant domain allocation.

## Attributes Reference

The following attributes are exported:

- `domain_ext_id`: External identifier of the domain.
- `external_network_config`: External network connection settings for the domain allocation.
- `projects`: A project owns the resource allocation and placement settings within a domain allocation.
- `state_info`: The state details of the domain allocation.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `links`: A HATEOAS style link for the response.

### projects

- `ext_id`: External identifier of the project.
- `name`: The name of the project.
- `quota`: The resource quota of the project.
- `resource_group`: The resource group of the project.
- `state`: The state of the project.

See detailed information in [Nutanix Get Tenant Domain Allocation V4](https://developers.nutanix.com/).
