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


# fetch operations
data "nutanix_operations_v2" "operation-list" {
  filter = "startswith(displayName, 'Create_')"
}

# create role
resource "nutanix_roles_v2" "role" {
  display_name = "role_auth_example"
  description  = "role for authorization policy"
  operations = [
    data.nutanix_operations_v2.operation-list.operations[0].ext_id,
    data.nutanix_operations_v2.operation-list.operations[1].ext_id,
    data.nutanix_operations_v2.operation-list.operations[2].ext_id,
    data.nutanix_operations_v2.operation-list.operations[3].ext_id
  ]
}

resource "nutanix_authorization_policy_v2" "ap-example" {
  role                      = nutanix_roles_v2.role.id
  display_name              = "auth_policy_example"
  description               = "authorization policy example"
  authorization_policy_type = "USER_DEFINED"
  identities {
    reserved = "{\"user\":{\"uuid\":{\"anyof\":[\"00000000-0000-0000-0000-000000000000\"]}}}"
  }
  entities {
    reserved = "{\"images\":{\"*\":{\"eq\":\"*\"}}}"
  }
  entities {
    reserved = "{\"marketplace_item\":{\"owner_uuid\":{\"eq\":\"SELF_OWNED\"}}}"
  }
}

#get authorization policy by id
data "nutanix_authorization_policy_v2" "example" {
  ext_id = nutanix_authorization_policy_v2.ap-example.id
}


#list of authorization policies, with limit and filter
data "nutanix_authorization_policies_v2" "filtered-ap" {
  filter = "displayName eq '${nutanix_authorization_policy_v2.ap-example.display_name}'"
  limit  = 2
}

# list of authorization policies, with select
data "nutanix_authorization_policies_v2" "select-ap" {
  select     = "extId,displayName,description,authorizationPolicyType"
  depends_on = [nutanix_authorization_policy_v2.ap-example]
}
