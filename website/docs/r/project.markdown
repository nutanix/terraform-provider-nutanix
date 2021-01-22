---
layout: "nutanix"
page_title: "NUTANIX: nutanix_project"
sidebar_current: "docs-nutanix-resource-project"
description: |-
  Provides a Nutanix Category key resource to Create a Project.
---

# nutanix_project

Provides a Nutanix Project resource to Create a Project.

## Example Usage

```hcl
resource "nutanix_subnet" "subnet" {
  cluster_uuid       = "<YOUR_CLUSTER_ID>"
  name               = "sunet_test_name"
  description        = "Description of my unit test VLAN"
  vlan_id            = 31
  subnet_type        = "VLAN"
  subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length      = 24

  dhcp_options = {
    boot_file_name   = "bootfile"
    domain_name      = "nutanix"
    tftp_server_name = "10.250.140.200"
  }

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

resource "nutanix_project" "project_test" {
  name        = "my-project"
  description = "This is my project"

  categories {
    name  = "Environment"
    value = "Staging"
  }

  resource_domain {
    resources {
      limit         = 4
      resource_type = "STORAGE"
    }
  }

  default_subnet_reference {
    uuid = nutanix_subnet.subnet.metadata.uuid
  }

  api_version = "3.1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the project.
* `description` - (Required) A description for project.

### Resource Domain
* `resource_domain` - (Required) The status for a resource domain (limits and values)
* `resource_domain.resources` - (Required) Array of the utilization/limit for resource types
* `resource_domain.resources.#.limit` - (Required) The resource consumption limit (unspecified is unlimited)
* `resource_domain.resources.#.resource_type` - (Required) The type of resource (for example storage, CPUs)

### Account Reference List
* `account_reference_list` - (Optional/Computed) List of accounts associated with the project.
* `account_reference_list.#.kind` - (Optional) The kind name. Default value is `account`
* `account_reference_list.#.uuid` - (Required) The UUID of an account.
* `account_reference_list.#.name` - (Optional/Computed) The name of an account.

### Environment Reference List
* `environment_reference_list` - (Optional/Computed) List of environments associated with the project.
* `environment_reference_list.#.kind` - (Optional) The kind name. Default value is `environment`
* `environment_reference_list.#.uuid` - (Required) The UUID of an environment.
* `environment_reference_list.#.name` - (Optional/Computed) The name of an environment.

### Default Subnet Reference Map
* `default_subnet_reference` - (Required) Reference to a subnet.
* `default_subnet_reference.kind` - (Optional) The kind name. Default value is `subnet`
* `default_subnet_reference.uuid` - (Required) The UUID of a subnet.
* `default_subnet_reference.name` - (Optional/Computed) The name of a subnet.

### user_reference_list
* `user_reference_list` - (Optional/Computed) List of users in the project.
* `user_reference_list.#.kind` - (Optional) The kind name. Default value is `user`
* `user_reference_list.#.uuid` - (Required) The UUID of a user
* `user_reference_list.#.name` - (Optional/Computed) The name of a user.

### External User Group Reference List
* `external_user_group_reference_list` - (Optional/Computed) List of directory service user groups. These groups are not managed by Nutanix.
* `external_user_group_reference_list.#.kind` - (Optional) The kind name. Default value is `user_group`
* `external_user_group_reference_list.#.uuid` - (Required) The UUID of a user_group
* `external_user_group_reference_list.#.name` - (Optional/Computed) The name of a user_group

### Subnet Reference List
* `subnet_reference_list` - (Optional/Computed) List of subnets for the project.
* `subnet_reference_list.#.kind` - (Optional) The kind name. Default value is `subnet`
* `subnet_reference_list.#.uuid` - (Required) The UUID of a subnet
* `subnet_reference_list.#.name` - (Optional/Computed) The name of a subnet.

### External Network List
* `external_network_list` - (Optional/Computed) List of external networks associated with the project.
* `external_network_list.#.uuid` - (Required) The UUID of a network.
* `external_network_list.#.name` - (Optional/Computed) The name of a network.


## Attributes Reference
The following attributes are exported:

### Resource Domain
* `resource_domain.resources.#.units` - The units of the resource type
* `resource_domain.resources.#.value` - The amount of resource consumed

### Metadata
The metadata attribute exports the following:

* `last_update_time` - UTC date and time in RFC-3339 format when the project was last updated.
* `uuid` - Project UUID.
* `creation_time` - UTC date and time in RFC-3339 format when the project was created.
* `spec_version` - Version number of the latest spec.
* `spec_hash` - Hash of the spec. This will be returned from server.
* `name` - Project name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories
The categories attribute supports the following:

* `name` - the key name.
* `value` - value of the key.

### Reference
The `project_reference`, `owner_reference` attributes supports the following:

* `kind` - (Required) The kind name (Default value: `project`).
* `name` - (Optional) the name.
* `uuid` - (Required) the UUID.

See detailed information in [Nutanix Project](https://www.nutanix.dev/reference/prism_central/v3/api/projects/postprojects/).
