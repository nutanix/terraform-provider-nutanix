package v3

import (
	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

//Client manages the V3 API
type Client struct {
	client *client.Client
	V3     Service
}

// NewV3Client return a client to operate V3 resources
func NewV3Client(credentials client.Credentials) (*Client, error) {
	c, err := client.NewClient(&credentials)

	if err != nil {
		return nil, err
	}

	f := &Client{
		client: c,
		V3: Operations{
			client: c,
		},
	}

	return f, nil
}
