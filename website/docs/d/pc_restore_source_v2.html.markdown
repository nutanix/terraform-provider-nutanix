---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pc_restore_source_v2 "
sidebar_current: "docs-nutanix-datasource-pc-restore-source-v2"
description: |-
  Retrieves the restore source from the PE cache store and returns the restore source configuration and external identifier to the user.
---

# nutanix_pc_restore_source_v2

Retrieves the restore source from the PE cache store and returns the restore source configuration and external identifier to the user.

## Example Usage

```hcl
#defining nutanix configuration for PE
provider "nutanix" {
  alias    = "pe"
  username = var.nutanix_pe_username
  password = var.nutanix_pe_password
  endpoint = var.nutanix_pe_endpoint # PE endpoint
  insecure = true
  port     = 9440
}


// get the restore source
data "nutanix_pc_restore_source_v2" "restore-source" {
  provider = nutanix.pe
  ext_id   = "323860ca-bd10-411e-9fe0-1430b62eaf45"
}

```

## Argument Reference

The following arguments are supported:

- `ext_id`: -(Required) A globally unique identifier of an instance that is suitable for external consumption.

## Attributes Reference

The following attributes are exported:

- `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `location`: - Location of the backup target. For example, a cluster or an object store endpoint, such as AWS s3.

### Location

The location argument exports the following:

- `cluster_location`: - A boolean value indicating whether to enable lockdown mode for a cluster.
- `object_store_location`: - Currently representing the build information to be used for the cluster creation.

#### Cluster Location

The `cluster_location` argument exports the following:

- `config`: - Cluster reference of the remote cluster to be connected.

##### Config

The `config` argument exports the following:

- `ext_id`: - Cluster UUID of a remote cluster.
- `name`: - Name of the cluster.

#### Object Store Location

The `object_store_location` argument exports the following:

- `provider_config`: - The base model of S3 object store endpoint where domain manager is backed up.
- `backup_policy`: - Backup policy for the object store provided.

##### Provider Config

The `provider_config` argument exports the following:

- `bucket_name`: - The bucket name of the object store endpoint where backup data of domain manager is to be stored.
- `region`: - The region name of the object store endpoint where backup data of domain manager is stored. Default is `us-east-1`.
- `credentials`: - Secret credentials model for the object store containing access key ID and secret access key.

###### Credentials

The `credentials` argument exports the following:

- `access_key_id`: - Access key ID for the object store provided for backup target.
- `secret_access_key`: - Secret access key for the object store provided for backup target.

##### Backup Policy

The `backup_policy` argument exports the following:

- `rpo_in_minutes`: - RPO interval in minutes at which the backup will be taken. The Value should be in the range of 60 to 1440.

See detailed information in [Nutanix Restore Source V4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/getRestoreSourceById).
