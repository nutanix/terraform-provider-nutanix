terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1.0"
    }
  }
}

#defining nutanix configuration for PE
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}


// deploy pc, this is PE resource, make sure you configure the provider with PE endpoint and credentials
resource "nutanix_pc_deploy_v2" "example" {
  timeouts {
    create = "120m"
  }
  config {
    build_info {
      version = "pc.2024.3"
    }
    size = "STARTER"
    name = "PC_EXAMPLE"
  }
  network {
    external_networks {
      network_ext_id = "ba416f8d-00f2-499d-bc4c-19da8d104af9"
      default_gateway {
        ipv4 {
          value = "10.97.64.1"
        }
      }
      subnet_mask {
        ipv4 {
          value = "255.255.252.0"
        }
      }
      ip_ranges {
        begin {
          ipv4 {
            value = "10.97.64.91"
          }
        }
        end {
          ipv4 {
            value = "10.97.64.91"
          }
        }
      }
    }
    name_servers {
      ipv4 {
        value = "10.40.64.16"
      }
    }
    name_servers {
      ipv4 {
        value = "10.40.64.15"
      }
    }
    ntp_servers {
      fqdn {
        value = "2.centos.pool.ntp.org"
      }
    }
    ntp_servers {
      fqdn {
        value = "3.centos.pool.ntp.org"
      }
    }
  }
}
