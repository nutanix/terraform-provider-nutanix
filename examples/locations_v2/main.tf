terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
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
  filter   = "startswith(country, 'United')"
  order_by = "name desc"
}
