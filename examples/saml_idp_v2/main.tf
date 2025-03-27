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

resource "nutanix_saml_identity_providers_v2" "example" {
  name                        = "example_idp_name"
  idp_metadata {
    entity_id = "entity_id"
    login_url = "login_url"
    logout_url = "logout_url"
    error_url = "error_url"
    certificate = "certificate"
  }
  username_attribute          = "username"
  email_attribute             = "email"
  groups_attribute            = "groups"
  groups_delim                = "," # such as ',' or ';'
  idp_metadata_xml            = "<IDENTITY_PROVIDER_METADATA_XML content>"
  entity_issuer               = "entity_issuer_issuer"
  is_signed_authn_req_enabled = true
  custom_attributes           = ["custom1", "custom2"]
}

# get saml identity provider by id
data "nutanix_saml_identity_provider_v2" "example" {
  ext_id = nutanix_saml_identity_providers_v2.example.id
}

# list saml identity providers
data "nutanix_saml_identity_providers_v2" "examples" {
  limit  = 2
  filter = "name eq '<IDENTITY_PROVIDER_NAME>'"
}
