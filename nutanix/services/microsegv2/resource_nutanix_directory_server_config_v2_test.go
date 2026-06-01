package microsegv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameDSC = "nutanix_directory_server_config_v2.test"
const datasourceNameDSC = "data.nutanix_directory_server_config_v2.test"
const datasourceNameDSCs = "data.nutanix_directory_server_configs_v2.test"

func TestAccV2NutanixDirectoryServerConfigResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServerConfigV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDSC, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameDSC, "directory_service_reference"),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServerConfigResource_WithUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServerConfigV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDSC, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameDSC, "directory_service_reference"),
				),
			},
			{
				Config: testDirectoryServerConfigV2Update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDSC, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDSC, "is_default_category_enabled", "true"),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServerConfigResource_ImportState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServerConfigV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDSC, "ext_id"),
				),
			},
			{
				ResourceName:      resourceNameDSC,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccV2NutanixDirectoryServerConfigDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServerConfigV2DatasourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameDSC, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameDSC, "directory_service_reference"),
				),
			},
		},
	})
}

func TestAccV2NutanixDirectoryServerConfigsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDirectoryServerConfigsV2DatasourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAttrListNotEmpty(datasourceNameDSCs, "directory_server_configs"),
				),
			},
		},
	})
}

func testCheckAttrListNotEmpty(resourceName, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		attrKey := fmt.Sprintf("%s.#", attr)
		val, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}
		if val == "0" {
			return fmt.Errorf("expected %s to be non-empty", attrKey)
		}
		return nil
	}
}

func testDirectoryServerConfigV2Basic() string {
	return `
    variable "directory_service_ext_id" {
      type    = string
      default = ""
    }

    resource "nutanix_directory_server_config_v2" "test" {
      directory_service_reference = var.directory_service_ext_id
      matching_criterias {
        match_entity = "VM"
        match_field  = "NAME"
        match_type   = "ALL"
      }
    }
  `
}

func testDirectoryServerConfigV2Update() string {
	return `
    variable "directory_service_ext_id" {
      type    = string
      default = ""
    }

    resource "nutanix_directory_server_config_v2" "test" {
      directory_service_reference = var.directory_service_ext_id
      is_default_category_enabled = true
      matching_criterias {
        match_entity = "VM"
        match_field  = "NAME"
        match_type   = "ALL"
      }
    }
  `
}

func testDirectoryServerConfigV2DatasourceBasic() string {
	return `
    variable "directory_service_ext_id" {
      type    = string
      default = ""
    }

    resource "nutanix_directory_server_config_v2" "test" {
      directory_service_reference = var.directory_service_ext_id
      matching_criterias {
        match_entity = "VM"
        match_field  = "NAME"
        match_type   = "ALL"
      }
    }

    data "nutanix_directory_server_config_v2" "test" {
      ext_id = nutanix_directory_server_config_v2.test.ext_id
    }
  `
}

func testDirectoryServerConfigsV2DatasourceBasic() string {
	return `
    variable "directory_service_ext_id" {
      type    = string
      default = ""
    }

    resource "nutanix_directory_server_config_v2" "test" {
      directory_service_reference = var.directory_service_ext_id
      matching_criterias {
        match_entity = "VM"
        match_field  = "NAME"
        match_type   = "ALL"
      }
    }

    data "nutanix_directory_server_configs_v2" "test" {
      depends_on = [nutanix_directory_server_config_v2.test]
    }
  `
}
