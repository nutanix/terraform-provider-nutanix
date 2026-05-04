---
layout: "nutanix"
page_title: "NUTANIX: nutanix_add_snmp_transport_v2"
sidebar_current: "docs-nutanix-resource-add-snmp-transport-v2"
description: |-
  Adds transport ports and protocol details to the SNMP configuration associated with the cluster identified by {clusterExtId}.
---

# nutanix_add_snmp_transport_v2

Adds transport ports and protocol details to the SNMP configuration associated with the cluster identified by the cluster UUID.

## Example Usage

```hcl
resource "nutanix_add_snmp_transport_v2" "example" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  port           = 162
  protocol       = "UDP"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: -(Required) Indicates the UUID of a cluster.
* `port`: -(Optional) SNMP port.
* `protocol`: -(Optional) SNMP transport protocol. Valid values are:
  - `UDP` - UDP protocol.
  - `UDP6` - UDP6 protocol.
  - `TCP` - TCP protocol.
  - `TCP6` - TCP6 protocol.

## Attributes Reference

No additional attributes are exported.

See detailed information in [Nutanix Add SNMP Transport V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.2#tag/Clusters/operation/addSnmpTransport).
