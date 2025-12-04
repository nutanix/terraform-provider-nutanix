---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_profiles"
sidebar_current: "docs-nutanix-datasource-ndb-profiles"
description: |-
 List profiles in Nutanix Database Service
---

# nutanix_ndb_profile

List profiles in Nutanix Database Service

## Example Usage

```hcl
data "nutanix_ndb_profiles" "profiles" {}

output "profiles_list" {
 value = data.nutanix_ndb_profiles.profiles
}

```

## Argument Reference

The following arguments are supported:

* `engine`: Database engine. For eg. postgres_database
* `profile_type`: profile type. Types: Software, Compute, Network and Database_Parameter

## Attribute Reference

The following attributes are exported:

* `profiles`: List of profiles 

## profiles

The following attributes are exported for each profile:

* `id`: - id of profile
* `name`: - profile name
* `description`: - description of profile
* `status`: - status of profile
* `owner`: - owner name
* `engine_type`: - database engine type
* `db_version`: - database version
* `topology`: - topology
* `system_profile`: - if system profile or not
* `assoc_db_servers`: - associated DB servers
* `assoc_databases`: - associated databases
* `latest_version`: - latest version for engine software
* `latest_version_id`: - ID of latest version for engine software
* `versions`: - profile's different version config
* `cluster_availability`: - list of clusters availability
* `nx_cluster_id`: - era cluster ID

See detailed information in [Nutanix Database Service Profiles](https://www.nutanix.dev/api_references/ndb/#/74ae456d63b24-get-all-profiles).
