package datapolicies

import (
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

	pcClient := datapolicies.NewApiClient()
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
		ProtectionPolicies: api.NewProtectionPoliciesApi(baseClient),
		StoragePolicies:    api.NewStoragePoliciesApi(baseClient),
	}, nil
}
