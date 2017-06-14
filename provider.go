package nutanix

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// NutanixV3Client struct is for defining provider
type NutanixV3Client struct {
	Endpoint string
	Username string
	Password string
	Port     string
	URL      string
	Insecure bool
}

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

// List of supported configuration fields for your provider
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_USERNAME", nil),
			Description: "User name for Nutanix Prism Element. Could be local cluster auth (e.g. 'admin') or directory auth.",
		},
		"password": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PASSWORD", nil),
			Description: "Password for provided user name.",
		},
		"insecure": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_INSECURE", false),
			Description: "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted default value is `false`",
		},
		"port": &schema.Schema{
			Type:        schema.TypeString,
			Default:     "9440",
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PORT", false),
			Description: "Port for Nutanix Prism Element",
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("NUTANIX_ENDPOINT", nil),
			Description: "IP address for Nutanix Prism Element",
		},
	}

}

// This function used to fetch the configuration params given to our provider which
// we will use to initialize a dummy client that interacts with API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	port := d.Get("port").(string)
	endpoint := d.Get("endpoint").(string)
	url := "https://" + endpoint + ":" + port + "/api/nutanix/v3"
	client := NutanixV3Client{
		Endpoint: endpoint,
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Insecure: d.Get("insecure").(bool),
		Port:     port,
		URL:      url,
	}

	return &client, nil
}
