package tenancyv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_tenant_external_network_connection_v2:
//   1. Create a tenant + domain allocation (prerequisites).
//   2. Configure a no-NAT external network connection and assert the config.
//   3. Read the connection back through the singular datasource.
// Note: this resource has no delete API (it lives with the domain allocation);
// deleting only clears Terraform state, so no CheckDestroy is defined for it.
// Import not supported for this resource.

const resourceNameExtNetConn = "nutanix_tenant_external_network_connection_v2.test"
const datasourceNameExtNetConn = "data.nutanix_tenant_external_network_connection_v2.test"

func TestAccV2NutanixTenantExternalNetworkConnectionResource_NoNat(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-tenant-enc-%d", r)
	transitVpc := "00000000-0000-0000-0000-000000000000"
	overlaySubnet := "11111111-1111-1111-1111-111111111111"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTenantExternalNetworkConnectionConfig(name, transitVpc, overlaySubnet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameExtNetConn, "tenant_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameExtNetConn, "domain_allocation_ext_id"),
					resource.TestCheckResourceAttr(resourceNameExtNetConn, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameExtNetConn, "config.0.no_nat.#", "1"),
					resource.TestCheckResourceAttr(resourceNameExtNetConn, "config.0.no_nat.0.transit_vpc_ext_id", transitVpc),
					resource.TestCheckResourceAttr(resourceNameExtNetConn, "config.0.no_nat.0.overlay_subnet_ext_id", overlaySubnet),
				),
			},
			{
				Config: testTenantExternalNetworkConnectionDataSourceConfig(name, transitVpc, overlaySubnet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceNameExtNetConn, "tenant_ext_id", resourceNameExtNetConn, "tenant_ext_id"),
					resource.TestCheckResourceAttrPair(datasourceNameExtNetConn, "domain_allocation_ext_id", resourceNameExtNetConn, "domain_allocation_ext_id"),
					resource.TestCheckResourceAttr(datasourceNameExtNetConn, "config.0.no_nat.0.transit_vpc_ext_id", transitVpc),
				),
			},
		},
	})
}

func testTenantExternalNetworkConnectionPrereq(name string) string {
	return fmt.Sprintf(`
resource "nutanix_tenant_v2" "test" {
  name        = "%[1]s"
  description = "external network connection prerequisite tenant"

  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = "Nutanix.123"
  }
}

resource "nutanix_tenant_domain_allocation_v2" "test" {
  tenant_ext_id = nutanix_tenant_v2.test.ext_id
  domain_ext_id = "00000000-0000-0000-0000-000000000000"
}
`, name)
}

func testTenantExternalNetworkConnectionConfig(name, transitVpc, overlaySubnet string) string {
	return testTenantExternalNetworkConnectionPrereq(name) + fmt.Sprintf(`
resource "nutanix_tenant_external_network_connection_v2" "test" {
  tenant_ext_id            = nutanix_tenant_v2.test.ext_id
  domain_allocation_ext_id = nutanix_tenant_domain_allocation_v2.test.ext_id

  config {
    no_nat {
      transit_vpc_ext_id    = "%[1]s"
      overlay_subnet_ext_id = "%[2]s"
      externally_routable_prefixes {
        ipv4 {
          ip {
            value = "10.0.0.0"
          }
          prefix_length = 24
        }
      }
    }
  }
}
`, transitVpc, overlaySubnet)
}

func testTenantExternalNetworkConnectionDataSourceConfig(name, transitVpc, overlaySubnet string) string {
	return testTenantExternalNetworkConnectionConfig(name, transitVpc, overlaySubnet) + `
data "nutanix_tenant_external_network_connection_v2" "test" {
  tenant_ext_id            = nutanix_tenant_v2.test.ext_id
  domain_allocation_ext_id = nutanix_tenant_domain_allocation_v2.test.ext_id
  depends_on               = [nutanix_tenant_external_network_connection_v2.test]
}
`
}
