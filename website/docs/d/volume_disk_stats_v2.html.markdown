---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_disk_stats_v2"
sidebar_current: "docs-nutanix-datasource-volume-disk-stats-v2"
description: |-
  Query the Volume Disk stats identified by {diskExtId}.
---

# nutanix_volume_disk_stats_v2

Query the Volume Disk stats identified by {diskExtId}.

## Example Usage

```hcl
# Query the stats of a Volume Disk over a time window.
data "nutanix_volume_disk_stats_v2" "example" {
  volume_group_ext_id = "3770be9d-06be-4e25-b85d-3457d9b0ceb1"
  ext_id              = "1d92110d-26b5-46c0-8c93-20b8171373e0"
  start_time          = "2024-01-01T00:00:00Z"
  end_time            = "2024-01-01T01:00:00Z"
  sampling_interval   = 30
  stat_type           = "AVG"
}
```

## Argument Reference

The following arguments are supported:

* `volume_group_ext_id`: -(Required) The external identifier of a Volume Group.
* `ext_id`: -(Required) The external identifier of a Volume Disk.
* `start_time`: -(Required) The start time in RFC-3339 format from which the stats should be reported.
* `end_time`: -(Required) The end time in RFC-3339 format until which the stats should be reported.
* `sampling_interval`: -(Optional) The sampling interval in seconds at which the stats should be reported.
* `stat_type`: -(Optional) The operator to use while performing down-sampling on stats data. Allowed values are `SUM`, `MIN`, `MAX`, `AVG`, `COUNT` and `LAST`.
* `select`: -(Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type.

## Attributes Reference

The following attributes are exported:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `volume_disk_ext_id`: - Uuid of the Volume Disk.
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

### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Stat Time Value Pairs

Each of the `controller_*` attributes is a list where each element supports the following:

* `timestamp`: - Timestamp is returned in Epoch format.
* `value`: - Value of the stat at the corresponding timestamp value represented in Int64 format.

See detailed information in [Nutanix Get Volume Disk Stats V4](https://developers.nutanix.com/api-reference?namespace=volumes&version=v4.2#tag/VolumeGroups/operation/getVolumeDiskStats).
