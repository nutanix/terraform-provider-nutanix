package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameUserGroups = "data.nutanix_user_groups_v2.test"

func TestAccV2NutanixUserGroupsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserGroupsDatasourceV4Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUserGroups, "user_groups.#"),
					checkAttributeLength(datasourceNameUserGroups, "user_groups", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixUserGroupsDatasource_WithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserGroupsDatasourceV4WithFilterConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUserGroups, "user_groups.#"),
					resource.TestCheckResourceAttr(datasourceNameUserGroups, "user_groups.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameUserGroups, "user_groups.0.distinguished_name", testVars.Iam.UserGroups.DistinguishedName),
					resource.TestCheckResourceAttr(datasourceNameUserGroups, "user_groups.0.name", testVars.Iam.UserGroups.Name),
					resource.TestCheckResourceAttr(datasourceNameUserGroups, "user_groups.0.idp_id", testVars.Iam.Users.DirectoryServiceID),
					resource.TestCheckResourceAttr(datasourceNameUserGroups, "user_groups.0.group_type", "LDAP"),
					resource.TestCheckResourceAttrSet(datasourceNameUserGroups, "user_groups.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameUserGroups, "user_groups.0.created_time"),
					resource.TestCheckResourceAttrSet(datasourceNameUserGroups, "user_groups.0.created_by"),
				),
			},
		},
	})
}

func TestAccV2NutanixUserGroupsDatasource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserGroupsDatasourceV4WithLimitConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUserGroups, "user_groups.#"),
					resource.TestCheckResourceAttr(datasourceNameUserGroups, "user_groups.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixUserGroupsDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserGroupsDatasourceV4WithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUserGroups, "user_groups.#"),
					resource.TestCheckResourceAttr(datasourceNameUserGroups, "user_groups.#", "0"),
				),
			},
		},
	})
}

func testUserGroupsDatasourceV4Config(filepath string) string {
	return fmt.Sprintf(`
	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
		user_groups = local.config.iam.user_groups
	}

	resource "nutanix_user_groups_v2" "test" {
		group_type = "LDAP"
		idp_id =  local.users.directory_service_id
		name = local.user_groups.name
		distinguished_name = local.user_groups.distinguished_name
	  }

	data "nutanix_user_groups_v2" "test"{
		depends_on = [resource.nutanix_user_groups_v2.test]
	}
	`, filepath)
}

func testUserGroupsDatasourceV4WithFilterConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
		user_groups = local.config.iam.user_groups
	}

	resource "nutanix_user_groups_v2" "test" {
		group_type = "LDAP"
		idp_id =  local.users.directory_service_id
		name = local.user_groups.name
		distinguished_name = local.user_groups.distinguished_name
	  }

	  data "nutanix_user_groups_v2" "test" {
		filter     = "name eq '${local.user_groups.name}'"
		depends_on = [resource.nutanix_user_groups_v2.test]
	  }
	`, filepath)
}

func testUserGroupsDatasourceV4WithLimitConfig(filepath string) string {
	return fmt.Sprintf(`
		locals{
			config = (jsondecode(file("%s")))
			users = local.config.iam.users
			user_groups = local.config.iam.user_groups
		}

		resource "nutanix_user_groups_v2" "test" {
			group_type = "LDAP"
			idp_id =  local.users.directory_service_id
			name = local.user_groups.name
			distinguished_name = local.user_groups.distinguished_name
		  }

		data "nutanix_user_groups_v2" "test" {
			limit      = 1
			depends_on = [resource.nutanix_user_groups_v2.test]
		}
	`, filepath)
}

func testUserGroupsDatasourceV4WithInvalidFilterConfig() string {
	return `
	  data "nutanix_user_groups_v2" "test" {
		filter     = "name eq 'invalid_filter'"
	  }
	`
}
