---
layout: "nutanix"
page_title: "NUTANIX: nutanix_clusters"
sidebar_current: "docs-nutanix-datasource-clusters-x"
description: |-
 Describes a Clusters
---

# nutanix_clusters

Describes Clusters

## Example Usage

```hcl
data "nutanix_clusters" "clusters" {}
```

## Argument Reference

No arguments are supported:

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when image was last updated.
* `uuid`: - image UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when image was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - image name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

## Attribute Reference

The following attributes are exported:

* `entities`: List of Clusters

# Entities

The entities attribute element contains the followings attributes:

* `name`: -  The name for the image.
* `categories`: - Categories for the image.
* `project_reference`: - The reference to a project.
* `owner_reference`: - The reference to a user.
* `availability_zone_reference`: - The reference to a availability_zone.
* `api_version` - The API version.
* `description`: - A description for image.
* `metadata`: - The image kind metadata.
* `state`: - The state of the cluster entity.
* `gpu_driver_version`: - GPU driver version.
* `client_auth`: - Client authentication config.
* `authorized_piblic_key_list`: - List of valid ssh keys for the cluster.
* `software_map_ncc`: - Map of software on the cluster with software type as the key.
* `software_map_nos`: - Map of software on the cluster with software type as the key.
* `encryption_status`: - Cluster encryption status.
* `ssl_key_type`: - SSL key type. Key types with RSA_2048, ECDSA_256 and ECDSA_384 are supported for key generation and importing.
* `ssl_key_signing_info`: - Customer information used in Certificate Signing Request for creating digital certificates.
* `ssl_key_expire_datetime`: - UTC date and time in RFC-3339 format when the key expires
* `service_list`: - Array of enabled cluster services. For example, a cluster can function as both AOS and cloud data gateway. - 'AOS': Regular Prism Element - 'PRISM_CENTRAL': Prism Central - 'CLOUD_DATA_GATEWAY': Cloud backup and DR gateway - 'AFS': Cluster for file server - 'WITNESS' : Witness cluster - 'XI_PORTAL': Xi cluster.
* `supported_information_verbosity`: - Verbosity level settings for populating support information. - 'Nothing': Send nothing - 'Basic': Send basic information - skip core dump and hypervisor stats information - 'BasicPlusCoreDump': Send basic and core dump information - 'All': Send all information (Default value: BASIC_PLUS_CORE_DUMP)
* `certification_signing_info`: - Customer information used in Certificate Signing Request for creating digital certificates.
* `operation_mode`: - Cluster operation mode. - 'NORMAL': Cluster is operating normally. - 'READ_ONLY': Cluster is operating in read only mode. - 'STAND_ALONE': Only one node is operational in the cluster. This is valid only for single node or two node clusters. - 'SWITCH_TO_TWO_NODE': Cluster is moving from single node to two node cluster. - 'OVERRIDE': Valid only for single node cluster. If the user wants to run vms on a single node cluster in read only mode, he can set the cluster peration mode to override. Writes will be allowed in override mode.
* `ca_certificate_list`: - Zone name used in value of TZ environment variable.
* `enabled_feature_list`: - Array of enabled features.
* `is_available`: - Indicates if cluster is available to contact. (Readonly)
* `build`: - Cluster build details.
* `timezone`: - Zone name used in value of TZ environment variable.
* `cluster_arch`: - Cluster architecture. (Readonly, Options: Options : X86_64 , PPC64LE)
* `management_server_list`: - List of cluster management servers. (Readonly)
* `masquerading_port`: - Port used together with masquerading_ip to connect to the cluster.
* `masquerading_ip`: - The cluster NAT'd or proxy IP which maps to the cluster local IP.
* `external_ip`: - The local IP of cluster visible externally.
* `http_proxy_list`: - List of proxies to connect to the service centers.
* `smtp_server_type`: - SMTP Server type.
* `smtp_server_email_address`: - SMTP Server Email Address.
* `smtp_server_credentials`: - SMTP Server Credentials.
* `smtp_server_proxy_type_list`: - SMTP Server Proxy Type List
* `smtp_server_address`: - SMTP Server Address.
* `ntp_server_ip_list`: - The list of IP addresses or FQDNs of the NTP servers.
* `external_subnet`: - External subnet for cross server communication. The format is IP/netmask. (default 172.16.0.0/255.240.0.0)
* `external_data_services_ip`: - The cluster IP address that provides external entities access to various cluster data services.
* `internal_subnet`: - The internal subnet is local to every server - its not visible outside.iSCSI requests generated internally within the appliance (by user VMs or VMFS) are sent to the internal subnet. The format is IP/netmask.
* `domain_server_nameserver`: -  The IP of the nameserver that can resolve the domain name. Must set when joining the domain.
* `domain_server_name`: - Joined domain name. In 'put' request, empty name will unjoin the cluster from current domain.
* `domain_server_credentials`: - Cluster domain credentials.
* `nfs_subnet_whitelist`: - Comma separated list of subnets (of the form 'a.b.c.d/l.m.n.o') that are allowed to send NFS requests to this container. If not specified, the global NFS whitelist will be looked up for access permission. The internal subnet is always automatically considered part of the whitelist, even if the field below does not explicitly specify it. Similarly, all the hypervisor IPs are considered part of the whitelist. Finally, to permit debugging, all of the SVMs local IPs are considered to be implicitly part of the whitelist.
* `name_server_ip_list`: - The list of IP addresses of the name servers.
* `http_proxy_whitelist`: - HTTP proxy whitelist.
* `analysis_vm_efficiency_map`: - Map of cluster efficiency which includes numbers of inefficient vms. The value is populated by analytics on PC. (Readonly)

## Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, attributes supports the following:

* `kind`: - The kind name (Default value: project).
* `name`: - the name.
* `UUID`: - the UUID.

### Version

The version attribute supports the following:

* `product_name`: - Name of the producer/distribution of the image. For example windows or red hat.
* `product_version`: - Version string for the disk image.

See detailed information in [Nutanix Clusters](https://www.nutanix.dev/api_references/prism-central-v3/#/d93c30e04327e-get-a-list-of-existing-clusters).
