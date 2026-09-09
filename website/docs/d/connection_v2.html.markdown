---
layout: "nutanix"
page_title: "NUTANIX: nutanix_connection_v2"
sidebar_current: "docs-nutanix-datasource-connection-v2"
description: |-
  Returns details of a hardware provider connection identified by its external ID.
---

# nutanix_connection_v2

Returns details of a hardware provider connection identified by its external ID.

## Example

```hcl
data "nutanix_connection_v2" "example" {
  hardware_provider_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id                   = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

* `hardware_provider_ext_id`: (Required) External ID of the hardware provider.
* `ext_id`: (Required) External ID of the hardware provider connection.

## Attributes Reference

The following attributes are exported:

* `name`: Name of the connection.
* `region`: Region for the connection.
* `deployment_type`: Type of deployment for the hardware provider connection.
* `created_time`: Creation time of the connection.
* `access_details`: Access details including endpoint and authentication information.
* `links`: A HATEOAS style list of links for the response.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
