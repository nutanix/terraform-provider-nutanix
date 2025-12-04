package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceTMsName = "data.nutanix_ndb_time_machines.test"

func TestAccEraTimeMachinesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTimeMachinesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceTMsName, "time_machines.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceTMsName, "time_machines.0.description"),
					resource.TestCheckResourceAttr(dataSourceTMsName, "time_machines.0.clone", "false"),
					resource.TestCheckResourceAttr(dataSourceTMsName, "time_machines.0.sla.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTMsName, "time_machines.0.schedule.#", "1"),
				),
			},
		},
	})
}

func testAccEraTimeMachinesDataSourceConfig() string {
	return `
		data "nutanix_ndb_time_machines" "test" {}
	`
}
