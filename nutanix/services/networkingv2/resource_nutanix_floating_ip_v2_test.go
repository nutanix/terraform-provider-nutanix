package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameFIP = "nutanix_floating_ip_v2.test"

func TestAccV2NutanixFloatingIPResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	updatedName := fmt.Sprintf("updated-fip-%d", r)
	updatedDesc := "updated fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFloatingIPv2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameFIP, "name", name),
					resource.TestCheckResourceAttr(resourceNameFIP, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "external_subnet_reference"),
				),
			},
			{
				Config: testFloatingIPv2Config(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameFIP, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameFIP, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "external_subnet_reference"),
				),
			},
		},
	})
}

func TestAccV2NutanixFloatingIPResource_WithVmNICAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFloatingIPv2ConfigWithVMNic(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameFIP, "name", name),
					resource.TestCheckResourceAttr(resourceNameFIP, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "association.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "external_subnet_reference"),
				),
			},
		},
	})
}

func TestAccV2NutanixFloatingIPResource_WithPrivateIpAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFloatingIPv2ConfigWithPrivateIP(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameFIP, "name", name),
					resource.TestCheckResourceAttr(resourceNameFIP, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "association.#"),
					resource.TestCheckResourceAttrSet(resourceNameFIP, "external_subnet_reference"),
				),
			},
		},
	})
}

func testFloatingIPv2Config(name, desc string) string {
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
`, name, desc)
}

func testFloatingIPv2ConfigWithVMNic(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			  ][0]
			config = jsondecode(file("%[3]s"))
			preEnv            = local.config.pre_env
		}

		data "nutanix_subnets_v2" "test" {
			filter = "name eq '${local.preEnv.external_nat_subnet.name}'"
		}

		data "nutanix_virtual_machines_v2" "test" {
			filter = "name eq '${local.preEnv.integration_vm.name}'"
		}

		resource "nutanix_floating_ip_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = data.nutanix_subnets_v2.test.subnets.0.ext_id
			floating_ip {
				ipv4 {
					value = "10.44.3.206"
				}
			}
			association{
				vm_nic_association{
					vm_nic_reference = data.nutanix_virtual_machines_v2.test.vms.0.nics.0.ext_id
				}
			  }
		  }
`, name, desc, filepath)
}

func testFloatingIPv2ConfigWithPrivateIP(name, desc string) string {
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
			depends_on = [data.nutanix_clusters_v2.clusters]
		}

		resource "nutanix_vpc_v2" "test" {
			name =  "terraform-test-vpc-floating-ip"
			description = "test vpc description"
			external_subnets{
			  subnet_reference = nutanix_subnet_v2.test.id
			}
			common_dhcp_options{
				domain_name_servers{
					ipv4{
						value = "8.8.8.9"
						prefix_length = 32
					}
				}
				domain_name_servers{
					ipv4{
						value = "8.8.8.8"
						prefix_length = 32
					}
				}
			}
			depends_on = [nutanix_subnet_v2.test]
		}

		resource "nutanix_floating_ip_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = nutanix_subnet_v2.test.id
			association{
				private_ip_association{
					vpc_reference = nutanix_vpc_v2.test.id
					private_ip{
						ipv4{
							value = "8.8.10.13"
						}
					}
				}
			  }
			depends_on = [nutanix_vpc_v2.test]
		  }
`, name, desc)
}
