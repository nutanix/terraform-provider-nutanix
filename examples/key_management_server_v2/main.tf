terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.0"
    } 
  }
}

resource "nutanix_key_management_server_v2" "example" {
  name = "tf-kms"
  access_information {
    endpoint_url           = var.access_information.endpoint_url
    key_id                 = var.access_information.key_id
    tenant_id              = var.access_information.tenant_id
    client_id              = var.access_information.client_id
    client_secret          = var.access_information.client_secret
    credential_expiry_date = var.access_information.credential_expiry_date
  }
  lifecycle {
    ignore_changes = [
      access_information[0].client_secret,
      access_information[0].key_id
    ]
  }
}


data "nutanix_key_management_server_v2" "get-kms" {
  ext_id = nutanix_key_management_server_v2.example.id
}

data "nutanix_key_management_servers_v2" "kms-list" {
  depends_on = [nutanix_key_management_server_v2.example]
}
