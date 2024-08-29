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

#pull all clusters data
data "nutanix_clusters" "clusters"{}

#create local variable pointing to desired cluster
locals {
	cluster1 = [
	  for cluster in data.nutanix_clusters.clusters.entities :
	  cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
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
  network_id = 112
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
      pool_list{
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

#output the subnet info
output "subnet" {
  value = nutanix_subnet_v2.vlan-112
}