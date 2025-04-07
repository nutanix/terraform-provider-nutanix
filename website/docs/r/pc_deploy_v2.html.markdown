---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pc_deploy_v2 "
sidebar_current: "docs-nutanix-resource-pc-deploy-v2"
description: |-
  This operation Deploys a Prism Central using the provided details. Prism Central Size, Network Config are mandatory fields to deploy Prism Central. The response from this endpoint contains the URL in the task object location header that can be used to track the request status.
---

# nutanix_pc_deploy_v2


> - The Pc Deploy V2 resource is an action-only resource that supports creating actions. The update and delete operations have no effect. To run it again, destroy and reapply the resource.
> - We need to increase the timeout for deploying the PC resource because the deployment takes longer than the default timeout allows for the operation to complete.

Deploys a Prism Central using the provided details. Prism Central Size, Network Config are mandatory fields to deploy Prism Central. The response from this endpoint contains the URL in the task object location header that can be used to track the request status.

## Example Usage

```hcl

resource "nutanix_pc_deploy_v2" "example"{
  # we need to increase the timeout for deploying the PC resource because the deployment takes longer than the default timeout allows for the operation to complete.
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
```

## Argument Reference

The following arguments are supported:

- `config`: -(Required) Domain manager (Prism Central) cluster configuration details.
- `network`: -(Required) Domain manager (Prism Central) network configuration details.
- `should_enable_high_availability`: -(Optional) This configuration enables Prism Central to be deployed in scale-out mode. Default is `false`.

### Config

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

- `version`: -(Required) Software version.

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

See detailed information in [Nutanix Deploy PC V4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/createDomainManager).
