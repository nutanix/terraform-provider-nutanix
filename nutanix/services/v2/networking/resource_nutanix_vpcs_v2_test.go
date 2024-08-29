package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVpc = "nutanix_vpc_v2.test"

func TestAccNutanixVpcV2_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vpc-%d", r)
	desc := "test vpc description"
	updatedName := fmt.Sprintf("updated-vpc-%d", r)
	updatedDesc := "updated vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
				),
			},
			{
				Config: testVpcConfig(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
				),
			},
		},
	})
}

func TestAccNutanixVpcV2_WithExternallyRoutablePrefixes(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vpc-%d", r)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithExtRoutablePrefix(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
				),
			},
		},
	})
}

func TestAccNutanixVpcV2_WithDHCP(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vpc-%d", r)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithDHCP(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "common_dhcp_options.#"),
				),
			},
		},
	})
}

func TestAccNutanixVpcV2_WithTransitType(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vpc-%d", r)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithTransitType(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "TRANSIT"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "common_dhcp_options.#"),
				),
			},
		},
	})
}

func testVpcConfig(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets{
		  subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
		}
	}
`, name, desc)
}

func testVpcConfigWithExtRoutablePrefix(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets{
		  subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
		  external_ips{
			ipv4{
			  value = "10.44.44.6"
			  prefix_length = 32
			}
		  }
		}
		externally_routable_prefixes{
		  ipv4{
			ip{
			  value = "172.30.0.0"
			  prefix_length = 32
			}
			prefix_length = 16
		  }
		}
	}
`, name, desc)
}

func testVpcConfigWithDHCP(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets{
		  subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
		}
		common_dhcp_options{
			domain_name_servers{
				ipv4{
					value = "8.8.8.9"
					prefix_length = 32
				}
			}
			domain_name_servers{
				ipv4{
					value = "8.8.8.8"
					prefix_length = 32
				}
			}
		}	
	}
`, name, desc)
}

func testVpcConfigWithTransitType(name, desc string) string {
	return fmt.Sprintf(`
		
	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets{
		  subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
		}
		vpc_type = "TRANSIT"
	}
`, name, desc)
}
