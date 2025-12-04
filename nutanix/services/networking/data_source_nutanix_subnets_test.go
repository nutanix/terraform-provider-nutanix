package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixSubnetsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_subnets.test", "entities.0.cluster_reference.name"),
				),
			},
		},
	})
}

func TestAccNutanixSubnetsDataSource_WithFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetsDataSourceConfigWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_subnets.test", "entities.0.name", "vlan0_test_2"),
				),
			},
		},
	})
}

func testAccSubnetsDataSourceConfig() string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "test" {
	count = 21
	name = "dou_vlan0_test_${count.index}_%[1]d"
	cluster_uuid = local.cluster1

	vlan_id = count.index + 1
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

data "nutanix_subnets" "test" {}

	`, randIntBetween(1, 25))
}

func testAccSubnetsDataSourceConfigWithFilters() string {
	return `
	data "nutanix_clusters" "clusters" {}
	locals{
		cluster1 = [
			for cluster in data.nutanix_clusters.clusters.entities:
			cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}
	
	resource "nutanix_subnet" "test-subnets" {
		count = 5
		name = "vlan0_test_${count.index}" 
		cluster_uuid = local.cluster1
	
		vlan_id = count.index + 1
		subnet_type = "VLAN"
		
	}
	
	data "nutanix_subnets" "test" {
		metadata {
		  filter = "name==vlan0_test_2"
		}
		depends_on = [
			nutanix_subnet.test-subnets
		]
	}`
}
