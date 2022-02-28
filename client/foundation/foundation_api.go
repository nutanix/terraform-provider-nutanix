package foundation

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	absolutePath = "foundation"
	userAgent    = "foundation"
)

//Foundation client with its services
type Client struct {

	//base client
	client *client.Client

	//Service for Imaging Nodes and Cluster Creation
	NodeImaging NodeImagingService

	//Service for File Management in foundation VM
	FileManagement FileManagementService
}

//This routine returns new Foundation API Client
func NewFoundationAPIClient(credentials client.Credentials) (*Client, error) {
	//for foundation client, url should be based on foundation's endpoint and port
	credentials.URL = fmt.Sprintf("%s:%s", credentials.FoundationEndpoint, credentials.FoundationPort)
	client, err := client.NewBaseClient(&credentials, absolutePath, true)

	if err != nil {
		return nil, err
	}

	//Fill user agent details
	client.UserAgent = userAgent

	foundationClient := &Client{
		client: client,
		NodeImaging: NodeImagingOperations{
			client: client,
		},
		FileManagement: FileManagementOperations{
			client: client,
		},
	}
	return foundationClient, nil
}
