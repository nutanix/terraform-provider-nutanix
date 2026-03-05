package dataprotection

import (
	"strconv"

	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v17/api"
	dataprotection "github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	RecoveryPoint     *api.RecoveryPointsServiceApi
	ProtectedResource *api.ProtectedResourcesServiceApi
}

func NewDataProtectionClient(credentials client.Credentials) (*Client, error) {
	var baseClient *dataprotection.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := dataprotection.NewApiClient()

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
		RecoveryPoint:     api.NewRecoveryPointsServiceApi(baseClient),
		ProtectedResource: api.NewProtectedResourcesServiceApi(baseClient),
	}

	return f, nil
}
