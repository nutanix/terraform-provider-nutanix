---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_policy_rules_v2"
sidebar_current: "docs-nutanix-datasource-network_security_policy_rules_v2"
description: |-
  List Network Security Policy rules for a given policy
---

# nutanix_network_security_policy_rules_v2

Gets the list of Network Security Policy rules for a given policy ExtID.

### Example

```hcl
# List all rules for a policy
data "nutanix_network_security_policy_rules_v2" "rules" {
  policy_ext_id = "4f9b9e0c-e473-4f2a-b20a-3e8c8e90236d"
}

# With pagination
data "nutanix_network_security_policy_rules_v2" "rules" {
  policy_ext_id = "4f9b9e0c-e473-4f2a-b20a-3e8c8e90236d"
  page         = 0
  limit        = 50
}

# With filter and order_by
data "nutanix_network_security_policy_rules_v2" "rules" {
  policy_ext_id = "4f9b9e0c-e473-4f2a-b20a-3e8c8e90236d"
  filter        = "type eq 'APPLICATION'"
  order_by      = "description asc"
}
```

## Argument Reference

The following arguments are supported:

- `policy_ext_id`: (Required) ExtId of the network security policy to list rules for.
- `page`: (Optional) Page number for pagination (0-based).
- `limit`: (Optional) Maximum number of rules to return (1–100). Default is 50 if not set.
- `filter`: (Optional) Filter expression for the list. The filter can be applied to the following fields:
  - `type`, Example: `filter = "type eq Microseg.Config.RuleType'QUARANTINE'"`
- `order_by`: (Optional) Order by clause. The order_by can be applied to the following fields:
  - `type`, Example: `order_by = "type desc"`
- `select`: (Optional) Comma-separated list of fields to return. The select can be applied to the following fields:
  - `type`, Example: `select = "type"`
  - `extId`, Example: `select = "extId"`
  - `description`, Example: `select = "description"`

## Attributes Reference

The following attributes are exported:

- `network_security_policy_rules`: List of network security policy rules.

### network_security_policy_rules

Each rule in the list has:

- `ext_id`: Globally unique identifier of the rule.
- `description`: User-defined description for the rule.
- `tenant_id`: Tenant that owns the rule.
- `type`: Rule type (e.g. `TWO_ENV_ISOLATION`, `APPLICATION`, `INTRA_GROUP`, `MULTI_ENV_ISOLATION`).
- `links`: HATEOAS-style links (e.g. `href`, `rel`).
- `spec`: Rule specification (one of the following blocks).

### spec

One of:

- `two_env_isolation_rule_spec`: Two-environment isolation rule (`first_isolation_group`, `second_isolation_group`).
- `application_rule_spec`: Application rule (secured groups, src/dest allow, categories, subnets, address/service groups, TCP/UDP/ICMP services, etc.).
- `intra_entity_group_rule_spec`: Intra-entity group rule (`secured_group_action`, `secured_group_category_references`).
- `multi_env_isolation_rule_spec`: Multi-environment isolation rule (`spec` → `all_to_all_isolation_group` → `isolation_group` with `group_category_references`).

See [Nutanix List Network Security Policy Rules v4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.2) for full API details.
