package prism

import (
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/api"
	prism "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	TaskRefAPI                      *api.TasksApi
	CategoriesAPIInstance           *api.CategoriesApi
	DomainManagerAPIInstance        *api.DomainManagerApi
	DomainManagerBackupsAPIInstance *api.DomainManagerBackupsApi
}

func NewPrismClient(credentials client.Credentials) (*Client, error) {
	var baseClient *prism.ApiClient

	pcClient := prism.NewApiClient()
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
		TaskRefAPI:                      api.NewTasksApi(baseClient),
		CategoriesAPIInstance:           api.NewCategoriesApi(baseClient),
		DomainManagerAPIInstance:        api.NewDomainManagerApi(baseClient),
		DomainManagerBackupsAPIInstance: api.NewDomainManagerBackupsApi(baseClient),
	}, nil
}
