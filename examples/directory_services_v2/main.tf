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

# Add Directory Service .
resource "nutanix_directory_services_v2" "active-directory" {
  name           = "example_active_directory"
  url            = "ldap://10.xx.xx.xx:xxxx"
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = "nutanix.com"
  service_account {
    username = "username"
    password = "password"
  }
  white_listed_groups = ["example"]
  lifecycle {
    ignore_changes = [
      service_account.0.password,
    ]
  }
}

# Add Directory Service and share with all projects.
resource "nutanix_directory_services_v2" "active-directory-share-with-all-projects" {
  name           = "example_active_directory"
  url            = "ldap://10.xx.xx.xx:xxxx"
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = "nutanix.com"
  share_with_all_projects = true
  service_account {
    username = "username"
    password = "password"
  }
  white_listed_groups = ["example"]
  lifecycle {
    ignore_changes = [
      service_account.0.password,
    ]
  }
}

# Add Directory Service and share with projects.
resource "nutanix_directory_services_v2" "active-directory-share-with-all-projects" {
  name           = "example_active_directory"
  url            = "ldap://10.xx.xx.xx:xxxx"
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = "nutanix.com"
  shared_with_projects = ["Project1_UUID", "Project2_UUID"]
  service_account {
    username = "username"
    password = "password"
  }
  white_listed_groups = ["example"]
  lifecycle {
    ignore_changes = [
      service_account.0.password,
    ]
  }
}


# List all  Directory Services.
data "nutanix_directory_services_v2" "example" {}

# List all  Directory Services with filter.
data "nutanix_directory_services_v2" "list-active-directory" {
  filter = "name eq '${nutanix_directory_services_v2.active-directory.name}'"
}

# Get a Directory Service.
data "nutanix_directory_service_v2" "get-active-directory" {
  ext_id = nutanix_directory_services_v2.active-directory.ext_id
}
