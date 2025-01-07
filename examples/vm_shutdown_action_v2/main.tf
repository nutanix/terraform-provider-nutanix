#Here we will create a vm clone 
#the variable "" present in terraform.tfvars file.
#Note - Replace appropriate values of variables in terraform.tfvars file as per setup

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

#create a virtual machine with minium configuration
resource "nutanix_virtual_machine_v2" "example-1" {
  name                 = "vm-example"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
  }
  power_state = "ON"
}

# shutdown the virtual machine
resource "nutanix_vm_shutdown_action_v2" "vmShuts" {
  ext_id = nutanix_virtual_machine_v2.example-1.id
  action = "shutdown"
}

# restart the virtual machine
resource "nutanix_vm_shutdown_action_v2" "vmReboot" {
  ext_id = nutanix_virtual_machine_v2.example-1.id
  action = "reboot"
}

# guest-shutdown the virtual machine
resource "nutanix_vm_shutdown_action_v2" "vmGuestShuts" {
  ext_id = nutanix_virtual_machine_v2.example-1.id
  action = "guest_shutdown"
}

# guest-restart the virtual machine
resource "nutanix_vm_shutdown_action_v2" "vmGuestReboot" {
  ext_id = nutanix_virtual_machine_v2.example-1.id
  action = "guest_reboot"
}
