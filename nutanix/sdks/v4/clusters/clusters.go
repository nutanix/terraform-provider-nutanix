package clusters

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v16/api"
	network "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v16/client"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	ClusterEntityAPI     *api.ClustersApi
	StorageContainersAPI *api.StorageContainersApi
}

func NewClustersClient(credentials client.Credentials) (*Client, error) {
	var baseClient *network.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := network.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		ClusterEntityAPI:     api.NewClustersApi(baseClient),
		StorageContainersAPI: api.NewStorageContainersApi(baseClient),
	}

	return f, nil
}
