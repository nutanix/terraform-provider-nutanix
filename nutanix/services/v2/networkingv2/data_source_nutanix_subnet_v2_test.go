package networkingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSubnet = "data.nutanix_subnet_v2.test"

func TestAccNutanixSubnetDataSourceV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameSubnet, "is_external", "true"),
					resource.TestCheckResourceAttr(datasourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttrSet(datasourceNameSubnet, "cluster_reference"),
					resource.TestCheckResourceAttrSet(datasourceNameSubnet, "links.#"),
				),
			},
		},
	})
}

func testAccSubnetDataSourceConfig() string {
	return (`

		data "nutanix_clusters" "clusters" {}

		locals {
			cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
		
		resource "nutanix_subnet_v2" "test" {
			name = "terraform_test_subnets_datasource"
			description = "terraform test subnets datasource description"
			cluster_reference = local.cluster0
			subnet_type = "VLAN"
			network_id = 112
			is_external = true
			ip_config {
				ipv4 {
					ip_subnet {
						ip {
							value = "192.168.0.0"
						}
						prefix_length = 24
					}
					default_gateway_ip {
						value = "192.168.0.1"
					}
					pool_list{
						start_ip {
							value = "192.168.0.20"
						}
						end_ip {
							value = "192.168.0.30"
						}
					}
				}
			}
			depends_on = [data.nutanix_clusters.clusters]
		}

		data "nutanix_subnet_v2" "test" {
			ext_id = nutanix_subnet_v2.test.id
		}
`)
}
