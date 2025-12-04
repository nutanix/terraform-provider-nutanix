---
layout: "nutanix"
page_title: "NUTANIX: nutanix_protected_resource_v2"
sidebar_current: "docs-nutanix-datasource-protected-resource-v2"
description: |-
  Get a protected resource


---

# nutanix_protected_resource_v2

Get the details of the specified protected resource such as the restorable time ranges available on the local Prism Central and the state of replication to the targets specified in the applied protection policies. This applies only if the entity is protected in a minutely or synchronous schedule. Other protection schedules are not served by this endpoint yet, and are considered not protected.


## Example 1: Get Protected Virtual Machine

```hcl

data "nutanix_protected_resource_v2" "protected-vm" {
  ext_id = "d22529bb-f02d-4710-894b-d1de772d7832" # protected vm ext_id
}

```

## Example 2: Get Protected Volume Group

```hcl

data "nutanix_protected_resource_v2" "protected-vg" {
  ext_id = "246c651a-1b16-4983-b5ff-204840f85e07" # protected volume group ext_id
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.


## Attributes Reference
The following attributes are exported:

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `entity_ext_id`: The external identifier of the VM or the volume group associated with the protected resource.
* `ext_id`: - The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.
* `entity_type`: Protected resource entity type. Possible values are: VM, VOLUME_GROUP.
* `source_site_reference`: Details about the data protection site in the Prism Central.
* `site_protection_info`: The data protection details for the protected resource that are relevant to any of the sites in the local Prism Central, like the time ranges available for recovery.
* `replication_states`: Replication related information about the protected resource.
* `consistency_group_ext_id`: External identifier of the Consistency group which the protected resource is part of.
* `category_fq_names`: Category key-value pairs associated with the protected resource at the time of protection. The category key and value are separated by '/'. For example, a category with key 'dept' and value 'hr' will be represented as 'dept/hr'.

### Links
The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Source Site Reference
The source_site_reference attribute supports the following:

* `mgmt_cluster_ext_id`: External identifier of the Prism Central.
* `cluster_ext_id`: External identifier of the cluster.

### Site Protection Info
The site_protection_info attribute supports the following:

* `recovery_info`:  The restorable time range details that can be used to recover the protected resource.
* `location_reference`: Details about the data protection site in the Prism Central.
* `synchronous_replication_role`: Synchronous Replication role related information of the protected resource. Possible values are:
  - `DECOUPLED`: VM is no longer in Synchronous Replication, and all the actions are blocked on VM, except a delete operation.
  - `SECONDARY`: This is the target site for VM in Synchronous Replication.
  - `INDEPENDENT`: VM is no longer in Synchronous Replication, and not replicating to the configured recovery cluster.
  - `PRIMARY`: VM is in Synchronous Replication, and is active on the primary site.

#### Recovery Info
The recovery_info attribute supports the following:

* `restorable_time_ranges`: The restorable time range details that can be used to recover the protected resource.

#### Restorable Time Range
The restorable_time_ranges attribute supports the following:

* `start_time`: UTC date and time in ISO 8601 format representing the time when the restorable time range for the entity starts.
* `end_time`: UTC date and time in ISO 8601 format representing the time when the restorable time range for the entity starts.

#### Location Reference
The location_reference attribute supports the following:

* `mgmt_cluster_ext_id`: External identifier of the Prism Central.
* `cluster_ext_id`: External identifier of the cluster.

### Replication States
The replication_states attribute supports the following:

* `protection_policy_ext_id`: The external identifier of the Protection policy associated with the protected resource.
* `recovery_point_objective_seconds`: The recovery point objective of the schedule in seconds.
* `replication_status`: Status of replication to a specified target site. Possible values are:
    - `IN_SYNC`: The specified recovery point objective is met on the target site and failover can be performed.
    - `SYNCING`: The system is trying to meet the specified recovery point objective for the target site via ongoing replications and failover can't yet be performed.
    - `OUT_OF_SYNC`: The replication schedule is disabled and there are no ongoing replications. Manual action might be needed by the user to meet the recovery point objective.
* `target_site_reference`: Details about the data protection site in the Prism Central.

#### Target Site Reference
The target_site_reference attribute supports the following:

* `mgmt_cluster_ext_id`: External identifier of the Prism Central.
* `cluster_ext_id`: External identifier of the cluster.


See detailed information in [Nutanix Get Protected Resource v4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/ProtectedResources/operation/getProtectedResourceById).

