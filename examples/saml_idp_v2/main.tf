#Here we will get and list permissions
#the variable "" present in terraform.tfvars file.
#Note - Replace appropriate values of variables in terraform.tfvars file as per setup

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.3.0"
    }
  }
}

#definig nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

resource "nutanix_saml_identity_providers_v2" "example" {
  name = "<IDENTITY_PROVIDER_NAME>"
  username_attribute = "<IDENTITY_PROVIDER_USERNAME>"
  email_attribute = "<IDENTITY_PROVIDER_EMAIL>"
  groups_attribute = "<IDENTITY_PROVIDER_GROUPS>"
  groups_delim = "<IDENTITY_PROVIDER_GROUPS_DELIM>" # such as ',' or ';'
  idp_metadata_xml = "<IDENTITY_PROVIDER_METADATA_XML>"
  entity_issuer = "<IDENTITY_PROVIDER_ENTITY_ISSUER>"
  is_signed_authn_req_enabled = "<IDENTITY_PROVIDER_IS_SIGNED_AUTHN_REQ_ENABLED>"
  custom_attributes = "<IDENTITY_PROVIDER_CUSTOM_ATTRIBUTES>"
}

#get saml identity provider by external id
data "nutanix_operation_v2" "permission" {
  ext_id = nutanix_saml_identity_providers_v2.example.id
}


#list permissions
data "nutanix_operations_v2" "permissions" {
  page   = 0
  limit  = 2
}
