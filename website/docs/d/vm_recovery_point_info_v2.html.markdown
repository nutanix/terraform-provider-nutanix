---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_recovery_point_info_v2"
sidebar_current: "docs-nutanix-datasource-vm-recovery-point-info-v2"
description: |-
  Provides a datasource to Query the VM recovery point identified by ex_id.
---

# nutanix_vm_recovery_point_info_v2

Get the VM recovery point identified by ex_id.

## Example Usage

```hcl
# vm recovery point details
data "nutanix_vm_recovery_point_info_v2" "rp-vm-info" {
  recovery_point_ext_id = "af1070f7-c946-49da-9b17-e337e06e0a18"
  ext_id                = "85ac418e-c847-45ab-9816-40a3c4de148c"
}

```

## Argument Reference

The following arguments are supported:

* `recovery_point_ext_id`: (Required) The external identifier that can be used to retrieve the recovery point using its URL.
* `ext_id`: (Required) The external identifier that can be used to identify a VM recovery point.


## Attribute Reference

The following attributes are exported:

* `ext_id`: recovery point UUID
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `consistency_group_ext_id`: External identifier of the Consistency group which the VM was part of at the time of recovery point creation.
* `location_agnostic_id`: Location agnostic identifier of the Recovery point.
* `disk_recovery_points`: array of disk recovery points.
* `vm_ext_id`: VM external identifier which is captured as a part of this recovery point.
* `vm_categories`: Category key-value pairs associated with the VM at the time of recovery point creation. The category key and value are separated by '/'. For example, a category with key 'dept' and value 'hr' is displayed as 'dept/hr'.
* `application_consistent_properties`: User-defined application-consistent properties for the recovery point.
*
### Links
The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.



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


See detailed information in [Nutanix Get Vm Recovery Point V4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/RecoveryPoints/operation/getVmRecoveryPointById).
