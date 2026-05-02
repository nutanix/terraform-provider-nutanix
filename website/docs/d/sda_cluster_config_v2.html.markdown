---
layout: "nutanix"
page_title: "NUTANIX: nutanix_sda_cluster_config_v2"
sidebar_current: "docs-nutanix-datasource-sda-cluster-config-v2"
description: |-
  Retrieves the cluster specific configuration associated with a System-Defined Alert Policy identified by external identifier of the System-Defined Alert Policy for a cluster identified by cluster identifier.
---

# nutanix_sda_cluster_config_v2

Retrieves the cluster specific configuration associated with a System-Defined Alert Policy identified by external identifier of the System-Defined Alert Policy for a cluster identified by cluster identifier.

## Example

```hcl
data "nutanix_sda_cluster_config_v2" "example" {
  system_defined_policy_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id                       = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

- `system_defined_policy_ext_id`: (Required) Unique ID of the System-Defined Alert Policy.
- `ext_id`: (Required) Cluster UUID.

## Attribute Reference

The following attributes are exported:

- `alert_config`: Alert configuration for this cluster.
- `configurable_parameters`: Parameters of the SDA that are configurable by a user.
- `is_enabled`: Indicates whether the SDA policy is enabled or not on the cluster.
- `last_modified_by_user`: Name of the user who made the latest update to this policy. Its value will be Nutanix if the last update is due to an upgrade event.
- `last_modified_time`: Time in ISO 8601 format when the SDA policy was last modified.
- `links`: A HATEOAS style link for the response.
- `schedule_interval_seconds`: Interval in seconds for periodically executing the SDA policy. This will not be set for policies with the type NOT_SCHEDULED & EVENT_DRIVEN.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

### alert_config

- `auto_resolve`: Auto resolve state (ENABLED, DISABLED, NOT_SUPPORTED).
- `critical_severity`: Critical severity configuration.
- `info_severity`: Info severity configuration.
- `warning_severity`: Warning severity configuration.

### severity_config (critical_severity, info_severity, warning_severity)

- `state`: Property state (ENABLED, DISABLED, NOT_SUPPORTED).
- `threshold_parameters`: Captures alert-related thresholds that correspond to a particular severity.

### threshold_parameters / configurable_parameters

- `display_name`: Equivalent name for the parameter used to display it on Prism UI.
- `name`: Unique identifier name for the parameter.
- `unit`: Unit for the parameter. For example, sec, %, MB, GB, and so on.
- `param_value`: Captures the parameter value.

### param_value

The param_value contains one of the following:

- `int_value`: Integer parameter value.
  - `current_int_value`: Captures the current value of the parameter.
  - `default_int_value`: Captures the default value of the parameter.
- `float_value`: Float parameter value.
  - `current_float_value`: Captures the current value of the parameter.
  - `default_float_value`: Captures the default value of the parameter.
- `bool_value`: Boolean parameter value.
  - `current_bool_value`: Captures the current value of the parameter.
  - `default_bool_value`: Captures the default value of the parameter.
- `str_value`: String parameter value.
  - `current_str_value`: Captures the current value of the parameter.
  - `default_str_value`: Captures the default value of the parameter.

### links

- `href`: The URL at which the entity described by the link can be accessed.
- `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.
