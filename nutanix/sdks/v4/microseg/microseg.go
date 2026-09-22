package microseg

import (
	"github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/api"
	microseg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	AddressGroupAPIInstance           *api.AddressGroupsServiceApi
	ServiceGroupAPIInstance           *api.ServiceGroupsServiceApi
	NetworkingSecurityInstance        *api.NetworkSecurityPoliciesServiceApi
	EntityGroupsAPIInstance           *api.EntityGroupsServiceApi
	DirectoryServerConfigsAPIInstance *api.DirectoryServerConfigsServiceApi
}

func NewMicrosegClient(credentials client.Credentials) (*Client, error) {
	var baseClient *microseg.ApiClient

	pcClient := microseg.NewApiClient()
	cfg, err := sdkconfig.ConfigureV4Client(credentials, pcClient)
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		pcClient.Host = cfg.Host
		pcClient.Port = cfg.Port
		pcClient.Username = cfg.Username
		pcClient.Password = cfg.Password
		pcClient.VerifySSL = cfg.VerifySSL
		pcClient.AllowVersionNegotiation = cfg.AllowVersionNegotiation
		baseClient = pcClient
	}
	f := &Client{
		AddressGroupAPIInstance:           api.NewAddressGroupsServiceApi(baseClient),
		ServiceGroupAPIInstance:           api.NewServiceGroupsServiceApi(baseClient),
		NetworkingSecurityInstance:        api.NewNetworkSecurityPoliciesServiceApi(baseClient),
		EntityGroupsAPIInstance:           api.NewEntityGroupsServiceApi(baseClient),
		DirectoryServerConfigsAPIInstance: api.NewDirectoryServerConfigsServiceApi(baseClient),
	}

	return f, nil
}
