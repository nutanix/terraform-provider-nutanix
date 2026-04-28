---
layout: "nutanix"
page_title: "NUTANIX: nutanix_audits_v2"
sidebar_current: "docs-nutanix-datasource-audits-v2"
description: |-
  Fetches a list of audits.

---

# nutanix_audits_v2

Fetches a list of audits.

## Example Usage

```hcl
# List all audits
data "nutanix_audits_v2" "audits" {}

# List audits with filter
data "nutanix_audits_v2" "filtered_audits" {
  filter = "serviceName eq 'Nutanix'"
}

# List audits with limit
data "nutanix_audits_v2" "limited_audits" {
  limit = 10
}
```

## Argument Reference

The following arguments are supported:

* `page`: -(Optional) A URL query parameter that specifies the page number of the result set.
* `limit`: -(Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter`: -(Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`: -(Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select`: -(Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attributes Reference

The following attributes are exported:

* `audits`: - List of audits.

### Audits

The audits attribute supports the following:

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `audit_type`: - The unique name for a given audit type. For example, VMCloneAudit or VMDeleteAudit.
* `message`: - Additional message associated with the audit.
* `service_name`: - The service which raised the event or audit. For internal Nutanix services, this value is set to "Nutanix".
* `operation_type`: - The operation type of the audit.
* `status`: - The status of the audit.
* `creation_time`: - The time in ISO 8601 format when the audit was created.
* `operation_start_time`: - The audit operation start time in ISO 8601 format.
* `operation_end_time`: - The audit operation end time in ISO 8601 format.
* `affected_entities`: - List of all the entities that are affected by the event or audit.
* `cluster_reference`: - The cluster reference associated with the audit.
* `source_entity`: - The source entity associated with the audit.
* `user_reference`: - The user reference associated with the audit.
* `parameters`: - Additional parameters associated with the audit. These parameters can be used to indicate custom key-value pairs for a given audit instance. For example, a service down audit in Prism Central can have the service name as a parameter.

### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Affected Entities

The affected_entities attribute supports the following:

* `ext_id`: - UUID of the entity.
* `name`: - The name of the entity.
* `type`: - The type of entity. For example, VM, node, or cluster.

### Cluster Reference

The cluster_reference attribute supports the following:

* `ext_id`: - UUID of the entity.
* `name`: - The name of the entity.
* `type`: - The type of entity. For example, VM, node, or cluster.

### Source Entity

The source_entity attribute supports the following:

* `ext_id`: - UUID of the entity.
* `name`: - The name of the entity.
* `type`: - The type of entity. For example, VM, node, or cluster.

### User Reference

The user_reference attribute supports the following:

* `ext_id`: - Unique UUID of the user who initiated the operation.
* `name`: - The name of the user who initiated the operation.
* `ip_address`: - The IP address from where the operation was triggered.

### Parameters

The parameters attribute supports the following:

* `param_name`: - Name or key of additional parameter for an instance.
* `param_value`: - Value of additional parameter for an instance.

#### Param Value

The param_value attribute supports the following:

* `string_value`: - Denotes a value of type string.
* `bool_value`: - Denotes a value of type boolean.
* `int_value`: - Denotes a value of type integer.

See detailed information in [Nutanix List Audits V4](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.2#tag/Audits/operation/listAudits).
