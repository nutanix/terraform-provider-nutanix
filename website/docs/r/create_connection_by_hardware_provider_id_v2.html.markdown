---
layout: "nutanix"
page_title: "NUTANIX: nutanix_create_connection_by_hardware_provider_id_v2"
sidebar_current: "docs-nutanix-resource-create-connection-by-hardware-provider-id-v2"
description: |-
  Creates a connection to an endpoint for a specific hardware provider.
---

# nutanix_create_connection_by_hardware_provider_id_v2

Creates a connection to an endpoint for a specific hardware provider. This is a standalone create-action resource. The full lifecycle (read/update/delete) of a connection is managed by [`nutanix_connection_v2`](connection_v2.html).

## Example

```hcl
resource "nutanix_create_connection_by_hardware_provider_id_v2" "example" {
  hardware_provider_ext_id = "00000000-0000-0000-0000-000000000000"
  name                     = "my-connection"
  region                   = "us-west"

  access_details {
    auth {
      basic_auth {
        username = "connection-user"
        password = "connection-password"
      }
    }
    endpoint {
      url_endpoint {
        url = "https://hardware-provider.example.com"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `hardware_provider_ext_id`: (Required) External ID of the hardware provider.
* `name`: (Required) Name of the connection.
* `region`: (Optional) Region for the connection.
* `access_details`: (Required) Access details including endpoint and authentication information. See the [`nutanix_connection_v2`](connection_v2.html) documentation for the full `access_details` structure.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id`: A globally unique identifier of the created connection.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
