---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_authorize_dbserver"
sidebar_current: "docs-nutanix-resource-ndb-authorize-dbserver"
description: |-
  This operation submits a request to authorize db server VMs for cloning of the database instance in Nutanix database service (NDB).
---

# nutanix_ndb_authorize_dbserver

Provides a resource to authorize db server VMs for cloning of database instance based on the input parameters. 

## Example Usage

```hcl

    resource "nutanix_ndb_authorize_dbserver" "name" {
        time_machine_name = "test-pg-inst"
        dbservers_id = [
            "{{ dbServer_IDs}}"
        ]
    }
```

## Arguments Reference

* `time_machine_id`: (Optional)
* `time_machine_name`: (Optional)
* `dbservers_id `: (Required)
