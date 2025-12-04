---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vpc"
sidebar_current: "docs-nutanix-datasource-vpc"
description: |-
   This operation retrieves a vpc based on the input parameters.
---

# nutanix_vpc

Provides a datasource to retrieve VPC with vpc_uuid or vpc_name .

## Example Usage

```hcl
data "nutanix_vpc" "test1"{
    vpc_uuid = <vpc_uuid>
}

data "nutanix_vpc" "test2"{
    vpc_name = <vpc_name>
}
```

## Argument Reference

The following arguments are supported:

* `vpc_uuid` - vpc UUID
* `vpc_name` - vpc Name

## Attribute Reference

The following attributes are exported:

* `metadata`: - The vpc kind metadata.
* `api_version` - The version of the API.
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

See detailed information in [Nutanix VPC](https://www.nutanix.dev/api_references/prism-central-v3/#/1f75cfa326e6c-get-a-existing-vpc).