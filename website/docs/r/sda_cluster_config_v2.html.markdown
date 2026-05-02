---
layout: "nutanix"
page_title: "NUTANIX: nutanix_sda_cluster_config_v2"
sidebar_current: "docs-nutanix-resource-sda-cluster-config-v2"
description: |-
  Modifies cluster-specific configuration associated with a System-Defined Alert Policy identified by external identifier of the System-Defined Alert Policy for a cluster identified by cluster identifier.
---

# nutanix_sda_cluster_config_v2

Modifies cluster-specific configuration associated with a System-Defined Alert Policy identified by external identifier of the System-Defined Alert Policy for a cluster identified by cluster identifier.

## Example

```hcl
resource "nutanix_sda_cluster_config_v2" "example" {
  system_defined_policy_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id                       = "00000000-0000-0000-0000-000000000000"

  alert_config {
    auto_resolve = "ENABLED"
    critical_severity {
      state = "ENABLED"
    }
    warning_severity {
      state = "ENABLED"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `system_defined_policy_ext_id`: (Required) Unique ID of the System-Defined Alert Policy.
- `ext_id`: (Required) Cluster UUID.
- `alert_config`: (Optional) Alert configuration for this cluster.

### alert_config

- `auto_resolve`: (Optional) Auto resolve state (ENABLED, DISABLED, NOT_SUPPORTED).
- `critical_severity`: (Optional) Critical severity configuration.
- `info_severity`: (Optional) Info severity configuration.
- `warning_severity`: (Optional) Warning severity configuration.

### severity_config (critical_severity, info_severity, warning_severity)

- `state`: (Optional) Property state (ENABLED, DISABLED, NOT_SUPPORTED).
- `threshold_parameters`: (Optional) Captures alert-related thresholds that correspond to a particular severity.

### threshold_parameters

- `name`: (Optional) Unique identifier name for the parameter.
- `param_value`: (Optional) Captures the parameter value.

### param_value

- `int_value`: (Optional) Integer parameter value.
  - `current_int_value`: (Optional) Captures the current value of the parameter.
- `float_value`: (Optional) Float parameter value.
  - `current_float_value`: (Optional) Captures the current value of the parameter.
- `bool_value`: (Optional) Boolean parameter value.
  - `current_bool_value`: (Optional) Captures the current value of the parameter.
- `str_value`: (Optional) String parameter value.
  - `current_str_value`: (Optional) Captures the current value of the parameter.

## Attribute Reference

The following attributes are exported:

- `configurable_parameters`: Parameters of the SDA that are configurable by a user.
- `is_enabled`: Indicates whether the SDA policy is enabled or not on the cluster.
- `last_modified_by_user`: Name of the user who made the latest update to this policy.
- `last_modified_time`: Time in ISO 8601 format when the SDA policy was last modified.
- `links`: A HATEOAS style link for the response.
- `schedule_interval_seconds`: Interval in seconds for periodically executing the SDA policy.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

See `nutanix_sda_cluster_config_v2` datasource documentation for the full attribute reference of computed fields.
