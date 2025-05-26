---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pbrs_v2"
sidebar_current: "docs-nutanix-datasource-pbrs-v2"
description: |-
  Provides a datasource to get the list of Routing Policies.
---

# nutanix_pbrs_v2

Get a list of Routing Policies.

## Example Usage

```hcl
data "nutanix_pbrs_v2" "pbrs"{}

data "nutanix_pbrs_v2" "pbrs-filter"{
  filter = "name eq 'pbr_example'"
}

data "nutanix_pbrs_v2" "pbrs-limit"{
  limit = 10
}

data "nutanix_pbrs_v2" "pbrs-filter-limit"{
  filter = "name eq 'pbr_example'"
  limit = 10
}
```

## Argument Reference

The following arguments are supported:

- `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources. The filter can be applied to the following fields:
  - name
  - policies/policyAction/actionType
  - policies/policyMatch/protocolType
  - policies/policyMatch/source
  - priority
  - vpcExtId
- `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
  - name
  - priority
- `expand`: (Optional) A URL query parameter that allows clients to request related resources when a resource that satisfies a particular request is retrieved. The expand can be applied to the following fields:
  - vpc
- `select`: (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. The select can be applied to the following fields:
  - description
  - extId
  - links
  - metadata
  - name
  - policies
  - priority
  - tenantId
  - vpc
  - vpcExtId

## Attribute Reference

The following attributes are exported:

- `routing_policies`: List all of routing policies.

### Routing Policies

The `routing_policies` object contains the following attributes:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `metadata`: Metadata associated with this resource.
- `name`: Name of the routing policy.
- `description`: A description of the routing policy.
- `priority`: Priority of the routing policy.
- `policies`: Routing Policies
- `vpc_ext_id`: ExtId of the VPC extId to which the routing policy belongs.
- `vpc`: VPC name for projections

### Links

The `links` attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### metadata

The `metadata` object contains the following attributes:

- `owner_reference_id` : A globally unique identifier that represents the owner of this resource.
- `owner_user_name` : The userName of the owner of this resource.
- `project_reference_id` : A globally unique identifier that represents the project this resource belongs to.
- `project_name` : The name of the project this resource belongs to.
- `category_ids` : A list of globally unique identifiers that represent all the categories the resource is associated with.

### policies

- `policy_match`: Match condition for the traffic that is entering the VPC.
- `policy_action`: The action to be taken on the traffic matching the routing policy.
- `is_bidirectional`: If True, policies in the reverse direction will be installed with the same action but source and destination will be swapped.

### policy_match

- `source`: Address Type like "EXTERNAL" or "ANY".
- `destination`: Address Type like "EXTERNAL" or "ANY".
- `protocol_type`: Routing Policy IP protocol type.
- `protocol_parameters`: Protocol Params Object.

### policy_match.source, policy_match.destination

- `address_type`: Address Type like "EXTERNAL" or "ANY".
- `subnet_prefix`: Subnet Prefix

### subnet_prefix

- `ip`: IP of address
- `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.

### protocol_parameters

- `layer_four_protocol_object`: Layer Four Protocol Object.
- `icmp_object`: ICMP object
- `protocol_number_object`: Protocol Number Object.

### layer_four_protocol_object

- `source_port_ranges`: Start and end port ranges object.
- `destination_port_ranges`: Start and end port ranges object.

### icmp_object

- `icmp_type`: icmp type
- `icmp_code`: icmp code

### protocol_number_object

- `protocol_number`: protocol number

### policy_action

- `action_type`: Routing policy action type.
- `reroute_params`: Routing policy Reroute params.
- `nexthop_ip_address`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

### reroute_params

- `service_ip`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `reroute_fallback_action`: Type of fallback action in reroute case when service VM is down.
- `ingress_service_ip`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `egress_service_ip`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

### ipv4,ipv6 Configuration format

- `value`: ip value
- `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.

See detailed information in [Nutanix List Routing Policies v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/RoutingPolicies/operation/listRoutingPolicies).
