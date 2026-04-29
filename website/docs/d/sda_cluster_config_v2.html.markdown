---
layout: "nutanix"
page_title: "NUTANIX: nutanix_sda_cluster_config_v2"
sidebar_current: "docs-nutanix-datasource-sda-cluster-config-v2"
description: |-
  Retrieves the cluster specific configuration associated with a System-Defined Alert Policy identified by external identifier of the System-Defined Alert Policy for a cluster identified by cluster identifier.
---

# nutanix_sda_cluster_config_v2

Retrieves the cluster specific configuration associated with a System-Defined Alert Policy identified by external identifier of the System-Defined Alert Policy for a cluster identified by cluster identifier.

## Example Usage

```hcl
data "nutanix_sda_cluster_config_v2" "example" {
  system_defined_policy_ext_id = "<system-defined-policy-ext-id>"
  ext_id                       = "<cluster-uuid>"
}
```

## Argument Reference

The following arguments are supported:

* `system_defined_policy_ext_id`: (Required) Unique ID of the System-Defined Alert Policy.
* `ext_id`: (Required) Cluster UUID.

## Attribute Reference

The following attributes are exported:

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
  * `rel`: A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.
  * `href`: The URL at which the entity described by the link can be accessed.
* `is_enabled`: Indicates whether the SDA policy is enabled or not on the cluster.
* `last_modified_by_user`: Name of the user who made the latest update to this policy. Its value will be Nutanix if the last update is due to an upgrade event.
* `last_modified_time`: Time in ISO 8601 format when the SDA policy was last modified. It gets automatically updated by the Nutanix service from the user context during an update event.
* `schedule_interval_seconds`: Interval in seconds for periodically executing the SDA policy. This will not be set for policies with the type NOT_SCHEDULED & EVENT_DRIVEN.
* `alert_config`: Alert configuration for the cluster.
  * `auto_resolve`: Auto-resolve state for the alert policy.
  * `critical_severity`: Configuration for critical severity.
    * `state`: Property state (ENABLED, DISABLED, NOT_SUPPORTED).
    * `threshold_parameters`: Captures alert-related thresholds that correspond to a particular severity.
      * `display_name`: Equivalent name for the parameter used to display it on Prism UI.
      * `name`: Unique identifier name for the parameter.
      * `param_value`: Captures the parameter value.
        * `int_value`: Integer parameter value.
          * `current_int_value`: Captures the current value of the parameter.
          * `default_int_value`: Captures the default value of the parameter.
        * `float_value`: Float parameter value.
          * `current_float_value`: Captures the current value of the parameter.
          * `default_float_value`: Captures the default value of the parameter.
        * `bool_value`: Boolean parameter value.
          * `current_bool_value`: Captures the current value of the parameter.
          * `default_bool_value`: Captures the default value of the parameter.
        * `string_value`: String parameter value.
          * `current_str_value`: Captures the current value of the parameter.
          * `default_str_value`: Captures the default value of the parameter.
      * `unit`: Unit for the parameter. For example, sec, %, MB, GB, and so on.
  * `info_severity`: Configuration for info severity. Same structure as `critical_severity`.
  * `warning_severity`: Configuration for warning severity. Same structure as `critical_severity`.
* `configurable_parameters`: Parameters of the SDA that are configurable by a user. Same structure as `threshold_parameters`.

See detailed information in [Nutanix Monitoring v4 API Reference](https://developers.nutanix.com/).
