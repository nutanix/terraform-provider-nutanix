package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameSLA = "nutanix_ndb_sla.acctest-managed"

func TestAccEra_Slabasic(t *testing.T) {
	name := "test-sla-tf"
	desc := "this is sla desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSLAConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSLA, "name", name),
					resource.TestCheckResourceAttr(resourceNameSLA, "description", desc),
					resource.TestCheckResourceAttr(resourceNameSLA, "continuous_retention", "30"),
					resource.TestCheckResourceAttr(resourceNameSLA, "daily_retention", "3"),
					resource.TestCheckResourceAttr(resourceNameSLA, "weekly_retention", "2"),
					resource.TestCheckResourceAttr(resourceNameSLA, "monthly_retention", "1"),
					resource.TestCheckResourceAttr(resourceNameSLA, "quarterly_retention", "1"),
				),
			},
		},
	})
}

func TestAccEra_SlaUpdate(t *testing.T) {
	name := "test-sla-tf"
	desc := "this is sla desc"
	updatedName := "test-sla-updated"
	updatedDesc := "desc is updated"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraSLAConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSLA, "name", name),
					resource.TestCheckResourceAttr(resourceNameSLA, "description", desc),
					resource.TestCheckResourceAttr(resourceNameSLA, "continuous_retention", "30"),
					resource.TestCheckResourceAttr(resourceNameSLA, "daily_retention", "3"),
					resource.TestCheckResourceAttr(resourceNameSLA, "weekly_retention", "2"),
					resource.TestCheckResourceAttr(resourceNameSLA, "monthly_retention", "1"),
					resource.TestCheckResourceAttr(resourceNameSLA, "quarterly_retention", "1"),
				),
			},
			{
				Config: testAccEraSLAConfigUpdated(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSLA, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameSLA, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameSLA, "continuous_retention", "25"),
					resource.TestCheckResourceAttr(resourceNameSLA, "daily_retention", "1"),
					resource.TestCheckResourceAttr(resourceNameSLA, "weekly_retention", "3"),
					resource.TestCheckResourceAttr(resourceNameSLA, "monthly_retention", "1"),
					resource.TestCheckResourceAttr(resourceNameSLA, "quarterly_retention", "3"),
				),
			},
		},
	})
}

func testAccEraSLAConfig(name, desc string) string {
	return fmt.Sprintf(`
	resource "nutanix_ndb_sla" "acctest-managed" {
		name= "%[1]s"
		description = "%[2]s"
		continuous_retention = 30
		daily_retention = 3
		weekly_retention = 2
		monthly_retention= 1
		quarterly_retention=1
	  }
	`, name, desc)
}

func testAccEraSLAConfigUpdated(name, desc string) string {
	return fmt.Sprintf(`
	resource "nutanix_ndb_sla" "acctest-managed" {
		name= "%[1]s"
		description = "%[2]s"
		continuous_retention = 25
		daily_retention = 1
		weekly_retention = 3
		monthly_retention= 1
		quarterly_retention=3
	  }
	`, name, desc)
}
