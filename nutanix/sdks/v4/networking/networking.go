package networking

import (
	"github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/api"
	network "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	Routes                      *api.RoutesServiceApi
	RoutesTable                 *api.RouteTablesServiceApi
	APIClientInstance           *network.ApiClient
	RoutingPolicy               *api.RoutingPoliciesServiceApi
	NetworkFunctionAPI          *api.NetworkFunctionsServiceApi
	NicProfilesAPI              *api.NicProfilesServiceApi
	SubnetAPIInstance           *api.SubnetsServiceApi
	VpcAPIInstance              *api.VpcsServiceApi
	FloatingIPAPIInstance       *api.FloatingIpsServiceApi
	VirtualSwitchAPI            *api.VirtualSwitchesServiceApi
	VpcVirtualSwitchMappingsAPI *api.VpcVirtualSwitchMappingsServiceApi
	VirtualSwitchNodesInfoAPI   *api.VirtualSwitchNodesInfoServiceApi
	// BridgesAPI exposes the `$actions/migrate` operation that converts an
	// existing OVS bridge into a Virtual Switch. The standard
	// VirtualSwitchAPI.CreateVirtualSwitch endpoint silently ignores any
	// pre-existing bridge hint, so the migrate API is the only way to bind
	// a Virtual Switch to a specific brN.
	BridgesAPI *api.BridgesServiceApi
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
		Routes:                      api.NewRoutesServiceApi(baseClient),
		RoutesTable:                 api.NewRouteTablesServiceApi(baseClient),
		RoutingPolicy:               api.NewRoutingPoliciesServiceApi(baseClient),
		NetworkFunctionAPI:          api.NewNetworkFunctionsServiceApi(baseClient),
		NicProfilesAPI:              api.NewNicProfilesServiceApi(baseClient),
		SubnetAPIInstance:           api.NewSubnetsServiceApi(baseClient),
		VpcAPIInstance:              api.NewVpcsServiceApi(baseClient),
		FloatingIPAPIInstance:       api.NewFloatingIpsServiceApi(baseClient),
		VirtualSwitchAPI:            api.NewVirtualSwitchesServiceApi(baseClient),
		VpcVirtualSwitchMappingsAPI: api.NewVpcVirtualSwitchMappingsServiceApi(baseClient),
		VirtualSwitchNodesInfoAPI:   api.NewVirtualSwitchNodesInfoServiceApi(baseClient),
		BridgesAPI:                  api.NewBridgesServiceApi(baseClient),
	}

	return f, nil
}
