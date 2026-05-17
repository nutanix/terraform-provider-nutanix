---
layout: "nutanix"
page_title: "NUTANIX: nutanix_rsyslog_server_v2"
sidebar_current: "docs-nutanix-resource-rsyslog-server-v2"
description: |-
  Adds RSYSLOG server configuration to the cluster identified by {clusterExtId}.
---

# nutanix_rsyslog_server_v2

Adds RSYSLOG server configuration to the cluster identified by {clusterExtId}. Supports create, read, update, and delete operations.

## Example Usage

```hcl
resource "nutanix_rsyslog_server_v2" "example" {
  cluster_ext_id   = "00000000-0000-0000-0000-000000000000"
  server_name      = "example-rsyslog-server"
  port             = 514
  network_protocol = "UDP"

  ip_address {
    ipv4 {
      value = "10.0.0.100"
    }
  }

  modules {
    name                     = "CASSANDRA"
    log_severity_level       = "INFO"
    should_log_monitor_files = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id` - (Required) Indicates the UUID of a cluster.
* `server_name` - (Required) RSYSLOG server name. This is the primary key of the entity and cannot be updated.
* `port` - (Required) RSYSLOG server port.
* `network_protocol` - (Required) Network protocol for RSYSLOG server. Valid values are `UDP`, `TCP`, `RELP`.
* `ip_address` - (Required) IP address of the RSYSLOG server.
* `modules` - (Optional) List of modules registered to RSYSLOG server.

### ip_address

* `ipv4` - (Optional) IPv4 address.
* `ipv6` - (Optional) IPv6 address.

### ipv4, ipv6

* `value` - (Required) The IP address value.
* `prefix_length` - (Optional) The prefix length of the network to which this host address belongs.

### modules

* `name` - (Required) Module name. Valid values are `CASSANDRA`, `CEREBRO`, `CURATOR`, `GENESIS`, `PRISM`, `STARGATE`, `SYSLOG_MODULE`, `ZOOKEEPER`, `UHARA`, `LAZAN`, `API_AUDIT`, `AUDIT`, `CALM`, `EPSILON`, `ACROPOLIS`, `MINERVA_CVM`, `FLOW`, `FLOW_SERVICE_LOGS`, `LCM`, `APLOS`, `NCM_AIOPS`.
* `log_severity_level` - (Required) Log severity level. Valid values are `EMERGENCY`, `ALERT`, `CRITICAL`, `ERROR`, `WARNING`, `NOTICE`, `INFO`, `DEBUG`.
* `should_log_monitor_files` - (Optional) Option to log, monitor/output files of a module.

## Attributes Reference

The following attributes are exported:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.

### links

* `href` - The URL at which the entity described by the link can be accessed.
* `rel` - A name that identifies the relationship of the link to the object that is returned by the URL.

## Import

Rsyslog Server can be imported using the format `cluster_ext_id/ext_id`:

```shell
terraform import nutanix_rsyslog_server_v2.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000001
```
