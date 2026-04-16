terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.9.5"
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


# Create a category key
resource "nutanix_category_key" "example-category-key" {
  name        = "app-support"              # this value can be updated
  description = "App Support Category Key" # this value can be updated
}

# Create a category value
resource "nutanix_category_value" "example-category-value" {
  name        = nutanix_category_key.example-category-key.id
  description = "Example Category Value" # this value can be updated
  value       = "example-value"          # this value can be updated
}

# Get the category key by name
data "nutanix_category_key" "get_key" {
  name = nutanix_category_key.example-category-key.name
}
