package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVpc = "nutanix_vpc_v2.test"

func TestAccV2NutanixVpcResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	vlanID := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf-test-vpc-%d", r)
	desc := "test vpc description"
	updatedName := fmt.Sprintf("updated-vpc-%d", r)
	updatedDesc := "updated vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfig(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.0.category_ids.#"),
					resource.TestCheckResourceAttrPair(resourceNameVpc, "metadata.0.category_ids.0", "nutanix_category_v2.test", "id"),
				),
			},
			{
				Config: testVpcConfig(updatedName, updatedDesc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.0.category_ids.#"),
					resource.TestCheckResourceAttrPair(resourceNameVpc, "metadata.0.category_ids.0", "nutanix_category_v2.test", "id"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcResource_WithExternallyRoutablePrefixes(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vpc-%d", r)
	vlanID := acctest.RandIntRange(1, 999)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithExtRoutablePrefix(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.0.category_ids.#"),
					resource.TestCheckResourceAttrPair(resourceNameVpc, "metadata.0.category_ids.0", "nutanix_category_v2.test", "id"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcResource_WithDHCP(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vpc-%d", r)
	vlanID := acctest.RandIntRange(1, 999)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithDHCP(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "common_dhcp_options.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcResource_WithTransitType(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vpc-%d", r)
	vlanID := acctest.RandIntRange(1, 999)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithTransitType(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "TRANSIT"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "common_dhcp_options.#"),
				),
			},
		},
	})
}

func testVpcConfig(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name = "terraform-test-subnet-vpc"
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
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_category_v2" "test" {
		key = "tf-test-category-key-%[3]d"
		value = "tf-test-category-value-%[3]d"
		description = "test category for vpc"
	}

	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets{
		  subnet_reference = nutanix_subnet_v2.test.id
		}
		metadata {
			category_ids = [nutanix_category_v2.test.id]
		}
		vpc_type = "REGULAR"
		depends_on = [nutanix_subnet_v2.test]
	}
`, name, desc, vlanID)
}

func testVpcConfigWithExtRoutablePrefix(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name = "terraform-test-subnet-vpc"
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
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_category_v2" "test" {
		key = "tf-test-category-key-%[1]s"
		value = "tf-test-category-value-%[1]s"
		description = "test category for vpc"
	}

	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
			external_ips {
			  ipv4 {
				value         = "192.168.0.24"
				prefix_length = 32
			  }
			}
			external_ips {
			  ipv4 {
				value         = "192.168.0.25"
				prefix_length = 32
			  }
			}
	   	}
		externally_routable_prefixes{
		  ipv4{
			ip{
			  value = "172.30.0.0"
			  prefix_length = 32
			}
			prefix_length = 16
		  }
		}
		metadata {
			category_ids = [nutanix_category_v2.test.id]
		}
		vpc_type = "REGULAR"
		depends_on = [nutanix_subnet_v2.test]
	}
`, name, desc, vlanID)
}

func testVpcConfigWithDHCP(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
	 	name              = "terraform-test-subnet-vpc"
	  	description       = "test subnet description"
		  cluster_reference = local.cluster0
		  subnet_type       = "VLAN"
		  network_id        = %[3]d
		  is_external       = true
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
			  pool_list {
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
		  name        = "%[1]s"
		  description = "%[2]s"
		  external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
			external_ips {
			  ipv4 {
				value         = "192.168.0.24"
				prefix_length = 32
			  }
			}
			external_ips {
			  ipv4 {
				value         = "192.168.0.25"
				prefix_length = 32
			  }
			}
		  }

		  externally_routable_prefixes {
			ipv4 {
			  ip {
				value         = "172.30.0.0"
				prefix_length = 32
			  }
			  prefix_length = 16
			}
		  }
		  depends_on = [nutanix_subnet_v2.test]
		}

`, name, desc, vlanID)
}

func testVpcConfigWithTransitType(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name = "terraform-test-subnet-vpc"
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
		depends_on = [data.nutanix_clusters_v2.clusters]
	}
	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets{
		  subnet_reference = nutanix_subnet_v2.test.id
		}
		vpc_type = "TRANSIT"
		depends_on = [nutanix_subnet_v2.test]
	}
`, name, desc, vlanID)
}
