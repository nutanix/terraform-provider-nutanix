package nutanix

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixSubnetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "prefix_length", "24"),
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "subnet_type", "VLAN"),
				),
			},
		},
	})
}

func TestAccNutanixSubnetDataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfigName(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "prefix_length", "24"),
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "subnet_type", "VLAN"),
				),
			},
		},
	})
}

func TestAccNutanixSubnetDataSource_conflicts(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccSubnetDataSourceConfigNameDuplicated(acctest.RandIntRange(0, 500)),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
		},
	})
}

func testAccSubnetDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"
	cluster_uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"

	vlan_id = %d
	subnet_type = "VLAN"

	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

	dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

data "nutanix_subnet" "test" {
	subnet_id = "${nutanix_subnet.test.id}"
}
`, r, r)
}

func testAccSubnetDataSourceConfigName(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"
	cluster_uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
	vlan_id = %d
	subnet_type = "VLAN"

	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

	dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

data "nutanix_subnet" "test" {
	subnet_name = "${nutanix_subnet.test.name}"
}
`, r, r)
}

func testAccSubnetDataSourceConfigNameDuplicated(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"
	cluster_uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
	vlan_id = %d
	subnet_type = "VLAN"

	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

	dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

resource "nutanix_subnet" "test1" {
	name = "${nutanix_subnet.test.name}"
	cluster_uuid= "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
	vlan_id = %d
	subnet_type = "VLAN"
	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

	dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

data "nutanix_subnet" "test" {
	subnet_id   = "${nutanix_subnet.test1.id}"
	subnet_name = "${nutanix_subnet.test1.name}"
}
`, r, r, r+2)
}
