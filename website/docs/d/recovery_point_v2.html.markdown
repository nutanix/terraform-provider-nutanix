---
layout: "nutanix"
page_title: "NUTANIX: nutanix_recovery_point_v2"
sidebar_current: "docs-nutanix-datasource-recovery-point-v2"
description: |-
  Provides a datasource to get a single Recovery Point corresponding to the extId.
---

# nutanix_recovery_point_v2

Get a single recovery point  corresponding to the extId.

## Example Usage

```hcl
data "nutanix_recovery_point_v2" "get-recovery-point" {
  ext_id = "1cefd0f5-6d38-4c9b-a07c-bdd2db004224"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) recovery point UUID


## Attribute Reference

The following attributes are exported:

* `ext_id`: recovery point UUID
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `location_agnostic_id`: Location agnostic identifier of the Recovery point.
* `name`: The name of the Recovery point.
* `creation_time`: The UTC date and time in ISO-8601 format when the Recovery point is created.
* `expiration_time`: The UTC date and time in ISO-8601 format when the current Recovery point expires and will be garbage collected.
* `status`: The status of the Recovery point, which indicates whether this Recovery point is fit to be consumed.
  * supported values:
      * `COMPLETE`: -  The Recovery point is in a complete state and ready to be consumed.
* `recovery_point_type`: Type of the Recovery point.
    * supported values:
      * `CRASH_CONSISTENT`: -  capture all the VM and application level details.
      * `APPLICATION_CONSISTENT`: -  stored in the memory and also the in-progress transaction details.
* `owner_ext_id`: A read only field inserted into recovery point at the time of recovery point creation, indicating the external identifier of the user who created this recovery point.
* `location_references`: List of location references where the VM or volume group recovery point are a part of the specified recovery point.
* `vm_recovery_points`: List of VM recovery point that are a part of the specified top-level recovery point. Note that a recovery point can contain a maximum number of 30 entities. These entities can be a combination of VM(s) and volume group(s).
* `volume_group_recovery_points`: List of volume group recovery point that are a part of the specified top-level recovery point. Note that a recovery point can contain a maximum number of 30 entities. These entities can be a combination of VM(s) and volume group(s).

### Links
The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.


### location_references

* `location_ext_id`: External identifier of the cluster where the recovery point is present.

### vm_recovery_points
* `ext_id`: recovery point UUID
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name` : The name of the Recovery point.
* `creation_time`: The UTC date and time in ISO-8601 format when the Recovery point is created.
* `expiration_time`: The UTC date and time in ISO-8601 format when the current Recovery point expires and will be garbage collected.
* `status`: The status of the Recovery point, which indicates whether this Recovery point is fit to be consumed.
  * supported values:
      * `COMPLETE`: -  The Recovery point is in a complete state and ready to be consumed.
* `recovery_point_type`: Type of the Recovery point.
* `consistency_group_ext_id`: External identifier of the Consistency group which the VM was part of at the time of recovery point creation.
* `location_agnostic_id`: Location agnostic identifier of the Recovery point.
* `disk_recovery_points`: array of disk recovery points.
* `vm_ext_id`: VM external identifier which is captured as a part of this recovery point.
* `vm_categories`: Category key-value pairs associated with the VM at the time of recovery point creation. The category key and value are separated by '/'. For example, a category with key 'dept' and value 'hr' is displayed as 'dept/hr'.
* `application_consistent_properties`: User-defined application-consistent properties for the recovery point.

### volume_group_recovery_points
* `ext_id`: recovery point UUID
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `consistency_group_ext_id`: External identifier of the Consistency group which the entity was part of at the time of recovery point creation.
* `location_agnostic_id`: Location agnostic identifier of the recovery point. This identifier is used to identify the same instances of a recovery point across different sites.
* `volume_group_ext_id`: Volume Group external identifier which is captured as part of this recovery point.
* `volume_group_categories`: Category key-value pairs associated with the volume group at the time of recovery point creation. The category key and value are separated by '/'. For example, a category with key 'dept' and value 'hr' will be represented as 'dept/hr'.
* `disk_recovery_points`: array of disk recovery points.


### disk_recovery_points
* `disk_recovery_point_ext_id`: External identifier of the disk recovery point.
* `disk_ext_id`: External identifier of the disk.


### application_consistent_properties
* `backup_type`: The backup type specifies the criteria for identifying the files to be backed up. This property should be specified to the application-consistent recovery points for Windows VMs/agents. The following backup types are supported for the application-consistent recovery points:
  * supported values:
    * `FULL_BACKUP`: -  All the files are backed up irrespective of their last backup date/time or state. Also, this backup type updates the backup history of each file that participated in the recovery point. If not explicitly specified, this is the default backup type.
    * `COPY_BACKUP`: -  this backup type does not update the backup history of individual files involved in the recovery point.
* `should_include_writers`: Indicates whether the given set of VSS writers' UUIDs should be included or excluded from the application consistent recovery point. By default, the value is set to false, indicating that all listed VSS writers' UUIDs will be excluded.
* `writers`: List of VSS writer UUIDs that are used in an application consistent recovery point. The default values are the system and the registry writer UUIDs.
* `should_store_vss_metadata`: Indicates whether to store the VSS metadata if the user is interested in application-specific backup/restore. The VSS metadata consists of VSS writers and requester metadata details. These are compressed into a cabinet file(.cab file) during a VSS backup operation. This cabinet file must be saved to the backup media during a backup operation, as it is required during the restore operation.
* `object_type`: value: `dataprotection.v4.common.VssProperties`


See detailed information in [Nutanix Get Recovery Point V4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/RecoveryPoints/operation/getRecoveryPointById).
