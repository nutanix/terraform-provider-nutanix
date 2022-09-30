package calm

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

const (
	libraryVersion = "v3.0"
	absolutePath   = "api/calm/" + libraryVersion
	userAgent      = "nutanix/" + libraryVersion
	clientName     = "calm"
)

// Client manages the Calm API
type Client struct {
	client *client.Client
	V3     Service
}

// NewV3Client return a client to operate V3 resources
func NewV3Client(credentials client.Credentials) (*Client, error) {
	var baseClient *client.Client

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		c, err := client.NewClient(&credentials, userAgent, absolutePath, false)
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		errorMsg := fmt.Sprintf("Prism Central (PC) Client is missing. "+
			"Please provide required details - %s in provider configuration.", strings.Join(credentials.RequiredFields[clientName], ", "))

		baseClient = &client.Client{UserAgent: userAgent, ErrorMsg: errorMsg}
	}

	f := &Client{
		client: baseClient,
		V3: Operations{
			client: baseClient,
		},
	}

	return f, nil
}
