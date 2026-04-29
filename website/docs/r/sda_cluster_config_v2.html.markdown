---
layout: "nutanix"
page_title: "NUTANIX: nutanix_sda_cluster_config_v2"
sidebar_current: "docs-nutanix-resource-sda-cluster-config-v2"
description: |-
  Modifies cluster-specific configuration associated with a System-Defined Alert Policy identified by external identifier of the System-Defined Alert Policy for a cluster identified by cluster identifier.
---

# nutanix_sda_cluster_config_v2

Modifies cluster-specific configuration associated with a System-Defined Alert Policy identified by external identifier of the System-Defined Alert Policy for a cluster identified by cluster identifier.

~> **Note:** This resource does not support delete. Destroying this resource will only remove it from state.

## Example Usage

### Basic Usage

```hcl
resource "nutanix_sda_cluster_config_v2" "example" {
  system_defined_policy_ext_id = "<system-defined-policy-ext-id>"
  ext_id                       = "<cluster-uuid>"
  is_enabled                   = true
}
```

### With Alert Config

```hcl
resource "nutanix_sda_cluster_config_v2" "example" {
  system_defined_policy_ext_id = "<system-defined-policy-ext-id>"
  ext_id                       = "<cluster-uuid>"
  is_enabled                   = true

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

* `system_defined_policy_ext_id`: (Required) Unique ID of the System-Defined Alert Policy.
* `ext_id`: (Required) Cluster UUID.
* `is_enabled`: (Optional) Indicates whether the SDA policy is enabled or not on the cluster.
* `schedule_interval_seconds`: (Optional) Interval in seconds for periodically executing the SDA policy. This will not be set for policies with the type NOT_SCHEDULED & EVENT_DRIVEN.
* `alert_config`: (Optional) Alert configuration for the cluster.
  * `auto_resolve`: (Optional) Auto-resolve state for the alert policy. Valid values: `ENABLED`, `DISABLED`, `NOT_SUPPORTED`.
  * `critical_severity`: (Optional) Configuration for critical severity.
    * `state`: (Optional) Property state. Valid values: `ENABLED`, `DISABLED`, `NOT_SUPPORTED`.
    * `threshold_parameters`: (Optional) Captures alert-related thresholds that correspond to a particular severity.
      * `name`: (Optional) Unique identifier name for the parameter.
      * `param_value`: (Optional) Captures the parameter value.
        * `int_value`: (Optional) Integer parameter value.
          * `current_int_value`: (Optional) Captures the current value of the parameter.
        * `float_value`: (Optional) Float parameter value.
          * `current_float_value`: (Optional) Captures the current value of the parameter.
        * `bool_value`: (Optional) Boolean parameter value.
          * `current_bool_value`: (Optional) Captures the current value of the parameter.
        * `string_value`: (Optional) String parameter value.
          * `current_str_value`: (Optional) Captures the current value of the parameter.
  * `info_severity`: (Optional) Configuration for info severity. Same structure as `critical_severity`.
  * `warning_severity`: (Optional) Configuration for warning severity. Same structure as `critical_severity`.
* `configurable_parameters`: (Optional) Parameters of the SDA that are configurable by a user. Same structure as `threshold_parameters`.

## Attribute Reference

The following attributes are exported:

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.
  * `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.
  * `href`: The URL at which the entity described by the link can be accessed.
* `last_modified_by_user`: Name of the user who made the latest update to this policy.
* `last_modified_time`: Time in ISO 8601 format when the SDA policy was last modified.

See detailed information in [Nutanix Monitoring v4 API Reference](https://developers.nutanix.com/).
