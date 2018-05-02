package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixSubnetDataSource_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.nutanix_subnet", "prefix_length", "24"),
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.nutanix_subnet", "subnet_type", "VLAN"),
				),
			},
		},
	})
}

func testAccSubnetDataSourceConfig(r int) string {
	return fmt.Sprintf(`
variable clusterid {
	default = "000567f3-1921-c722-471d-0cc47ac31055"
}

resource "nutanix_subnet" "test" {
	metadata = {
		kind = "subnet"
	}

	name = "dou_vlan0_test_%d"
	description = "Dou Vlan 0"

	cluster_reference = {
	  kind = "cluster"
	  uuid = "${var.clusterid}"
  	}

	vlan_id = 201
	subnet_type = "VLAN"
	
	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]
	
	dhcp_options {
		boot_file_name = "bootfile"
		tftp_server_name = "192.168.0.252"
		domain_name = "nutanix"
	}

	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list = ["nutanix.com", "calm.io"]
	
}
`, r)
}
