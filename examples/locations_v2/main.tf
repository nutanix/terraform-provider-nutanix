terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.1"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username # NC Admin Username
  password = var.nutanix_password # NC Admin Password
  nutanix_central_endpoint = var.nutanix_central_endpoint # NC Central Environment FQDN
  port = "443"
}

# List all locations (default page/limit applied by the API)
data "nutanix_locations_v2" "all" {}

# List locations with pagination
data "nutanix_locations_v2" "paged" {
  page  = 0
  limit = 20
}

# Filter locations by country and sort them by name
data "nutanix_locations_v2" "filtered" {
  limit    = 10
  filter   = "startswith(country, 'India')"
  order_by = "name desc"
}
