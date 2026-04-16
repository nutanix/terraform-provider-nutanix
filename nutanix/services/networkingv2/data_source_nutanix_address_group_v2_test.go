package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameaddgrp = "data.nutanix_address_group_v2.test"

func TestAccV2NutanixAddressGroupDataSource_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-service-%d", r)
	desc := "test service description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressGrpDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameaddgrp, "name", name),
					resource.TestCheckResourceAttr(datasourceNameaddgrp, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameaddgrp, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameaddgrp, "ipv4_addresses.#"),
					resource.TestCheckResourceAttr(datasourceNameaddgrp, "ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameaddgrp, "ext_id"),
				),
			},
		},
	})
}

func testAccAddressGrpDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`

	resource "nutanix_address_groups_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		ipv4_addresses{
		  value = "10.0.0.0"
		  prefix_length = 24
		}
	  }

		data "nutanix_address_group_v2" "test" {
			ext_id = nutanix_address_groups_v2.test.ext_id
			depends_on = [
				resource.nutanix_address_groups_v2.test
			]
		}
	`, name, desc)
}
