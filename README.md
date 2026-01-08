# Terraform Nutanix Provider

Terraform provider plugin to integrate with Nutanix Cloud Platform.

NOTE: The latest version of the Nutanix provider is [v2.4.0](https://github.com/nutanix/terraform-provider-nutanix/releases/tag/v2.4.0).

Modules based on Terraform Nutanix Provider can be found here : [Modules](https://github.com/nutanix/terraform-provider-nutanix/tree/master/modules)

## Build, Quality Status

 [![Go Report Card](https://goreportcard.com/badge/github.com/nutanix/terraform-provider-nutanix)](https://goreportcard.com/report/github.com/nutanix/terraform-provider-nutanix)
<!-- [![Maintainability](https://api.codeclimate.com/v1/badges/8b9e61df450276bbdbdb/maintainability)](https://codeclimate.com/github/nutanix/terraform-provider-nutanix/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/8b9e61df450276bbdbdb/test_coverage)](https://codeclimate.com/github/nutanix/terraform-provider-nutanix/test_coverage) -->

| Master                                                                                                                                                          | Develop                                                                                                                                                           |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [![Build Status](https://travis-ci.org/nutanix/terraform-provider-nutanix.svg?branch=master)](https://travis-ci.org/nutanix/terraform-provider-nutanix) | [![Build Status](https://travis-ci.org/nutanix/terraform-provider-nutanix.svg?branch=develop)](https://travis-ci.org/nutanix/terraform-provider-nutanix) |


### Requirements
* [Terraform](https://www.terraform.io/downloads.html) 0.12+
* [Go](https://golang.org/doc/install) 1.17+ (to build the provider plugin)
* This provider uses [SDKv2](https://www.terraform.io/plugin/sdkv2/sdkv2-intro) from release 1.3.0

## Introducing Nutanix Terraform Provider Version v2.4.0

We're excited to announce the release of Nutanix Terraform Provider Version 2.4.0!

### What's New in v2.4.0

- **New Resource Support**
  - **Key Management Server (Security)**: Create, Update, Read and Delete Key Management Servers secure data encryption keys when encryption is enabled.
  - **Security Technical Implementation Guide controls details (Security)**: Fetch the STIG controls details for STIG rules on each cluster.
  - **SSL Certification (Cluster Management)**: Provides the ability to manage SSL certificates for clusters. This includes the ability to retrieve and update SSL certificates for clusters.
  - **Cluster Profile (Cluster Management)**: Create, Update, Read and Delete cluster configuration profiles for consistent deployments.
  - **Associate/Disassociate Cluster from Cluster Profile (Cluster Management)**: Associate or Disassociate clusters to profiles for streamlined management.
  - **Associate/Disassociate Categories to Cluster (Cluster Management)**: Associate or Disassociate categories to clusters.
  - **Storage Policies (Data Policies)**: Create, Update, Read and Delete Storage Policy which helps in ease of storage management at scale.

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
   - From PC version 7.5 onwards, the resource domain is not supported by Projects API. As a result, Terraform support for this functionality (resource_doamin attribute) has been removed starting with the 2.4.0 release. [#1049](https://github.com/nutanix/terraform-provider-nutanix/issues/1049)


### Software Requirements
The provider is used to interact with the many resources and data sources supported by Nutanix, using Prism Central as the provider endpoint. To fully utilize the capabilities of version 2.4.0, ensure your Nutanix environment meets the following software requirements:
- Self Service version: 4.3.0 (Required only for running Self Service based resource and data source)
- AOS Version: 7.5 or later
- Prism Central Version: pc 7.5 or later
- Nutanix Terraform Provider Version: 2.4.0


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
| 2.2.0 | | | Self Service v4.1.0 | yes |
| 2.1.1 | 7.0.1, 7.0 | pc2024.3, pc2024.3.1 or later | | yes |
| 2.1.0 | 7.0.1, 7.0 | pc2024.3, pc2024.3.1 or later | | yes |
| 2.0.0 | 7.0 | pc2024.3 or later  | ndb v2.7, nke v2.8, foundation v5.7 | Yes |
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


## Example Usage

See the Examples folder for a handful of main.tf demos as well as some pre-compiled binaries.

We'll be refreshing these examples and binaries as we work through tech preview.

Long term, once this is upstream, no pre-compiled binaries will be needed, as terraform will automatically download on use.

## Configuration Reference

The following keys can be used to configure the provider.

* **endpoint** - (Required) IP address for the Nutanix Prism Central.
* **username** - (Required) Username for Nutanix Prism Central. Could be local cluster auth (e.g. `auth`) or directory auth.
* **password** - (Required) Password for the provided username.
* **port** - (Optional) Port for the Nutanix Prism Central. Default port is 9440.
* **insecure** - (Optional) Explicitly allow the provider to perform insecure SSL requests. If omitted, default value is false.
* **wait_timeout** - (optional) Set if you know that the creation o update of a resource may take long time (minutes).

```hcl
provider "nutanix" {
  username     = "admin"
  password     = "myPassword"
  port         = 9440
  endpoint     = "10.36.7.201"
  insecure     = true
  wait_timeout = 10
}
```

## From terraform-provider-nutanix v1.5.0-beta :

The following keys can be used to configure the provider.

* **endpoint** - (Optional) IP address for the Nutanix Prism Central.
* **username** - (Optional) Username for Nutanix Prism Central. Could be local cluster auth (e.g. `auth`) or directory auth.
* **password** - (Optional) Password for the provided username.
* **port** - (Optional) Port for the Nutanix Prism Central. Default port is 9440.
* **insecure** - (Optional) Explicitly allow the provider to perform insecure SSL requests. If omitted, default value is false.
* **wait_timeout** - (optional) Set if you know that the creation or update of a resource may take long time (minutes).
* **foundation_endpoint** - (optional) IP address of foundation vm.
* **foundation_port** - (optional) Port of foundation vm. Default port is 8000.

```hcl
provider "nutanix" {
  username            = "admin"
  password            = "myPassword"
  port                = 9440
  endpoint            = "10.36.7.201"
  insecure            = true
  wait_timeout        = 10
  foundation_endpoint = "10.xx.xx.xx"
  foundation_port     = 8000
}
```

## Additional fields for using Nutanix Database Service:

* **ndb_username** - (Optional) Username of Nutanix Database Service server
* **ndb_password** - (Optional) Password of Nutanix Database Service server
* **ndb_endpoint** - (Optional) IP of Nutanix Database Service server

```hcl
provider "nutanix" {
  ndb_username = var.ndb_username
  ndb_password = var.ndb_password
  ndb_endpoint = var.ndb_endpoint
}
```

### Provider Configuration Requirements & Warnings
From foundation getting released in 1.5.0-beta, provider configuration will accomodate prism central and foundation apis connection details. **It will show warnings for disabled api connections as per the attributes given in provider configuration in above mentioned format**. The below are the required attributes for corresponding provider componenets :
* endpoint, username and password are required fields for using Prism Central & Karbon based resources and data sources
* foundation_endpoint is required field for using Foundation based resources and data sources
* ndb_username, ndb_password and ndb_endpoint are required fields for using NDB based resources and data sources


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



## Developing the provider

The Nutanix Provider for Terraform is the work of many contributors. We appreciate your help!

* [Contribution Guidelines](./CONTRIBUTING.md)
* [Code of Conduct](./CODE_OF_CONDUCT.md)


## Support

-> **Note:** We now have a brand new developer-centric Support Program designed for organizations that require a deeper level of developer support to manage their Nutanix environment and build applications quickly and efficiently. As part of this new Advanced API/SDK Support Program, you will get access to trusted technical advisors who specialize in developer tools including Nutanix Terraform Provider and receive support for your unique development needs and custom integration queries. Visit our Support Portal - [Premium Add-On Support Programs](https://www.nutanix.com/support-services/product-support/premium-support-programs) to learn more about this program.

Customers not taking advantage of the  Advanced API/SDK Support Program will continue to receive the support through our standard, community-supported model. This community model also provides support for contributions to the open-sourceNutanix Terraform Provider repository .Visit https://portal.nutanix.com/kb/13424   for more details.


## Community

Nutanix is taking an inclusive approach to developing this new feature and welcomes customer feedback. Please see our development project on GitHub (you're here!), comment on requirements, design, code, and/or feel free to join us on Slack. Instructions on commenting, contributing, and joining our community Slack channel are all located within our GitHub Readme.

For a slack invite, please contact terraform@nutanix.com from your business email address, and we'll add you.
