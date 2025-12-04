terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.7.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Add a User group to the system.

resource "nutanix_user_groups_v2" "ldap-ug" {
  group_type         = "LDAP"
  idp_id             = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
  name               = "group_0664229e"
  distinguished_name = "cn=group_0664229e,ou=group,dc=devtest,dc=local"
}

# Saml User group
resource "nutanix_user_groups_v2" "saml-ug" {
  group_type = "SAML"
  idp_id     = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
  name       = "adfs19admingroup"
}


# List all the user groups in the system.
data "nutanix_user_groups_v2" "example" {
  depends_on = [ nutanix_user_groups_v2.ldap-ug, nutanix_user_groups_v2.saml-ug ]
}

# List user groups with a filter.
data "nutanix_user_groups_v2" "example" {
  filter = "name eq '${nutanix_user_groups_v2.ldap-ug.name}'"
}

# Get the details of a user group.
data "nutanix_user_group_v2" "example" {
  ext_id = nutanix_user_groups_v2.ldap-ug.id
}
