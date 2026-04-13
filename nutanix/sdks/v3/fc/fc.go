package foundationcentral

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

const (
	libraryVersion = "v1"
	absolutePath   = "api/fc/" + libraryVersion
	userAgent      = "nutanix/" + libraryVersion
)

// Client manages the foundation central API
type Client struct {
	client  *client.Client
	Service Service
}

// NewFoundationCentralClient return a client to operate foundation central resources
func NewFoundationCentralClient(credentials client.Credentials) (*Client, error) {
	var baseClient *client.Client

	// check if all required fields are present. Else create an empty client
	// Accept either (username + password) OR api_key for authentication
	hasBasicAuth := credentials.Username != "" && credentials.Password != ""
	hasAPIKey := credentials.APIKey != ""
	hasEndpoint := credentials.Endpoint != ""

	if hasEndpoint && (hasBasicAuth || hasAPIKey) {
		c, err := client.NewClient(&credentials, userAgent, absolutePath, false)
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		errorMsg := fmt.Sprintf("Foundation Central Client is missing. " +
			"Please provide endpoint and either (username + password) or api_key in provider configuration.")

		baseClient = &client.Client{UserAgent: userAgent, ErrorMsg: errorMsg}
	}

	fc := &Client{
		client: baseClient,
		Service: Operations{
			client: baseClient,
		},
	}
	return fc, nil
}
