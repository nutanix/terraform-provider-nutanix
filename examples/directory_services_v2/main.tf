terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.7.0"
    }
  }
}

#definig nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}



# Add Directory Service .
resource "nutanix_directory_services_v2" "example" {
  name           = "<name of directory service>"
  url            = "<URL for the Directory Service>"
  directory_type = "<Type of Directory Service.>"
  domain_name    = "<Domain name for the Directory Service.>"
  service_account {
    username = "<Username to connect to the Directory Service>"
    password = "<Password to connect to the Directory Service>"
  }
  white_listed_groups = ["example"]
}

# List all  Directory Services.
data "nutanix_directory_services_v2" "example" {}

# Get a Directory Service.
data "nutanix_directory_service_v2" "example" {
  ext_id = "<Directory Service UUID>"
}
