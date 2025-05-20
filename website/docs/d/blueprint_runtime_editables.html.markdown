---
layout: "nutanix"
page_title: "NUTANIX: nutanix_blueprint_runtime_editables"
sidebar_current: "docs-nutanix_self_service_app"
description: |-
  Describes runtime editables that are present in a blueprint.
---

# nutanix_blueprint_runtime_editables

Describes runtime editables that are present in a blueprint.

## Example Usage

```hcl
data "nutanix_blueprint_runtime_editables" "example" {
    bp_name = "NAME OF BLUEPRINT"
}

# dumps read value into a readable json file
resource "local_file" "dump_runtime_value" {
    content  = jsonencode(data.nutanix_blueprint_runtime_editables.example.runtime_editables)
    filename = "runtime_value.json"
}
```

## Argument Reference

The following arguments are supported:

* `bp_uuid`: - (Optional) The UUID of the blueprint for which runtime editables will be listed. If this is provided, it will return runtime editables for the specified blueprint.
* `bp_name`: - (Optional) The name of the blueprint for which runtime editables will be listed. If this is provided, it will return runtime editables for the specified blueprint.

Both (`bp_uuid` and `bp_name`) are optional but atleast one of them to be provided for this data source to work.

## Attribute Reference

The following attributes are exported:

### runtime_editables

The runtime_editables block contains a list of runtime-editable items associated with the blueprint. Each runtime-editable item contains the following attributes:

* `action_list`: -  A list of actions associated with the blueprint. Each action contains runtime specifications (defined in RuntimeSpecDS).

* `service_list`: - A list of services associated with the blueprint. Each service contains runtime specifications (defined in RuntimeSpecDS).

* `credential_list`: -  A list of credentials associated with the blueprint. Each credential contains runtime specifications (defined in RuntimeSpecDS).

* `substrate_list`: - A list of substrates associated with the blueprint. Each substrate contains runtime specifications (defined in RuntimeSpecDS).

* `package_list`: -  A list of packages associated with the blueprint. Each package contains runtime specifications (defined in RuntimeSpecDS).

* `snapshot_config_list`: - A list of snapshot configurations associated with the blueprint. Each snapshot configuration contains runtime specifications (defined in RuntimeSpecDS).

* `app_profile`: -  A list of application profiles associated with the blueprint. Each application profile contains runtime specifications (defined in RuntimeSpecDS).

* `task_list`: - A list of tasks associated with the blueprint. Each task contains runtime specifications (defined in RuntimeSpecDS).

* `restore_config_list`: -  A list of restore configurations associated with the blueprint. Each restore configuration contains runtime specifications (defined in RuntimeSpecDS).

* `variable_list`: - A list of variables associated with the blueprint. Each variable contains runtime specifications (defined in RuntimeSpecDS).

* `deployment_list`: -  A list of deployments associated with the blueprint. Each deployment contains runtime specifications (defined in RuntimeSpecDS).

### RuntimeSpecDS

The RuntimeSpecDS function defines the runtime specifications for each of the entities listed in the runtime_editables block. Each runtime specification contains the following attributes:

- `description`: (Optional, Computed) A textual description of the runtime specification. This field provides additional information or context about the entity.
- `value`: (Optional, Computed) The value associated with the runtime specification. This can be a string value representing a configuration or setting.
- `name`: (Optional, Computed) The name of the runtime specification. This could be the name of an action, service, credential, or other runtime-editable resource.
- `type`: (Optional, Computed) The type of runtime specification. This field indicates the category or classification of the runtime resource, such as an action, service, or credential.
- `uuid`: (Optional, Computed) The unique identifier (UUID) associated with the runtime specification. This is useful for identifying specific resources or entities.
- `context`: (Optional, Computed) The context in which the runtime specification is applied. It is full address of where the entity in target is present.

See detailed information in [Runtime Editables in Blueprint](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/Blueprints/paths/~1blueprints~1%7Buuid%7D~1runtime_editables/get).
