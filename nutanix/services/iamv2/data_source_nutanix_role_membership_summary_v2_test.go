package iamv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixRoleMembershipSummaryV2Datasource_basic(t *testing.T) {
	datasourceName := "data.nutanix_role_membership_summary_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		ProviderFactories: acc.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipSummaryV2DatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "summaries.#"),
				),
			},
		},
	})
}

func testAccNutanixRoleMembershipSummaryV2DatasourceConfig() string {
	return `
data "nutanix_role_membership_summary_v2" "test" {}
`
}
