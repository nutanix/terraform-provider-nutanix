---
layout: "nutanix"
page_title: "NUTANIX: nutanix_protection_policies_v2"
sidebar_current: "docs-nutanix-datasource-protection-policies-v2"
description: |-
  List the protection policies defined on the system. This operation supports filtering, sorting, selection and pagination.



---

# nutanix_protection_policies_v2

List the protection policies defined on the system. This operation supports filtering, sorting, selection and pagination.



## Example Usage

```hcl
// with filter
data "nutanix_protection_policies_v2" "pps-filter" {
  filter = "name eq 'example_protection_policy'"
}

// with limit
data "nutanix_protection_policies_v2" "pp-limit" {
  limit = 4
}

// with filter and limit
data "nutanix_protection_policies_v2" "example"{
  filter = "startswith(name, 'C')"
  limit = "10"
}
```

## Argument Reference

The following arguments are supported:
* `page`: -(Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: -(Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: -(Optional) A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions. For example, filter '$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - extId
  - name
  - ownerExtId
  - replicationLocations/domainManagerExtId
* `order_by`: -(Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - name
* `select`: -(Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions. If a $select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned.
  - extId
  - name

## Attributes Reference
The following attributes are exported:

* `protection_policies`: - List of protection policies.

## Protection Policies
The `protection_policies` is a list of protection policies. Each protection policy supports the following attributes:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`: - Name of the protection policy.
* `description`: - Description of the protection policy.
* `replication_locations`: - Hypervisor details.
* `replication_configurations`: - Cluster reference for an entity.
* `category_ids`: - Host entity with its attributes.
* `is_approval_policy_needed`: - Disks attached to host.
* `owner_ext_id`: - Node degraded status.


### Links
The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Replication Locations
The replication_locations attribute supports the following:

* `label`: - This is a unique user defined label of the replication location. It is used to identify the location in the replication configurations.
* `domain_manager_ext_id`: - External identifier of the domain manager.
* `replication_sub_location`: - Specifies the replication sub-locations where recovery points can be created or replicated.
* `is_primary`: - One of the locations must be specified as the primary location. All the other locations must be connected to the primary location.

#### Replication Sub Location
The replication_sub_location attribute supports the following:
> One of `cluster_ext_ids` :
* `cluster_ext_ids` :  - External identifier of the clusters.

##### Cluster Ext Ids
The cluster_ext_ids attribute supports the following:

* `cluster_ext_id`: - List of Prism Element cluster external identifiers whose associated VMs and volume groups are protected. Only the primary location can have multiple clusters configured, while the other locations can specify only one cluster. Clusters must be specified for replication within the same Prism Central and cannot be specified for an MST type location. All clusters are considered if the cluster external identifier list is empty.

### Replication Configurations
The replication_configurations attribute supports the following:

* `source_location_label`: - Label of the source location from the replication locations list, where the entity is running. The location of type MST can not be specified as the replication source.
* `remote_location_label`: - Label of the source location from the replication locations list, where the entity will be replicated.
* `schedule`: - Schedule for protection. The schedule specifies the recovery point objective and the retention policy for the participating locations.

#### Schedule
The schedule attribute supports the following:

* `recovery_point_type`: - Type of recovery point.
   * `CRASH_CONSISTENT`: Crash-consistent Recovery points capture all the VM and application level details.
   * `APP_CONSISTENT`: Application-consistent Recovery points can capture all the data stored in the memory and also the in-progress transaction details.
* `recovery_point_objective_time_seconds`: - The Recovery point objective of the schedule in seconds and specified in multiple of 60 seconds. Only following RPO values can be provided for rollup retention type:
  - Minute(s): 1, 2, 3, 4, 5, 6, 10, 12, 15
  - Hour(s): 1, 2, 3, 4, 6, 8, 12
  - Day(s): 1
  - Week(s): 1, 2
* `retention`: - Specifies the retention policy for the recovery point schedule.
* `start_time`: - Represents the protection start time for the new entities added to the policy after the policy is created in h:m format. The values must be between 00h:00m and 23h:59m and in UTC timezone. It specifies the time when the first snapshot is taken and replicated for any entity added to the policy. If this is not specified, the snapshot is taken immediately and replicated for any new entity added to the policy.
* `sync_replication_auto_suspend_timeout_seconds`: - Auto suspend timeout if there is a connection failure between locations for synchronous replication. If this value is not set, then the policy will not be suspended.

#### Retention
> One of `linear_retention` or `auto_rollup_retention` :

* `linear_retention`: - Linear retention policy.
* `auto_rollup_retention`: - Auto rollup retention policy.

##### Linear Retention
The linear_retention attribute supports the following:

* `local`: - Specifies the number of recovery points to retain on the local location.
* `remote`: - Specifies the number of recovery points to retain on the remote location.

##### Auto Rollup Retention
The auto_rollup_retention attribute supports the following:

* `local`: - Specifies the auto rollup retention details.
* `remote`: - Specifies the auto rollup retention details.

###### Local, Remote
The local, remote attribute in the auto_rollup_retention supports the following:

* `snapshot_interval_type`: - Snapshot interval period.
   * `YEARLY`: Specifies the number of latest yearly recovery points to retain.
   * `WEEKLY`: Specifies the number of latest weekly recovery points to retain.
   * `DAILY`: Specifies the number of latest daily recovery points to retain.
   * `MONTHLY`: Specifies the number of latest monthly recovery points to retain.
   * `HOURLY`: Specifies the number of latest hourly recovery points to retain.
* `frequency`: - Multiplier to 'snapshotIntervalType'. For example, if 'snapshotIntervalType' is 'YEARLY' and 'multiple' is 5, then 5 years worth of rollup snapshots will be retained.




See detailed information in [Nutanix List Protection Policies V4](https://developers.nutanix.com/api-reference?namespace=datapolicies&version=v4.0#tag/ProtectionPolicies/operation/listProtectionPolicies).
