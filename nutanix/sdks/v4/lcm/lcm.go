package lcm

import (
	"github.com/nutanix/ntnx-api-golang-clients/lcm-go-client/v4/api"
	lcm "github.com/nutanix/ntnx-api-golang-clients/lcm-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	LcmConfigAPIInstance    *api.ConfigApi
	LcmInventoryAPIInstance *api.InventoryApi
	LcmPreChecksAPIInstance *api.PrecheckApi
	LcmStatusAPIInstance    *api.StatusApi
	LcmEntitiesAPIInstance  *api.EntityApi
}

func NewLcmClient(credentials client.Credentials) (*Client, error) {
	var baseClient *lcm.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := lcm.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		LcmInventoryAPIInstance: api.NewInventoryApi(baseClient),
		LcmConfigAPIInstance:    api.NewConfigApi(baseClient),
		LcmPreChecksAPIInstance: api.NewPrecheckApi(baseClient),
		LcmStatusAPIInstance:    api.NewStatusApi(baseClient),
		LcmEntitiesAPIInstance:  api.NewEntityApi(baseClient),
	}

	return f, nil
}
