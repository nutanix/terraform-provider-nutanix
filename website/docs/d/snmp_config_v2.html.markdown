---
layout: "nutanix"
page_title: "NUTANIX: nutanix_snmp_config_v2"
sidebar_current: "docs-nutanix-datasource-snmp-config-v2"
description: |-
  Fetches the full SNMP configuration of the cluster identified by {clusterExtId}: the enabled flag, transports, traps and users in one call.
---

# nutanix_snmp_config_v2

Fetches the full SNMP configuration of the cluster identified by
`cluster_ext_id`. Returns the cluster-wide enabled flag along with every
configured transport, trap and user in a single read — useful when you need to
introspect the SNMP surface of a cluster without paging through individual
trap / user GETs.

## Example Usage

```hcl
data "nutanix_snmp_config_v2" "current" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
}

output "snmp_traps" {
  value = data.nutanix_snmp_config_v2.current.traps
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: - (Required) Indicates the UUID of a cluster.

## Attribute Reference

The following attributes are exported:

* `ext_id`: - A globally unique identifier of an instance that is suitable for
  external consumption.
* `tenant_id`: - A globally unique identifier that represents the tenant that
  owns this entity.
* `links`: - A HATEOAS style link for the response. Each link contains a
  user-friendly name identifying the link and an address for retrieving the
  particular resource.
* `is_enabled`: - SNMP status on the cluster.
* `transports`: - List of SNMP transports configured on the cluster. Each
  entry exports `port` and `protocol`.
* `traps`: - List of SNMP traps configured on the cluster. See **Trap** below.
* `users`: - List of SNMP users configured on the cluster. See **User** below.

### Trap

Each entry in `traps` exports:

* `ext_id`: - A globally unique identifier of the trap entity.
* `tenant_id`: - Tenant identifier.
* `links`: - HATEOAS links for the trap.
* `address`: - Address of the SNMP trap receiver. Includes nested `ipv4` /
  `ipv6` blocks, each with `value` and `prefix_length`.
* `community_string`: - Community string (plaintext) for SNMP V2 traps.
* `engine_id`: - SNMP engine Id.
* `port`: - SNMP port.
* `protocol`: - SNMP protocol. One of `TCP`, `TCP6`, `UDP`, `UDP6`.
* `receiver_name`: - SNMP receiver name.
* `should_inform`: - SNMP information status.
* `username`: - SNMP username (required for V3 traps).
* `version`: - SNMP version. One of `V2`, `V3`.

### User

Each entry in `users` exports:

* `ext_id`: - A globally unique identifier of the user entity.
* `tenant_id`: - Tenant identifier.
* `links`: - HATEOAS links for the user.
* `username`: - SNMP username.
* `auth_type`: - SNMP user authentication type. One of `MD5`, `SHA`, `SHA224`, `SHA256`, `SHA384`, `SHA512`.
* `priv_type`: - SNMP user encryption type. One of `DES`, `AES`, `AES192`, `AES256`.

### Links

The `links` block exports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object
  that is returned by the URL. The unique value of `self` identifies the URL
  for the object.

See detailed information in [Nutanix Get SNMP Config](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.3).
