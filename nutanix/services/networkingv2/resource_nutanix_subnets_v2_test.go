package networkingv2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameSubnet = "nutanix_subnet_v2.test"

func TestAccV2NutanixSubnetResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-subnet-%d", r)
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

func TestAccV2NutanixSubnetResource_WithIPPool(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-subnet-%d", r)
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

func TestAccV2NutanixSubnetResource_WithExternalSubnet(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-subnet-%d", r)
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
					resource.TestCheckResourceAttr(resourceNameSubnet, "network_id", "122"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "ip_usage.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "cluster_reference"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "is_external", "true"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "is_nat_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.default_gateway_ip.0.value", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.ip_subnet.0.ip.0.value", "192.168.0.0"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.ip_subnet.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.pool_list.0.start_ip.0.value", "192.168.0.20"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.pool_list.0.end_ip.0.value", "192.168.0.30"),
				),
			},
		},
	})
}

func TestAccV2NutanixSubnetResource_isNatEnableFalse(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-subnet-%d", r)
	desc := "test subnet description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSubnetV2ConfigIsNatEnableFalse(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", name),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", desc),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "OVERLAY"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "is_external", "true"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "is_nat_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.default_gateway_ip.0.value", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.ip_subnet.0.ip.0.value", "192.168.0.0"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.ip_subnet.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.pool_list.0.start_ip.0.value", "192.168.0.20"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "ip_config.0.ipv4.0.pool_list.0.end_ip.0.value", "192.168.0.30"),
				),
			},
		},
	})
}

func TestAccV2NutanixSubnetResource_WithMetadata(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-subnet-%d", r)
	desc := "test subnet description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSubnetV2ConfigwithMetadata(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", name),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", desc),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "network_id", "112"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "ip_usage.#"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "cluster_reference"),
					resource.TestCheckResourceAttrSet(resourceNameSubnet, "metadata.#"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "metadata.0.category_ids.#", "1"),
					testCheckMetadataCategoryIDsContain(resourceNameSubnet, "nutanix_category_v2.test"),
					// data source check
					resource.TestCheckResourceAttr(datasourceNameSubnet, "metadata.0.category_ids.#", "1"),
					testCheckMetadataCategoryIDsContain(datasourceNameSubnet, "nutanix_category_v2.test"),
				),
			},
			{
				Config: testSubnetV2ConfigwithMetadataUpdate(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", name),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", desc),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "network_id", "112"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "metadata.0.category_ids.#", "2"),
					testCheckMetadataCategoryIDsContain(resourceNameSubnet, "nutanix_category_v2.test", "nutanix_category_v2.test2"),
					// data source check
					resource.TestCheckResourceAttr(datasourceNameSubnet, "metadata.0.category_ids.#", "2"),
					testCheckMetadataCategoryIDsContain(datasourceNameSubnet, "nutanix_category_v2.test", "nutanix_category_v2.test2"),
				),
			},
		},
	})
}

func testCheckMetadataCategoryIDsContain(target string, expectedCategoryResources ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		targetRs, ok := s.RootModule().Resources[target]
		if !ok {
			return fmt.Errorf("resource %q not found in state", target)
		}

		var categoryPrefix string
		for attrKey := range targetRs.Primary.Attributes {
			if strings.HasSuffix(attrKey, ".category_ids.#") {
				categoryPrefix = strings.TrimSuffix(attrKey, "#")
				break
			}
		}
		if categoryPrefix == "" {
			return fmt.Errorf("resource %q has no metadata category_ids", target)
		}

		actualIDs := make(map[string]struct{})
		for attrKey, attrVal := range targetRs.Primary.Attributes {
			if strings.HasPrefix(attrKey, categoryPrefix) && !strings.HasSuffix(attrKey, ".#") && attrVal != "" {
				actualIDs[attrVal] = struct{}{}
			}
		}

		for _, expectedResource := range expectedCategoryResources {
			expectedRs, exists := s.RootModule().Resources[expectedResource]
			if !exists {
				return fmt.Errorf("expected category resource %q not found in state", expectedResource)
			}

			expectedID := expectedRs.Primary.ID
			if expectedID == "" {
				expectedID = expectedRs.Primary.Attributes["id"]
			}
			if expectedID == "" {
				return fmt.Errorf("expected category resource %q has empty id", expectedResource)
			}

			if _, present := actualIDs[expectedID]; !present {
				return fmt.Errorf("resource %q metadata category_ids does not contain expected id %q from %q", target, expectedID, expectedResource)
			}
		}

		return nil
	}
}

func testSubnetV2Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
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
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
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
			depends_on = [data.nutanix_clusters_v2.clusters]
		}
`, name, desc)
}

func testSubnetV2ConfigWithExternalSubnet(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
		}

		resource "nutanix_subnet_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			cluster_reference = local.cluster0
			subnet_type = "VLAN"
			network_id = 122
			is_external = true
			is_nat_enabled = true
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
`, name, desc)
}

func testSubnetV2ConfigIsNatEnableFalse(name, desc string) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "clusters" {}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

resource "nutanix_vpc_v2" "test" {
  name        = "test_vpc_%[1]s"
  description = "test vpc %[2]s"
  vpc_type   = "TRANSIT"
}

resource "nutanix_subnet_v2" "test" {
  name 				= "%[1]s"
  description		= "%[2]s"
  cluster_reference = local.clusterExtId
  vpc_reference     = nutanix_vpc_v2.test.id
  subnet_type       = "OVERLAY"
  is_nat_enabled    = false
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
}
`, name, desc)
}

func testSubnetV2ConfigwithMetadata(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
		}

		resource "nutanix_category_v2" "test" {
			key = "tf-test-category-key-%[1]s"
			value = "tf-test-category-value-%[1]s"
			description = "test category for subnet"
		}

		resource "nutanix_subnet_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			cluster_reference = local.cluster0
			subnet_type = "VLAN"
			network_id = 112
			metadata {
				category_ids = [nutanix_category_v2.test.id]
			}
		}
		data "nutanix_subnet_v2" "test" {
			ext_id = nutanix_subnet_v2.test.id
		}
`, name, desc)
}

func testSubnetV2ConfigwithMetadataUpdate(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
		}

		resource "nutanix_category_v2" "test" {
			key = "tf-test-category-key-%[1]s"
			value = "tf-test-category-value-%[1]s"
			description = "test category for subnet"
		}

		resource "nutanix_category_v2" "test2" {
			key = "tf-test-category-key-%[1]s-2"
			value = "tf-test-category-value-%[1]s-2"
			description = "test category for subnet 2"
		}

		resource "nutanix_subnet_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			cluster_reference = local.cluster0
			subnet_type = "VLAN"
			network_id = 112
			metadata {
				category_ids = [nutanix_category_v2.test.id, nutanix_category_v2.test2.id]
			}
		}
		data "nutanix_subnet_v2" "test" {
			ext_id = nutanix_subnet_v2.test.id
		}
`, name, desc)
}
