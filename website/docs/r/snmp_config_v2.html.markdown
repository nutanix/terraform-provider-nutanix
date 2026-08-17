---
layout: "nutanix"
page_title: "NUTANIX: nutanix_snmp_config_v2"
sidebar_current: "docs-nutanix-resource-snmp-config-v2"
description: |-
  Manages the SNMP configuration of the cluster identified by {clusterExtId}: the cluster-wide enabled flag and a single SNMP transport entry.
---

# nutanix_snmp_config_v2

Manages the SNMP configuration of the cluster identified by `cluster_ext_id`.
The resource is polymorphic — depending on which arguments you set it manages
either the cluster-wide SNMP enabled flag, or a single SNMP transport entry on
the cluster:

* **Status mode** (only `is_enabled` is set): the resource calls
  `UpdateSnmpStatus` on create / update and is a no-op on destroy. The cluster's
  global SNMP toggle is the only thing managed in this mode.
* **Transport mode** (`port` and `protocol` are set): the resource calls
  `AddSnmpTransport` on create and `RemoveSnmpTransport` on destroy. `port` and
  `protocol` are immutable; changing either forces resource replacement.

The SNMP traps and SNMP users associated with a cluster are managed by the
dedicated [`nutanix_snmp_trap_v2`](snmp_trap_v2.html) and
[`nutanix_snmp_user_v2`](snmp_user_v2.html) resources.

## Example Usage — Status mode

```hcl
resource "nutanix_snmp_config_v2" "status" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  is_enabled     = true
}
```

## Example Usage — Transport mode

```hcl
resource "nutanix_snmp_config_v2" "transport" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  port           = 161
  protocol       = "UDP"
}
```

You can declare both modes side-by-side as separate resource instances on the
same cluster — the resource only mutates the attributes you set on each
instance.

## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: - (Required, ForceNew) Indicates the UUID of a cluster.
* `is_enabled`: - (Optional, Computed) SNMP status. Set to `true` to enable
  SNMP on the cluster, `false` to disable. When omitted, the resource only
  reads the current value back and never issues `UpdateSnmpStatus`.
* `port`: - (Optional, ForceNew, RequiredWith `protocol`) SNMP transport port.
  When set together with `protocol` the resource manages a single SNMP
  transport on the cluster. Both fields are immutable; changing either forces
  resource replacement.
* `protocol`: - (Optional, ForceNew, RequiredWith `port`) SNMP transport
  protocol. One of `TCP`, `TCP6`, `UDP`, `UDP6`.

## Attribute Reference

The following attributes are exported in addition to the arguments above:

* `ext_id`: - A globally unique identifier of an instance that is suitable for
  external consumption.
* `tenant_id`: - A globally unique identifier that represents the tenant that
  owns this entity. The system automatically assigns it, and it is immutable
  from an API consumer perspective.
* `links`: - A HATEOAS style link for the response. Each link contains a
  user-friendly name identifying the link and an address for retrieving the
  particular resource.

### Links

The `links` block exports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object
  that is returned by the URL. The unique value of `self` identifies the URL
  for the object.

## Import

`nutanix_snmp_config_v2` resources do not support import. The cluster's SNMP
config is a singleton with no per-mode server-assigned identifier, and the
Terraform id is synthesized at create time as either `<cluster_ext_id>:status`
(status mode) or `<cluster_ext_id>:<port>:<protocol>` (transport mode).

See detailed information in [Nutanix SNMP Config](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.3).
