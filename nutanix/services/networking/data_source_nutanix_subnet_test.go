package networking_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixSubnetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfig(randIntBetween(1, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "prefix_length", "24"),
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "subnet_type", "VLAN"),
					resource.TestCheckResourceAttrSet("data.nutanix_subnet.test", "cluster_name"),
				),
			},
		},
	})
}

func TestAccNutanixSubnetDataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfigName(randIntBetween(11, 20)),
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
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccSubnetDataSourceBadConfig(),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
		},
	})
}

func TestAccNutanixSubnetDataSource_nameWithFilters(t *testing.T) {
	r := randIntBetween(11, 20)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfigNameWithFilters(r, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "subnet_type", "VLAN"),
				),
			},
			{
				Config:      testAccSubnetDataSourceConfigNameWithFilters(r, 30),
				ExpectError: regexp.MustCompile("subnet with the given name, not found"),
			},
		},
	})
}

func testAccSubnetDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"
	cluster_uuid = local.cluster1

	vlan_id = %d
	subnet_type = "VLAN"

	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

	dhcp_options = {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

data "nutanix_subnet" "test" {
	subnet_id = nutanix_subnet.test.id
}
`, r, r)
}

func testAccSubnetDataSourceConfigName(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"
	cluster_uuid = local.cluster1
	vlan_id = %d
	subnet_type = "VLAN"

	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	ip_config_pool_list_ranges = ["192.168.0.10 192.168.0.100"]

	dhcp_options = {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

data "nutanix_subnet" "test" {
	subnet_name = nutanix_subnet.test.name
}
`, r, r)
}

func testAccSubnetDataSourceConfigNameWithFilters(r int, f int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "test" {
	name = "dnd_acc_test_vlan0_test_%d"
	cluster_uuid = local.cluster1
	vlan_id = %d
	subnet_type = "VLAN"
}

data "nutanix_subnet" "test" {
	subnet_name = nutanix_subnet.test.name
	additional_filter {
		name = "vlan_id"
		values = ["%d"]
	}
}
`, r, r, f)
}

func testAccSubnetDataSourceBadConfig() string {
	return `
data "nutanix_subnet" "test" {
	subnet_id   = "test-subnet-id"
	subnet_name = "test-subnet-name"
}
`
}
