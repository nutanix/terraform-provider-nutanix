package nutanix

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// Config ...
type Config struct {
	Endpoint string
	Username string
	Password string
	Port     string
	Insecure bool
}

// Client ...
func (c *Config) Client() (*OutscaleClient, error) {

	config := client.Config{
		Credentials: client.Credentials{
			Endpoint: c.Endpoint,
			Username: c.Username,
			Password: c.Password,
			Port: c.Port,
			Insecure: c.Insecure,
			URL: fmt.Sprintf(client.DefaultBaseURL, c.Endpoint, c.Port)
		}
	}

	v3, err := v3.NewV3Client(config)
	if err != nil {
		return nil, err
	}
	client := &NutanixClient{
		API: v3,
	}

	return client, nil
}

// Nutanix client
type NutanixClient struct {
	API *client.Client
}
