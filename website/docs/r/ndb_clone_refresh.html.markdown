---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_clone_refresh"
sidebar_current: "docs-nutanix-resource-ndb-clone-refresh"
description: |-
  NDB allows you to create and refresh clones to a point in time either by using transactional logs or by using snapshots. This operation submits a request to perform refresh clone of the database in Nutanix database service (NDB).
---

# nutanix_ndb_clone_refresh

Provides a resource to perform the refresh clone of database based on the input parameters. 

## Example Usage

### resource to refresh clone with snapshot id

```hcl
    resource "nutanix_ndb_clone_refresh" "acctest-managed"{
        clone_id = "{{ clone_id }}"
        snapshot_id = "{{ snapshot_id }}"
        timezone = "Asia/Calcutta"
    }
```

### resource to refresh clone with user pitr timestamp

```hcl
    resource "nutanix_ndb_clone_refresh" "acctest-managed"{
        clone_id = "{{ clone_id }}"
        user_pitr_stamp = "{{ timestamp }}"
        timezone = "Asia/Calcutta"
    }
```

## Argument Reference
* `clone_id`: (Required) clone id
* `snapshot_id`: (Optional) snapshot id where clone has to be refreshed
* `user_pitr_stamp`: (Optional) Point in time recovery where clone has to be refreshed
* `timezone`: (Optional) timezone. Default is Asia/Calcutta. 

See detailed information in [NDB Clone Refresh](https://www.nutanix.dev/api_references/ndb/#/d4e53fff274fa-start-refresh-operation-for-the-given-clone).
