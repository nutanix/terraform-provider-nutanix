package objects

import(
	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/api"
	object "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	APIClientInstance *= *object.ApiClient
}

func NewObjectsClient(credentials client.Credentials) (*Client, error) {
	var baseClient *object.ApiClient

	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := object.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		APIClientInstance: object.NewApiClient(),
	}

	return f, nil
}