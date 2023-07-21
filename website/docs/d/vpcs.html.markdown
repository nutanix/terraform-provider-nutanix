---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vpcs"
sidebar_current: "docs-nutanix-datasource-vpcs"
description: |-
  This operation retrieves a list of all the vpcs.
---

# nutanix_vpcs

Provides a datasource to retrieve all the vpcs.

## Example Usage

```hcl
    data "nutanix_vpcs" "test"{ }
```

## Attribute Reference

The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of VPCs

### Entities

The entities attribute element contains the following attributes:

* `metadata`: - The vpc kind metadata.
* `status` - VPC output status
* `spec` - VPC input spec


### spec

* `name` - Name of VPC .
* `resources` - VPC resources .

### status

* `name` - The name of the VPC
* `state` - The state of the VPC
* `resources` - VPC resources status
* `execution_context` - Execution Context of VPC. 

### resources

* `external_subnet_list` - List of external subnets attached to this VPC.
* `externally_routable_prefix_list` - List of external routable ip and prefix . 
* `common_domain_name_server_ip_list` - List of domain name server IPs. 

### external_subnet_list

* `external_subnet_reference` - Reference to a subnet. 
* `external_ip_list` - List of external subnets attached to this VPC. Only present in VPC Status Resources .
* `active_gateway_node` - Active Gateway Node. Only present in VPC Status Resources. 

### externally_routable_prefix_list

* `ip` - ip address . 
* `prefix_length` - prefix length of routable ip .

### common_domain_name_server_ip_list

* `ip` - ip address of domain name server. 

#### active_gateway_node

* `host_reference` - Reference to host.
* `ip_address` - ip address. 

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

The `host_reference`, `external_subnet_reference`  attributes supports the following:

* `kind`: - The kind name (Default value: project).
* `name`: - the name.
* `uuid`: - the UUID.

See detailed information in [Nutanix VPCs](https://www.nutanix.dev/api_references/prism-central-v3/#/66f54e0e5ae08-get-a-list-of-existing-vp-cs).