package nutanix

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixSubnet_basic(t *testing.T) {
	r := acctest.RandIntRange(3500, 3900)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists("nutanix_subnet.acctest-managed"),
				),
			},
			{
				Config: testAccNutanixSubnetConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists("nutanix_subnet.acctest-managed"),
					resource.TestCheckResourceAttr("nutanix_subnet.acctest-managed", "name", "acctest-managed-updateName"),
				),
			},
		},
	})
}

func testAccCheckNutanixSubnetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixSubnetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_subnet" {
			continue
		}
		if _, err := resourceNutanixSubnetExists(conn.API, rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccNutanixSubnetConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "acctest-managed" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed"
	description = "here's a good description of my unit test VLAN"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
		domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		tftp_server_name = "10.250.140.200"
  }
}
`, r)
}

func testAccNutanixSubnetConfigUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "acctest-managed" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed-updateName"
	description = "here's a good description of my unit test VLAN"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
		domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		tftp_server_name = "10.250.140.200"
  }
}
`, r)
}
