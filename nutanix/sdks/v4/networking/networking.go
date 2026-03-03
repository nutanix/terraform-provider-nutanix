package networking

import (
	"strconv"

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
		Routes:                api.NewRoutesServiceApi(baseClient),
		RoutesTable:           api.NewRouteTablesServiceApi(baseClient),
		RoutingPolicy:         api.NewRoutingPoliciesServiceApi(baseClient),
		SubnetAPIInstance:     api.NewSubnetsServiceApi(baseClient),
		VpcAPIInstance:        api.NewVpcsServiceApi(baseClient),
		FloatingIPAPIInstance: api.NewFloatingIpsServiceApi(baseClient),
	}

	return f, nil
}
