---
layout: "nutanix"
page_title: "NUTANIX: nutanix_subnet_v2"
sidebar_current: "docs-nutanix-datasource-subnet-v2"
description: |-
  This operation retrieves a subnet with the specified UUID.
---

# nutanix_subnet_v2

Get a subnet with the specified UUID.

### Example

```hcl
data "nutanix_subnet_v2" "example"{
  ext_id = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: - (Required) The UUID of the subnet.

## Attributes Reference

The following attributes are exported:

- `name`: Name of the subnet.
- `description`: Description of the subnet.
- `subnet_type`: Type of subnet.
- `network_id`: or VLAN subnet, this field represents VLAN Id, valid range is from 0 to 4095; For overlay subnet, this field represents 24-bit VNI, this field is read-only.
- `dhcp_options`: List of DHCP options to be configured.
- `ip_config`: IP configuration for the subnet.
- `cluster_reference`: UUID of the cluster this subnet belongs to.
- `virtual_switch_reference`: UUID of the virtual switch this subnet belongs to (type VLAN only).
- `vpc_reference`: UUID of Virtual Private Cloud this subnet belongs to (type Overlay only).
- `is_nat_enabled`: Indicates whether NAT must be enabled for VPCs attached to the subnet. This is supported only for external subnets. NAT is enabled by default on external subnets.
- `is_external`: Indicates whether the subnet is used for external connectivity.
- `reserved_ip_addresses`: List of IPs that are excluded while allocating IP addresses to VM ports.
- `dynamic_ip_addresses`: List of IPs, which are a subset from the reserved IP address list, that must be advertised to the SDN gateway.
- `network_function_chain_reference`: UUID of the Network function chain entity that this subnet belongs to (type VLAN only).
- `bridge_name`: Name of the bridge on the host for the subnet.
- `is_advanced_networking`: Indicates whether the subnet is used for advanced networking.
- `cluster_name`: Cluster Name
- `hypervisor_type`: Hypervisor Type
- `virtual_switch`: Schema to configure a virtual switch
- `vpc`: Networking common base object
- `ip_prefix`: IP Prefix in CIDR format.
- `ip_usage`: IP usage statistics.
- `migration_state`: Migration state of the subnet. This field is read-only.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

### dhcp_options

- `domain_name_servers`: List of Domain Name Server addresses.
- `domain_name`: The DNS domain name of the client.
- `search_domains`: The DNS domain search list.
- `tftp_server_name`: TFTP server name
- `boot_file_name`: Boot file name
- `ntp_servers`: List of NTP server addresses

### domain_name_servers, ntp_servers

- `ipv4`: IPv4 Object
- `ipv6`: IPv6 Object

### ip_config

- `ipv4`: IP V4 configuration.
- `ipv6`: IP V6 configuration

### ip_config.ipv4, ip_config.ipv6

- `ip_subnet`: subnet ip
- `default_gateway_ip`: Reference to address configuration
- `dhcp_server_address`: Reference to address configuration
- `pool_list`: Pool of IP addresses from where IPs are allocated.

### ip_subnet

- `ip`: Reference to address configuration
- `prefix_length`: The prefix length of the network to which this host IPv4 address belongs.

### pool_list

- `start_ip`: Reference to address configuration
- `end_ip`: Reference to address configuration

### ip_usage

- `num_macs`: Number of MAC addresses.
- `num_free_ips`: Number of free IPs.
- `num_assigned_ips`: Number of assigned IPs.
- `ip_pool_usages`: IP Pool usages

### ip_pool_usages

- `num_free_ips`: Number of free IPs
- `num_total_ips`: Total number of IPs in this pool.
- `range`: Start/end IP address range.

### ipv4, ipv6 (Reference to address configuration)

- `value`: value of address
- `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.

See detailed information in [Nutanix Get Subnet v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/Subnets/operation/getSubnetById).
