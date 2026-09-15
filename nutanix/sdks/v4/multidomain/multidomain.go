package multidomain

import (
	"github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/api"
	multidomainClient "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	Projects          *api.ProjectsServiceApi
	ResourceGroups    *api.ResourceGroupsServiceApi
	APIClientInstance *multidomainClient.ApiClient
}

func NewMultidomainClient(credentials client.Credentials) (*Client, error) {
	var baseClient *multidomainClient.ApiClient

	pcClient := multidomainClient.NewApiClient()
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
		Projects:          api.NewProjectsServiceApi(baseClient),
		ResourceGroups:    api.NewResourceGroupsServiceApi(baseClient),
		APIClientInstance: baseClient,
	}, nil

}
