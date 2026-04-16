package microseg

import (
	"github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/api"
	microseg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	AddressGroupAPIInstance    *api.AddressGroupsApi
	ServiceGroupAPIInstance    *api.ServiceGroupsApi
	NetworkingSecurityInstance *api.NetworkSecurityPoliciesApi
	EntityGroupsAPIInstance    *api.EntityGroupsApi
}

func NewMicrosegClient(credentials client.Credentials) (*Client, error) {
	var baseClient *microseg.ApiClient

	pcClient := microseg.NewApiClient()
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
		AddressGroupAPIInstance:    api.NewAddressGroupsApi(baseClient),
		ServiceGroupAPIInstance:    api.NewServiceGroupsApi(baseClient),
		NetworkingSecurityInstance: api.NewNetworkSecurityPoliciesApi(baseClient),
		EntityGroupsAPIInstance:    api.NewEntityGroupsApi(baseClient),
	}

	return f, nil
}
