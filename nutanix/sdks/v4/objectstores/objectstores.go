package objectstores

import (
	"strconv"

	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/api"
	object "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ObjectStoresAPIInstance *api.ObjectStoresServiceApi
}

func NewObjectStoresClient(credentials client.Credentials) (*Client, error) {
	var baseClient *object.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := object.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = sdkconfig.DefaultPort
		if credentials.Port != "" {
			if p, err := strconv.Atoi(credentials.Port); err == nil {
				pcClient.Port = p
			}
		}
		pcClient.VerifySSL = false
		pcClient.AllowVersionNegotiation = sdkconfig.AllowVersionNegotiation
		baseClient = pcClient
	}

	f := &Client{
		ObjectStoresAPIInstance: api.NewObjectStoresServiceApi(baseClient),
	}

	return f, nil
}
