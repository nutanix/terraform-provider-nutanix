# Nutanix Central Admin examples

Workflow examples for a **Nutanix Central Admin**.

## Provider configuration

### Manage infrastructure on a registered PC domain

Use `nutanix_central_endpoint` with `domain_name`. Do **not** set `tenant_name`.

```hcl
provider "nutanix" {
  username                 = var.nutanix_username
  password                 = var.nutanix_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  domain_name              = var.domain_name
  port                     = 443
}
```

```hcl
nutanix_username         = "<nc-admin-username>"
nutanix_password         = "<nc-admin-password>"
nutanix_central_endpoint = "<nutanix-central-fqdn>"
domain_name              = "<registered-domain-name>"
```

### NC–PC registration

Registration uses two providers: NC for registered-domain / locations APIs, and PC for local-domain registration.

```hcl
provider "nutanix" {
  alias                    = "nc"
  username                 = var.nutanix_central_username
  password                 = var.nutanix_central_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  port                     = 443
}

provider "nutanix" {
  alias    = "pc"
  username = var.nutanix_pc_username
  password = var.nutanix_pc_password
  endpoint = var.nutanix_pc_endpoint
  port     = 9440
  insecure = true
}
```

```hcl
nutanix_central_username = "<nc-admin-username>"
nutanix_central_password = "<nc-admin-password>"
nutanix_central_endpoint = "<nutanix-central-fqdn>"
nutanix_pc_username      = "<pc-username>"
nutanix_pc_password      = "<pc-password>"
nutanix_pc_endpoint      = "<prism-central-ip-or-fqdn>"
```

## Examples

1. **[nc_pc_registration](./nc_pc_registration)** — Register a Prism Central with Nutanix Central:
   - Discover locations (`nutanix_locations_v2`)
   - Register the domain on NC (`nutanix_registered_domain_v2` + trust / credential helpers)
   - Complete local-domain registration on PC (`nutanix_local_domain_v2`)
2. **[create_infra](./create_infra)** — Create infrastructure on a registered PC domain via `nutanix_central_endpoint` + `domain_name` (no `tenant_name`).
