---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pc_restore_source_v2"
sidebar_current: "docs-nutanix-resource-pc-restore-source-v2"
description: |-
  Creates a restore source pointing to a cluster or object store to restore the domain manager. The created restore source is intended to be deleted after use. If the restore source is not deleted using the deleteRestoreSource API, then it is auto-deleted after sometime. Also note that a restore source will not contain a backup policy. It is only used to access the backup data at the location from where the Prism Central may be restored. Credentials used to access the restore source are not validated at the time of creation of the restore source. They are validated when the restore source is used to fetch data.
---

# nutanix_pc_restore_source_v2

> The restore source is auto-deleted after sometime, nutanix_pc_restore_source_v2 resource is auto-create new restore source if it was deleted, so notice that the id of the restore source resource will be different after recreation. and for update restore source, it will create a new restore source with the new configuration

Create a restore source pointing to a cluster or object store to restore the domain manager. The created restore source is intended to be deleted after use. If the restore source is not deleted using the deleteRestoreSource API, then it is auto-deleted after sometime. Also note that a restore source will not contain a backup policy. It is only used to access the backup data at the location from where the Prism Central may be restored. Credentials used to access the restore source are not validated at the time of creation of the restore source. They are validated when the restore source is used to fetch data.


## Example Usage - Cluster Location

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

# restore source is auto-deleted after sometime, nutanix_pc_restore_source_v2 resource is auto-create
# new restore source if it was deleted, so notice that the id of the restore source will be different
# after recreation
resource "nutanix_pc_restore_source_v2" "cluster-location" {
  provider = nutanix.pe
  location {
    cluster_location {
      config {
        ext_id = "323860ca-bd10-411e-9fe0-1430b62eaf45"
      }
    }
  }
}

```

## Example Usage - Object Store Location

```hcl
// using object store location
resource "nutanix_pc_restore_source_v2" "object-store-location"{
  location {
    object_store_location {
      provider_config {
        bucket_name = "nutanix-terraform-bucket"
        region      = "us-west-1"
        credentials {
          access_key_id     = "IHSAJHDHADFWYTKJHFGCJKHASGJHKDSA"
          secret_access_key = "JGSDHJYHGFHGHDS+JKBASDF/HSDAFHJ+SjkfbdsASDfdJFdSDFJfk"
        }
      }
      backup_policy {
        rpo_in_minutes = 120
      }
    }
  }
  lifecycle {
    ignore_changes = [
      location[0].object_store_location[0].provider_config[0].credentials
    ]
  }
}

```

## Argument Reference

The following arguments are supported:

- `location`: -(Required) Location of the backup target. For example, a cluster or an object store endpoint, such as AWS s3.

### Location

The location argument supports the following:

> one of the following is required:

- `cluster_location`: -(Optional) A boolean value indicating whether to enable lockdown mode for a cluster.
- `object_store_location`: -(Optional) Currently representing the build information to be used for the cluster creation.

#### Cluster Location

The `cluster_location` argument supports the following:

- `config`: -(Required) Cluster reference of the remote cluster to be connected.

##### Config

The `config` argument supports the following:

- `ext_id`: -(Required) Cluster UUID of a remote cluster.

#### Object Store Location

The `object_store_location` argument supports the following:

- `provider_config`: -(Required) The base model of S3 object store endpoint where domain manager is backed up.
- `backup_policy`: -(Optional) Backup policy for the object store provided.

##### Provider Config

The `provider_config` argument supports the following:

- `bucket_name`: -(Required) The bucket name of the object store endpoint where backup data of domain manager is to be stored.
- `region`: -(Optional) The region name of the object store endpoint where backup data of domain manager is stored. Default is `us-east-1`.
- `credentials`: -(Optional) Secret credentials model for the object store containing access key ID and secret access key.

###### Credentials

The `credentials` argument supports the following:

- `access_key_id`: -(Required) Access key ID for the object store provided for backup target.
- `secret_access_key`: -(Required) Secret access key for the object store provided for backup target.

##### Backup Policy

The `backup_policy` argument supports the following:

- `rpo_in_minutes`: -(Required) RPO interval in minutes at which the backup will be taken. The Value should be in the range of 60 to 1440.

## Import

This helps to manage existing entities which are not created through terraform. Restore source can be imported using the `UUID`. (ext_id in v4 terms).  eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_pc_restore_source_v2" "import_restore_source" {}

// execute this command
terraform import nutanix_pc_restore_source_v2.import_restore_source <UUID>
```

See detailed information in [Nutanix Restore Source V4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/createRestoreSource).
