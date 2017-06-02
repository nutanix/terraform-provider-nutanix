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
}

// Machine struct is for defining virtual machine
type Machine st.VirtualMachine

// Provider function returns the object that implements the terraform.ResourceProvider interface, specifically a schema.Provider
func Provider() terraform.ResourceProvider {

	// Nutanix provider schema
	return &schema.Provider{
		Schema: providerSchema(),
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
			Description: "Username for authentiaction",
		},
		"password": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Password for authentiaction",
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Endpoint for API call",
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
	}

	return &client, nil
}
