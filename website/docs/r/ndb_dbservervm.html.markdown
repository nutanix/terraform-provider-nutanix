---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_dbserver_vm"
sidebar_current: "docs-nutanix-resource-ndb-dbserver-vm"
description: |-
  This operation submits a request to create, update and delete database server VMs in Nutanix database service (NDB).
  Note: For 1.8.0 release, only postgress database type is qualified and officially supported.
---

# nutanix_ndb_dbserver_vm

Provides a resource to create database server VMs based on the input parameters. For 1.8.0 release, only postgress database type is qualified and officially supported.

## Example Usage

### dbserver vm resource with software profile

```hcl
    resource nutanix_ndb_dbserver_vm acctest-managed {
		database_type = "postgres_database"
        nx_cluster_id = {{ nx_cluster_id }}
		software_profile_id = {{ software_profile_id }}
		software_profile_version_id =  {{ software_profile_version_id }}

        // dbserver details 
        description = "{{ description }}"
		compute_profile_id = {{ compute_profile_id }}
		network_profile_id = {{ network_profile_id }}
		vm_password = "{{ vm_password }}"
		postgres_database {
			vm_name = "test-vm"
			client_public_key = "{{ public_key }}"
		}
	}
```

### dbserver vm resource with time machine
```hcl
    resource nutanix_ndb_dbserver_vm acctest-managed {
		database_type = "postgres_database"
        nx_cluster_id = {{ nx_cluster_id }}
		time_machine_id = {{ time_machine_id }}

        // dbserver details 
        description = "{{ description }}"
		compute_profile_id = {{ compute_profile_id }}
		network_profile_id = {{ network_profile_id }}
		vm_password = "{{ vm_password }}"
		postgres_database {
			vm_name = "test-vm"
			client_public_key = "{{ public_key }}"
		}
	}
```

## Argument Reference

The following arguments are supported:

* `database_type`: (Required) database type. Valid values: postgres_database 
* `software_profile_id`: (Optional) software profile id you want to provision a database server VM from an existing software profile.Required with software_profile_version_id. Conflicts with time_machine_id . 
* `software_profile_version_id`: (Optional) SOftware Profile Version Id. 
* `time_machine_id`: (Optional) Time Machine id you want to provision a database server VM by using the database and operating system software stored in a time machine. Conflicts with software_profile_id. 
* `snapshot_id`: (Optional) Snapshot id. If not given, it will use latest snapshot to provision db server vm. 

* `description`: (Optional) Type a description for the database server VM.
* `compute_profile_id`: (Optional) Compute profile id.
* `network_profile_id`: (Optioanl) Network profile id.
* `vm_password`: (Optional) password of the NDB drive user account.
* `postgres_database`: (Optional) Postgres database server vm
* `maintenance_tasks`: (Optional) maintenance window configured to enable automated patching.


* `delete`:- (Optional) Delete the VM and associated storage. Default value is true
* `remove`:- (Optional) Unregister the database from NDB. Default value is false
* `soft_remove`:- (Optional) Soft remove. Default will be false
* `delete_vgs`:- (Optional) Delete volume grous. Default value is true
* `delete_vm_snapshots`:- (Optional) Delete the vm snapshots. Default is true


### postgres_database
* `vm_name`: (Required) name for the database server VM.
* `client_public_key`: (Required) use SSH public keys to access the database server VM.

### maintenance_tasks
* `maintenance_window_id`: Associate an existing maintenance window id. NDB starts OS patching or database patching as per the schedule defined in the maintenance window.
* `tasks`: Tasks for the maintenance.
* `tasks.task_type`: use this option if you want NDB to perform database patching or OS patching automatically. Supports [ OS_PATCHING, DB_PATCHING ]. 
* `tasks.pre_command`: add pre (operating system and database patching) commands.
* `tasks.post_command`:add post (operating system and database patching) commands.


### actionarguments

Structure for each action argument in actionarguments list:

* `name`: name of the dbserver vm
* `properties`: Properties of dbserver vm
* `dbserver_cluster_id`: dbserver cluster id.
* `vm_cluster_name`: cluster name for dbserver vm
* `vm_cluster_uuid`: clusetr uuid for dbserver vm
* `ip_addresses`: IP addresses of the dbserver vm
* `mac_addresses`: Mac addresses of dbserver vm
* `type`: Type of entity. i.e. Dbserver
* `status`: Status of Dbserver . Active or not.
* `client_id`:  client id
* `era_drive_id`: era drive id
* `era_version`: era version
* `vm_timezone`:  timezone of dbserver vm


See detailed information in [NDB Provision Database Server VM](https://www.nutanix.dev/api_references/ndb/#/c9126257bc0fc-provision-database-server).