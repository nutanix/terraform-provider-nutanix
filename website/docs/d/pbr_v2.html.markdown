---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pbr_v2"
sidebar_current: "docs-nutanix-datasource-pbr-v2"
description: |-
   Provides a datasource to get a single Routing Policy corresponding to the extId.
---

# nutanix_pbr_v2

Get a single Routing Policy corresponding to the extId.

## Example Usage

```hcl
    data "nutanix_pbr_v2" "pbr"{
        ext_id = "<pbr_uuid>"
    }
```

## Argument Reference

The following arguments are supported:

* `pbr_uuid`: (Required) pbr UUID


## Attribute Reference

The following attributes are exported:

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `metadata`: Metadata associated with this resource.
* `name`: Name of the routing policy.
* `description`: A description of the routing policy.
* `priority`: Priority of the routing policy.
* `policies`: Routing Policies
* `vpc_ext_id`: ExtId of the VPC extId to which the routing policy belongs.
* `vpc`: VPC name for projections


### policies

* `policy_match`: Match condition for the traffic that is entering the VPC.
* `policy_action`: The action to be taken on the traffic matching the routing policy.
* `is_bidirectional`: If True, policies in the reverse direction will be installed with the same action but source and destination will be swapped.


### policy_match
* `source`: Address Type like "EXTERNAL" or "ANY".
* `destination`: Address Type like "EXTERNAL" or "ANY".
* `protocol_type`: Routing Policy IP protocol type.
* `protocol_parameters`: Protocol Params Object.

### policy_match.source, policy_match.destination
* `address_type`: Address Type like "EXTERNAL" or "ANY".
* `subnet_prefix`: Subnet Prefix

### subnet_prefix
* `ip`: IP of address
* `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.


### protocol_parameters
* `layer_four_protocol_object`: Layer Four Protocol Object. 
* `icmp_object`: ICMP object
* `protocol_number_object`: Protocol Number Object. 

### layer_four_protocol_object
* `source_port_ranges`: Start and end port ranges object.
* `destination_port_ranges`: Start and end port ranges object.

### icmp_object
* `icmp_type`: icmp type
* `icmp_code`: icmp code

### protocol_number_object
* `protocol_number`: protocol number


### policy_action
* `action_type`: Routing policy action type.
* `reroute_params`: Routing policy Reroute params.
* `nexthop_ip_address`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

### reroute_params
* `service_ip`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `reroute_fallback_action`: Type of fallback action in reroute case when service VM is down.
* `ingress_service_ip`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `egress_service_ip`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.


### ipv4,ipv6 Configuration format
* `value`: ip value
* `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.


See detailed information in [Nutanix Routing Policy v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0.b1).