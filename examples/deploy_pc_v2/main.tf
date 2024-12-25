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


// deploy pc
resource "nutanix_deploy_pc_v2" "pc"{
  config {
    build_info {
      version = "5.17.0"
    }
    size = "SMALL"
    name = "pc_example"
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
        value = "10.0.0.43"
      }
    }
  }
}

// list pcs
data "nutanix_pcs_v2" "pcs"{}

// get pc details
data "nutanix_pc_v2" "pc"{
  ext_id = nutanix_deploy_pc_v2.pc.id
}