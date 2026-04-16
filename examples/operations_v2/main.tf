#Here we will get and list permissions
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



#list operations
data "nutanix_operations_v2" "operation-list" {}

# filtered list operation
data "nutanix_operations_v2" "operation-list-filtered" {
  filter = "displayName eq 'Create_Role'"
}

# list operations withe page and limit
data "nutanix_operations_v2" "operation-list-paginated" {
  page  = 1
  limit = 10
}

#get permission by ext-id
data "nutanix_operation_v2" "get-operation" {
  ext_id = data.nutanix_operations_v2.operation-list.operations.0.ext_id
}


