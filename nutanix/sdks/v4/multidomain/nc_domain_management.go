package multidomain

import (
	"github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/api"
	multidomainClient "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type NCDomainManagementClient struct {
	DomainServiceAPIInstance *api.DomainServiceApi
	Locations                *api.LocationsServiceApi
	RegisteredDomains        *api.RegisteredDomainsServiceApi
	APIClientInstance        *multidomainClient.ApiClient
	Tasks                    *NCTaskClient
}

func NewNCDomainManagementClient(credentials client.Credentials) (*NCDomainManagementClient, error) {
	var baseClient *multidomainClient.ApiClient

	pcClient := multidomainClient.NewApiClient()
	if cfg := sdkconfig.ConfigureV4ClientForDomainManagement(credentials, pcClient); cfg != nil {
		pcClient.Host = cfg.Host
		pcClient.Port = cfg.Port
		pcClient.Username = cfg.Username
		pcClient.Password = cfg.Password
		pcClient.VerifySSL = cfg.VerifySSL
		pcClient.AllowVersionNegotiation = cfg.AllowVersionNegotiation
		baseClient = pcClient
	}

	f := &NCDomainManagementClient{
		DomainServiceAPIInstance: api.NewDomainServiceApi(baseClient),
		Locations:                api.NewLocationsServiceApi(baseClient),
		RegisteredDomains:        api.NewRegisteredDomainsServiceApi(baseClient),
		APIClientInstance:        baseClient,
		Tasks:                    NewNCTaskClient(baseClient, credentials),
	}
	return f, nil
}
