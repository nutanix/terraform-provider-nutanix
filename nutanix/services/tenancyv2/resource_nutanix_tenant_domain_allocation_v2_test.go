package tenancyv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_tenant_domain_allocation_v2:
//   1. Create a tenant (prerequisite) and a domain allocation with projects and
//      assert every attribute.
//   2. Update the projects list and assert the new values.
//   3. Read the domain allocation back through the singular and list datasources.
//   4. Verify CheckDestroy removes the allocation.
// External network connection settings are exercised in a dedicated test.

const resourceNameDomainAllocation = "nutanix_tenant_domain_allocation_v2.test"
const datasourceNameDomainAllocation = "data.nutanix_tenant_domain_allocation_v2.test"
const datasourceNameDomainAllocations = "data.nutanix_tenant_domain_allocations_v2.list_test"

func TestAccV2NutanixTenantDomainAllocationResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-tenant-da-%d", r)
	domainExtID := "00000000-0000-0000-0000-000000000000"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixTenantDomainAllocationV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testTenantDomainAllocationConfig(name, domainExtID, "project-one"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDomainAllocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameDomainAllocation, "tenant_ext_id"),
					resource.TestCheckResourceAttr(resourceNameDomainAllocation, "domain_ext_id", domainExtID),
					resource.TestCheckResourceAttr(resourceNameDomainAllocation, "projects.#", "1"),
					resource.TestCheckResourceAttr(resourceNameDomainAllocation, "projects.0.name", "project-one"),
					resource.TestCheckResourceAttrSet(resourceNameDomainAllocation, "links.#"),
				),
			},
			{
				Config: testTenantDomainAllocationConfig(name, domainExtID, "project-two"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDomainAllocation, "projects.#", "1"),
					resource.TestCheckResourceAttr(resourceNameDomainAllocation, "projects.0.name", "project-two"),
				),
			},
			{
				Config: testTenantDomainAllocationDataSourceConfig(name, domainExtID, "project-two"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceNameDomainAllocation, "ext_id", resourceNameDomainAllocation, "ext_id"),
					resource.TestCheckResourceAttrPair(datasourceNameDomainAllocation, "domain_ext_id", resourceNameDomainAllocation, "domain_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameDomainAllocations, "domain_allocations.#"),
				),
			},
		},
	})
}

func testTenantDomainAllocationPrereq(name string) string {
	return fmt.Sprintf(`
resource "nutanix_tenant_v2" "test" {
  name        = "%[1]s"
  description = "domain allocation prerequisite tenant"

  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = "Nutanix.123"
  }
}
`, name)
}

func testTenantDomainAllocationConfig(name, domainExtID, projectName string) string {
	return testTenantDomainAllocationPrereq(name) + fmt.Sprintf(`
resource "nutanix_tenant_domain_allocation_v2" "test" {
  tenant_ext_id = nutanix_tenant_v2.test.ext_id
  domain_ext_id = "%[1]s"

  projects {
    name = "%[2]s"
    quota {
      vcpu    = 4
      memory  = 8589934592
      storage = 53687091200
    }
  }
}
`, domainExtID, projectName)
}

func testTenantDomainAllocationDataSourceConfig(name, domainExtID, projectName string) string {
	return testTenantDomainAllocationConfig(name, domainExtID, projectName) + `
data "nutanix_tenant_domain_allocation_v2" "test" {
  tenant_ext_id = nutanix_tenant_v2.test.ext_id
  ext_id        = nutanix_tenant_domain_allocation_v2.test.ext_id
}

data "nutanix_tenant_domain_allocations_v2" "list_test" {
  tenant_ext_id = nutanix_tenant_v2.test.ext_id
  depends_on    = [nutanix_tenant_domain_allocation_v2.test]
}
`
}
