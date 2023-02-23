---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_tag"
sidebar_current: "docs-nutanix-resource-ndb-tag"
description: |-
  NDB allows you to assign metadata to entities (clones, time machines, databases, and database servers) by using tags. When you are cloning a database, you can associate tags with the database that you are creating. This operation submits a request to create, update and delete tags in Nutanix database service (NDB).
---

# nutanix_ndb_tag

Provides a resource to create tags based on the input parameters. 

## Example Usage

### resource to create tag
```hcl
    resource "nutanix_ndb_tag" "name" {
        name= "testst-up"
        description = "this is desc ok"
        entity_type = "DATABASE"
        required=true
    }
```

### resource to update tag with status
```hcl
    resource "nutanix_ndb_tag" "name" {
        name= "testst-up"
        description = "this is desc ok"
        entity_type = "DATABASE"
        required=true
        status = "DEPRECATED"
    }
```

## Argument Reference
* `name`: (Required) name for the tag
* `description`: (Optional) description for the tag
* `entity_type`: (Required) entity for the tag to be associated with. Supported values [ DATABASE, TIME_MACHINE, CLONE, DATABASE_SERVER ]. 
* `required`: (Optional) provide a tag value for entities.

* `status`: (Optional)Status of the tag. Supported values are [ ENABLED, DEPRECATED ]


## Attributes Reference
* `owner`: owner id of the tag
* `values`: value for the tag
* `date_created`: date created of the tag
* `date_modified`: modified date of tha tag


See detailed information in [NDB Tag](https://www.nutanix.dev/api_references/ndb/#/5d6a2dc1bc153-create-a-tag).