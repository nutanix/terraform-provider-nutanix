---
layout: "nutanix"
page_title: "NUTANIX: nutanix_alerts_v2"
sidebar_current: "docs-nutanix-datasource-alerts-v2"
description: |-
  Fetches a list of alerts.
---

# nutanix_alerts_v2

Fetches a list of alerts.

## Example Usage

```hcl
data "nutanix_alerts_v2" "example" {}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attribute Reference

The following attributes are exported:

* `alerts`: List of alert objects.

### alerts
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.
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

See detailed information in [Nutanix Monitoring v4 Alerts](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.0).
