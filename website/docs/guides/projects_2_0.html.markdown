---
layout: "nutanix"
page_title: "Terraform: Projects 2.0"
sidebar_current: "docs-nutanix-guides-projects-2-0"
description: |-
  Introduction to Projects 2.0, entity sharing, and how to migrate existing Terraform configurations after a Prism Central upgrade to 7.6.
---

# Terraform: Projects 2.0

This guide describes Projects 2.0 in the Nutanix Terraform provider and how to migrate existing configurations after Prism Central is upgraded to 7.6.

## Introduction to Projects 2.0

As part of Projects 2.0, the project itself is a thin entity. It holds only a name, a description, and an ID. Entities across IAM, Networking, and Flow Management can now be associated with a project.

Create and manage that thin project with [`nutanix_project_v2`](/docs/providers/nutanix/r/project_v2.html).

## Entity Association and Sharing

### Sharing and unsharing with a specific project

The following entities can be shared or unshared with a specific project through `shared_with_projects`:

- Prism: Categories (`nutanix_category_v2`)
- IAM: SAML Identity Provider (`nutanix_saml_idps_v2`)
- IAM: Directory Service (`nutanix_directory_services_v2`)
- Networking: Subnet (`nutanix_subnet_v2`)
- Networking: VPC (`nutanix_vpc_v2`)

`shared_with_projects` is a list of project external identifiers.

### Global sharing

IAM SAML Identity Providers and Directory Services can be shared with all projects, including projects created later, by setting `share_with_all_projects = true`.

### User management

Users associated with a project are managed with Role Membership APIs. Use [`nutanix_role_membership_v2`](/docs/providers/nutanix/r/role_membership_v2.html) and [`nutanix_role_memberships_v2`](/docs/providers/nutanix/d/role_memberships_v2.html).

### Cluster management

Clusters associated with a project are managed with Resource Group APIs. Use [`nutanix_resource_group_v2`](/docs/providers/nutanix/r/resource_group_v2.html) and [`nutanix_resource_groups_v2`](/docs/providers/nutanix/d/resource_groups_v2.html).

## New Resource Support and Entity Mapping

The following resources support Projects 2.0. Project association is expressed with `project_ext_id` unless noted otherwise.

### IAM and identity

| Entity | Terraform resource |
| --- | --- |
| Directory Service | `nutanix_directory_services_v2` |
| SAML Identity Provider | `nutanix_saml_idps_v2` |
| Authorization Policies | `nutanix_authorization_policy_v2` |
| Role Memberships (new) | `nutanix_role_membership_v2` |
| Roles | `nutanix_roles_v2` |

### Prism and MultiDomain

| Entity | Terraform resource |
| --- | --- |
| Categories | `nutanix_category_v2` |
| Project (new) | `nutanix_project_v2` |
| Resource Groups (new) | `nutanix_resource_group_v2` |

### VMM

| Entity | Terraform resource |
| --- | --- |
| Templates | `nutanix_templates_v2` |
| Image | `nutanix_images_v2` |
| OVA | `nutanix_ova_v2` |

### Networking

| Entity | Terraform resource |
| --- | --- |
| Floating IPs | `nutanix_floating_ip_v2` |
| Routes | `nutanix_routes_v2` |
| Routing Policies | `nutanix_pbrs_v2` |
| Subnets | `nutanix_subnet_v2` |
| Virtual Switch | `nutanix_virtual_switch_v2` |
| VPC | `nutanix_vpc_v2` |
| Network Functions | `nutanix_network_function_v2` |

### Flow Management

| Entity | Terraform resource |
| --- | --- |
| Address Groups | `nutanix_address_group_v2` |
| Service Groups | `nutanix_service_groups_v2` |
| Network Security Policy | `nutanix_network_security_policy_v2` |
| Entity Groups | `nutanix_entity_group_v2` |

### Data Protection and Volumes

| Entity | Terraform resource |
| --- | --- |
| Recovery Points | `nutanix_recovery_points_v2` |
| Protection Policies | `nutanix_protection_policy_v2` |
| Volume Groups | `nutanix_volume_group_v2` |

## Migration Strategy for Terraform Users

This section describes how to move existing infrastructure to a provider release that includes Projects 2.0.

In Projects 1.0, no entity other than the VM supported project association.

### Example scenario

A typical Projects 1.0 setup:

1. Create an Active Directory.
2. Create a subnet.
3. Create an image.
4. Create a project, and add one user (Project Admin) and one user group (Developer).
5. Create a category.
6. Create a VPC. The VPC is not consumed by another resource in this example.
7. Create a VM, and associate it with the category and the project.

```hcl
terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.2" # PC 7.5
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.4"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_roles_v2" "fetch_roles" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]

  project_admin_role_name = "Project Admin"
  developer_role_name     = "Developer"
  user_name               = "<user_name>"
  user_group_name         = "<user_group_name>"

  project_admin_role_id = [
    for role in data.nutanix_roles_v2.fetch_roles.roles :
    role.ext_id if role.display_name == local.project_admin_role_name
  ][0]

  developer_role_id = [
    for role in data.nutanix_roles_v2.fetch_roles.roles :
    role.ext_id if role.display_name == local.developer_role_name
  ][0]

  base_url    = "https://${var.nutanix_endpoint}:${var.nutanix_port}"
  auth_header = "Basic ${base64encode("${var.nutanix_username}:${var.nutanix_password}")}"
}

data "nutanix_users_v2" "fetch_users" {
  filter = "username eq '${local.user_name}'"
}

data "nutanix_user_groups_v2" "fetch_user_groups" {
  filter = "name eq '${local.user_group_name}'"
}

locals {
  user_uuid       = data.nutanix_users_v2.fetch_users.users[0].ext_id
  user_group_uuid = data.nutanix_user_groups_v2.fetch_user_groups.user_groups[0].ext_id
}

resource "nutanix_category_v2" "runner_category" {
  key         = "devops"
  value       = "githubRunners"
  description = "Category to associate with Runner VM"
}

resource "nutanix_images_v2" "centos_image" {
  name = "CentOs"
  type = "DISK_IMAGE"
  source {
    url_source {
      url = "<image_url>"
    }
  }
}

resource "nutanix_directory_services_v2" "active_directory" {
  name           = "activedirectorytest"
  url            = "ldap://<directory_host>:389"
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = "<domain_name>"

  service_account {
    username = "<service_account_username>"
    password = "<service_account_password>"
  }

  lifecycle {
    ignore_changes = [
      service_account.0.password,
    ]
  }
}

data "http" "list_accounts" {
  url    = "${local.base_url}/api/nutanix/v3/accounts/list"
  method = "POST"

  request_headers = {
    Content-Type  = "application/json"
    Accept        = "application/json"
    Authorization = local.auth_header
  }

  request_body = jsonencode({
    kind   = "account"
    filter = "name==NTNX_LOCAL_AZ"
  })

  insecure = true
}

locals {
  accounts_result = jsondecode(data.http.list_accounts.response_body)
  account_uuid    = local.accounts_result.entities[0].metadata.uuid
}

resource "nutanix_subnet_v2" "vlan_112" {
  name              = "vlan-112"
  description       = "Subnet Created on PC 7.5"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = 122
}

resource "nutanix_subnet_v2" "vlan_113" {
  name              = "vlan-test-2"
  description       = "subnet VLAN 113 managed by Terraform with IP pool"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = 122
  is_external       = true

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "<subnet_ip>"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "<gateway_ip>"
      }
      pool_list {
        start_ip {
          value = "<pool_start_ip>"
        }
        end_ip {
          value = "<pool_end_ip>"
        }
      }
    }
  }
}

resource "nutanix_vpc_v2" "vpc" {
  name        = "vpc-example"
  description = "VPC for example"

  external_subnets {
    subnet_reference = nutanix_subnet_v2.vlan_113.id
  }
}

resource "nutanix_project" "runner_project" {
  name                 = "Devops Project"
  description          = "Devops Project"
  use_project_internal = true
  api_version          = "3.1"

  cluster_reference_list {
    uuid = local.cluster_ext_id
  }

  account_reference_list {
    uuid = local.account_uuid
  }

  subnet_reference_list {
    uuid = nutanix_subnet_v2.vlan_112.id
  }

  default_subnet_reference {
    uuid = nutanix_subnet_v2.vlan_112.id
  }

  user_reference_list {
    name = local.user_name
    kind = "user"
    uuid = local.user_uuid
  }

  external_user_group_reference_list {
    name = local.user_group_name
    kind = "user_group"
    uuid = local.user_group_uuid
  }

  acp {
    role_reference {
      kind = "role"
      uuid = local.project_admin_role_id
      name = local.project_admin_role_name
    }
    user_reference_list {
      name = local.user_name
      kind = "user"
      uuid = local.user_uuid
    }
  }

  acp {
    role_reference {
      kind = "role"
      uuid = local.developer_role_id
      name = local.developer_role_name
    }
    user_group_reference_list {
      name = local.user_group_name
      kind = "user_group"
      uuid = local.user_group_uuid
    }
  }

  cluster_uuid = local.cluster_ext_id
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "Github runner VM"
  num_sockets          = 2
  num_cores_per_socket = 1
  memory_size_bytes    = 3 * 1024 * 1024 * 1024

  cluster {
    ext_id = local.cluster_ext_id
  }

  project {
    ext_id = nutanix_project.runner_project.id
  }

  categories {
    ext_id = nutanix_category_v2.runner_category.id
  }

  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        data_source {
          reference {
            image_reference {
              image_ext_id = nutanix_images_v2.centos_image.id
            }
          }
        }
        disk_size_bytes = 20 * 1024 * 1024 * 1024
      }
    }
  }

  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = nutanix_subnet_v2.vlan_112.id
        }
        vlan_mode                 = "ACCESS"
        should_allow_unknown_macs = false
      }
    }
  }

  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }

  power_state = "ON"
}
```

## Brownfield Migration Steps (PC Upgrade to 7.6)

When Prism Central is upgraded to 7.6, the backend migration runs as follows:

- Existing entities are associated with the default project (the system-defined project).
- The existing category is associated with the default project and shared with all projects.
- Entities such as subnets, VPCs, and directory services that were associated with the project are updated with `shared_with_projects` set to that project's UUID.
- Users associated with the project are migrated to role membership entities. One user on the project becomes one role membership.
- Clusters associated with the project are migrated to resource group entities. One cluster on the project becomes one resource group.

## Terraform Migration Flow

After the Prism Central upgrade, upgrade the Terraform provider to a release that includes Projects 2.0.

### Step 1: Refresh state

Run `terraform refresh` so entity state is updated with the new project fields, including `project_ext_id` and `shared_with_projects`.

```bash
terraform refresh
```

#### Subnet state before the PC upgrade

`project_ext_id` and `shared_with_projects` are absent. `metadata.project_reference_id` is empty.

```hcl
# nutanix_subnet_v2.vlan_112
name                       = "vlan-112"
description                = "Subnet Created on PC 7.5"
id                         = "<subnet_uuid>"
ext_id                     = "<subnet_uuid>"
cluster_reference          = "<cluster_uuid>"
subnet_type                = "VLAN"
network_id                 = 122
virtual_switch_reference   = "<virtual_switch_uuid>"
metadata {
  project_name         = ""
  project_reference_id = ""
}
```

#### Subnet state after the PC upgrade and provider upgrade

After migration, the subnet is associated with the default project and shared with the existing project.

```hcl
# nutanix_subnet_v2.vlan_112
name                       = "vlan-112"
description                = "Subnet Created on PC 7.5"
id                         = "<subnet_uuid>"
ext_id                     = "<subnet_uuid>"
cluster_reference          = "<cluster_uuid>"
subnet_type                = "VLAN"
network_id                 = 122
virtual_switch_reference   = "<virtual_switch_uuid>"
project_ext_id             = "00000000-0000-0000-0000-000000000000"
shared_with_projects       = ["<devops_project_uuid>"]
metadata {
  project_name         = null
  project_reference_id = null
}
```

#### Image state before the PC upgrade

```hcl
# nutanix_images_v2.centos_image
name                     = "CentOs"
id                       = "<image_uuid>"
ext_id                   = "<image_uuid>"
type                     = "DISK_IMAGE"
cluster_location_ext_ids = ["<cluster_uuid>"]
owner_ext_id             = "00000000-0000-0000-0000-000000000000"
```

#### Image state after the PC upgrade, provider upgrade, and refresh

After migration, the image is associated with the default project and shared with all projects.

```hcl
# nutanix_images_v2.centos_image
name                     = "CentOs"
id                       = "<image_uuid>"
ext_id                   = "<image_uuid>"
type                     = "DISK_IMAGE"
cluster_location_ext_ids = ["<cluster_uuid>"]
owner_ext_id             = "00000000-0000-0000-0000-000000000000"
project_ext_id           = "00000000-0000-0000-0000-000000000000"
share_with_all_projects  = true
```

~> **Note:** Later `terraform plan` runs do not propose changes for these project association fields. `project_ext_id` and `shared_with_projects` are optional and computed. Terraform reads them from the API, and it manages them only when you set them in configuration.

### Step 2: Migrate the project resource

Import the existing project into `nutanix_project_v2`, then remove the v3 `nutanix_project` resource from state.

Declare an empty resource block:

```hcl
resource "nutanix_project_v2" "devops_project" {}
```

Import the project by its UUID:

```bash
terraform import nutanix_project_v2.devops_project <project_uuid>
```

After the import succeeds, remove the v3 project from state:

```bash
terraform state rm nutanix_project.runner_project
```

Remove the `nutanix_project` block from configuration and manage the project with `nutanix_project_v2`.

### Step 3: Manage users associated with the project (role memberships)

List the role memberships created by the backend migration, then import each one you want Terraform to manage.

```hcl
data "nutanix_role_memberships_v2" "fetch_role_memberships" {
  filter = "projectExtId eq '${nutanix_project_v2.devops_project.ext_id}'"
}
```

Each membership UUID is available at `data.nutanix_role_memberships_v2.fetch_role_memberships.role_memberships[*].ext_id`.

Declare an empty resource block for each membership:

```hcl
resource "nutanix_role_membership_v2" "import_role_membership" {}
```

Import it:

```bash
terraform import nutanix_role_membership_v2.import_role_membership <role_membership_uuid>
```

Repeat the import for every role membership you want to manage, including user groups.

### Step 4: Manage clusters associated with the project (resource groups)

List the resource groups created for the project, then import each one.

```hcl
data "nutanix_resource_groups_v2" "fetch_resource_groups" {
  filter = "projectExtId eq '${nutanix_project_v2.devops_project.ext_id}'"
}
```

Each resource group UUID is available at `data.nutanix_resource_groups_v2.fetch_resource_groups.resource_groups[*].ext_id`.

Declare an empty resource block:

```hcl
resource "nutanix_resource_group_v2" "import_rg" {}
```

Import it:

```bash
terraform import nutanix_resource_group_v2.import_rg <resource_group_uuid>
```

Importing the resource group lets you manage the clusters and storage containers associated with the project.

### Step 5: Update sharing configuration

To change which projects a subnet, VPC, category, SAML identity provider, or directory service is shared with, set `shared_with_projects` in configuration. The current value is already in state after the refresh in Step 1.

To share an Active Directory or a SAML identity provider with every project, including projects created later, set `share_with_all_projects = true`.

## Configuration After Migration

After the provider upgrade, the project import, and the role membership and resource group imports, configuration looks like the following. The v3 account lookup is no longer required, because `nutanix_project_v2` does not take account, subnet, user, or cluster references. Those associations live on role memberships, resource groups, and `shared_with_projects`.

```hcl
terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0-dev" # PC 7.6
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_roles_v2" "fetch_roles" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]

  project_admin_role_name = "Project Admin"
  developer_role_name     = "Developer"
  user_name               = "<user_name>"
  user_group_name         = "<user_group_name>"

  project_admin_role_id = [
    for role in data.nutanix_roles_v2.fetch_roles.roles :
    role.ext_id if role.display_name == local.project_admin_role_name
  ][0]

  developer_role_id = [
    for role in data.nutanix_roles_v2.fetch_roles.roles :
    role.ext_id if role.display_name == local.developer_role_name
  ][0]
}

data "nutanix_users_v2" "fetch_users" {
  filter = "username eq '${local.user_name}'"
}

data "nutanix_user_groups_v2" "fetch_user_groups" {
  filter = "name eq '${local.user_group_name}'"
}

locals {
  user_uuid       = data.nutanix_users_v2.fetch_users.users[0].ext_id
  user_group_uuid = data.nutanix_user_groups_v2.fetch_user_groups.user_groups[0].ext_id
}

resource "nutanix_category_v2" "runner_category" {
  key         = "devops"
  value       = "githubRunners"
  description = "Category to associate with Runner VM"
}

resource "nutanix_images_v2" "centos_image" {
  name = "CentOs"
  type = "DISK_IMAGE"
  source {
    url_source {
      url = "<image_url>"
    }
  }
}

resource "nutanix_directory_services_v2" "active_directory" {
  name                    = "activedirectorytest"
  url                     = "ldap://<directory_host>:389"
  directory_type          = "ACTIVE_DIRECTORY"
  domain_name             = "<domain_name>"
  share_with_all_projects = true

  service_account {
    username = "<service_account_username>"
    password = "<service_account_password>"
  }

  lifecycle {
    ignore_changes = [
      service_account.0.password,
    ]
  }
}

resource "nutanix_project_v2" "runner_project" {
  name        = "Devops Project"
  description = "Devops Project"
  project_id  = "devops-project"
}

resource "nutanix_project_v2" "runner_project_extra" {
  name        = "Devops Project 2"
  description = "Devops Project 2"
  project_id  = "devops-project-2"
}

resource "nutanix_subnet_v2" "vlan_112" {
  name                 = "vlan-112"
  description          = "Subnet Created on PC 7.5"
  cluster_reference    = local.cluster_ext_id
  subnet_type          = "VLAN"
  network_id           = 122
  shared_with_projects = [
    nutanix_project_v2.runner_project.ext_id,
    nutanix_project_v2.runner_project_extra.ext_id,
  ]
}

resource "nutanix_subnet_v2" "vlan_113" {
  name              = "vlan-test-2"
  description       = "subnet VLAN 113 managed by Terraform with IP pool"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = 122
  is_external       = true

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "<subnet_ip>"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "<gateway_ip>"
      }
      pool_list {
        start_ip {
          value = "<pool_start_ip>"
        }
        end_ip {
          value = "<pool_end_ip>"
        }
      }
    }
  }
}

resource "nutanix_vpc_v2" "vpc" {
  name        = "vpc-example"
  description = "VPC for example"

  external_subnets {
    subnet_reference = nutanix_subnet_v2.vlan_113.id
  }
}

data "nutanix_resource_groups_v2" "fetch_resource_groups" {
  filter = "projectExtId eq '${nutanix_project_v2.runner_project.ext_id}'"
}

data "nutanix_role_memberships_v2" "fetch_role_memberships" {
  filter = "projectExtId eq '${nutanix_project_v2.runner_project.ext_id}'"
}

resource "nutanix_resource_group_v2" "runner_project_resource_group" {
  name           = "devops-project-infra"
  project_ext_id = nutanix_project_v2.runner_project.ext_id

  placement_targets {
    cluster_ext_id = "<cluster_uuid>"
    storage_containers {
      ext_id = "<storage_container_ext_id>"
    }
  }
}

resource "nutanix_role_membership_v2" "project_admin_user" {
  identity_ext_id     = local.user_uuid
  identity_type       = "USER"
  idp_ext_id          = nutanix_directory_services_v2.active_directory.id
  project_ext_id      = nutanix_project_v2.runner_project.ext_id
  role_ext_id         = local.project_admin_role_id
  scope_template_name = "ProjectsScopeTemplate"

  scope_template_name_values {
    name  = "projectExtId"
    value = nutanix_project_v2.runner_project.ext_id
  }
}

resource "nutanix_role_membership_v2" "developer_user_group" {
  identity_ext_id     = local.user_group_uuid
  identity_type       = "GROUP"
  idp_ext_id          = nutanix_directory_services_v2.active_directory.id
  project_ext_id      = nutanix_project_v2.runner_project.ext_id
  role_ext_id         = local.developer_role_id
  scope_template_name = "ProjectsScopeTemplate"

  scope_template_name_values {
    name  = "projectExtId"
    value = nutanix_project_v2.runner_project.ext_id
  }
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "Github runner VM"
  num_sockets          = 2
  num_cores_per_socket = 1
  memory_size_bytes    = 3 * 1024 * 1024 * 1024

  cluster {
    ext_id = local.cluster_ext_id
  }

  project {
    ext_id = nutanix_project_v2.runner_project.id
  }

  categories {
    ext_id = nutanix_category_v2.runner_category.id
  }

  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        data_source {
          reference {
            image_reference {
              image_ext_id = nutanix_images_v2.centos_image.id
            }
          }
        }
        disk_size_bytes = 20 * 1024 * 1024 * 1024
      }
    }
  }

  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = nutanix_subnet_v2.vlan_112.id
        }
        vlan_mode                 = "ACCESS"
        should_allow_unknown_macs = false
      }
    }
  }

  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }

  power_state = "ON"
}
```

## Continuing to Use the Legacy Project Resource

Configurations that keep using [`nutanix_project`](/docs/providers/nutanix/r/project.html) still apply. `directory_reference_list` and `identity_providers_reference_list` are populated in state by `terraform refresh` or `terraform plan`.

Add those attributes to configuration when you create or update a project that has user associations. On Prism Central 7.6, the v3 API rejects create and update calls that omit directory service and identity provider references while project-user associations are present. Set `enable_directory_and_identity_provider_shortlist = true` together with the reference lists.

Change the configuration only when you want Terraform to update the directory service or identity provider references.

```hcl
resource "nutanix_project" "runner_project" {
  name                 = "Devops Project"
  description          = "Devops Project"
  use_project_internal = true
  api_version          = "3.1"

  cluster_reference_list {
    uuid = "<cluster_uuid>"
  }

  account_reference_list {
    uuid = "<account_uuid>"
  }

  subnet_reference_list {
    uuid = "<subnet_uuid>"
  }

  default_subnet_reference {
    uuid = "<subnet_uuid>"
  }

  enable_directory_and_identity_provider_shortlist = true

  directory_reference_list {
    uuid = "<directory_service_uuid>"
  }

  # Add this block when identity provider users are added to the project.
  # identity_providers_reference_list {
  #   uuid = "<identity_provider_uuid>"
  # }

  user_reference_list {
    name = "<user_name>"
    kind = "user"
    uuid = "<user_uuid>"
  }

  external_user_group_reference_list {
    name = "<user_group_name>"
    kind = "user_group"
    uuid = "<user_group_uuid>"
  }

  acp {
    role_reference {
      kind = "role"
      uuid = "<project_admin_role_uuid>"
      name = "Project Admin"
    }
    user_reference_list {
      name = "<user_name>"
      kind = "user"
      uuid = "<user_uuid>"
    }
  }

  acp {
    role_reference {
      kind = "role"
      uuid = "<developer_role_uuid>"
      name = "Developer"
    }
    user_group_reference_list {
      name = "<user_group_name>"
      kind = "user_group"
      uuid = "<user_group_uuid>"
    }
  }

  cluster_uuid = "<cluster_uuid>"
}
```
