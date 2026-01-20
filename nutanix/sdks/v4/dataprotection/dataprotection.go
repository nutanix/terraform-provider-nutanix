package dataprotection

import (
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/api"
	dataprotection "github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	RecoveryPoint     *api.RecoveryPointsApi
	ProtectedResource *api.ProtectedResourcesApi
}

func NewDataProtectionClient(credentials client.Credentials) (*Client, error) {
	var baseClient *dataprotection.ApiClient

	pcClient := dataprotection.NewApiClient()
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
		RecoveryPoint:     api.NewRecoveryPointsApi(baseClient),
		ProtectedResource: api.NewProtectedResourcesApi(baseClient),
	}, nil
}
