package sdkconfig

import (
	"fmt"
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

// V4ApiClient is the common interface for all Nutanix v4 SDK ApiClient types.
// All v4 SDK ApiClient types implement this method.
type V4ApiClient interface {
	AddDefaultHeader(headerName string, headerValue string)
}

// V4ApiClientConfig holds the settable fields common to all v4 ApiClient types.
type V4ApiClientConfig struct {
	Host                    string
	Port                    int
	Username                string
	Password                string
	VerifySSL               bool
	AllowVersionNegotiation bool
}

// ConfigureV4Client checks credentials and returns configuration for a v4 SDK client.
// Returns nil if credentials are insufficient (no endpoint, or neither basic auth nor API key).
// Applies API key and custom headers to the apiClient via AddDefaultHeader.
//
// An error is returned only for a malformed port: falling back to the default port would
// silently talk to a different endpoint than the operator configured.
func ConfigureV4Client(credentials client.Credentials, apiClient V4ApiClient) (*V4ApiClientConfig, error) {
	hasBasicAuth := credentials.Username != "" && credentials.Password != ""
	hasAPIKey := credentials.APIKey != ""
	hasEndpoint := credentials.Endpoint != ""

	if !hasEndpoint || (!hasBasicAuth && !hasAPIKey) {
		return nil, nil
	}

	port := DefaultPort
	if credentials.Port != "" {
		p, err := strconv.Atoi(credentials.Port)
		if err != nil {
			return nil, fmt.Errorf("invalid port %q: must be a number", credentials.Port)
		}
		if p < minPort || p > maxPort {
			return nil, fmt.Errorf("invalid port %d: must be between %d and %d", p, minPort, maxPort)
		}
		port = p
	}

	cfg := &V4ApiClientConfig{
		Host: credentials.Endpoint,
		Port: port,
		// VerifySSL is the inverse of the provider's `insecure` flag, which defaults to
		// false. Disabling verification exposes credentials to man-in-the-middle attacks,
		// so it must stay an explicit operator opt-in rather than a built-in default.
		VerifySSL:               !credentials.Insecure,
		AllowVersionNegotiation: AllowVersionNegotiation,
	}

	// Set authentication - API key takes precedence if both are provided
	if hasAPIKey {
		apiClient.AddDefaultHeader("X-Ntnx-Api-Key", credentials.APIKey)
	} else {
		cfg.Username = credentials.Username
		cfg.Password = credentials.Password
	}

	// Add custom headers (e.g., for Cloudflare Access)
	for key, value := range credentials.CustomHeaders {
		apiClient.AddDefaultHeader(key, value)
	}

	return cfg, nil
}
