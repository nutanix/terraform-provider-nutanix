#Here we will get and list permissions
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

#get permission by ext-id
data "nutanix_operation_v2" "permission" {
  ext_id = var.permission_ext_id
}


#list permissions
data "nutanix_operations_v2" "permissions" {
  page   = 0
  limit  = 2
  filter = "displayName eq 'test-Permission-filter'"
}
