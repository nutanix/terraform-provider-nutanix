package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameCloneRefresh = "nutanix_ndb_clone_refresh.acctest-managed"

func TestAccEra_CloneRefreshbasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraCloneRefreshConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameCloneRefresh, "snapshot_id"),
					resource.TestCheckResourceAttrSet(resourceNameCloneRefresh, "timezone"),
				),
			},
		},
	})
}

func testAccEraCloneRefreshConfig() string {
	return `
        data "nutanix_ndb_clones" "clones"{ }
      
        data "nutanix_ndb_time_machines" "test1" {}

        data "nutanix_ndb_tms_capability" "test"{
            time_machine_id = data.nutanix_ndb_time_machines.test1.time_machines.0.id
        }

        resource "nutanix_ndb_clone_refresh" "acctest-managed"{
            clone_id = data.nutanix_ndb_clones.clones.clones.0.id
            snapshot_id = data.nutanix_ndb_tms_capability.test.capability.1.snapshots.0.id
            timezone = "Asia/Calcutta"
        }
    `
}
