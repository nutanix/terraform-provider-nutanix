terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}
provider "nutanix" {
  username = "admin"
  password = "Nutanix/123456"
  endpoint = "10.xx.xx.xx"
  insecure = true
  port     = 9440
}



#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#create local variable pointing to desired cluster
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

#creating subnet without IP pool
resource "nutanix_subnet_v2" "vlan-112" {
  name              = "vlan-112"
  description       = "subnet VLAN 112 managed by Terraform with IP pool"
  cluster_reference = local.clusterExtId
  subnet_type       = "VLAN"
  network_id        = 122
  is_external       = true
  ip_config {

    ipv4 {
      ip_subnet {
        ip {
          value = "192.168.0.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "192.168.0.1"
      }
      pool_list {
        start_ip {
          value = "192.168.0.20"
        }
        end_ip {
          value = "192.168.0.30"
        }
      }
    }
  }
}

// creating VPC
resource "nutanix_vpc_v2" "vpc" {
  name        = "vpc-example"
  description = "VPC for example"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.vlan-112.id
  }
}

# creating VPC with external routable prefixes
resource "nutanix_vpc_v2" "external-vpc-routable-vpc" {
  name        = "tf-vpc-example"
  description = "VPC "
  external_subnets {
    subnet_reference = nutanix_subnet_v2.vlan-112.id
    external_ips {
      ipv4 {
        value         = "192.168.0.24"
        prefix_length = 32
      }
    }
    external_ips {
      ipv4 {
        value         = "192.168.0.25"
        prefix_length = 32
      }
    }
  }
  externally_routable_prefixes {
    ipv4 {
      ip {
        value         = "172.30.0.0"
        prefix_length = 32
      }
      prefix_length = 16
    }
  }
}

// creating VPC with transit type
resource "nutanix_vpc_v2" "transit-vpc" {
  name        = "vpc-transit"
  description = "VPC for transit type"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.vlan-112.id
  }
  vpc_type = "TRANSIT"
}


//dataSource to get details for an entity with vpc uuid
data "nutanix_vpc_v2" "get-vpc" {
  ext_id = nutanix_vpc_v2.external-vpc-routable-vpc.id
}


// vpc list with filter
data "nutanix_vpcs_v2" "filter-vpc" {
  filter     = "name eq '${nutanix_vpc_v2.vpc.name}'"
  depends_on = [nutanix_vpc_v2.vpc, nutanix_vpc_v2.external-vpc-routable-vpc, nutanix_vpc_v2.transit-vpc]
}


// vpc list with limit and page
data "nutanix_vpcs_v2" "list-vpc-limit" {
  limit      = 3
  page       = 2
  depends_on = [nutanix_vpc_v2.vpc, nutanix_vpc_v2.external-vpc-routable-vpc, nutanix_vpc_v2.transit-vpc]
}
