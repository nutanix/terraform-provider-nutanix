package lifecycle

import (
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/api"
	lifecycle "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

// Client is the SDK wrapper for the Nutanix lifecycle (Foundation Central) v4 APIs.
type Client struct {
	NodesAPIInstance  *api.NodesApi
	APIClientInstance *lifecycle.ApiClient
}

// NewLifecycleClient builds a lifecycle API client from the provided credentials.
func NewLifecycleClient(credentials client.Credentials) (*Client, error) {
	var baseClient *lifecycle.ApiClient

	pcClient := lifecycle.NewApiClient()
	if cfg := sdkconfig.ConfigureV4Client(credentials, pcClient); cfg != nil {
		pcClient.Host = cfg.Host
		pcClient.Port = cfg.Port
		pcClient.Username = cfg.Username
		pcClient.Password = cfg.Password
		pcClient.VerifySSL = cfg.VerifySSL
		pcClient.AllowVersionNegotiation = cfg.AllowVersionNegotiation
		baseClient = pcClient
	}

	return &Client{
		NodesAPIInstance:  api.NewNodesApi(baseClient),
		APIClientInstance: baseClient,
	}, nil
}
