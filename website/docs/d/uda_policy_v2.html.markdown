---
layout: "nutanix"
page_title: "NUTANIX: nutanix_uda_policy_v2"
sidebar_current: "docs-nutanix-datasource-uda-policy-v2"
description: |-
  Fetches the details of a User-Defined Alert policy identified by external identifier.

---

# nutanix_uda_policy_v2

Fetches the details of a User-Defined Alert policy identified by external identifier.

## Example Usage

```hcl
data "nutanix_uda_policy_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id` - (Required) A globally unique identifier of an instance that is suitable for external consumption.

## Attribute Reference

* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.
* `title` - Title of the policy.
* `description` - Description of the policy.
* `entity_type` - Entity type associated with the User-Defined Alert policy. Allowed values are VM, node and cluster.
* `trigger_conditions` - Trigger conditions for the policy. If there are multiple trigger conditions, all of them will be considered during the operation.
* `filters` - Filter criteria for narrowing down the entities on which User-Defined Alert policies can be set up.
* `impact_types` - Impact types for the associated resulting alert.
* `is_auto_resolved` - Indicates whether the auto-resolve feature is enabled for this policy.
* `is_enabled` - Enable/Disable flag for the policy.
* `trigger_wait_period` - Waiting duration in seconds before triggering the alert, when the specified condition is met.
* `created_by` - Username of the user who created the policy.
* `last_updated_time` - Last updated time of the policy in ISO 8601 format.
* `policies_to_override` - List of IDs of the alert policies that should be overridden.
* `related_policies` - List of alert policies that are related to the entities of the current policy.
* `is_expected_to_error_on_conflict` - Error when conflicting alert policies are found.

### links

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object that is returned by the URL.

### trigger_conditions

* `condition` - The condition for the trigger.
* `condition_type` - The condition type.
* `severity_level` - The severity level.

#### condition

* `metric_name` - The metric key.
* `operator` - Comparison operator.
* `threshold_value` - The threshold value that was used for the condition evaluation.

##### threshold_value

* `int_value` - Denotes a value of type integer.
* `double_value` - Denotes a value of type double.

### filters

* `entity_filter` - List of entity filters.
* `group_filter` - List of group filters.

#### entity_filter

* `ext_id` - Entity UUID on which the User-Defined Alert policy should be set up.

#### group_filter

* `ext_id` - Entity UUID of the group entity type on which the User-Defined Alert policy should be set up.
* `type` - The group entity type.

### related_policies

* `entity_uuid` - UUID of the entity the User-Defined Alert policy is associated with.
* `policy_ids` - Policy IDs associated with the specified entity.

See detailed information in [Nutanix User-Defined Alert Policies](https://developers.nutanix.com/).
