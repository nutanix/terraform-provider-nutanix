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



#creating category
resource "nutanix_category_v2" "example" {
  key         = "category_example_key"
  value       = "category_example_value"
  description = "category example description"
}


#pull all categories data
data "nutanix_categories_v2" "categories-list" {}

# pull all categories with limit and filter
data "nutanix_categories_v2" "filtered-categories" {
  limit  = 2
  filter = "key eq '${nutanix_category_v2.example.key}'"
}

# get category by ext id
data "nutanix_category_v2" "get-category" {
  ext_id = nutanix_category_v2.example.id
}
