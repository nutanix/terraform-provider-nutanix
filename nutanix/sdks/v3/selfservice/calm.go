package selfservice

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

const (
	libraryVersion = "v3"
	absolutePath   = "api/nutanix/" + libraryVersion
	userAgent      = "nutanix/" + libraryVersion
	clientName     = "Self Service"
)

// Client manages the foundation central API
type Client struct {
	client  *client.Client
	Service Service
}

func NewCalmClient(credentials client.Credentials) (*Client, error) {
	var baseClient *client.Client

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		c, err := client.NewClient(&credentials, userAgent, absolutePath, false)
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		errorMsg := fmt.Sprintf("Self Service Client is missing. "+
			"Please provide required details - %s in provider configuration.", strings.Join(credentials.RequiredFields[clientName], ", "))

		baseClient = &client.Client{UserAgent: userAgent, ErrorMsg: errorMsg}
	}

	fc := &Client{
		client: baseClient,
		Service: Operations{
			client: baseClient,
		},
	}
	return fc, nil
}
