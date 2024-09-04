---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vpc"
sidebar_current: "docs-nutanix-resource-vpc"
description: |-
  Create Virtual Private Cloud .
---

# nutanix_vpc

Provides Nutanix resource to create VPC.

## Example Usage

### vpc creation with external subnet name

```hcl
resource "nutanix_vpc" "vpc" {
  name = "testtNew-1"

  external_subnet_reference_name = [
    "test-Ext1",
    "test-ext2"
  ]

  common_domain_name_server_ip_list{
          ip = "8.8.8.8"
  }
  common_domain_name_server_ip_list{
          ip = "8.8.8.9"
  }

  externally_routable_prefix_list{
    ip=  "192.43.0.0"
    prefix_length= 16
  }
}
```

### vpc creation with external subnet uuid

```hcl
resource "nutanix_vpc" "vpc" {
  name = "testtNew-1"

  external_subnet_reference_uuid = [
    "<subnet_uuid>"
  ]

  common_domain_name_server_ip_list{
          ip = "8.8.8.8"
  }

  externally_routable_prefix_list{
    ip=  "192.43.0.0"
    prefix_length= 16
  }
  externally_routable_prefix_list{
    ip=  "192.42.0.0"
    prefix_length= 16
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the VPC.
* `external_subnet_reference_uuid` - (Optional) List of external subnets uuid attached to this VPC. Should not be used with external_subnet_reference_name.
* `external_subnet_reference_name` - (Optional) List of external subnets name attached to this VPC. Should not be used with external_subnet_reference_uuid.
* `externally_routable_prefix_list` - (Optional) List Externally Routable IP Addresses. Required when external subnet with NoNAT is used.
* `common_domain_name_server_ip_list` - (Optional) List of domain name server IPs.

## externally_routable_prefix_list
Externally Routable IP Addresses

* `ip` - (Required) The name for the VPC.
* `prefix_length` - (Required) prefix length.


## common_domain_name_server_ip_list
List of domain name server IPs

* `ip` - (Required) ip address.


## Attributes Reference

The following attributes are exported:

* `metadata` - The vpc kind metadata.
* `api_version` - The version of the API.
* `external_subnet_list_status` - Status of List of external subnets attached to this VPC

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: - subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

See detailed information in [Nutanix VPC](https://www.nutanix.dev/api_references/prism-central-v3/#/1b537be26b12f-create-a-new-vpc).
