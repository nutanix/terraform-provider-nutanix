package iamv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameUsers = "data.nutanix_users_v2.test"

func TestAccNutanixUsersV4Datasource_Basic(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUsersDatasourceV4Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUsers, "users.#"),
					resource.TestCheckResourceAttrSet(datasourceNameUsers, "users.0.username"),
					resource.TestCheckResourceAttrSet(datasourceNameUsers, "users.0.user_type"),
					resource.TestCheckResourceAttrSet(datasourceNameUsers, "users.0.ext_id"),
				),
			},
		},
	})
}

func TestAccNutanixUsersV4Datasource_WithFilter(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUsersDatasourceV4WithFilterConfig(filepath, "userType eq Schema.Enums.UserType'LOCAL'"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUsers, "users.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameUsers, "users.0.user_type", "LOCAL"),
				),
			},
			{
				Config: testUsersDatasourceV4WithFilterConfig(filepath, "username eq '"+testVars.Iam.Users.Username+"'"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUsers, "users.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameUsers, "users.0.username", testVars.Iam.Users.Username),
				),
			},
		},
	})
}

func TestAccNutanixUsersV4Datasource_WithLimit(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUsersDatasourceV4WithLimitConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameUsers, "users.#", strconv.Itoa(testVars.Iam.Users.Limit)),
				),
			},
		},
	})
}

func testUsersDatasourceV4Config(filepath string) string {
	return fmt.Sprintf(`
		locals{
			config = (jsondecode(file("%s")))
			users = local.config.iam.users
		}
		
		resource "nutanix_users_v2" "test" {
			username = local.users.username
			first_name = local.users.first_name
			middle_initial = local.users.middle_initial
			last_name = local.users.last_name
			email_id = local.users.email_id
			locale = local.users.locale
			region = local.users.region
			display_name = local.users.display_name
			password = local.users.password
			user_type = "LOCAL"
			status = "ACTIVE"  
			force_reset_password = local.users.force_reset_password   
		}

		data "nutanix_users_v2" "test"{
			depends_on = [nutanix_users_v2.test]
		}
	`, filepath)
}

func testUsersDatasourceV4WithFilterConfig(filepath, userQuery string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		users = local.config.iam.users
	}
	
	resource "nutanix_users_v2" "test" {
		username = local.users.username
		first_name = local.users.first_name
		middle_initial = local.users.middle_initial
		last_name = local.users.last_name
		email_id = local.users.email_id
		locale = local.users.locale
		region = local.users.region
		display_name = local.users.display_name
		password = local.users.password
		user_type = "LOCAL"
		status = "ACTIVE"  
		force_reset_password = local.users.force_reset_password   
	}
	
	data "nutanix_users_v2" "test" {
		filter = "%s"
		depends_on = [nutanix_users_v2.test]
	}

	
	`, filepath, userQuery)
}

func testUsersDatasourceV4WithLimitConfig(filepath string) string {
	return fmt.Sprintf(`
		locals{
			config = (jsondecode(file("%s")))
			users = local.config.iam.users
		}
		
		resource "nutanix_users_v2" "test" {
			username = local.users.username
			first_name = local.users.first_name
			middle_initial = local.users.middle_initial
			last_name = local.users.last_name
			email_id = local.users.email_id
			locale = local.users.locale
			region = local.users.region
			display_name = local.users.display_name
			password = local.users.password
			user_type = "LOCAL"
			status = "ACTIVE"  
			force_reset_password = local.users.force_reset_password   
		}
		
		data "nutanix_users_v2" "test" {
			limit     = local.users.limit
			depends_on = [nutanix_users_v2.test]
		}
	`, filepath)
}
