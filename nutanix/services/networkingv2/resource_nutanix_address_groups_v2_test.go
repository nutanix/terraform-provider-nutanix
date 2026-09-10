package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameAddressGroup = "nutanix_address_groups_v2.test"

func TestAccV2NutanixAddressGroupResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-address-group-%d", r)
	desc := "test address group description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAddressGroupV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "ipv4_addresses.#"),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixAddressGroupResource_WithUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-address-group-%d", r)
	desc := "test address group description"
	updatedName := fmt.Sprintf("tf-test-address-group-%d-updated", r)
	updatedDesc := "test address group description updated"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAddressGroupV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "ipv4_addresses.#"),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "ipv4_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "ext_id"),
				),
			},
			{
				Config: testAddressGroupV2ConfigWithUpdate(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "ipv4_addresses.#"),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "ipv4_addresses.#", "2"),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixAddressGroupResource_WithIPRanges(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-address-group-%d", r)
	desc := "test address group description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAddressGroupV2ConfigWithIPRanges(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "ip_ranges.#"),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "ip_ranges.0.start_ip", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceNameAddressGroup, "ip_ranges.0.end_ip", "10.0.0.10"),
					resource.TestCheckResourceAttrSet(resourceNameAddressGroup, "ext_id"),
				),
			},
		},
	})
}

func testAddressGroupV2Config(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_address_groups_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		ipv4_addresses{
		  value = "10.0.0.0"
		  prefix_length = 24
		}
	  }
`, name, desc)
}

func testAddressGroupV2ConfigWithUpdate(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_address_groups_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		ipv4_addresses{
		  value = "10.0.0.0"
		  prefix_length = 24
		}
		ipv4_addresses{
			value = "172.0.0.0"
			prefix_length = 24
		  }
	  }
`, name, desc)
}

func testAddressGroupV2ConfigWithIPRanges(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_address_groups_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		ip_ranges{
			start_ip = "10.0.0.1"
			end_ip = "10.0.0.10"
		}
	  }
`, name, desc)
}
