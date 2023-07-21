---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_rule"
sidebar_current: "docs-nutanix-datasource-network_security_rule"
description: |-
 Describes a Network security rule
---

# nutanix_network_security_rule

Describes a Network security rule

> NOTE: The use of network_security_rule is only applicable in AHV clusters and requires Microsegmentation to be enabled. This feature is a function of the Flow product and requires a Flow license. For more information on Flow and Microsegmentation please visit https://www.nutanix.com/products/flow

## Example Usage (Isolate Development VMs from Production VMs and get its information)

```hcl
resource "nutanix_network_security_rule" "isolation" {
	name        = "example-isolation-rule"
	description = "Isolation Rule Example"
	
	isolation_rule_action = "APPLY"
	
	isolation_rule_first_entity_filter_kind_list = ["vm"]
	isolation_rule_first_entity_filter_type      = "CATEGORIES_MATCH_ALL"
	isolation_rule_first_entity_filter_params {
		name   = "Environment"
		values = ["Dev"]
	}
	
	isolation_rule_second_entity_filter_kind_list = ["vm"]
	isolation_rule_second_entity_filter_type      = "CATEGORIES_MATCH_ALL"
	isolation_rule_second_entity_filter_params {
		name   = "Environment"
		values = ["Production"]
	}
}

data "nutanix_network_security_rule" "test" {
  network_security_rule_id = "${nutanix_network_security_rule.isolation.id}"
}
```

## Argument Reference

The following arguments are supported:

* `network_security_rule_id`: Represents network security rule UUID

## Attribute Reference

The following attributes are exported:

* `network_security_rule_id` - (Required) The ID for the rule you want to retrieve.
* `name`: - The name for the network_security_rule.
* `categories`: Categories for the network_security_rule.
* `project_reference`: The reference to a project.
* `owner_reference`: The reference to a user.
* `api_version`
* `description`: A description for network_security_rule.
* `quarantine_rule_action`: These rules are used for quarantining suspected VMs. Target group is a required attribute. Empty inbound_allow_list will not allow anything into target group. Empty outbound_allow_list will allow everything from target group.
* `quarantine_rule_outbound_allow_list`:
* `quarantine_rule_target_group_default_internal_policy`: - Default policy for communication within target group.
* `quarantine_rule_target_group_peer_specification_type`: - Way to identify the object for which rule is applied.
* `quarantine_rule_target_group_filter_kind_list`: - List of kinds associated with this filter.
* `quarantine_rule_target_group_filter_type`: - The type of the filter being used.
* `quarantine_rule_target_group_filter_params`: - A list of category key and list of values.
* `quarantine_rule_inbound_allow_list`:
* `app_rule_action`: - These rules govern what flows are allowed. Target group is a required attribute. Empty inbound_allow_list will not anything into target group. Empty outbound_allow_list will allow everything from target group.
* `app_rule_outbound_allow_list`:
* `app_rule_target_group_default_internal_policy`: - Default policy for communication within target group.
* `app_rule_target_group_peer_specification_type`: - Way to identify the object for which rule is applied.
* `app_rule_target_group_filter_kind_list`: - List of kinds associated with this filter.
* `app_rule_target_group_filter_type`: - The type of the filter being used.
* `app_rule_target_group_filter_params`: - A list of category key and list of values.
* `app_rule_inbound_allow_list`: The set of categories that matching VMs need to have.
* `ad_rule_action`: - These rules govern what flows are allowed. Target group is a required attribute. Empty inbound_allow_list will not anything into target group. Empty outbound_allow_list will allow everything from target group.
* `ad_rule_outbound_allow_list`:
* `ad_rule_target_group_default_internal_policy`: - Default policy for communication within target group.
* `ad_rule_target_group_peer_specification_type`: - Way to identify the object for which rule is applied.
* `ad_rule_target_group_filter_kind_list`: - List of kinds associated with this filter.
* `ad_rule_target_group_filter_type`: - The type of the filter being used.
* `ad_rule_target_group_filter_params`: - A list of category key and list of values.
* `ad_rule_inbound_allow_list`: The set of categories that matching VMs need to have.
* `isolation_rule_action`: - These rules are used for environmental isolation.
* `app_rule_inbound_allow_list`:
* `isolation_rule_first_entity_filter_kind_list`: - List of kinds associated with this filter.
* `isolation_rule_first_entity_filter_type`: - The type of the filter being used.
* `isolation_rule_first_entity_filter_params`: - A list of category key and list of values.
* `isolation_rule_second_entity_filter_kind_list`: - List of kinds associated with this filter.
* `isolation_rule_second_entity_filter_type`: - The type of the filter being used.
* `isolation_rule_second_entity_filter_params`: - A list of category key and list of values.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when network_security_rule was last updated.
* `UUID`: - network_security_rule UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when network_security_rule was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - network_security_rule name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, attributes supports the following:

* `kind`: - The kind name (Default value: project).
* `name`: - the name.
* `uuid`: - the UUID.

### Version

The version attribute supports the following:

* `product_name`: - Name of the producer/distribution of the image. For example windows or red hat.
* `product_version`: - Version string for the disk image.

See detailed information in [Nutanix Security Rules](https://www.nutanix.dev/api_references/prism-central-v3/#/064cd0a641d8d-get-a-existing-network-security-rule).
