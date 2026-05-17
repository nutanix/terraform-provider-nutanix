---
layout: "nutanix"
page_title: "NUTANIX: nutanix_rsyslog_server_v2"
sidebar_current: "docs-nutanix-datasource-rsyslog-server-v2"
description: |-
  Fetches the RSYSLOG server configuration identified by {extId} associated with the cluster identified by {clusterExtId}.
---

# nutanix_rsyslog_server_v2

Describes an RSYSLOG server configuration identified by {extId} associated with the cluster identified by {clusterExtId}.

## Example

```hcl
data "nutanix_rsyslog_server_v2" "example" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id         = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id` - (Required) Indicates the UUID of a cluster.
* `ext_id` - (Required) RSYSLOG server UUID.

## Attribute Reference

The following attributes are exported:

* `server_name` - RSYSLOG server name.
* `port` - RSYSLOG server port.
* `network_protocol` - RSYSLOG server protocol type.
* `ip_address` - IP address of the RSYSLOG server.
* `modules` - List of modules registered to RSYSLOG server.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

### IP Address

The `ip_address` attribute has the following:

* `ipv4` - IPv4 address.
  * `value` - The IPv4 address of the host.
  * `prefix_length` - The prefix length of the network to which this host IPv4 address belongs.
* `ipv6` - IPv6 address.
  * `value` - The IPv6 address of the host.
  * `prefix_length` - The prefix length of the network to which this host IPv6 address belongs.

### Modules

The `modules` attribute has the following:

* `name` - RSYSLOG module name.
* `log_severity_level` - RSYSLOG module log severity level.
* `should_log_monitor_files` - Option to log, monitor/output files of a module.

See detailed information in [Nutanix Cluster Management RSYSLOG Server v4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.2).
