package tenancyv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for the tenant action resources (suspend / restore / scheduled-delete):
//   - Suspend a tenant, then restore it, then schedule its deletion. Each action
//     is a standalone mutating resource keyed by the tenant ext_id. These actions
//     have no read API, so we assert the resource id is set. Import not supported.

const resourceNameSuspendTenant = "nutanix_suspend_tenant_v2.test"
const resourceNameRestoreTenant = "nutanix_restore_tenant_v2.test"
const resourceNameScheduledDeleteTenant = "nutanix_scheduled_delete_tenant_v2.test"

func TestAccV2NutanixTenantActions_SuspendRestore(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-tenant-actions-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTenantSuspendConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSuspendTenant, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameSuspendTenant, "id"),
				),
			},
			{
				Config: testTenantRestoreConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreTenant, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreTenant, "id"),
				),
			},
		},
	})
}

func TestAccV2NutanixTenantActions_ScheduledDelete(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-tenant-sched-del-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTenantScheduledDeleteConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameScheduledDeleteTenant, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameScheduledDeleteTenant, "retention_seconds", "3600"),
					resource.TestCheckResourceAttrSet(resourceNameScheduledDeleteTenant, "id"),
				),
			},
		},
	})
}

func testTenantActionPrereq(name string) string {
	return fmt.Sprintf(`
resource "nutanix_tenant_v2" "test" {
  name        = "%[1]s"
  description = "tenant action prerequisite tenant"

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

func testTenantSuspendConfig(name string) string {
	return testTenantActionPrereq(name) + `
resource "nutanix_suspend_tenant_v2" "test" {
  ext_id = nutanix_tenant_v2.test.ext_id
}
`
}

func testTenantRestoreConfig(name string) string {
	return testTenantActionPrereq(name) + `
resource "nutanix_suspend_tenant_v2" "test" {
  ext_id = nutanix_tenant_v2.test.ext_id
}

resource "nutanix_restore_tenant_v2" "test" {
  ext_id     = nutanix_tenant_v2.test.ext_id
  depends_on = [nutanix_suspend_tenant_v2.test]
}
`
}

func testTenantScheduledDeleteConfig(name string) string {
	return testTenantActionPrereq(name) + `
resource "nutanix_scheduled_delete_tenant_v2" "test" {
  ext_id            = nutanix_tenant_v2.test.ext_id
  retention_seconds = 3600
}
`
}
