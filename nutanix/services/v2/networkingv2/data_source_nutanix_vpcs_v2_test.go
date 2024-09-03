package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNamevpcs = "data.nutanix_vpcs_v2.test"

func TestAccNutanixVpcsDataSourceV2_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vpc-%d", r)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.#"),
					resource.TestCheckResourceAttr(datasourceNamevpcs, "vpcs.0.name", name),
					resource.TestCheckResourceAttr(datasourceNamevpcs, "vpcs.0.description", desc),
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.0.metadata.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.0.snat_ips.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.0.external_subnets.#"),
				),
			},
		},
	})
}

func testAccVpcsDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`

		data "nutanix_clusters" "clusters" {}

		locals {
			cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
		
		resource "nutanix_subnet_v2" "test" {
			name = "terraform-test-subnet-vpc"
			description = "test subnet description"
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
		resource "nutanix_vpc_v2" "rtest" {
			name =  "%[1]s"
			description = "%[2]s"
			external_subnets{
			  subnet_reference = nutanix_subnet_v2.test.id
			}
			depends_on = [nutanix_subnet_v2.test]
		}

		data "nutanix_vpcs_v2" "test" {
			filter = "name eq '%[1]s'"
			depends_on = [
				resource.nutanix_vpc_v2.rtest
			]
		}
	`, name, desc)
}
