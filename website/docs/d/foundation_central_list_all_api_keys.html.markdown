---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_list_api_keys"
sidebar_current: "docs-nutanix-datasource-foundation-central-list-api-keys"
description: |-
 List all the api keys created in Foundation Central.
---

# nutanix_foundation_central_list_api_keys

List all the api keys created in Foundation Central.

## Example Usage

```hcl
data "nutanix_foundation_central_list_api_keys" "api_keys_list" {}
```

## Argument Reference

No arguments are supported

## Attribute Reference

The following attributes are exported:

* `metadata`: List metadata output for all list apis.

### metadata
* `total_matches`: Total matches found.
* `length`: The number of records retrieved.
* `offset`: Offset from the start of the object list.

### api_keys
* `created_timestamp`: Time when the api key was created.
* `alias`: Alias of the api key.
* `key_uuid`: UUID of the api key.
* `api_key`: Api key in string format.
* `current_time`: Current time of Foundation Central.


See detailed information in [Nutanix Foundation Central List all the API keys](https://www.nutanix.dev/api_references/foundation-central/#/91806fd4d9abc-list-all-the-api-keys).