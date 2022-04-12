package foundation_central

import (
	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	libraryVersion = "v1"
	absolutePath   = "api/fc/" + libraryVersion
	userAgent      = "nutanix/" + libraryVersion
)

// Client manages the foundation central API
type Client struct {
	client *client.Client
	Service
}

// NewFoundationCentralClient return a client to operate foundation central resources
func NewFoundationCentralClient(credentials client.Credentials) (*Client, error) {
	c, err := client.NewClient(&credentials, userAgent, absolutePath, false)

	if err != nil {
		return nil, err
	}

	fc := &Client{
		client: c,
		Service: Operations{
			client: c,
		},
	}

	return fc, nil
}
