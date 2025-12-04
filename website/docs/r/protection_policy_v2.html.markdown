---
layout: "nutanix"
page_title: "NUTANIX: nutanix_protection_policy_v2"
sidebar_current: "docs-nutanix-resource-protection-policy-v2"
description: |-
  Creates a protection policy to automate the recovery point creation and replication process.

---

# nutanix_protection_policy_v2

Creates a protection policy to automate the recovery point creation and replication process.



## Example—Synchronous Protection Policy

```hcl

resource "nutanix_protection_policy_v2" "synchronous-protection-policy"{
  name        = "synchronous_protection_policy"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }

  replication_locations {
    domain_manager_ext_id = "6a44b05e-cb9b-4e7e-8d75-b1b4715369c4" # Local Domain Manager UUID
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17" # Remote Domain Manager UUID
    label                 = "target"
    is_primary            = false
  }

  category_ids = ["b08ed184-6b0c-42c1-8179-7b9026fe2676"]
}
```

## Example—Linear Retention Protection Policy

```hcl
resource "nutanix_protection_policy_v2" "linear-retention-protection-policy" {
  name = "linear-retention-protection-policy"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds = 7200
      recovery_point_type                   = "CRASH_CONSISTENT"
      retention {
        linear_retention {
          local  = 1
          remote = 1
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds = 7200
      recovery_point_type                   = "CRASH_CONSISTENT"
      retention {
        linear_retention {
          local  = 1
          remote = 1
        }
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = "6a44b05e-cb9b-4e7e-8d75-b1b4715369c4" # Local Domain Manager UUID
    label                 = "source"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17" # Remote Domain Manager UUID
    label      = "target"
    is_primary = false
  }

  category_ids = ["b08ed184-6b0c-42c1-8179-7b9026fe2676"]
}
```

## Example—Auto Rollup Retention Protection Policy

```hcl

# Create Auto Rollup Retention Protection Policy
resource "nutanix_protection_policy_v2" "auto-rollup-retention-protection-policy" {
  name = "auto_rollup_retention_protection_policy"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 20
      start_time                                    = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 30
      start_time                                    = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = "6a44b05e-cb9b-4e7e-8d75-b1b4715369c4" # Local Domain Manager UUID
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17" # Remote Domain Manager UUID
    label      = "target"
    is_primary = false
  }

  category_ids = ["b08ed184-6b0c-42c1-8179-7b9026fe2676"]
}
```
## Argument Reference

The following arguments are supported:

* `name`: -(Required) Name of the protection policy.
* `description`: -(Optional) Description of the protection policy.
* `replication_locations`: -(Required) Hypervisor details.
* `replication_configurations`: -(Required) Cluster reference for an entity.
* `category_ids`: -(Optional) Host entity with its attributes.


### Replication Locations
The replication_locations attribute supports the following:

* `label`: -(Required) This is a unique user defined label of the replication location. It is used to identify the location in the replication configurations.
* `domain_manager_ext_id`: -(Required) External identifier of the domain manager.
* `replication_sub_location`: -(Optional) Specifies the replication sub-locations where recovery points can be created or replicated.
* `is_primary`: -(Optional) One of the locations must be specified as the primary location. All the other locations must be connected to the primary location.

#### Replication Sub Location
The replication_sub_location attribute supports the following:
> One of `cluster_ext_ids` :
* `cluster_ext_ids` :  -(Optional) External identifier of the clusters.

##### Cluster Ext Ids
The cluster_ext_ids attribute supports the following:

* `cluster_ext_id`: -(Optional) List of Prism Element cluster external identifiers whose associated VMs and volume groups are protected. Only the primary location can have multiple clusters configured, while the other locations can specify only one cluster. Clusters must be specified for replication within the same Prism Central and cannot be specified for an MST type location. All clusters are considered if the cluster external identifier list is empty.

### Replication Configurations
The replication_configurations attribute supports the following:

* `source_location_label`: -(Required) Label of the source location from the replication locations list, where the entity is running. The location of type MST can not be specified as the replication source.
* `remote_location_label`: -(Optional) Label of the source location from the replication locations list, where the entity will be replicated.
* `schedule`: -(Required) Schedule for protection. The schedule specifies the recovery point objective and the retention policy for the participating locations.

#### Schedule
The schedule attribute supports the following:

* `recovery_point_type`: -(Optional) Type of recovery point.
    * `CRASH_CONSISTENT`: Crash-consistent Recovery points capture all the VM and application level details.
    * `APP_CONSISTENT`: Application-consistent Recovery points can capture all the data stored in the memory and also the in-progress transaction details.
* `recovery_point_objective_time_seconds`: -(Required) The Recovery point objective of the schedule in seconds and specified in multiple of 60 seconds. Only following RPO values can be provided for rollup retention type:
    - Minute(s): 1, 2, 3, 4, 5, 6, 10, 12, 15
    - Hour(s): 1, 2, 3, 4, 6, 8, 12
    - Day(s): 1
    - Week(s): 1, 2
* `retention`: -(Optional) Specifies the retention policy for the recovery point schedule.
* `start_time`: -(Optional) Represents the protection start time for the new entities added to the policy after the policy is created in h:m format. The values must be between 00h:00m and 23h:59m and in UTC timezone. It specifies the time when the first snapshot is taken and replicated for any entity added to the policy. If this is not specified, the snapshot is taken immediately and replicated for any new entity added to the policy.
* `sync_replication_auto_suspend_timeout_seconds`: -(Optional) Auto suspend timeout if there is a connection failure between locations for synchronous replication. If this value is not set, then the policy will not be suspended.

#### Retention
> One of `linear_retention` or `auto_rollup_retention` must be specified.

* `linear_retention`: -(Optional) Linear retention policy.
* `auto_rollup_retention`: -(Optional) Auto rollup retention policy.

##### Linear Retention
The linear_retention attribute supports the following:

* `local`: -(Required) Specifies the number of recovery points to retain on the local location.
* `remote`: -(Optional) Specifies the number of recovery points to retain on the remote location.

##### Auto Rollup Retention
The auto_rollup_retention attribute supports the following:

* `local`: -(Required) Specifies the auto rollup retention details.
* `remote`: -(Optional) Specifies the auto rollup retention details.

###### Local, Remote
The local, remote attribute in the auto_rollup_retention supports the following:

* `snapshot_interval_type`: -(Required) Snapshot interval period.
    * `YEARLY`: Specifies the number of latest yearly recovery points to retain.
    * `WEEKLY`: Specifies the number of latest weekly recovery points to retain.
    * `DAILY`: Specifies the number of latest daily recovery points to retain.
    * `MONTHLY`: Specifies the number of latest monthly recovery points to retain.
    * `HOURLY`: Specifies the number of latest hourly recovery points to retain.
* `frequency`: -(Required) Multiplier to 'snapshotIntervalType'. For example, if 'snapshotIntervalType' is 'YEARLY' and 'multiple' is 5, then 5 years worth of rollup snapshots will be retained.

## Import

This helps to manage existing entities which are not created through terraform. protection policy can be imported using the `UUID`. (ext_id in v4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_protection_policy_v2" "import_pp" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_protection_policies_v2" "fetch_policies"{}
terraform import nutanix_protection_policy_v2.import_pp <UUID>
```

See detailed information in [Nutanix Protection Policy v4](https://developers.nutanix.com/api-reference?namespace=datapolicies&version=v4.0#tag/ProtectionPolicies/operation/createProtectionPolicy).
