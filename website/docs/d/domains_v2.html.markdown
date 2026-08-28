---
layout: "nutanix"
page_title: "NUTANIX: nutanix_domains_v2"
sidebar_current: "docs-nutanix-datasource-domains-v2"
description: |-
  Lists allowed domains in Nutanix Central tenancy.
---

# nutanix_domains_v2

Lists all allowed domains in Nutanix Central tenancy.

## Example

```hcl
data "nutanix_domains_v2" "example" {}
```

## Attributes Reference

- `domains`: List of allowed domains.
- `domains.domain_ext_id`: External identifier of the domain.
- `domains.ext_id`: External identifier of the allowed-domain association.
- `domains.tenant_id`: Identifier of the tenant that owns the domain association.
- `domains.links`: HATEOAS links for the domain association.
