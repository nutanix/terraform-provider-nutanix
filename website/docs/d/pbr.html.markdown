---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pbr"
sidebar_current: "docs-nutanix-datasource-pbr"
description: |-
   Provides a datasource to retrieve Policy Based Routing with pbr_uuid .
---

# nutanix_pbr

Provides a datasource to retrieve PBR with pbr_uuid .

## Example Usage

```hcl
    data "nutanix_pbr" "test"{
        pbr_uuid = <pbr_uuid>
    }
```

## Argument Reference

The following arguments are supported:

* `pbr_uuid` - (Required) pbr UUID

## Attribute Reference

The following attributes are exported:

* `metadata`: - The routing policies kind metadata.
* `api_version` - The version of the API.
* `status` - PBR output status
* `spec` - PBR input spec

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

See detailed information in [Nutanix Policy Based Routing](https://www.nutanix.dev/api_references/prism-central-v3/#/3506dc2d5ec27-get-a-existing-routing-policy).