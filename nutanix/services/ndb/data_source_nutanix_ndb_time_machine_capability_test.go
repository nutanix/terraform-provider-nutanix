package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNDBTmsCapability = "data.nutanix_ndb_tms_capability.test"

func TestAccEraTmsCapabilityDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraTmsCapabilityDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNDBTmsCapability, "output_time_zone"),
					resource.TestCheckResourceAttrSet(dataSourceNDBTmsCapability, "type"),
					resource.TestCheckResourceAttrSet(dataSourceNDBTmsCapability, "nx_cluster_id"),
					resource.TestCheckResourceAttrSet(dataSourceNDBTmsCapability, "sla_id"),
					resource.TestCheckResourceAttrSet(dataSourceNDBTmsCapability, "capability.#"),
					resource.TestCheckResourceAttrSet(dataSourceNDBTmsCapability, "capability.0.mode"),
				),
			},
		},
	})
}

func testAccEraTmsCapabilityDataSourceConfig() string {
	return `
		data "nutanix_ndb_time_machines" "test1" {}

		data "nutanix_ndb_tms_capability" "test"{
			time_machine_id = data.nutanix_ndb_time_machines.test1.time_machines.0.id
		}
	`
}
