package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceTMName = "data.nutanix_ndb_time_machine.test"

func TestAccEraTimeMachineDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTimeMachineDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceTMName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceTMName, "description"),
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
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTimeMachineDataSourceConfigWithID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceTMName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceTMName, "description"),
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
