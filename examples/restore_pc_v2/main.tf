terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1"
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


resource "nutanix_restore_pc_v2" "example"{
  restorable_domain_manager_ext_id = "<domain_manager_uuid>"
  restore_source_ext_id            = "<restore_source_uuid>"
  ext_id                           = "<restore_point_uuid>"
  domain_manager {
    config {
      name = "example-domain-manager"
      size = "SMALL"
    }
    network {
      external_address {
        ipv4 {
          value = "10.0.0.2"
        }
      }
      ntp_servers {
        ipv4 {
          value = "10.0.0.22"
        }
      }
      name_servers {
        ipv4 {
          value = "10.0.0.33"
        }
      }
    }
    should_enable_high_availability = true
  }
}
