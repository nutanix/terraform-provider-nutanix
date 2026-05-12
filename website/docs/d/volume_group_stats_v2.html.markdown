---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_stats_v2"
sidebar_current: "docs-nutanix-datasource-volume-group-stats-v2"
description: |-
  Query the Volume Group stats identified by {extId}.
---

# nutanix_volume_group_stats_v2

Query the Volume Group stats identified by {extId}.

## Example Usage

```hcl
data "nutanix_volume_group_stats_v2" "example" {
  ext_id     = "d09aeec9-5bb7-4bfd-9717-a051178f6e7c"
  start_time = "2024-01-01T00:00:00Z"
  end_time   = "2024-01-02T00:00:00Z"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a Volume Group.
* `start_time`: -(Required) The start time for the stats query in RFC3339 format (e.g. 2024-01-01T00:00:00Z).
* `end_time`: -(Optional) The end time for the stats query in RFC3339 format (e.g. 2024-01-02T00:00:00Z). If not provided, defaults to current time.
* `sampling_interval`: -(Optional) The sampling interval in seconds.
* `stat_type`: -(Optional) The down sampling operator for the stats query.
* `select`: -(Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.

## Attributes Reference

The following attributes are exported:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity.
* `links`: - A HATEOAS style link for the response.
* `volume_group_ext_id`: - Uuid of the Volume Group.
* `controller_avg_io_latency_usecs`: - Controller average I/O latency measured in microseconds.
* `controller_avg_read_io_latency_usecs`: - Controller average read I/O latency measured in microseconds.
* `controller_avg_write_io_latency_usecs`: - Controller average write I/O latency measured in microseconds.
* `controller_io_bandwidth_kbps`: - Controller I/O bandwidth measured in Kbps.
* `controller_num_iops`: - Controller I/O rate measured in iops.
* `controller_num_read_iops`: - Controller read I/O measured in iops.
* `controller_num_write_iops`: - Controller write I/O measured in iops.
* `controller_read_io_bandwidth_kbps`: - Controller read I/O bandwidth measured in Kbps.
* `controller_user_bytes`: - Controller user bytes.
* `controller_write_io_bandwidth_kbps`: - Controller write I/O bandwidth measured in Kbps.
* `hydration_remaining_bytes`: - Number of bytes that are left to hydrate the Volume Group.

### Time Value Pair

Each stat attribute is a list of time value pairs with the following fields:

* `timestamp`: - Timestamp is returned in Epoch format.
* `value`: - Value of the stat at the corresponding timestamp value represented in Int64 format.

### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL.
