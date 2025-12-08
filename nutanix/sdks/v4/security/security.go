package security

import (
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/api"
	prism "github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	KeyManagementServersAPIInstance *api.KeyManagementServersApi
	STIGsAPI                        *api.STIGsApi
}

func NewSecurityClient(credentials client.Credentials) (*Client, error) {
	var baseClient *prism.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := prism.NewApiClient()

		port, err := strconv.Atoi(credentials.Port)
		if err != nil {
			pcClient.Port = 9440
		}
		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = port
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		KeyManagementServersAPIInstance: api.NewKeyManagementServersApi(baseClient),
		STIGsAPI:                        api.NewSTIGsApi(baseClient),
	}

	return f, nil
}
