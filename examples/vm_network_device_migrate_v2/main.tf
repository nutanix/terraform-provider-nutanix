#Here we will create a vm clone 
#the variable "" present in terraform.tfvars file.
#Note - Replace appropriate values of variables in terraform.tfvars file as per setup

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster0 = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# migrate the network device of the vm
resource "nutanix_vm_network_device_migrate_v2" "example" {
  vm_ext_id = "<vm_ext_id>"
  ext_id    = "<vm_nic_ext_id>"
  subnet {
    ext_id = "<subnet_ext_id>"
  }
  migrate_type = "RELEASE_IP"
}
