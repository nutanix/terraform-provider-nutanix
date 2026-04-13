package networking

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/api"
	network "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	Routes                *api.RoutesServiceApi
	RoutesTable           *api.RouteTablesServiceApi
	APIClientInstance     *network.ApiClient
	RoutingPolicy         *api.RoutingPoliciesServiceApi
	SubnetAPIInstance     *api.SubnetsServiceApi
	VpcAPIInstance        *api.VpcsServiceApi
	FloatingIPAPIInstance *api.FloatingIpsServiceApi
}

func NewNetworkingClient(credentials client.Credentials) (*Client, error) {
	var baseClient *network.ApiClient

	pcClient := network.NewApiClient()
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
		Routes:                api.NewRoutesServiceApi(baseClient),
		RoutesTable:           api.NewRouteTablesServiceApi(baseClient),
		RoutingPolicy:         api.NewRoutingPoliciesServiceApi(baseClient),
		SubnetAPIInstance:     api.NewSubnetsServiceApi(baseClient),
		VpcAPIInstance:        api.NewVpcsServiceApi(baseClient),
		FloatingIPAPIInstance: api.NewFloatingIpsServiceApi(baseClient),
	}

	return f, nil
}
