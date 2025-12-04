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

data "nutanix_clusters_v2" "clusters" {}
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# create a subnet
resource "nutanix_subnet_v2" "example" {
  name              = "subnet_for_route"
  description       = "subnet to test create route"
  cluster_reference = local.clusterExtId
  subnet_type       = "VLAN"
  network_id        = "198"
  is_external       = true
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "10.44.3.192"
        }
        prefix_length = "27"
      }
      default_gateway_ip {
        value = "10.44.3.193"
      }
      pool_list {
        start_ip {
          value = "10.44.3.198"
        }
        end_ip {
          value = "10.44.3.207"
        }
      }
    }
  }
}

# crete a vpc
resource "nutanix_vpc_v2" "example" {
  name        = "terraform_test_vpc_1"
  description = "terraform test vpc 1 to test create route"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.example.id
  }
  depends_on = [nutanix_subnet_v2.example]
}

# List all route tables
data "nutanix_route_tables_v2" "list-route-tables" {
  depends_on = [nutanix_vpc_v2.example]
}

# List all route tables with order by vpcReference
data "nutanix_route_tables_v2" "route-tables-with-orderby" {
  order_by   = "vpcReference"
  depends_on = [nutanix_vpc_v2.example]
}


# list all route tables with limit
data "nutanix_route_tables_v2" "limited-tables" {
  limit      = 3
  depends_on = [nutanix_vpc_v2.example]
}

# fetch route table by id
data "nutanix_route_table_v2" "table-by-id" {
  ext_id     = data.nutanix_route_tables_v2.list-route-tables.route_tables.0.ext_id
  depends_on = [data.nutanix_route_tables_v2.list-route-tables]
}

