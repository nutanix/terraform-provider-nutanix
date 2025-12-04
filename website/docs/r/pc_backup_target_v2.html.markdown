---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pc_backup_target_v2 "
sidebar_current: "docs-nutanix-resource-pc-backup-target-v2"
description: |-
  Creates a cluster or object store as the backup target. For a given Prism Central, there can be up to 3 clusters as backup targets and 1 object store as backup target. If any cluster or object store is not eligible for backup or lacks appropriate permissions, the API request will fail. For object store backup targets, specifying backup policy is mandatory along with the location of the object store.
---

# nutanix_pc_backup_target_v2

Create a cluster or object store as the backup target. For a given Prism Central, there can be up to 3 clusters as backup targets and 1 object store as backup target. If any cluster or object store is not eligible for backup or lacks appropriate permissions, the API request will fail. For object store backup targets, specifying backup policy is mandatory along with the location of the object store.

## Example Usage - Cluster Location

```hcl

resource "nutanix_pc_backup_target_v2" "cluster-location"{
  domain_manager_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17"
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

//using object store location
resource "nutanix_pc_backup_target_v2" "object-store-location"{
  domain_manager_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17"
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

- `domain_manager_ext_id`: -(Required) A unique identifier for the domain manager.
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

This helps to manage existing entities which are not created through terraform. Backup Target can be imported using the `domainManagerUUID/backupTargetUUID`. (ext_id in v4 API context). eg,

**Note**:To import Backup Target, you need to have the Backup Target UUID, and provide it in the format mentioned above while importing.

```hcl
// create its configuration in the root module. For example:
resource "nutanix_pc_backup_target_v2" "imported"{}


// execute the below command. UUID can be fetched using datasource. Example:

// list pcs
data "nutanix_clusters_v2" "pcs" {
	filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  domainManagerExtId = data.nutanix_clusters_v2.pcs.cluster_entities.0.ext_id
}

// list backup targets for a pc
data "nutanix_pc_backup_targets_v2" "backup_targets" {
  domain_manager_ext_id = local.domainManagerExtId
}


terraform import nutanix_pc_backup_target_v2.imported <domainManagerUUID/backupTargetUUID>
```

See detailed information in [Nutanix Create Backup Target V4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/createBackupTarget).
