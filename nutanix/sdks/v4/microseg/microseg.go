package microseg

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/api"
	microseg "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	AddressGroupAPIInstance    *api.AddressGroupsServiceApi
	ServiceGroupAPIInstance    *api.ServiceGroupsServiceApi
	NetworkingSecurityInstance *api.NetworkSecurityPoliciesServiceApi
	EntityGroupsAPIInstance    *api.EntityGroupsServiceApi
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
		AddressGroupAPIInstance:    api.NewAddressGroupsServiceApi(baseClient),
		ServiceGroupAPIInstance:    api.NewServiceGroupsServiceApi(baseClient),
		NetworkingSecurityInstance: api.NewNetworkSecurityPoliciesServiceApi(baseClient),
		EntityGroupsAPIInstance:    api.NewEntityGroupsServiceApi(baseClient),
	}

	return f, nil
}
