package Era

import (
	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	libraryVersion = "v0.9"
	absolutePath   = "era/" + libraryVersion
	//userAgent      = "nutanix/" + libraryVersion // Check whether user-agent will be same as that of nutanix.
)

type Client struct {
	client  *client.Client
	Service Service
}

func NewV3Client(credentials client.Credentials) (*Client, error) {
	c, err := client.NewClient(&credentials, "", absolutePath, false)

	if err != nil {
		return nil, err
	}

	f := &Client{
		client: c,
		Service: ServiceClient{
			c: c,
		},
	}

	return f, nil
}
