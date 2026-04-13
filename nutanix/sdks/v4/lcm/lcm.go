package lcm

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/lifecycle-go-client/v17/api"
	lcm "github.com/nutanix-core/ntnx-api-golang-sdk-internal/lifecycle-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	LcmConfigAPIInstance    *api.ConfigServiceApi
	LcmInventoryAPIInstance *api.InventoryServiceApi
	LcmPreChecksAPIInstance *api.PrechecksServiceApi
	LcmStatusAPIInstance    *api.StatusServiceApi
	LcmEntitiesAPIInstance  *api.EntitiesServiceApi
	LcmUpgradeAPIInstance   *api.UpgradesServiceApi
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

	f := &Client{
		LcmInventoryAPIInstance: api.NewInventoryServiceApi(baseClient),
		LcmConfigAPIInstance:    api.NewConfigServiceApi(baseClient),
		LcmPreChecksAPIInstance: api.NewPrechecksServiceApi(baseClient),
		LcmStatusAPIInstance:    api.NewStatusServiceApi(baseClient),
		LcmEntitiesAPIInstance:  api.NewEntitiesServiceApi(baseClient),
		LcmUpgradeAPIInstance:   api.NewUpgradesServiceApi(baseClient),
	}
	return f, nil
}
