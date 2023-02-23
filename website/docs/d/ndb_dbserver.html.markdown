---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_dbserver"
sidebar_current: "docs-nutanix-datasource-ndb-dbserver"
description: |-
 Describes Database Server VM in Nutanix Database Service
---

# nutanix_ndb_dbserver

Describes Database Server VM in Nutanix Database Service

## Example Usage

```hcl
    data "nutanix_ndb_dbserver" "dbserver" {
        id = "{{ dbserver_vm_id }}"
    }   

    data "nutanix_ndb_dbserver" "dbserver" {
        name = "{{ dbserver_vm_name }}"
    }

    data "nutanix_ndb_dbserver" "dbserver" {
        ip = "{{ dbserver_vm_ip }}"
    }
```

## Argument Reference

The following arguments are supported:

* `id`: (Optional) id of database server vm
* `name`: (Optional) name of database server vm
* `ip`: (Optional) ip of database server vm
* `vm_cluster_name`: (Optional) vm cluster name of database server
* `vm_cluster_id`: (Optional) vm cluster id of database server

## Attribute Reference

The following attributes are exported:

* `description`: description of db server vm
* `date_created`: date created of db server vm
* `date_modified`: date modified of db server vm
* `access_level`: access level
* `properties`: properties of db server vm
* `tags`: tags for db server vm
* `vm_cluster_uuid`: clusetr uuid for dbserver vm
* `ip_addresses`: IP addresses of the dbserver vm
* `mac_addresses`: Mac addresses of dbserver vm
* `type`: Type of entity. i.e. Dbserver
* `status`: Status of Dbserver . Active or not.
* `client_id`:  client id
* `era_drive_id`: era drive id
* `era_version`: era version
* `vm_timezone`:  timezone of dbserver vm
* `vm_info`: info of dbserver vm
* `clustered`: clustered or not
* `is_server_driven`: is server down or not
* `protection_domain_id`: protection domain id
* `query_count`: query count
* `database_type`: database type
* `dbserver_invalid_ea_state`: dbserver invalid ea state
* `working_directory`: working directory of db server vm
* `valid_diagnostic_bundle_state`: valid diagnostic bundle state
* `windows_db_server`: window db server
* `associated_time_machine_ids`: associated time machines ids
* `access_key_id`: access key id of dbserver vm


See detailed information in [Database Server VM](https://www.nutanix.dev/api_references/ndb/#/c531f4158d5f5-get-a-database-server-by-value-type).
