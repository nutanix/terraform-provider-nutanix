package lifecycle

import (
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/api"
	lifecycle "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

// Client wraps the lifecycle-go-client v4 API instances used by the provider.
//
// The FoundationCentralConfig APIs are served by the Foundation Central (FCVM)
// service, NOT Prism Central. When the provider is configured with a Foundation
// endpoint/port, this client is pointed at that host; otherwise it falls back to
// the standard Prism Central endpoint.
type Client struct {
	FoundationCentralConfigAPIInstance *api.FoundationCentralConfigApi
}

// NewLifecycleClient returns a client to operate lifecycle (Foundation Central
// management) resources. It targets the Foundation/FCVM endpoint when configured.
func NewLifecycleClient(credentials client.Credentials) (*Client, error) {
	var baseClient *lifecycle.ApiClient

	pcClient := lifecycle.NewApiClient()

	// FoundationCentralConfig is an FCVM/Foundation-scoped API. Prefer the
	// Foundation endpoint (host) when the provider supplies it so the API talks
	// to the correct host. The v4 lifecycle API is served over the standard v4
	// HTTPS port (9440), NOT the legacy Foundation VM port (8000); keep the
	// Prism-Central-style port here instead of the legacy foundation_port.
	if credentials.FoundationEndpoint != "" {
		credentials.Endpoint = credentials.FoundationEndpoint
	}

	if cfg := sdkconfig.ConfigureV4Client(credentials, pcClient); cfg != nil {
		pcClient.Host = cfg.Host
		pcClient.Port = cfg.Port
		pcClient.Username = cfg.Username
		pcClient.Password = cfg.Password
		pcClient.VerifySSL = cfg.VerifySSL
		pcClient.AllowVersionNegotiation = cfg.AllowVersionNegotiation
		baseClient = pcClient
	}

	return &Client{
		FoundationCentralConfigAPIInstance: api.NewFoundationCentralConfigApi(baseClient),
	}, nil
}
