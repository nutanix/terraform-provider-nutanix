package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraSLADataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSLADataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test1", "slas.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test1", "slas.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test1", "slas.0.unique_name"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_sla.test", "system_sla", "true"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_sla.test", "yearly_retention", "0"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_sla.test", "pitr_enabled", "true"),
				),
			},
		},
	})
}

func TestAccEraSLADataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSLADataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test1", "slas.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test1", "slas.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_slas.test1", "slas.0.unique_name"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_sla.test", "system_sla", "true"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_sla.test", "yearly_retention", "0"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_sla.test", "pitr_enabled", "true"),
				),
			},
		},
	})
}

func testAccEraSLADataSourceConfig() string {
	return `
		data "nutanix_ndb_slas" "test1" {}

		data "nutanix_ndb_sla" "test"{
			sla_id = data.nutanix_ndb_slas.test1.slas.0.id
		}
	`
}

func testAccEraSLADataSourceConfigByName() string {
	return `
		data "nutanix_ndb_slas" "test1" {}

		data "nutanix_ndb_sla" "test"{
			sla_name = data.nutanix_ndb_slas.test1.slas.0.name
		}
	`
}
