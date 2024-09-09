terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.7.0"
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

# Add Address group.
resource "nutanix_address_groups_v2" "example_1" {
  name = "address_group"
  description = "address group description"
  ipv4_addresses{
    value = "10.0.0.0"
    prefix_length = 24
  }
}

# Add Address group. with ip range
resource "nutanix_address_groups_v2" "example_2" {
  name = "address_group"
  description = "address group description"
  ip_ranges{
    start_ip = "10.0.0.1"
    end_ip = "10.0.0.10"
  }
}