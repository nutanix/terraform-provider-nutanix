package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixSubnetsDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetsDataSourceConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetsExists("data.nutanix_subnets.test1"),
				),
			},
		},
	})
}

func testAccCheckNutanixSubnetsExists(n string) resource.TestCheckFunc {
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

func testAccSubnetsDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"

	cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  	}

	vlan_id = %d
	subnet_type = "VLAN"

	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

	dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
		domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		tftp_server_name = "10.250.140.200"
  }
}

data "nutanix_subnet" "test" {
	subnet_id = "${nutanix_subnet.test.id}"
}

data "nutanix_subnets" "test1" {
	metadata {
		length = 1
	}
}`, r, r)
}
