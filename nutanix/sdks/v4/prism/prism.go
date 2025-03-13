package prism

import (
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/api"
	prism "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	TaskRefAPI               *api.TasksApi
	CategoriesAPIInstance    *api.CategoriesApi
	DomainManagerAPIInstance *api.DomainManagerApi
}

func NewPrismClient(credentials client.Credentials) (*Client, error) {
	var baseClient *prism.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := prism.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		TaskRefAPI:               api.NewTasksApi(baseClient),
		CategoriesAPIInstance:    api.NewCategoriesApi(baseClient),
		DomainManagerAPIInstance: api.NewDomainManagerApi(baseClient),
	}

	return f, nil
}
