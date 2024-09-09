---
layout: "nutanix"
page_title: "NUTANIX: nutanix_permission_v2"
sidebar_current: "docs-nutanix-datasource-permission-v2"
description: |-
  This operation retrieves a list of all the permission.
---

# nutanix_permission_v2
Lists the Permission defined on the system. List of permission can be further filtered out using various filtering options.

## Example

```hcl

     data "nutanix_operation_v2" "test" {
        ext_id = "<ext-id>"
    }

```

## Argument Reference

The following arguments are supported:

* `ext_id`:(Required) ExtId of the Operation.

## Attributes Reference

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `display_name`: Permission name.
* `dedescription`: Permission description
* `create_time`: Permission creation time
* `last_updated_time`: Permission last updated time.

* `entity_type`: Type of entity associated with this Operation.
* `operation_type`: The Operation type. Currently we support INTERNAL, EXTERNAL and SYSTEM_DEFINED_ONLY.
* `client_name`: Client that created the entity.
* `related_operation_list`: List of related Operations. These are the Operations which might need to be given access to, along with the current Operation, for certain workflows to succeed.
* `associated_endpoint_list`: List of associated endpoint objects for the Operation.

### associated_endpoint_list
* `api_version`: Version of the API for the provided associated endpoint.
* `endpoint_url`: Endpoint URL.
* `http_method`: HTTP method for the provided associated endpoint.

See detailed information in [Nutanix Operations](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0.b1).