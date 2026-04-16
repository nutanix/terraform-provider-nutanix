package foundation_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFoundationIPMIConfigResource(t *testing.T) {
	name := "ipmi_configure"
	resourcePath := "nutanix_foundation_ipmi_config." + name
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIPMIConfigResource(name, foundationVars.IpmiConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "blocks.0.nodes.0.ipmi_configure_successful", "true"),
					// verify that again apply would again do create due to "ipmi_configure_now" = true
					resource.TestCheckResourceAttr(resourcePath, "blocks.0.nodes.0.ipmi_configure_now", "true"),
					resource.TestCheckResourceAttr(resourcePath, "blocks.0.nodes.#", "1"),
				),
			},
		},
	})
}

func TestAccFoundationIPMIConfigResource_Error(t *testing.T) {
	name := "ipmi_configure"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testIPMIConfigResourceWithWrongPassword(name, foundationVars.IpmiConfig),
				ExpectError: regexp.MustCompile("IPMI config failed for IPMI IP"),
			},
		},
	})
}

func testIPMIConfigResource(name string, i IPMIConfig) string {
	return fmt.Sprintf(`
	resource "nutanix_foundation_ipmi_config" "%[1]s" {
		ipmi_gateway = "%[2]s"
		ipmi_netmask = "%[3]s"
		ipmi_user = "%[4]s"
		ipmi_password = "%[5]s"
		blocks{
			nodes {
				ipmi_mac = "%[6]s"
				ipmi_configure_now = true
				ipmi_ip = "%[7]s"
			}
		}

	}`, name, i.IpmiGateway, i.IpmiNetmask, i.IpmiUser, i.IpmiPassword, i.IpmiMac, i.IpmiIP)
}

func testIPMIConfigResourceWithWrongPassword(name string, i IPMIConfig) string {
	return fmt.Sprintf(`
	resource "nutanix_foundation_ipmi_config" "%[1]s" {
		ipmi_gateway = "%[2]s"
		ipmi_netmask = "%[3]s"
		ipmi_user = "%[4]s"
		ipmi_password = "%[5]s"
		blocks{
			nodes {
				ipmi_mac = "%[6]s"
				ipmi_configure_now = true
				ipmi_ip = "%[7]s"
			}
		}

	}`, name, i.IpmiGateway, i.IpmiNetmask, i.IpmiUser, "ironman", i.IpmiMac, i.IpmiIP)
}
