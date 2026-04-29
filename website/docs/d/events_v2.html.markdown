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
data "nutanix_events_v2" "example" {
  limit = 10
}
```

## Argument Reference

The following arguments are supported:

* `page`: -(Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: -(Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: -(Optional) A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response.
* `order_by`: -(Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default.
* `select`: -(Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.

## Attributes Reference

The following attributes are exported:

* `events`: - List of events.

### Events

The events attribute supports the following (each element has the same schema as the `nutanix_event_v2` datasource):

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `event_type`: - A preconfigured or dynamically generated unique value for each event type.
* `message`: - Additional message associated with the event.
* `creation_time`: - The time in ISO 8601 format when the event was created.
* `cluster_name`: - Name of the cluster associated with the entity.
* `cluster_uuid`: - Cluster UUID associated with the cluster where the event was first raised.
* `service_name`: - The service which raised the event or audit. For internal Nutanix services, this value is set to "Nutanix".
* `source_cluster_uuid`: - Cluster UUID associated with the source entity of the event.
* `operation_type`: - The operation type associated with the event.
* `classifications`: - Various categories into which this event type can be classified. For example, hardware, storage, or license.
* `source_entity`: - The source entity associated with the event.
* `affected_entities`: - List of all the entities that are affected by the event or audit.
* `metric_details`: - Details of the metric for a metric-based event.
* `parameters`: - Additional parameters associated with the event.

See detailed information in [Nutanix List Events V4](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.2#tag/EventsService/operation/ListEvents).
