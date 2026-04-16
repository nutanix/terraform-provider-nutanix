terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.7.0"
        }
    }
}

#defining nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}

# Add a User group to the system.

# Add Directory Service user group.
resource "nutanix_user_groups" "t1" {
	directory_service_user_group{
		distinguished_name = "<distinguished name of the user group>"
	}
}

# Add Directory Service organizational unit.
resource "nutanix_user_groups" "t1" {
  directory_service_ou{
		distinguished_name = "<distinguished name of the organizational group>"
	}
}

# Add SAML Service user group.
resource "nutanix_user_groups" "t1" {
  saml_user_group{
    name = "<group name>"
    idp_uuid = "<idp uuid of user group>"
  }
}