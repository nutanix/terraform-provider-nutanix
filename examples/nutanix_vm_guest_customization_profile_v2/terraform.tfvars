# Provider
nutanix_username = "admin"
nutanix_password = "password"
nutanix_endpoint = "10.xx.xx.xx"

# Infrastructure
subnet_name    = "your-subnet-name"
source_vm_uuid = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

# GC Profile settings
admin_password    = "YourAdminPassword"
domain_name       = "example.com"
domain_username   = "Administrator@example.com"
domain_password   = "YourDomainPassword"
domain_dns_server = "10.xx.xx.xx"

# Override settings (for template deploy / clone)
override_admin_password  = "YourOverridePassword"
override_domain_name     = "override.example.com"
override_domain_username = "administrator@override.example.com"
override_domain_password = "YourOverrideDomainPassword"
override_dns_server      = "10.xx.xx.xx"

# Static IPs for deploy and clone
deploy_ip_address    = "10.xx.xx.xx"
clone_ip_address     = "10.xx.xx.xx"
subnet_prefix_length = 27
subnet_gateway       = "10.xx.xx.xx"

# Windows product key (optional)
windows_product_key = "XXXXX-XXXXX-XXXXX-XXXXX-XXXXX"
