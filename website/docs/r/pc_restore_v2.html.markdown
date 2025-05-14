---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pc_restore_v2"
sidebar_current: "docs-nutanix-resource-pc-restore-v2"
description: |-
  The restore domain manager is a task-driven operation to restore a domain manager from a cluster or object store backup location based on the selected restore point.
---

# nutanix_pc_restore_v2


> - The Pc Restore V2 resource is an action-only resource that supports creating actions. The update and delete operations have no effect. To run it again, destroy and reapply the resource.
> -  We need to increase the timeout for restoring the PC, because the restore pc takes longer than the default timeout allows for the operation to complete.


The restore domain manager is a task-driven operation to restore a domain manager from a cluster or object store backup location based on the selected restore point.

## Example Usage

```hcl

# define another alias for the provider,  PE
provider "nutanix" {
  alias    = "pe"
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
resource "nutanix_pc_restore_source_v2" "cluster-location" {
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
  restore_source_ext_id = nutanix_pc_restore_source_v2.cluster-location.ext_id
}

locals {
  restorablePcExtId = data.nutanix_restorable_pcs_v2.restorable-pcs.restorable_pcs.0.ext_id
}

data "nutanix_pc_restore_points_v2" "restore-points" {
  provider                         = nutanix.remote
  restorable_domain_manager_ext_id = local.restorablePcExtId
  restore_source_ext_id            = nutanix_pc_restore_source_v2.cluster-location.id
}

data "nutanix_pc_restore_point_v2" "restore-point" {
  provider = nutanix.remote
  restore_source_ext_id = nutanix_pc_restore_source_v2.cluster-location.id
  restorable_domain_manager_ext_id = local.restorablePcExtId
  ext_id   = data.nutanix_pc_restore_points_v2.restore-points.restore_points[0].ext_id
}

locals {
  restorePoint = data.nutanix_pc_restore_point_v2.restore-point
}


# define the restore pc resource
# you can get these values from the data source nutanix_pc_v2, this data source is on PC provider
resource "nutanix_pc_restore_v2" "test" {
  provider = nutanix.remote
  # we need to increase the timeout for restoring the PC, because the restore pc takes longer than the default timeout allows for the operation to complete
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

  # after restore pc, you need to reset the password of the admin user
  provisioner "local-exec" {
    command    = "sshpass -p 'nutanix/4u' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null nutanix@10.44.76.16 '/home/nutanix/prism/cli/ncli user reset-password user-name=admin password=o.P.5.#.s.U.Z.f ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=n.L.9.@.P.Y ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=g.B.1.$.U.$.2.@ ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=r.B.7.$.V.9.W ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=l.H.2.$.2.a.a.P ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=q.F.4.#.u.t ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=n.T.0.#.r ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=s.K.0.$.w ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=o.K.7.@.j ; /home/nutanix/prism/cli/ncli user reset-password user-name=admin password=Nutanix.123'"
    on_failure = continue
  }
}

```

## Argument Reference

The following arguments are supported:

- `restore_source_ext_id`: -(Required) A unique identifier obtained from the restore source API that corresponds to the details provided for the restore source.
- `restorable_domain_manager_ext_id`: -(Required) A unique identifier for the domain manager.
- `ext_id`: -(Required) Restore point ID for the backup created in cluster/object store.
- `domain_manager`: -(Required) Domain manager (Prism Central) details.

### Domain Manager

The location argument supports the following:

- `config`: -(Required) Domain manager (Prism Central) cluster configuration details.
- `network`: -(Required) Domain manager (Prism Central) network configuration details.
- `should_enable_high_availability`: -(Optional) This configuration enables Prism Central to be deployed in scale-out mode. Default is `false`.

#### Config

The config argument supports the following:

- `should_enable_lockdown_mode`: -(Optional) A boolean value indicating whether to enable lockdown mode for a cluster.
- `build_info`: -(Required) Currently representing the build information to be used for the cluster creation.
- `name`: -(Required) Name of the domain manager (Prism Central).
- `size`: - (Required) Domain manager (Prism Central) size is an enumeration of starter, small, large, or extra large starter values. The allowed values are:
  - `SMALL` : Domain manager (Prism Central) of size small.
  - `LARGE` : Domain manager (Prism Central) of size large.
  - `EXTRALARGE` : Domain manager (Prism Central) of size extra large.
  - `STARTER` : Domain manager (Prism Central) of size starter.
- `bootstrap_config`: - (Optional) Bootstrap configuration details for the domain manager (Prism Central).
- `credentials`: - (Optional) The credentials consist of a username and password for a particular user like admin. Users can pass the credentials of admin users currently which will be configured in the create domain manager operation.
- `resource_config`: -(Optional) This configuration is used to provide the resource-related details like container external identifiers, number of VCPUs, memory size, data disk size of the domain manager (Prism Central). In the case of a multi-node setup, the sum of resources like number of VCPUs, memory size and data disk size are provided.

#### Build Info

The `build_info` argument supports the following:

- `version`: -(Optional) Software version.

#### Bootstrap Config

The `bootstrap_config` argument supports the following:

- `cloud_init_config`: -(Optional) Cloud-init configuration for the domain manager (Prism Central) cluster.
- `environment_info`: -(Optional) Environment information for the domain manager (Prism Central) cluster.

##### Cloud Init Config

The `cloud_init_config` argument supports the following:

- `datasource_type`: -(Optional) Type of datasource. Default: CONFIG_DRIVE_V2
- `metadata`: -(Optional)The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded. Default value is 'CONFIG_DRIVE_V2'.
- `cloud_init_script`: -(Optional) The script to use for cloud-init.
- `cloud_init_script.user_data`: -(Optional) user data object
- `cloud_init_script.custom_keys`: -(Optional) The list of the individual KeyValuePair elements.

##### Environment Info

The `environment_info` argument supports the following:

- `type`: -(Optional) Enums denoting the environment type of the PC, that is, on-prem PC or cloud PC.
  Following are the supported entity types:
  - `ONPREM` : On-prem environment.
  - `NTNX_CLOUD` : Nutanix cloud environment.
- `provider_type`: -(Optional) Enums denoting the provider type of the PC, that is, AHV or ESXi.
  Following are the supported provider types:
  - `VSPHERE` : Vsphere cloud provider.
  - `AZURE` : Azure cloud provider.
  - `NTNX` : Nutanix cloud provider.
  - `GCP` : GCP cloud provider.
  - `AWS` : AWS cloud provider.
- `provisioning_type`: -(Optional) Enums denoting the instance type of the cloud PC. It indicates whether the PC is created on bare-metal or on a cloud-provisioned VM. Hence, it supports two possible values:
  - `NTNX` : Nutanix instance.
  - `NATIVE` : Native instance.

#### Credentials

The `credentials` argument supports the following:

- `username`: -(Required) Username required for the basic auth scheme. As per RFC 2617 usernames might be case sensitive.
- `password`: -(Required) Password required for the basic auth scheme.

#### Resource Config

The `resource_config` argument supports the following:

- `container_ext_ids`: -(Optional) The external identifier of the container that will be used to create the domain manager (Prism Central) cluster.

### Network

the `network` argument supports the following:

- `external_address`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `name_servers`: -(Required) List of name servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
- `ntp_servers`: -(Required) List of NTP servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
- `internal_networks`: -(Required) This configuration is used to internally manage Prism Central network.
- `external_networks`: -(Required) This configuration is used to manage Prism Central.

#### External Address

The `external_address` argument supports the following:

- `ipv4`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv6`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv6 format.

#### Name Servers, NTP Servers

The `name_servers` and `ntp_servers` arguments support the following:

- `ipv4`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv6`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv6 format.
- `fqdn`: -(Optional) A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

#### Internal Networks

The `internal_networks` and `external_networks` arguments support the following:

- `default_gateway`: -(Required) An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
- `subnet_mask`: -(Required) An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
- `ip_ranges`: -(Required) Range of IPs used for Prism Central network setup.

#### External Networks

The `external_networks` argument supports the following:

- `default_gateway`: -(Required) An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
- `subnet_mask`: -(Required) An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
- `ip_ranges`: -(Required) Range of IPs used for Prism Central network setup.
- `network_ext_id`: -(Required) The network external identifier to which Domain Manager (Prism Central) is to be deployed or is already configured.

#### Default Gateway, Subnet Mask

The `default_gateway`and `subnet_mask` arguments support the following:

- `ipv4`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv6`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv6 format.
- `fqdn`: -(Optional) A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

#### IP Ranges

The `ip_ranges` argument supports the following:

- `begin`: -(Optional) The beginning IP address of the range.
- `end`: -(Optional) The ending IP address of the range.

#### begin, end

The `begin` and `end` arguments support the following:

- `ipv4`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv6`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv6 format.

#### IpV4, IpV6

The `ipv4` and `ipv6` arguments support the following:

- `value`: -(Required) The IPv4/IPv6 address of the host.
- `prefix_length`: -(Optional) The prefix length of the network to which this host IPv4/IPv6 address belongs.

#### FQDN

The `fqdn` argument supports the following:

- `value`: -(Optional) The fully qualified domain name of the host.

See detailed information in [Nutanix Restore PC V4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/restore).
