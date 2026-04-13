---
layout: "nutanix"
page_title: "NUTANIX: nutanix_storage_policies_v2"
sidebar_current: "docs-nutanix-datasource-storage-policies-v2"
description: |-
   This operation retrieves a List of the Storage Policies present in the system.
---

# nutanix_storage_policies_v2

Provides a datasource to Lists the Storage Policies present in the system.

## Example Usage

```hcl
data "nutanix_storage_policies_v2" "storage-policies"{ }
```

## Argument Reference

The following arguments are supported:


* `page`:- (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`:- (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`:- (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by`:- (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default.
* `select`:- A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions.


## Attribute Reference

The following attributes are exported:

* `storage_policies`:- Lists the Storage Policies present in the system.

## Storage Policies
The `storage_policies` contains list of Storage Policy objects. Each Storage Policy object contains the following attributes:

* `ext_id`:- External identifier of the Storage Policy.
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity.
* `links`:- A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`:- Storage Policy name.
* `category_ext_ids`:- List of external identifiers for Categories included in the Storage Policy.
* `compression_spec`:- Compression parameters for entities governed by the Storage Policy.
  * `compression_state`:- Compression state value.
* `encryption_spec`:- Encryption parameters for entities governed by the Storage Policy.
  * `encryption_state`:- Encryption state value.
* `qos_spec`:- Storage Quality of Service (QOS) parameters for the entities.
  * `throttled_iops`:- Throttled IOPS value.
* `fault_tolerance_spec`:- Fault Tolerance parameters for the entities.
  * `replication_factor`:- Replication factor value.
* `policy_type`:- Indicates whether the policy is user-created or system-created. Valid values: `"USER"`, `"SYSTEM"`.


See detailed information in [Nutanix Get Storage Policies v4](https://developers.nutanix.com/api-reference?namespace=datapolicies&version=v4.1#tag/StoragePolicies/operation/listStoragePolicies).