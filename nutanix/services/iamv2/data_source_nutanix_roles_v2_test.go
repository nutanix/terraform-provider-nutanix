package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRoles = "data.nutanix_roles_v2.test"

func TestAccNutanixRolesV2Datasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRolesDatasourceV4Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRoles, "roles.#"),
					resource.TestCheckResourceAttrSet(datasourceNameRoles, "roles.0.display_name"),
					resource.TestCheckResourceAttrSet(datasourceNameRoles, "roles.0.operations.#"),
				),
			},
		},
	})
}

func TestAccNutanixRolesV2Datasource_WithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRolesDatasourceV4WithFilterConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRoles, "roles.#"),
					resource.TestCheckResourceAttr(datasourceNameRoles, "roles.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameRoles, "roles.0.display_name", testVars.Iam.Roles.DisplayName),
					resource.TestCheckResourceAttr(datasourceNameRoles, "roles.0.description", testVars.Iam.Roles.Description),
				),
			},
		},
	})
}

func TestAccNutanixRolesV2Datasource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRolesDatasourceV4WithLimitConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRoles, "roles.#"),
					resource.TestCheckResourceAttr(datasourceNameRoles, "roles.#", "1"),
				),
			},
		},
	})
}

func testRolesDatasourceV4Config() string {
	return `
	data "nutanix_roles_v2" "test"{}
	`
}

func testRolesDatasourceV4WithFilterConfig(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		roles = local.config.iam.roles
	}

	data "nutanix_operations_v2" "test" {
	  filter = "startswith(displayName, 'Create_')"
	}
	
	resource "nutanix_roles_v2" "test" {
		display_name = local.roles.display_name
		description  = local.roles.description
		operations = [
			data.nutanix_operations_v2.test.operations[0].ext_id,
			data.nutanix_operations_v2.test.operations[1].ext_id,
			data.nutanix_operations_v2.test.operations[2].ext_id,
			data.nutanix_operations_v2.test.operations[3].ext_id
	  	]
		depends_on = [data.nutanix_operations_v2.test]
  	}
	  
	  data "nutanix_roles_v2" "test" {
		filter     = "displayName eq '${local.roles.display_name}'"
		depends_on = [resource.nutanix_roles_v2.test]
	  }
	`, filepath)
}

func testRolesDatasourceV4WithLimitConfig(filepath string) string {
	return fmt.Sprintf(`
		locals{
			config = (jsondecode(file("%s")))
			roles = local.config.iam.roles
		}

		data "nutanix_operations_v2" "test" {
		  filter = "startswith(displayName, 'Create_')"
		}
		resource "nutanix_roles_v2" "test" {
			display_name = local.roles.display_name
			description  = local.roles.description
			operations = [
				data.nutanix_operations_v2.test.operations[0].ext_id,
				data.nutanix_operations_v2.test.operations[1].ext_id,
				data.nutanix_operations_v2.test.operations[2].ext_id,
				data.nutanix_operations_v2.test.operations[3].ext_id
			]
			depends_on = [data.nutanix_operations_v2.test]
		}

		data "nutanix_roles_v2" "test" {
			limit     = 1
			depends_on = [resource.nutanix_roles_v2.test]
		}
	`, filepath)
}
