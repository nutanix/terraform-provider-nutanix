---
layout: "nutanix"
page_title: "NUTANIX: nutanix_self_service_app_provision"
sidebar_current: "docs-nutanix_self_service_app"
description: |-
  Launches a blueprint to create an application and perform system actions on application.
---

# nutanix_self_service_app_provision

Launches a blueprint to create an application and perform system actions on application.

## Example 1: Provision Application

```hcl
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION TO SET"
    app_description = "DESCRIPTION OF APPLICATION"
}
```

## Example 2: Provision Application with runtime editable
<b>Remark: 

Requires installation of jq in machine.</b>
For ubuntu: sudo apt-get install jq. For detail installation [Read here](https://jqlang.org/download/)

Runtime editable is currently supported for updating `num_sockets, num_vcpus_per_socket, memory_size_mib` Follow below steps:

- Step 1: Extract runtime editable from data source nutanix_blueprint_runtime_editables [Read here](../d/blueprint_runtime_editables.html.markdown)
- Step 2: Dump extracted value in a json [Read here](../d/blueprint_runtime_editables.html.markdown)
- Step 3: Open extracted json file and copy exact value string from substrate_list["value"] 
- Step 4: In terminal. Use jq to format the string using `echo '<copied-value>' | jq -r | jq` and copy the content from terminal.
- Step 5: Make `variables.tf` similar to the one mentioned below and put this jq formatted value in place of `<updated-value>` show below.
- Step 6: Replace value of `num_sockets, num_vcpus_per_socket, memory_size_mib` as per your need.
- Step 7: Provision application

```hcl
data "nutanix_blueprint_runtime_editables" "example" {
    bp_name = "NAME OF BLUEPRINT"
}

# dumps read value into a readable json file
resource "local_file" "dump_runtime_value" {
    content  = jsonencode(data.nutanix_blueprint_runtime_editables.example.runtime_editables)
    filename = "runtime_value.json"
}

[SAMPLE variables.tf]
variable "substrate" {
  type = string
  default = <<EOT
  <updated-value>
EOT 
}

# Launch blueprint and provision your application
resource "nutanix_self_service_app_provision" "test" {
   bp_name         = "NAME OF BLUEPRINT"
   app_name        = "NAME OF APPLICATION TO SET"
   app_description = "DESCRIPTION OF APPLICATION"

   runtime_editables {
    substrate_list {
       name= "VM1"
       value = var.substrate
     }
   }
}
```

## Example 3: Run system action

Step 1: Provision application 

Step 2: use external id of resource (uuid of app) created as input to run system actions on this application.

```hcl
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION TO SET"
    app_description = "DESCRIPTION OF APPLICATION"
}

resource "nutanix_self_service_app_provision" "test" {
    app_uuid        = nutanix_self_service_app_provision.test.id
    action = "NAME OF SYSTEM ACTION TO RUN"
}

# Alternatively you can also run system action by using app name
resource "nutanix_self_service_app_provision" "test" {
    app_name        = "NAME OF APPLICATION"
    action = "NAME OF SYSTEM ACTION TO RUN"
}
```

## Example 4: Soft delete Application

Step 1: Provision application 

Step 2: use external id of resource (uuid of app) created as input to run system actions on this application.

Step 3: set soft_delete attribute as true

Step 4: Run terraform destroy to soft delete application.

```hcl
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION TO SET"
    app_description = "DESCRIPTION OF APPLICATION"
}

resource "nutanix_self_service_app_provision" "test" {
    app_uuid        = nutanix_self_service_app_provision.test.id
    soft_delete     = true
}

```

## Argument Reference

The following arguments are supported:

* `bp_name`: - (Optional) The name of the blueprint to launch.
* `bp_uuid`: - (Optional) The UUID of the blueprint to launch.
* `app_name`: - (Required) The name of the application you want to set.
* `app_description`: - (Optional) The description of application.

Both are `bp_name` and `bp_uuid` are optional but atleast one of them to be provided for this resource to work.

## Attribute Reference

The following attributes are exported:

* `action`: - (Optional) System action to trigger after provisioning. Valid values: ["start", "stop", "restart"]
* `soft_delete`: - (Optional, Default: false) If true, the application is soft-deleted when the resource is destroyed.
* `spec`: - (Computed) Application specification (JSON string).
* `status`: - (Computed) Application status coming as response from server.
* `api_version`: - (Computed) API version used.
* `state`: - (Computed) state of the application (e.g. running, stopped)

### runtime_editables

(Optional) A block of runtime-editable values to override blueprint configuration. Available lists:

* `action_list`: -  A list of actions associated with the blueprint. Each action contains runtime specifications (defined in RuntimeSpec).
* `service_list`: - A list of services associated with the blueprint. Each service contains runtime specifications (defined in RuntimeSpec).
* `credential_list`: -  A list of credentials associated with the blueprint. Each credential contains runtime specifications (defined in RuntimeSpec).
* `substrate_list`: - A list of substrates associated with the blueprint. Each substrate contains runtime specifications (defined in RuntimeSpec).
* `package_list`: -  A list of packages associated with the blueprint. Each package contains runtime specifications (defined in RuntimeSpec).
* `snapshot_config_list`: - A list of snapshot configurations associated with the blueprint. Each snapshot configuration contains runtime specifications (defined in RuntimeSpec).
* `app_profile`: -  A list of application profiles associated with the blueprint. Each application profile contains runtime specifications (defined in RuntimeSpec).
* `task_list`: - A list of tasks associated with the blueprint. Each task contains runtime specifications (defined in RuntimeSpec).
* `restore_config_list`: -  A list of restore configurations associated with the blueprint. Each restore configuration contains runtime specifications (defined in RuntimeSpec).
* `variable_list`: - A list of variables associated with the blueprint. Each variable contains runtime specifications (defined in RuntimeSpec).
* `deployment_list`: -  A list of deployments associated with the blueprint. Each deployment contains runtime specifications (defined in RuntimeSpec).

### RuntimeSpec

The RuntimeSpec function defines the runtime specifications for each of the entities listed in the runtime_editables block. Each runtime specification contains the following attributes:

- `description`: (Optional, Computed) A textual description of the runtime specification. This field provides additional information or context about the entity.
- `value`: (Optional, Computed) The value associated with the runtime specification. This can be a string value representing a configuration or setting.
- `name`: (Optional, Computed) The name of the runtime specification. This could be the name of an action, service, credential, or other runtime-editable resource.
- `type`: (Optional, Computed) The type of runtime specification. This field indicates the category or classification of the runtime resource, such as an action, service, or credential.
- `uuid`: (Optional, Computed) The unique identifier (UUID) associated with the runtime specification. This is useful for identifying specific resources or entities.
- `context`: (Optional, Computed) The context in which the runtime specification is applied. It is full address of where the entity in target is present.


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

Provides an overview of the provisioned application:

- `application_uuid`: The UUID of the application.
- `blueprint`: The blueprint associated with the application.
- `application_profile`: The profile assigned to the application.
- `project`: The project associated with the application.
- `owner`: The owner of the application.
- `created_on`: The timestamp when the application was created.
- `last_updated_on`: The timestamp when the application was last updated.

### actions

List of available actions on the app:

- `name`: name of action.
- `uuid`: UUID of the action.
- `description`: description of the action

See detailed information in [Launch a Blueprint](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/Blueprints/paths/~1blueprints~1%7Buuid%7D~1simple_launch/post).