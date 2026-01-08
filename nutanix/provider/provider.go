package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/internal"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/clusters"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/clustersv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/datapoliciesv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/dataprotectionv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/foundation"
	foundationCentral "github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/foundationCentral"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/iam"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/iamv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/lcmv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/ndb"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/networking"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/networkingv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/nke"
	objectstoresv2 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/objectsv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/passwordmanagerv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/prismv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/securityv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/selfservice"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/storagecontainersv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/vmmv2"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/volumesv2"
)

var requiredProviderFields map[string][]string = map[string][]string{
	"prism_central":      {"username", "password", "endpoint"},
	"karbon":             {"username", "password", "endpoint"},
	"foundation":         {"foundation_endpoint"},
	"foundation_central": {"username", "password", "endpoint"},
	"ndb":                {"ndb_endpoint", "ndb_username", "ndb_password"},
}

// Provider function returns the object that implements the terraform.ResourceProvider interface, specifically a schema.Provider
func Provider() *schema.Provider {
	// defines descriptions for ResourceProvider schema definitions
	descriptions := map[string]string{
		"username": "User name for Nutanix Prism. Could be\n" +
			"local cluster auth (e.g. 'admin') or directory auth.",

		"password": "Password for provided user name.",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",

		"session_auth": "Use session authentification instead of basic auth for each request",

		"port": "Port for Nutanix Prism.",

		"wait_timeout": "Set if you know that the creation o update of a resource may take long time (minutes)",

		"endpoint": "URL for Nutanix Prism (e.g IP or FQDN for cluster VIP\n" +
			"note, this is never the data services VIP, and should not be an\n" +
			"individual CVM address, as this would cause calls to fail during\n" +
			"cluster lifecycle management operations, such as AOS upgrades.",

		"foundation_endpoint": "endpoint for foundation VM (eg. Foundation VM IP)",

		"foundation_port": "Port for foundation VM",

		"ndb_endpoint": "endpoint for Era VM (era ip)",
	}

	// Nutanix provider schema
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_USERNAME", nil),
				Description: descriptions["username"],
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PASSWORD", nil),
				Description: descriptions["password"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_INSECURE", false),
				Description: descriptions["insecure"],
			},
			"session_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_SESSION_AUTH", false),
				Description: descriptions["session_auth"],
			},
			"port": {
				Type:        schema.TypeString,
				Default:     "9440",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PORT", false),
				Description: descriptions["port"],
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_ENDPOINT", nil),
				Description: descriptions["endpoint"],
			},
			"wait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_WAIT_TIMEOUT", nil),
				Description: descriptions["wait_timeout"],
			},
			"proxy_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PROXY_URL", nil),
				Description: descriptions["proxy_url"],
			},
			"foundation_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FOUNDATION_ENDPOINT", nil),
				Description: descriptions["foundation_endpoint"],
			},
			"foundation_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "8000",
				DefaultFunc: schema.EnvDefaultFunc("FOUNDATION_PORT", nil),
				Description: descriptions["foundation_port"],
			},
			"ndb_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NDB_ENDPOINT", nil),
				Description: descriptions["ndb_endpoint"],
			},
			"ndb_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NDB_USERNAME", nil),
				Description: descriptions["ndb_username"],
			},
			"ndb_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NDB_PASSWORD", nil),
				Description: descriptions["ndb_password"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nutanix_image":                                   vmm.DataSourceNutanixImage(),
			"nutanix_subnet":                                  networking.DataSourceNutanixSubnet(),
			"nutanix_subnets":                                 networking.DataSourceNutanixSubnets(),
			"nutanix_cluster":                                 clusters.DataSourceNutanixCluster(),
			"nutanix_clusters":                                clusters.DataSourceNutanixClusters(),
			"nutanix_virtual_machine":                         vmm.DataSourceNutanixVirtualMachine(),
			"nutanix_category_key":                            prism.DataSourceNutanixCategoryKey(),
			"nutanix_network_security_rule":                   networking.DataSourceNutanixNetworkSecurityRule(),
			"nutanix_host":                                    clusters.DataSourceNutanixHost(),
			"nutanix_hosts":                                   clusters.DataSourceNutanixHosts(),
			"nutanix_access_control_policy":                   iam.DataSourceNutanixAccessControlPolicy(),
			"nutanix_access_control_policies":                 iam.DataSourceNutanixAccessControlPolicies(),
			"nutanix_project":                                 prism.DataSourceNutanixProject(),
			"nutanix_projects":                                prism.DataSourceNutanixProjects(),
			"nutanix_role":                                    iam.DataSourceNutanixRole(),
			"nutanix_roles":                                   iam.DataSourceNutanixRoles(),
			"nutanix_user":                                    iam.DataSourceNutanixUser(),
			"nutanix_user_group":                              iam.DataSourceNutanixUserGroup(),
			"nutanix_users":                                   iam.DataSourceNutanixUsers(),
			"nutanix_user_groups":                             iam.DataSourceNutanixUserGroups(),
			"nutanix_permission":                              iam.DataSourceNutanixPermission(),
			"nutanix_permissions":                             iam.DataSourceNutanixPermissions(),
			"nutanix_karbon_cluster_kubeconfig":               nke.DataSourceNutanixKarbonClusterKubeconfig(),
			"nutanix_karbon_cluster":                          nke.DataSourceNutanixKarbonCluster(),
			"nutanix_karbon_clusters":                         nke.DataSourceNutanixKarbonClusters(),
			"nutanix_karbon_cluster_ssh":                      nke.DataSourceNutanixKarbonClusterSSH(),
			"nutanix_karbon_private_registry":                 nke.DataSourceNutanixKarbonPrivateRegistry(),
			"nutanix_karbon_private_registries":               nke.DataSourceNutanixKarbonPrivateRegistries(),
			"nutanix_protection_rule":                         prism.DataSourceNutanixProtectionRule(),
			"nutanix_protection_rules":                        prism.DataSourceNutanixProtectionRules(),
			"nutanix_recovery_plan":                           prism.DataSourceNutanixRecoveryPlan(),
			"nutanix_recovery_plans":                          prism.DataSourceNutanixRecoveryPlans(),
			"nutanix_address_groups":                          networking.DataSourceNutanixAddressGroups(),
			"nutanix_address_group":                           networking.DataSourceNutanixAddressGroup(),
			"nutanix_service_group":                           networking.DataSourceNutanixServiceGroup(),
			"nutanix_service_groups":                          networking.DataSourceNutanixServiceGroups(),
			"nutanix_foundation_hypervisor_isos":              foundation.DataSourceFoundationHypervisorIsos(),
			"nutanix_foundation_discover_nodes":               foundation.DataSourceFoundationDiscoverNodes(),
			"nutanix_foundation_nos_packages":                 foundation.DataSourceFoundationNOSPackages(),
			"nutanix_foundation_node_network_details":         foundation.DataSourceNodeNetworkDetails(),
			"nutanix_assert_helper":                           internal.DataSourceAssertHelper(),
			"nutanix_foundation_central_api_keys":             foundationCentral.DataSourceNutanixFCAPIKeys(),
			"nutanix_foundation_central_list_api_keys":        foundationCentral.DataSourceNutanixFCListAPIKeys(),
			"nutanix_foundation_central_imaged_nodes_list":    foundationCentral.DataSourceNutanixFCImagedNodesList(),
			"nutanix_foundation_central_imaged_clusters_list": foundationCentral.DataSourceNutanixFCImagedClustersList(),
			"nutanix_foundation_central_cluster_details":      foundationCentral.DataSourceNutanixFCClusterDetails(),
			"nutanix_foundation_central_imaged_node_details":  foundationCentral.DataSourceFCImagedNodeDetails(),
			"nutanix_vpc":                                     networking.DataSourceNutanixVPC(),
			"nutanix_vpcs":                                    networking.DataSourceNutanixVPCs(),
			"nutanix_pbr":                                     networking.DataSourceNutanixPbr(),
			"nutanix_pbrs":                                    networking.DataSourceNutanixPbrs(),
			"nutanix_floating_ip":                             networking.DataSourceNutanixFloatingIP(),
			"nutanix_floating_ips":                            networking.DataSourceNutanixFloatingIPs(),
			"nutanix_static_routes":                           networking.DataSourceNutanixStaticRoute(),
			"nutanix_ndb_sla":                                 ndb.DataSourceNutanixEraSLA(),
			"nutanix_ndb_slas":                                ndb.DataSourceNutanixEraSLAs(),
			"nutanix_ndb_profile":                             ndb.DataSourceNutanixEraProfile(),
			"nutanix_ndb_profiles":                            ndb.DataSourceNutanixEraProfiles(),
			"nutanix_ndb_cluster":                             ndb.DataSourceNutanixEraCluster(),
			"nutanix_ndb_clusters":                            ndb.DataSourceNutanixEraClusters(),
			"nutanix_ndb_database":                            ndb.DataSourceNutanixEraDatabase(),
			"nutanix_ndb_databases":                           ndb.DataSourceNutanixEraDatabases(),
			"nutanix_ndb_time_machine":                        ndb.DataSourceNutanixNDBTimeMachine(),
			"nutanix_ndb_time_machines":                       ndb.DataSourceNutanixNDBTimeMachines(),
			"nutanix_ndb_clone":                               ndb.DataSourceNutanixNDBClone(),
			"nutanix_ndb_clones":                              ndb.DataSourceNutanixNDBClones(),
			"nutanix_ndb_snapshot":                            ndb.DataSourceNutanixNDBSnapshot(),
			"nutanix_ndb_snapshots":                           ndb.DataSourceNutanixNDBSnapshots(),
			"nutanix_ndb_tms_capability":                      ndb.DataSourceNutanixNDBTmsCapability(),
			"nutanix_ndb_maintenance_window":                  ndb.DataSourceNutanixNDBMaintenanceWindow(),
			"nutanix_ndb_maintenance_windows":                 ndb.DataSourceNutanixNDBMaintenanceWindows(),
			"nutanix_ndb_tag":                                 ndb.DataSourceNutanixNDBTag(),
			"nutanix_ndb_tags":                                ndb.DataSourceNutanixNDBTags(),
			"nutanix_ndb_network":                             ndb.DataSourceNutanixEraNetwork(),
			"nutanix_ndb_networks":                            ndb.DataSourceNutanixEraNetworks(),
			"nutanix_ndb_dbserver":                            ndb.DataSourceNutanixNDBDBServer(),
			"nutanix_ndb_dbservers":                           ndb.DataSourceNutanixNDBDBServers(),
			"nutanix_ndb_network_available_ips":               ndb.DataSourceNutanixNDBProfileAvailableIPs(),
			"nutanix_self_service_snapshot_policy_list":       selfservice.DataSourceNutanixSnapshotPolicy(),
			"nutanix_self_service_app":                        selfservice.DatsourceNutanixCalmApp(),
			"nutanix_blueprint_runtime_editables":             selfservice.DatsourceNutanixCalmRuntimeEditables(),
			"nutanix_self_service_app_snapshots":              selfservice.DataSourceNutanixCalmSnapshots(),
			"nutanix_subnet_v2":                               networkingv2.DataSourceNutanixSubnetV2(),
			"nutanix_subnets_v2":                              networkingv2.DataSourceNutanixSubnetsV2(),
			"nutanix_vpc_v2":                                  networkingv2.DataSourceNutanixVPCv2(),
			"nutanix_vpcs_v2":                                 networkingv2.DataSourceNutanixVPCsv2(),
			"nutanix_floating_ip_v2":                          networkingv2.DatasourceNutanixFloatingIPV2(),
			"nutanix_floating_ips_v2":                         networkingv2.DatasourceNutanixFloatingIPsV2(),
			"nutanix_network_security_policy_v2":              networkingv2.DataSourceNutanixNetworkSecurityPolicyV2(),
			"nutanix_network_security_policies_v2":            networkingv2.DataSourceNutanixNetworkSecurityPoliciesV2(),
			"nutanix_route_table_v2":                          networkingv2.DatasourceNutanixRouteTableV2(),
			"nutanix_route_tables_v2":                         networkingv2.DatasourceNutanixRouteTablesV2(),
			"nutanix_route_v2":                                networkingv2.DatasourceNutanixRouteV2(),
			"nutanix_routes_v2":                               networkingv2.DatasourceNutanixRoutesV2(),
			"nutanix_pbr_v2":                                  networkingv2.DatasourceNutanixPbrV2(),
			"nutanix_pbrs_v2":                                 networkingv2.DatasourceNutanixPbrsV2(),
			"nutanix_service_group_v2":                        networkingv2.DatasourceNutanixServiceGroupV2(),
			"nutanix_service_groups_v2":                       networkingv2.DatasourceNutanixServiceGroupsV2(),
			"nutanix_address_group_v2":                        networkingv2.DatasourceNutanixAddressGroupV2(),
			"nutanix_address_groups_v2":                       networkingv2.DatasourceNutanixAddressGroupsV2(),
			"nutanix_directory_service_v2":                    iamv2.DatasourceNutanixDirectoryServiceV2(),
			"nutanix_directory_services_v2":                   iamv2.DatasourceNutanixDirectoryServicesV2(),
			"nutanix_saml_identity_provider_v2":               iamv2.DatasourceNutanixSamlIDPV2(),
			"nutanix_saml_identity_providers_v2":              iamv2.DatasourceNutanixSamlIDPsV2(),
			"nutanix_user_group_v2":                           iamv2.DatasourceNutanixUserGroupV2(),
			"nutanix_user_groups_v2":                          iamv2.DatasourceNutanixUserGroupsV2(),
			"nutanix_roles_v2":                                iamv2.DatasourceNutanixRolesV2(),
			"nutanix_role_v2":                                 iamv2.DatasourceNutanixRoleV2(),
			"nutanix_operation_v2":                            iamv2.DatasourceNutanixOperationV2(),
			"nutanix_operations_v2":                           iamv2.DatasourceNutanixOperationsV2(),
			"nutanix_user_v2":                                 iamv2.DatasourceNutanixUserV2(),
			"nutanix_users_v2":                                iamv2.DatasourceNutanixUsersV2(),
			"nutanix_authorization_policy_v2":                 iamv2.DatasourceNutanixAuthorizationPolicyV2(),
			"nutanix_authorization_policies_v2":               iamv2.DatasourceNutanixAuthorizationPoliciesV2(),
			"nutanix_user_keys_v2":                            iamv2.DatasourceNutanixUserKeysV2(),
			"nutanix_user_key_v2":                             iamv2.DatasourceNutanixUserKeyV2(),
			"nutanix_storage_container_v2":                    storagecontainersv2.DatasourceNutanixStorageContainerV2(),
			"nutanix_storage_containers_v2":                   storagecontainersv2.DatasourceNutanixStorageContainersV2(),
			"nutanix_storage_container_stats_info_v2":         storagecontainersv2.DatasourceNutanixStorageStatsInfoV2(),
			"nutanix_category_v2":                             prismv2.DatasourceNutanixCategoryV2(),
			"nutanix_categories_v2":                           prismv2.DatasourceNutanixCategoriesV2(),
			"nutanix_pc_v2":                                   prismv2.DatasourceNutanixFetchPcV2(),
			"nutanix_pcs_v2":                                  prismv2.DatasourceNutanixListPcsV2(),
			"nutanix_restorable_pcs_v2":                       prismv2.DatasourceNutanixListRestorablePcsV2(),
			"nutanix_pc_restore_points_v2":                    prismv2.DatasourceNutanixFetchRestorePointsV2(),
			"nutanix_pc_restore_point_v2":                     prismv2.DatasourceNutanixFetchRestorePointV2(),
			"nutanix_pc_backup_target_v2":                     prismv2.DatasourceNutanixBackupTargetV2(),
			"nutanix_pc_backup_targets_v2":                    prismv2.DatasourceNutanixBackupTargetsV2(),
			"nutanix_pc_restore_source_v2":                    prismv2.DatasourceNutanixRestoreSourceV2(),
			"nutanix_volume_groups_v2":                        volumesv2.DatasourceNutanixVolumeGroupsV2(),
			"nutanix_volume_group_v2":                         volumesv2.DatasourceNutanixVolumeGroupV2(),
			"nutanix_volume_group_disks_v2":                   volumesv2.DatasourceNutanixVolumeDisksV2(),
			"nutanix_volume_group_disk_v2":                    volumesv2.DatasourceNutanixVolumeDiskV2(),
			"nutanix_volume_iscsi_clients_v2":                 volumesv2.DatasourceNutanixVolumeIscsiClientsV2(),
			"nutanix_volume_iscsi_client_v2":                  volumesv2.DatasourceNutanixVolumeIscsiClientV2(),
			"nutanix_recovery_point_v2":                       dataprotectionv2.DatasourceNutanixRecoveryPointV2(),
			"nutanix_recovery_points_v2":                      dataprotectionv2.DatasourceNutanixRecoveryPointsV2(),
			"nutanix_vm_recovery_point_info_v2":               dataprotectionv2.DatasourceNutanixVMRecoveryPointInfoV2(),
			"nutanix_protected_resource_v2":                   dataprotectionv2.DatasourceNutanixGetProtectedResourceV2(),
			"nutanix_protection_policy_v2":                    datapoliciesv2.DatasourceNutanixProtectionPolicyV2(),
			"nutanix_protection_policies_v2":                  datapoliciesv2.DatasourceNutanixProtectionPoliciesV2(),
			"nutanix_storage_policy_v2":                       datapoliciesv2.DataSourceNutanixStoragePolicyV2(),
			"nutanix_storage_policies_v2":                     datapoliciesv2.DataSourceNutanixStoragePoliciesV2(),
			"nutanix_image_v2":                                vmmv2.DatasourceNutanixImageV4(),
			"nutanix_images_v2":                               vmmv2.DatasourceNutanixImagesV4(),
			"nutanix_ova_v2":                                  vmmv2.DatasourceNutanixOvaV2(),
			"nutanix_ovas_v2":                                 vmmv2.DatasourceNutanixOvasV2(),
			"nutanix_virtual_machine_v2":                      vmmv2.DatasourceNutanixVirtualMachineV4(),
			"nutanix_virtual_machines_v2":                     vmmv2.DatasourceNutanixVirtualMachinesV4(),
			"nutanix_template_v2":                             vmmv2.DatasourceNutanixTemplateV2(),
			"nutanix_templates_v2":                            vmmv2.DatasourceNutanixTemplatesV2(),
			"nutanix_ngt_configuration_v2":                    vmmv2.DatasourceNutanixNGTConfigurationV4(),
			"nutanix_image_placement_policy_v2":               vmmv2.DatasourceNutanixImagePlacementV4(),
			"nutanix_image_placement_policies_v2":             vmmv2.DatasourceNutanixImagePlacementsV4(),
			"nutanix_cluster_v2":                              clustersv2.DatasourceNutanixClusterEntityV2(),
			"nutanix_clusters_v2":                             clustersv2.DatasourceNutanixClusterEntitiesV2(),
			"nutanix_system_user_passwords_v2":                passwordmanagerv2.DataSourceNutanixPasswordManagersV2(),
			"nutanix_host_v2":                                 clustersv2.DatasourceNutanixHostEntityV2(),
			"nutanix_hosts_v2":                                clustersv2.DatasourceNutanixHostEntitiesV2(),
			"nutanix_ssl_certificate_v2":                      clustersv2.DatasourceNutanixSSLCertificateV2(),
			"nutanix_cluster_profile_v2":                      clustersv2.DatasourceNutanixClusterProfileV2(),
			"nutanix_cluster_profiles_v2":                     clustersv2.DatasourceNutanixClusterProfilesV2(),
			"nutanix_lcm_status_v2":                           lcmv2.DatasourceNutanixLcmStatusV2(),
			"nutanix_lcm_entities_v2":                         lcmv2.DatasourceNutanixLcmEntitiesV2(),
			"nutanix_lcm_entity_v2":                           lcmv2.DatasourceNutanixLcmEntityV2(),
			"nutanix_lcm_config_v2":                           lcmv2.DatasourceNutanixLcmConfigV2(),
			"nutanix_object_store_v2":                         objectstoresv2.DatasourceNutanixObjectStoreV2(),
			"nutanix_object_stores_v2":                        objectstoresv2.DatasourceNutanixObjectStoresV2(),
			"nutanix_certificate_v2":                          objectstoresv2.DatasourceNutanixObjectStoreCertificateV2(),
			"nutanix_certificates_v2":                         objectstoresv2.DatasourceNutanixObjectStoreCertificatesV2(),
			"nutanix_key_management_server_v2":                securityv2.DatasourceNutanixKeyManagementServerV2(),
			"nutanix_key_management_servers_v2":               securityv2.DatasourceNutanixKeyManagementServersV2(),
			"nutanix_stigs_v2":                                securityv2.DatasourceNutanixStigsControlsV2(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nutanix_virtual_machine":                         vmm.ResourceNutanixVirtualMachine(),
			"nutanix_image":                                   vmm.ResourceNutanixImage(),
			"nutanix_subnet":                                  networking.ResourceNutanixSubnet(),
			"nutanix_category_key":                            prism.ResourceNutanixCategoryKey(),
			"nutanix_category_value":                          prism.ResourceNutanixCategoryValue(),
			"nutanix_network_security_rule":                   networking.ResourceNutanixNetworkSecurityRule(),
			"nutanix_access_control_policy":                   prism.ResourceNutanixAccessControlPolicy(),
			"nutanix_project":                                 prism.ResourceNutanixProject(),
			"nutanix_role":                                    iam.ResourceNutanixRole(),
			"nutanix_user":                                    iam.ResourceNutanixUser(),
			"nutanix_karbon_cluster":                          nke.ResourceNutanixKarbonCluster(),
			"nutanix_karbon_private_registry":                 nke.ResourceNutanixKarbonPrivateRegistry(),
			"nutanix_karbon_worker_nodepool":                  nke.ResourceNutanixKarbonWorkerNodePool(),
			"nutanix_protection_rule":                         prism.ResourceNutanixProtectionRule(),
			"nutanix_recovery_plan":                           prism.ResourceNutanixRecoveryPlan(),
			"nutanix_service_group":                           networking.ResourceNutanixServiceGroup(),
			"nutanix_address_group":                           networking.ResourceNutanixAddressGroup(),
			"nutanix_foundation_image_nodes":                  foundation.ResourceFoundationImageNodes(),
			"nutanix_foundation_ipmi_config":                  foundation.ResourceNutanixFoundationIPMIConfig(),
			"nutanix_foundation_image":                        foundation.ResourceNutanixFoundationImage(),
			"nutanix_foundation_central_image_cluster":        foundationCentral.ResourceNutanixFCImageCluster(),
			"nutanix_foundation_central_api_keys":             foundationCentral.ResourceNutanixFCAPIKeys(),
			"nutanix_foundation_central_onboard_nodes":        foundationCentral.ResourceNutanixFCOnboardNodes(),
			"nutanix_vpc":                                     networking.ResourceNutanixVPC(),
			"nutanix_pbr":                                     networking.ResourceNutanixPbr(),
			"nutanix_floating_ip":                             networking.ResourceNutanixFloatingIP(),
			"nutanix_static_routes":                           networking.ResourceNutanixStaticRoute(),
			"nutanix_user_groups":                             iam.ResourceNutanixUserGroups(),
			"nutanix_ndb_database":                            ndb.ResourceDatabaseInstance(),
			"nutanix_ndb_sla":                                 ndb.ResourceNutanixNDBSla(),
			"nutanix_ndb_database_restore":                    ndb.ResourceNutanixNDBDatabaseRestore(),
			"nutanix_ndb_log_catchups":                        ndb.ResourceNutanixNDBLogCatchUps(),
			"nutanix_ndb_profile":                             ndb.ResourceNutanixNDBProfile(),
			"nutanix_ndb_software_version_profile":            ndb.ResourceNutanixNDBSoftwareVersionProfile(),
			"nutanix_ndb_scale_database":                      ndb.ResourceNutanixNDBScaleDatabase(),
			"nutanix_ndb_database_scale":                      ndb.ResourceNutanixNDBScaleDatabase(),
			"nutanix_ndb_register_database":                   ndb.ResourceNutanixNDBRegisterDatabase(),
			"nutanix_ndb_database_snapshot":                   ndb.ResourceNutanixNDBDatabaseSnapshot(),
			"nutanix_ndb_clone":                               ndb.ResourceNutanixNDBClone(),
			"nutanix_ndb_authorize_dbserver":                  ndb.ResourceNutanixNDBAuthorizeDBServer(),
			"nutanix_ndb_linked_databases":                    ndb.ResourceNutanixNDBLinkedDB(),
			"nutanix_ndb_maintenance_window":                  ndb.ResourceNutanixNDBMaintenanceWindow(),
			"nutanix_ndb_maintenance_task":                    ndb.ResourceNutanixNDBMaintenanceTask(),
			"nutanix_ndb_tms_cluster":                         ndb.ResourceNutanixNDBTmsCluster(),
			"nutanix_ndb_tag":                                 ndb.ResourceNutanixNDBTags(),
			"nutanix_ndb_network":                             ndb.ResourceNutanixNDBNetwork(),
			"nutanix_ndb_dbserver_vm":                         ndb.ResourceNutanixNDBServerVM(),
			"nutanix_ndb_register_dbserver":                   ndb.ResourceNutanixNDBRegisterDBServer(),
			"nutanix_ndb_stretched_vlan":                      ndb.ResourceNutanixNDBStretchedVlan(),
			"nutanix_ndb_clone_refresh":                       ndb.ResourceNutanixNDBCloneRefresh(),
			"nutanix_ndb_cluster":                             ndb.ResourceNutanixNDBCluster(),
			"nutanix_self_service_app_provision":              selfservice.ResourceNutanixCalmAppProvision(),
			"nutanix_self_service_app_patch":                  selfservice.ResourceNutanixCalmAppPatch(),
			"nutanix_self_service_app_recovery_point":         selfservice.ResourceNutanixCalmAppRecoveryPoint(),
			"nutanix_self_service_app_custom_action":          selfservice.ResourceNutanixCalmAppCustomAction(),
			"nutanix_self_service_app_restore":                selfservice.ResourceNutanixCalmAppRestore(),
			"nutanix_subnet_v2":                               networkingv2.ResourceNutanixSubnetV2(),
			"nutanix_floating_ip_v2":                          networkingv2.ResourceNutanixFloatingIPv2(),
			"nutanix_vpc_v2":                                  networkingv2.ResourceNutanixVPCsV2(),
			"nutanix_network_security_policy_v2":              networkingv2.ResourceNutanixNetworkSecurityPolicyV2(),
			"nutanix_routes_v2":                               networkingv2.ResourceNutanixRoutesV2(),
			"nutanix_pbr_v2":                                  networkingv2.ResourceNutanixPbrsV2(),
			"nutanix_service_groups_v2":                       networkingv2.ResourceNutanixServiceGroupsV2(),
			"nutanix_address_groups_v2":                       networkingv2.ResourceNutanixAddressGroupsV2(),
			"nutanix_directory_services_v2":                   iamv2.ResourceNutanixDirectoryServicesV2(),
			"nutanix_user_groups_v2":                          iamv2.ResourceNutanixUserGroupsV2(),
			"nutanix_roles_v2":                                iamv2.ResourceNutanixRolesV2(),
			"nutanix_users_v2":                                iamv2.ResourceNutanixUserV2(),
			"nutanix_authorization_policy_v2":                 iamv2.ResourceNutanixAuthPoliciesV2(),
			"nutanix_saml_identity_providers_v2":              iamv2.ResourceNutanixSamlIdpV2(),
			"nutanix_user_key_v2":                             iamv2.ResourceNutanixUserKeyV2(),
			"nutanix_user_key_revoke_v2":                      iamv2.ResourceNutanixUserRevokeKeyV2(),
			"nutanix_storage_containers_v2":                   storagecontainersv2.ResourceNutanixStorageContainersV2(),
			"nutanix_category_v2":                             prismv2.ResourceNutanixCategoriesV2(),
			"nutanix_pc_deploy_v2":                            prismv2.ResourceNutanixDeployPcV2(),
			"nutanix_pc_backup_target_v2":                     prismv2.ResourceNutanixBackupTargetV2(),
			"nutanix_pc_restore_source_v2":                    prismv2.ResourceNutanixRestoreSourceV2(),
			"nutanix_pc_restore_v2":                           prismv2.ResourceNutanixRestorePcV2(),
			"nutanix_pc_registration_v2":                      prismv2.ResourceNutanixClusterPCRegistrationV2(),
			"nutanix_pc_unregistration_v2":                    prismv2.ResourceNutanixUnregisterClusterV2(),
			"nutanix_volume_group_v2":                         volumesv2.ResourceNutanixVolumeGroupV2(),
			"nutanix_associate_category_to_volume_group_v2":   volumesv2.ResourceNutanixAssociateCategoryToVolumeGroupV2(),
			"nutanix_volume_group_disk_v2":                    volumesv2.ResourceNutanixVolumeGroupDiskV2(),
			"nutanix_volume_group_iscsi_client_v2":            volumesv2.ResourceNutanixVolumeGroupIscsiClientV2(),
			"nutanix_volume_group_vm_v2":                      volumesv2.ResourceNutanixVolumeAttachVMToVolumeGroupV2(),
			"nutanix_recovery_points_v2":                      dataprotectionv2.ResourceNutanixRecoveryPointsV2(),
			"nutanix_recovery_point_replicate_v2":             dataprotectionv2.ResourceNutanixRecoveryPointReplicateV2(),
			"nutanix_recovery_point_restore_v2":               dataprotectionv2.ResourceNutanixRecoveryPointRestoreV2(),
			"nutanix_promote_protected_resource_v2":           dataprotectionv2.ResourceNutanixPromoteProtectedResourceV2(),
			"nutanix_restore_protected_resource_v2":           dataprotectionv2.ResourceNutanixRestoreProtectedResourceV2(),
			"nutanix_protection_policy_v2":                    datapoliciesv2.ResourceNutanixProtectionPoliciesV2(),
			"nutanix_storage_policy_v2":                       datapoliciesv2.ResourceNutanixStoragePoliciesV2(),
			"nutanix_vm_revert_v2":                            vmmv2.ResourceNutanixRevertVMRecoveryPointV2(),
			"nutanix_virtual_machine_v2":                      vmmv2.ResourceNutanixVirtualMachineV2(),
			"nutanix_vm_shutdown_action_v2":                   vmmv2.ResourceNutanixVmsShutdownActionV2(),
			"nutanix_vm_cdrom_insert_eject_v2":                vmmv2.ResourceNutanixVmsCdRomsInsertEjectV2(),
			"nutanix_deploy_templates_v2":                     vmmv2.ResourceNutanixTemplateDeployV2(),
			"nutanix_template_v2":                             vmmv2.ResourceNutanixTemplatesV2(),
			"nutanix_template_guest_os_actions_v2":            vmmv2.ResourceNutanixTemplateActionsV2(),
			"nutanix_ngt_installation_v2":                     vmmv2.ResourceNutanixNGTInstallationV2(),
			"nutanix_ngt_upgrade_v2":                          vmmv2.ResourceNutanixNGTUpgradeV2(),
			"nutanix_ngt_insert_iso_v2":                       vmmv2.ResourceNutanixNGTInsertIsoV2(),
			"nutanix_vm_clone_v2":                             vmmv2.ResourceNutanixVMCloneV2(),
			"nutanix_vm_gc_update_v2":                         vmmv2.ResourceNutanixVMGCUpdateV2(),
			"nutanix_images_v2":                               vmmv2.ResourceNutanixImageV4(),
			"nutanix_ova_v2":                                  vmmv2.ResourceNutanixOvaV2(),
			"nutanix_ova_vm_deploy_v2":                        vmmv2.ResourceNutanixOvaVMDeploymentV2(),
			"nutanix_ova_download_v2":                         vmmv2.ResourceNutanixOvaDownloadV2(),
			"nutanix_vm_network_device_assign_ip_v2":          vmmv2.ResourceNutanixVmsNetworkDeviceAssignIPV2(),
			"nutanix_vm_network_device_migrate_v2":            vmmv2.ResourceNutanixVmsNetworkDeviceMigrateV2(),
			"nutanix_image_placement_policy_v2":               vmmv2.ResourceNutanixImagePlacementV2(),
			"nutanix_cluster_v2":                              clustersv2.ResourceNutanixClusterV2(),
			"nutanix_cluster_add_node_v2":                     clustersv2.ResourceNutanixClusterAddNodeV2(),
			"nutanix_clusters_discover_unconfigured_nodes_v2": clustersv2.ResourceNutanixClusterDiscoverUnconfiguredNodesV2(),
			"nutanix_clusters_unconfigured_node_networks_v2":  clustersv2.ResourceNutanixClusterUnconfiguredNodeNetworkV2(),
			"nutanix_ssl_certificate_v2":                      clustersv2.ResourceNutanixSSLCertificateV2(),
			"nutanix_cluster_profile_v2":                      clustersv2.ResourceNutanixClusterProfileV2(),
			"nutanix_password_change_request_v2":              passwordmanagerv2.ResourceNutanixPasswordManagerV2(),
			"nutanix_lcm_perform_inventory_v2":                lcmv2.ResourceNutanixLcmPerformInventoryV2(),
			"nutanix_lcm_prechecks_v2":                        lcmv2.ResourceNutanixPreChecksV2(),
			"nutanix_lcm_upgrade_v2":                          lcmv2.ResourceLcmUpgradeV2(),
			"nutanix_lcm_config_v2":                           lcmv2.ResourceNutanixLcmConfigV2(),
			"nutanix_object_store_v2":                         objectstoresv2.ResourceNutanixObjectStoresV2(),
			"nutanix_object_store_certificate_v2":             objectstoresv2.ResourceNutanixObjectStoreCertificateV2(),
			"nutanix_key_management_server_v2":                securityv2.ResourceNutanixKeyManagementServerV2(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// This function used to fetch the configuration params given to our provider which
// we will use to initialize a dummy client that interacts with API.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[DEBUG] config wait_timeout %d", d.Get("wait_timeout").(int))

	disabledProviders := make([]string, 0)
	// create warnings for disabled provider services
	var diags diag.Diagnostics
	for k, v := range requiredProviderFields {
		// check if any field is not provided
		for _, attr := range v {
			// for string fields
			if _, ok := d.GetOk(attr); !ok {
				disabledProviders = append(disabledProviders, k)
				break
			}
		}
	}

	if len(disabledProviders) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Disabled Providers: %s. Please provide required fields in provider configuration to enable them. Refer docs.", strings.Join(disabledProviders, ", ")),
		})
	}

	config := conns.Config{
		Endpoint:           d.Get("endpoint").(string),
		Username:           d.Get("username").(string),
		Password:           d.Get("password").(string),
		Insecure:           d.Get("insecure").(bool),
		SessionAuth:        d.Get("session_auth").(bool),
		Port:               d.Get("port").(string),
		WaitTimeout:        int64(d.Get("wait_timeout").(int)),
		ProxyURL:           d.Get("proxy_url").(string),
		FoundationEndpoint: d.Get("foundation_endpoint").(string),
		FoundationPort:     d.Get("foundation_port").(string),
		NdbEndpoint:        d.Get("ndb_endpoint").(string),
		NdbUsername:        d.Get("ndb_username").(string),
		NdbPassword:        d.Get("ndb_password").(string),
		RequiredFields:     requiredProviderFields,
	}
	c, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
