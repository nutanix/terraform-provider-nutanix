---
layout: "nutanix"
page_title: "NUTANIX: nutanix_self_service_app"
sidebar_current: "docs-nutanix_self_service_app"
description: |-
  Describes a self service app.
---

# nutanix_self_service_app

Describes an Application (app) in NCM Self Service.

## Example Usage

```hcl
resource "nutanix_self_service_app_provision" "test" {
	bp_name         = "NAME OF BLUEPRINT IN SERVER"
	app_name        = "NAME OF APP YOU WANT TO SET"
	app_description = "DESCRIPTION OF APP"
}

data "nutanix_self_service_app" "test"{
	app_uuid = nutanix_self_service_app_provision.test.id
}
```

## Argument Reference

The following arguments are supported:

* `app_uuid`: - (Required) The UUID of an app you want to get data from.

## Attribute Reference

The following attributes are exported:

* `app_name`: - The name of the Self Service Application.
* `app_description`: - The description of the Self Service Application.
* `spec`: - The specification of the Self Service Application.
* `status`: - The current status of the application.
* `api_version`: - The API version used for the application.
* `state`: - The state of the application (e.g., Running, Provisioning, etc.).

### vm

The vm block contains a list of virtual machines associated with the Self Service application. Each virtual machine has the following attributes:

* `configuration`: -  Configuration details for the VM.
	- `name`: -   Name of the VM.
	- `ip_address`: -   IP address of the VM.
	- `vcpus`: -   Number of virtual CPUs assigned to the VM.
	- `cores`: -   Number of CPU cores.
	- `memory`: -  Memory allocated to the VM.
	- `vm_uuid`: -   The UUID of the VM.
	- `image`: -   The VM image used.

* `nics`: -  A list of network interfaces attached to the VM.
	- `mac_address`: -   The MAC address of the VM's network interface.
	- `type`: -   The type of network interface.
	- `subnet`: -   The subnet the VM's network interface is attached to.

* `cluster_info`: -  Cluster-related information for the VM.
	- `cluster_name`: -   The name of the cluster.
	- `cluster_uuid`: -   The UUID of the cluster.

* `categories`: -  A map of categories applied to the VM. Each key is a category name, and the value is the category value.

### app_summary

- `application_uuid`: The UUID of the application.
- `blueprint`: The blueprint associated with the application.
- `application_profile`: The profile assigned to the application.
- `project`: The project associated with the application.
- `owner`: The owner of the application.
- `created_on`: The timestamp when the application was created.
- `last_updated_on`: The timestamp when the application was last updated.

### actions
   - `name`: The name of the action.
   - `uuid`: The UUID of the action.
   - `description`: A description of the action.


See detailed information in [NCM Self Service Apps](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/Apps).
