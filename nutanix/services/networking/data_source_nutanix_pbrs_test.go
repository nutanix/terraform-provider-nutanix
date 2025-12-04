package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixPbrsDataSource_basic(t *testing.T) {
	r := randIntBetween(151, 160)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPbrsDataSourceConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_pbrs.test", "entities.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_pbrs.test", "api_version"),
					resource.TestCheckResourceAttr("data.nutanix_pbrs.test", "entities.0.spec.0.resources.0.action.0.action", "DENY"),
					resource.TestCheckResourceAttr("data.nutanix_pbrs.test", "entities.0.spec.0.resources.0.protocol_type", "ALL"),
					resource.TestCheckResourceAttr("data.nutanix_pbrs.test", "entities.0.spec.0.resources.0.priority", "1"),
				),
			},
		},
	})
}

func testAccPbrsDataSourceConfig(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
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
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.35.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "pbr-test" {
		name = "acctest-%[1]d"
		priority = %[1]d
		protocol_type = "ALL"
		action = "PERMIT"
		vpc_reference_uuid = resource.nutanix_vpc.test-vpc.id
		source{
		  address_type = "ALL"
		}
		destination{
		  address_type = "ALL"
		}
	}

	data "nutanix_pbrs" "test"{
		depends_on = [
			resource.nutanix_pbr.pbr-test
		]
	}
	`, r)
}
