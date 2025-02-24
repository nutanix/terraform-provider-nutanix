
# define another alias for the provider,  PE
provider "nutanix" {
  alias    = "remote"
  username = var.nutanix_remote_username
  password = var.nutanix_remote_password
  endpoint = var.nutanix_remote_endpoint # PE endpoint
  insecure = true
  port     = 9440
}

# Fetch Cluster Ext ID from PC
data "nutanix_clusters_v2" "clusters" {
  provider = nutanix
}
locals {
  domainManagerExtId = data.nutanix_clusters_v2.cls.cluster_entities.0.ext_id
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# Create a restore source, before make sure to get the cluster ext_id from PC and create backup target
# wait until backup target is synced, you can check the last_sync_time from the backup target data source
resource "nutanix_restore_source_v2" "cluster-location" {
  provider = nutanix.remote
  location {
    cluster_location {
      config {
        # clusterExtID, get it from the PC
        ext_id = local.clusterExtId
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

data "nutanix_restore_point_v2" "restore-point" {
  provider = nutanix.remote
  restore_source_ext_id = nutanix_restore_source_v2.cluster-location.id
  restorable_domain_manager_ext_id = local.restorablePcExtId
  ext_id   = data.nutanix_restore_points_v2.restore-points.restore_points[0].ext_id
}

locals {
  restorePoint = data.nutanix_restore_point_v2.restore-point
}


# define the restore pc resource
# you can get these values from the data source nutanix_pc_v2, this data source is on PC provider
resource "nutanix_restore_pc_v2" "test" {
  provider = nutanix.remote
  timeouts {
    create = "120m"
  }
  ext_id                           = local.restorePoint.ext_id
  restore_source_ext_id            = nutanix_restore_source_v2.cluster-location.id
  restorable_domain_manager_ext_id = local.restorablePcExtId

  domain_manager {
    config {
      should_enable_lockdown_mode = local.restorePoint.domain_manager[0].config[0].should_enable_lockdown_mode

      build_info {
        version = local.restorePoint.domain_manager[0].config[0].build_info[0].version
      }

      name = local.restorePoint.domain_manager[0].config[0].name
      size = local.restorePoint.domain_manager[0].config[0].size

      resource_config {
        container_ext_ids    = local.restorePoint.domain_manager[0].config[0].resource_config[0].container_ext_ids
        data_disk_size_bytes = local.restorePoint.domain_manager[0].config[0].resource_config[0].data_disk_size_bytes
        memory_size_bytes    = local.restorePoint.domain_manager[0].config[0].resource_config[0].memory_size_bytes
        num_vcpus            = local.restorePoint.domain_manager[0].config[0].resource_config[0].num_vcpus
      }
    }

    network {
      external_address {
        ipv4 {
          value = local.restorePoint.domain_manager[0].network[0].external_address[0].ipv4[0].value
        }
      }

      # Dynamically create a block for each name server
      dynamic "name_servers" {
        for_each = local.restorePoint.domain_manager[0].network[0].name_servers
        content {
          ipv4 {
            value = name_servers.value.ipv4[0].value
          }
        }
      }

      # Dynamically create a block for each NTP server
      dynamic "ntp_servers" {
        for_each = local.restorePoint.domain_manager[0].network[0].ntp_servers
        content {
          fqdn {
            value = ntp_servers.value.fqdn[0].value
          }
        }
      }

      external_networks {
        network_ext_id = local.restorePoint.domain_manager[0].network[0].external_networks[0].network_ext_id

        default_gateway {
          ipv4 {
            value = local.restorePoint.domain_manager[0].network[0].external_networks[0].default_gateway[0].ipv4[0].value
          }
        }

        subnet_mask {
          ipv4 {
            value = local.restorePoint.domain_manager[0].network[0].external_networks[0].subnet_mask[0].ipv4[0].value
          }
        }

        ip_ranges {
          begin {
            ipv4 {
              value = local.restorePoint.domain_manager[0].network[0].external_networks[0].ip_ranges[0].begin[0].ipv4[0].value
            }
          }
          end {
            ipv4 {
              value = local.restorePoint.domain_manager[0].network[0].external_networks[0].ip_ranges[0].end[0].ipv4[0].value
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
