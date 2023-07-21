---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_api_keys"
sidebar_current: "docs-nutanix-datasource-foundation-central-api-keys"
description: |-
 Details of the api key.
---

# nutanix_foundation_central_api_keys

Get an api key given its UUID.

## Example Usage

```hcl
data "nutanix_foundation_central_api_keys" "api_keys_list" {
    key_uuid = "<KEY_UUID>"
}
```

## Argument Reference

*`key_uuid`: UUID of the key which needs to be fetched. 

## Attribute Reference

The following attributes are exported:

* `created_timestamp`: Time when the api key was created.
* `alias`: Alias of the api key.
* `key_uuid`: UUID of the api key.
* `api_key`: Api key in string format.
* `current_time`: Current time of Foundation Central.


See detailed information in [Nutanix Foundation Central Get an API key](https://www.nutanix.dev/api_references/foundation-central/#/92553fa628770-get-an-api-key).