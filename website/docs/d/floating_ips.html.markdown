---
layout: "nutanix"
page_title: "NUTANIX: nutanix_floating_ips"
sidebar_current: "docs-nutanix-datasource-floating-ips"
description: |-
   Provides a datasource to retrieve list of all floating ips.
---

# nutanix_floating_ips

Provides a datasource to retrieve all the floating IPs .

## Example Usage

```hcl
    data "nutanix_floating_ips" "test"{ }
```

## Attribute Reference
The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of Floating IPs. 

### Entities

The entities attribute element contains the following attributes:

* `metadata`: - The floating_ip kind metadata.
* `status` - Floating IP output status
* `spec` - Floating IP spec

### spec
An intentful representation of a floating_ip spec

* `resources` - Floating IP Resources. 

### status
An intentful representation of a floating_ip status

* `state` - The state of the floating_ip.
* `name` - floating_ip Name.
* `resources` - Floating IP allocation status.
* `execution_context` - Execution Context of Floating IP. 

### resources

* `external_subnet_reference` - The reference to a subnet 
* `floating_ip` - Private IP with which the floating IP is associated.
* `vm_nic_reference` - The reference to a vm_nic
* `vpc_reference` - The reference to a vpc

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: - subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Reference

The `vpc_reference`, `vm_nic_reference`, `external_subnet_reference` attributes supports the following:

* `kind`: - The kind name (Default value: project).
* `name`: - the name.
* `uuid`: - the UUID.

See detailed information in [Nutanix Floating IPs](https://www.nutanix.dev/api_references/prism-central-v3/#/5f65d87a3d014-get-a-list-of-existing-floating-i-ps).