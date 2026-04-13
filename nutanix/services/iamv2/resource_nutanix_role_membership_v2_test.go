package iamv2_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/rolemembership"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestAccNutanixRoleMembershipV2Resource_basic(t *testing.T) {
	resourceName := "nutanix_role_membership_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		ProviderFactories: acc.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNutanixRoleMembershipV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_ext_id"),
					resource.TestCheckResourceAttr(resourceName, "identity_type", "USER"),
					resource.TestCheckResourceAttrSet(resourceName, "identity_ext_id"),
				),
			},
		},
	})
}

func TestAccNutanixRoleMembershipV2Resource_import(t *testing.T) {
	resourceName := "nutanix_role_membership_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		ProviderFactories: acc.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNutanixRoleMembershipV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipV2Config(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixRoleMembershipV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_role_membership_v2" {
			continue
		}
		getRequest := import1.GetRoleMembershipByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := conn.IamAPI.RoleMembershipAPIInstance.GetRoleMembershipById(ctx, &getRequest)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "not found") || strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
		deleteRequest := import1.DeleteRoleMembershipByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		if _, err := conn.IamAPI.RoleMembershipAPIInstance.DeleteRoleMembershipById(ctx, &deleteRequest); err != nil {
			return err
		}
	}
	return nil
}

func testAccNutanixRoleMembershipV2Config() string {
	return `
data "nutanix_roles_v2" "roles" {}

data "nutanix_users_v2" "users" {}

resource "nutanix_role_membership_v2" "test" {
  role_ext_id      = data.nutanix_roles_v2.roles.roles[0].ext_id
  identity_type    = "USER"
  identity_ext_id  = data.nutanix_users_v2.users.users[0].ext_id
}
`
}
