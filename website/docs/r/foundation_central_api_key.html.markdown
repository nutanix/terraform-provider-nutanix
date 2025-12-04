---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_api_keys"
sidebar_current: "docs-nutanix-resource-foundation-central-api-keys"
description: |-
  Create a new api key which will be used by remote nodes to authenticate with Foundation Central .
---

# nutanix_foundation_central_api_keys

Provides a resource to create a new API key for nodes registration with Foundation Central. 

## Example Usage

```hcl
resource "nutanix_foundation_central_api_keys" "new_api_key" {
	alias = "<NAME-FOR-API-KEY>"
}
```


## Argument Reference

The following arguments are supported:

* `alias`: - (Required) Alias for the api key to be created.

## Attributes Reference

The following attributes are exported:

* `created_timestamp`: Time when the api key was created.
* `alias`: Alias of the api key.
* `key_uuid`: UUID of the api key.
* `api_key`: Api key in string format.
* `current_time`: Current time of Foundation Central.


See detailed information in [Nutanix Foundation Central Create an API Key](https://www.nutanix.dev/api_references/foundation-central/#/c2e963769f299-create-an-api-key).