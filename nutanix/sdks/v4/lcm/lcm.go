package lcm

import (
	"strconv"

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

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := lcm.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = sdkconfig.DefaultPort
		if credentials.Port != "" {
			if p, err := strconv.Atoi(credentials.Port); err == nil {
				pcClient.Port = p
			}
		}
		pcClient.VerifySSL = false
		pcClient.AllowVersionNegotiation = sdkconfig.AllowVersionNegotiation
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
