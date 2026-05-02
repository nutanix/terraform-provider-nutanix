---
layout: "nutanix"
page_title: "NUTANIX: nutanix_sda_policy_v2"
sidebar_current: "docs-nutanix-datasource-sda-policy-v2"
description: |-
  Get details of a System-Defined Alert Policy identified by external identifier.
---

# nutanix_sda_policy_v2

Get details of a System-Defined Alert Policy identified by external identifier.

## Example

```hcl
data "nutanix_sda_policy_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: (Required) Unique ID of the System-Defined Alert Policy.

## Attribute Reference

The following attributes are exported:

- `name`: Name of the System-Defined Alert Policy.
- `description`: System-defined alert policy description.
- `title`: Title of a System-Defined Alert Policy.
- `policy_id`: Unique ID associated with the policy.
- `publisher`: Publisher of the policy. For example, NCC for all health check policies.
- `entity_type`: Entity type of the policy.
- `scope`: Scope of the policy.
- `sub_type`: Sub type of the policy.
- `sda_type`: Type of the policy.
- `classifications`: Various categories into which this alert type can be classified. For example, hardware, storage, or license.
- `impact_types`: Impact types to which this rule applies.
- `kb_articles`: List of knowledge base article links.
- `target_clusters`: Indicates the cluster type against which this policy can be executed.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `links`: A HATEOAS style link for the response.
- `cluster_configs`: SDA policy parameters that may differ across clusters since each cluster can run on different NCC versions.

### links

- `href`: The URL at which the entity described by the link can be accessed.
- `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.

### cluster_configs

- `alert_config`: Alert configuration for this cluster.
- `configurable_parameters`: Parameters of the SDA that are configurable by a user.
- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `is_enabled`: Indicates whether the SDA policy is enabled or not on the cluster.
- `last_modified_by_user`: Name of the user who made the latest update to this policy.
- `last_modified_time`: Time in ISO 8601 format when the SDA policy was last modified.
- `links`: A HATEOAS style link for the response.
- `schedule_interval_seconds`: Interval in seconds for periodically executing the SDA policy.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

See detailed attribute reference for `alert_config` and `configurable_parameters` in the `nutanix_sda_cluster_config_v2` datasource documentation.
