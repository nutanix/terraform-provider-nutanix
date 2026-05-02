---
layout: "nutanix"
page_title: "NUTANIX: nutanix_uda_policy_v2"
sidebar_current: "docs-nutanix-resource-uda-policy-v2"
description: |-
  Creates a new User-Defined Alert policy.

---

# nutanix_uda_policy_v2

Creates a new User-Defined Alert policy.

## Example Usage

```hcl
resource "nutanix_uda_policy_v2" "example" {
  title       = "example-uda-policy"
  entity_type = "VM"
  description = "Example User-Defined Alert policy"
  is_enabled  = true
  trigger_wait_period = 600

  trigger_conditions {
    condition {
      metric_name = "hypervisor_cpu_usage_ppm"
      operator    = "GREATER_THAN"
      threshold_value {
        int_value = 900000
      }
    }
    condition_type = "STATIC_THRESHOLD"
    severity_level = "CRITICAL"
  }
}
```

## Argument Reference

* `title` - (Required) Title of the policy.
* `entity_type` - (Required) Entity type associated with the User-Defined Alert policy. Allowed values are VM, node and cluster.
* `trigger_conditions` - (Required) Trigger conditions for the policy. If there are multiple trigger conditions, all of them will be considered during the operation.
* `description` - (Optional) Description of the policy.
* `filters` - (Optional) Filter criteria for narrowing down the entities on which User-Defined Alert policies can be set up.
* `impact_types` - (Optional) Impact types for the associated resulting alert.
* `is_auto_resolved` - (Optional) Indicates whether the auto-resolve feature is enabled for this policy.
* `is_enabled` - (Optional) Enable/Disable flag for the policy.
* `is_expected_to_error_on_conflict` - (Optional) Error when conflicting alert policies are found.
* `trigger_wait_period` - (Optional) Waiting duration in seconds before triggering the alert, when the specified condition is met. It is set to 600s by default.
* `policies_to_override` - (Optional) List of IDs of the alert policies that should be overridden.

### trigger_conditions

* `condition` - (Required) The condition for the trigger.
* `condition_type` - (Required) The condition type.
* `severity_level` - (Required) The severity level.

#### condition

* `metric_name` - (Required) The metric key.
* `operator` - (Required) Comparison operator.
* `threshold_value` - (Required) The threshold value that was used for the condition evaluation.

##### threshold_value

* `int_value` - (Optional) Denotes a value of type integer.
* `double_value` - (Optional) Denotes a value of type double.

### filters

* `entity_filter` - (Optional) List of entity filters.
* `group_filter` - (Optional) List of group filters.

#### entity_filter

* `ext_id` - (Required) Entity UUID on which the User-Defined Alert policy should be set up.

#### group_filter

* `ext_id` - (Required) Entity UUID of the group entity type on which the User-Defined Alert policy should be set up.
* `type` - (Required) The group entity type.

## Attribute Reference

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.
* `created_by` - Username of the user who created the policy.
* `last_updated_time` - Last updated time of the policy in ISO 8601 format.
* `related_policies` - List of alert policies that are related to the entities of the current policy.

### links

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object that is returned by the URL.

### related_policies

* `entity_uuid` - UUID of the entity the User-Defined Alert policy is associated with.
* `policy_ids` - Policy IDs associated with the specified entity.

## Import

`nutanix_uda_policy_v2` can be imported using the `ext_id`:

```shell
terraform import nutanix_uda_policy_v2.example 00000000-0000-0000-0000-000000000000
```

See detailed information in [Nutanix User-Defined Alert Policies](https://developers.nutanix.com/).
