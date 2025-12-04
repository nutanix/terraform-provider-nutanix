---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pbrs"
sidebar_current: "docs-nutanix-datasource-pbrs"
description: |-
    This operation retrieves a list of all the policy based routing.
---

# nutanix_pbr

Provides a datasource to retrieve all the pbrs.

## Example Usage

```hcl
    data "nutanix_pbrs" "test"{ }
```

## Attribute Reference
The following attributes are exported:

* `api_version`: version of the API
* `entities`: List of PBRs. 

### Entities

The entities attribute element contains the following attributes:

* `metadata`: - The routing policies kind metadata.
* `status` - PBR output status
* `spec` - PBR spec

### spec

* `name` - Name of PBR
* `resources` - PBR resources

### status

* `name` - The name of the PBR
* `state` - The state of the PBR
* `resources` - PBR resources status
* `execution_context` - Execution Context of PBR. 

### resources

* `is_bidirectional` - Policy in reverse direction.
* `vpc_reference` - Reference to VPC
* `destination` - destination address of an IP.
* `source` - source address of an IP. 
* `priority` - priority of routing policy
* `protocol_parameters` - Routing policy IP protocol parameters
* `action` - Routing policy action
* `protocol_type` - Protocol type of routing policy

### source , destination
source/destination address of an IP.

*`address` - address type of source.
*`subnet_ip` - IP subnet provided as an address.
*`prefix_length` - prefix length of provided subnet. 

### protocol_parameters
Routing policy IP protocol parameters

*`tcp` -  TCP parameters in routing policy
*`udp` -  UDP parameters in routing policy
*`icmp` -  ICMP parameters in routing policy.
*`protocol_number` - Protocol number in routing policy

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

The `vpc_reference`  attributes supports the following:

* `kind`: - The kind name (Default value: project).
* `name`: - the name.
* `uuid`: - the UUID.

See detailed information in [Nutanix Policy Based Routings](https://www.nutanix.dev/api_references/prism-central-v3/#/81005f2996866-get-a-list-of-existing-routing-policies).