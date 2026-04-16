package lcm

import (
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/api"
	lcm "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	LcmConfigAPIInstance    *api.ConfigApi
	LcmInventoryAPIInstance *api.InventoryApi
	LcmPreChecksAPIInstance *api.PrechecksApi
	LcmStatusAPIInstance    *api.StatusApi
	LcmEntitiesAPIInstance  *api.EntitiesApi
	LcmUpgradeAPIInstance   *api.UpgradesApi
}

func NewLcmClient(credentials client.Credentials) (*Client, error) {
	var baseClient *lcm.ApiClient

	pcClient := lcm.NewApiClient()
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
		LcmInventoryAPIInstance: api.NewInventoryApi(baseClient),
		LcmConfigAPIInstance:    api.NewConfigApi(baseClient),
		LcmPreChecksAPIInstance: api.NewPrechecksApi(baseClient),
		LcmStatusAPIInstance:    api.NewStatusApi(baseClient),
		LcmEntitiesAPIInstance:  api.NewEntitiesApi(baseClient),
		LcmUpgradeAPIInstance:   api.NewUpgradesApi(baseClient),
	}, nil
}
