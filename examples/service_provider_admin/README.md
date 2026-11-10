# Service Provider Admin examples

Workflow examples for a **Service Provider Admin** operating against Nutanix Central (`nutanix_central_endpoint`).

## Provider configuration

Service Provider Admin calls use the NC Service Provider endpoint. Do **not** set `domain_name` or `tenant_name`.

```hcl
provider "nutanix" {
  username                 = var.nutanix_central_username
  password                 = var.nutanix_central_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  port                     = 443
}
```

Set values in each subdirectory's `terraform.tfvars` (replace placeholders with your environment):

```hcl
nutanix_central_username = "<sp-admin-username>"
nutanix_central_password = "<sp-admin-password>"
nutanix_central_endpoint = "<nutanix-central-fqdn>"
```

## Typical sequence

1. **[add_domain](./add_domain)** — Allow a registered Prism Central domain for tenancy (`nutanix_add_domain_v2` + `nutanix_domains_v2`).
2. **[create_tenant](./create_tenant)** — Create a tenant with a local admin and optional quota (`nutanix_tenant_v2` + list DS).
3. **[domain_allocation](./domain_allocation)** — Allocate domain capacity (and optional networking) to the tenant (`nutanix_tenant_domain_allocation_v2`).
