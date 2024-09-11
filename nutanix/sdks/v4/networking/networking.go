package networking

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/api"
	network "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/client"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	RoutesTable           *api.RouteTablesApi
	APIClientInstance     *network.ApiClient
	RoutingPolicy         *api.RoutingPoliciesApi
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
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		RoutesTable:           api.NewRouteTablesApi(baseClient),
		RoutingPolicy:         api.NewRoutingPoliciesApi(baseClient),
		SubnetAPIInstance:     api.NewSubnetsApi(baseClient),
		VpcAPIInstance:        api.NewVpcsApi(baseClient),
		FloatingIPAPIInstance: api.NewFloatingIpsApi(baseClient),
	}

	return f, nil
}
