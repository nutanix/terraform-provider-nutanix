# Service Provider Tenant Admin examples

Workflow examples for a **Tenant Admin** after a Service Provider has created the tenant and allocated a domain.

> **Note (KB-22890):** Calls to Nutanix Central may return **404** when the SP FQDN Host header includes an explicit `:443` port. This can affect Terraform and other API clients. Until a fix ships, add both `<host>` and `<host>:443` to `ALLOWED_CALLBACK_HOSTS` on the Dion service and restart it. Details: [KB article](https://portal.nutanix.com/page/documents/kbs/details?targetId=kA0VO000000EVaD0AW).

## Provider configuration

Tenant Admin calls route through Nutanix Central to the allocated PC domain. Set `domain_name` and `tenant_name`. Use `use_nc_iam = true` for membership / IAM workflows.

```hcl
provider "nutanix" {
  username                 = var.nutanix_username
  password                 = var.nutanix_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  domain_name              = var.domain_name
  tenant_name              = var.tenant_name
  use_nc_iam               = true # membership / IAM
  port                     = 443
}
```

Set values in each subdirectory's `terraform.tfvars` (replace placeholders with your environment):

```hcl
nutanix_username         = "<tenant-admin-username>"
nutanix_password         = "<tenant-admin-password>"
nutanix_central_endpoint = "<nutanix-central-fqdn>"
domain_name              = "<allocated-domain-name>"
tenant_name              = "<tenant-name>"
```

## Examples

1. **[tenant_membership](./tenant_membership)** — Assign users or groups to roles within the allocated domain (`nutanix_tenant_membership_v2`).
2. **[create_infra_on_allocated_domain](./create_infra_on_allocated_domain)** — Create infrastructure (image + VM) on the allocated domain using the tenant project and an existing VLAN.
