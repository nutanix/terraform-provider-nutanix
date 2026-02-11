package networking

import (
	"strconv"

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

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := network.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = sdkconfig.DefaultPort
		if credentials.Port != "" {
			if p, err := strconv.Atoi(credentials.Port); err == nil {
				pcClient.Port = p
			}
		}
		pcClient.VerifySSL = false
		pcClient.AllowVersionNegotiation = sdkconfig.AllowVersionNegotiation
		baseClient = pcClient
	}

	f := &Client{
		Routes:                api.NewRoutesApi(baseClient),
		RoutesTable:           api.NewRouteTablesApi(baseClient),
		RoutingPolicy:         api.NewRoutingPoliciesApi(baseClient),
		NetworkFunctionAPI:    api.NewNetworkFunctionsApi(baseClient),
		SubnetAPIInstance:     api.NewSubnetsApi(baseClient),
		VpcAPIInstance:        api.NewVpcsApi(baseClient),
		FloatingIPAPIInstance: api.NewFloatingIpsApi(baseClient),
	}

	return f, nil
}
