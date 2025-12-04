---
layout: "nutanix"
page_title: "NUTANIX: nutanix_storage_container_v2"
sidebar_current: "docs-nutanix-datasource-storage-container-v2"
description: |-
   This operation retrieves a Storage Container configuration.
---

# nutanix_storage_container_v2

Provides a datasource to Fetch the configuration details of the existing Storage Container identified by the {containerExtId}.

## Example Usage

```hcl
data "nutanix_storage_container_v2" "get-storage-container"{
  ext_id = "1891fd3a-1ef7-4947-af56-9ee4b973c6fd"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) storage container UUID

## Attribute Reference

The following attributes are exported:

* `ext_id`: - the storage container uuid
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

* `container_ext_id`: - the storage container ext id
* `owner_ext_id`: - owner ext id
* `name`: Name of the storage container.  Note that the name of Storage Container should be unique per cluster.
* `cluster_ext_id`: - ext id for the cluster owning the storage container.
* `storage_pool_ext_id`: - extId of the Storage Pool owning the Storage Container instance.
* `is_marked_for_removal`: - Indicates if the Storage Container is marked for removal. This field is set when the Storage Container is about to be destroyed.
* `max_capacity_bytes`: - Maximum physical capacity of the Storage Container in bytes.
* `logical_explicit_reserved_capacity_bytes`: - Total reserved size (in bytes) of the container (set by Admin). This also accounts for the container's replication factor. The actual reserved capacity of the container will be the maximum of explicitReservedCapacity and implicitReservedCapacity.
* `logical_implicit_reserved_capacity_bytes`: - This is the summation of reservations provisioned on all vdisks in the container. The actual reserved capacity of the container will be the maximum of explicitReservedCapacity and implicitReservedCapacity
* `logical_advertised_capacity_bytes`: - Max capacity of the Container as defined by the user.
* `replication_factor`: - Replication factor of the Storage Container.
* `nfs_whitelist_addresses`: - List of NFS addresses which need to be whitelisted.
* `is_nfs_whitelist_inherited`: - Indicates whether the NFS whitelist is inherited from global config.
* `erasure_code`: - Indicates the current status value for Erasure Coding for the Container. available values:  `NONE`,    `OFF`,    `ON`

* `is_inline_ec_enabled`: - Indicates whether data written to this container should be inline erasure coded or not. This field is only considered when ErasureCoding is enabled.
* `has_higher_ec_fault_domain_preference`: - Indicates whether to prefer a higher Erasure Code fault domain.
* `erasure_code_delay_secs`: - Delay in performing ErasureCode for the current Container instance.
* `cache_deduplication`: - Indicates the current status of Cache Deduplication for the Container. available values:  `NONE`,    `OFF`,    `ON`
* `on_disk_dedup`: - Indicates the current status of Disk Deduplication for the Container. available values:  `NONE`,    `OFF`,    `POST_PROCESS`
* `is_compression_enabled`: - Indicates whether the compression is enabled for the Container.
* `compression_delay_secs`: - The compression delay in seconds.
* `is_internal`: - Indicates whether the Container is internal and is managed by Nutanix.
* `is_software_encryption_enabled`: - Indicates whether the Container instance has software encryption enabled.
* `is_encrypted`: - Indicates whether the Container is encrypted or not.
* `affinity_host_ext_id`: - Affinity host extId for RF 1 Storage Container.
* `cluster_name`: - Corresponding name of the Cluster owning the Storage Container instance.


### nfs_whitelist_addresses

* `ipv4`: Reference to address configuration
* `ipv6`: Reference to address configuration
* `fqdn`: Reference to address configuration

### ipv4, ipv6 (Reference to address configuration)

* `value`: value of address
* `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.

### fqdn (Reference to address configuration)

* `value`: value of fqdn address


See detailed information in [Nutanix Get Storage Container v4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0#tag/StorageContainers/operation/getStorageContainerById).
