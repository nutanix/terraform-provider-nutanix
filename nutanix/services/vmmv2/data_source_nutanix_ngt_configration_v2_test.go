package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameNGTConfiguration = "data.nutanix_ngt_configuration_v2.test"

func TestAccV2NutanixNGTConfigurationDatasource_GetNGTConfigurationForVM(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// get NGT configuration for the VM
			{
				Config: testNGTConfigurationDatasource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNGTConfiguration, "guest_os_version"),
					resource.TestCheckResourceAttrSet(datasourceNameNGTConfiguration, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_installed", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_reachable", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_enabled", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.#", "2"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccV2NutanixNGTConfigurationV4Datasource_GetNGTConfigurationForVM_NGTNotInstalled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNGTConfigurationNotInstalledDatasource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_enabled", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_installed", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_reachable", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_vm_mobility_drivers_installed", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_vss_snapshot_capable", "false"),
				),
			},
		},
	})
}

func testNGTConfigurationDatasource() string {
	return fmt.Sprintf(`
	locals {
	  config            = jsondecode(file("%[1]s"))
	  preEnv            = local.config.pre_env
	}

	data "nutanix_virtual_machines_v2" "test" {
		filter = "name eq '${local.preEnv.ngt_vm.name}'"
	}
	data "nutanix_ngt_configuration_v2" "test" {
		ext_id = data.nutanix_virtual_machines_v2.test.vms[0].ext_id
	}
  `, filepath)
}

func testNGTConfigurationNotInstalledDatasource() string {
	return fmt.Sprintf(`
	locals {
	  config            = jsondecode(file("%[1]s"))
	  preEnv            = local.config.pre_env
	}

	data "nutanix_virtual_machines_v2" "test" {
		filter = "name eq '${local.preEnv.integration_vm.name}'"
	}
	data "nutanix_ngt_configuration_v2" "test" {
		ext_id = data.nutanix_virtual_machines_v2.test.vms[0].ext_id
	}
  `, filepath)
}
