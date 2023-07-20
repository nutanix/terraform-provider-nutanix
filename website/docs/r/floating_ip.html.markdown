---
layout: "nutanix"
page_title: "NUTANIX: nutanix_floating_ip"
sidebar_current: "docs-nutanix-resource-floating-ip"
description: |-
  Create Floating IPs .
---

# nutanix_floating_ip

Provides Nutanix resource to create Floating IPs. 

## Example Usage

## create Floating IP with External Subnet UUID

```hcl
resource "nutanix_floating_ip" "fip1" {
    external_subnet_reference_uuid = "{{ext_sub_uuid}}"
}
```

## create Floating IP with vpc name with external subnet name

```hcl
resource "nutanix_floating_ip" "fip2" {
    external_subnet_reference_name = "{{ext_sub_name}}"
    vpc_reference_name= "{{vpc_name}}"
    private_ip = "{{ip_address}}"
}
```

## Argument Reference

The following arguments are supported:

* `external_subnet_reference_uuid` - (Optional) The reference to a subnet. Should not be used with {external_subnet_reference_name} .
* `external_subnet_reference_name` - (Optional) The reference to a subnet. Should not be used with 
{external_subnet_reference_uuid} . 
* `vm_nic_reference_uuid` - (Optional) The reference to a vm_nic .
* `vpc_reference_uuid` - (Optional) The reference to a vpc. Should not be used with {vpc_reference_name}.
* `vpc_reference_name` - (Optional) The reference to a vpc. Should not be used with {vpc_reference_uuid}.
* `private_ip` - (Optional) Private IP with which floating IP is associated. Should be used with vpc_reference .

## Attributes Reference

The following attributes are exported:

* `metadata` - The floating_ips kind metadata.
* `api_version` - The version of the API.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: - subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

See detailed information in [Nutanix Floating IP](https://www.nutanix.dev/api_references/prism-central-v3/#/a9e06d3bba013-create-a-new-floating-ip).