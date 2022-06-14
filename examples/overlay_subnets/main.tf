terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.6.0"
    }
  }
}
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  insecure  = true
  port      = 9440
}

# pull all clusters
data "nutanix_clusters" "clusters" {}

# create local variable pointing to desired cluster
locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}
  

# overlay subnet resource

resource "nutanix_subnet" "acc" {
    # General Information
	name        = "test-overlay-subnet"
	description = "Description of my overlay subnet"

    # subnet type should be overlay
	subnet_type = "OVERLAY"
	subnet_ip          = "10.xx.xx.xx"
    default_gateway_ip = "10.xx.xx.xx"
    prefix_length = 24

    dhcp_options = {
        domain_name = "lab.fr"
        tftp_server_name = "tftp.lab.fr"
        boot_file_name = "pxelinux.0"
    }

    ip_config_pool_list_ranges = ["10.xx.xx.xx 10.xx.xx.xx"]

    # vpc reference uuid is required for overlay subnet type
    vpc_reference_uuid = var.vpc_reference_uuid
 
}

output "accSub"{
  value = resource.nutanix_subnet.acc
}