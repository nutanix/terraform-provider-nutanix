---
layout: "nutanix"
page_title: "NUTANIX: nutanix_protection_rules"
sidebar_current: "docs-nutanix-datasource-protection-rules"
description: |-
 Describes a protection rules
---

# nutanix_protection_rules

Describes Protection Rules

## Example Usage

```hcl
data "nutanix_protection_rules" "protection_rules" {}
```


## Attribute Reference

The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of Protection Rules

# Entities

The entities attribute element contains the followings attributes:

* `name` The name for the protection rule.
* `description` A description for protection rule.

### Availability Zone Connectivity List
* `availability_zone_connectivity_list` - (Required) This encodes the datapipes between various availability zones and\nthe backup policy of the pipes.
* `availability_zone_connectivity_list.destination_availability_zone_index` - (Optional/Computed) Index of the availability zone.
* `availability_zone_connectivity_list.source_availability_zone_index` - (Optional/Computed) Index of the availability zone.
* `availability_zone_connectivity_list.snapshot_schedule_list` - (Optional/Computed) Snapshot schedules for the pair of the availability zones.
* `availability_zone_connectivity_list.snapshot_schedule_list.#.recovery_point_objective_secs` - (Required) "A recovery point objective (RPO) is the maximum acceptable amount of data loss.
* `availability_zone_connectivity_list.snapshot_schedule_list.#.local_snapshot_retention_policy` - (Optional/Computed) This describes the snapshot retention policy for this availability zone.
* `availability_zone_connectivity_list.snapshot_schedule_list.#.local_snapshot_retention_policy.0.num_snapshots` - (Optional/Computed) Number of snapshots need to be retained.
* `availability_zone_connectivity_list.snapshot_schedule_list.#.local_snapshot_retention_policy.0.rollup_retention_policy_multiple` - (Optional/Computed) Multiplier to 'snapshot_interval_type'.
* `availability_zone_connectivity_list.snapshot_schedule_list.#.local_snapshot_retention_policy.0.rollup_retention_policy_snapshot_interval_type` - (Optional/Computed)
* `availability_zone_connectivity_list.snapshot_schedule_list.#.auto_suspend_timeout_secs` - (Optional/Computed) Auto suspend timeout in case of connection failure between the sites.
* `availability_zone_connectivity_list.snapshot_schedule_list.#.snapshot_type` - (Optional/Computed) Crash consistent or Application Consistent snapshot.
* `availability_zone_connectivity_list.snapshot_schedule_list.#.remote_snapshot_retention_policy` - (Optional/Computed) This describes the snapshot retention policy for this availability zone.

### Ordered Availability Zone List
* `ordered_availability_zone_list` - (Required) A list of availability zones, each of which, receives a replica\nof the data for the entities protected by this protection rule.
* `ordered_availability_zone_list.#.cluster_uuid` - (Optional/Computed) UUID of specific cluster to which we will be replicating.
* `ordered_availability_zone_list.#.availability_zone_url` - (Optional/Computed) The FQDN or IP address of the availability zone. 

### Category Filter
* `category_filter` - (Optional/Computed)
* `category_filter.0.type` - (Optional/Computed) The type of the filter being used.
* `category_filter.0.kind_list` - (Optional/Computed) List of kinds associated with this filter.
* `category_filter.0.params` - (Optional/Computed) A list of category key and list of values.

### Metadata
The metadata attribute exports the following:

* `last_update_time` - UTC date and time in RFC-3339 format when vm was last updated.
* `uuid` - vm UUID.
* `creation_time` - UTC date and time in RFC-3339 format when vm was created.
* `spec_version` - Version number of the latest spec.
* `spec_hash` - Hash of the spec. This will be returned from server.
* `name` - vm name.

### Categories
The categories attribute supports the following:

* `name` - the key name.
* `value` - value of the key.

### Reference
The `project_reference`, `owner_reference` attributes supports the following:

* `kind` - (Required) The kind name (Default value: `project`).
* `name` - (Optional) the name.
* `uuid` - (Required) the UUID.


See detailed information in [Nutanix Protection Rules](https://www.nutanix.dev/api_references/prism-central-v3/#/23ee168cb0310-get-protection-rules).
