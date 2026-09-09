---
layout: "nutanix"
page_title: "NUTANIX: nutanix_connection_v2"
sidebar_current: "docs-nutanix-resource-connection-v2"
description: |-
  Manages a hardware provider connection, including endpoint and authentication details.
---

# nutanix_connection_v2

Manages the lifecycle of a hardware provider connection. A connection stores the endpoint and authentication information used to reach an external hardware provider. Create, Update, and Delete are asynchronous, task-based operations.

## Example

```hcl
resource "nutanix_connection_v2" "example" {
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
* `access_details`: (Required) Access details including endpoint and authentication information. See [Access Details](#access-details) below.

### Access Details

The `access_details` block supports the following:

* `auth`: (Optional) Authentication configuration for the connection. See [Auth](#auth) below.
* `endpoint`: (Optional) Endpoint configuration for the connection. See [Endpoint](#endpoint) below.

### Auth

Exactly one of the following authentication blocks should be set:

* `api_key_auth`: (Optional) API key authentication.
  * `api_key_id`: (Required) API key ID for authentication.
  * `api_key_secret`: (Optional) API key secret for authentication.
* `basic_auth`: (Optional) Username / password authentication.
  * `username`: (Required) Username for authentication.
  * `password`: (Optional) Password for authentication.

### Endpoint

Exactly one of the following endpoint blocks should be set:

* `url_endpoint`: (Optional) URL endpoint.
  * `url`: (Required) URL for the endpoint.
* `ip_range_endpoint`: (Optional) IP range endpoint.
  * `ip_ranges`: (Optional) List of IP ranges for the endpoint. Each range has `start_ip` and `end_ip`, each an `ipv4`/`ipv6` block with a `value` and optional `prefix_length`.
* `ip_address_endpoint`: (Optional) IP address endpoint.
  * `ip_addresses`: (Optional) List of IP addresses for the endpoint. Each entry supports `ipv4`, `ipv6`, and `fqdn` blocks.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id`: A globally unique identifier of the connection.
* `deployment_type`: Type of deployment for the hardware provider connection.
* `created_time`: Creation time of the connection.
* `links`: A HATEOAS style list of links for the response.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
