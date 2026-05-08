---
layout: "nutanix"
page_title: "NUTANIX: nutanix_events_v2"
sidebar_current: "docs-nutanix-datasource-events-v2"
description: |-
  Fetches a list of events.

---

# nutanix_events_v2

Fetches a list of events.

## Example Usage

```hcl
// list all events
data "nutanix_events_v2" "events-list" {}

// with limit
data "nutanix_events_v2" "events-limit" {
  limit = 4
}

// with filter and limit
data "nutanix_events_v2" "example"{
  filter = "startswith(eventType, 'A')"
  limit = 10
}
```

## Argument Reference

The following arguments are supported:
* `page`: -(Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: -(Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: -(Optional) A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions.
* `order_by`: -(Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default.
* `select`: -(Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned.

## Attributes Reference
The following attributes are exported:

* `events`: - List of events.

## Events
The `events` is a list of events. Each event supports the following attributes:

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `affected_entities`: - List of all the entities that are affected by the event or audit.
* `classifications`: - Various categories into which this event type can be classified. For example, hardware, storage, or license.
* `cluster_name`: - Name of the cluster associated with the entity.
* `cluster_uuid`: - Cluster UUID associated with the cluster where the event was first raised.
* `creation_time`: - The time in ISO 8601 format when the event was created.
* `event_type`: - A preconfigured or dynamically generated unique value for each event type.
* `message`: - Additional message associated with the event.
* `metric_details`: - Details of the metric for a metric-based event.
* `operation_type`: - The operation type associated with the audit. For example, create, update, or delete.
* `parameters`: - Additional parameters associated with the event. These parameters can be used to indicate custom key-value pairs for a given event instance. For example, a service down event in Prism Central can have the service name as a parameter.
* `service_name`: - The service which raised the event or audit. For internal Nutanix services, this value is set to "Nutanix".
* `source_cluster_uuid`: - Cluster UUID associated with the source entity of the event.
* `source_entity`: - The source entity associated with the event.


### Links
The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Affected Entities
The affected_entities attribute supports the following:

* `ext_id`: - UUID of the entity.
* `name`: - The name of the entity.
* `type`: - The type of entity. For example, VM, node, or cluster.

### Metric Details
The metric_details attribute supports the following:

* `comparison_operator`: - The comparison operator used for the condition evaluation.
* `condition_type`: - Indicating if this symptom is caused by a static threshold or anomaly (dynamic threshold) evaluation.
* `data_type`: - Data type of the metric value as stored in the database.
* `metric_category`: - Broad category under which this metric falls. For example, disk, CPU, or memory.
* `metric_display_name`: - Readable name of the metric in English.
* `metric_name`: - The metric key.
* `metric_value`: - The raw value of the metric when the condition threshold was exceeded.
* `threshold_value`: - The threshold value that was used for the condition evaluation.
* `trigger_time`: - The time in ISO 8601 format when the event was triggered.
* `trigger_wait_time_seconds`: - How long the metric breached the given condition before raising an event.
* `unit`: - Unit of the metric. For example, percentage, ms or usecs.

#### Metric Value / Threshold Value
The metric_value and threshold_value attributes support the following:

* `string_value`: - Denotes a value of type string.
* `bool_value`: - Denotes a value of type boolean.
* `int_value`: - Denotes a value of type integer.
* `double_value`: - Denotes a value of type double.

### Parameters
The parameters attribute supports the following:

* `param_name`: - Name or key of additional parameter for an instance.
* `param_value`: - Value of additional parameter for an instance.

#### Param Value
The param_value attribute supports the following:

* `string_value`: - Denotes a value of type string.
* `bool_value`: - Denotes a value of type boolean.
* `int_value`: - Denotes a value of type integer.

### Source Entity
The source_entity attribute supports the following:

* `ext_id`: - UUID of the entity.
* `name`: - The name of the entity.
* `type`: - The type of entity. For example, VM, node, or cluster.


See detailed information in [Nutanix List Events V4](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.2#tag/EventsService/operation/ListEvents).
