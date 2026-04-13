package datapolicies

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/datapolicies-go-client/v17/api"
	datapolicies "github.com/nutanix-core/ntnx-api-golang-sdk-internal/datapolicies-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ProtectionPolicies *api.ProtectionPoliciesServiceApi
	StoragePolicies    *api.StoragePoliciesServiceApi
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

	f := &Client{
		ProtectionPolicies: api.NewProtectionPoliciesServiceApi(baseClient),
		StoragePolicies:    api.NewStoragePoliciesServiceApi(baseClient),
	}

	return f, nil
}
