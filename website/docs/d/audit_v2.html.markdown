---
layout: "nutanix"
page_title: "NUTANIX: nutanix_audit_v2"
sidebar_current: "docs-nutanix-datasource-audit-v2"
description: |-
  Fetches the details of an audit identified by external identifier.

---

# nutanix_audit_v2

Fetches the details of an audit identified by external identifier.

## Example Usage

```hcl
data "nutanix_audit_v2" "example" {
   ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:
* `ext_id`: - (Required) UUID of the generated audit.

## Attributes Reference
The following attributes are exported:

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `affected_entities`: - List of all the entities that are affected by the event or audit.
* `audit_type`: - The unique name for a given audit type. For example, VMCloneAudit or VMDeleteAudit.
* `cluster_reference`: - Reference to the cluster entity.
* `creation_time`: - The time in ISO 8601 format when the audit was created.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `message`: - Additional message associated with the audit.
* `operation_end_time`: - The audit operation end time in ISO 8601 format.
* `operation_start_time`: - The audit operation start time in ISO 8601 format.
* `operation_type`: - Type of operation performed.
* `parameters`: - Additional parameters associated with the audit. These parameters can be used to indicate custom key-value pairs for a given audit instance. For example, a service down audit in Prism Central can have the service name as a parameter.
* `service_name`: - The service which raised the event or audit. For internal Nutanix services, this value is set to "Nutanix".
* `source_entity`: - The entity that initiated the operation.
* `status`: - Status of the audit operation.
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `user_reference`: - Reference to the user who initiated the operation.

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

### Links
The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

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

### User Reference
The user_reference attribute supports the following:

* `ext_id`: - Unique UUID of the user who initiated the operation.
* `ip_address`: - The IP address from where the operation was triggered.
* `name`: - The name of the user who initiated the operation.

See detailed information in [Nutanix Get Audit V4](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.2#tag/AuditsService/operation/GetAuditById).
