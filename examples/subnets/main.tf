terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.3.0"
        }
    }
}

#defining nutanix configuration
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
resource "nutanix_subnet" "vlan-112" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information
  name        = "vlan-112-managed"
  vlan_id     = 112
  subnet_type = "VLAN"

  # Managed L3 Networks
  # This bit is only needed if you intend to turn on IPAM
  prefix_length = 24

  default_gateway_ip = "10.xx.xx.xx"
  subnet_ip          = "10.xx.xx.xx"

  ip_config_pool_list_ranges = ["10.xx.xx.xx 10.xx.xx.xx"]

  dhcp_domain_name_server_list = ["10.xx.xx.xx"]
  dhcp_domain_search_list      = ["nxlab.fr"]

  dhcp_options = {
      domain_name = "lab.fr"
      tftp_server_name = "tftp.lab.fr"
      boot_file_name = "pxelinux.0"
  }
}

#output the subnet info
output "subnet" {
  value = nutanix_subnet.vlan-112
}