package tenancyv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameAddDomain = "nutanix_add_domain_v2.test"
const dataSourceNameDomains = "data.nutanix_domains_v2.test"

// TestAccV2NutanixAddDomainResource_Basic adds an existing domain to tenancy and
// verifies that it is represented in Terraform state. The tenancy API does not
// currently offer a delete operation, so the allowed-domain association remains
// after this test completes.
func TestAccV2NutanixAddDomainResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAddDomainConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameAddDomain, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameAddDomain, "tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameDomains, "id"),
					resource.TestCheckResourceAttrSet(dataSourceNameDomains, "domains.#"),
				),
			},
		},
	})
}

func testAddDomainConfig() string {
	return fmt.Sprintf(`
locals {
  config = jsondecode(file(%q))
}

resource "nutanix_add_domain_v2" "test" {
  domain_ext_id = local.config.tenancy.domain_ext_id
}

data "nutanix_domains_v2" "test" {
  depends_on = [nutanix_add_domain_v2.test]
}
`, filepath)
}
