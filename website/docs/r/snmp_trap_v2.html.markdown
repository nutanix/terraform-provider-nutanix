---
layout: "nutanix"
page_title: "NUTANIX: nutanix_snmp_trap_v2"
sidebar_current: "docs-nutanix-resource-snmp-trap-v2"
description: |-
  Manages an SNMP trap (V2 or V3) on the cluster identified by {clusterExtId}.
---

# nutanix_snmp_trap_v2

Manages an SNMP trap on the cluster identified by `cluster_ext_id`. The
resource handles the full lifecycle of a trap (`CreateSnmpTrap`,
`GetSnmpTrapById`, `UpdateSnmpTrapById`, `DeleteSnmpTrapById`) and supports
both V2 and V3 traps:

* **V2 traps** authenticate via `community_string`.
* **V3 traps** reference an existing SNMP user on the same cluster by
  `username`. Provision the user first with
  [`nutanix_snmp_user_v2`](snmp_user_v2.html) — Terraform's dependency graph
  will order create / destroy correctly when the trap references the user.

`version` and `cluster_ext_id` are immutable. The upstream API has no
in-place mechanism for switching SNMP versions or moving a trap between
clusters; changing either attribute forces resource replacement. All other
attributes (`address`, `port`, `protocol`, `community_string`, `username`,
`engine_id`, `receiver_name`, `should_inform`) can be updated in place.

## Example Usage — V2 trap

```hcl
resource "nutanix_snmp_trap_v2" "v2" {
  cluster_ext_id   = "00000000-0000-0000-0000-000000000000"
  version          = "V2"
  community_string = "public"
  port             = 162
  protocol         = "UDP"

  address {
    ipv4 {
      value = "10.0.0.100"
    }
  }
}
```

## Example Usage — V3 trap referencing an SNMP user

```hcl
resource "nutanix_snmp_user_v2" "v3_user" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  username       = "tf-snmp-v3-user"
  auth_type      = "SHA"
  auth_key       = "auth-key-12345678"
  priv_type      = "AES"
  priv_key       = "priv-key-12345678"
}

resource "nutanix_snmp_trap_v2" "v3" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  version        = "V3"
  username       = nutanix_snmp_user_v2.v3_user.username
  port           = 163
  protocol       = "UDP"

  address {
    ipv4 {
      value = "10.0.0.101"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: - (Required, ForceNew) Indicates the UUID of a cluster.
* `version`: - (Required, ForceNew) SNMP version. One of `V2`, `V3`. Switching
  versions is not supported in place; changing this attribute forces resource
  replacement.
* `address`: - (Required) Address block of the SNMP trap receiver. Exactly one
  of `ipv4` or `ipv6` should be set. See **Address** below.
* `port`: - (Required) SNMP port.
* `protocol`: - (Required) SNMP protocol. One of `TCP`, `TCP6`, `UDP`, `UDP6`.
* `community_string`: - (Optional, Sensitive) Community string (plaintext) for
  SNMP V2 traps.
* `username`: - (Optional, Computed) SNMP username. Required when `version` is
  `V3` and must reference an existing
  [`nutanix_snmp_user_v2`](snmp_user_v2.html) on the same cluster.
* `engine_id`: - (Optional, Computed) SNMP engine Id.
* `receiver_name`: - (Optional, Computed) SNMP receiver name.
* `should_inform`: - (Optional, Computed) SNMP information status.

### Address

The `address` block supports the following nested blocks:

* `ipv4`: - (Optional) IPv4 address block. Contains:
  * `value`: - (Required) The IPv4 address of the host.
  * `prefix_length`: - (Optional, Computed) The prefix length of the network
    to which this host IPv4 address belongs.
* `ipv6`: - (Optional) IPv6 address block. Contains:
  * `value`: - (Required) The IPv6 address of the host.
  * `prefix_length`: - (Optional, Computed) The prefix length of the network
    to which this host IPv6 address belongs.

## Attribute Reference

The following attributes are exported in addition to the arguments above:

* `ext_id`: - A globally unique identifier of an instance that is suitable for
  external consumption.
* `tenant_id`: - A globally unique identifier that represents the tenant that
  owns this entity.
* `links`: - A HATEOAS style link for the response.

### Links

The `links` block exports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object
  that is returned by the URL. The unique value of `self` identifies the URL
  for the object.

## Import

`nutanix_snmp_trap_v2` supports import via the trap's server-assigned ext_id.
Note that the trap's `cluster_ext_id` must be set in the configuration prior
to import.

```shell
terraform import nutanix_snmp_trap_v2.example <trap-ext-id>
```

See detailed information in [Nutanix Create SNMP Trap](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.3#tag/Clusters/operation/createSnmpTrap).
