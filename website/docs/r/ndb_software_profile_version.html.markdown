---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_software_version_profile"
sidebar_current: "docs-nutanix-resource-ndb-software-version-profile"
description: |-
  This operation submits a request to create, update and delete software profile versions in Nutanix database service (NDB).
---

# nutanix_ndb_software_version_profile

Provides a resource to create software profile versions based on the input parameters. 

## Example Usage

```hcl
    
    resource "nutanix_ndb_software_version_profile" "name" {
      engine_type = "postgres_database"
      profile_id= resource.nutanix_ndb_profile.name12.id
      name = "test-tf"
      description= "made  by tf"
      postgres_database{
        source_dbserver_id = "{{ DB_Server_ID }}"
      }
      available_cluster_ids =  ["{{ cluster_ids }}"]
      status = "published"
    }
```

## Argument Reference

* `profile_id`: (Required) profile id
* `name`: Name of profile
* `description`: description of profile
* `engine_type`: engine type of profile
* `status`: status of profile. Allowed Values are "deprecated", "published", "unpublished"
* `postgres_database`: postgres database info
* `available_cluster_ids`: available cluster ids

### postgres_database

* `source_dbserver_id`: (Optional) source dbserver id
* `os_notes`: (Optional) os notes for software profile
* `db_software_notes`: (Optional) db software notes

## Attributes Reference

* `status`: status of profile
* `owner`: owner  of profile
* `db_version`: Db version of software profile
* `topology`: topology of software profile
* `system_profile`: system profile or not. 
* `version`: Version of software profile
* `published`: Published or not
* `deprecated`: deprecated or not
* `properties`: properties of profile
* `properties_map`: properties map of profile 
* `version_cluster_association`: version cluster association

### version_cluster_association

* `nx_cluster_id`: nutanix cluster id
* `date_created`: date created of profile
* `date_modified`: date modified of profile
* `owner_id`: owner id
* `status`: status of version
* `profile_version_id`: profile version id
* `properties`: properties of software profile
* `optimized_for_provisioning`: version optimized for provisioning


### properties
* `name`: name of property
* `value`: value of property
* `secure`: secure or not

See detailed information in [NDB Profile version](https://www.nutanix.dev/api_references/ndb/#/351a7caf34bbb-create-profile-version).

