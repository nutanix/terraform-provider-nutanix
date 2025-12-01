---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_policy_v2"
sidebar_current: "docs-nutanix-resource-network-security-policy-v2"
description: |-
  Create a Network Security Policy
---

# nutanix_network_security_policy_v2

Create a Network Security Policy

## Example

```hcl

# Network Security Policy TWO_ENV_ISOLATION Rule
resource "nutanix_network_security_policy_v2" "isolation-nsp" {
  name        = "isolation_policy"
  description = "isolation policy example"
  state       = "SAVE"
  type        = "ISOLATION"
  rules {
    type = "TWO_ENV_ISOLATION"
    spec {
      two_env_isolation_rule_spec {
        first_isolation_group = ["ba250e3e-1db1-4950-917f-a9e2ea35b8e3"]
        second_isolation_group = ["ab520e1d-4950-1db1-917f-a9e2ea35b8e3"]
      }
    }
  }
  is_hitlog_enabled = true
}

```

## Argument Reference

The following arguments are supported:

- `name`: (Required) Name of the Flow Network Security Policy.
- `type`: (Required) Defines the type of rules that can be used in a policy. Acceptable values are "QUARANTINE", "ISOLATION", "APPLICATION".
- `description`: (Optional) A user defined annotation for a policy.
- `state`: (Optional) Whether the policy is applied or monitored; can be omitted or set null to save the policy without applying or monitoring it. Acceptable values are "SAVE", "MONITOR", "ENFORCE".
- `rules`: (Optional) A list of rules that form a policy. For isolation policies, use isolation rules; for application or quarantine policies, use application rules.
- `is_ipv6_traffic_allowed`: (Optional) If Ipv6 Traffic is allowed.
- `is_hitlog_enabled`: (Optional) If Hitlog is enabled.
- `scope`: Defines the scope of the policy. Currently, only ALL_VLAN and VPC_LIST are supported. If scope is not provided, the default is set based on whether vpcReferences field is provided or not.
- `vpc_reference`: (Optional) A list of external ids for VPCs, used only when the scope of policy is a list of VPCs.

### rules

- `description`: (Optional) A user defined annotation for a rule.
- `type`: (Required) The type for a rule—the value chosen here restricts which specification can be chosen. Acceptable values are "QUARANTINE", "TWO_ENV_ISOLATION", "APPLICATION", "INTRA_GROUP".
- `spec`: (Required) Spec for rules.

### spec

One of below rules spec.

- `two_env_isolation_rule_spec`: (Optional) Two Environment Isolation Rule Spec.
- `application_rule_spec`: (Optional) Application Rule Spec.
- `intra_entity_group_rule_spec`: (Optional) Intra entity group Rule Spec
- `multi_env_isolation_rule_spec`: (Optional) Multi Environment Isolation Rule Spec.

### two_env_isolation_rule_spec

- `first_isolation_group`: (Required) Denotes the first group of category uuids that will be used in an isolation policy.
- `second_isolation_group`: (Required) Denotes the second group of category uuids that will be used in an isolation policy.

### application_rule_spec

- `secured_group_category_references`: (Required) A set of network endpoints which is protected by a Network Security Policy and defined as a list of categories.
- `src_allow_spec`: (Optional) A specification to how allow mode traffic should be applied, either ALL or NONE.
- `dest_allow_spec`: (Optional) A specification to how allow mode traffic should be applied, either ALL or NONE.
- `src_category_references`: (Optional) List of categories that define a set of network endpoints as inbound.
- `dest_category_references`: (Optional) List of categories that define a set of network endpoints as outbound.
- `src_subnet`: (Optional) source subnet value
- `dest_subnet`: (Optional) destination subnet value
- `src_address_group_references`: (Optional) A list of address group references.
- `dest_address_group_references`: (Optional) A list of address group references.
- `service_group_references`: (Optional) A list of service group references.
- `is_all_protocol_allowed`: (Optional) Denotes if rule allows traffic for all protocol.
- `tcp_services`: (Optional) tcp services
- `udp_services`: (Optional) udp services
- `icmp_services`: (Optional) icmp services
- `network_function_chain_reference`: (Optional) A reference to the network function chain in the rule.

### intra_entity_group_rule_spec

- `secured_group_action`: (Required) List of secured group action.
- `secured_group_category_references`: (Required) A specification to whether traffic between intra secured group entities should be allowed or denied.

### multi_env_isolation_rule_spec

- `spec`: (Required) Multi Environment Isolation Rule Spec.

#### spec

- `all_to_all_isolation_group`: all to all isolation groups

#### all_to_all_isolation_group

- `isolation_group`: (Required) Denotes the list of secured groups that will be used in All to All mutual isolation.

#### isolation_groups

- `group_category_references`: (Required) External identifiers of categories belonging to the isolation group.

### tcp_services, tcp_services

- `start_port`: (Required) start port
- `end_port`: (Required) end port

### icmp_services

- `is_all_allowed`: (Optional) Set this field to true if both Type and Code is ANY.
- `type`: (Optional) Icmp service Type. Ignore this field if Type has to be ANY.
- `code`: (Optional) Icmp service Code. Ignore this field if Code has to be ANY.

## Attributes Reference

The following attributes are exported:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `secured_groups`: Uuids of the secured groups in the NSP.
- `is_system_defined`: Is system defined NSP
- `created_by`: created by.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
- `last_update_time`: last updated time
- `creation_time`: creation time of NSP

## Import

This helps to manage existing entities which are not created through terraform. Network Security Policy can be imported using the `UUID`. (ext_id in v4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_network_security_policy_v2" "import_nsp" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_network_security_policies_v2" "list-nsps"{}
terraform import nutanix_network_security_policy_v2.import_nsp <UUID>
```

See detailed information in [Nutanix Security Policy v4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.0#tag/NetworkSecurityPolicies/operation/createNetworkSecurityPolicy).
