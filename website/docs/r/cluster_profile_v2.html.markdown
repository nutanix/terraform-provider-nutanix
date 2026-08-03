---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster_profile_v2"
sidebar_current: "docs-nutanix-resource-cluster-profile-v2"
description: |-
  Creates a cluster profile with the settings provided in the request body.
---

# nutanix_cluster_profile_v2

## Example Usage

```hcl
resource "nutanix_cluster_profile_v2" "example" {
  name              = "tf-cluster-profile"
  description       = "Example Cluster Profile created via Terraform"
  allowed_overrides = ["NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG"]

  name_server_ip_list {
    ipv4 {
      value = "240.29.254.180"
    }
    ipv6 {
      value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
    }
  }
  ntp_server_ip_list {
    ipv4 {
      value = "240.29.254.180"
    }
    ipv6 {
      value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
    }
    fqdn {
      value = "ntp.example.com"
    }
  }
  smtp_server {
    email_address = "email@example.com"
    type          = "SSL"
    server {
      ip_address {
        ipv4 {
          value = "240.29.254.180"
        }
        ipv6 {
          value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
        }
        fqdn {
          value = "smtp.example.com"
        }
      }
    }
  }
  nfs_subnet_white_list = ["10.110.106.45/255.255.255.255"]
  snmp_config {
    is_enabled = true
    users {
      username  = "snmpuser1"
      auth_type = "MD5"
      auth_key  = "Example_SNMP_user_authentication_key"
      priv_type = "DES"
      priv_key  = "Example_SNMP_user_encryption_key"
    }
    transports {
      protocol = "UDP"
      port     = 21
    }
    traps {
      address {
        ipv4 {
          value         = "240.29.254.180"
          prefix_length = 24
        }
        ipv6 {
          value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
        }
      }
      username         = "trapuser"
      protocol         = "UDP"
      port             = 59
      should_inform    = false
      engine_id        = "0x1234567890abcdef12"
      version          = "V2"
      receiver_name    = "trap-receiver"
      community_string = "snmp-server community public RO 192.168.1.0 255.255.255.0"
    }
  }
  rsyslog_server_list {
    server_name      = "exampleServer1"
    port             = 29
    network_protocol = "UDP"
    ip_address {
      ipv4 {
        value = "240.29.254.180"
      }
      ipv6 {
        value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"
      }
    }
    modules {
      name                     = "CASSANDRA"
      log_severity_level       = "EMERGENCY"
      should_log_monitor_files = true
    }
    modules {
      name                     = "CURATOR"
      log_severity_level       = "ERROR"
      should_log_monitor_files = false
    }
  }
  pulse_status {
    is_enabled          = false
    pii_scrubbing_level = "DEFAULT"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) Name of the cluster profile.
* `description`: - (Optional) Detailed description of a cluster profile.
* `allowed_overrides`: - (Optional) Indicates if a configuration of attached clusters can be skipped from monitoring.

    | Enum                      | Description                                |
    |---------------------------|--------------------------------------------|
    | NFS_SUBNET_WHITELIST_CONFIG | NFS subnet whitelist configuration         |
    | NTP_SERVER_CONFIG           | NTP server configuration                   |
    | SNMP_SERVER_CONFIG          | SNMP server configuration                  |
    | SMTP_SERVER_CONFIG          | SMTP server configuration                  |
    | PULSE_CONFIG                | Pulse status for a cluster                 |
    | NAME_SERVER_CONFIG          | Name server configuration                  |
    | RSYSLOG_SERVER_CONFIG       | RSYSLOG server configuration               |
    | HTTP_PROXY_CONFIG           | HTTP proxy configuration                   |
    | FAULT_TOLERANCE_CONFIG      | Fault tolerance configuration              |
    | REBUILD_RESERVATION_CONFIG  | Rebuild reservation configuration          |
    | RESILIENT_CAPACITY_WARNING_THRESHOLD_CONFIG | Resilient capacity warning threshold |

* `name_server_ip_list`: - (Optional) List of name servers on a cluster. This is a part of payload for both clusters create and update operations. Currently, only IPv4 address and FQDN (fully qualified domain name) values are supported for the create operation.

* `ntp_server_ip_list`: - (Optional) List of NTP servers on a cluster. This is a part of payload for both cluster create and update operations. Currently, only IPv4 address and FQDN (fully qualified domain name) values are supported for the create operation.

* `smtp_server`: - (Optional) SMTP servers on a cluster. This is part of payload for cluster update operation only.

* `nfs_subnet_white_list`: - (Optional) NFS subnet allowlist addresses. This is part of the payload for cluster update operation only.

* `snmp_config`: - (Optional) SNMP information.

* `rsyslog_server_list`: - (Optional) RSYSLOG Server.

* `pulse_status`: - (Optional) Pulse status for a cluster.

* `ntp_server_config_list`: - (Optional) List of NTP server configurations with encryption authentication.

* `http_proxy_config`: - (Optional) HTTP proxy configuration for the cluster.

* `fault_tolerance_config`: - (Optional) Fault tolerant state of cluster.

* `rebuild_reservation_config`: - (Optional) Rebuild capacity reservation configuration.

* `resilient_capacity_warning_threshold_config`: - (Optional) Resilient capacity warning threshold configuration.

### Name Server IP List

The name_server_ip_list attribute supports the following:

* `ipv4`: - (Optional) ip v4 address params.

  * `value`: - (Required) Ip V4 address.
  * `prefix_length`: - (Optional, default 32) The prefix length of the network to which this host IPv4 address belongs.

* `ipv6`: - (Optional) ip v6 address params.

  * `value`: - (Required) Ip V6 address.
  * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

### NTP Server IP List

The ntp_server_ip_list attribute supports the following:

* `ipv4`: - (Optional) ip address params.

  * `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
  * `value`: - (Required) Ip address.

* `ipv6`: - (Optional) Ip address params.

  * `value`: - (Required) Ip address.
  * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

* `fqdn`: - (Optional) A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

  * `value`: - (Required) FQDN value.

### SMTP Server

The smtp_server attribute supports the following:

* `email_address` (Required) SMTP email address.

* `server` (Required) SMTP network details.

  * `ip_address` An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

    * `ipv4`: - ip address params.

      * `value`: - (Required) Ip address.
      * `prefix_length`: - (Optional, default 32) The prefix length of the network to which this host IPv4 address belongs.

    * `ipv6`: - Ip address params.

      * `value`: - (Required) Ip address.
      * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

    * `fqdn`: - A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

      * `value`: - (Required) FQDN value.

  * `port` SMTP port.
  * `username` SMTP server user name.
  * `password` SMTP server password.

* `type` Type of SMTP server.

    | Enum      | Description                |
    |-----------|----------------------------|
    | PLAIN     | Plain type SMTP server     |
    | STARTTLS  | Start TLS type SMTP server |
    | SSL       | SSL type SMTP server       |

### SNMP Config

The snmp_config attribute supports the following:

* `is_enabled` (Optional) - (Boolean) SNMP status. Whether SNMP is enabled.

* `users` (Optional) - (List) SNMP user information. Each user object supports:

  * `username` (Required) - (String, max 64 chars) SNMP username. Required for SNMP trap v3 version.
  * `auth_type` (Required) - (String) SNMP user authentication type. Allowed values:

      | Enum      | Description                  |
      |-----------|------------------------------|
      | SHA       | SHA SNMP authentication      |
      | MD5       | MD5 SNMP authentication      |
      | SHA224    | SHA-224 SNMP authentication  |
      | SHA256    | SHA-256 SNMP authentication  |
      | SHA384    | SHA-384 SNMP authentication  |
      | SHA512    | SHA-512 SNMP authentication  |

  * `auth_key` (Required) - (String) SNMP user authentication key (must not contain single quotes).
  * `priv_type` (Optional) - (String) SNMP user encryption type. Allowed values:

      | Enum      | Description        |
      |-----------|--------------------|
      | DES       | DES SNMP key       |
      | AES       | AES SNMP key       |
      | AES192    | AES-192 SNMP key   |
      | AES256    | AES-256 SNMP key   |

  * `priv_key` (Optional) - (String) SNMP user encryption key (must not contain single quotes).

* `transports` (Optional) - (List) SNMP transport details. Each transport object supports:

  * `protocol` (Required) - (String) SNMP protocol type. Allowed values:

      | Enum      | Description        |
      |-----------|--------------------|
      | UDP       | UDP protocol       |
      | TCP       | TCP protocol       |
      | UDP6      | UDP6 protocol      |
      | TCP6      | TCP6 protocol      |

  * `port` (Required) - (Integer) SNMP port.

* `traps` (Optional) - (List) SNMP trap details. Each trap object supports:

  * `address` (Required) An unique address block that supports:

    * `ipv4`: - ip address params.

      * `value`: - (Required) Ip address.
      * `prefix_length`: - (Optional, default 32) The prefix length of the network to which this host IPv4 address belongs.

    * `ipv6`: - Ip address params.

      * `value`: - (Required) Ip address.
      * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

  * `username` (Optional) - (String, max 64 chars) SNMP username. Required for SNMP trap v3 version.
  * `protocol` (Optional) - (String) SNMP protocol type. Allowed values:

      | Enum      | Description        |
      |-----------|--------------------|
      | UDP       | UDP protocol       |
      | TCP       | TCP protocol       |
      | UDP6      | UDP6 protocol      |
      | TCP6      | TCP6 protocol      |

  * `port` (Optional) - (Integer) SNMP port.
  * `should_inform` (Optional) - (Boolean) SNMP inform mode status.
  * `engine_id` (Optional) - (String) SNMP engine ID (hexadecimal string, e.g. 0x12345678).
  * `version` (Required) - (String) SNMP version. Allowed values:

      | Enum      | Description        |
      |-----------|--------------------|
      | V2        | V2 SNMP version    |
      | V3        | V3 SNMP version    |

  * `receiver_name` (Optional) - (String, max 64 chars) SNMP receiver name.
  * `community_string` (Optional) - (String) Community string (plaintext) for SNMP version 2.0.

### RSYSLOG Server List

 The rsyslog_server_list attribute supports the following:

* `server_name` - (Required) Name of the RSYSLOG server.
* `port` - (Required) Port number for the RSYSLOG server.
* `network_protocol` - (Required) Network protocol for the RSYSLOG server. Allowed values:

    | Enum      | Description        |
    |-----------|--------------------|
    | UDP       | UDP protocol       |
    | TCP       | TCP protocol       |
    | RELP      | RELP protocol      |

* `ip_address` - (Required) IP address of the RSYSLOG server.

  * `ipv4`: - (Optional) ip address params.

    * `value`: - (Required) Ip address.
    * `prefix_length`: - (Optional, default 32) The prefix length of the network to which this host IPv4 address belongs.

  * `ipv6`: - (Optional) Ip address params.

    * `value`: - (Required) Ip address.
    * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

* `modules` - (Optional) List of modules for the RSYSLOG server. Each module object supports:

  * `name` - (Required) Name of the module. Allowed values:

    | Enum      | Description        |
    |-----------|--------------------|
    | AUDIT     | Audit module       |
    | CALM      | Calm module       |
    | MINERVA_CVM      | Minerva CVM module       |
    | STARGATE      | Stargate module       |
    | FLOW_SERVICE_LOGS      | Flow service logs module       |
    | SYSLOG_MODULE      | Syslog module       |
    | CEREBRO      | Cerebro module       |
    | API_AUDIT      | API audit module       |
    | GENESIS      | Genesis module       |
    | PRISM      | Prism module       |
    | ZOOKEEPER      | Zookeeper module       |
    | FLOW      | Flow module       |
    | EPSILON      | Epsilon module       |
    | ACROPOLIS      | Acropolis module       |
    | UHARA      | Uhara module       |
    | LCM      | LCM module       |
    | APLOS      | Aplos module       |
    | NCM_AIOPS      | NCM AIOPS module       |
    | CURATOR      | Curator module       |
    | CASSANDRA      | Cassandra module       |
    | LAZAN      | Lazan module       |

  * `log_severity_level` - (Required) Log severity level for the module. Allowed values:

    | Enum      | Description        |
    |-----------|--------------------|
    | EMERGENCY     | Emergency level       |
    | NOTICE      | Notice level       |
    | ERROR      | Error level       |
    | ALERT      | Alert level       |
    | INFO      | Info level       |
    | WARNING      | Warning level       |
    | DEBUG      | Debug level       |
    | CRITICAL      | Critical level       |

  * `should_log_monitor_files` - (Optional, default true) Boolean flag to indicate if log monitor files should be logged.
  
### Pulse Status

The pulse_status attribute supports the following:

* `is_enabled`: - (Optional) Flag to indicate if pulse is enabled or not.
* `pii_scrubbing_level`: - (Optional) PII scrubbing level.

    | Enum     | Description                                                                                    |
    |----------|------------------------------------------------------------------------------------------------|
    | ALL      | Scrub All PII Information from Pulse including data like entity names and IP addresses        |
    | DEFAULT  | Default PII Scrubbing level. Data like entity names and IP addresses will not be scrubbed from Pulse |

### NTP Server Config List

The ntp_server_config_list attribute supports the following:

* `ntp_server_address`: - (Required) NTP server address (supports IPv4, IPv6, and FQDN).
* `encryption_algorithm`: - (Optional) Encryption algorithm used for NTP server authentication.

    | Enum     | Description         |
    |----------|---------------------|
    | SHA256   | SHA-256 algorithm   |
    | SHA384   | SHA-384 algorithm   |
    | SHA512   | SHA-512 algorithm   |

* `encryption_key`: - (Optional, Sensitive) Encryption key in hexadecimal format used for NTP server authentication.
* `encryption_key_id`: - (Optional) Encryption key ID used for NTP server authentication.

### HTTP Proxy Config

The http_proxy_config attribute supports the following:

* `proxy_list`: - (Optional) List of HTTP proxy server configurations. Each proxy object supports:
  * `name` (Required) - HTTP proxy server name.
  * `ip_address` (Optional) - IP address of the proxy server.
  * `port` (Optional) - HTTP proxy server port.
  * `username` (Optional) - HTTP proxy server username.
  * `password` (Optional, Sensitive) - HTTP proxy server password.
  * `proxy_types` (Optional) - List of HTTP proxy types. Allowed values:

      | Enum     | Description         |
      |----------|---------------------|
      | HTTP     | HTTP proxy          |
      | HTTPS    | HTTPS proxy         |
      | SOCKS    | SOCKS proxy         |

* `proxy_white_list`: - (Optional) Targets exempted from going through the configured HTTP proxy. Each whitelist entry supports:
  * `target` (Required) - Target identifier exempted from HTTP proxy.
  * `target_type` (Required) - Type of the whitelist target. Allowed values:

      | Enum                 | Description              |
      |----------------------|--------------------------|
      | IPV4_ADDRESS         | IPv4 address             |
      | IPV6_ADDRESS         | IPv6 address             |
      | IPV4_NETWORK_MASK    | IPv4 network mask        |
      | DOMAIN_NAME_SUFFIX   | Domain name suffix       |
      | HOST_NAME            | Host name                |

### Fault Tolerance Config

The fault_tolerance_config attribute supports the following:

* `desired_cluster_fault_tolerance`: - (Optional) Desired cluster fault tolerance level. Allowed values:

    | Enum             | Description                         |
    |------------------|-------------------------------------|
    | CFT_0N_AND_0D    | No fault tolerance                  |
    | CFT_1N_OR_1D     | Tolerate 1 node or 1 disk failure   |
    | CFT_2N_OR_2D     | Tolerate 2 node or 2 disk failures  |
    | CFT_1N_AND_1D    | Tolerate 1 node and 1 disk failure  |
    | CFT_1N_OR_2D     | Tolerate 1 node or 2 disk failures  |

### Rebuild Reservation Config

The rebuild_reservation_config attribute supports the following:

* `is_rebuild_reservation_enabled`: - (Optional, default `false`) Enable rebuild capacity reservation to maintain free space for fault tolerance recovery operations.

### Resilient Capacity Warning Threshold Config

The resilient_capacity_warning_threshold_config attribute supports the following:

* `resilient_capacity_warning_threshold_percentage`: - (Optional, default `75`) Warning threshold percentage for resilient storage capacity.

## Import

This helps to manage existing entities which are not created through terraform. Cluster profile can be imported using the `UUID`.  eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_cluster_profile_v2" "import_cluster_profile" {}

// execute this cli command
terraform import nutanix_cluster_v2.import_cluster <UUID>
```

See detailed information in [Nutanix Create Cluster Profile V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.2#tag/ClusterProfiles/operation/createClusterProfile).
