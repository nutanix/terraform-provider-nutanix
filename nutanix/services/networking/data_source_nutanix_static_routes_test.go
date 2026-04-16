package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixStaticRoutesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRoutesDataSourceConfig(randIntBetween(25, 45)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_static_routes.test", "spec.0.resources.0.static_routes_list.0.destination", "10.2.2.0/24"),
					resource.TestCheckResourceAttr(
						"data.nutanix_static_routes.test", "status.0.resources.0.static_routes_list.0.destination", "10.2.2.0/24"),
				),
			},
		},
	})
}

func testAccStaticRoutesDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "ext-sub" {
	cluster_uuid = local.cluster1
	name        = "acctest-managed-%[1]d"
	description = "Description of my unit test VLAN"
	vlan_id     = %[1]d
	subnet_type = "VLAN"
	subnet_ip          = "10.250.140.0"
	default_gateway_ip = "10.250.140.1"
	prefix_length = 24
	is_external = true
	ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
}

resource "nutanix_vpc" "test-vpc" {
	name = "acctest-managed-%[1]d"
  
  
	external_subnet_reference_uuid = [
		resource.nutanix_subnet.ext-sub.id
	]
  
	common_domain_name_server_ip_list{
			ip = "8.8.8.9"
	}
  
	externally_routable_prefix_list{
	  ip=  "176.40.0.0"
	  prefix_length= 16
	}
}

resource "nutanix_static_routes" "acctest-managed"{
	vpc_uuid = resource.nutanix_vpc.test-vpc.id
	static_routes_list{
		destination= "10.2.2.0/24"
		external_subnet_reference_uuid = resource.nutanix_subnet.ext-sub.id
	}
}

data "nutanix_static_routes" "test" {
	vpc_reference_uuid = resource.nutanix_vpc.test-vpc.id

	depends_on = [
		resource.nutanix_static_routes.acctest-managed
	]
}
`, r)
}
