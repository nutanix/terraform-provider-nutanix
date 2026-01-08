---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster_profiles_v2"
sidebar_current: "docs-nutanix-datasource-cluster-profiles-v2"
description: |-
  Fetches a list of cluster profile entities. A profile consists of different cluster settings like Network Time Protocol(NTP), Domain Name System(DNS), and so on.


---

# nutanix_cluster_profiles_v2

Fetches a list of cluster profile entities. A profile consists of different cluster settings like Network Time Protocol(NTP), Domain Name System(DNS), and so on.

## Example Usage

```hcl
data "nutanix_cluster_profile_v2" "list-cluster-profiles"{
}

data "nutanix_cluster_profiles_v2" "filtered-cluster-profiles"{
  filter = "clusterCount eq 62"
}

data "nutanix_cluster_profiles_v2" "paged-cluster-profiles"{
  page = 1
  limit = 10
}

data "nutanix_cluster_profiles_v2" "ordered-cluster-profiles"{
  order_by = "name"
}

data "nutanix_cluster_profiles_v2" "selected-cluster-profiles"{
  select = "name,description"
}
```

## Argument Reference

The following arguments are supported:

* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. The filter can be applied to the following fields:
  * clusterCount : `filter="clusterCount eq 62"`
  * createTime : `filter="createTime eq '2009-09-23T14:30:00-07:00'"`
  * createdBy : `filter="createdBy eq '984ff444-b6aa-4fa9-a8bd-37ef221cf1b0'"`
  * driftedClusterCount : `filter="driftedClusterCount eq 3"`
  * lastUpdateTime : `filter="lastUpdateTime eq '2009-09-23T14:30:00-07:00'"`
  * lastUpdatedBy : `filter="lastUpdatedBy eq '6debfc73-b4a9-4d68-a702-f0cfc1d721fb'"`
  * name : `filter="name eq 'Test Cluster Profile'"`

* `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
  * clusterCount : `orderby="clusterCount"`
  * createTime : `orderby="createTime desc"`
  * createdBy : `orderby="createdBy desc"`
  * driftedClusterCount : `orderby="driftedClusterCount desc"`
  * lastUpdateTime : `orderby="lastUpdateTime desc"`
  * lastUpdatedBy : `orderby="lastUpdatedBy desc"`
  * name : `orderby="name desc"`

* `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned.

  The following selection keys are supported:

  * allowedOverrides : `select="allowedOverrides"`
  * clusterCount : `select="clusterCount"`
  * clusters : `select="clusters"`
  * createTime : `select="createTime"`
  * createdBy : `select="createdBy"`
  * description : `select="description"`
  * driftedClusterCount : `select="driftedClusterCount"`
  * extId : `select="extId"`
  * lastUpdateTime : `select="lastUpdateTime"`
  * lastUpdatedBy : `select="lastUpdatedBy"`
  * links : `select="links"`
  * name : `select="name"`
  * nameServerIpList : `select="nameServerIpList"`
  * nfsSubnetWhitelist : `select="nfsSubnetWhitelist"`
  * ntpServerIpList : `select="ntpServerIpList"`
  * pulseStatus : `select="pulseStatus"`
  * rsyslogServerList : `select="rsyslogServerList"`
  * smtpServer : `select="smtpServer"`
  * snmpConfig : `select="snmpConfig"`
  * tenantId : `select="tenantId"`

## Attribute Reference

The following attribute are exported:

* `cluster_profiles`: - List of cluster profiles.

### Cluster Profiles

The cluster_profiles attribute supports the following:

* `tenant_id`: -  globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: -  A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`: - Name of the cluster profile.
* `description`: - Detailed description of a cluster profile.
* `create_time`: - Creation time of a cluster profile.
* `last_update_time`: - The last updated time of a cluster profile.
* `created_by`: - Details of the user who created the cluster profile.
* `last_updated_by`: - Details of the user who has recently updated the cluster profile.
* `cluster_count`: - Count of clusters associated to a cluster profile.
* `drifted_cluster_count`: - The count indicates the number of clusters associated with a cluster profile that has experienced drift. Drifted clusters are those in which the configuration differs from the defined profile. For example, the NTP server has different values on a cluster as compared to the profile it is attached.
* `clusters`: - Managed cluster information.
* `allowed_overrides`: - Indicates if a configuration of attached clusters can be skipped from monitoring.

    | Enum                      | Description                                |
    |---------------------------|--------------------------------------------|
    | NFS_SUBNET_WHITELIST_CONFIG | NFS subnet whitelist configuration         |
    | NTP_SERVER_CONFIG           | NTP server configuration                   |
    | SNMP_SERVER_CONFIG          | SNMP server configuration                  |
    | SMTP_SERVER_CONFIG          | SMTP server configuration                  |
    | PULSE_CONFIG                | Pulse status for a cluster                 |
    | NAME_SERVER_CONFIG          | Name server configuration                  |
    | RSYSLOG_SERVER_CONFIG       | RSYSLOG server configuration               |

* `name_server_ip_list`: - List of name servers on a cluster. This is a part of payload for both clusters create and update operations. Currently, only IPv4 address and FQDN (fully qualified domain name) values are supported for the create operation.
* `ntp_server_ip_list`: - List of NTP servers on a cluster. This is a part of payload for both cluster create and update operations. Currently, only IPv4 address and FQDN (fully qualified domain name) values are supported for the create operation.
* `smtp_server`: - SMTP servers on a cluster. This is part of payload for cluster update operation only.
* `nfs_subnet_white_list`: - NFS subnet allowlist addresses. This is part of the payload for cluster update operation only.
* `snmp_config`: - SNMP information.
* `rsyslog_server_list`: - RSYSLOG Server.
* `pulse_status`: - Pulse status for a cluster.

### Clusters

The clusters attribute supports the following:

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `is_compliant`: - Indicates if a cluster is compliant with the cluster profile.
* `last_synced_time`: - The last synced time of a cluster.
* `config_drifts`: - The configuration drifts of a cluster.
  
    | Enum                      | Description                                |
    |---------------------------|--------------------------------------------|
    | NFS_SUBNET_WHITELIST_CONFIG | NFS subnet whitelist configuration         |
    | NTP_SERVER_CONFIG           | NTP server configuration                   |
    | SNMP_SERVER_CONFIG          | SNMP server configuration                  |
    | SMTP_SERVER_CONFIG          | SMTP server configuration                  |
    | PULSE_CONFIG                | Pulse status for a cluster                 |
    | NAME_SERVER_CONFIG          | Name server configuration                  |
    | RSYSLOG_SERVER_CONFIG       | RSYSLOG server configuration               |

### Name Server IP List

The name_server_ip_list attribute supports the following:

* `ipv4`: - ip address params.
  * `value`: - (Required) Ip address.
  * `prefix_length`: - (Optional, default 32) The prefix length of the network to which this host IPv4 address belongs.

* `ipv6`: - Ip address params.

  * `value`: - (Required) Ip address.
  * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

### Ntp Server IP List

The ntp_server_ip_list attribute supports the following:

* `ipv4`: - (Optional) ip address params.

  * `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
  * `value`: - (Required) Ip address.

* `ipv6`: - (Optional) Ip address params.

  * `value`: - (Required) Ip address.
  * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

* `fqdn`: - (Optional) A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

  * `value`: - (Required) FQDN value.

### SMTP server

The smtp_server attribute supports the following:

* `email_address` SMTP email address.

* `server` SMTP network details.
  * `ip_address` An unique address that identifies a device on the internet or a local network in IPv4, IPv6 or format.
    * `ipv4`: - (Optional) ip address params.

      * `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
      * `value`: - (Required) Ip address.

    * `ipv6`: - (Optional) Ip address params.

      * `value`: - (Required) Ip address.
      * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

    * `fqdn`: - (Optional) A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.
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

### SNMP config

The snmp_config attribute supports the following:

* `is_enabled` - (Boolean) SNMP status. Whether SNMP is enabled.

* `users` - (List) SNMP user information. Each user object supports:
  * `username` - (String, max 64 chars) SNMP username. Required for SNMP trap v3 version.
  * `auth_type` - (String) SNMP user authentication type. Allowed values:
  
      | Enum      | Description                  |
      |-----------|------------------------------|
      | SHA       | SHA SNMP authentication      |
      | MD5       | MD5 SNMP authentication      |

  * `auth_key` - (String) SNMP user authentication key (must not contain single quotes).
  * `priv_type` - (String) SNMP user encryption type. Allowed values:

      | Enum      | Description        |
      |-----------|--------------------|
      | DES       | DES SNMP key       |
      | AES       | AES SNMP key       |

  * `priv_key` - (String) SNMP user encryption key (must not contain single quotes).

* `transports` - (List) SNMP transport details. Each transport object supports:
  * `protocol` - (String) SNMP protocol type. Allowed values:

      | Enum      | Description        |
      |-----------|--------------------|
      | UDP       | UDP protocol       |
      | TCP       | TCP protocol       |
      | UDP6      | UDP6 protocol      |
      | TCP6      | TCP6 protocol      |

  * `port` - (Integer) SNMP port.

* `traps` - (List) SNMP trap details. Each trap object supports:
  * `address` - (Block) An unique address block that supports:
    * `ipv4` - (Optional) ip address params.
      * `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
      * `value`: - (Required) Ip address.

    * `ipv6` - (Optional) Ip address params.
      * `value`: - (Required) Ip address.
      * `prefix_length`: - (Optional, default 128) The prefix length of the network to which this host IPv6 address belongs.

  * `username` - (String, max 64 chars) SNMP username. Required for SNMP trap v3 version.
  * `protocol` - (String) SNMP protocol type. Allowed values:
  
      | Enum      | Description        |
      |-----------|--------------------|
      | UDP       | UDP protocol       |
      | TCP       | TCP protocol       |
      | UDP6      | UDP6 protocol      |
      | TCP6      | TCP6 protocol      |

  * `port` - (Integer) SNMP port.
  * `should_inform` - (Boolean) SNMP inform mode status.
  * `engine_id` - (String) SNMP engine ID (hexadecimal string, e.g. 0x12345678).
  * `version` - (String) SNMP version. Allowed values:

      | Enum      | Description        |
      |-----------|--------------------|
      | V2        | V2 SNMP version    |
      | V3        | V3 SNMP version    |

  * `receiver_name` - (String, max 64 chars) SNMP receiver name.
  * `community_string` - (String) Community string (plaintext) for SNMP version 2.0.

### Pulse Status

The pulse_status attribute supports the following:

* `is_enabled`: - Flag to indicate if pulse is enabled or not.
* `pii_scrubbing_level`: - PII scrubbing level.

    | Enum     | Description                                                                                    |
    |----------|------------------------------------------------------------------------------------------------|
    | ALL      | Scrub All PII Information from Pulse including data like entity names and IP addresses        |
    | DEFAULT  | Default PII Scrubbing level. Data like entity names and IP addresses will not be scrubbed from Pulse |

See detailed information in [Nutanix Create Cluster Profile V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.1#tag/ClusterProfiles/operation/createClusterProfile).
