package dataprotection

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v17/api"
	dataprotection "github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	RecoveryPoint     *api.RecoveryPointsServiceApi
	ProtectedResource *api.ProtectedResourcesServiceApi
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

	f := &Client{
		RecoveryPoint:     api.NewRecoveryPointsServiceApi(baseClient),
		ProtectedResource: api.NewProtectedResourcesServiceApi(baseClient),
	}

	return f, nil
}
