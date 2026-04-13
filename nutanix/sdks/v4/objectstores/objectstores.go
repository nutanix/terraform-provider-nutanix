package objectstores

import (
	objects "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/client"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/api"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ObjectStoresAPIInstance *api.ObjectStoresServiceApi
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

	f := &Client{
		ObjectStoresAPIInstance: api.NewObjectStoresServiceApi(baseClient),
	}

	return f, nil
}
