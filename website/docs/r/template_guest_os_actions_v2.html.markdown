---
layout: "nutanix"
page_title: "NUTANIX: nutanix_template_version_v4"
sidebar_current: "docs-nutanix-resource-template-version-v2"
description: |-
  Performs Guest OS actions on given template. 
---

# nutanix_template_version_v4

Performs Guest OS actions on given template. It Initiates, Completes and Cancels the Guest OS operation. 

## Example 

```hcl
resource "nutanix_template_guest_os_actions_v2" "example-1"{
    ext_id = {{ template uuid }}
    action = "initiate"
    version_id = {{  template version id}}
}

resource "nutanix_template_guest_os_actions_v2" "example-2"{
    ext_id = {{ template uuid }}
    action = "complete"
    version_name = "version name"
    version_description = "version desc"
    is_active_version = true
}

resource "nutanix_template_guest_os_actions_v2" "example-3"{
    ext_id = {{ template uuid }}
    action = "cancel"
}
```


## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) The identifier of a Template.
* `action`: (Required) Actions to be performed. Acceptable values are "initiate", "complete", "cancel" . 

* `version_id`: (Required) The identifier of a Template Version. Only applicable with `Initiate` action. 
* `version_name`: (Required) The user defined name of a Template Version. Only applicable with `complete` action.
* `version_description`: The user defined description of a Template Version. (Required) Only applicable with `complete` action.
* `is_active_version`: (Optional) Specify whether to mark the Template Version as active or not. The newly created Version during Template Creation, Updating or Guest OS Updating is set to Active by default unless specified otherwise. Default is true. Only applicable with `complete` action.


See detailed information in [Nutanix Template Guest OS Action V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0).