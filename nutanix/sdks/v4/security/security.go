package security

import (
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/api"
	security "github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	KeyManagementServersAPIInstance *api.KeyManagementServersApi
	STIGsAPI                        *api.STIGsApi
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

	return &Client{
		KeyManagementServersAPIInstance: api.NewKeyManagementServersApi(baseClient),
		STIGsAPI:                        api.NewSTIGsApi(baseClient),
	}, nil
}
