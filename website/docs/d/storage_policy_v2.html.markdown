---
layout: "nutanix"
page_title: "NUTANIX: nutanix_storage_policy_v2"
sidebar_current: "docs-nutanix-datasource-storage-policy-v2"
description: |-
   This operation retrieves a Storage Policy configuration.
---

# nutanix_storage_policy_v2

Provides a datasource to Fetch the configuration details of the existing Storage Policy identified by the {policyExtId}.

## Example Usage

```hcl
data "nutanix_storage_policy_v2" "get-storage-policy"{
  # Identifier of storage policy
  ext_id = "1891fd3a-1ef7-4947-af56-9ee4b973c6fd"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) The external identifier of the Storage Policy.

## Attribute Reference

The following attributes are exported:

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

See detailed information in [Nutanix Get Storage Policy v4](https://developers.nutanix.com/api-reference?namespace=datapolicies&version=v4.1#tag/StoragePolicies/operation/getStoragePolicyById).