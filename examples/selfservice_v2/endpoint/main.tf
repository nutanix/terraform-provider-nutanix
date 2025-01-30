terraform {
  required_providers {
    nutanix = {
      source  = "nutanixtemp/nutanix"
      version = "1.99.99"
    }
  }
}

provider "nutanix" {
  username = "admin"
  password = "Nutanix.123"
  endpoint = "10.101.176.123"
  insecure = true
  port     = 9440
}

resource "nutanix_calm_endpoint" "endpoint" {
  name = "a"
  description = "d"
  ip_address = "10.10.10.10"
  cred_username = "nutanix"
  cred_password = "nutanix/4u"
}
