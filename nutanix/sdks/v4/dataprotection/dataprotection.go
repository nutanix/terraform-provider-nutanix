package dataprotection

import (
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/api"
	dataprotection "github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	APIClientInstance *dataprotection.ApiClient
	RecoveryPoint     *api.RecoveryPointsApi
}

func NewDataProtectionClient(credentials client.Credentials) (*Client, error) {
	var baseClient *dataprotection.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := dataprotection.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		RecoveryPoint: api.NewRecoveryPointsApi(baseClient),
	}

	return f, nil
}
