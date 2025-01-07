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

# create authorization policy
resource "nutanix_authorization_policy_v2" "auth_policy_example" {
  role         = "<role_uuid>"
  display_name = "<acp name>"
  description  = "<acp description>"
  authorization_policy_type = "<acp type>"
  # identity and entity will defined as a json string
  identities {
    reserved = "<identity_uuid>" # ex : "{\"user\":{\"uuid\":{\"anyof\":[\"00000000-0000-0000-0000-000000000000\"]}}}"
  }
  entities {
    reserved = "<entity_uuid>" # ex : "{\"images\":{\"*\":{\"eq\":\"*\"}}}"
  }
}

#get authorization policy by id
data "nutanix_authorization_policy_v2" "example" {
  ext_id = nutanix_authorization_policy_v2.auth_policy_example.id
}


#list of authorization policies, with limit and filter
data "nutanix_authorization_policies_v2" "examples" {
  limit  = 2
  filter = "display_name eq '<acp name>'"
}
