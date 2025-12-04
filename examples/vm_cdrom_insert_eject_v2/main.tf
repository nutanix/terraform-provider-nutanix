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
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# pull image data
data "nutanix_images_v2" "vm-image" {
  filter = "name eq '${var.image_name}'"
  limit  = 1
}

#create a virtual machine with minium configuration
resource "nutanix_virtual_machine_v2" "example-1" {
  name                 = "vm-example"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster_ext_id
  }
  cd_roms {
    disk_address {
      bus_type = "IDE"
      index    = 0
    }
  }
  power_state = "OFF"
}

#insert a cdrom into the vm/ when delete this resource, it will eject the cdrom
resource "nutanix_vm_cdrom_insert_eject_v2" "test" {
  vm_ext_id = nutanix_virtual_machine_v2.example-1.id
  ext_id    = nutanix_virtual_machine_v2.example-1.cd_roms.0.ext_id
  backing_info {
    data_source {
      reference {
        image_reference {
          image_ext_id = data.nutanix_images_v2.vm-image.images[0].ext_id
        }
      }
    }
  }
}


# Eject the ISO which is inserted through Terraform, can be done in two ways:
# 1. By setting `action = "eject"` → triggers eject operation explicitly.
# 2. By deleting this resource → automatically ejects the ISO.
# resource "nutanix_vm_cdrom_insert_eject_v2" "insert-iso" {
#   vm_ext_id = nutanix_virtual_machine_v2.example-1.id
#   ext_id    = nutanix_virtual_machine_v2.example-1.cd_roms.0.ext_id
#   action    = "eject"
# }


# Incase if users need to eject the ISO which is not inserted through Terraform, they can do so by importing the resource.(Example: GUEST CUSTOMIZATION ISO, this is mounted on the CD-ROM by default during the Guest Customization process)
# // Step 1: Create a placeholder resource in your root module. For example:
# resource "nutanix_vm_cdrom_insert_eject_v2" "import_cdrom_inserted" {}

# // Step 2: execute this command in cli
# terraform import nutanix_vm_cdrom_insert_eject_v2.import_cdrom_inserted vm_ext_id/cdrom_ext_id

# // Step 3: Once imported, update the resource configuration(resource placeholder added in Step 1) to perform the eject operation
# resource "nutanix_vm_cdrom_insert_eject_v2" "import_cdrom_inserted" {
#   vm_ext_id = <Virtual_Machine_UUID>
#   ext_id    = <CD_ROM_UUID>
#   action    = "eject"
# }
