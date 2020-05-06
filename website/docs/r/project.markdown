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
* `account_reference_list` - (Optional/Computed)
* `account_reference_list.#.kind` - (Optional) The efault value is `account`
* `account_reference_list.#.uuid` - (Required)
* `account_reference_list.#.name` - (Optional/Computed)

### Enviroment Reference List
* `environment_reference_list` - (Optional/Computed)
* `environment_reference_list.#.kind` - (Optional) The efault value is `enviroment`
* `environment_reference_list.#.uuid` - (Required)
* `environment_reference_list.#.name` - (Optional/Computed)

### Default Subnet Reference Map
* `default_subnet_reference` - (Required)
* `default_subnet_reference.kind` - (Optional) The efault value is `subnet`
* `default_subnet_reference.uuid` - (Required)
* `default_subnet_reference.name` - (Optional/Computed)

### user_reference_list
* `user_reference_list` - (Optional/Computed)
* `user_reference_list.#.kind` - (Optional) The efault value is `user`
* `user_reference_list.#.uuid` - (Required)
* `user_reference_list.#.name` - (Optional/Computed)

### External User Group Reference List
* `external_user_group_reference_list` - (Optional/Computed)
* `external_user_group_reference_list.#.kind` - (Optional) The efault value is `user_group`
* `external_user_group_reference_list.#.uuid` - (Required)
* `external_user_group_reference_list.#.name` - (Optional/Computed)

### Subnet Reference List
* `subnet_reference_list` - (Optional/Computed)
* `subnet_reference_list.#.kind` - (Optional) The efault value is `subnet`
* `subnet_reference_list.#.uuid` - (Required)
* `subnet_reference_list.#.name` - (Optional/Computed)

### External Network List
* `subnet_reference_list` - (Optional/Computed)
* `subnet_reference_list.#.uuid` - (Required)
* `subnet_reference_list.#.name` - (Optional/Computed)


## Attributes Reference
The following attributes are exported:

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
