package tenancyv2_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// testAccCheckNutanixTenantV2Destroy verifies that every managed tenant has been
// removed by calling GetTenantById and asserting it returns an error.
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
		_, err := conn.ServiceProviderManagementAPI.TenantManagementAPIInstance.GetTenantById(ctx, &getTenantRequest)
		if err == nil {
			return fmt.Errorf("tenant %s still exists", rs.Primary.ID)
		}
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
		_, err := conn.ServiceProviderManagementAPI.TenantManagementAPIInstance.GetTenantDomainAllocationById(ctx, &getRequest)
		if err == nil {
			return fmt.Errorf("tenant domain allocation %s still exists", rs.Primary.ID)
		}
	}
	return nil
}
