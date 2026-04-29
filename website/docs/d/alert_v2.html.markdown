---
layout: "nutanix"
page_title: "NUTANIX: nutanix_alert_v2"
sidebar_current: "docs-nutanix-datasource-alert-v2"
description: |-
  Fetches the details of an alert identified by external identifier.
---

# nutanix_alert_v2

Fetches the details of an alert identified by external identifier.

## Example Usage

```hcl
data "nutanix_alert_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) UUID of the generated alert.

## Attribute Reference

The following attributes are exported:

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `acknowledged_by_username`: Name of the user who acknowledged this alert.
* `acknowledged_time`: The time in ISO 8601 format when the alert was acknowledged.
* `affected_entities`: List of all the entities that are affected by the alert.
* `alert_type`: A preconfigured or dynamically generated unique value for each alert type.
* `classifications`: Various categories into which this alert type can be classified.
* `cluster_name`: Name of the cluster associated with the entity.
* `cluster_uuid`: Cluster UUID associated with the source entity of the alert.
* `creation_time`: Time in ISO 8601 format when the alert was created.
* `impact_types`: The impact this alert or event will have on the system.
* `is_acknowledged`: Indicates whether the alert is acknowledged or not.
* `is_auto_resolved`: Indicates whether the alert is auto-resolved or not.
* `is_resolved`: Indicates whether the alert is resolved or not.
* `is_runnable`: Indicates whether the policy associated with the alert is runnable or not.
* `is_user_defined`: Flag to indicate if the alert was generated from a User-Defined Alert policy.
* `kb_articles`: List of knowledge base article links.
* `last_updated_time`: Time in ISO 8601 format when the alert was last updated.
* `message`: Additional message associated with the alert.
* `metric_details`: Details of the metric for a metric-based event.
* `originating_cluster_uuid`: Cluster UUID associated with the cluster where the alert was first raised.
* `parameters`: Additional parameters associated with the alert.
* `resolved_by_username`: Name of the user who resolved this alert.
* `resolved_time`: The time in ISO 8601 format when the alert was resolved.
* `root_cause_analysis`: Possible causes, resolutions and additional details to troubleshoot this alert.
* `service_name`: The service that raised the alert.
* `severity`: Severity of the alert.
* `severity_trails`: Contains information on the severity change history for alerts.
* `source_entity`: Source entity of the alert.
* `title`: Title of the alert.

### affected_entities
* `ext_id`: UUID of the entity.
* `name`: The name of the entity.
* `type`: The type of entity. For example, VM, node, or cluster.

### source_entity
* `ext_id`: UUID of the entity.
* `name`: The name of the entity.
* `type`: The type of entity. For example, VM, node, or cluster.

### metric_details
* `comparison_operator`: Comparison operator used for the condition.
* `condition_type`: Type of the condition.
* `data_type`: Data type of the metric.
* `metric_category`: Broad category under which this metric falls.
* `metric_display_name`: Readable name of the metric in English.
* `metric_name`: The metric key.
* `metric_value`: The raw value of the metric when the condition threshold was exceeded.
* `threshold_value`: The threshold value that was used for the condition evaluation.
* `trigger_time`: The time in ISO 8601 format when the event was triggered.
* `trigger_wait_time_seconds`: How long the metric breached the given condition before raising an event.
* `unit`: Unit of the metric.

### parameters
* `param_name`: Name or key of additional parameter for an instance.
* `param_value`: Value of additional parameter for an instance.

### root_cause_analysis
* `cause`: Possible causes of this alert.
* `detail`: Additional details to troubleshoot this alert.
* `resolution`: Possible resolutions to troubleshoot this alert.

### severity_trails
* `severity`: Severity level.
* `severity_change_time`: The time in ISO 8601 format when the severity of the alert was changed.

### links
* `href`: The URL at which the entity described by the link can be accessed.
* `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.

See detailed information in [Nutanix Monitoring v4 Alerts](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.0).
