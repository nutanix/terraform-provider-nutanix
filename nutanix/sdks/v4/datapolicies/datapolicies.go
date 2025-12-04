package datapolicies

import (
	"github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/api"
	datapolicies "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	ProtectionPolicies *api.ProtectionPoliciesApi
}

func NewDataPoliciesClient(credentials client.Credentials) (*Client, error) {
	var baseClient *datapolicies.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := datapolicies.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		ProtectionPolicies: api.NewProtectionPoliciesApi(baseClient),
	}

	return f, nil
}
