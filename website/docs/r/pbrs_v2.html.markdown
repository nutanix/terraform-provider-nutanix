---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pbr_v2"
sidebar_current: "docs-nutanix-resource-pbr-v2"
description: |-
  Create Routing Policy within VPCs .
---

# nutanix_pbr_v2

Create a Routing Policy.


## Example

```hcl

# create PBR with vpc name with any source or destination or protocol with permit action
resource "nutanix_pbr_v2" "any-source-destination"{
  name        = "routing_policy_any_source_destination"
  description = "routing policy with any source and destination"
  vpc_ext_id  = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
  priority    = 11
  policies {
    policy_match {
      source {
        address_type = "ANY"
      }
      destination {
        address_type = "ANY"
      }
      protocol_type = "UDP"
    }
    policy_action {
      action_type = "PERMIT"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name`: (Required) Name of the routing policy.
* `description`: A description of the routing policy.
* `priority`: (Required) Priority of the routing policy.
* `policies`: (Required) Routing Policies.
* `vpc_ext_id`: (Required) ExtId of the VPC extId to which the routing policy belongs.


### policies

* `policy_match`: (Required) Match condition for the traffic that is entering the VPC.
* `policy_action`: (Required) The action to be taken on the traffic matching the routing policy.
* `is_bidirectional`: (Optional) If True, policies in the reverse direction will be installed with the same action but source and destination will be swapped.


### policy_match
* `source`: (Required) Address Type like "EXTERNAL" or "ANY".
* `destination`: (Required) Address Type like "EXTERNAL" or "ANY".
* `protocol_type`: (Required) Routing Policy IP protocol type. Acceptable values are "TCP", "UDP", "PROTOCOL_NUMBER", "ANY", "ICMP" .
* `protocol_parameters`: (Optional) Protocol Params Object.

### policy_match.source, policy_match.destination
* `address_type`: (Required) Address Type. Acceptable values are "SUBNET", "EXTERNAL", "ANY" .
* `subnet_prefix`: (Optional) Subnet Prefix

### subnet_prefix
* `ipv4`: (Optional) IPv4 Object.
* `ipv6`: (Optional) IPv6 Object.

### subnet_prefix.ipv4. subnet_prefix.ipv6
* `ip`: (Required) IP of address
* `prefix_length`: (Optional) The prefix length of the network to which this host IPv4/IPv6 address belongs.


### protocol_parameters
* `layer_four_protocol_object`: (Optional) Layer Four Protocol Object.
* `icmp_object`: (Optional) ICMP object
* `protocol_number_object`: (Optional) Protocol Number Object.

### layer_four_protocol_object
* `source_port_ranges`: (Optional) Start and end port ranges object.
* `destination_port_ranges`: (Optional) Start and end port ranges object.

### source_port_ranges, destination_port_ranges
* `start_port`: (Required) Start Port.
* `end_port`: (Required) End Port.


### icmp_object
* `icmp_type`: (Optional) icmp type
* `icmp_code`: (Optional) icmp code

### protocol_number_object
* `protocol_number`: (Required) protocol number


### policy_action
* `action_type`: (Required) Routing policy action type.
* `reroute_params`: (Optional) Routing policy Reroute params.
* `nexthop_ip_address`: (Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

### reroute_params
* `service_ip`: (Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `reroute_fallback_action`: (Optional) Type of fallback action in reroute case when service VM is down. Acceptable values are "PASSTHROUGH", "NO_ACTION", "ALLOW", "DENY".
* `ingress_service_ip`: (Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `egress_service_ip`: (Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.


### ipv4,ipv6 Configuration format
* `value`: (Required) ip value
* `prefix_length`: (Optional) The prefix length of the network to which this host IPv4/IPv6 address belongs.



## Attributes Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `metadata`: Metadata associated with this resource.
* `vpc`: VPC name for projections

## Import
This helps to manage existing entities which are not created through terraform. Routing Policy can be imported using the `UUID`.  eg,

`
terraform import nutanix_pbr_v2.pbr_import <UUID>
`

See detailed information in [Nutanix Routing Policy v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0).
