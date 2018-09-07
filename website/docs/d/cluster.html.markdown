---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster"
sidebar_current: "docs-nutanix-datasource-cluster"
description: |-
 Describes a Cluster
---

# nutanix_cluster
Describes Clusters

## Example Usage

```hcl
data "nutanix_cluster" "cluster" {
   cluster_id = "${data.nutanix_clusters.clusters.entities.1.metadata.uuid}"
}`
```

## Argument Reference

The following arguments are supported:

* `cluster_id`: Represents clusters uuid

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when image was last updated.
* `uuid`: - image uuid.
* `creation_time`: - UTC date and time in RFC-3339 format when image was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - image name.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

## Attribute Reference

The following attributes are exported:

* `name`: -  The name for the image.
* `categories`: - Categories for the image.
* `project_reference`: - The reference to a project.
* `owner_reference`: - The reference to a user.
* `availability_zone_reference`: - The reference to a availability_zone.
* `api_version` -
* `description`: - A description for image.
* `metadata`: - The image kind metadata.
* `state`: - The state of the cluster entity.
* `gpu_driver_version`: -
* `client_auth`: -
* `authorized_piblic_key_list`: -
* `software_map_ncc`: -
* `software_map_nos`: -
* `encryption_status`: -
* `ssl_key_type`: -
* `ssl_key_signing_info`: -
* `ssl_key_expire_datetime`: -
* `service_list`: -
* `supported_information_verbosity`: -
* `certification_signing_info`: -
* `operation_mode`: -
* `ca_certificate_list`: -
* `enabled_feature_list`: -
* `is_available`: -
* `build`: -
* `timezone`: -
* `cluster_arch`: -
* `management_server_list`: -
* `masquerading_port`: -
* `masquerading_ip`: -
* `external_ip`: -
* `http_proxy_list`: -
* `smtp_server_type`: -
* `smtp_server_email_address`: -
* `smtp_server_credentials`: -
* `smtp_server_proxy_type_list`: -
* `smtp_server_address`: -
* `ntp_server_ip_list`: -
* `external_subnet`: -
* `external_data_services_ip`: -
* `internal_subnet`: -
* `domain_server_nameserver`: -
* `domain_server_name`: -
* `domain_server_credentials`: -
* `nfs_subnet_whitelist`: -
* `name_server_ip_list`: -
* `http_proxy_whitelist`: -
* `analysis_vm_efficiency_map`: -


### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, attributes supports the following:

* `kind`: - The kind name (Default value: project).
* `name`: - the name.
* `uuid`: - the uuid.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.


### Version

The version attribute supports the following:

* `product_name`: - Name of the producer/distribution of the image. For example windows or red hat.
* `product_version`: - Version string for the disk image.

See detailed information in [Nutanix Image](https://nutanix.github.io/Automation/experimental/swagger-redoc-sandbox/#tag/clusters/paths/~1clusters~1multicluster_config/post).