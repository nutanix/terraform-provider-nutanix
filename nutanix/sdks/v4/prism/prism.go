package prism

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/api"
	prism "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	TaskRefAPI                      *api.TasksServiceApi
	CategoriesAPIInstance           *api.CategoriesServiceApi
	DomainManagerAPIInstance        *api.DomainManagerServiceApi
	DomainManagerBackupsAPIInstance *api.DomainManagerBackupsServiceApi
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

	f := &Client{
		TaskRefAPI:                      api.NewTasksServiceApi(baseClient),
		CategoriesAPIInstance:           api.NewCategoriesServiceApi(baseClient),
		DomainManagerAPIInstance:        api.NewDomainManagerServiceApi(baseClient),
		DomainManagerBackupsAPIInstance: api.NewDomainManagerBackupsServiceApi(baseClient),
	}

	return f, nil
}
