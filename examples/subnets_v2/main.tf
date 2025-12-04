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


#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#create local variable pointing to desired cluster
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

#creating subnet with IP pool
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

#creating subnet without IP pool
resource "nutanix_subnet_v2" "vlan-113" {
  name              = "vlan-113"
  description       = "subnet VLAN 113 managed by Terraform"
  cluster_reference = local.clusterExtId
  subnet_type       = "VLAN"
  network_id        = 113
}

# creating subnet with IP pool and DHCP options
resource "nutanix_subnet_v2" "van-114" {
  name              = "vlan-114"
  description       = "subnet VLAN 114 managed by Terraform"
  cluster_reference = local.clusterExtId
  subnet_type       = "VLAN"
  network_id        = 114
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

  dhcp_options {
    domain_name_servers {
      ipv4 {
        value = "8.8.8.8"
      }
    }
    search_domains   = ["eng.nutanix.com"]
    domain_name      = "nutanix.com"
    tftp_server_name = "10.5.0.10"
    boot_file_name   = "pxelinux.0"
  }
}

# pull all subnets data
data "nutanix_subnets_v2" "subnets" {
  depends_on = [nutanix_subnet_v2.vlan-112, nutanix_subnet_v2.vlan-113, nutanix_subnet_v2.van-114]
}

# list all subnets with filter
data "nutanix_subnets_v2" "filtered-subnets" {
  filter     = "clusterReference eq '${local.clusterExtId}'"
  depends_on = [nutanix_subnet_v2.vlan-112, nutanix_subnet_v2.vlan-113, nutanix_subnet_v2.van-114]
}

# pull all subnets data with filter and limit
data "nutanix_subnets_v2" "subnets-filter-limit" {
  filter     = "isExternal eq true"
  limit      = 2
  depends_on = [nutanix_subnet_v2.vlan-112, nutanix_subnet_v2.vlan-113, nutanix_subnet_v2.van-114]
}

# fetch the subnet data by ID
data "nutanix_subnet_v2" "subnet" {
  ext_id     = nutanix_subnet_v2.vlan-112.id
  depends_on = [nutanix_subnet_v2.vlan-112, nutanix_subnet_v2.vlan-113, nutanix_subnet_v2.van-114]
}
