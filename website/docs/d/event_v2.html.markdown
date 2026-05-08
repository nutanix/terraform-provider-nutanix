---
layout: "nutanix"
page_title: "NUTANIX: nutanix_event_v2"
sidebar_current: "docs-nutanix-datasource-event-v2"
description: |-
  Fetches the details of an event identified by external identifier.

---

# nutanix_event_v2

Fetches the details of an event identified by external identifier.

## Example Usage

```hcl
data "nutanix_event_v2" "example"{
   ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:
* `ext_id`: -(Required) UUID of the generated event.

## Attributes Reference
The following attributes are exported:

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


See detailed information in [Nutanix Get Event By Id V4](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.2#tag/EventsService/operation/GetEventById).
