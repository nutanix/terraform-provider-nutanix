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



# Create Address group with ipv4 addresses
resource "nutanix_address_groups_v2" "ipv4-address" {
  name        = "address_group_ipv4_address"
  description = "address group description"
  ipv4_addresses {
    value         = "10.0.0.0"
    prefix_length = 24
  }
  ipv4_addresses {
    value         = "172.0.0.0"
    prefix_length = 24
  }
}

# Create Address group. with ip range
resource "nutanix_address_groups_v2" "ip-ranges" {
  name        = "address_group_ip_ranges"
  description = "address group description"
  ip_ranges {
    start_ip = "10.0.0.1"
    end_ip   = "10.0.0.10"
  }
}

# list add address group
data "nutanix_address_groups_v2" "list-address-groups" {
  depends_on = [nutanix_address_groups_v2.ipv4-address, nutanix_address_groups_v2.ip-ranges]
}

# list add address group with filter
data "nutanix_address_groups_v2" "example-filter" {
  filter     = "name eq '${nutanix_address_groups_v2.ipv4-address.name}'"
  depends_on = [nutanix_address_groups_v2.ipv4-address, nutanix_address_groups_v2.ip-ranges]
}

# list add address group with order by
data "nutanix_address_groups_v2" "example-order-by" {
  order_by   = "name desc"
  depends_on = [nutanix_address_groups_v2.ipv4-address, nutanix_address_groups_v2.ip-ranges]
}

# list add address group select
data "nutanix_address_groups_v2" "example-select" {
  select     = "name,description,ipRanges"
  depends_on = [nutanix_address_groups_v2.ipv4-address, nutanix_address_groups_v2.ip-ranges]
}
