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

resource "nutanix_users_v2" "example" {
  username       = "<username>"
  first_name     = "<first name>"
  middle_initial = "<middle initial>"
  last_name      = "<last name>"
  email_id       = "<email id>"
  locale         = "<locale>"
  region         = "<region>"
  display_name   = "<user display name>"
  password       = "<user password>"
  # Type of the User LOCAL, LDAP, SAML
  user_type = "LOCAL"
  # Status of the User ACTIVE, INACTIVE
  status = "ACTIVE"
}


# List all the users in the system.
data "nutanix_users_v2" "test" {}

# List all the users with a filter.
data "nutanix_users_v2" "test" {
  filter = "username eq '<username>'"
}

# Get the details of a user.
data "nutanix_user_v2" "test" {
  ext_id = nutanix_users_v2.example.id
}
