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

# set use_project_internal flag to create project with acps

resource "nutanix_project" "project_test" {
  name        = "my-project"
  description = "This is my project"

  cluster_uuid = "<YOUR_CLUSTER_ID>"

  use_project_internal = true

  default_subnet_reference {
    uuid = nutanix_subnet.subnet.metadata.uuid
  }

  user_reference_list{
    name= "{{user_name}}"
    kind= "user"
    uuid= "{{user_uuid}}"
    }
    subnet_reference_list{
      uuid=resource.nutanix_subnet.sub.id
  }
  acp{
    # acp name consists name_uuid string, it should be different for each acp. 
    name="{{acp_name}}"
    role_reference{
      kind= "role"
      uuid= "{{role_uuid}}"
      name="Developer"
    }
    user_reference_list{
      name= "{{user_name}}"
      kind= "user"
      uuid= "{{user_uuid}}"
    }
    description= "{{description}}"
  }
  api_version = "3.1"
}

## Create a project with user which not added in the PC

resource "nutanix_project" "project_test" {
  name        = "my-project"
  description = "This is my project"

  cluster_uuid = "<YOUR_CLUSTER_ID>"

  use_project_internal = true

  default_subnet_reference {
    uuid = nutanix_subnet.subnet.metadata.uuid
  }

  user_reference_list{
    name= "{{user_name}}"
    kind= "user"
    uuid= "{{user_uuid}}"
    }
    subnet_reference_list{
      uuid=resource.nutanix_subnet.sub.id
  }
  acp{
    # acp name consists name_uuid string, it should be different for each acp. 
    name="{{acp_name}}"
    role_reference{
      kind= "role"
      uuid= "{{role_uuid}}"
      name="Developer"
    }
    user_reference_list{
      name= "{{user_name}}"
      kind= "user"
      uuid= "{{user_uuid}}"
    }
    description= "{{description}}"
  }
  user_list{
    metadata={
      kind="user"
      uuid= "{{ UUID of the USER }}"
    }
    directory_service_user{
      user_principal_name= "{{ Name of user }}"
      directory_service_reference{
        uuid="{{ DIRECTORY SERVICE UUID }}"
        kind="directory_service"
      }
    }
  }
  api_version = "3.1"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the project.
* `description` - (Required) A description for project.

* `use_project_internal` - (Optional) flag to use project internal for user role mapping
* `cluster_uuid` - (Optional) The UUID of cluster. (Required when using project_internal flag).
* `enable_collab` - (Optional) flag to allow collaboration of projects. (Use with project_internal flag)

### Resource Domain
* `resource_domain` - (Optional) The status for a resource domain (limits and values)
* `resource_domain.resources` - (Required) Array of the utilization/limit for resource types
* `resource_domain.resources.#.limit` - (Required) The resource consumption limit.
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

### Tunnel Reference List
* `tunnel_reference_list` - (Optional/Computed) List of tunnels associated with the project.
* `tunnel_reference_list.#.kind` - (Optional) The kind name. Default value is `tunnel`
* `tunnel_reference_list.#.uuid` - (Required) The UUID of a tunnel
* `tunnel_reference_list.#.name` - (Optional/Computed) The name of a tunnel.

### Cluster Reference List
* `cluster_reference_list` - (Optional/Computed) List of clusters associated with the project..
* `cluster_reference_list.#.kind` - (Optional) The kind name. Default value is `cluster`
* `cluster_reference_list.#.uuid` - (Required) The UUID of a cluster
* `cluster_reference_list.#.name` - (Optional/Computed) The name of a cluster.

### VPC Reference List
* `vpc_reference_list` - (Optional/Computed) List of VPCs associated with the project..
* `vpc_reference_list.#.kind` - (Optional) The kind name. Default value is `vpc`
* `vpc_reference_list.#.uuid` - (Required) The UUID of a vpc
* `vpc_reference_list.#.name` - (Optional/Computed) The name of a vpc.

### Default Environment Reference Map
* `default_environment_reference` - (Optional/Computed) Reference to a environment.
* `default_environment_reference.kind` - (Optional) The kind name. Default value is `environment`
* `default_environment_reference.uuid` - (Required) The UUID of a environment
* `default_environment_reference.name` - (Optional/Computed) The name of a environment.


### ACP
* `acp` - (Optional) The list of ACPs to be attached to the users belonging to a project. It is mandate to provide cluster_uuid while using ACP. It helps to get the context list based on user role. 
* `acp.#.name` - (Required) Name of the Access Control Policy.  
* `acp.#.description` -  The description of the association of a role to a user in a given context.

* `acp.#.user_reference_list` - The User(s) being assigned a given role.
* `acp.#.user_reference_list.#.kind` - The kind name. Default value is `user`
* `acp.#.user_reference_list.#.uuid` - (Required) The UUID of a user
* `acp.#.user_reference_list.#.name` - (Optional/Computed) The name of a user.

* `acp.#.user_group_reference_list` - The User group(s) being assigned a given role
* `acp.#.user_group_reference_list.#.kind` - The kind name. Default value is `user_group`
* `acp.#.user_group_reference_list.#.uuid` - (Required) The UUID of a user group
* `acp.#.user_group_reference_list.#.name` - (Optional/Computed) The name of a user group.

* `acp.#.role_reference` - Reference to a role.
* `acp.#.role_reference.kind` - The kind name. Default value is `role`
* `acp.#.role_reference.uuid` - (Required) The UUID of a role
* `acp.#.role_reference.name` - (Optional/Computed) The name of a role.

### User List
* `user_list` - (Optional) The list of user specification to be associated with the project. It is only required when user is not added in the PC. 
* `user_list.#.directory_service_user` - (Optional) A Directory Service user.
* `user_list.#.directory_service_user.user_principal_name` - (Required) The UserPrincipalName of the user from the directory service. 
* `user_list.#.directory_service_user.directory_service_reference` - (Required) Reference to a directory_service . 
* `user_list.#.directory_service_user.directory_service_reference.uuid` - (Required) The uuid to a directory_service. 
* `user_list.#.directory_service_user.directory_service_reference.kind` - (Optional) The kind to a directory_service.

* `user_list.#.identity_provider_user` - (Optional) An Identity Provider user.
* `user_list.#.identity_provider_user.username` - (Required) The username from the identity provider. Name Id for SAML Identity Provider.
* `user_list.#.identity_provider_user.identity_provider_reference` - (Required) The reference to a identity_provider. 
* `user_list.#.identity_provider_user.identity_provider_reference.uuid` - (Required) The uuid to a identity_provider. 
* `user_list.#.identity_provider_user.identity_provider_reference.kind` - (Optional) The kind to a identity_provider.
* `user_list.#.metadata` - (Required) Metadata Reference for user
* `user_list.#.metadata.uuid` - (Required) UUID of the USER
* `user_list.#.metadata.Kind` - Kind of the USER. 


### User Group

* `user_group` - (Optional) The list of user group specification to be associated with the project. It is only Required when user group is not added in the PC. 
* `user_group.#.directory_service_user_group` - (Optional) A Directory Service user group.
* `user_group.#.directory_service_user_group.distinguished_name` -  (Required) The Distinguished name for the user group. 

* `user_group.#.saml_user_group` - (Optional) A SAML Service user group.
* `user_group.#.saml_user_group.idp_uuid` - (Required) The UUID of the Identity Provider that the group belongs to. 
* `user_group.#.saml_user_group.name` - (Required) The name of the SAML group which the IDP provides as attribute in SAML response. 

* `user_group.#.directory_service_ou` - (Optional) A Directory Service user group. 
* `user_group.#.directory_service_ou.distinguished_name` -  (Required) The Distinguished name for the user group. 
* `user_group.#.metadata` - (Required) Metadata Reference for user group
* `user_group.#.metadata.uuid` - (Required) UUID of the USER Group
* `user_group.#.metadata.Kind` - Kind of the USER Group. 


## Attributes Reference
The following attributes are exported:

### Resource Domain
* `resource_domain.resources.#.units` - The units of the resource type
* `resource_domain.resources.#.value` - The amount of resource consumed

### ACP
ACPs will be exported if use_project_internal flag is set.
* `name` - Name of ACP
* `description` - Description of ACP
* `user_reference_list` - List of Reference of users.
* `user_group_reference_list` - List of Reference of users groups.
* `role_reference` - Reference to role.
* `context_filter_list` - The list of context filters. These are OR filters. The scope-expression-list defines the context, and the filter works in conjunction with the entity-expression-list.

The context_list attribute supports the following:

* `scope_filter_expression_list`: - (Optional) Filter the scope of an Access Control Policy.
* `entity_filter_expression_list` - (Required) A list of Entity filter expressions.

### Scope Filter Expression List

The scope_filter_expression_list attribute supports the following.

* `left_hand_side`: - (Optional)  The LHS of the filter expression - the scope type.
* `operator`: - (Required) The operator of the filter expression.
* `right_hand_side`: - (Required) The right hand side (RHS) of an scope expression.


### Entity Filter Expression List

The scope_filter_expression_list attribute supports the following.

* `left_hand_side_entity_type`: - (Optional)  The LHS of the filter expression - the entity type.
* `operator`: - (Required) The operator in the filter expression.
* `right_hand_side`: - (Required) The right hand side (RHS) of an scope expression.

### Right Hand Side

The right_hand_side attribute supports the following.

* `collection`: - (Optional)  A representative term for supported groupings of entities. ALL = All the entities of a given kind.
* `categories`: - (Optional) The category values represented as a dictionary of key -> list of values.
* `uuid_list`: - (Optional) The explicit list of UUIDs for the given kind.


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

Note: Few attributes which are added to support ACPs for Project are dependent on PC version. Features such as VPC, Cluster Reference requires pc2022.4 while Tunnel Reference requires pc2022.6 . 

See detailed information in [Nutanix Project](https://www.nutanix.dev/api_references/prism-central-v3/#/8411486d42e4a-create-a-new-project).
