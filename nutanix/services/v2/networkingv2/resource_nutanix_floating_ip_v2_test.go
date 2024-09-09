package networkingv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamefip = "nutanix_floating_ip_v2.test"

func TestAccNutanixFloatingIPv2_Basic(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceNamefip, "name", name),
					resource.TestCheckResourceAttr(resourceNamefip, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "external_subnet_reference"),
				),
			},
			{
				Config: testFloatingIPv2Config(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamefip, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNamefip, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "external_subnet_reference"),
				),
			},
		},
	})
}

func TestAccNutanixFloatingIPv2_WithVmNICAssociation(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFloatingIPv2ConfigwithVMNic(filepath, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamefip, "name", name),
					resource.TestCheckResourceAttr(resourceNamefip, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "association.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "external_subnet_reference"),
				),
			},
		},
	})
}

func TestAccNutanixFloatingIPv2_WithPrivateipAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFloatingIPv2ConfigwithPrivateIP(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamefip, "name", name),
					resource.TestCheckResourceAttr(resourceNamefip, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "association.#"),
					resource.TestCheckResourceAttrSet(resourceNamefip, "external_subnet_reference"),
				),
			},
		},
	})
}

func testFloatingIPv2Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
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

func testFloatingIPv2ConfigwithVMNic(filepath, name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
			config = (jsondecode(file("%s")))
			floating_ip = local.config.networking.floating_ip
		}
		
		resource "nutanix_subnet_v2" "test" {
			name = "terraform-test-subnet-floating-ip-1"
			description = "test subnet floating ip description"
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
			association{
				vm_nic_association{
					vm_nic_reference = local.floating_ip.vm_nic_reference
				}
			  }
		  }
`, filepath, name, desc)
}

func testFloatingIPv2ConfigwithPrivateIP(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
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
			depends_on = [data.nutanix_clusters.clusters]
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
