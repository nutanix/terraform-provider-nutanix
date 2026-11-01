package tenancyv2_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import1 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/request/membership"

	import4 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func checkAttributeLength(resourceName, attribute string, minLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}
		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}
		if count < minLength {
			return fmt.Errorf("expected %s to be >= %d, got %d", attrKey, minLength, count)
		}
		return nil
	}
}

// tenantMembershipTestInputs holds the infrastructure identifiers required to
// run the live tenant-membership lifecycle test. These reference SP-tenancy
// entities (an allocated domain, a role, and a principal) that cannot be
// created by this provider, so they are supplied via environment variables.
type tenantMembershipTestInputs struct {
	DomainExtID           string
	RoleExtID             string
	IdentityExtID         string
	IdentityProviderExtID string
	PrincipalValue        string
	IdentityType          string
}

// tenantMembershipInputsFromEnv reads the membership test inputs from the
// environment and skips the test when the mandatory ones are missing.
func tenantMembershipInputsFromEnv(t *testing.T) tenantMembershipTestInputs {
	t.Helper()

	in := tenantMembershipTestInputs{
		DomainExtID:           os.Getenv("NUTANIX_DOMAIN_EXT_ID"),
		RoleExtID:             os.Getenv("NUTANIX_MEMBERSHIP_ROLE_EXT_ID"),
		IdentityExtID:         os.Getenv("NUTANIX_MEMBERSHIP_IDENTITY_EXT_ID"),
		IdentityProviderExtID: os.Getenv("NUTANIX_MEMBERSHIP_IDP_EXT_ID"),
		PrincipalValue:        os.Getenv("NUTANIX_MEMBERSHIP_PRINCIPAL_VALUE"),
		IdentityType:          os.Getenv("NUTANIX_MEMBERSHIP_IDENTITY_TYPE"),
	}

	if in.DomainExtID == "" || in.RoleExtID == "" {
		t.Skip("skipping tenant membership lifecycle test: set NUTANIX_DOMAIN_EXT_ID and NUTANIX_MEMBERSHIP_ROLE_EXT_ID (plus either NUTANIX_MEMBERSHIP_IDENTITY_EXT_ID or the identity_details env vars) to run it")
	}
	if in.IdentityExtID == "" && (in.IdentityProviderExtID == "" || in.PrincipalValue == "") {
		t.Skip("skipping tenant membership lifecycle test: provide NUTANIX_MEMBERSHIP_IDENTITY_EXT_ID or both NUTANIX_MEMBERSHIP_IDP_EXT_ID and NUTANIX_MEMBERSHIP_PRINCIPAL_VALUE")
	}
	if in.IdentityType == "" {
		in.IdentityType = "USER"
	}
	return in
}

func testAccCheckNutanixTenantMembershipV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).ServiceProviderManagementAPI
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_tenant_membership_v2" {
			continue
		}
		getRequest := import3.GetDomainMembershipByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		if _, err := conn.TenantMembershipAPIInstance.GetDomainMembershipById(ctx, &getRequest); err == nil {
			return fmt.Errorf("tenant membership %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// testAccCheckNutanixTenantV2Destroy verifies that every managed tenant has been
// removed. Deleting a tenant is a soft delete: GetTenantById keeps returning the
// tenant with its state set to DELETED (this is the terminal state the resource
// delete itself waits for). So the tenant is considered destroyed when either
// the GET fails (not found) or the returned tenant is in the DELETED state.
func testAccCheckNutanixTenantV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_tenant_v2" {
			continue
		}

		getTenantRequest := import4.GetTenantByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		resp, err := conn.ServiceProviderManagementAPI.TenantManagementAPIInstance.GetTenantById(ctx, &getTenantRequest)
		if err != nil {
			// Not found: the tenant has been fully removed.
			continue
		}
		if resp != nil && resp.Data != nil {
			if tenant, ok := resp.Data.GetValue().(import1.Tenant); ok {
				if tenant.StateInfo != nil && tenant.StateInfo.State != nil &&
					tenant.StateInfo.State.GetName() == "DELETED" {
					// Soft-deleted: treat as destroyed.
					continue
				}
			}
		}
		return fmt.Errorf("tenant %s still exists", rs.Primary.ID)
	}
	return nil
}

// testAccCheckNutanixTenantDomainAllocationV2Destroy verifies that every managed
// domain allocation has been removed by calling GetTenantDomainAllocationById.
func testAccCheckNutanixTenantDomainAllocationV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_tenant_domain_allocation_v2" {
			continue
		}

		getRequest := import4.GetTenantDomainAllocationByIdRequest{
			TenantExtId: utils.StringPtr(rs.Primary.Attributes["tenant_ext_id"]),
			ExtId:       utils.StringPtr(rs.Primary.ID),
		}
		resp, err := conn.ServiceProviderManagementAPI.TenantManagementAPIInstance.GetTenantDomainAllocationById(ctx, &getRequest)
		if err != nil {
			// A hard error (e.g. not found) means the allocation is gone.
			continue
		}
		// The allocation is soft-deleted: the GET still succeeds but the state
		// transitions to DELETED. Treat that as destroyed.
		if resp != nil && resp.Data != nil {
			if allocation, ok := resp.Data.GetValue().(import1.DomainAllocation); ok {
				if allocation.StateInfo != nil && allocation.StateInfo.State != nil &&
					allocation.StateInfo.State.GetName() == "DELETED" {
					continue
				}
			}
		}
		return fmt.Errorf("tenant domain allocation %s still exists", rs.Primary.ID)
	}
	return nil
}
