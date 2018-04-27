package nutanix

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixSubnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNutanixVMConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists("nutanix_virtual_machine.vm1"),
					resource.TestCheckResourceAttrSet("nutanix_subnet.resource", "ip_config"),
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
	for i := range s.RootModule().Resources {
		if s.RootModule().Resources[i].Type != "nutanix_subnet" {
			continue
		}
		id := string(s.RootModule().Resources[i].Primary.ID)
		if id == "" {
			err := errors.New("ID is already set to the null string")
			return err
		}
		return nil
	}
	return nil
}

const testAccNutanixSubnetConfig = `
provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
	insecure = true
	port = 9440
}

resource "nutanix_subnet" "my-image" {
	metadata = {
		kind = "subnet"
	}

	name = "sarath_vlan0"
	description = "Sarath Vlan 0"
	resources = {
		vlan_id = 0 
		subnet_type = "VLAN"
		ip_config {
			prefix_length = 24
			default_gateway_ip = "192.168.0.1"
			pool_list = [
				{range = "192.168.0.5"},
				{range = 192.168.0.100"}
			]
			subnet_ip = "192.168.0.0"
			dhcp_options {
				boot_file_name = "bootfile"
				domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
				domain_search_list = ["nutanix.com", "calm.io"]
				tftp_server_name = "192.168.0.252"
				domain_name = "nutanix"
			}
    	}
		
	}  
}


`
