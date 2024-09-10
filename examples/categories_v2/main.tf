terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.3.0"
        }
    }
}

#definig nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}



#creating category
resource "nutanix_category_v2" "example" {
  key = "category_example_key"
  value = "category_example_value"
  description = "category example description"
}


#pull all categories data
data "nutanix_categories_v2" "clusters"{}


# get category by ext id
data "nutanix_category_v2" "example" {
  ext_id = resource.nutanix_category_v2.example.ext_id
}