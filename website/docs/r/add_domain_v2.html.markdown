---
layout: "nutanix"
page_title: "NUTANIX: nutanix_add_domain_v2"
sidebar_current: "docs-nutanix-resource-add-domain-v2"
description: |-
  Adds an allowed domain to Nutanix Central tenancy.
---

# nutanix_add_domain_v2

Adds an existing domain to the allowed-domain list for Nutanix Central tenancy.

The tenancy API does not provide update or delete operations for domains.
Removing this resource from Terraform configuration removes only its Terraform
state; it does not remove the allowed-domain association from Nutanix Central.

## Prerequisites

The domain must have the following services enabled before it can be added to
the tenancy allowed-domain list:

- Flow Network Security
- Flow Virtual Networking
- Policy Engine

## Example

```hcl
resource "nutanix_add_domain_v2" "example" {
  domain_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following argument is supported:

- `domain_ext_id`: (Required) External identifier of the existing domain to add.

## Attributes Reference

The following attributes are exported:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption..
- `domain_ext_id`: External identifier of the domain.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

See detailed information in [Nutanix Add Domain V4](https://developers.nutanix.com/api-reference?namespace=tenancy&version=v4.0.a1#tag/Tenants/operation/addDomain).
