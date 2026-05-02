package monitoring

import (
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/api"
	monitoring "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	SystemDefinedChecksAPI  *api.SystemDefinedChecksApi
	SystemDefinedPoliciesAPI *api.SystemDefinedPoliciesApi
	APIClientInstance        *monitoring.ApiClient
}

func NewMonitoringClient(credentials client.Credentials) (*Client, error) {
	var baseClient *monitoring.ApiClient

	pcClient := monitoring.NewApiClient()
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
		SystemDefinedChecksAPI:  api.NewSystemDefinedChecksApi(baseClient),
		SystemDefinedPoliciesAPI: api.NewSystemDefinedPoliciesApi(baseClient),
		APIClientInstance:        baseClient,
	}, nil
}
