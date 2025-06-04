terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.3.0"
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


# Create Service Account
resource "nutanix_users_v2" "service_account" {
	username = "service_account_terraform_example"
	description = "service account tf"
	email_id = "terraform_plugin@domain.com"
	user_type = "SERVICE_ACCOUNT"
}

# Create key under service account, never expires
resource "nutanix_user_key_v2" "create_key" {
  user_ext_id = nutanix_users_v2.service_account.ext_id
  name = "api_key_developers"
  key_type = "API_KEY"
  expiry_time = "2125-01-01T00:00:00Z"
  assigned_to = "developer_user_1"
}

// Get key details
data "nutanix_user_key_v2" "get_key"{
  user_ext_id = nutanix_users_v2.service_account.ext_id
  ext_id = nutanix_user_key_v2.create_key.ext_id
}

// To fetch the list of keys under service account
data "nutanix_user_keys_v2" "get_keys_filter" {
  user_ext_id = nutanix_users_v2.service_account.ext_id
  filter = "name eq '${nutanix_user_key_v2.create_key.name}'"
}


// Revoke the key
resource "nutanix_user_key_revoke_v2" "revoke-key"{
  user_ext_id = nutanix_users_v2.service_account.ext_id
  ext_id = nutanix_user_key_v2.create_key.ext_id
}