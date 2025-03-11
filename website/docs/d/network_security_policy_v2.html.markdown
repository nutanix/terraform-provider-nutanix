---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_policy_v2"
sidebar_current: "docs-nutanix-datasource-network_security_policy_v2"
description: |-
  Get a Network Security Policy
---

# nutanix_network_security_policy_v2

Get a Network Security Policy by ExtID

### Example

```hcl

data "nutanix_network_security_policy_v2" "get-ns-policy"{
    ext_id = "0d717fa1-21da-4ccc-a719-92d51489c0f9"
}

```

## Argument Reference

The following arguments are supported:

- `ext_id`: (Required) Network security policy UUID.

## Attribute Reference

The following attributes are exported:

- `name`: Name of the Flow Network Security Policy.
- `type`: Defines the type of rules that can be used in a policy.
- `description`: A user defined annotation for a policy.
- `state`: Whether the policy is applied or monitored; can be omitted or set null to save the policy without applying or monitoring it.
- `rules`: A list of rules that form a policy. For isolation policies, use isolation rules; for application or quarantine policies, use application rules.
- `is_ipv6_traffic_allowed`: If Ipv6 Traffic is allowed.
- `is_hitlog_enabled`: If Hitlog is enabled.
- `scope`: Defines the scope of the policy. Currently, only ALL_VLAN and VPC_LIST are supported. If scope is not provided, the default is set based on whether vpcReferences field is provided or not.
- `vpc_reference`: A list of external ids for VPCs, used only when the scope of policy is a list of VPCs.
- `secured_groups`: Uuids of the secured groups in the NSP.
- `last_update_time`: last updated time
- `creation_time`: creation time of NSP
- `is_system_defined`: Is system defined NSP
- `created_by`: created by.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

### rules

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `description`: A user defined annotation for a rule.
- `type`: The type for a rule - the value chosen here restricts which specification can be chosen.
- `spec`: Spec for rules.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

### spec

- `two_env_isolation_rule_spec`: Two Environment Isolation Rule Spec.
- `application_rule_spec`: Application Rule Spec.
- `intra_entity_group_rule_spec`: Intra entity group Rule Spec
- `multi_env_isolation_rule_spec`: Multi Environment Isolation Rule Spec.

### two_env_isolation_rule_spec

- `first_isolation_group`: Denotes the first group of category uuids that will be used in an isolation policy.
- `second_isolation_group`: Denotes the second group of category uuids that will be used in an isolation policy.

### application_rule_spec

- `secured_group_category_references`: A set of network endpoints which is protected by a Network Security Policy and defined as a list of categories.
- `src_allow_spec`: A specification to how allow mode traffic should be applied, either ALL or NONE.
- `dest_allow_spec`: A specification to how allow mode traffic should be applied, either ALL or NONE.
- `src_category_references`: List of categories that define a set of network endpoints as inbound.
- `dest_category_references`: List of categories that define a set of network endpoints as outbound.
- `src_subnet`: source subnet value
- `dest_subnet`: destination subnet value
- `src_address_group_references`: A list of address group references.
- `dest_address_group_references`: A list of address group references.
- `service_group_references`: A list of service group references.
- `is_all_protocol_allowed`: Denotes if rule allows traffic for all protocol.
- `tcp_services`: tcp services
- `udp_services`: udp services
- `icmp_services`: icmp services
- `network_function_chain_reference`: A reference to the network function chain in the rule.

### intra_entity_group_rule_spec

- `secured_group_action`: List of secured group action.
- `secured_group_category_references`: A specification to whether traffic between intra secured group entities should be allowed or denied.

### multi_env_isolation_rule_spec

- `spec`: Multi Environment Isolation Rule Spec.

#### spec

- `all_to_all_isolation_group`: all to all isolation group

#### all_to_all_isolation_group

- `isolation_group`: Denotes the list of secured groups that will be used in All to All mutual isolation.

#### isolation_groups

- `group_category_references`: External identifiers of categories belonging to the isolation group.

### tcp_services, tcp_services

- `start_port`: start port
- `end_port`: end port

### icmp_services

- `is_all_allowed`: Set this field to true if both Type and Code is ANY.
- `type`: Icmp service Type. Ignore this field if Type has to be ANY.
- `code`: Icmp service Code. Ignore this field if Code has to be ANY.

### Links

The `links` attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix Security Policy v4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.0).
