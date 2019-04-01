package nutanix

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider function returns the object that implements the terraform.ResourceProvider interface, specifically a schema.Provider
func Provider() terraform.ResourceProvider {

	// Nutanix provider schema
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_USERNAME", nil),
				Description: descriptions["username"],
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PASSWORD", nil),
				Description: descriptions["password"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_INSECURE", false),
				Description: descriptions["insecure"],
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
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_ENDPOINT", nil),
				Description: descriptions["endpoint"],
			},
			"wait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_WAIT_TIMEOUT", nil),
				Description: descriptions["wait_timeout"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nutanix_image":           dataSourceNutanixImage(),
			"nutanix_subnet":          dataSourceNutanixSubnet(),
			"nutanix_cluster":         dataSourceNutanixCluster(),
			"nutanix_clusters":        dataSourceNutanixClusters(),
			"nutanix_virtual_machine": dataSourceNutanixVirtualMachine(),
			// "nutanix_network_security_rule": dataSourceNutanixNetworkSecurityRule(),
			// "nutanix_volume_group":           dataSourceNutanixVolumeGroup(),
			// "nutanix_volume_groups":          dataSourceNutanixVolumeGroups(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nutanix_virtual_machine": resourceNutanixVirtualMachine(),
			"nutanix_image":           resourceNutanixImage(),
			"nutanix_subnet":          resourceNutanixSubnet(),
			"nutanix_category_key":    resourceNutanixCategoryKey(),
			"nutanix_category_value":  resourceNutanixCategoryValue(),
			// "nutanix_network_security_rule": resourceNutanixNetworkSecurityRule(),
			// "nutanix_volume_group":          resourceNutanixVolumeGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// defines descriptions for ResourceProvider schema definitions
var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"username": "User name for Nutanix Prism. Could be\n" +
			"local cluster auth (e.g. 'admin') or directory auth.",

		"password": "Password for provided user name.",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",

		"port": "Port for Nutanix Prism.",

		"endpoint": "URL for Nutanix Prism (e.g IP or FQDN for cluster VIP\n" +
			"note, this is never the data services VIP, and should not be an\n" +
			"individual CVM address, as this would cause calls to fail during\n" +
			"cluster lifecycle management operations, such as AOS upgrades.",

		"wait_timeout": "Set if you know that the creation o update of a resource may take long time (minutes)",
	}
}

// This function used to fetch the configuration params given to our provider which
// we will use to initialize a dummy client that interacts with API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[DEBUG] config wait_timeout %d", d.Get("wait_timeout").(int))

	config := Config{
		Endpoint:    d.Get("endpoint").(string),
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		Insecure:    d.Get("insecure").(bool),
		Port:        d.Get("port").(string),
		WaitTimeout: int64(d.Get("wait_timeout").(int)),
	}

	return config.Client()
}
