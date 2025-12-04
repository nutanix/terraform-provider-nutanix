terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

# Defining Nutanix provider configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# -------------------------------------------------
# Create an ACTIVE local user in the system
# -------------------------------------------------
resource "nutanix_users_v2" "active-user" {
  username             = "example_user"
  first_name           = "first-name"
  middle_initial       = "middle-initial"
  last_name            = "last-name"
  email_id             = "example_user@email.com"
  locale               = "en-us"
  region               = "en-us"
  display_name         = "display-name"
  password             = "example.password"
  user_type            = "LOCAL"
  status               = "ACTIVE"
  force_reset_password = true
}

# -------------------------------------------------
# Create an INACTIVE local user in the system
# -------------------------------------------------
resource "nutanix_users_v2" "inactive-user" {
  username             = "inactive_user"
  first_name           = "first-name"
  middle_initial       = "middle-initial"
  last_name            = "last-name"
  email_id             = "example_user@email.com"
  locale               = "en-us"
  region               = "en-us"
  display_name         = "display-name"
  password             = "example.password"
  user_type            = "LOCAL"
  status               = "INACTIVE"
  force_reset_password = true
}

# -------------------------------------------------
# Create an LDAP user in the system
# -------------------------------------------------
# This resource creates a user that is authenticated via an LDAP Identity Provider (IdP).
# You must provide the IdP UUID (idp_id) associated with your LDAP configuration.
resource "nutanix_users_v2" "ldap-user" {
  username  = "ldap_user"
  user_type = "LDAP"
  idp_id    = var.ldap_idp_id
}

# -------------------------------------------------
# Create a SAML user in the system
# -------------------------------------------------
# This resource creates a user that is authenticated via a SAML Identity Provider (IdP).
# You must provide the IdP UUID (idp_id) associated with your SAML configuration.
resource "nutanix_users_v2" "saml-user" {
  username  = "saml_user"
  user_type = "SAML"
  idp_id    = var.sam_idp_id
}

# -------------------------------------------------
# Retrieve a list of all users in the system
# -------------------------------------------------
data "nutanix_users_v2" "list-users" {
  depends_on = [
    nutanix_users_v2.active-user,
    nutanix_users_v2.inactive-user,
    nutanix_users_v2.ldap-user,
    nutanix_users_v2.saml-user
  ]
}

# -------------------------------------------------
# Retrieve a filtered list of users based on username
# -------------------------------------------------
data "nutanix_users_v2" "test" {
  filter = "username eq '${nutanix_users_v2.active-user.username}'"
}

# -------------------------------------------------
# Retrieve details of a specific user using ext_id
# -------------------------------------------------
data "nutanix_user_v2" "get-user" {
  ext_id = nutanix_users_v2.active-user.id
}

# -------------------------------------------------
# Create a Service Account user
# -------------------------------------------------
resource "nutanix_users_v2" "service_account" {
  username    = "service_account_terraform_example"
  description = "service account tf"
  email_id    = "terraform_plugin@domain.com"
  user_type   = "SERVICE_ACCOUNT"
}

# -------------------------------------------------
# Retrieve Service Account details using ext_id
# -------------------------------------------------
data "nutanix_user_v2" "get_service_account" {
  ext_id = nutanix_users_v2.service_account.id
}

# -------------------------------------------------
# Retrieve list of Service Accounts using a filter
# -------------------------------------------------
data "nutanix_users_v2" "list_service_account" {
  filter = "userType eq Schema.Enums.UserType'SERVICE_ACCOUNT' and username contains '${nutanix_users_v2.service_account.username}'"
}
