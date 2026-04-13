package security

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/security-go-client/v17/api"
	security "github.com/nutanix-core/ntnx-api-golang-sdk-internal/security-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	KeyManagementServersAPIInstance *api.KeyManagementServersServiceApi
	STIGsAPI                        *api.STIGsServiceApi
}

func NewSecurityClient(credentials client.Credentials) (*Client, error) {
	var baseClient *security.ApiClient

	pcClient := security.NewApiClient()
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
		KeyManagementServersAPIInstance: api.NewKeyManagementServersServiceApi(baseClient),
		STIGsAPI:                        api.NewSTIGsServiceApi(baseClient),
	}

	return f, nil
}
