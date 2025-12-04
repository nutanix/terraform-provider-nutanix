package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameLinkedDB = "nutanix_ndb_linked_databases.acctest-managed"

func TestAccEraLinkedDB_basic(t *testing.T) {
	name := "test-linked-tf"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraLinkedDB(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameLinkedDB, "name", name),
					resource.TestCheckResourceAttrSet(resourceNameLinkedDB, "id"),
					resource.TestCheckResourceAttrSet(resourceNameLinkedDB, "status"),
					resource.TestCheckResourceAttrSet(resourceNameLinkedDB, "parent_database_id"),
				),
			},
		},
	})
}

func testAccEraLinkedDB(name string) string {
	return fmt.Sprintf(
		`
		data "nutanix_ndb_databases" "test1" {}

		resource "nutanix_ndb_linked_databases" "acctest-managed" {
			database_id= data.nutanix_ndb_databases.test1.database_instances.0.id
			database_name = "%[1]s"
		  }
		`, name)
}
