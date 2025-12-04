---
layout: "nutanix"
page_title: "NUTANIX: nutanix_template_version_v2"
sidebar_current: "docs-nutanix-resource-template-version-v2"
description: |-
  Performs Guest OS actions on given template.
---

# nutanix_template_version_v2

Performs Guest OS actions on given template. It Initiates, Completes and Cancels the Guest OS operation.

## Example

```hcl
resource "nutanix_template_guest_os_actions_v2" "example-1"{
    ext_id = "ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
    action = "initiate"
    version_id = "c2c249b0-98a0-43fa-9ff6-dcde578d3936"
}

resource "nutanix_template_guest_os_actions_v2" "example-2"{
    ext_id = "8a938cc5-282b-48c4-81be-de22de145d07"
    action = "complete"
    version_name = "version_name"
    version_description = "version desc"
    is_active_version = true
}

resource "nutanix_template_guest_os_actions_v2" "example-3"{
    ext_id = "1cefd0f5-6d38-4c9b-a07c-bdd2db004224"
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


See detailed information in [Nutanix Template Guest OS Action Initiate V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Templates/operation/initiateGuestUpdate).
See detailed information in [Nutanix Template Guest OS Action Complete V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Templates/operation/completeGuestUpdate).
See detailed information in [Nutanix Template Guest OS Action Cancel V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Templates/operation/cancelGuestUpdate).

