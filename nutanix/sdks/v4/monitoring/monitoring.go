package monitoring

import (
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/api"
	monitoringClient "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	AuditsAPI        *api.AuditsApi
	APIClientInstance *monitoringClient.ApiClient
}

func NewMonitoringClient(credentials client.Credentials) (*Client, error) {
	var baseClient *monitoringClient.ApiClient

	pcClient := monitoringClient.NewApiClient()
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
		AuditsAPI:        api.NewAuditsApi(baseClient),
		APIClientInstance: baseClient,
	}, nil
}
