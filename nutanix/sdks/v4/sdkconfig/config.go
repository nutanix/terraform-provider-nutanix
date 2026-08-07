package sdkconfig

import (
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
// Uses ServiceTypePC (standard PC calls) for endpoint resolution.
func ConfigureV4Client(credentials client.Credentials, apiClient V4ApiClient) *V4ApiClientConfig {
	return ConfigureV4ClientForService(credentials, apiClient, ServiceTypePC)
}

// ConfigureV4ClientForIAM configures a v4 SDK client for IAM-related calls.
// Uses the appropriate IAM endpoint based on admin type and use_nc_iam setting.
func ConfigureV4ClientForIAM(credentials client.Credentials, apiClient V4ApiClient) *V4ApiClientConfig {
	return ConfigureV4ClientForService(credentials, apiClient, ServiceTypeIAM)
}

// ConfigureV4ClientForDomainManagement configures a v4 SDK client for Domain Management calls.
func ConfigureV4ClientForDomainManagement(credentials client.Credentials, apiClient V4ApiClient) *V4ApiClientConfig {
	return ConfigureV4ClientForService(credentials, apiClient, ServiceTypeDomainManagement)
}

// ConfigureV4ClientForSP configures a v4 SDK client for Service Provider calls.
// Used for tenant creation, domain shortlist, and tenant membership management.
func ConfigureV4ClientForSP(credentials client.Credentials, apiClient V4ApiClient) *V4ApiClientConfig {
	return ConfigureV4ClientForService(credentials, apiClient, ServiceTypeSP)
}

// ConfigureV4ClientForService is the core configuration function that handles endpoint
// resolution based on service type and admin type.
func ConfigureV4ClientForService(credentials client.Credentials, apiClient V4ApiClient, serviceType ServiceType) *V4ApiClientConfig {
	hasBasicAuth := credentials.Username != "" && credentials.Password != ""
	hasAPIKey := credentials.APIKey != ""
	hasEndpoint := credentials.Endpoint != ""
	hasNCEndpoint := credentials.NCEndpoint != ""
	isTenant := credentials.TenantName != ""

	// For NC modes, we need NCEndpoint; for PC Admin mode, we need Endpoint
	if !hasNCEndpoint && !hasEndpoint {
		return nil
	}
	if !hasBasicAuth && !hasAPIKey {
		return nil
	}

	// Resolve the endpoint based on admin type and service type
	host := ResolveEndpoint(credentials, serviceType)

	// Determine port - use 443 for NC mode (FQDN-based), otherwise use configured or default
	port := DefaultPort
	if hasNCEndpoint {
		// NC mode uses HTTPS on port 443
		port = 443
	}
	if credentials.Port != "" {
		if p, err := strconv.Atoi(credentials.Port); err == nil {
			port = p
		}
	}

	cfg := &V4ApiClientConfig{
		Host:                    host,
		Port:                    port,
		VerifySSL:               false,
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

	// Set tenant header - use configured tenant_name
	if isTenant {
		apiClient.AddDefaultHeader("x-tenant-name", credentials.TenantName)
	}

	return cfg
}
