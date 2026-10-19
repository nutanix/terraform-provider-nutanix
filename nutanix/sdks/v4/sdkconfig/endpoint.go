package sdkconfig

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

// ServiceType represents the type of service being called
type ServiceType int

const (
	// ServiceTypePC is for standard Prism Central API calls
	ServiceTypePC ServiceType = iota
	// ServiceTypeIAM is for IAM-related calls (users, roles, auth policies, etc.)
	ServiceTypeIAM

	// ServiceTypeSP is for Service Provider calls (tenant creation, domain shortlist)
	ServiceTypeSP

	// ServiceTypeDomainManagement is for Domain Management calls (domain creation, domain shortlist, Projects, Resource Groups)
	ServiceTypeDomainManagement
)

// IsNCMode returns true if the provider is configured for Nutanix Cloud (NC) mode
func IsNCMode(creds client.Credentials) bool {
	return creds.NutanixCentralEndpoint != ""
}

// ShouldUseNCIAM returns true if IAM calls should use the NC IAM endpoint.
// This is true when:
// - TenantName is set (Tenant Admin always uses NC IAM), OR
// - UseNCIAM is explicitly set to true (NC Admin opted into NC IAM)
func ShouldUseNCIAM(creds client.Credentials) bool {
	return creds.TenantName != "" || creds.UseNCIAM
}

// ResolveEndpoint determines the correct endpoint based on NC mode and service type.
//
// If NC mode is OFF (nutanix_central_endpoint not set):
//   - Returns endpoint as-is (traditional PC Admin flow)
//
// If NC mode is ON (nutanix_central_endpoint is set):
//   - Service Provider calls: sp.services.<nutanix_central_endpoint>
//   - IAM calls:
//   - If tenant_name is set OR use_nc_iam is true: iam.services.<nutanix_central_endpoint>
//   - Otherwise: <domain_name>.domains.<nutanix_central_endpoint>
//   - All other calls: <domain_name>.domains.<nutanix_central_endpoint> (if domain_name is set)
func ResolveEndpoint(creds client.Credentials, serviceType ServiceType) string {
	// If not in NC mode, return default endpoint as-is
	if !IsNCMode(creds) {
		return creds.Endpoint
	}

	// NC mode is ON
	switch serviceType {
	case ServiceTypeSP:
		// Service Provider calls always go to sp.services
		return fmt.Sprintf("sp.services.%s", creds.NutanixCentralEndpoint)

	case ServiceTypeIAM:
		// IAM calls: use NC IAM if tenant_name is set OR use_nc_iam is true
		if ShouldUseNCIAM(creds) {
			return fmt.Sprintf("iam.services.%s", creds.NutanixCentralEndpoint)
		} else {
			return fmt.Sprintf("%s.domains.%s", creds.DomainName, creds.NutanixCentralEndpoint)
		}

	case ServiceTypeDomainManagement:
		// Domain Management calls always go to Nutanix Central
		return creds.NutanixCentralEndpoint

	default:
		// All other PC calls go to domain if set
		if creds.DomainName != "" {
			return fmt.Sprintf("%s.domains.%s", creds.DomainName, creds.NutanixCentralEndpoint)
		}
		return creds.NutanixCentralEndpoint
	}
}
