---
layout: "nutanix"
page_title: "NUTANIX: nutanix_rsyslog_servers_v2"
sidebar_current: "docs-nutanix-datasource-rsyslog-servers-v2"
description: |-
  Lists the RSYSLOG server configurations associated with the cluster identified by {clusterExtId}.
---

# nutanix_rsyslog_servers_v2

Lists the RSYSLOG server configurations associated with the cluster identified by {clusterExtId}.

## Example

```hcl
data "nutanix_rsyslog_servers_v2" "example" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id` - (Required) Indicates the UUID of a cluster.

## Attribute Reference

The following attributes are exported:

* `rsyslog_servers` - List of RSYSLOG server configurations.

### Rsyslog Server

Each item in `rsyslog_servers` has the following attributes:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `server_name` - RSYSLOG server name.
* `port` - RSYSLOG server port.
* `network_protocol` - RSYSLOG server protocol type.
* `ip_address` - IP address of the RSYSLOG server.
  * `ipv4` - IPv4 address.
    * `value` - The IPv4 address of the host.
    * `prefix_length` - The prefix length of the network to which this host IPv4 address belongs.
  * `ipv6` - IPv6 address.
    * `value` - The IPv6 address of the host.
    * `prefix_length` - The prefix length of the network to which this host IPv6 address belongs.
* `modules` - List of modules registered to RSYSLOG server.
  * `name` - RSYSLOG module name.
  * `log_severity_level` - RSYSLOG module log severity level.
  * `should_log_monitor_files` - Option to log, monitor/output files of a module.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

See detailed information in [Nutanix Cluster Management RSYSLOG Server v4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.2).
