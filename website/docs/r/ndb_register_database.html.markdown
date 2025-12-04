---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_register_database"
sidebar_current: "docs-nutanix-resource-ndb-database-register"
description: |-
    It helps to register a source (production) database running on a Nutanix cluster with NDB. When you register a database with NDB, the database server VM (VM that hosts the source database) is also registered with NDB. After you have registered a database with NDB, a time machine is created for that database.
    This operation submits a request to register the database in Nutanix database service (NDB).
---

# nutanix_ndb_register_database

Provides a resource to register the database based on the input parameters. 

## Example Usage

```hcl

    // register PostgreSQL database with registered DBServer VM

    resource "nutanix_ndb_register_database" "name" {
        database_type = "postgres_database"
        database_name=  "test-inst"
        description = "added by terraform"
        category =  "DEFAULT"

        // registered vm IP
        vm_ip = "{{ vm_ip }}"

        // optional 
        working_directory= "/tmp"
        
        reset_description_in_nx_cluster= false

        // time Machine Info
        time_machine_info {
            name= "test-pg-inst-regis"
            description= "description of tms"
            slaid=" {{ SLA ID}}"
            schedule {
                snapshottimeofday{
                    hours= 16
                    minutes= 0
                    seconds= 0
                }
                continuousschedule{
                    enabled=true
                    logbackupinterval= 30
                    snapshotsperday=1
                }
                weeklyschedule{
                    enabled=true
                    dayofweek= "WEDNESDAY"
                }
                monthlyschedule{
                    enabled = true
                    dayofmonth= "27"
                }
                quartelyschedule{
                    enabled=true
                    startmonth="JANUARY"
                    dayofmonth= 27
                }
                yearlyschedule{
                    enabled= false
                    dayofmonth= 31
                    month="DECEMBER"
                }
            }
        }
        postgress_info{

            // required args
            listener_port= "5432"
            db_password ="pass"
            db_name= "testdb1"

            // Optional with default values
            db_user= "postgres"
            backup_policy= "prefer_secondary"
            postgres_software_home= "{{ directory where the PostgreSQL database software is installed.}}"
            software_home= "{{ directory where the PostgreSQL database software is installed. }}"
        
        }
    }


    // register PostgreSQL database with instance not registered on VM
    resource "nutanix_ndb_register_database" "name" {
        database_type = "postgres_database"
        database_name=  "test-inst"
        description = "added by terraform"
        category =  "DEFAULT"
        nx_cluster_id = "{{ cluster_ID }}"

        // registered vm info
        vm_ip = "{{ vm_ip }}"
        vm_username = "{{ vm_username }}"
        vm_password = "{{ vm_password }}"

        // optional 
        working_directory= "/tmp"
        
        reset_description_in_nx_cluster= false

        // time Machine Info
        time_machine_info {
            name= "test-pg-inst-regis"
            description= "description of tms"
            slaid=" {{ SLA ID}}"
            schedule {
            snapshottimeofday{
                hours= 16
                minutes= 0
                seconds= 0
            }
            continuousschedule{
                enabled=true
                logbackupinterval= 30
                snapshotsperday=1
            }
            weeklyschedule{
                enabled=true
                dayofweek= "WEDNESDAY"
            }
            monthlyschedule{
                enabled = true
                dayofmonth= "27"
            }
            quartelyschedule{
                enabled=true
                startmonth="JANUARY"
                dayofmonth= 27
            }
            yearlyschedule{
                enabled= false
                dayofmonth= 31
                month="DECEMBER"
            }
            }
        }
        postgress_info{

            // required args
            listener_port= "5432"
            db_password ="pass"
            db_name= "testdb1"

            // Optional with default values
            db_user= "postgres"
            backup_policy= "prefer_secondary"
            postgres_software_home= "{{ directory where the PostgreSQL database software is installed }}"
        }
    }

```


## Argument Reference

* `database_type`: (Required) type of database. Required value: postgres_database 
* `database_name`: (Required) name of database
* `description`: (Optional) description
* `clustered`: (Optional) clustered or not. Default is false
* `forced_install`: (Optional) forced install. Default:  true
* `category`: (Optional) category of database. Default is "DEFAULT"
* `vm_ip`: (Required) IP address of dbserver VM
* `vm_username`: (Optional) username of the NDB drive user account that has sudo access.
* `vm_password`: (Optional) password of the NDB drive user account.
* `vm_sshkey`: (Optional) ssh key for vm
* `vm_description`: (Optional) description for VM
* `nx_cluster_id`:(Optional) cluster on which NDB is present
* `reset_description_in_nx_cluster`: (Optional) Reset description in cluster
* `auto_tune_staging_drive`: (Optional) auto tune staging drive. Default is true
* `working_directory`: (Optional) working directory. Default is /tmp 
* `time_machine_info`: (Required) Time Machine info
* `tags`: (Optional) tags 
* `actionarguments`: (Optional) action arguments
* `postgress_info`:  (Optional) Postgress_Info for registering. 


* `delete`:- (Optional) Delete the database from the VM. Default value is false
* `remove`:- (Optional) Unregister the database from NDB. Default value is true
* `soft_remove`:- (Optional) Soft remove. Default will be false
* `forced`:- (Optional) Force delete of instance. Default is false
* `delete_time_machine`:- (Optional) Delete the database's Time Machine (snapshots/logs) from the NDB. Default value is true
* `delete_logical_cluster`:- (Optional) Delete the logical cluster. Default is true

### postgress_info

* `listener_port`: (Required) listner port of database
* `db_password`: (Required) database password
* `db_name`: (Required) name of the database server VM on which the instance you want to register is running.
* `db_user`: (Optional) username of the NDB drive user account that has sudo access. 
* `switch_log`: (Optional) switch log of database. Default is true
* `allow_multiple_databases`: (Optional) allow multiple databases. Default is true
* `backup_policy`: (Optional) backup policy of database. Default is prefer_secondary.
* `vm_ip`: (Optional) VM IP of the database server VM on which the instance you want to register is running.
* `postgres_software_home`: (Required) path to the PostgreSQL home directory in which the PostgreSQL software is installed.
* `software_home`: (Optional) path to the directory in which the PostgreSQL software is installed.

### time_machine_info

The timemachineinfo attribute supports the following:

* `name`: - (Required) name of time machine
* `description`: - (Optional) description of time machine
* `slaid`: - (Optional) SLA ID for single instance 
* `sla_details`:-  (optional) SLA details for HA instance
* `autotunelogdrive`: - (Optional) enable auto tune log drive. Default: true
* `schedule`: - (Optional) schedule for snapshots
* `tags`: - (Optional) tags

### sla_details

* `primary_sla`:- (Required) primary sla details
* `primary_sla.sla_id` :- (Required) sla id
* `primary_sla.nx_cluster_ids` -: (Optioanl) cluster ids


### schedule

The schedule attribute supports the following:

* `snapshottimeofday`: - (Optional) daily snapshot config
* `continuousschedule`: - (Optional) snapshot freq and log config
* `weeklyschedule`: - (Optional) weekly snapshot config
* `monthlyschedule`: - (Optional) monthly snapshot config
* `quartelyschedule`: - (Optional) quaterly snapshot config
* `yearlyschedule`: - (Optional) yearly snapshot config


### actionarguments

Structure for each action argument in actionarguments list:

* `name`: - (Required) name of argument
* `value`: - (Required) value for argument


## Attributes Reference

* `name`: Name of database instance
* `description`: description of database instance
* `databasetype`: type of database
* `properties`: properties of database created
* `date_created`: date created for db instance
* `date_modified`: date modified for instance
* `tags`: allows you to assign metadata to entities (clones, time machines, databases, and database servers) by using tags.
* `clone`: whether instance is cloned or not
* `database_name`: name of database
* `type`: type of database
* `database_cluster_type`: database cluster type
* `status`: status of instance
* `database_status`: status of database
* `dbserver_logical_cluster_id`: dbserver logical cluster id
* `time_machine_id`: time machine id of instance 
* `parent_time_machine_id`: parent time machine id
* `time_zone`: timezone on which instance is created xw
* `info`: info of instance
* `metric`: Stores storage info regarding size, allocatedSize, usedSize and unit of calculation that seems to have been fetched from PRISM.
* `category`: category of instance
* `parent_database_id`: parent database id
* `parent_source_database_id`: parent source database id
* `lcm_config`: LCM config of instance
* `time_machine`: Time Machine details of instance
* `dbserver_logical_cluster`: dbserver logical cluster
* `database_nodes`: database nodes associated with database instance 
* `linked_databases`: linked databases within database instance


See detailed information in [NDB Register Database](https://www.nutanix.dev/api_references/ndb/#/40355715d4188-register-an-existing-database).