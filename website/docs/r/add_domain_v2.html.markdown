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

- `ext_id`: External identifier of the allowed-domain association.
- `tenant_id`: Identifier of the tenant that owns the domain association.
- `links`: HATEOAS links for the domain association.
