---
layout: "nutanix"
page_title: "NUTANIX: nutanix_self_service_app_snapshots"
sidebar_current: "docs-nutanix_self_service_app"
description: |-
  Describes recovery points (snapshots taken) present in an application.
---

# nutanix_self_service_app_snapshots

Describes recovery points (snapshots taken) present in an NCM Self Service Application.

## Example Usage

```hcl
data "nutanix_self_service_app_snapshots" "test" {
    app_name = "NAME OF APPLICATION"
    length = 250
    offset = 0
}
```

## Argument Reference

The following arguments are supported:

* `app_name`: - (Optional) The name of the application.
* `app_uuid`: - (Optional) The UUID of the application.

Both (`app_name` and `app_uuid`) are optional but atleast one of them to be provided for this data source to work.

## Attribute Reference

The following attributes are exported:

* `length`: - (Required) The number of snapshots to retrieve.
* `offset`: - (Required) The index of the first snapshot to return (for pagination). Default value: 0
* `api_version`: - (Computed) The API version used to fetch the snapshot data.
* `total_matches`: - The total number of recovery points available for the application.
* `kind`: - The kind of the resource.

### entities

The entities block contains a list of snapshots associated with the specified application. Each snapshot has the following attributes:

* `type`: -  The type of the snapshot.
* `name`: - The name of the snapshot.
* `uuid`: -  The UUID of the snapshot.
* `description`: - The description of the snapshot.
* `action_name`: - The name of the action to run to create the snapshot.
* `recovery_point_info_list`: - The recovery_point_info_list contains information about recovery points for the snapshots. Each recovery point has the following attributes:
    * `id`: -  The ID of the recovery point.
    * `name`: - The name of the recovery point.
    * `snapshot_id`: - The ID of the snapshot associated with the recovery point.
    * `kind`: -  The kind of recovery point.
    * `creation_time`: - The creation time of the recovery point.
    * `recovery_point_type`: -  The type of recovery point.
    * `expiration_time`: -  The expiration time of the recovery point.
    * `location_agnostic_uuid`: - The UUID for the location-agnostic reference of the recovery point.
    * `service_references`: -  A list of service references related to the recovery point.
    * `config_spec_reference`: -  A map containing configuration specification references for the recovery point.
* `api_version`: - The API version used to retrieve the snapshot data.
* `spec`: -  The spec block contains the specification details for the snapshot
* `creation_time`: -  The creation time of the snapshot.
* `last_update_time`: - The last update time of the snapshot.
* `spec_version`: - The version of the snapshot specification.
* `kind`: -  The type of resource represented by the snapshot specification.


See detailed information in [List recovery groups](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/Apps/paths/~1apps~1%7Buuid%7D~1recovery_groups~1list/post).
