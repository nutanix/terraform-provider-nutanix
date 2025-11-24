package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameFIPps = "data.nutanix_floating_ips_v2.test"

func TestAccV2NutanixFloatingIPsDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFipsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameFIPps, "floating_ips.#"),
					checkAttributeLength(datasourceNameFIPps, "floating_ips", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixFloatingIPsDataSource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFipsDataSourceWithFilterConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameFIPps, "floating_ips.#"),
					resource.TestCheckResourceAttr(datasourceNameFIPps, "floating_ips.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameFIPps, "floating_ips.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameFIPps, "floating_ips.0.description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameFIPps, "floating_ips.0.metadata.#"),
					resource.TestCheckResourceAttrSet(datasourceNameFIPps, "floating_ips.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameFIPps, "floating_ips.0.association.#"),
					resource.TestCheckResourceAttrSet(datasourceNameFIPps, "floating_ips.0.external_subnet_reference"),
				),
			},
		},
	})
}

func TestAccV2NutanixFloatingIPsDataSource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFipsDataSourceWithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameFIPps, "floating_ips.#"),
					resource.TestCheckResourceAttr(datasourceNameFIPps, "floating_ips.#", "0"),
				),
			},
		},
	})
}

func testAccFipsDataSourceConfig(name, desc string) string {
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
		}
		resource "nutanix_floating_ip_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = nutanix_subnet_v2.test.id
	    }

		data "nutanix_floating_ips_v2" "test" {
			depends_on = [
				resource.nutanix_floating_ip_v2.test
			]
		}
	`, name, desc)
}

func testAccFipsDataSourceWithFilterConfig(name, desc string) string {
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
		}
		resource "nutanix_floating_ip_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = nutanix_subnet_v2.test.id
	    }

		data "nutanix_floating_ips_v2" "test" {
			filter = "name eq '%[1]s'"
			depends_on = [
				resource.nutanix_floating_ip_v2.test
			]
		}
	`, name, desc)
}

func testAccFipsDataSourceWithInvalidFilterConfig() string {
	return `
		data "nutanix_floating_ips_v2" "test" {
			filter = "name eq 'invalid_name'"
		}
	`
}
