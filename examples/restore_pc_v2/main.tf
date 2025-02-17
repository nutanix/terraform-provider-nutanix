# define another alias for the provider,  PE
provider "nutanix" {
  alias    = "remote"
  username = var.nutanix_pe_username
  password = var.nutanix_pe_password
  endpoint = var.nutanix_pe_endpoint # PE endpoint
  insecure = true
  port     = 9440
}


# Create a restore source, before make sure to get the cluster ext_id from PC and create backup target
# wait until backup target is synced, you can check the last_sync_time from the backup target data source
# power off the PC VM before restore it
resource "nutanix_restore_source_v2" "cluster-location" {
  provider = nutanix.remote
  location {
    cluster_location {
      config {
        # clusterExtID, get it from the PC
        ext_id = "00062d4c-42b3-20b8-185b-ac1f6b6f97e2"
      }
    }
  }
}

data "nutanix_restorable_pcs_v2" "restorable-pcs" {
  provider              = nutanix.remote
  restore_source_ext_id = nutanix_restore_source_v2.cluster-location.ext_id
}

locals {
  restorablePcExtId = data.nutanix_restorable_pcs_v2.restorable-pcs.restorable_pcs.0.ext_id
}

data "nutanix_restore_points_v2" "restore-points" {
  provider                         = nutanix.remote
  restorable_domain_manager_ext_id = local.restorablePcExtId
  restore_source_ext_id            = nutanix_restore_source_v2.cluster-location.id
}

locals {
  restorePointId = data.nutanix_restore_points_v2.restore-points.restore_points[0].ext_id
}


# define the restore pc resource
# you can get these values from the data source nutanix_pc_v2, this data source is on PC provider
resource "nutanix_restore_pc_v2" "test" {
  provider = nutanix.remote
  timeouts {
    create = "120m"
  }
  ext_id                           = local.restorePointId
  restore_source_ext_id            = nutanix_restore_source_v2.cluster-location.id
  restorable_domain_manager_ext_id = local.restorablePcExtId
  domain_manager {
    config {
      should_enable_lockdown_mode = false
      build_info {
        version = "pc.2024.3"
      }
      name = "PC_10.44.76.17"
      size = "SMALL"
      resource_config {
        container_ext_ids = ["3cec211e-3b16-4832-9b93-f299bcc328fc"]
        data_disk_size_bytes = 536870912000
        memory_size_bytes    = 39728447488
        num_vcpus            = 10
      }
    }
    network {
      external_address {
        ipv4 {
          value = "10.44.76.17"
        }
      }

      # name servers
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

      # ntp servers
      ntp_servers {
        fqdn {
          value = "0.centos.pool.ntp.org"
        }
      }

      ntp_servers {
        fqdn {
          value = "1.centos.pool.ntp.org"
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

      external_networks {
        network_ext_id = "f48ef43f-5d52-4dfd-8123-90ff672c4b1d"
        default_gateway {
          ipv4 {
            value = "10.44.76.1"
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
              value = "10.44.76.16"
            }
          }
          end {
            ipv4 {
              value = "10.44.76.16"
            }
          }
        }
      }
    }
  }

  # after restore pc, you need to reset the password of the admin user
  provisioner "local-exec" {
    command    = "sshpass -p 'nutanix/4u' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null nutanix@10.44.76.16 '/home/nutanix/prism/cli/ncli user reset-password user-name=admin password=o.P.5.#.s.U.Z.f ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=n.L.9.@.P.Y ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=g.B.1.$.U.$.2.@ ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=r.B.7.$.V.9.W ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=l.H.2.$.2.a.a.P ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=q.F.4.#.u.t ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=n.T.0.#.r ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=s.K.0.$.w ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=o.K.7.@.j ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=Nutanix.123'"
    on_failure = continue
  }
}
