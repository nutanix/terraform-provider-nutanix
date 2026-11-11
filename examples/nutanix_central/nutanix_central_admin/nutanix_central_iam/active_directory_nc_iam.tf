# Stage 2: configure NC IAM Active Directory and Auth policies.
data "nutanix_roles_v2" "nc" {
  provider = nutanix.nc_iam
}

locals {
  nc_admin_role_ext_id = one([
    for role in data.nutanix_roles_v2.nc.roles :
    role.ext_id if role.display_name == var.nc_admin_role_name
  ])

  service_provider_admin_role_ext_id = one([
    for role in data.nutanix_roles_v2.nc.roles :
    role.ext_id if role.display_name == var.service_provider_admin_role_name
  ])
}

resource "nutanix_directory_services_v2" "nc_active_directory" {
  provider       = nutanix.nc_iam
  name           = var.nc_directory_service.name
  url            = var.nc_directory_service.url
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = var.nc_directory_service.domain_name

  // Stage 2 begins only after the PC LDAP administrator and its Prism Admin
  # policy have been configured in Stage 1.
  depends_on = [nutanix_authorization_policy_v2.pc_admin]

  service_account {
    username = var.nc_directory_service.service_account_username
    password = var.nc_directory_service.service_account_password
  }

  lifecycle {
    ignore_changes = [service_account[0].password]
  }
}

resource "nutanix_users_v2" "nc_ldap_user" {
  provider  = nutanix.nc_iam
  username  = var.nc_ldap_username
  user_type = "LDAP"
  idp_id    = nutanix_directory_services_v2.nc_active_directory.id
}

resource "nutanix_authorization_policy_v2" "nc_admin" {
  provider                  = nutanix.nc_iam
  role                      = local.nc_admin_role_ext_id
  display_name              = "Terraform Nutanix Central Admin policy"
  description               = "Grants the configured LDAP user Nutanix Central administrator access."
  authorization_policy_type = "USER_DEFINED"

  identities {
    reserved = jsonencode({
      user = { uuid = { anyof = [nutanix_users_v2.nc_ldap_user.ext_id] } }
    })
  }

  entities {
    reserved = jsonencode({ "*" = { "*" = { eq = "*" } } })
  }
}

resource "nutanix_authorization_policy_v2" "service_provider_admin" {
  provider                  = nutanix.nc_iam
  role                      = local.service_provider_admin_role_ext_id
  display_name              = "Terraform Service Provider Admin policy"
  description               = "Grants the configured LDAP user Service Provider administrator access."
  authorization_policy_type = "USER_DEFINED"

  identities {
    reserved = jsonencode({
      user = { uuid = { anyof = [nutanix_users_v2.nc_ldap_user.ext_id] } }
    })
  }

  entities {
    reserved = jsonencode({ "*" = { "*" = { eq = "*" } } })
  }
}
