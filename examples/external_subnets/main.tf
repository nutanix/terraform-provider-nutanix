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

# External Subnets with NAT

resource "nutanix_subnet" "accNat" {
  # General Info
  name        = "test-ext-sub-with-nat"
  description = "Description of my test VLAN updated"

  # subnet_type should be VLAN for external subnet with NAT
  subnet_type = "VLAN"
  cluster_uuid = local.cluster1
  vlan_id = 121
  subnet_ip          = "10.xx.xx.xx"
  default_gateway_ip = "10.xx.xx.xx"
  prefix_length = 24

  # required to be set true for external connectivity
  is_external = true
  # set true if NAT reuired
  enable_nat = true

  ip_config_pool_list_ranges = ["10.xx.xx.xx 10.xx.xx.xx"]
 
}

# External Subnet with No NAT

resource "nutanix_subnet" "accNoNat" {
  # General Info
  name        = "test-ext-sub-with-no-nat"
  description = "Description of my test VLAN updated"

  # subnet_type should be VLAN for external subnet with No NAT
  subnet_type = "VLAN"
  cluster_uuid = local.cluster1
  vlan_id = 121
  subnet_ip          = "10.xx.xx.xx"
  default_gateway_ip = "10.xx.xx.xx"
  prefix_length = 24

  # required to be set true for external connectivity
  is_external = true
  # set fasle for No NAT
  enable_nat = false

  ip_config_pool_list_ranges = ["10.xx.xx.xx 10.xx.xx.xx"]
 
}