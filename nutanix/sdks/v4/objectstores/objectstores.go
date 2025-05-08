package objectstores

import (
	"fmt"
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/api"
	object "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	ObjectStoresAPIInstance *api.ObjectStoresApi
}

func NewObjectStoresClient(credentials client.Credentials) (*Client, error) {
	var baseClient *object.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := object.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		port, err := strconv.Atoi(credentials.Port)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}
		pcClient.Port = port
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		ObjectStoresAPIInstance: api.NewObjectStoresApi(baseClient),
	}

	return f, nil
}
