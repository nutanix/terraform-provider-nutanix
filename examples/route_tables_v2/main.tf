terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "2.0"
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

# list all route tables
data "nutanix_route_tables_v2" "all-tables" {
}

# list all route tables with filter
data "nutanix_route_tables_v2" "filtered-tables" {
    filter = "vpcReference eq '<vpc_uuid>'"
}

# list all route tables with limit
data "nutanix_route_tables_v2" "limited-tables" {
    limit = 3
}

# fetch route table by id
data "nutanix_route_table_v2" "table-by-id" {
  ext_id = "<route_table_uuid>"
}

