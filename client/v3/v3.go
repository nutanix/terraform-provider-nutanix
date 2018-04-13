package v3

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/client/handler"
)

//Client manages the V3 API
type Client struct {
	client *client.Client
	V3     Service
}

// NewV3Client return a client to operate V3 resources
func NewV3Client(config client.Config) (*Client, error) {

	u, err := url.Parse(fmt.Sprintf(client.DefaultBaseURL, config.Credentials.Endpoint, config.Credentials.Port))
	if err != nil {
		return nil, err
	}

	config.BaseURL = u
	config.UserAgent = client.UserAgent
	config.Client = &http.Client{}

	c := client.Client{
		Config:                config,
		MarshalHander:         handler.URLEncodeMarshalHander,
		BuildRequestHandler:   handler.BuildURLEncodedRequest,
		UnmarshalHandler:      handler.UnmarshalXML,
		UnmarshalErrorHandler: handler.UnmarshalErrorHandler,
	}

	f := &Client{
		client: &c,
		V3: Operations{
			client: &c,
		},
	}

	return f, nil
}
