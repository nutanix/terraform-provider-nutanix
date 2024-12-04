---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_categories_v4"
sidebar_current: "docs-nutanix-resource-vm-categories-v2"
description: |-
   Associate/Disassociate categories to a Virtual Machine.
---

# nutanix_vm_categories_v4

Associate and Disassociate categories to a Virtual Machine.


## Example

```hcl

    data "nutanix_categories_v4" "ctg"{}

    resource "nutanix_vm_categories_v4" "ctgRes" {
        ext_id = {{ vm uuid }}
        categories{
            ext_id = data.nutanix_categories_v4.ctg.categories.2.ext_id
        }
        categories{
            ext_id = data.nutanix_categories_v4.ctg.categories.5.ext_id
        }
    }

```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) The globally unique identifier of a VM. It should be of type UUID
* `categories`: List of category ext_ids.


See detailed information in [Nutanix VM Categories](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).

