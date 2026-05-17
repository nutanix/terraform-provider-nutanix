---
layout: "nutanix"
page_title: "NUTANIX: nutanix_rsyslog_server_v2"
sidebar_current: "docs-nutanix-datasource-rsyslog-server-v2"
description: |-
  Fetches the RSYSLOG server configuration identified by {extId} associated with the cluster identified by {clusterExtId}.
---

# nutanix_rsyslog_server_v2

Fetches the RSYSLOG server configuration identified by {extId} associated with the cluster identified by {clusterExtId}.

## Example Usage

```hcl
data "nutanix_rsyslog_server_v2" "example" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id         = "00000000-0000-0000-0000-000000000001"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id` - (Required) Indicates the UUID of a cluster.
* `ext_id` - (Required) RSYSLOG server UUID.

## Attributes Reference

The following attributes are exported:

* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.
* `ip_address` - IP address of the RSYSLOG server.
* `server_name` - RSYSLOG server name.
* `port` - RSYSLOG server port.
* `network_protocol` - Network protocol for RSYSLOG server. Valid values are `UDP`, `TCP`, `RELP`.
* `modules` - List of modules registered to RSYSLOG server.

### links

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object that is returned by the URL.

### ip_address

* `ipv4` - IPv4 address.
* `ipv6` - IPv6 address.

### ipv4, ipv6

* `value` - The IP address value.
* `prefix_length` - The prefix length of the network to which this host address belongs.

### modules

* `name` - Module name. Valid values are `CASSANDRA`, `CEREBRO`, `CURATOR`, `GENESIS`, `PRISM`, `STARGATE`, `SYSLOG_MODULE`, `ZOOKEEPER`, `UHARA`, `LAZAN`, `API_AUDIT`, `AUDIT`, `CALM`, `EPSILON`, `ACROPOLIS`, `MINERVA_CVM`, `FLOW`, `FLOW_SERVICE_LOGS`, `LCM`, `APLOS`, `NCM_AIOPS`.
* `log_severity_level` - Log severity level. Valid values are `EMERGENCY`, `ALERT`, `CRITICAL`, `ERROR`, `WARNING`, `NOTICE`, `INFO`, `DEBUG`.
* `should_log_monitor_files` - Option to log, monitor/output files of a module.
