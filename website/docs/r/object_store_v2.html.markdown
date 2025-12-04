---
layout: "nutanix"
page_title: "NUTANIX: nutanix_object_store_v2 "
sidebar_current: "docs-nutanix-resource-object-store-v2"
description: |-
  Run the prechecks, create and start the deployment of an Object store on Prism Central.
---

# nutanix_object_store_v2

Run the prechecks, create and start the deployment of an Object store on Prism Central.

> ⚠️ **Warning:** Before deleting the Object Store, make sure to delete all buckets inside it manually.
> Currently, the Terraform provider does not support the Delete Bucket API.

> ⚠️ **Warning:** The Object Store **update** operation is intended **only** to resume a failed deployment.
> It should be used when the Object Store is in the `OBJECT_STORE_DEPLOYMENT_FAILED` state.
> Triggering an update in this state will attempt to resume the deployment process.

## Example Usage

```hcl
resource "nutanix_object_store_v2" "example"{
  name                     = "tf-example-os"
  description              = "terraform create object store example"
  deployment_version       = "5.1.1"
  domain                   = "msp.pc-idbc.nutanix.com"
  num_worker_nodes         = 1
  cluster_ext_id           = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
  total_capacity_gib       = 20 * pow(1024, 3)
  public_network_reference = "57c4caf1-67e3-457e-8265-6d872f2a3135"
  public_network_ips {
    ipv4 {
      value = "10.44.77.123"
    }
  }
  storage_network_reference = "57c4caf1-67e3-457e-8265-6d872f2a3135"
  storage_network_dns_ip {
    ipv4 {
      value = "10.44.77.124"
    }
  }
  storage_network_vip {
    ipv4 {
      value = "10.44.77.125"
    }
  }
}

# Deploying Object Store in draft state
resource "nutanix_object_store_v2" "example-draft" {
  name                     = "tf-draft-os"
  description              = "terraform deploy object store draft example"
  deployment_version       = "5.1.1"
  domain                   = "msp.pc-idbc.nutanix.com"
  num_worker_nodes         = 1
  cluster_ext_id           = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
  total_capacity_gib       = 20 * pow(1024, 3)
  public_network_reference = "57c4caf1-67e3-457e-8265-6d872f2a3135"
  state                    = "UNDEPLOYED_OBJECT_STORE"
  public_network_ips {
    ipv4 {
      value = "10.44.77.126"
    }
  }
  storage_network_reference = "57c4caf1-67e3-457e-8265-6d872f2a3135"
  storage_network_dns_ip {
    ipv4 {
      value = "10.44.77.127"
    }
  }
  storage_network_vip {
    ipv4 {
      value = "10.44.77.128"
    }
  }
}

```

## Argument Reference

The following arguments are supported:

- `metadata`: -(Optional) Metadata associated with this resource.
- `name`: -(Required) The name of the Object store.
- `description`: -(Optional) A brief description of the Object store.
- `deployment_version`: -(Optional) The deployment version of the Object store.
- `domain`: -(Optional) The DNS domain/subdomain the Object store belongs to. All the Object stores under one Prism Central must have the same domain name. The domain name must consist of at least 2 parts separated by a '.'. Each part can contain upper and lower case letters, digits, hyphens, or underscores. Each part can be up to 63 characters long. The domain must begin and end with an alphanumeric character. For example - 'objects-0.pc_nutanix.com'.
- `region`: -(Optional) The region in which the Object store is deployed.
- `num_worker_nodes`: -(Optional) The number of worker nodes (VMs) to be created for the Object store. Each worker node requires 10 vCPUs and 32 GiB of memory.
- `cluster_ext_id`: -(Optional) UUID of the AHV or ESXi cluster.
- `storage_network_reference`: -(Optional) Reference to the Storage Network of the Object store. This is the subnet UUID for an AHV cluster or the IPAM name for an ESXi cluster.
- `storage_network_vip`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `storage_network_dns_ip`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `public_network_reference`: -(Optional) Public network reference of the Object store. This is the subnet UUID for an AHV cluster or the IPAM name for an ESXi cluster.
- `public_network_ips`: -(Optional) A list of static IP addresses used as public IPs to access the Object store.
- `total_capacity_gib`: -(Optional) Size of the Object store in GiB.
- `state`: -(Optional) Enum for the state of the Object store.
  | Enum | Description |
  |----------------------------------------|-----------------------------------------------------------------|
  | `DEPLOYING_OBJECT_STORE` | The Object store is being deployed. |
  | `OBJECT_STORE_DEPLOYMENT_FAILED` | The Object store deployment has failed. |
  | `DELETING_OBJECT_STORE` | A deployed Object store is being deleted. |
  | `OBJECT_STORE_OPERATION_FAILED` | There was an error while performing an operation on the Object store. |
  | `UNDEPLOYED_OBJECT_STORE` | The Object store is not deployed. |
  | `OBJECT_STORE_OPERATION_PENDING` | There is an ongoing operation on the Object store. |
  | `OBJECT_STORE_AVAILABLE` | There are no ongoing operations on the deployed Object store. |
  | `OBJECT_STORE_CERT_CREATION_FAILED` | Creating the Object store certificate has failed. |
  | `CREATING_OBJECT_STORE_CERT` | A certificate is being created for the Object store. |
  | `OBJECT_STORE_DELETION_FAILED` | There was an error deleting the Object store. |

### Metadata

The `metadata` argument supports the following:

- `owner_reference_id`: -(Optional) A globally unique identifier that represents the owner of this resource.
- `owner_user_name`: -(Optional) The userName of the owner of this resource.
- `project_reference_id`: -(Optional) A globally unique identifier that represents the project this resource belongs to.
- `project_name`: -(Optional) The name of the project this resource belongs to.
- `category_ids`: -(Optional) A list of globally unique identifiers that represent all the categories the resource is associated with.

## Attributes Reference

The following attributes are exported:

- `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `metadata`: - Metadata associated with this resource.
- `name`: - The name of the Object store.
- `creation_time`: - The time when the Object store was created.
- `last_update_time`: - The time when the Object store was last updated.
- `description`: - A brief description of the Object store.
- `deployment_version`: - The deployment version of the Object store.
- `domain`: - The DNS domain/subdomain the Object store belongs to. All the Object stores under one Prism Central must have the same domain name. The domain name must consist of at least 2 parts separated by a '.'. Each part can contain upper and lower case letters, digits, hyphens, or underscores. Each part can be up to 63 characters long. The domain must begin and end with an alphanumeric character. For example - 'objects-0.pc_nutanix.com'.
- `region`: - The region in which the Object store is deployed.
- `num_worker_nodes`: - The number of worker nodes (VMs) to be created for the Object store. Each worker node requires 10 vCPUs and 32 GiB of memory.
- `cluster_ext_id`: - UUID of the AHV or ESXi cluster.
- `storage_network_reference`: - Reference to the Storage Network of the Object store. This is the subnet UUID for an AHV cluster or the IPAM name for an ESXi cluster.
- `storage_network_vip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `storage_network_dns_ip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `public_network_reference`: - Public network reference of the Object store. This is the subnet UUID for an AHV cluster or the IPAM name for an ESXi cluster.
- `public_network_ips`: - A list of static IP addresses used as public IPs to access the Object store.
- `total_capacity_gib`: - Size of the Object store in GiB.
- `state`: - Enum for the state of the Object store.
  | Enum | Description |
  |----------------------------------------|-----------------------------------------------------------------|
  | `"DEPLOYING_OBJECT_STORE"` | The Object store is being deployed. |
  | `"OBJECT_STORE_DEPLOYMENT_FAILED"` | The Object store deployment has failed. |
  | `"DELETING_OBJECT_STORE"` | A deployed Object store is being deleted. |
  | `"OBJECT_STORE_OPERATION_FAILED"` | There was an error while performing an operation on the Object store. |
  | `"UNDEPLOYED_OBJECT_STORE"` | The Object store is not deployed. |
  | `"OBJECT_STORE_OPERATION_PENDING"` | There is an ongoing operation on the Object store. |
  | `"OBJECT_STORE_AVAILABLE"` | There are no ongoing operations on the deployed Object store. |
  | `"OBJECT_STORE_CERT_CREATION_FAILED"` | Creating the Object store certificate has failed. |
  | `"CREATING_OBJECT_STORE_CERT"` | A certificate is being created for the Object store. |
  | `"OBJECT_STORE_DELETION_FAILED"` | There was an error deleting the Object store. |

- `certificate_ext_ids`: - A list of the UUIDs of the certificates of an Object store.

### Links

The `links` argument exports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Metadata

The `metadata` argument exports the following:

- `owner_reference_id`: - A globally unique identifier that represents the owner of this resource.
- `owner_user_name`: - The userName of the owner of this resource.
- `project_reference_id`: - A globally unique identifier that represents the project this resource belongs to.
- `project_name`: - The name of the project this resource belongs to.
- `category_ids`: - A list of globally unique identifiers that represent all the categories the resource is associated with.

### Storage Network VIP

The `storage_network_vip` argument exports the following:

- `ipv4`: An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv6`: An unique address that identifies a device on the internet or a local network in IPv6 format.

### Storage Network DNS IP

The `storage_network_dns_ip` argument exports the following:

- `ipv4`: An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv6`: An unique address that identifies a device on the internet or a local network in IPv6 format.

### Public Network IPs

The `public_network_ips` argument exports the following:

- `ipv4`: An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv6`: An unique address that identifies a device on the internet or a local network in IPv6 format.

### IPv4, IPv6

The `ipv4` and `ipv6` argument exports the following:

- `value`: - The IPv4/IPv6 address of the host.
- `prefix_length`: - The prefix length of the network to which this host IPv4 address belongs. Default for IPv4 is 32 and for IPv6 is 128.

## Import

This helps to manage existing entities which are not created through terraform. Object store can be imported using the `UUID`. (ext_id in v4 terms). eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_object_store_v2" "imported" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_object_stores_v2" "fetch_objects"{}
terraform import nutanix_object_store_v2.imported <UUID>
```

See detailed information in [Nutanix Get Object Store V4 ](https://developers.nutanix.com/api-reference?namespace=objects&version=v4.0#tag/ObjectStores/operation/getObjectstoreById).
