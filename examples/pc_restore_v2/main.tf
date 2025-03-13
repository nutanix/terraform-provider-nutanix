terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1.0"
    }
  }
}

#Define provider configuration to connect to Nutanix PC

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Define another alias provider configuration to connect to Nutanix PE
provider "nutanix" {
  alias    = "pe"
  username = var.nutanix_pe_username
  password = var.nutanix_pe_endpoint
  endpoint = var.nutanix_pe_endpoint # PE endpoint
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
resource "nutanix_pc_restore_source_v2" "cluster-location" {
  provider = nutanix.pe
  location {
    cluster_location {
      config {
        # clusterExtID, get it from the PC
        ext_id = local.clusterExtId
      }
    }
  }
}

# Get the restorable pcs using the restore source extID
data "nutanix_restorable_pcs_v2" "restorable-pcs" {
  provider              = nutanix.pe
  restore_source_ext_id = nutanix_pc_restore_source_v2.cluster-location.ext_id
}

# Take the first restorable pc ext_id since we are restoring only one pc
# if you are having multiple pcs, you can loop through the restorable pcs and take the desired pc ext_id
locals {
  restorablePcExtId = data.nutanix_restorable_pcs_v2.restorable-pcs.restorable_pcs.0.ext_id
}

# Get the restore points of the pc using the restorable pc ext_id and restore source ext_id
data "nutanix_pc_restore_points_v2" "restore-points" {
  provider                         = nutanix.pe
  restorable_domain_manager_ext_id = local.restorablePcExtId
  restore_source_ext_id            = nutanix_pc_restore_source_v2.cluster-location.id
}

# Take the first restore point ext_id since we are define backup target of type cluster location
# so we will have only one restore point
# for backup of type Object store, we will have multiple restore points
# if you are having multiple restore points, you can loop through the restore points and take the desired restore point ext_id
data "nutanix_pc_restore_point_v2" "restore-point" {
  provider                         = nutanix.pe
  restore_source_ext_id            = nutanix_pc_restore_source_v2.cluster-location.id
  restorable_domain_manager_ext_id = local.restorablePcExtId
  ext_id                           = data.nutanix_pc_restore_points_v2.restore-points.restore_points[0].ext_id
}

# Capture the restore point and store it in a local variable.
locals {
  restorePoint = data.nutanix_pc_restore_point_v2.restore-point
}


# define the restore pc resource
# you can get these values from the data source nutanix_pc_v2, this data source is on PC provider
resource "nutanix_pc_restore_v2" "restore-pc" {
  provider = nutanix.pe
  timeouts {
    create = "120m"
  }
  ext_id                           = local.restorePoint.ext_id
  restore_source_ext_id            = nutanix_pc_restore_source_v2.cluster-location.id
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
}
