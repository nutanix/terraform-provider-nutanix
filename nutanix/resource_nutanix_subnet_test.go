package nutanix

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixSubnet_basic(t *testing.T) {
	r := rand.Int31()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNutanixSubnetConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists("nutanix_subnet.test"),
				),
			},
		},
	})
}

func testAccCheckNutanixSubnetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixSubnetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*NutanixClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_subnet" {
			continue
		}
		for {
			_, err := conn.API.V3.GetSubnet(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}

	}

	return nil
}

func testAccNutanixSubnetConfig(r int32) string {
	return fmt.Sprintf(`
provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
	insecure = true
	port = 9440
}

resource "nutanix_subnet" "test" {
	metadata = {
		kind = "subnet"
	}

	name = "dou_vlan0_test_%d"
	description = "Dou Vlan 0"

	cluster_reference = {
	  kind = "cluster"
	  uuid = "000567f3-1921-c722-471d-0cc47ac31055" 
  	}

	vlan_id = 0
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
