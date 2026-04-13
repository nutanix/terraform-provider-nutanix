package karbon

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

const (
	absolutePath = "karbon"
	userAgent    = "nutanix"
)

// Client manages the V3 API
type Client struct {
	client          *client.Client
	Cluster         ClusterService
	PrivateRegistry PrivateRegistryService
	Meta            MetaService
}

// NewKarbonAPIClient return a client to operate Karbon resources
func NewKarbonAPIClient(credentials client.Credentials) (*Client, error) {
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
		errorMsg := fmt.Sprintf("Karbon Client is missing. " +
			"Please provide endpoint and either (username + password) or api_key in provider configuration.")
		baseClient = &client.Client{UserAgent: userAgent, ErrorMsg: errorMsg}
	}

	f := &Client{
		client: baseClient,
		Cluster: ClusterOperations{
			client: baseClient,
		},
		PrivateRegistry: PrivateRegistryOperations{
			client: baseClient,
		},
		Meta: MetaOperations{
			client: baseClient,
		},
	}

	return f, nil
}
