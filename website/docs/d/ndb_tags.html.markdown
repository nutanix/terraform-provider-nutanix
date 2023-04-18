---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_tags"
sidebar_current: "docs-nutanix-datasource-ndb-tag"
description: |-
 List of tags in Nutanix Database Service
---

# nutanix_ndb_tags

List of tags in Nutanix Database Service

## Example Usage

```hcl
    data "nutanix_ndb_tags" "tags"{ }
```

## Attribute Reference

The following attributes are exported:

* `tags`: List of tags present in NDB.
* `entity_type`: (Optional) entity type of specific tag. Valid values are DATABASE, TIME_MACHINE, CLONE,DATABASE_SERVER. 

### tags
* `name`:  name for the tag
* `description`: description for the tag
* `entity_type`:  entity for the tag to be associated with.
* `required`: tag value for entities.
* `status`: Status of the tag
* `owner`: owner id of the tag
* `values`: value for the tag
* `date_created`: date created of the tag
* `date_modified`: modified date of tha tag


See detailed information in [NDB Tags](https://www.nutanix.dev/api_references/ndb/#/0a7bf3bdeed86-get-list-of-all-tags).