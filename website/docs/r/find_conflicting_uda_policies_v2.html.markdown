---
layout: "nutanix"
page_title: "NUTANIX: nutanix_find_conflicting_uda_policies_v2"
sidebar_current: "docs-nutanix-resource-find-conflicting-uda-policies-v2"
description: |-
  Fetches all the existing policies with conflicting criteria to a User-Defined Alert policy identified by external identifier.

---

# nutanix_find_conflicting_uda_policies_v2

Fetches all the existing policies with conflicting criteria to a User-Defined Alert policy identified by external identifier.

## Example Usage

```hcl
resource "nutanix_find_conflicting_uda_policies_v2" "example" {
  title       = "conflict-check-policy"
  entity_type = "VM"

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
* `trigger_conditions` - (Required) Trigger conditions for the policy.
* `description` - (Optional) Description of the policy.
* `filters` - (Optional) Filter criteria for narrowing down the entities on which User-Defined Alert policies can be set up.
* `impact_types` - (Optional) Impact types for the associated resulting alert.
* `is_auto_resolved` - (Optional) Indicates whether the auto-resolve feature is enabled for this policy.
* `is_enabled` - (Optional) Enable/Disable flag for the policy.
* `trigger_wait_period` - (Optional) Waiting duration in seconds before triggering the alert, when the specified condition is met.

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

* `conflicting_policies` - List of conflicting policies.

### conflicting_policies

* `ext_id` - Unique UUID associated with the User-Defined Alert policy, that conflicts with the given policy.

See detailed information in [Nutanix User-Defined Alert Policies](https://developers.nutanix.com/).
