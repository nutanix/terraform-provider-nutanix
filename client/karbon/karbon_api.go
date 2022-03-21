package karbon

import (
	"github.com/terraform-providers/terraform-provider-nutanix/client"
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
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		c, err := client.NewClient(&credentials, userAgent, absolutePath, false)
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		baseClient = &client.Client{UserAgent: userAgent}
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
