---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_register_dbserver"
sidebar_current: "docs-nutanix-resource-ndb-dbservervm-register"
description: |-
  This operation submits a request to register the database server VM in Nutanix database service (NDB).
  Note: For 1.8.0 release, only postgress database type is qualified and officially supported.
---

# nutanix_ndb_register_dbserver

Provides a resource to register database server VMs based on the input parameters. For 1.8.0 release, only postgress database type is qualified and officially supported.

## Example Usage

```hcl
    resource "nutanix_ndb_register_dbserver" "name" {
        database_type = "postgres_database"
        vm_ip= "{{ vmip to register}}"
        nxcluster_id = "{{ cluster_id }}"
        username= "{{ username of the NDB drive user account }}"
        password="{{ password of the NDB drive user account }}"
        postgres_database{
          listener_port  = {{ listner_port }}
          postgres_software_home= "{{ path to the PostgreSQL home directory in which the PostgreSQL software is installed }}"
        }
    }
```


## Argument Reference

The following arguments are supported:
* `database_type`: (Required) database type i.e. postgres_database
* `vm_ip`: (Required) IP address of the database server VM
* `nxcluster_id`: (Required) cluster on which you want to register the database server VM.
* `username`: (Required) username of the NDB drive user account that has sudo access
* `password`: (Optional) password of the NDB drive user account. Conflicts with ssh_key.
* `ssh_key`: (Optional) the private key. Conflicts with password.
* `postgres_database`: (Optional) postgres info for dbserver

* `name`: (Optional) Name of db server vm. Should be used in Update Method only. 
* `description`: (Optional) description of db server vm. Should be used in update Method only . 
* `update_name_description_in_cluster`: (Optional) Updates the name and description in cluster. Should be used in Update Method only. 
* `working_directory`: (Optional) working directory of postgres. Default is "/tmp"
* `forced_install`: (Optional) forced install the packages. Default is true 

* `delete`:- (Optional) Delete the VM and associated storage. Default value is false
* `remove`:- (Optional) Unregister the database from NDB. Default value is true
* `soft_remove`:- (Optional) Soft remove. Default will be false
* `delete_vgs`:- (Optional) Delete volume grous. Default value is true
* `delete_vm_snapshots`:- (Optional) Delete the vm snapshots. Default is true


### postgres_database
* `listener_port`: (Optional) listener port of db server
* `postgres_software_home`: (Required) path to the PostgreSQL home directory in which the PostgreSQL software is installed 


### actionarguments

Structure for each action argument in actionarguments list:

* `name`: name of the dbserver vm
* `properties`: Properties of dbserver vm
* `era_created`: created by era or not.
*  `internal`: is internal or not.
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


See detailed information in [NDB Register Database Server VM](https://www.nutanix.dev/api_references/ndb/#/5bd6f03bd6ed7-register-an-existing-database-server).