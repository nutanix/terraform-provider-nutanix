terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
    }
  }
}

#definig nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}



// DomainManagerRemoteClusterSpec
resource "nutanix_pe_pc_registration_v2" "pc1" {
  pc_ext_id = "00000000-0000-0000-0000-000000000000"
  remote_cluster {
    domain_manager_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = "0.0.0.0"
          }
        }
        credentials {
          authentication {
            username = "example"
            password = "example.123"
          }
        }
      }
      cloud_type = "NUTANIX_HOSTED_CLOUD"
    }
  }
}

// AOSRemoteClusterSpec
resource "nutanix_pe_pc_registration_v2" "pc1" {
  pc_ext_id = "00000000-0000-0000-0000-000000000000"
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = "0.0.0.0"
          }
        }
        credentials {
          authentication {
            username = "example"
            password = "example.123"
          }
        }
      }
    }
  }
}

// ClusterReference
resource "nutanix_pe_pc_registration_v2" "pc1" {
  pc_ext_id = "00000000-0000-0000-0000-000000000000"
  remote_cluster {
    cluster_reference {
      ext_id = "11111111-1111-1111-1111-111111111111"
    }
  }
}
