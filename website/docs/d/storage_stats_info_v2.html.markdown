---
layout: "nutanix"
page_title: "NUTANIX: nutanix_storage_container_stats_info_v2"
sidebar_current: "docs-nutanix-datasource-storage-stats-info-v2"
description: |-
   This operation retrieves a Stats for a Storage Container.
---

# nutanix_storage_container_stats_info_v2

Provides a datasource to Fetches the stats information of the Storage Container identified by {containerExtId}.


## Example Usage

```hcl
data "nutanix_storage_container_stats_info_v2" "example"{
   ext_id = "1891fd3a-1ef7-4947-af56-9ee4b973c6fd"
   start_time = "2024-08-01T00:00:00Z"
   end_time = "2024-08-30T00:00:00Z"
   sampling_interval = 1
   stat_type = "SUM"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) storage container UUID
* `start_time`: (Required) storage container UUID
* `end_time`: (Required) storage container UUID
* `sampling_interval`: (Optional) storage container UUID
* `stat_type`: (Optional) storage container UUID
    * available values:
        * `AVG`: - Aggregation indicating mean or average of all values.
        * `MIN`: - Aggregation containing lowest of all values.
        * `MAX`: - 	Aggregation containing highest of all values.
        * `LAST`: - Aggregation containing only the last recorded value.
        * `SUM`: - Aggregation with sum of all values.
        * `COUNT`: - Aggregation containing total count of values.

## Attribute Reference

The following attributes are exported:

* `ext_id`: - the storage container uuid
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `container_ext_id`: - the storage container uuid
* `controller_num_iops`: - Number of I/O per second.
* `controller_io_bandwidth_kbps`: - Total I/O bandwidth - kB per second.
* `controller_avg_io_latencyu_secs`: - Average I/O latency in micro secs.
* `controller_num_read_iops`: - Number of read I/O per second.
* `controller_num_write_iops`: - Number of write I/O per second.
* `controller_read_io_bandwidth_kbps`: - Read I/O bandwidth - kB per second.
* `controller_write_io_bandwidth_kbps`: - Write I/O bandwidth - kB per second.
* `controller_avg_read_io_latencyu_secs`: - Average read I/O latency in microseconds.
* `controller_avg_write_io_latencyu_secs`: - Average read I/O latency in microseconds.
* `storage_reserved_capacity_bytes`: - Implicit physical reserved capacity(aggregated on vDisk level due to thick provisioning) in bytes.
* `storage_actual_physical_usage_bytes`: - Actual physical disk usage of the container without accounting for the reservation.
* `data_reduction_saving_ratio_ppm`: - Saving ratio in PPM as a result of Deduplication, compression and Erasure Coding.
* `data_reduction_total_saving_ratio_ppm`: - Saving ratio in PPM consisting of Deduplication, Compression, Erasure Coding, Cloning, and Thin Provisioning.
* `storage_free_bytes`: - Free storage in bytes.
* `storage_capacity_bytes`: - Storage capacity in bytes.
* `data_reduction_saved_bytes`: - Storage savings in bytes as a result of all the techniques.
* `data_reduction_overall_pre_reduction_bytes`: - Usage in bytes before reduction of Deduplication, Compression, Erasure Coding, Cloning, and Thin provisioning.
* `data_reduction_overall_post_reduction_bytes`: - Usage in bytes after reduction of Deduplication, Compression, Erasure Coding, Cloning, and Thin provisioning.
* `data_reduction_compression_saving_ratio_ppm`: - Saving ratio in PPM as a result of the Compression technique.
* `data_reduction_dedup_saving_ratio_ppm`: - Saving ratio in PPM as a result of the Deduplication technique.
* `data_reduction_erasure_coding_saving_ratio_ppm`: - Saving ratio in PPM as a result of the Erasure Coding technique.
* `data_reduction_thin_provision_saving_ratio_ppm`: - Saving ratio in PPM as a result of the Thin Provisioning technique.
* `data_reduction_clone_saving_ratio_ppm`: - Saving ratio in PPM as a result of the Cloning technique.
* `data_reduction_snapshot_saving_ratio_ppm`: - Saving ratio in PPM as a result of Snapshot technique.
* `data_reduction_zero_write_savings_bytes`: - Total amount of savings in bytes as a result of zero writes.
* `controller_read_io_ratio_ppm`: - Ratio of read I/O to total I/O in PPM.
* `controller_write_io_ratio_ppm`: - Ratio of read I/O to total I/O in PPM.
* `storage_replication_factor`: - Replication factor of Container.
* `storage_usage_bytes`: - Used storage in bytes.
* `storage_tier_das_sata_usage_bytes`: - Total usage on HDD tier for the Container in bytes.
* `storage_tier_ssd_usage_bytes`: - Total usage on SDD tier for the Container in bytes
* `health`: - Health of the container is represented by an integer value in the range 0-100. Higher value is indicative of better health.

### controller_num_iops,controller_io_bandwidth_kbps,controller_io_bandwidth_kbps ....,health

* `value`: Value of the stat at the recorded date and time in extended ISO-8601 format."
* `timestamp`: The date and time at which the stat was recorded.The value should be in extended ISO-8601 format. For example, start time of 2022-04-23T01:23:45.678+09:00 would consider all stats starting at 1:23:45.678 on the 23rd of April 2022. Details around ISO-8601 format can be found at https://www.iso.org/standard/70907.html



See detailed information in [Nutanix Get Stats for a Storage Container v4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0#tag/StorageContainers/operation/getStorageContainerStats).
