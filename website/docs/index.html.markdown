---
layout: "nutanix"
page_title: "Provider: Nutanix"
sidebar_current: "docs-nutanix-index"
description: |-
  The provider is used to interact with the many resources supported by Nutanix. The provider needs to be configured with the proper credentials before it can be used.
---

# Nutanix Provider

The provider is used to interact with the many resources and data sources supported by Nutanix, using Prism Central as the provider endpoint.

Use the navigation on the left to read about the available resources and data sources this provider can use.


## Introducing Nutanix Terraform Provider Version v2.4.0

We're excited to announce the release of Nutanix Terraform Provider Version 2.4.0!

### What's New in v2.4.0

- **New Resource Support**
  - **Key Management Server (Security)**: Manage and configure external Key Management Servers for securing workloads.
  - **Security Technical Implementation Guide controls details (Security)**: View compliance with technical security controls.
  - **SSL Certification (Cluster Management)**: Add and manage SSL certificates for secure cluster communications.
  - **Cluster Profile (Cluster Management)**: Define and manage cluster configuration profiles for consistent deployments.
  - **Associate/Disassociate Cluster from Cluster Profile (Cluster Management)**: Link or unlink clusters to profiles for streamlined management.
  - **Associate/Disassociate Categories to Cluster (Cluster Management)**: Assign or remove custom categories to clusters.
  - **Storage Policies (Data Policies)**: Create and manage storage policy rules to optimize resource allocation.

- **Enhancements:**
  - Add Support for Package-Specific Acceptance Tests via /ok-to-test -p Command [#1014](https://github.com/nutanix/terraform-provider-nutanix/issues/1014)
  - Centralize task entity type and completion detail constants for reliable UUID extraction [#1029](https://github.com/nutanix/terraform-provider-nutanix/issues/1029)

- **Fixed Bugs:**
   - Unable to list VPC using data "nutanix_vpcs_v2" "list_vpcs" [#1000](https://github.com/nutanix/terraform-provider-nutanix/issues/1000)
   - virtual_machine_v2: VM creation fails with multiple NICs ("invalid input arguments") [#994](https://github.com/nutanix/terraform-provider-nutanix/issues/994)
   - V3: Project: Revisit the Project Module resources [#962](https://github.com/nutanix/terraform-provider-nutanix/issues/962)
      - Projects: ACP: Order changes in API response lead to data inconsistency in state file. [#1042](https://github.com/nutanix/terraform-provider-nutanix/issues/1042)
      - Projects: ACP: Removing a ACP causing index shifting issues. [#1044](https://github.com/nutanix/terraform-provider-nutanix/issues/1044)
      - Project: ACP: Adding a new user or new user group to existing ACP is failed. [#1043](https://github.com/nutanix/terraform-provider-nutanix/issues/1043)
   - Bug Report: resource "nutanix_user_groups_v2" [#947](https://github.com/nutanix/terraform-provider-nutanix/issues/947)

- **Breaking Chnages:**
   - From PC version 7.5 onwards, the resource domain is not supported for Projects resources. As a result, Terraform support for this functionality has been removed starting with the 2.4.0 release. [#1049](https://github.com/nutanix/terraform-provider-nutanix/issues/1049)


~> **Important Notice:** Upcoming Deprecation of Legacy Nutanix Terraform Provider Resources. Starting with the Nutanix Terraform Provider release planned for Q4-CY2026, legacy resources which are based on v0.8,v1,v2 and v3 APIs will be deprecated and no longer supported. For more information, visit [Legacy API Deprecation Announcement](https://portal.nutanix.com/page/documents/eol/list?type=announcement) [Legacy API Deprecation - FAQs](https://portal.nutanix.com/page/documents/kbs/details?targetId=kA0VO0000005rgP0AQ). Nutanix strongly encourages you to migrate your scripts and applications to the latest v2 version of the Nutanix Terraform Provider resources, which are built on our v4 APIs/SDKs. By adopting the latest v2 version based on v4 APIs and SDKs, our users can leverage the enhanced capabilities and latest innovations from Nutanix. We understand that this transition may require some effort, and we are committed to supporting you throughout the process. Please refer to our documentation and support channels for guidance and assistance.

## Support

-> **Note:** We now have a brand new developer-centric Support Program designed for organizations that require a deeper level of developer support to manage their Nutanix environment and build applications quickly and efficiently. As part of this new Advanced API/SDK Support Program, you will get access to trusted technical advisors who specialize in developer tools including Nutanix Terraform Provider and receive support for your unique development needs and custom integration queries. Visit our Support Portal - [Premium Add-On Support Programs](https://www.nutanix.com/support-services/product-support/premium-support-programs) to learn more about this program.

Customers not taking advantage of the  Advanced API/SDK Support Program will continue to receive the support through our standard, community-supported model. This community model also provides support for contributions to the open-sourceNutanix Terraform Provider repository .Visit https://portal.nutanix.com/kb/13424   for more details. 

## Compatibility Matrix
| Terraform Version |  AOS Version | PC version  | Other software versions | Supported |
|  :--- |  :--- | :--- | :--- | :--- |
| 2.4.0 | 7.5 | pc7.5 or later | Self Service  v4.3.0 | yes |
| 2.3.4 | 7.3 | pc7.3 or later | Self Service  v4.2.0, v4.1.0 | yes |
| 2.3.3 | 7.3 | pc7.3 or later | Self Service  v4.2.0, v4.1.0 | yes |
| 2.3.2 | 7.3 | pc7.3 or later | Self Service  v4.2.0, v4.1.0 | yes |
| 2.3.1 | 7.3 | pc7.3 or later | Self Service  v4.2.0, v4.1.0 | yes |
| 2.3.0 | 7.3 | pc7.3 or later | Self Service  v4.2.0, v4.1.0 | yes |
| 2.2.3 | 7.0.1, 7.0 | pc2024.3, pc2024.3.1 or later | | yes |
| 2.2.2 (⚠️ Deprecated/Invalid) | 7.0.1, 7.0 | pc2024.3, pc2024.3.1 or later | | yes |
| 2.2.1 | 7.0.1, 7.0 | pc2024.3, pc2024.3.1 or later | | yes |
| 2.2.0 | | | Self Service  v4.1.0 | yes | 
| 2.1.1 | 7.0.1, 7.0 | pc2024.3, pc2024.3.1 or later | | yes |
| 2.1.0 | 7.0.1, 7.0 | pc2024.3, pc2024.3.1 or later | | yes |
| 2.0.0   |  7.0  | pc2024.3 or later  | ndb v2.7, nke v2.8, foundation v5.7 | Yes |
| 1.9.5 | | pc2023.1.0.2 | ndb v2.5.1.1, v2.5.1,  v2.5 |  Yes |
| 1.9.4 | | pc2023, pc2023.1.0.2, pc2023.1.0.1 |  | Yes |
| 1.9.3 | | pc2023.1.0.1 | | No |
| 1.9.2 | | pc2023.1.0.1 | | No |
| 1.9.1 | | pc2023.1.0.1 | ndb v2.5.1,  v2.5 | No |
| 1.9.0 | | pc2023.1.0.1, pc2022.9 | ndb v2.5.1, v2.5 | No |
| 1.8.0 | | pc2022.6 | ndb v2.5.1.1, v2.5.1 and v2.5 | No |
| 1.8.1 | | pc2022.6 | ndb v2.5.1.1, v2.5.1 and v2.5 | No |
| 1.7.0 | | pc2022.6, pc2022.4 and pc2022.1.0.2 | | No |
| 1.7.1 | | pc2022.6, pc2022.4.0.1 and pc2022.1.0.2 | | No |
| 1.6.1 | | pc2022.4 pc2022.1.0.2 and pc2021.9.0.4| | No |
| 1.5.0 | | pc2022.1.0.2 pc.2021.9.0.4 and pc.2021.8.0.1 | foundation v5.2, v5.1.1 , foundation central v1.3, v1.2 | No |
| 1.4.0 | | pc2022.1 pc.2021.9.0.4 and pc.2021.8.0.1 | | No |
| 1.3.0 | | pc.2021.9.0.4, pc.2021.8.0.1 and pc.2021.7 | | No |
| 1.2.0 | 5.18, 5.19 | pc2020.9 and pc2020.11| | No |

## Resources

| v1 Resources| v2 Resources |
|  :--- |  :--- |
| nutanix_subnet | nutanix_subnet_v2 |
| nutanix_vpc | nutanix_vpc_v2 |
| nutanix_floating_ip | nutanix_floating_ip_v2 |
| nutanix_pbr | nutanix_pbr_v2 |
| nutanix_static_routes | nutanix_routes_v2 |
| nutanix_address_group | nutanix_address_groups_v2 |
| nutanix_service_group | nutanix_service_groups_v2 |
| nutanix_network_security_rule | nutanix_network_security_policy_v2 |
| nutanix_role | nutanix_roles_v2 |
| nutanix_user | nutanix_users_v2 |
| nutanix_user_groups | nutanix_user_groups_v2 |
| nutanix_access_control_policy | nutanix_authorization_policy_v2 |
| - | nutanix_saml_identity_providers_v2 |
| - | nutanix_directory_services_v2 |
| nutanix_category_key | nutanix_category_v2 |
| nutanix_category_value | - |
| nutanix_image |nutanix_images_v2 |
| - | nutanix_image_placement_policy_v2 |
| nutanix_virtual_machine | nutanix_virtual_machine_v2 |
| - | nutanix_ova_v2 |
| - | nutanix_ova_vm_deploy_v2 |
| - | nutanix_ova_download_v2 |
| - | nutanix_vm_clone_v2 |
| - | nutanix_vm_cdrom_insert_eject_v2 |
| - | nutanix_vm_shutdown_action_v2 |
| - | nutanix_vm_gc_update_v2 |
| - | nutanix_vm_network_device_assign_ip_v2 |
| - | nutanix_vm_network_device_migrate_v2 |
| - | nutanix_template_v2 |
| - | nutanix_deploy_templates_v2 |
| - | nutanix_template_guest_os_actions_v2 |
| - | nutanix_ngt_installation_v2 |
| - | nutanix_ngt_upgrade_v2 |
| - | nutanix_ngt_insert_iso_v2 |
| - | nutanix_vm_revert_v2 |
| - | nutanix_recovery_points_v2 |
| - | nutanix_recovery_point_replicate_v2 |
| - | nutanix_recovery_point_restore_v2 |
| - | nutanix_volume_group_v2 |
| - | nutanix_volume_group_disk_v2 |
| - | nutanix_volume_group_iscsi_client_v2 |
| - | nutanix_volume_group_vm_v2 |
| - | nutanix_storage_containers_v2 |
| - | nutanix_cluster_v2 |
| - | nutanix_cluster_add_node_v2 |
| - | nutanix_pc_registration_v2 |
| - | nutanix_clusters_discover_unconfigured_nodes_v2 |
| - | nutanix_clusters_unconfigured_node_networks_v2 |
| nutanix_project | - |
| nutanix_protection_rule | - |
| nutanix_recovery_plan | - |
| nutanix_karbon_cluster | - |
| nutanix_karbon_private_registry | - |
| nutanix_foundation_image_nodes | - |
| nutanix_foundation_ipmi_config | - |
| nutanix_foundation_image | - |
| nutanix_foundation_central_image_cluster | - |
| nutanix_foundation_central_api_keys | - |
| nutanix_ndb_database | - |
| nutanix_ndb_sla | - |
| nutanix_ndb_database_restore | - |
| nutanix_ndb_log_catchups | - |
| nutanix_ndb_profile | - |
| nutanix_ndb_software_version_profile | - |
| nutanix_ndb_scale_database | - |
| nutanix_ndb_database_scale | - |
| nutanix_ndb_register_database | - |
| nutanix_ndb_database_snapshot | - |
| nutanix_ndb_clone | - |
| nutanix_ndb_authorize_dbserver | - |
| nutanix_ndb_linked_databases | - |
| nutanix_ndb_maintenance_window | - |
| nutanix_ndb_maintenance_task | - |
| nutanix_ndb_tms_cluster | - |
| nutanix_ndb_tag | - |
| nutanix_ndb_network | - |
| nutanix_ndb_dbserver_vm | - |
| nutanix_ndb_register_dbserver | - |
| nutanix_ndb_stretched_vlan | - |
| nutanix_ndb_clone_refresh | - |
| nutanix_ndb_cluster | - |
| - | nutanix_pc_deploy_v2 |
| - | nutanix_pc_backup_target_v2 |
| - | nutanix_pc_restore_source_v2 |
| - | nutanix_pc_restore_v2 |
| - | nutanix_pc_unregistration_v2 |
| - | nutanix_promote_protected_resource_v2 |
| - | nutanix_restore_protected_resource_v2 |
| - | nutanix_protection_policy_v2 |
| - | nutanix_lcm_perform_inventory_v2 |
| - | nutanix_lcm_prechecks_v2 |
| - | nutanix_lcm_upgrade_v2 |
| - | nutanix_lcm_config_v2 |
| nutanix_self_service_app_provision | - |
| nutanix_self_service_app_patch | - |
| nutanix_self_service_app_recovery_point | - |
| nutanix_self_service_app_custom_action | - |
| nutanix_self_service_app_restore | - |
| - | nutanix_user_key_v2 |
| - | nutanix_user_key_revoke_v2 |
| - | nutanix_object_store_v2 |
| - | nutanix_object_store_certificate_v2 |
| - | nutanix_password_change_request_v2 |
| - | nutanix_key_management_server_v2 |
| - | nutanix_ssl_certificate_v2 |
| - | nutanix_cluster_profile_v2 |
| - | nutanix_storage_policy_v2 |

## Data Sources

| v1 datasources | v2 datasources |
|  :--- |  :--- |
| nutanix_cluster | nutanix_cluster_v2 |
| nutanix_clusters | nutanix_clusters_v2 |
| nutanix_host | nutanix_host_v2 |
| nutanix_hosts | nutanix_hosts_v2 |
| nutanix_subnet | nutanix_subnet_v2 |
| nutanix_subnets | nutanix_subnets_v2 |
| nutanix_vpc | nutanix_vpc_v2 |
| nutanix_vpcs | nutanix_vpcs_v2 |
| nutanix_pbr | nutanix_pbr_v2 |
| nutanix_pbrs | nutanix_pbrs_v2 |
| nutanix_floating_ip | nutanix_floating_ip_v2 |
| nutanix_floating_ips | nutanix_floating_ips_v2 |
| nutanix_address_group | nutanix_address_group_v2 |
| nutanix_address_groups | nutanix_address_groups_v2 |
| nutanix_service_group | nutanix_service_group_v2 |
| nutanix_service_groups | nutanix_service_groups_v2 |
| nutanix_network_security_rule | nutanix_network_security_policy_v2 |
| - | nutanix_network_security_policies_v2 |
| nutanix_role | nutanix_role_v2 |
| nutanix_roles | nutanix_roles_v2 |
| nutanix_permission | nutanix_operation_v2 |
| nutanix_permissions | nutanix_operations_v2 |
| nutanix_user | nutanix_user_v2 |
| nutanix_users | nutanix_users_v2 |
| nutanix_user_group | nutanix_user_group_v2 |
| nutanix_user_groups | nutanix_user_groups_v2 |
| nutanix_access_control_policy | nutanix_authorization_policy_v2 |
| nutanix_access_control_policies | nutanix_authorization_policies_v2 |
| - | nutanix_saml_identity_provider_v2 |
| - | nutanix_saml_identity_providers_v2 |
| - | nutanix_directory_service_v2 |
| - | nutanix_directory_services_v2 |
| nutanix_category_key | nutanix_category_v2 |
| - | nutanix_categories_v2 |
| nutanix_image | nutanix_image_v2 |
| - | nutanix_images_v2 |
| nutanix_virtual_machine | nutanix_virtual_machine_v2 |
| - | nutanix_virtual_machines_v2 |
| - | nutanix_ova_v2 |
| - | nutanix_ovas_v2 |
| - | nutanix_template_v2 |
| - | nutanix_templates_v2 |
| - | nutanix_ngt_configuration_v2 |
| - | nutanix_image_placement_policy_v2 |
| - | nutanix_image_placement_policies_v2 |
| - | nutanix_volume_group_v2 |
| - | nutanix_volume_groups_v2 |
| - | nutanix_volume_group_disk_v2 |
| - | nutanix_volume_group_disks_v2 |
| - | nutanix_volume_group_iscsi_clients_v2 |
| - | nutanix_volume_group_category_details_v2 |
| - | nutanix_volume_group_vms_v2 |
| - | nutanix_volume_iscsi_client_v2 |
| - | nutanix_volume_iscsi_clients_v2 |
| - | nutanix_recovery_point_v2 |
| - | nutanix_recovery_points_v2 |
| - | nutanix_vm_recovery_point_info_v2 |
| - | nutanix_storage_container_v2 |
| - | nutanix_storage_containers_v2 |
| - | nutanix_storage_container_stats_info_v2 |
| nutanix_project | - |
| nutanix_projects | - |
| nutanix_karbon_cluster_kubeconfig | - |
| nutanix_karbon_cluster | - |
| nutanix_karbon_clusters | - |
| nutanix_karbon_cluster_ssh | - |
| nutanix_karbon_private_registry | - |
| nutanix_karbon_private_registries | - |
| nutanix_protection_rule | - |
| nutanix_protection_rules | - |
| nutanix_recovery_plan | - |
| nutanix_recovery_plans | - |
| nutanix_foundation_hypervisor_isos | - |
| nutanix_foundation_discover_nodes | - |
|nutanix_foundation_nos_packages | - |
| nutanix_foundation_node_network_details | - |
| nutanix_foundation_central_api_keys | - |
| nutanix_foundation_central_list_api_keys | - |
| nutanix_foundation_central_imaged_nodes_list | - |
| nutanix_foundation_central_imaged_clusters_list | - |
| nutanix_foundation_central_cluster_details | - |
| nutanix_foundation_central_imaged_node_details | - |
| nutanix_ndb_sla | - |
| nutanix_ndb_slas | - |
| nutanix_ndb_profile | - |
| nutanix_ndb_profiles | - |
| nutanix_ndb_cluster | - |
| nutanix_ndb_clusters | - |
| nutanix_ndb_database | - |
| nutanix_ndb_databases | - |
| nutanix_ndb_time_machine | - |
| nutanix_ndb_time_machines | - |
| nutanix_ndb_clone | - |
| nutanix_ndb_clones | - |
| nutanix_ndb_snapshot | - |
| nutanix_ndb_snapshots | - |
| nutanix_ndb_tms_capability | - |
| nutanix_ndb_maintenance_window | - |
| nutanix_ndb_maintenance_windows | - |
| nutanix_ndb_tag | - |
| nutanix_ndb_tags | - |
| nutanix_ndb_network | - |
| nutanix_ndb_networks | - |
| nutanix_ndb_dbserver | - |
| nutanix_ndb_dbservers | - |
| nutanix_ndb_network_available_ips | - |
| - | nutanix_pc_v2 |
| - | nutanix_pcs_v2 |
| - | nutanix_restorable_pcs_v2 |
| - | nutanix_pc_restore_points_v2 |
| - | nutanix_pc_restore_point_v2 |
| - | nutanix_pc_backup_target_v2 |
| - | nutanix_pc_backup_targets_v2 |
| - | nutanix_pc_restore_source_v2
| - | nutanix_protected_resource_v2 |
| - | nutanix_protection_policy_v2 |
| - | nutanix_protection_policies_v2 |
| - | nutanix_lcm_status_v2 |
| - | nutanix_lcm_entities_v2 |
| - | nutanix_lcm_entity_v2 |
| - | nutanix_lcm_config_v2 |
| nutanix_self_service_app | - |
| nutanix_blueprint_runtime_editables | - |
| nutanix_self_service_snapshot_policy_list | - |
| nutanix_self_service_app_snapshots | - |
| - | nutanix_user_keys_v2 |
| - | nutanix_user_key_v2 |
| - | nutanix_object_store_v2 |
| - | nutanix_object_stores_v2 |
| - | nutanix_certificate_v2 |
| - | nutanix_certificates_v2 |
| - | nutanix_system_user_passwords_v2 |
| - | nutanix_key_management_server_v2 |
| - | nutanix_key_management_servers_v2 |
| - | nutanix_stigs_v2 |
| - | nutanix_ssl_certificate_v2 |
| - | nutanix_cluster_profile_v2 |
| - | nutanix_cluster_profiles_v2 |
| - | nutanix_storage_policy_v2 |
| - | nutanix_storage_policies_v2 |

## Example Usage

### Terraform 0.12 and below

```terraform
provider "nutanix" {
  username     = var.nutanix_username
  password     = var.nutanix_password
  endpoint     = var.nutanix_endpoint
  port         = var.nutanix_port
  insecure     = true
  wait_timeout = 10
}
```

### Terraform 0.13+

```terraform
terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

provider "nutanix" {
  username     = var.nutanix_username
  password     = var.nutanix_password
  endpoint     = var.nutanix_endpoint
  port         = var.nutanix_port
  insecure     = true
  wait_timeout = 10
}
```

## Argument Reference

The following arguments are used to configure the Nutanix Provider:
* `username` - **(Required)** This is the username for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_USERNAME` environment variable.
* `password` - **(Required)** This is the password for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_PASSWORD` environment variable.
* `endpoint` - **(Required)** This is the endpoint for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_ENDPOINT` environment variable.
* `insecure` - (Optional) This specifies whether to allow verify ssl certificates. This can also be specified with `NUTANIX_INSECURE`. Defaults to `false`.
* `port` - (Optional) This is the port for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_PORT` environment variable. Defaults to `9440`.
* `session_auth` - (Optional) This specifies whether to use [session authentication](#session-based-authentication). This can also be specified with the `NUTANIX_SESSION_AUTH` environment variable. Defaults to `true`
* `wait_timeout` - (Optional) This specifies the timeout on all resource operations in the provider in minutes. This can also be specified with the `NUTANIX_WAIT_TIMEOUT` environment variable. Defaults to `1`. Also see [resource timeouts](#resource-timeouts).
* `proxy_url` - (Optional) This specifies the url to proxy through to access the Prism Elements or Prism Central endpoint. This can also be specified with the `NUTANIX_PROXY_URL` environment variable.

### Session based Authentication

Session based authentication can be used which authenticates only once with basic authentication and uses a cookie for all further attempts.
The main benefit is a reduction in the time API calls take to complete. Sessions are only valid for 15 minutes.

Usage:

```terraform
provider "nutanix" {
  ...
  session_auth = true
  ...
}
```

## Notes

### Resource Timeouts
Currently, the only way to set a timeout is using the `wait_timeout` argument or `NUTANIX_WAIT_TIMEOUT` environment variable. This will set a timeout for all operations on all resources. This provider currently doesn't support specifying [operation timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts).

## Nutanix Foundation (>=v1.5.0-beta)

Going from 1.5.0-beta release of nutanix provider, two more params are added to provider configuration to support foundation components :

* `foundation_endpoint` - (Optional) This is the endpoint for foundation vm. This can also be specified with the `FOUNDATION_ENDPOINT` environment variable.
* `foundation_port` - (Optional) This is the port for foundation vm. This can also be specified with the `FOUNDATION_PORT` environment variable. Default is `8000`.

```terraform
terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = ">=1.5.0-beta"
    }
  }
}

provider "nutanix" {
  username            = var.nutanix_username
  password            = var.nutanix_password
  endpoint            = var.nutanix_endpoint
  port                = var.nutanix_port
  insecure            = true
  wait_timeout        = 10
  foundation_endpoint = var.foundation_endpoint
  foundation_port     = var.foundation_port
}
```

Foundation based examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/examples/foundation/

Foundation based modules & examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/modules/foundation/

## Nutanix Database Service (NDB) (>=v1.8.0)

Going from 1.8.0 release of nutanix provider, some params are added to provider configuration to support Nutanix Database Service (NDB) components :

* `ndb_username` - (Optional) This is the username for the NDB instance. This can also be specified with the `NDB_USERNAME` environment variable.
* `ndb_password` - (Optional) This is the password for the NDB instance. This can also be specified with the `NDB_PASSWORD` environment variable.
* `ndb_endpoint` - (Optional) This is the endpoint for the NDB instance. This can also be specified with the `NDB_ENDPOINT` environment variable.

```terraform
terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = ">=1.8.0"
    }
  }
}

provider "nutanix" {
  username            = var.nutanix_username
  password            = var.nutanix_password
  endpoint            = var.nutanix_endpoint
  port                = var.nutanix_port
  insecure            = true
  wait_timeout        = 10
  ndb_endpoint        = var.ndb_endpoint 
  ndb_username        = var.ndb_username
  ndb_password        = var.ndb_password
}
```

NDB based examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/examples/ndb/

## Provider configuration required details

Going from 1.8.0-beta release of nutanix provider, fields inside provider configuration would be mandatory as per the usecase : 

* `Prism Central & Karbon` : For prism central and karbon related resources and data sources, `username`, `password` & `endpoint` are manadatory.
* `Foundation` : For foundation related resources and data sources, `foundation_endpoint` in manadatory.
* `NDB` : For Nutanix Database Service (NDB) related resources and data sources. 

