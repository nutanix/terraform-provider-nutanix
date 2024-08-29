package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameSubnet = "nutanix_subnet_v2.test"

func TestAccNutanixSubnetV2_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-subnet-%d", r)
	desc := "test subnet description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSubnetV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", name),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", desc),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "network_id", "112"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "ip_usage.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "cluster_reference"),
				),
			},
			{
				Config: testSubnetV2Config("updated-name", "updated-description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", "updated-name"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "updated-description"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "network_id", "112"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "ip_usage.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "cluster_reference"),
				),
			},
		},
	})
}

func TestAccNutanixSubnetV2_WithIPPool(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-subnet-%d", r)
	desc := "test subnet description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSubnetV2ConfigWithIPPool(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", name),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", desc),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "network_id", "112"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "ip_usage.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "cluster_reference"),
				),
			},
		},
	})
}

func TestAccNutanixSubnetV2_WithExternalSubnet(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-subnet-%d", r)
	desc := "test subnet description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSubnetV2ConfigWithExternalSubnet(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", name),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", desc),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "network_id", "112"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "ip_usage.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "cluster_reference"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "is_external", "true"),
				),
			},
		},
	})
}

func testSubnetV2Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
		
		resource "nutanix_subnet_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			cluster_reference = local.cluster0
			subnet_type = "VLAN"
			network_id = 112
		}
`, name, desc)
}

func testSubnetV2ConfigWithIPPool(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
		
		resource "nutanix_subnet_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			cluster_reference = local.cluster0
			subnet_type = "VLAN"
			network_id = 112
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

			dhcp_options {
				domain_name_servers {
					ipv4{
						value = "8.8.8.8"
					}
				}
				search_domains = ["eng.nutanix.com"]
				domain_name      = "nutanix.com"
				tftp_server_name = "10.5.0.10"
				boot_file_name = "pxelinux.0"
			}
			depeds_on = [data.nutanix_clusters.clusters]
		}
`, name, desc)
}

func testSubnetV2ConfigWithExternalSubnet(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
		
		resource "nutanix_subnet_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
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
		depeds_on = [data.nutanix_clusters.clusters]
		}
`, name, desc)
}
