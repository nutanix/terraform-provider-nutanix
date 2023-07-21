---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pbr"
sidebar_current: "docs-nutanix-resource-pbr"
description: |-
  Create Policy Based Routing within VPCs .
---

# nutanix_pbr

Provides Nutanix resource to create Policy Based Routing inside VPCs.

## Example Usage

### pbr creation with vpc name with any source or destination or protocol with permit action

```hcl
resource "nutanix_pbr" "pbr" {
  name = "test-policy-1"
  priority = 123
  protocol_type = "ALL"
  action = "PERMIT"
  vpc_name = "test123"
  source{
    address_type = "ALL"
  }
  destination{
     address_type = "ALL"
  }
}
```

### pbr creation with vpc uuid with source external and destination network with reroute action and  tcp port rangelist

```hcl
resource "nutanix_pbr" "pbr2" {
    name = "test2"
    priority = 132
 
 
    vpc_reference_uuid = <vpc_reference_uuid>
    source{
        address_type = "INTERNET"
    }
    destination{
        subnet_ip=  "1.2.2.0"
        prefix_length=  24
    }

    protocol_type = "TCP"
    protocol_parameters{
        tcp{
            source_port_range_list{
                end_port  = 50
                start_port = 50
            }
            destination_port_range_list{
                end_port  = 40
                start_port = 40
            }
        }
    }

    action = "REROUTE"
    service_ip_list = ["10.x.x.xx"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) name of policy
* `priority` - (Required) priority of policy
* `protocol_type` - (Required) Protocol type of policy based routing. Must be one of {TCP, UDP, ICMP, PROTOCOL_NUMBER, ALL} .
* `action` - (Required) Routing policy action. Must be one of {DENY, PERMIT, REROUTE} .
* `service_ip_list` - (Optional) IP addresses of network services. This field is valid only when action is REROUTE.
* `vpc_reference_uuid` - (Required) The reference to a vpc . Should not be used with {vpc_name} .
* `vpc_name` - (Required) The reference to a vpc. Should not be used with {vpc_reference_uuid}
* `is_bidirectional` - (Optional) Additionally create Policy in reverse direction. Should be used with {TCP, UDP with start and end port ranges and ICMP with icmp code and type}. Supported with 2022.x.

## source
source address of an IP packet. This could be either an ip prefix or an address_type . 

* `address` - (Optional) address type of source. Should be one of {INTERNET, ALL}.
* `subnet_ip` - (Optional) IP subnet provided as an address.
* `prefix_length` - (Optional) prefix length of provided subnet. 

## destination
destination address of an IP packet. This could be either an ip prefix or an address_type . 

* `address` - (Optional) address type of source. Should be one of {INTERNET, ALL}.
* `subnet_ip` - (Optional) IP subnet provided as an address.
* `prefix_length` - (Optional) prefix length of provided subnet. 

## protocol_parameters
Routing policy IP protocol parameters

* `tcp` - (Optional) TCP parameters in routing policy
* `udp` - (Optional) UDP parameters in routing policy
* `icmp` - (Optional) ICMP parameters in routing policy.
* `protocol_number` - (Optional) Protocol number in routing policy

## tcp, udp

* `source_port_range` - (Required) Range of TCP/UDP ports.
* `destination_port_range` - (Required) Range of TCP/UDP ports.

## source_port_range, destination_port_range

* `start_port` - (Required) start port number
* `end_port` - (Required) end port number


## Attributes Reference

The following attributes are exported:

* `metadata` - The routing policies kind metadata.
* `api_version` - The version of the API.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: - subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

See detailed information in [Nutanix Policy Based Routing](https://www.nutanix.dev/api_references/prism-central-v3/#/18a0dab82342c-create-a-new-routing-policy)