// Package licensing provides a client for interacting with the Nutanix licensing API.
package licensing

import (
	"fmt"
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/api"
	licensing "github.com/nutanix/ntnx-api-golang-clients/licensing-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	LicensesAPIInstance                *api.LicensesApi
	EndUserLicenseAgreementAPIInstance *api.EndUserLicenseAgreementApi
}

func NewLicensingClient(credentials client.Credentials) (*Client, error) {
	var baseClient *licensing.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := licensing.NewApiClient()

		port, err := strconv.Atoi(credentials.Port)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = port
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		LicensesAPIInstance:                api.NewLicensesApi(baseClient),
		EndUserLicenseAgreementAPIInstance: api.NewEndUserLicenseAgreementApi(baseClient),
	}

	return f, nil
}
