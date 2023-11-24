package networking

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/api"
	network "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type NetworkingClient struct {
	SubnetApiInstance     *api.SubnetApi
	VpcApiInstance        *api.VpcApi
	FloatingIpApiInstance *api.FloatingIpApi
}

func NewNetworkingClient(credentials client.Credentials) (*NetworkingClient, error) {
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

	f := &NetworkingClient{
		SubnetApiInstance:     api.NewSubnetApi(baseClient),
		VpcApiInstance:        api.NewVpcApi(baseClient),
		FloatingIpApiInstance: api.NewFloatingIpApi(baseClient),
	}

	return f, nil
}
