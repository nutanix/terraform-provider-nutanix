package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceTMName = "data.nutanix_ndb_time_machine.test"

func TestAccEraTimeMachineDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTimeMachineDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceTMName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceTMName, "description"),
					resource.TestCheckResourceAttr(dataSourceTMName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTMName, "clone", "false"),
					resource.TestCheckResourceAttr(dataSourceTMName, "sla.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTMName, "schedule.#", "1"),
				),
			},
		},
	})
}

func TestAccEraTimeMachineDataSource_basicWithID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTimeMachineDataSourceConfigWithID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceTMName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceTMName, "description"),
					resource.TestCheckResourceAttr(dataSourceTMName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTMName, "clone", "false"),
					resource.TestCheckResourceAttr(dataSourceTMName, "sla.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTMName, "schedule.#", "1"),
				),
			},
		},
	})
}

func testAccEraTimeMachineDataSourceConfig() string {
	return `
		data "nutanix_ndb_time_machines" "test1" {}

		data "nutanix_ndb_time_machine" "test"{
			time_machine_name = data.nutanix_ndb_time_machines.test1.time_machines.0.name
		}
	`
}

func testAccEraTimeMachineDataSourceConfigWithID() string {
	return `
		data "nutanix_ndb_time_machines" "test1" {}

		data "nutanix_ndb_time_machine" "test"{
			time_machine_id = data.nutanix_ndb_time_machines.test1.time_machines.0.id
		}
	`
}
