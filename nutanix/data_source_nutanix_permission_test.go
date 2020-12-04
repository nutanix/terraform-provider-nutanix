package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	PERMISSIONNAME = "Access_Console_Virtual_Machine"
	PERMISSINOUUID = "16b81a55-2bca-48c6-9fab-4f82c6bb4284"
)

func TestAccNutanixPermissionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionDataSourceConfig(PERMISSINOUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "name", PERMISSIONNAME),
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "operation", "console_access"),
					resource.TestCheckResourceAttr(
						"data.nutanix_permission.test", "fields.0.field_mode", "DISALLOWED"),
				),
			},
		},
	})
}

func testAccPermissionDataSourceConfig(uuid string) string {
	return fmt.Sprintf(`
data "nutanix_permission" "test" {
	permission_id = "%s"
}
`, uuid)
}

// func TestAccNutanixPermissionDataSource_ByName(t *testing.T) {

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccPermissionDataSourceConfigByName(displayName),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(
// 						"data.nutanix_permission.test", "display_name", displayName),
// 					resource.TestCheckResourceAttr(
// 						"data.nutanix_permission.test", "id", uuid),
// 					resource.TestCheckResourceAttrSet("data.nutanix_permission.test", "directory_service_user_group.#"),
// 					resource.TestCheckResourceAttr(
// 						"data.nutanix_permission.test", "directory_service_user_group.0.distinguished_name", distinguishedName),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccPermissionDataSourceConfigByName(dn string) string {
// 	return fmt.Sprintf(`
// data "nutanix_permission" "test" {
// 	user_group_name = "%s"
// }
// `, dn)
// }

// func TestAccNutanixPermissionDataSource_ByDistinguishedName(t *testing.T) {
// 	distinguishedName := "cn=dou-group-1,cn=users,dc=ntnxlab,dc=local"
// 	displayName := "dou-group-1"
// 	uuid := "d12fa0a3-13f1-4f5d-b773-c8e2f8144f0e"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccPermissionDataSourceConfigByDistinguishedName(distinguishedName),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(
// 						"data.nutanix_permission.test", "display_name", displayName),
// 					resource.TestCheckResourceAttr(
// 						"data.nutanix_permission.test", "id", uuid),
// 					resource.TestCheckResourceAttrSet("data.nutanix_permission.test", "directory_service_user_group.#"),
// 					resource.TestCheckResourceAttr(
// 						"data.nutanix_permission.test", "directory_service_user_group.0.distinguished_name", distinguishedName),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccPermissionDataSourceConfigByDistinguishedName(dn string) string {
// 	return fmt.Sprintf(`
// data "nutanix_permission" "test" {
// 	user_group_distinguished_name = "%s"
// }
// `, dn)
// }
