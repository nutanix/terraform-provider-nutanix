---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_rule"
sidebar_current: "docs-nutanix-datasource-network_security_rule"
description: |-
 Describes a Network security rule
---

# nutanix_network_security_rule

Describes a Network security rule

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

The following arguments are supported:

* `name`: - (Required) The name for the image.
* `categories`: - (Optional) Categories for the image.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.
* `api_version` - (Optional)
* `description`: - (Optional) A description for image.
* `quarantine_rule_action`: - (Optional) These rules are used for quarantining suspected VMs. Target group is a required attribute. Empty inbound_allow_list will not allow anything into target group. Empty outbound_allow_list will allow everything from target group.
* `quarantine_rule_outbound_allow_list`: - (Optional)
* `quarantine_rule_target_group_default_internal_policy`: - (Optional) - Default policy for communication within target group.
* `quarantine_rule_target_group_peer_specification_type`: - (Optional) - Way to identify the object for which rule is applied.
* `quarantine_rule_target_group_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `quarantine_rule_target_group_filter_type`: - (Optional) - The type of the filter being used.
* `quarantine_rule_target_group_filter_params`: - (Optional) - A list of category key and list of values.
* `quarantine_rule_inbound_allow_list`: - (Optional)
* `app_rule_action`: - (Optional) - These rules govern what flows are allowed. Target group is a required attribute. Empty inbound_allow_list will not anything into target group. Empty outbound_allow_list will allow everything from target group.
* `app_rule_outbound_allow_list`: - (Optional)
* `app_rule_target_group_default_internal_policy`: - (Optional) - Default policy for communication within target group.
* `app_rule_target_group_peer_specification_type`: - (Optional) - Way to identify the object for which rule is applied.
* `app_rule_target_group_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `app_rule_target_group_filter_type`: - (Optional) - The type of the filter being used.
* `app_rule_target_group_filter_params`: - (Optional) - A list of category key and list of values.
* `app_rule_inbound_allow_list`: - (Optional)
* `isolation_rule_action`: - (Optional) - These rules are used for environmental isolation.
* `app_rule_inbound_allow_list`: - (Optional)
* `isolation_rule_first_entity_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `isolation_rule_first_entity_filter_type`: - (Optional) - The type of the filter being used.
* `isolation_rule_first_entity_filter_params`: - (Optional) - A list of category key and list of values.
* `isolation_rule_second_entity_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `isolation_rule_second_entity_filter_type`: - (Optional) - The type of the filter being used.
* `isolation_rule_second_entity_filter_params`: - (Optional) - A list of category key and list of values.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when image was last updated.
* `UUID`: - image UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when image was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - image name.

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

See detailed information in [Nutanix Security Rules](https://www.nutanix.dev/reference/prism_central/v3/api/network-security-rules/getnetworksecurityrulesuuid).
