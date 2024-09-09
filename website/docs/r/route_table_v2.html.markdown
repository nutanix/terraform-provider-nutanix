---
layout: "nutanix"
page_title: "NUTANIX: nutanix_routes_table_v4"
sidebar_current: "docs-nutanix-resource-routes-table-v4"
description: |-
  Update route table.
---

# nutanix_routes_table_v4

Provides Nutanix resource to update route table


## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) Route table UUID
* `vpc_reference`: (Optional) VPC
* `external_routing_domain_reference`: (Optional) External routing domain associated with this route table
* `static_routes`: (Required) Static routes


### static_routes
* `destination`: (Required) Destination IP Subnet Configuration.
* `next_hop_type`: (Required) Next hop type. Acceptable values are "INTERNAL_SUBNET", "DIRECT_CONNECT_VIF", "VPN_CONNECTION", "IP_ADDRESS", "EXTERNAL_SUBNET"
* `next_hop_reference`: (Optional) The reference to a link, such as a VPN connection or a subnet.
* `next_hop_ip_address`: (Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.


### destination
* `ipv4`: (Optional) IPv4 Subnet Object
* `ipv4.ip`: (Required) IP address format
* `ipv4.prefix_length`: (Required) The prefix length of the network to which this host IPv4 address belongs.

* `ipv6`: (Optional) IPv6 Subnet Object
* `ipv6.ip`: (Required) IP address format
* `ipv6.prefix_length`: (Required) The prefix length of the network to which this host IPv6 address belongs.


### next_hop_ip_address
* `ipv4`: (Optional) IPv4 Address
* `ipv6`: (Optional) IPv6 Address


### IPv4/IPv6 Address
* `value`: (Optional) value of IP address
* `prefix_length`: (Optional) The prefix length of the network to which this host IPv4/IPv6 address belongs.


See detailed information in [Nutanix Route Table v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0.b1).