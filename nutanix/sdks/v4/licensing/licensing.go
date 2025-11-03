package licensing

import (
	"github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/api"
	licensing "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	LicensingEULAAPIInstance *api.EndUserLicenseAgreementApi
	LicensesAPIInstance      *api.LicensesApi
	LicenseKeysAPIInstance   *api.LicenseKeysApi
}

func NewLicensingClient(credentials client.Credentials) (*Client, error) {
	var baseClient *licensing.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := licensing.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		LicensingEULAAPIInstance: api.NewEndUserLicenseAgreementApi(baseClient),
		LicensesAPIInstance:      api.NewLicensesApi(baseClient),
		LicenseKeysAPIInstance:   api.NewLicenseKeysApi(baseClient),
	}
	return f, nil
}
