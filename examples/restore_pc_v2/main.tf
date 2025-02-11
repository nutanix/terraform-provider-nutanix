terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1"
    }
  }
}

#defining nutanix configuration for pe
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint # PE Endpoint
  port     = 9440
  insecure = true
}


resource "nutanix_restore_pc_v2" "test" {
  timeouts {
    create = "120m"
  }
  ext_id                           = "3f4017e8-85e2-3cf5-8dfe-2dbfdfe546bf"
  restore_source_ext_id            = "642c2af7-a38a-4085-99a6-5baedd02cdb5"
  restorable_domain_manager_ext_id = "36526df2-08ac-49c4-bb4f-13f769c2b7ed"
  domain_manager {
    config {
      should_enable_lockdown_mode = false
      build_info {
        version = "pc.2024.3"
      }
      name = "tf-test-deploy-pc-4681093599007582619"
      size = "STARTER"
      resource_config {
        container_ext_ids    = ["6a6da162-bd6c-418d-ad16-e99c5a6c4fb2"]
        data_disk_size_bytes = 289910292480
        memory_size_bytes    = 19327352832
        num_vcpus            = 4
      }
    }
    network {
      external_address {
        ipv4 {
          value = ""
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
          value = "3.centos.pool.ntp.org"
        }
      }
      ntp_servers {
        fqdn {
          value = "2.centos.pool.ntp.org"
        }
      }
      ntp_servers {
        fqdn {
          value = "1.centos.pool.ntp.org"
        }
      }
      ntp_servers {
        fqdn {
          value = "0.centos.pool.ntp.org"
        }
      }

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
    }
  }
  # after restore pc you need to reset the password 5 times before setting the desired password, you will need to use local-exec provisioner to do that
  provisioner "local-exec" {
    command    = "sshpass -p 'nutanix/4u' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null nutanix@10.97.64.91 '/home/nutanix/prism/cli/ncli user reset-password user-name=admin password=u.B.8.@.D.@.R ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=y.O.7.@.d.d.y ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=a.Z.5.$.5 ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=z.J.1.#.g ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=Nutanix.123'"
    on_failure = continue
  }
}