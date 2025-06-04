terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Add a User group to the system.

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

resource "nutanix_users_v2" "ldap-user" {
  username  = "ldap_user"
  user_type = "LDAP"
  idp_id    = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
}

resource "nutanix_users_v2" "saml-user" {
  username  = "saml_user"
  user_type = "SAML"
  idp_id    = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
}

# List all the users in the system.
data "nutanix_users_v2" "list-users" {
  depends_on = [nutanix_users_v2.active-user, nutanix_users_v2.inactive-user, nutanix_users_v2.ldap-user, nutanix_users_v2.saml-user]
}

# List all the users with a filter.
data "nutanix_users_v2" "test" {
  filter = "username eq '${nutanix_users_v2.active-user.username}'"
}

# Get the details of a user.
data "nutanix_user_v2" "get-user" {
  ext_id = nutanix_users_v2.active-user.id
}

# Create Service Account
resource "nutanix_users_v2" "service_account" {
  username = "service_account_terraform_example"
  description = "service account tf"
  email_id = "terraform_plugin@domain.com"
  user_type = "SERVICE_ACCOUNT"
}

# Get Service Account using the ext_id
data "nutanix_user_v2" "get_service_account" {
	ext_id = nutanix_users_v2.service_account.id
}

# Get list of Service Accounts with Filter
data "nutanix_users_v2" "list_service_account" {
	filter = "userType eq Schema.Enums.UserType'SERVICE_ACCOUNT' and username contains '${nutanix_users_v2.service_account.username}'"
}
