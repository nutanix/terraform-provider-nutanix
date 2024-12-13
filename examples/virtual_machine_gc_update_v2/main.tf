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

#update a virtual machine guest customization for next boot
resource "nutanix_vm_gc_update_v2" "test" {
  ext_id = var.vm_uuid
  config {
    cloud_init {
      cloud_init_script {
        user_data {
          value = var.user_data
        }
      }
    }
  }
}
