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
	fc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/v1/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/v1/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/v1/ndb"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/v1/nke"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/v1/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/v2/networking"
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
			"nutanix_image":                                   prism.DataSourceNutanixImage(),
			"nutanix_subnet":                                  prism.DataSourceNutanixSubnet(),
			"nutanix_subnets":                                 prism.DataSourceNutanixSubnets(),
			"nutanix_cluster":                                 prism.DataSourceNutanixCluster(),
			"nutanix_clusters":                                prism.DataSourceNutanixClusters(),
			"nutanix_virtual_machine":                         prism.DataSourceNutanixVirtualMachine(),
			"nutanix_category_key":                            prism.DataSourceNutanixCategoryKey(),
			"nutanix_network_security_rule":                   prism.DataSourceNutanixNetworkSecurityRule(),
			"nutanix_host":                                    prism.DataSourceNutanixHost(),
			"nutanix_hosts":                                   prism.DataSourceNutanixHosts(),
			"nutanix_access_control_policy":                   prism.DataSourceNutanixAccessControlPolicy(),
			"nutanix_access_control_policies":                 prism.DataSourceNutanixAccessControlPolicies(),
			"nutanix_project":                                 prism.DataSourceNutanixProject(),
			"nutanix_projects":                                prism.DataSourceNutanixProjects(),
			"nutanix_role":                                    prism.DataSourceNutanixRole(),
			"nutanix_roles":                                   prism.DataSourceNutanixRoles(),
			"nutanix_user":                                    prism.DataSourceNutanixUser(),
			"nutanix_user_group":                              prism.DataSourceNutanixUserGroup(),
			"nutanix_users":                                   prism.DataSourceNutanixUsers(),
			"nutanix_user_groups":                             prism.DataSourceNutanixUserGroups(),
			"nutanix_permission":                              prism.DataSourceNutanixPermission(),
			"nutanix_permissions":                             prism.DataSourceNutanixPermissions(),
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
			"nutanix_address_groups":                          prism.DataSourceNutanixAddressGroups(),
			"nutanix_address_group":                           prism.DataSourceNutanixAddressGroup(),
			"nutanix_service_group":                           prism.DataSourceNutanixServiceGroup(),
			"nutanix_service_groups":                          prism.DataSourceNutanixServiceGroups(),
			"nutanix_foundation_hypervisor_isos":              foundation.DataSourceFoundationHypervisorIsos(),
			"nutanix_foundation_discover_nodes":               foundation.DataSourceFoundationDiscoverNodes(),
			"nutanix_foundation_nos_packages":                 foundation.DataSourceFoundationNOSPackages(),
			"nutanix_foundation_node_network_details":         foundation.DataSourceNodeNetworkDetails(),
			"nutanix_assert_helper":                           internal.DataSourceAssertHelper(),
			"nutanix_foundation_central_api_keys":             fc.DataSourceNutanixFCAPIKeys(),
			"nutanix_foundation_central_list_api_keys":        fc.DataSourceNutanixFCListAPIKeys(),
			"nutanix_foundation_central_imaged_nodes_list":    fc.DataSourceNutanixFCImagedNodesList(),
			"nutanix_foundation_central_imaged_clusters_list": fc.DataSourceNutanixFCImagedClustersList(),
			"nutanix_foundation_central_cluster_details":      fc.DataSourceNutanixFCClusterDetails(),
			"nutanix_foundation_central_imaged_node_details":  fc.DataSourceFCImagedNodeDetails(),
			"nutanix_vpc":                                     prism.DataSourceNutanixVPC(),
			"nutanix_vpcs":                                    prism.DataSourceNutanixVPCs(),
			"nutanix_pbr":                                     prism.DataSourceNutanixPbr(),
			"nutanix_pbrs":                                    prism.DataSourceNutanixPbrs(),
			"nutanix_floating_ip":                             prism.DataSourceNutanixFloatingIP(),
			"nutanix_floating_ips":                            prism.DataSourceNutanixFloatingIPs(),
			"nutanix_static_routes":                           prism.DataSourceNutanixStaticRoute(),
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
			"nutanix_subnet_v4":                               networking.DataSourceNutanixSubnetv4(),
			"nutanix_subnets_v4":                              networking.DataSourceNutanixSubnetsv4(),
			"nutanix_vpc_v4":                                  networking.DataSourceNutanixVPCv4(),
			"nutanix_vpcs_v4":                                 networking.DataSourceNutanixVPCsv4(),
			"nutanix_floating_ip_v4":                          networking.DatasourceNutanixFloatingIPV4(),
			"nutanix_floating_ips_v4":                         networking.DatasourceNutanixFloatingIPsV4(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nutanix_virtual_machine":                  prism.ResourceNutanixVirtualMachine(),
			"nutanix_image":                            prism.ResourceNutanixImage(),
			"nutanix_subnet":                           prism.ResourceNutanixSubnet(),
			"nutanix_category_key":                     prism.ResourceNutanixCategoryKey(),
			"nutanix_category_value":                   prism.ResourceNutanixCategoryValue(),
			"nutanix_network_security_rule":            prism.ResourceNutanixNetworkSecurityRule(),
			"nutanix_access_control_policy":            prism.ResourceNutanixAccessControlPolicy(),
			"nutanix_project":                          prism.ResourceNutanixProject(),
			"nutanix_role":                             prism.ResourceNutanixRole(),
			"nutanix_user":                             prism.ResourceNutanixUser(),
			"nutanix_karbon_cluster":                   nke.ResourceNutanixKarbonCluster(),
			"nutanix_karbon_private_registry":          nke.ResourceNutanixKarbonPrivateRegistry(),
			"nutanix_protection_rule":                  prism.ResourceNutanixProtectionRule(),
			"nutanix_recovery_plan":                    prism.ResourceNutanixRecoveryPlan(),
			"nutanix_service_group":                    prism.ResourceNutanixServiceGroup(),
			"nutanix_address_group":                    prism.ResourceNutanixAddressGroup(),
			"nutanix_foundation_image_nodes":           foundation.ResourceFoundationImageNodes(),
			"nutanix_foundation_ipmi_config":           foundation.ResourceNutanixFoundationIPMIConfig(),
			"nutanix_foundation_image":                 foundation.ResourceNutanixFoundationImage(),
			"nutanix_foundation_central_image_cluster": fc.ResourceNutanixFCImageCluster(),
			"nutanix_foundation_central_api_keys":      fc.ResourceNutanixFCAPIKeys(),
			"nutanix_vpc":                              prism.ResourceNutanixVPC(),
			"nutanix_pbr":                              prism.ResourceNutanixPbr(),
			"nutanix_floating_ip":                      prism.ResourceNutanixFloatingIP(),
			"nutanix_static_routes":                    prism.ResourceNutanixStaticRoute(),
			"nutanix_user_groups":                      prism.ResourceNutanixUserGroups(),
			"nutanix_ndb_database":                     ndb.ResourceDatabaseInstance(),
			"nutanix_ndb_sla":                          ndb.ResourceNutanixNDBSla(),
			"nutanix_ndb_database_restore":             ndb.ResourceNutanixNDBDatabaseRestore(),
			"nutanix_ndb_log_catchups":                 ndb.ResourceNutanixNDBLogCatchUps(),
			"nutanix_ndb_profile":                      ndb.ResourceNutanixNDBProfile(),
			"nutanix_ndb_software_version_profile":     ndb.ResourceNutanixNDBSoftwareVersionProfile(),
			"nutanix_ndb_scale_database":               ndb.ResourceNutanixNDBScaleDatabase(),
			"nutanix_ndb_database_scale":               ndb.ResourceNutanixNDBScaleDatabase(),
			"nutanix_ndb_register_database":            ndb.ResourceNutanixNDBRegisterDatabase(),
			"nutanix_ndb_database_snapshot":            ndb.ResourceNutanixNDBDatabaseSnapshot(),
			"nutanix_ndb_clone":                        ndb.ResourceNutanixNDBClone(),
			"nutanix_ndb_authorize_dbserver":           ndb.ResourceNutanixNDBAuthorizeDBServer(),
			"nutanix_ndb_linked_databases":             ndb.ResourceNutanixNDBLinkedDB(),
			"nutanix_ndb_maintenance_window":           ndb.ResourceNutanixNDBMaintenanceWindow(),
			"nutanix_ndb_maintenance_task":             ndb.ResourceNutanixNDBMaintenanceTask(),
			"nutanix_ndb_tms_cluster":                  ndb.ResourceNutanixNDBTmsCluster(),
			"nutanix_ndb_tag":                          ndb.ResourceNutanixNDBTags(),
			"nutanix_ndb_network":                      ndb.ResourceNutanixNDBNetwork(),
			"nutanix_ndb_dbserver_vm":                  ndb.ResourceNutanixNDBServerVM(),
			"nutanix_ndb_register_dbserver":            ndb.ResourceNutanixNDBRegisterDBServer(),
			"nutanix_ndb_stretched_vlan":               ndb.ResourceNutanixNDBStretchedVlan(),
			"nutanix_ndb_clone_refresh":                ndb.ResourceNutanixNDBCloneRefresh(),
			"nutanix_ndb_cluster":                      ndb.ResourceNutanixNDBCluster(),
			"nutanix_subnet_v4":                        networking.ResourceNutanixSubnetv4(),
			"nutanix_floating_ip_v4":                   networking.ResourceNutanixFloatingIPv4(),
			"nutanix_vpc_v4":                           networking.ResourceNutanixVPCsV4(),
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
