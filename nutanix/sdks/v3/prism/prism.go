package prism

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

const (
	libraryVersion = "v3"
	absolutePath   = "api/nutanix/" + libraryVersion
	userAgent      = "nutanix/" + libraryVersion
	clientName     = "prism_central"
)

// Client manages the V3 API
type Client struct {
	client *client.Client
	V3     Service
}

// NewV3Client return a client to operate V3 resources
func NewV3Client(credentials client.Credentials) (*Client, error) {
	var baseClient *client.Client

	// check if all required fields are present. Else create an empty client
	// Accept either (username + password) OR api_key for authentication
	hasBasicAuth := credentials.Username != "" && credentials.Password != ""
	hasAPIKey := credentials.APIKey != ""
	hasEndpoint := credentials.Endpoint != ""

	if hasEndpoint && (hasBasicAuth || hasAPIKey) {
		c, err := client.NewClient(&credentials, userAgent, absolutePath, false)
		if err != nil {
			return nil, err
		}
		baseClient = c
	} else {
		errorMsg := fmt.Sprintf("Prism Central (PC) Client is missing. "+
			"Please provide endpoint and either (username + password) or api_key in provider configuration.")

		baseClient = &client.Client{UserAgent: userAgent, ErrorMsg: errorMsg}
	}

	f := &Client{
		client: baseClient,
		V3: Operations{
			client: baseClient,
		},
	}

	// f.client.OnRequestCompleted(func(req *http.Request, resp *http.Response, v interface{}) {
	// 	if v != nil {
	// 		utils.PrintToJSON(v, "[Debug] FINISHED REQUEST")
	// 		// TBD: How to print responses before all requests.
	// 	}
	// })

	return f, nil
}
