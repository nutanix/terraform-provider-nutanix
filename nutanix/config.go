package nutanix

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

//Version represents api version
const Version = "3.0"

// Config ...
type Config struct {
	Endpoint string
	Username string
	Password string
	Port     string
	Insecure bool
}

// Client ...
func (c *Config) Client() (*NutanixClient, error) {

	configCreds := client.Credentials{
		URL:      fmt.Sprintf("%s:%s", c.Endpoint, c.Port),
		Endpoint: c.Endpoint,
		Username: c.Username,
		Password: c.Password,
		Port:     c.Port,
		Insecure: c.Insecure,
	}

	v3, err := v3.NewV3Client(configCreds)
	if err != nil {
		return nil, err
	}
	client := &NutanixClient{
		API: v3,
	}

	return client, nil
}

//NutanixClient represents the nutanix API client
type NutanixClient struct {
	API *v3.Client
}
