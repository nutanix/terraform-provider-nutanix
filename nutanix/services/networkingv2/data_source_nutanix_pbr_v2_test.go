package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNamePbr = "data.nutanix_pbr_v2.test"

func TestAccV2NutanixPbrDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPbrDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNamePbr, "name", name),
					resource.TestCheckResourceAttr(datasourceNamePbr, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNamePbr, "metadata.#"),
					resource.TestCheckResourceAttr(datasourceNamePbr, "policies.0.is_bidirectional", "false"),
					resource.TestCheckResourceAttr(datasourceNamePbr, "policies.0.policy_match.0.protocol_type", "UDP"),
					resource.TestCheckResourceAttr(datasourceNamePbr, "priority", "14"),
				),
			},
		},
	})
}

func testAccPbrDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name = "terraform-test-subnet-vpc_%[1]s"
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
		depends_on = [data.nutanix_clusters_v2.clusters]
	}
	resource "nutanix_vpc_v2" "test" {
		name =  "pbr_vpc_%[1]s"
		description = "pbr_vpc_ %[2]s"
		external_subnets{
		  subnet_reference = nutanix_subnet_v2.test.id
		}
		depends_on = [nutanix_subnet_v2.test]
	}

	resource "nutanix_pbr_v2" "rtest" {
		name = "%[1]s"
		description = "%[2]s"
		vpc_ext_id = nutanix_vpc_v2.test.ext_id
		priority = 14
		policies{
			policy_match{
				source{
					address_type = "ANY"
				}
				destination{
					address_type = "ANY"
				}
				protocol_type = "UDP"
			}
			policy_action{
				action_type  = "PERMIT"
			}
		}
	}

	data "nutanix_pbr_v2" "test" {
		ext_id = resource.nutanix_pbr_v2.rtest.ext_id
		depends_on = [
			resource.nutanix_pbr_v2.rtest
		]
	}
	`, name, desc)
}
