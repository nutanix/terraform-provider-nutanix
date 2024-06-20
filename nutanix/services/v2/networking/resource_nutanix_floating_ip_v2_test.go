package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamefip = "nutanix_floating_ip_v4.test"

func TestAccNutanixFloatingIPV4_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	updatedName := fmt.Sprintf("updated-fip-%d", r)
	updatedDesc := "updated fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFloatingIPV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamefip, "name", name),
					resource.TestCheckResourceAttr(resourceNamefip, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "external_subnet_reference"),
				),
			},
			{
				Config: testFloatingIPV4Config(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamefip, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNamefip, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "external_subnet_reference"),
				),
			},
		},
	})
}

func TestAccNutanixFloatingIPV4_WithVmNICAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFloatingIPV4ConfigwithVMNic(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamefip, "name", name),
					resource.TestCheckResourceAttr(resourceNamefip, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "association.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "external_subnet_reference"),
				),
			},
		},
	})
}

func TestAccNutanixFloatingIPV4_WithPrivateipAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFloatingIPV4ConfigwithPrivateIP(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamefip, "name", name),
					resource.TestCheckResourceAttr(resourceNamefip, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "association.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "external_subnet_reference"),
				),
			},
		},
	})
}

func testFloatingIPV4Config(name, desc string) string {
	return fmt.Sprintf(`
		
		resource "nutanix_floating_ip_v4" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
		  }
`, name, desc)
}

func testFloatingIPV4ConfigwithVMNic(name, desc string) string {
	return fmt.Sprintf(`
		
		resource "nutanix_floating_ip_v4" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
			association{
				vm_nic_association{
					vm_nic_reference = "ba209e04-87a7-4dbe-a54b-b0b1e1430e48"
				}
			  }
		  }
`, name, desc)
}

func testFloatingIPV4ConfigwithPrivateIP(name, desc string) string {
	return fmt.Sprintf(`
		
		resource "nutanix_floating_ip_v4" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
			association{
				private_ip_association{
					vpc_reference = "5f79d5e2-5051-4dad-8079-82c2564bb2e1"
					private_ip{
						ipv4{
							value = "10.44.44.7"
						}
					}
				}
			  }
		  }
`, name, desc)
}
