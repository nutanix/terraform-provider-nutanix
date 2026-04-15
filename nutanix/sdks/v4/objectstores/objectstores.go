package objectstores

import (
	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/api"
	objects "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ObjectStoresAPIInstance *api.ObjectStoresApi
}

func NewObjectStoresClient(credentials client.Credentials) (*Client, error) {
	var baseClient *objects.ApiClient

	pcClient := objects.NewApiClient()
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
		ObjectStoresAPIInstance: api.NewObjectStoresApi(baseClient),
	}, nil
}
