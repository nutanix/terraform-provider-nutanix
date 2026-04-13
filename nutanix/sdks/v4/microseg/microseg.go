package microseg

import (
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/api"
	microseg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	AddressGroupAPIInstance    *api.AddressGroupsApi
	ServiceGroupAPIInstance    *api.ServiceGroupsApi
	NetworkingSecurityInstance *api.NetworkSecurityPoliciesApi
}

func NewMicrosegClient(credentials client.Credentials) (*Client, error) {
	var baseClient *microseg.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := microseg.NewApiClient()

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
		AddressGroupAPIInstance:    api.NewAddressGroupsApi(baseClient),
		ServiceGroupAPIInstance:    api.NewServiceGroupsApi(baseClient),
		NetworkingSecurityInstance: api.NewNetworkSecurityPoliciesApi(baseClient),
	}

	return f, nil
}
