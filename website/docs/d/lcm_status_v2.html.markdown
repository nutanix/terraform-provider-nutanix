---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_status_v2"
sidebar_current: "docs-nutanix-datasource-lcm-status-v2"
description: |-
  Get the LCM framework status.
---

# nutanix_lcm_status_v2

Get the LCM framework status. Represents the Status of LCM. Status represents details about a pending or ongoing action in LCM.

## Example

```hcl
data "nutanix_lcm_status_v2" "lcm_framework_status" {
  x_cluster_id = "0005a104-0b0b-4b0b-8005-0b0b0b0b0b0b"
}
```

## Argument Reference
The following arguments are supported:

* `x_cluster_id`: (Optional) Cluster uuid on which the resource is present or operation is being performed.

## Attributes Reference
The following attributes are exported:


* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `framework_version`: LCM framework version information.
* `in_progress_operation`: Operation type and UUID of an ongoing operation in LCM.
* `is_cancel_intent_set`: Boolean that indicates if cancel intent for LCM update is set or not.
* `upload_task_uuid`: Upload task UUID.

### FrameworkVersion
The `framework_version` attribute supports the following:

* `current_version`: - Current LCM Version.
* `available_version`: - LCM framework version present in the LCM URL.
* `is_update_needed`: - Boolean that indicates if LCM framework update is needed.

### InProgressOperation
The `in_progress_operation` attribute supports the following:

* `operation_type`: - Type of the operation tracked by the task. Values are:
  - `PRECHECKS`: Perform LCM prechecks for the intended update operation.
  - `INVENTORY`: Perform an LCM inventory operation.
  - `UPGRADE`: Perform upgrade operation to a specific target version for discovered LCM entity/entities.
  - `NONE`: Indicates that no operation is currently ongoing.
* `operation_id`: - Root task UUID of the operation, if it is in running state.

See detailed information in [Nutanix LCM Status v4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Status)
