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



# define filters
locals {
  systemTypePCFilter  = "systemType eq Clustermgmt.Config.SystemType'PC'"
  systemTypeAOSFilter = "systemType eq Clustermgmt.Config.SystemType'AOS'"

  adminPCFilter = "username eq 'admin' and systemType eq Clustermgmt.Config.SystemType'PC'"
  adminAOSFilter = "username eq 'admin' and systemType eq Clustermgmt.Config.SystemType'AOS'"
}


# List Password Status Of All System Users
data "nutanix_system_user_passwords_v2" "passwords" {
}


# List Password Status Of All System Users With Limit
data "nutanix_system_user_passwords_v2" "limited_passwords" {
  limit  = 10
}


# List Password Status Of All System Users With Filter
data "nutanix_system_user_passwords_v2" "filtered_passwords" {
  filter = local.systemTypeAOSFilter
}

# List Password Status Of Admin PC User
data "nutanix_system_user_passwords_v2" "admin_pc_passwords" {
  filter = local.adminPCFilter
}


# change password for admin AOS user
data "nutanix_system_user_passwords_v2" "admin_aos_passwords" {
  filter = local.adminAOSFilter
}

resource "nutanix_password_change_request_v2" "change_admin_aos_password" {
	ext_id = data.nutanix_system_user_passwords_v2.admin_aos_passwords.passwords.0.ext_id
	current_password = var.current_password
	new_password = var.new_password
}
