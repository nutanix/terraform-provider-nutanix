package lifecycle

import (
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/api"
	lifecycle "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	FoundationCentralConfigAPIInstance *api.FoundationCentralConfigApi
	HardwareProvidersAPIInstance       *api.HardwareProvidersApi
	APIClientInstance                  *lifecycle.ApiClient
}

func NewLifecycleClient(credentials client.Credentials) (*Client, error) {
	var baseClient *lifecycle.ApiClient

	pcClient := lifecycle.NewApiClient()
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
		FoundationCentralConfigAPIInstance: api.NewFoundationCentralConfigApi(baseClient),
		HardwareProvidersAPIInstance:       api.NewHardwareProvidersApi(baseClient),
		APIClientInstance:                  baseClient,
	}, nil
}
