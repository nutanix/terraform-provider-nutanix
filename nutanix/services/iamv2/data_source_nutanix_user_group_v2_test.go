package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameUserGroup = "data.nutanix_user_group_v2.test"

func TestAccV2NutanixUserGroupsDatasource_GetUserGroupByExtId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixUserGroupsV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testUserGroupDatasourceV4Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameUserGroup, "distinguished_name", testVars.Iam.UserGroups.DistinguishedName),
					resource.TestCheckResourceAttr(datasourceNameUserGroup, "name", testVars.Iam.UserGroups.Name),
					resource.TestCheckResourceAttr(datasourceNameUserGroup, "idp_id", testVars.Iam.Users.DirectoryServiceID),
					resource.TestCheckResourceAttr(datasourceNameUserGroup, "group_type", "LDAP"),
				),
			},
		},
	})
}

func testUserGroupDatasourceV4Config(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%[1]s")))
		users = local.config.iam.users
		user_groups = local.config.iam.user_groups
	}

	resource "nutanix_user_groups_v2" "test" {
		group_type = "LDAP"
		idp_id =  local.users.directory_service_id
		name = local.user_groups.name
		distinguished_name = local.user_groups.distinguished_name
	}

	data "nutanix_user_group_v2" "test" {
		ext_id = resource.nutanix_user_groups_v2.test.id
	}
	`, filepath)
}
