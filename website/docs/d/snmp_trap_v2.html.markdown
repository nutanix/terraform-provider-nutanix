---
layout: "nutanix"
page_title: "NUTANIX: nutanix_snmp_trap_v2"
sidebar_current: "docs-nutanix-datasource-snmp-trap-v2"
description: |-
  Fetches SNMP trap configuration details identified by {extId} associated with the cluster identified by {clusterExtId}.
---

# nutanix_snmp_trap_v2

Fetches SNMP trap configuration details identified by {extId} associated with the cluster identified by {clusterExtId}.

## Example Usage

```hcl
data "nutanix_snmp_trap_v2" "example" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id         = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: - (Required) Indicates the UUID of a cluster.
* `ext_id`: - (Required) SNMP trap UUID.

## Attribute Reference

The following attributes are exported:

* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `address`: - Address of the SNMP trap receiver. Includes `ipv4` and `ipv6` blocks.
* `community_string`: - Community string (plaintext) for SNMP version 2.0.
* `engine_id`: - SNMP engine Id.
* `port`: - SNMP port.
* `protocol`: - SNMP protocol. One of `TCP`, `TCP6`, `UDP`, `UDP6`.
* `receiver_name`: - SNMP receiver name.
* `should_inform`: - SNMP information status.
* `username`: - SNMP username. For SNMP trap v3 version, SNMP username is required parameter.
* `version`: - SNMP version. One of `V2`, `V3`.

### Address

The `address` block exports the following:

* `ipv4`: - IPv4 address. Contains:
  * `value`: - The IPv4 address of the host.
  * `prefix_length`: - The prefix length of the network to which this host IPv4 address belongs.
* `ipv6`: - IPv6 address. Contains:
  * `value`: - The IPv6 address of the host.
  * `prefix_length`: - The prefix length of the network to which this host IPv6 address belongs.

### Links

The `links` block exports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix Get SNMP Trap](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.3).
