package iamv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixRoleMembershipV2Datasource_basic(t *testing.T) {
	datasourceName := "data.nutanix_role_membership_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		ProviderFactories: acc.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipV2DatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "role_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "identity_type"),
				),
			},
		},
	})
}

func testAccNutanixRoleMembershipV2DatasourceConfig() string {
	return `
data "nutanix_role_memberships_v2" "list" {}

data "nutanix_role_membership_v2" "test" {
  ext_id = data.nutanix_role_memberships_v2.list.role_memberships[0].ext_id
}
`
}
