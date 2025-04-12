package selfservice_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamePatch = "nutanix_calm_app_patch.test"

func TestAccNutanixCalmAppVmUpdateResource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Testing for patch config"
	config_name_basic := "VmUpdate"
	config_name_editable := "VmUpdateEditable"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionWithUpdateConfig(name, desc) + testCalmAppVmUpdateBasic(config_name_basic),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						aJSON, _ := json.MarshalIndent(s.RootModule().Resources[resourceNamePatch].Primary.Attributes, "", "  ")
						fmt.Printf("################### %s #########################\n", resourceNamePatch)
						fmt.Printf("Resource Attributes: \n%v\n", string(aJSON))
						fmt.Printf("\n############################################\n")
						return nil
					},
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", config_name_basic),
				),
			},
			{
				Config: testCalmAppProvisionWithUpdateConfig(name, desc) + testCalmAppVmUpdateEditable(config_name_editable),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						aJSON, _ := json.MarshalIndent(s.RootModule().Resources[resourceNamePatch].Primary.Attributes, "", "  ")
						fmt.Printf("################### %s #########################\n", resourceNamePatch)
						fmt.Printf("Resource Attributes: \n%v\n", string(aJSON))
						fmt.Printf("\n############################################\n")
						return nil
					},
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", config_name_editable),
					resource.TestCheckResourceAttr(resourceNamePatch, "vm_config.0.memory_size_mib", "2048"),
					resource.TestCheckResourceAttr(resourceNamePatch, "vm_config.0.num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNamePatch, "vm_config.0.num_vcpus_per_socket", "2"),
				),
			},
		},
	})
}

func TestAccNutanixCalmAppCategoryUpdateResource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Testing for patch config"
	category_add_config := "CategoriesAdd"
	category_delete_config := "CategoriesDelete"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionWithUpdateConfig(name, desc) + testCalmAppCategoryAdd(category_add_config),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						aJSON, _ := json.MarshalIndent(s.RootModule().Resources[resourceNamePatch].Primary.Attributes, "", "  ")
						fmt.Printf("################### %s #########################\n", resourceNamePatch)
						fmt.Printf("Resource Attributes: \n%v\n", string(aJSON))
						fmt.Printf("\n############################################\n")
						return nil
					},
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", category_add_config),
				),
			},
			{
				Config: testCalmAppProvisionWithUpdateConfig(name, desc) + testCalmAppCategoryDelete(category_delete_config),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						aJSON, _ := json.MarshalIndent(s.RootModule().Resources[resourceNamePatch].Primary.Attributes, "", "  ")
						fmt.Printf("################### %s #########################\n", resourceNamePatch)
						fmt.Printf("Resource Attributes: \n%v\n", string(aJSON))
						fmt.Printf("\n############################################\n")
						return nil
					},
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", category_delete_config),
				),
			},
		},
	})
}

func TestAccNutanixCalmAppDiskAddResource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Testing for patch config"
	disk_add_config := "DiskAdd"
	disk_add_config_editables := "DiskAddEditables"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionWithUpdateConfig(name, desc) + testCalmAppDiskAddBasic(disk_add_config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", disk_add_config),
				),
			},
			{
				Config: testCalmAppProvisionWithUpdateConfig(name, desc) + testCalmAppDiskAddEditable(disk_add_config_editables),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", disk_add_config_editables),
					resource.TestCheckResourceAttr(resourceNamePatch, "disks.0.operation", "add"),
					resource.TestCheckResourceAttr(resourceNamePatch, "disks.0.disk_size_mib", "3072"),
				),
			},
		},
	})
}

func TestAccNutanixCalmAppNicAddResource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Testing for patch config"
	nic_add_config := "NicAdd"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionWithUpdateConfig(name, desc) + testCalmAppNicAdd(nic_add_config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", nic_add_config),
				),
			},
		},
	})
}

func testCalmAppProvisionWithUpdateConfig(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_provision" "test" {
		bp_name         = "demo_bp"
		app_name        = "%[1]s"
		app_description = "%[2]s"
		}
`, name, desc)
}

func testCalmAppVmUpdateBasic(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_patch" "test" {
		app_uuid = nutanix_calm_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
	}
`, name)
}

func testCalmAppVmUpdateEditable(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_patch" "test" {
		app_uuid = nutanix_calm_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
		vm_config {
			memory_size_mib = 2048
			num_sockets = 2
			num_vcpus_per_socket = 2
		}
	}
`, name)
}

func testCalmAppCategoryAdd(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_patch" "test" {
		app_uuid = nutanix_calm_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
		categories {
		value = "AppType: Default"
		operation = "add"
  		}
	}
`, name)
}

func testCalmAppCategoryDelete(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_patch" "test" {
		app_uuid = nutanix_calm_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
		categories {
			value = "AppType: Default"
			operation = "delete"
		}
	}
`, name)
}

func testCalmAppDiskAddBasic(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_patch" "test" {
		app_uuid = nutanix_calm_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
	}
`, name)
}

func testCalmAppDiskAddEditable(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_patch" "test" {
		app_uuid = nutanix_calm_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
		disks {
			disk_size_mib = 3072
			operation = "add"
  		}
	}
`, name)
}

func testCalmAppNicAdd(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_calm_app_patch" "test" {
		app_uuid = nutanix_calm_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
	}
`, name)
}
