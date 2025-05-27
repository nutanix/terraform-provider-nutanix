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
resource "nutanix_subnet_v2" "ext-subnet" {
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
resource "nutanix_vpc_v2" "vpc" {
  name        = "terraform_example_vpc_1"
  description = "terraform example vpc 1 to test create route"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.ext-subnet.id
  }
  depends_on = [nutanix_subnet_v2.ext-subnet]
}

# get route table
data "nutanix_route_tables_v2" "list-route-tables" {
  filter     = "vpcReference eq '${nutanix_vpc_v2.vpc.id}'"
  depends_on = [nutanix_vpc_v2.vpc]
}

# create a project
resource "nutanix_project" "example-project" {
  name        = "tf-example-project"
  description = "terraform example project"
  default_subnet_reference {
    kind = "subnet"
    uuid = nutanix_subnet_v2.ext-subnet.id
  }
  lifecycle {
    ignore_changes = [default_subnet_reference]
  }
}
# create a route
resource "nutanix_routes_v2" "route" {
  name               = "terraform_example_route"
  description        = "terraform example route to example create route"
  vpc_reference      = nutanix_vpc_v2.vpc.id
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
  destination {
    ipv4 {
      ip {
        value = "10.0.0.2"
      }
      prefix_length = 32
    }
  }
  next_hop {
    next_hop_type      = "EXTERNAL_SUBNET"
    next_hop_reference = nutanix_subnet_v2.ext-subnet.id
  }
  metadata {
    owner_reference_id   = nutanix_vpc_v2.vpc.id
    project_reference_id = nutanix_project.example-project.metadata.uuid
  }
  route_type = "STATIC"
}

# list all routes
data "nutanix_routes_v2" "all-routes" {
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
  depends_on         = [nutanix_routes_v2.route]
}

# list all routes with filter
data "nutanix_routes_v2" "filtered-routes" {
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
  filter             = "vpcReference eq '${nutanix_vpc_v2.vpc.id}' and isActive eq true"
  order_by           = "name desc" # list all routes sorted by name in descending order
  depends_on         = [nutanix_routes_v2.route]
}

# list all routes with limit
data "nutanix_routes_v2" "limited-routes" {
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
  limit              = 3
  depends_on         = [nutanix_routes_v2.route]
}

# fetch route by id
data "nutanix_route_v2" "route-by-id" {
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
  ext_id             = data.nutanix_routes_v2.filtered-routes.routes[0].ext_id
  depends_on         = [nutanix_routes_v2.route]
}
