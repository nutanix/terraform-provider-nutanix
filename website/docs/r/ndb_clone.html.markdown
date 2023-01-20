layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_clone"
sidebar_current: "docs-nutanix-resource-ndb-clone"
description: |-
  This operation submits a request to perform clone of the database instance in Nutanix database service (NDB).
---

# nutanix_ndb_clone

Provides a resource to perform the clone of database instance based on the input parameters. 

## Example Usage

```hcl
## resource for cloning using Point in time given time machine name

    resource "nutanix_ndb_clone" "name" {
        time_machine_name = "test-pg-inst"
        name = "test-inst-tf-check"
        nx_cluster_id = "{{ nx_Cluster_id }}"
        ssh_public_key = "{{ sshkey }}"
        user_pitr_timestamp=  "{{ point_in_time }}"
        time_zone = "Asia/Calcutta"
        create_dbserver = true
        compute_profile_id = "{{ compute_profile_id }}"
        network_profile_id ="{{ network_profile_id }}"
        database_parameter_profile_id =  "{{ databse_profile_id }}"
        nodes{
            vm_name= "test_vm_clone"
            compute_profile_id = "{{ compute_profile_id }}"
            network_profile_id ="{{ network_profile_id }}"
            nx_cluster_id = "{{ nx_Cluster_id }}"
        }
        postgresql_info{
            vm_name="test_vm_clone"
            db_password= "pass"
        }
    }
```

## Argument Reference

* `time_machine_id`: (Optional) time machine id 
* `time_machine_name`: (Optional) time machine name
* `snapshot_id`: (Optional) snapshot id from where clone is created
* `user_pitr_timestamp`:(Optional) point in time for clone to be created
* `time_zone`:(Optional) timezone
* `node_count`: Node count. Default is 1 for single instance
* `nodes`: Nodes contain info about dbservers vm
* `lcm_config`: LCM Config contains the expiry details and refresh details
* `name`: Clone name
* `description`: Clone description
* `nx_cluster_id`: cluster id on where clone will be present
* `ssh_public_key`: ssh public key
* `compute_profile_id`: specify the compute profile id
* `network_profile_id`: specify the network profile id
* `database_parameter_profile_id`: specify the database parameter profile id
* `vm_password`: vm password
* `create_dbserver`: create new dbserver
* `clustered`: clone will be clustered or not
* `dbserver_id`: Specify if you want to create a database server. This value can be set to true or false as required.
* `dbserver_cluster_id`: dbserver cluster id
* `dbserver_logical_cluster_id`: dbserver logical cluster id
* `latest_snapshot`: latest snapshot 
* `postgresql_info`: postgresql info for the clone
* `actionarguments`: (Optional) if any action arguments is required

### nodes

* `vm_name`: name for the database server VM.
* `compute_profile_id`: specify compute profile id
* `network_profile_id`: specify network profile id
* `new_db_server_time_zone`: dbserver time zone
* `nx_cluster_id`: cluster id
* `properties`: properties of vm
* `dbserver_id`: dberver id

### postgresql_info

* `vm_name`: name for the database server VM.
* `dbserver_description`: description for the dbserver.
* `db_password`:  password of the postgres superuser.
* `pre_clone_cmd`:  OS command that you want to run before the instance is created.
* `post_clone_cmd`: OS command that you want to run after the instance is created.

### actionarguments

Structure for each action argument in actionarguments list:

* `name`: - (Required) name of argument
* `value`: - (Required) value for argument


## Attributes Reference

* `owner_id`: owner id
* `date_created`: date created for clone
* `date_modified`: last modified date for clone
* `tags`: allows you to assign metadata to entities (clones, time machines, databases, and database servers) by using tags.
* `clone`: cloned or not
* `era_created`: era created or not
* `internal`: internal
* `placeholder`: placeholder
* `database_name`: database name
* `type`: type of clone
* `database_cluster_type`: database cluster type
* `status`: status of clone
* `database_status`: database status 
* `info`: info of clone
* `group_info`: group info of clone
* `metadata`: metadata about clone
* `metric`: Stores storage info regarding size, allocatedSize, usedSize and unit of calculation that seems to have been fetched from PRISM.
* `category`:category of clone
* `parent_database_id`: parent database id
* `parent_source_database_id`: parent source database id
* `dbserver_logical_cluster`: dbserver logical cluster
* `database_nodes`: database nodes associated with database instance 
* `linked_databases`: linked databases within database instance


See detailed information in [NDB Database Instance](https://www.nutanix.dev/api_references/ndb/#/9a50106e42347-create-clone-using-given-time-machine) .

