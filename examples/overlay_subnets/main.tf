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


# create a new overlay subnet with vpc and external subnet

# pull all clusters 
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

# create external subnet
resource "nutanix_subnet" "sub-ext" {
  cluster_uuid = local.cluster1
	name        = "test-ext-subnet"
	description = "Description of my unit test VLAN"
	vlan_id     = 434
	subnet_type = "VLAN"
	subnet_ip          = "10.xx.xx.xx"
  default_gateway_ip = "10.xx.xx.xx"
	ip_config_pool_list_ranges = ["10.xx.xx.xx 10.xx.xx.xx"]

  prefix_length = 24
	is_external = true
	enable_nat = false
}

# create vpc with external subnet reference
resource "nutanix_vpc" "acctest-managed-vpc" {
	name = "test-vpc"
	external_subnet_reference_uuid = [
	  resource.nutanix_subnet.sub-ext.id
	]
	common_domain_name_server_ip_list{
		ip = "x.x.x.x"
	}
	externally_routable_prefix_list{
	  ip=  "xx.xx.xx.xx"
	  prefix_length= 16
	}
}


# create overlay subnet with vpc reference
  resource "nutanix_subnet" "acctest-managed" {
	name        = "test-overlay-subnet"
	description = "Description of my unit test OVERLAY"
	vpc_reference_uuid = resource.nutanix_vpc.acctest-managed-vpc.id
	subnet_type = "OVERLAY"
	subnet_ip          = "10.xx.xx.xx"
	default_gateway_ip = "10.xx.xx.xx"
	ip_config_pool_list_ranges = ["10.xx.xx.xx 10.xx.xx.xx"]
	prefix_length = 24
}