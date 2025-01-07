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
  port     = 9440
  insecure = true
}


resource "nutanix_ngt_installation_v2" "example" {
  ext_id = "<VM UUID>"
  credential {
    username = var.username
    password = var.password
  }
  reboot_preference {
    schedule_type = "SKIP"
  }
  capablities = ["VSS_SNAPSHOT", "SELF_SERVICE_RESTORE"]
}


data "nutanix_ngt_configuration_v2" "example" {
  ext_id     = "<VM UUID>"
  depends_on = [nutanix_ngt_installation_v2.example]
}
