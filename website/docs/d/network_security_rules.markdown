---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_rules"
sidebar_current: "docs-nutanix-datasource-network-security-rules"
description: |-
 Describes a List of Network Security Rules
---

# nutanix_network_security_rules

Describes a List of Network Security Rules

## Example Usage

```hcl
resource "nutanix_category_key" "test-category-key"{
    name = "TIER-1"
	description = "TIER Category Key"
}


resource "nutanix_category_value" "WEB"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "WEB Category Value"
	 value = "WEB-1"
}

resource "nutanix_category_value" "APP"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "APP Category Value"
	 value = "APP-1"
}

resource "nutanix_category_value" "DB"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "DB Category Value"
	 value = "DB-1"
}

resource "nutanix_category_value" "ashwini"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "ashwini Category Value"
	 value = "ashwini-1"
}


resource "nutanix_network_security_rule" "TEST-TIER" {
  name        = "RULE-1-TIERS"
  description = "rule 1 tiers"

  app_rule_action = "APPLY"

  app_rule_inbound_allow_list = [
    {
      peer_specification_type = "FILTER"
      filter_type             = "CATEGORIES_MATCH_ALL"
      filter_kind_list        = ["vm"]

      filter_params = [
        {
          name   = "${nutanix_category_key.test-category-key.id}"
          values = ["${nutanix_category_value.WEB.id}"]
        },
      ]
    },
  ]

  app_rule_target_group_default_internal_policy = "DENY_ALL"

  app_rule_target_group_peer_specification_type = "FILTER"

  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"

  app_rule_target_group_filter_kind_list = ["vm"]

  app_rule_target_group_filter_params = [
    {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.APP.id}"]
    },
    {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.ashwini.id}"]
    },
  ]

  app_rule_outbound_allow_list = [
    {
      peer_specification_type = "FILTER"
      filter_type             = "CATEGORIES_MATCH_ALL"
      filter_kind_list        = ["vm"]

      filter_params = [
        {
          name   = "${nutanix_category_key.test-category-key.id}"
          values = ["${nutanix_category_value.DB.id}"]
        },
      ]
    },
  ]
}

data "nutanix_network_security_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `metadata`: Represents virtual machine UUID

### Metadata Argument

The metadata attribute supports the following:

* `kind`: - The kind name.
* `sort_attribute`: The attribute to perform sort on.
* `filter`: - The filter in FIQL syntax used for the results.
* `length`: - The number of records to retrieve relative to the offset.
* `sort_order`: - The sort order in which results are returned
* `offset`: - Offset from the start of the entity list

## Attribute Reference

The following attributes are exported:

* `entities`: - A list of virtual machines.

### Entities Attribute

The entities attribute supports the following:

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
* `UUID`: - the UUID.

### Version

The version attribute supports the following:

* `product_name`: - Name of the producer/distribution of the image. For example windows or red hat.
* `product_version`: - Version string for the disk image.

See detailed information in [Nutanix Network Security Rules](https://nutanix.github.io/Automation/experimental/swagger-redoc-sandbox/#tag/network_security_rules/paths/~1network_security_rules~1list/post).