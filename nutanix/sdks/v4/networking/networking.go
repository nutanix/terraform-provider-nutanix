package networking

import (
	"github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/api"
	network "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	Routes                *api.RoutesApi
	RoutesTable           *api.RouteTablesApi
	APIClientInstance     *network.ApiClient
	RoutingPolicy         *api.RoutingPoliciesApi
	NetworkFunctionAPI    *api.NetworkFunctionsApi
	SubnetAPIInstance     *api.SubnetsApi
	VpcAPIInstance        *api.VpcsApi
	FloatingIPAPIInstance *api.FloatingIpsApi
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

	return &Client{
		Routes:                api.NewRoutesApi(baseClient),
		RoutesTable:           api.NewRouteTablesApi(baseClient),
		RoutingPolicy:         api.NewRoutingPoliciesApi(baseClient),
		NetworkFunctionAPI:    api.NewNetworkFunctionsApi(baseClient),
		SubnetAPIInstance:     api.NewSubnetsApi(baseClient),
		VpcAPIInstance:        api.NewVpcsApi(baseClient),
		FloatingIPAPIInstance: api.NewFloatingIpsApi(baseClient),
	}, nil
}
