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

# Network Security Policy with GLOBAL scope (VMs resolved by category across all VPCs)
resource "nutanix_network_security_policy_v2" "global-nsp" {
  name        = "my-global-policy"
  description = "Application policy with global scope"
  state       = "SAVE"
  type        = "APPLICATION"
  scope       = "GLOBAL"
  rules {
    type = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [nutanix_category_v2.example.id]
        service_group_references          = [nutanix_service_groups_v2.example.id]
        src_address_group_references      = [nutanix_address_groups_v2.example.id]
      }
    }
  }
}

```

## Argument Reference

The following arguments are supported:

- `name`: (Required) Name of the Flow Network Security Policy.
- `type`: (Required) Defines the type of rules that can be used in a policy. Acceptable values are "QUARANTINE", "ISOLATION", "APPLICATION", "SHAREDSERVICE", "CRITICAL", "COREINFRASTRUCTURE", "ZONE", "WORKLOAD". The "WORKLOAD", "CRITICAL", "COREINFRASTRUCTURE" and "ZONE" types are used with Flex rules (rule-centric / SMSP mode) and require `priority` to be set.
- `description`: (Optional) A user defined annotation for a policy.
- `state`: (Optional) Whether the policy is applied or monitored; can be omitted or set null to save the policy without applying or monitoring it. Acceptable values are "SAVE", "MONITOR", "ENFORCE".
- `priority`: (Optional) Policy priority. Mandatory for the Flex policy types (WORKLOAD/CRITICAL/COREINFRASTRUCTURE/ZONE); a lower value means higher precedence. For user-defined WORKLOAD policies use 1-349 (350 is reserved for the system catch-all).
- `rules`: (Optional) A list of rules that form a policy. For isolation policies, use isolation rules; for application or quarantine policies, use application rules; for Flex policies, use flex rules only.
- `is_ipv6_traffic_allowed`: (Optional) If Ipv6 Traffic is allowed.
- `is_hitlog_enabled`: (Optional) If Hitlog is enabled.
- `scope`: (Optional) Defines the scope of the policy. Acceptable values are "ALL_VLAN", "ALL_VPC", "VPC_LIST", "GLOBAL", and "VPC_AS_CATEGORY".
- `vpc_reference`: (Optional) A list of external ids for VPCs, used only when the scope of policy is a list of VPCs.
- `project_ext_id`: (Optional) Project external ID to associate with the network security policy. Note: This field cannot be updated after creation.
- `is_shared_with_all_projects`: (Optional) Indicates whether the network security policy is shared with all projects.

### rules

- `description`: (Optional) A user defined annotation for a rule.
- `type`: (Required) The type for a rule—the value chosen here restricts which specification can be chosen. Acceptable values are "QUARANTINE", "TWO_ENV_ISOLATION", "APPLICATION", "INTRA_GROUP", "MULTI_ENV_ISOLATION", "SHARED_SERVICE", "FLEX".
- `name`: (Optional) A user-defined name for the rule. Primarily used for Flex policy rules.
- `is_logging_enabled`: (Optional) Specifies whether hit log is enabled for the rule.
- `spec`: (Required) Spec for rules.

### spec

One of below rules spec.

- `two_env_isolation_rule_spec`: (Optional) Two Environment Isolation Rule Spec.
- `application_rule_spec`: (Optional) Application Rule Spec.
- `intra_entity_group_rule_spec`: (Optional) Intra entity group Rule Spec
- `multi_env_isolation_rule_spec`: (Optional) Multi Environment Isolation Rule Spec.
- `flex_rule_spec`: (Optional) Flex Rule Spec.

### two_env_isolation_rule_spec

- `first_isolation_group`: (Required) Denotes the first group of category uuids that will be used in an isolation policy.
- `second_isolation_group`: (Required) Denotes the second group of category uuids that will be used in an isolation policy.

### application_rule_spec

- `secured_group_category_associated_entity_type`: (Optional) Entity type for the secured group category. Acceptable values are "SUBNET", "VM", "VPC". Default is "VM".
- `secured_group_category_references`: (Optional) A set of network endpoints which is protected by a Network Security Policy and defined as a list of categories. Exactly one of `secured_group_category_references` and `secured_group_entity_group_reference` must be set.
- `secured_group_entity_group_reference`: (Optional) Reference to the secured group entity group. Exactly one of `secured_group_category_references` and `secured_group_entity_group_reference` must be set.
- `src_allow_spec`: (Optional) A specification to how allow mode traffic should be applied, either ALL or NONE.
- `dest_allow_spec`: (Optional) A specification to how allow mode traffic should be applied, either ALL or NONE.
- `src_category_associated_entity_type`: (Optional) Entity type for the source category. Acceptable values are "SUBNET", "VM", "VPC". Default is "VM".
- `src_category_references`: (Optional) List of categories that define a set of network endpoints as inbound.
- `src_entity_group_reference`: (Optional) Reference to the source entity group.
- `dest_category_associated_entity_type`: (Optional) Entity type for the destination category. Acceptable values are "SUBNET", "VM", "VPC". Default is "VM".
- `dest_category_references`: (Optional) List of categories that define a set of network endpoints as outbound.
- `dest_entity_group_reference`: (Optional) Reference to the destination entity group.
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
- `network_function_reference`: (Optional) A reference to the network function in the rule.

### intra_entity_group_rule_spec

- `secured_group_category_associated_entity_type`: (Optional) Entity type for the secured group category. Acceptable values are "SUBNET", "VM", "VPC". Default is "VM".
- `secured_group_entity_group_reference`: (Optional) Reference to the secured group entity group.
- `secured_group_action`: (Required) Whether traffic between intra secured group entities should be allowed or denied. Acceptable values are "ALLOW", "DENY".
- `secured_group_category_references`: (Optional) List of category references for the secured group.
- `secured_group_service_references`: (Optional) List of service group references for the secured group.
- `tcp_services`: (Optional) TCP port ranges for the rule.
- `udp_services`: (Optional) UDP port ranges for the rule.
- `icmp_services`: (Optional) ICMP type/code for the rule.

### multi_env_isolation_rule_spec

- `spec`: (Required) Multi Environment Isolation Rule Spec.

#### spec

- `all_to_all_isolation_group`: all to all isolation groups

#### all_to_all_isolation_group

- `isolation_group`: (Required) Denotes the list of secured groups that will be used in All to All mutual isolation.

#### isolation_group

- `group_category_associated_entity_type`: (Optional) Entity type for the group category. Acceptable values are "SUBNET", "VM", "VPC". Default is "VM".
- `group_category_references`: (Required) External identifiers of categories belonging to the isolation group.
- `group_entity_group_reference`: (Optional) Reference to the entity group for the isolation group.

### flex_rule_spec

- `action`: (Required) Action for the flex rule. Acceptable values are "ALLOW", "DENY", "REJECT".
- `direction`: (Required) Direction of traffic. Acceptable values are "IN", "OUT", "IN_OUT".
- `applied_to_entity_group_references`: (Optional) Entity group references to which the flex rule is applied.
- `src_entity_group_references`: (Optional) Source entity group references.
- `dest_entity_group_references`: (Optional) Destination entity group references.
- `src_subnet`: (Optional) Source subnet (value, prefix_length).
- `dest_subnet`: (Optional) Destination subnet (value, prefix_length).
- `should_allow_any_src`: (Optional) Whether the rule should allow all sources.
- `should_allow_any_dst`: (Optional) Whether the rule should allow all destinations.
- `is_all_protocol_allowed`: (Optional) Whether traffic is allowed for all protocols.
- `service_group_references`: (Optional) A list of service group references.
- `tcp_services`: (Optional) TCP port ranges (start_port, end_port).
- `udp_services`: (Optional) UDP port ranges (start_port, end_port).
- `icmp_services`: (Optional) ICMP type/code specs.
- `icmpv6_services`: (Optional) ICMPv6 type/code specs.
- `network_function_reference`: (Optional) A reference to the network function.
- `priority`: (Optional) Priority for the flex rule. Lower numbers indicate higher priority.
- `ip_version`: (Optional) IP version scope. Acceptable values are "IPV4", "IPV6", "IPV4_IPV6".
- `is_system_rule`: (Computed) Whether the flex rule is system-defined.

### tcp_services, udp_services

- `start_port`: (Required) start port
- `end_port`: (Required) end port

### icmp_services, icmpv6_services

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
- `project_ext_id`: Project external ID associated with the network security policy.
- `is_shared_with_all_projects`: Whether the network security policy is shared with all projects.

## Import

This helps to manage existing entities which are not created through terraform. Network Security Policy can be imported using the `UUID`. (ext_id in v4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_network_security_policy_v2" "import_nsp" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_network_security_policies_v2" "list-nsps"{}
terraform import nutanix_network_security_policy_v2.import_nsp <UUID>
```

See detailed information in [Nutanix Security Policy v4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.3#tag/NetworkSecurityPolicies/operation/createNetworkSecurityPolicy).
