---
layout: "nutanix"
page_title: "NUTANIX: nutanix_subnet_v2"
sidebar_current: "docs-nutanix-resource-subnet-v2"
description: |-
  This operation submits a request to create a subnet based on the input parameters. A subnet is a block of IP addresses.
---

# nutanix_subnet_v2

Provides a resource to create a subnet based on the input parameters.

## Example

```hcl

#creating subnet with IP pool
resource "nutanix_subnet_v2" "vlan-112" {
  name              = "vlan-112"
  description       = "subnet VLAN 112 managed by Terraform with IP pool"
  cluster_reference = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
  subnet_type       = "VLAN"
  network_id        = 122
  is_external       = true
  ip_config {

    ipv4 {
      ip_subnet {
        ip {
          value = "192.168.0.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "192.168.0.1"
      }
      pool_list {
        start_ip {
          value = "192.168.0.20"
        }
        end_ip {
          value = "192.168.0.30"
        }
      }
    }
  }
}


#creating subnet without IP pool
resource "nutanix_subnet_v2" "vlan-113" {
  name              = "vlan-113"
  description       = "subnet VLAN 113 managed by Terraform"
  cluster_reference = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
  subnet_type       = "VLAN"
  network_id        = 113
}

# creating subnet with IP pool and DHCP options
resource "nutanix_subnet_v2" "van-114" {
  name              = "vlan-114"
  description       = "subnet VLAN 114 managed by Terraform"
  cluster_reference = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
  subnet_type       = "VLAN"
  network_id        = 114
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "192.168.0.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "192.168.0.1"
      }
      pool_list {
        start_ip {
          value = "192.168.0.20"
        }
        end_ip {
          value = "192.168.0.30"
        }
      }
    }
  }

  dhcp_options {
    domain_name_servers {
      ipv4 {
        value = "8.8.8.8"
      }
    }
    search_domains   = ["eng.nutanix.com"]
    domain_name      = "nutanix.com"
    tftp_server_name = "10.5.0.10"
    boot_file_name   = "pxelinux.0"
  }
}

```

## Argument Reference

- `name`: (Required) Name of the subnet.
- `description`: (Optional) Description of the subnet.
- `subnet_type`: (Required) Type of subnet. Acceptables values are "OVERLAY", "VLAN".
- `network_id`: (Optional) For VLAN subnet, this field represents VLAN Id, valid range is from 0 to 4095; For overlay subnet, this field represents 24-bit VNI, this field is read-only.
- `dhcp_options`: (Optional) List of DHCP options to be configured.
- `ip_config`: (Optional) IP configuration for the subnet.
- `cluster_reference`: (Optional) UUID of the cluster this subnet belongs to.
- `virtual_switch_reference`: (Optional) UUID of the virtual switch this subnet belongs to (type VLAN only).
- `vpc_reference`: (Optional) UUID of Virtual Private Cloud this subnet belongs to (type Overlay only).
- `is_nat_enabled`: (Optional) Indicates whether NAT must be enabled for VPCs attached to the subnet. This is supported only for external subnets. NAT is enabled by default on external subnets.
- `is_external`: (Optional) Indicates whether the subnet is used for external connectivity.
- `reserved_ip_addresses`: (Optional) List of IPs that are excluded while allocating IP addresses to VM ports. Reference to address configuration
- `dynamic_ip_addresses`: (Optional) List of IPs, which are a subset from the reserved IP address list, that must be advertised to the SDN gateway.
- `network_function_chain_reference`: (Optional) UUID of the Network function chain entity that this subnet belongs to (type VLAN only).
- `bridge_name`: (Optional) Name of the bridge on the host for the subnet.
- `is_advanced_networking`: (Optional) Indicates whether the subnet is used for advanced networking.
- `cluster_name`: (Optional) Cluster Name
- `hypervisor_type`: (Optional) Hypervisor Type
- `virtual_switch`: (Optional) Schema to configure a virtual switch
- `vpc`: (Optional) Networking common base object
- `ip_prefix`: (Optional) IP Prefix in CIDR format.

### dhcp_options

- `domain_name_servers`: (Optional) List of Domain Name Server addresses.
- `domain_name`: (Optional) The DNS domain name of the client.
- `search_domains`: (Optional) The DNS domain search list.
- `tftp_server_name`: (Optional) TFTP server name
- `boot_file_name`: (Optional) Boot file name
- `ntp_servers`: (Optional) List of NTP server addresses

### domain_name_servers, ntp_servers

- `ipv4`: (Optional) IPv4 Object. Reference to address configuration
- `ipv6`: (Optional) IPv6 Object. Reference to address configuration

### ip_config

- `ipv4`: (Optional) IP V4 configuration.
- `ipv6`: (Optional) IP V6 configuration

### ip_config.ipv4, ip_config.ipv6 (IP V4/v6 configuration)

- `ip_subnet`: (Required) subnet ip
- `default_gateway_ip`: (Optional) Reference to address configuration
- `dhcp_server_address`: (Optional) Reference to address configuration
- `pool_list`: (Optional) Pool of IP addresses from where IPs are allocated.

### ip_subnet

- `ip`: (Required) Reference to address configuration
- `prefix_length`: (Required) The prefix length of the network to which this host IPv4 address belongs.

### pool_list

- `start_ip`: (Required) Reference to address configuration
- `end_ip`: (Required) Reference to address configuration

### virtual_switch

- `name`: (Required) User-visible Virtual Switch name
- `description`: (Optional) Input body to configure a Virtual Switch
- `is_default`: (Optional) Indicates whether it is a default Virtual Switch which cannot be deleted
- `has_deployment_error`: (Optional) When true, the node is not put in maintenance mode during the create/update operation.
- `mtu`: (Optional) MTU
- `bond_mode`: (Required) The types of bond modes
- `clusters`: (Required) Cluster configuration list

### clusters

- `ext_id`: (Required) Reference ExtId for the cluster. This is a required parameter on Prism Element ; and is optional on Prism Central
- `hosts`: (Required) Host configuration array
- `gateway_ip_address`: (Optional) Reference to address configuration

### hosts

- `ext_id`: (Required) Reference to the host
- `host_nics`: (Required) Host NIC array
- `ip_address`: (Optional) Ip Address config.
- `ip_address.ip`: (Required) Reference to address configuration
- `ip_address.prefix_length`: (Required) prefix length of address.

### vpc

- `ext_id`: (Optional) A globally unique identifier of an instance that is suitable for external consumption.
- `name`: (Required) Name of the VPC.
- `description`: (Optional) Description of the VPC.
- `vpc_type`: (Optional) Type of VPC. Acceptables values are "REGULAR" , "TRANSIT".
- `common_dhcp_options`: (Optional) List of DHCP options to be configured.
- `external_subnets`: (Optional) List of external subnets that the VPC is attached to.
- `external_routing_domain_reference`: (Optional) External routing domain associated with this route table
- `externally_routable_prefixes`: (Optional) CIDR blocks from the VPC which can talk externally without performing NAT. This is applicable when connecting to external subnets which have disabled NAT.

### common_dhcp_options

- `domain_name_servers`: (Optional) List of Domain Name Server addresses.
- `domain_name_servers.ipv4`: (Optional) Reference to address configuration
- `domain_name_servers.ipv6`: (Optional) Reference to address configuration

### external_subnets

- `subnet_reference`: (Required) External subnet reference.
- `external_ips`: (Optional) List of IP Addresses used for SNAT, if NAT is enabled on the external subnet. If NAT is not enabled, this specifies the IP address of the VPC port connected to the external gateway.
- `gateway_nodes`: (Optional) List of gateway nodes that can be used for external connectivity.
- `active_gateway_node`: (Optional) Reference of gateway nodes

### external_ips

- `ipv4`: (Optional) Reference to address configuration
- `ipv6`: (Optional) Reference to address configuration

### active_gateway_node

- `node_id`: (Optional) Node id
- `node_ip_address`: (Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `node_ip_address.ipv4`: (Optional) Reference to address configuration
- `node_ip_address.ipv6`: (Optional) Reference to address configuration

### externally_routable_prefixes

- `ipv4`: (Optional) IP v4 subnet
- `ipv4.ip`: (Required) Reference to address configuration
- `ipv4.prefix_length`: (Required) The prefix length of the network.

- `ipv6`: (Optional) IP v6 subnet
- `ipv6.ip`: (Required) Reference to address configuration
- `ipv6.prefix_length`: (Required) The prefix length of the network.

### ipv4, ipv6 (Reference to address configuration)

- `value`: value of address
- `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs. Default value is 32.

## Import

This helps to manage existing entities which are not created through terraform. Subnet can be imported using the `UUID`. (ext_id in v4 API context). eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_subnet_v2" "import_subnet" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_subnet_v2" "fetch_subnet"{}
terraform import nutanix_subnet_v2.import_subnet <UUID>
```

See detailed information in [Nutanix Subnet v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/Subnets/operation/createSubnet).
