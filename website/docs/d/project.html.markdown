---
layout: "nutanix"
page_title: "NUTANIX: nutanix_project"
sidebar_current: "docs-nutanix-datasource-project"
description: |-
  Describe a Nutanix Project and its values (if it has them).
---

# nutanix_project

Describe a Nutanix Project and its values (if it has them).

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

data "nutanix_project" "test" {
    project_id = nutanix_project.project_test.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id`: - (Required) The `id` of the project.

## Attributes Reference

The following attributes are exported:

* `name` The name for the project.
* `description` A description for project.

### Resource Domain
* `resource_domain` The status for a resource domain (limits and values)
* `resource_domain.resources` Array of the utilization/limit for resource types
* `resource_domain.resources.#.limit` The resource consumption limit (unspecified is unlimited)
* `resource_domain.resources.#.resource_type` The type of resource (for example storage, CPUs)
* `resource_domain.resources.#.units` - The units of the resource type
* `resource_domain.resources.#.value` - The amount of resource consumed

### Account Reference List
* `account_reference_list`
* `account_reference_list.#.kind`
* `account_reference_list.#.uuid`
* `account_reference_list.#.name`

### Environment Reference List
* `environment_reference_list`
* `environment_reference_list.#.kind`
* `environment_reference_list.#.uuid`
* `environment_reference_list.#.name`

### Default Subnet Reference Map
* `default_subnet_reference`
* `default_subnet_reference.kind`
* `default_subnet_reference.uuid`
* `default_subnet_reference.name`

### user_reference_list
* `user_reference_list`
* `user_reference_list.#.kind`
* `user_reference_list.#.uuid`
* `user_reference_list.#.name`

### External User Group Reference List
* `external_user_group_reference_list`
* `external_user_group_reference_list.#.kind`
* `external_user_group_reference_list.#.uuid`
* `external_user_group_reference_list.#.name`

### Subnet Reference List
* `subnet_reference_list`
* `subnet_reference_list.#.kind`
* `subnet_reference_list.#.uuid`
* `subnet_reference_list.#.name`

### External Network List
* `subnet_reference_list`
* `subnet_reference_list.#.uuid`
* `subnet_reference_list.#.name`

### Resource Domain
* `resource_domain.resources.#.units` - The units of the resource type
* `resource_domain.resources.#.value` - The amount of resource consumed

### Metadata
The metadata attribute exports the following:

* `last_update_time` - UTC date and time in RFC-3339 format when vm was last updated.
* `uuid` - vm UUID.
* `creation_time` - UTC date and time in RFC-3339 format when vm was created.
* `spec_version` - Version number of the latest spec.
* `spec_hash` - Hash of the spec. This will be returned from server.
* `name` - vm name.
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

See detailed information in [Nutanix Project](https://www.nutanix.dev/reference/prism_central/v3/api/projects/getprojectsuuid/).
