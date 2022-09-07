package Era

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	libraryVersion = "v0.9"
	absolutePath   = "era/" + libraryVersion
	clientName     = "ERA"
)

type Client struct {
	client  *client.Client
	Service Service
}

func NewEraClient(credentials client.Credentials) (*Client, error) {
	var baseClient *client.Client

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.EraEndpoint != "" {
		credentials.URL = fmt.Sprintf(credentials.EraEndpoint)
		c, err := client.NewBaseClient(&credentials, absolutePath, false)
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		errorMsg := fmt.Sprintf("Era Client is missing. "+
			"Please provide required details - %s in provider configuration.", strings.Join(credentials.RequiredFields[clientName], ", "))

		baseClient = &client.Client{ErrorMsg: errorMsg}
	}

	era := &Client{
		client: baseClient,
		Service: ServiceClient{
			c: baseClient,
		},
	}
	return era, nil

}
