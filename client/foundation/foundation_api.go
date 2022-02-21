package foundation

import (
	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	absolutePath = "foundation"
	//To-Do: if its required
	userAgent = "foundation"
	mediaType = "application/json"
)

//Foundation client with its services
type Client struct {

	//base client
	client *client.Client

	//Service for Imaging and Cluster Creation
	ImageNodes ImageNodesService

	//Service for File Management in foundation VM
	FileManagement FileManagementService
}

//This routine returns new Foundation API Client
func NewFoundationAPIClient(credentials client.Credentials) (*Client, error) {

	client, err := client.NewClient(&credentials, userAgent, absolutePath)

	if err != nil {
		return nil, err
	}

	foundationClient := &Client{
		client: client,
		ImageNodes: ImageNodesOperations{
			client: client,
		},
		FileManagement: FileManagementOperations{
			client: client,
		},
	}
	return foundationClient, nil
}
