package nutanix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	st "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachinestruct"
)

// MyClient struct is for defining provider
type MyClient struct {
	Endpoint string
	Username string
	Password string
	Insecure bool
}

// Machine struct is for defining virtual machine
type Machine st.VirtualMachine

// Provider function returns the object that implements the terraform.ResourceProvider interface, specifically a schema.Provider
func Provider() terraform.ResourceProvider {

	// Nutanix provider schema
	return &schema.Provider{
		Schema:         providerSchema(),
		DataSourcesMap: map[string]*schema.Resource{
		//"nutanix_image": dataSourceImage(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nutanix_virtual_machine": resourceNutanixVirtualMachine(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// defines descriptions for ResourceProvider schema definitions
var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"user": "User name for Nutanix Prism Element. Could be\n" +
			"local cluster auth (e.g. 'admin') or directory auth.",

		"password": "Password for provided user name.",

		"endpoint": "URL for Nutanix Prism Element (e.g IP or FQDN for cluster VIP\n" +
			"note, this is never the data services VIP, and should not be an\n" +
			"individual CVM address, as this would cause calls to fail during\n" +
			"cluster lifecycle management operations, such as AOS upgrades.",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",
	}
}

// List of supported configuration fields for your provider
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_USERNAME", nil),
			Description: descriptions["username"],
		},
		"password": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PASSWORD", nil),
			Description: descriptions["password"],
		},
		"insecure": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_INSECURE", false),
			Description: descriptions["insecure"],
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_ENDPOINT", nil),
			Description: descriptions["endpoint"],
		},
	}

}

// This function used to fetch the configuration params given to our provider which
// we will use to initialize a dummy client that interacts with API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := MyClient{
		Endpoint: d.Get("endpoint").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Insecure: d.Get("insecure").(bool),
	}

	return &client, nil
}
