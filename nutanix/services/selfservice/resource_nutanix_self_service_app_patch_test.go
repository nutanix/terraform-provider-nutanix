package selfservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamePatch = "nutanix_self_service_app_patch.test"

func TestAccNutanixCalmAppVmUpdateResource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Testing for patch config"
	configNameBasic := "VmUpdate"
	configNameEditable := "VmUpdateEditable"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionWithUpdateConfig(testVars.SelfService.BlueprintName, name, desc) + testCalmAppVMUpdateBasic(configNameBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", configNameBasic),
				),
			},
			{
				Config: testCalmAppProvisionWithUpdateConfig(testVars.SelfService.BlueprintName, name, desc) + testCalmAppVMUpdateEditable(configNameEditable),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", configNameEditable),
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
	categoryAddConfig := "CategoriesAdd"
	categoryDeleteConfig := "CategoriesDelete"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionWithUpdateConfig(testVars.SelfService.BlueprintName, name, desc) + testCalmAppCategoryAdd(categoryAddConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", categoryAddConfig),
				),
			},
			{
				Config: testCalmAppProvisionWithUpdateConfig(testVars.SelfService.BlueprintName, name, desc) + testCalmAppCategoryDelete(categoryDeleteConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", categoryDeleteConfig),
				),
			},
		},
	})
}

func TestAccNutanixCalmAppDiskAddResource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-app-%d", r)
	desc := "Testing for patch config"
	diskAddConfig := "DiskAdd"
	diskAddConfigEditables := "DiskAddEditables"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionWithUpdateConfig(testVars.SelfService.BlueprintName, name, desc) + testCalmAppDiskAddBasic(diskAddConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", diskAddConfig),
				),
			},
			{
				Config: testCalmAppProvisionWithUpdateConfig(testVars.SelfService.BlueprintName, name, desc) + testCalmAppDiskAddEditable(diskAddConfigEditables),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", diskAddConfigEditables),
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
	nicAddConfig := "NicAdd"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCalmAppProvisionWithUpdateConfig(testVars.SelfService.BlueprintName, name, desc) + testCalmAppNicAdd(nicAddConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatch, "config_name", nicAddConfig),
				),
			},
		},
	})
}

func testCalmAppProvisionWithUpdateConfig(blueprintName, name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_provision" "test" {
		bp_name         = "%[1]s"
		app_name        = "%[2]s"
		app_description = "%[3]s"
		}
`, blueprintName, name, desc)
}

func testCalmAppVMUpdateBasic(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_patch" "test" {
		app_uuid = nutanix_self_service_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
	}
`, name)
}

func testCalmAppVMUpdateEditable(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_patch" "test" {
		app_uuid = nutanix_self_service_app_provision.test.id
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
		resource "nutanix_self_service_app_patch" "test" {
		app_uuid = nutanix_self_service_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
		categories {
			value = "AppType:Default"
			operation = "add"
  		}
	}
`, name)
}

func testCalmAppCategoryDelete(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_patch" "test" {
		app_uuid = nutanix_self_service_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
		categories {
			value = "AppType:Default"
			operation = "delete"
		}
	}
`, name)
}

func testCalmAppDiskAddBasic(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_patch" "test" {
		app_uuid = nutanix_self_service_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
	}
`, name)
}

func testCalmAppDiskAddEditable(name string) string {
	return fmt.Sprintf(`
		resource "nutanix_self_service_app_patch" "test" {
		app_uuid = nutanix_self_service_app_provision.test.id
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
		resource "nutanix_self_service_app_patch" "test" {
		app_uuid = nutanix_self_service_app_provision.test.id
		patch_name = "%[1]s"
		config_name = "%[1]s"
	}
`, name)
}
