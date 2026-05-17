---
layout: "nutanix"
page_title: "NUTANIX: nutanix_rsyslog_server_v2"
sidebar_current: "docs-nutanix-resource-rsyslog-server-v2"
description: |-
  Adds RSYSLOG server configuration to the cluster identified by {clusterExtId}.
---

# nutanix_rsyslog_server_v2

Provides Nutanix resource to add RSYSLOG server configuration to a cluster. Update RSYSLOG server configuration except RSYSLOG server name as it is a primary key of the entity.

## Example

### Basic Example

```hcl
resource "nutanix_rsyslog_server_v2" "example" {
  cluster_ext_id   = "00000000-0000-0000-0000-000000000000"
  server_name      = "my-rsyslog-server"
  port             = 514
  network_protocol = "UDP"

  ip_address {
    ipv4 {
      value = "10.0.0.1"
    }
  }

  modules {
    name                     = "CASSANDRA"
    log_severity_level       = "INFO"
    should_log_monitor_files = true
  }
}
```

### Example with Multiple Modules

```hcl
resource "nutanix_rsyslog_server_v2" "multi_module" {
  cluster_ext_id   = "00000000-0000-0000-0000-000000000000"
  server_name      = "multi-module-rsyslog"
  port             = 6514
  network_protocol = "TCP"

  ip_address {
    ipv4 {
      value = "10.0.0.2"
    }
  }

  modules {
    name                     = "CASSANDRA"
    log_severity_level       = "INFO"
    should_log_monitor_files = true
  }

  modules {
    name                     = "STARGATE"
    log_severity_level       = "WARNING"
    should_log_monitor_files = false
  }

  modules {
    name                     = "GENESIS"
    log_severity_level       = "ERROR"
    should_log_monitor_files = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id` - (Required) Indicates the UUID of a cluster.
* `server_name` - (Required) RSYSLOG server name. This is the primary key of the entity and cannot be updated.
* `port` - (Required) RSYSLOG server port.
* `network_protocol` - (Required) RSYSLOG server protocol type. Valid values: `"UDP"`, `"TCP"`, `"RELP"`.
* `ip_address` - (Required) IP address of the RSYSLOG server.
* `modules` - (Optional) List of modules registered to RSYSLOG server.

### IP Address

The `ip_address` block supports the following:

* `ipv4` - (Optional) IPv4 address.
* `ipv6` - (Optional) IPv6 address.

### IPv4

The `ipv4` block supports the following:

* `value` - (Required) The IPv4 address of the host.
* `prefix_length` - (Optional) The prefix length of the network to which this host IPv4 address belongs.

### IPv6

The `ipv6` block supports the following:

* `value` - (Required) The IPv6 address of the host.
* `prefix_length` - (Optional) The prefix length of the network to which this host IPv6 address belongs.

### Modules

The `modules` block supports the following:

* `name` - (Required) RSYSLOG module name. Valid values: `"CASSANDRA"`, `"CEREBRO"`, `"CURATOR"`, `"GENESIS"`, `"PRISM"`, `"STARGATE"`, `"SYSLOG_MODULE"`, `"ZOOKEEPER"`, `"UHARA"`, `"LAZAN"`, `"API_AUDIT"`, `"AUDIT"`, `"CALM"`, `"EPSILON"`, `"ACROPOLIS"`, `"MINERVA_CVM"`, `"FLOW"`, `"FLOW_SERVICE_LOGS"`, `"LCM"`, `"APLOS"`, `"NCM_AIOPS"`.
* `log_severity_level` - (Required) RSYSLOG module log severity level. Valid values: `"EMERGENCY"`, `"ALERT"`, `"CRITICAL"`, `"ERROR"`, `"WARNING"`, `"NOTICE"`, `"INFO"`, `"DEBUG"`.
* `should_log_monitor_files` - (Optional) Option to log, monitor/output files of a module.

## Attribute Reference

The following attributes are exported:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

## Import

RSYSLOG server resources can be imported using the format `cluster_ext_id:ext_id`. For example:

```hcl
resource "nutanix_rsyslog_server_v2" "import_rsyslog" {}

// execute this command in cli
terraform import nutanix_rsyslog_server_v2.import_rsyslog <cluster_ext_id>:<ext_id>
```

See detailed information in [Nutanix Cluster Management RSYSLOG Server v4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.2).
