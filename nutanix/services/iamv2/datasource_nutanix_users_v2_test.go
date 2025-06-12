package iamv2_test

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameUsers = "data.nutanix_users_v2.test"
const dataSourceServiceAccount = "data.nutanix_users_v2.service_account"

func TestAccV2NutanixUsersDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-user-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUsersDatasourceV4Config(filepath, name),
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

func TestAccV2NutanixUsersDatasource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-user-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUsersDatasourceV4WithFilterConfig(filepath, name, "userType eq Schema.Enums.UserType'LOCAL'"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUsers, "users.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameUsers, "users.0.user_type", "LOCAL"),
				),
			},
			{
				Config: testUsersDatasourceV4WithFilterConfig(filepath, name, "username eq '"+name+"'"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameUsers, "users.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameUsers, "users.0.username", name),
				),
			},
		},
	})
}

func TestAccV2NutanixUsersDatasource_WithLimit(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-user-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUsersDatasourceV4WithLimitConfig(filepath, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameUsers, "users.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixUsersDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUsersDatasourceV4WithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameUsers, "users.#", "0"),
				),
			},
		},
	})
}

func testUsersDatasourceV4Config(filepath, name string) string {
	return fmt.Sprintf(`
		locals{
			config = (jsondecode(file("%[1]s")))
			users = local.config.iam.users
		}

		resource "nutanix_users_v2" "test" {
			username = "%[2]s"
			first_name = "first-name-%[2]s"
			middle_initial = "middle-initial-%[2]s"
			last_name = "last-name-%[2]s"
			email_id = local.users.email_id
			locale = local.users.locale
			region = local.users.region
			display_name = "display-name-%[2]s"
			password = local.users.password
			user_type = "LOCAL"
			status = "ACTIVE"
			force_reset_password = local.users.force_reset_password
		}

		data "nutanix_users_v2" "test"{
			depends_on = [nutanix_users_v2.test]
		}
	`, filepath, name)
}

func testUsersDatasourceV4WithFilterConfig(filepath, name, userQuery string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%[1]s")))
		users = local.config.iam.users
	}

	resource "nutanix_users_v2" "test" {
		username = "%[2]s"
		first_name = "first-name-%[2]s"
		middle_initial = "middle-initial-%[2]s"
		last_name = "last-name-%[2]s"
		email_id = local.users.email_id
		locale = local.users.locale
		region = local.users.region
		display_name = "display-name-%[2]s"
		password = local.users.password
		user_type = "LOCAL"
		status = "ACTIVE"
		force_reset_password = local.users.force_reset_password
	}

	data "nutanix_users_v2" "test" {
		filter = "%[3]s"
		depends_on = [nutanix_users_v2.test]
	}


	`, filepath, name, userQuery)
}

func testUsersDatasourceV4WithLimitConfig(filepath, name string) string {
	return fmt.Sprintf(`
		locals{
			config = (jsondecode(file("%[1]s")))
			users = local.config.iam.users
		}

		resource "nutanix_users_v2" "test" {
			username = "%[2]s"
			first_name = "first-name-%[2]s"
			middle_initial = "middle-initial-%[2]s"
			last_name = "last-name-%[2]s"
			email_id = local.users.email_id
			locale = local.users.locale
			region = local.users.region
			display_name = "display-name-%[2]s"
			password = local.users.password
			user_type = "LOCAL"
			status = "ACTIVE"
			force_reset_password = local.users.force_reset_password
		}

		data "nutanix_users_v2" "test" {
			limit     = 1
			depends_on = [nutanix_users_v2.test]
		}
	`, filepath, name)
}

func TestAccV2NutanixUsersDataSourceServiceAccount(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("service-account-unique%d", r)
	expectedKeys := []map[string]string{
		{
			"username":    name,
			"description": "test service account tf",
			"email_id":    "terraform_plugin@domain.com",
			"user_type":   "SERVICE_ACCOUNT",
		},
		{
			"username":    name + "_another",
			"description": "test service account tf another",
			"email_id":    "terraform_plugin_another@domain.com",
			"user_type":   "SERVICE_ACCOUNT",
		},
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiecAccountDataSourceConfig(name, "userType eq Schema.Enums.UserType'SERVICE_ACCOUNT' and username contains '"+name+"'"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceServiceAccount, "users.#"),
					resource.TestCheckResourceAttr(dataSourceServiceAccount, "users.#", "2"),
					checkServiceAccountValues(expectedKeys, dataSourceServiceAccount),
				),
			},
		},
	})
}

func testServiecAccountDataSourceConfig(name string, query string) string {
	return fmt.Sprintf(`
	// Create service account
	resource "nutanix_users_v2" "service_account" {
		username = "%[2]s"
		description = "test service account tf"
		email_id = "terraform_plugin@domain.com"
		user_type = "SERVICE_ACCOUNT"
	}

	// Create another service account
	resource "nutanix_users_v2" "another_service_account" {
		username = "%[2]s_another"
		description = "test service account tf another"
		email_id = "terraform_plugin_another@domain.com"
		user_type = "SERVICE_ACCOUNT"
	}
	
	// Data source to fetch the service account
	data "nutanix_users_v2" "service_account" {
		filter = "%[3]s"
		depends_on = [nutanix_users_v2.service_account, nutanix_users_v2.another_service_account]
	}
	`, filepath, name, query)
}

func checkServiceAccountValues(expectedKeys []map[string]string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource Not found")
		}

		keys := rs.Primary.Attributes
		keyCount, err := strconv.Atoi(rs.Primary.Attributes["users.#"])
		if err != nil {
			return fmt.Errorf("error converting users count: %v", err)
		}

		if keyCount != len(expectedKeys) {
			return fmt.Errorf("expected %d keys, found %d", len(expectedKeys), keyCount)
		}
		log.Printf("found %d keys\n", keyCount)

		// Match each expected key
		for _, expected := range expectedKeys {
			found := false
			for i := 0; i < keyCount; i++ {
				if keys[fmt.Sprintf("users.%d.username", i)] == expected["username"] &&
					keys[fmt.Sprintf("users.%d.description", i)] == expected["description"] &&
					keys[fmt.Sprintf("users.%d.email_id", i)] == expected["email_id"] &&
					keys[fmt.Sprintf("users.%d.user_type", i)] == expected["user_type"] {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("expected key not found: %+v", expected)
			}
		}
		return nil
	}
}

func testUsersDatasourceV4WithInvalidFilterConfig() string {
	return `
	data "nutanix_users_v2" "test" {
		filter = "username eq 'invalid'"
	}
	`
}
