---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_tag"
sidebar_current: "docs-nutanix-datasource-ndb-tag"
description: |-
 Describes a tag in Nutanix Database Service
---

# nutanix_ndb_tag

Describes a tag in Nutanix Database Service

## Example Usage

```hcl
    data "nutanix_ndb_tag" "tag"{
        id = "{{ tag id }}"
    }
```


## Argument Reference

The following arguments are supported:
* `id` : (Optional) tag id. Conflicts with  name.
* `name`: (Optional) tag name. Conflicts with id.

## Attribute Reference

The following attributes are exported:

* `name`:  name for the tag
* `description`: description for the tag
* `entity_type`:  entity for the tag to be associated with.
* `required`: tag value for entities.
* `status`: Status of the tag
* `owner`: owner id of the tag
* `values`: value for the tag
* `date_created`: date created of the tag
* `date_modified`: modified date of tha tag


See detailed information in [NDB Tag](https://www.nutanix.dev/api_references/ndb/#/0a7bf3bdeed86-get-list-of-all-tags).