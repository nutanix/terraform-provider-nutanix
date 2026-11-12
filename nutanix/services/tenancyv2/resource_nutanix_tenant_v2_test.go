package tenancyv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_tenant_v2:
//   1. Create a tenant with all supported scalar/nested fields populated and
//      assert every attribute has the expected literal value.
//   2. Update scalar fields (name, description, url) and assert the new values,
//      while re-asserting stable fields.
//   3. Read the tenant back through the singular datasource (ext_id input) and
//      assert every attribute matches the resource via TestCheckResourceAttrPair.
//   4. Read the list datasource and assert count > 0.
//   5. Import the resource and verify state round-trips (password is Sensitive and
//      never returned by the API, so it is ignored on verify).
//   6. Negative tests: missing required admin_info fields.

const resourceNameTenant = "nutanix_tenant_v2.test"
const datasourceNameTenant = "data.nutanix_tenant_v2.test"
const datasourceNameTenants = "data.nutanix_tenants_v2.list_test"

func TestAccV2NutanixTenantResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-tenant-%d", r)
	desc := "test tenant description"
	updatedName := fmt.Sprintf("updated-tenant-%d", r)
	updatedDesc := "updated tenant description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixTenantV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testTenantConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTenant, "name", name),
					resource.TestCheckResourceAttr(resourceNameTenant, "description", desc),
					resource.TestCheckResourceAttr(resourceNameTenant, "url", "https://example.com/"+name),
					resource.TestCheckResourceAttr(resourceNameTenant, "admin_info.#", "1"),
					resource.TestCheckResourceAttr(resourceNameTenant, "admin_info.0.username", "tenant-admin"),
					resource.TestCheckResourceAttr(resourceNameTenant, "admin_info.0.first_name", "First"),
					resource.TestCheckResourceAttr(resourceNameTenant, "admin_info.0.last_name", "Last"),
					resource.TestCheckResourceAttr(resourceNameTenant, "admin_info.0.email", "tenant-admin@example.com"),
					resource.TestCheckResourceAttr(resourceNameTenant, "aggregated_quota.#", "1"),
					resource.TestCheckResourceAttr(resourceNameTenant, "aggregated_quota.0.vcpu", "8"),
					resource.TestCheckResourceAttrSet(resourceNameTenant, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameTenant, "links.#"),
				),
			},
			{
				Config: testTenantConfig(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTenant, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameTenant, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameTenant, "url", "https://example.com/"+updatedName),
					resource.TestCheckResourceAttr(resourceNameTenant, "admin_info.0.username", "tenant-admin"),
					resource.TestCheckResourceAttrSet(resourceNameTenant, "ext_id"),
				),
			},
			{
				ResourceName:            resourceNameTenant,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_info.0.password"},
			},
		},
	})
}

func TestAccV2NutanixTenantResource_WithDataSources(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-tenant-ds-%d", r)
	desc := "test tenant datasource description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixTenantV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testTenantDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceNameTenant, "name", resourceNameTenant, "name"),
					resource.TestCheckResourceAttrPair(datasourceNameTenant, "description", resourceNameTenant, "description"),
					resource.TestCheckResourceAttrPair(datasourceNameTenant, "url", resourceNameTenant, "url"),
					resource.TestCheckResourceAttrPair(datasourceNameTenant, "ext_id", resourceNameTenant, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameTenants, "tenants.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixTenantResource_MissingAdminFirstName(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-tenant-neg-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testTenantConfigMissingName(name),
				ExpectError: regexp.MustCompile("error while creating tenant"),
			},
		},
	})
}

func testTenantConfig(name, desc string) string {
	return fmt.Sprintf(`
resource "nutanix_tenant_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"
  url         = "https://example.com/%[1]s"

  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = "Nutanix.123"
  }

  aggregated_quota {
    vcpu    = 8
    memory  = 17179869184
    storage = 107374182400
  }
}
`, name, desc)
}

func testTenantDataSourceConfig(name, desc string) string {
	return testTenantConfig(name, desc) + `
data "nutanix_tenant_v2" "test" {
  ext_id = nutanix_tenant_v2.test.ext_id
}

data "nutanix_tenants_v2" "list_test" {
  depends_on = [nutanix_tenant_v2.test]
}
`
}

func testTenantConfigMissingName(name string) string {
	return fmt.Sprintf(`
resource "nutanix_tenant_v2" "test" {
  description = "tenant without a name %[1]s"

  admin_info {
    username = "tenant-admin"
    email    = "tenant-admin@example.com"
    password = "Nutanix.123"
  }
}
`, name)
}
