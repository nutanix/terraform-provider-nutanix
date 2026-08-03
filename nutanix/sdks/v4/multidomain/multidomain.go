package multidomain

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/api"
	multidomainClient "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/client"
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

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := multidomainClient.NewApiClient()

		if cfg := sdkconfig.ConfigureV4Client(credentials, pcClient); cfg != nil {
			pcClient.Host = cfg.Host
			pcClient.Port = cfg.Port
			pcClient.Username = cfg.Username
			pcClient.Password = cfg.Password
			pcClient.VerifySSL = cfg.VerifySSL
			pcClient.AllowVersionNegotiation = cfg.AllowVersionNegotiation
		}
		pcClient.VerifySSL = false
		pcClient.AllowVersionNegotiation = sdkconfig.AllowVersionNegotiation
		baseClient = pcClient
	}

	f := &Client{
		Projects:          api.NewProjectsServiceApi(baseClient),
		ResourceGroups:    api.NewResourceGroupsServiceApi(baseClient),
		APIClientInstance: baseClient,
	}

	return f, nil
}
