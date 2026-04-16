package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameFIP = "data.nutanix_floating_ip_v2.test"

func TestAccV2NutanixFloatingIPDataSource_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFipDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameFIP, "name", name),
					resource.TestCheckResourceAttr(datasourceNameFIP, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameFIP, "metadata.#"),
					resource.TestCheckResourceAttrSet(datasourceNameFIP, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameFIP, "association.#"),
					resource.TestCheckResourceAttrSet(datasourceNameFIP, "external_subnet_reference"),
				),
			},
		},
	})
}

func testAccFipDataSourceConfig(name, desc string) string {
	networkID := acctest.RandIntRange(0, 999)
	return fmt.Sprintf(`

		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
		}

		resource "nutanix_subnet_v2" "test" {
			name = "terraform-test-subnet-floating-ip"
			description = "test subnet description"
			cluster_reference = local.cluster0
			subnet_type = "VLAN"
			network_id = %[3]d
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
		}
		resource "nutanix_floating_ip_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = nutanix_subnet_v2.test.id
	    }

		data "nutanix_floating_ip_v2" "test" {
			ext_id = nutanix_floating_ip_v2.test.ext_id
		}
	`, name, desc, networkID)
}
