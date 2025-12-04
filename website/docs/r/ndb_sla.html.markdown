---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_sla"
sidebar_current: "docs-nutanix-resource-ndb-sla"
description: |-
  SLAs are data retention policies that allow you to specify how long the daily, weekly, monthly, and quarterly snapshots are retained in NDB. This operation submits a request to create, update and delete slas in Nutanix database service (NDB).
---

# nutanix_ndb_sla

Provides a resource to create SLAs based on the input parameters. 

## Example Usage

```hcl

    resource "nutanix_ndb_sla" "sla" {
        name= "test-sla"
        description = "here goes description"
        
        // Rentention args are optional with default values
        continuous_retention = 30
        daily_retention = 3
        weekly_retention = 2
        monthly_retention= 1
        quarterly_retention=1
    }
```


## Argument Reference
* `name` : (Required) Name of profile
* `description` : (Optional) Description of profile
* `continuous_retention`: (Optional) Duration in days for which transaction logs are retained in NDB.
* `daily_retention`: (Optional) Duration in days for which a daily snapshot must be retained in NDB.
* `weekly_retention`: (Optional) Duration in weeks for which a weekly snapshot must be retained in NDB.
* `monthly_retention`: (Optional) Duration in months for which a monthly snapshot must be retained in NDB
* `quarterly_retention`: (Optional) Duration in number of quarters for which a quarterly snapshot must be retained in NDB.
* `yearly_retention`: (Optional) Not supported as of now. 

## Attributes Reference

* `unique_name`: name of sla
* `owner_id`: owner id
* `system_sla`: refers whether sla is custom or built-in 
* `date_created`: sla created data
* `date_modified`: sla last modified date
* `reference_count`: reference count
* `pitr_enabled`: pitr enabled 
* `current_active_frequency`: slas current frequency 


See detailed information in [NDB SLA](https://www.nutanix.dev/api_references/ndb/#/2fbbab4326e22-create-sla-from-ndb-service).
