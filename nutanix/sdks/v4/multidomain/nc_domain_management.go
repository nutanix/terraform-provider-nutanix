package multidomain

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/api"
	multidomainClient "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type NCDomainManagementClient struct {
	DomainServiceAPIInstance    *api.DomainServiceApi
	APIClientInstance *multidomainClient.ApiClient
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
		APIClientInstance:        baseClient,
	}

	return f, nil
}
