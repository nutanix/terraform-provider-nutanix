terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Add an SNMP user to a cluster.
resource "nutanix_snmp_user_v2" "example" {
  cluster_ext_id = var.cluster_ext_id
  username       = "tf-snmp-user"
  auth_type      = "MD5"
  auth_key       = "auth-key-12345678"
  priv_type      = "DES"
  priv_key       = "priv-key-12345678"
}

# Read the user back via the singular datasource.
data "nutanix_snmp_user_v2" "example" {
  cluster_ext_id = nutanix_snmp_user_v2.example.cluster_ext_id
  ext_id         = nutanix_snmp_user_v2.example.ext_id
}

output "snmp_user" {
  value     = data.nutanix_snmp_user_v2.example
  sensitive = true
}
