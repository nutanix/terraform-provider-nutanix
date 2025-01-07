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
  cluster1 = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

#creating subnet
resource "nutanix_subnet_v2" "vlan-112" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information
  name        = "vlan-112-managed"
  description = "subnet VLAN 112 managed by Terraform"
  vlan_id     = 112
  subnet_type = "VLAN"
  network_id  = 112
  is_external = true
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
resource "nutanix_vpc_v2" "test" {
  name        = "test-New-1"
  description = "%[2]s"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.vlan-112.id
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



//dataSource to get details for an entity with vpc uuid
data "nutanix_vpc_v2" "vpc1" {
  vpc_uuid = nutanix_vpc_v2.test.id
}

output "vpcOut1" {
  value = data.nutanix_vpc_v2.vpc1
}



//dataSource to get details for an entity with vpc name
data "nutanix_vpc_v2" "vpc2" {
  ext_id = "{{vpc_uuid}}"
}

output "vpcOut2" {
  value = data.nutanix_vpc.vpc2
}

// vpc list with filter
data "nutanix_vpcs_v2" "vpc3" {
  filter = "name eq 'test-New-1'"
}

output "vpcOut3" {
  value = data.nutanix_vpcs.vpc3
}

// vpc list with limit and page
data "nutanix_vpcs_v2" "vpc4" {
  limit = 3
  page  = 2
}

output "vpcOut4" {
  value = data.nutanix_vpcs.vpc4
}