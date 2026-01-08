package datapolicies

import (
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/api"
	datapolicies "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ProtectionPolicies *api.ProtectionPoliciesApi
	StoragePolicies    *api.StoragePoliciesApi
}

func NewDataPoliciesClient(credentials client.Credentials) (*Client, error) {
	var baseClient *datapolicies.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := datapolicies.NewApiClient()

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
		ProtectionPolicies: api.NewProtectionPoliciesApi(baseClient),
		StoragePolicies:    api.NewStoragePoliciesApi(baseClient),
	}

	return f, nil
}
