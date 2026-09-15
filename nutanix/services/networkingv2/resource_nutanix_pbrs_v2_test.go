package networkingv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamePbr = "nutanix_pbr_v2.test"

func TestAccV2NutanixPbrResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-pbr-%d", r)
	desc := "test pbr description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPbrConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", name),
					resource.TestCheckResourceAttr(resourceNamePbr, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNamePbr, "metadata.#"),
					resource.TestCheckResourceAttr(resourceNamePbr, "policies.0.is_bidirectional", "false"),
					resource.TestCheckResourceAttr(resourceNamePbr, "policies.0.policy_match.0.protocol_type", "UDP"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", "14"),
				),
			},
		},
	})
}

func TestAccV2NutanixPbrResource_WithSourceDest(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-pbr-%d", r)
	desc := "test pbr description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPbrConfigWithSrcDstn(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", name),
					resource.TestCheckResourceAttr(resourceNamePbr, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNamePbr, "metadata.#"),
					resource.TestCheckResourceAttr(resourceNamePbr, "policies.0.is_bidirectional", "false"),
					resource.TestCheckResourceAttr(resourceNamePbr, "policies.0.policy_match.0.protocol_type", "ANY"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", "11"),
				),
			},
		},
	})
}

func TestAccV2NutanixPbrResource_ErrorWithPriority(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-pbr-%d", r)
	desc := "test pbr description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testPbrConfigWithDefaultPriority(name, desc),
				ExpectError: regexp.MustCompile("Modification of default routing policy with priority less than 10 is not allowed"),
			},
		},
	})
}

func testPbrConfig(name, desc string) string {
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

	resource "nutanix_pbr_v2" "test" {
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
`, name, desc)
}

func testPbrConfigWithSrcDstn(name, desc string) string {
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

	resource "nutanix_pbr_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		vpc_ext_id = nutanix_vpc_v2.test.ext_id
		priority = 11
		policies{
			policy_match{
				source{
					address_type = "EXTERNAL"
				}
				destination{
					address_type = "SUBNET"
					subnet_prefix{
						# In an IPv4Subnet the backend canonicalizes the address
						# (ip.prefix_length) to /32; the actual subnet mask goes in
						# the sibling prefix_length. Setting the mask on ip.prefix_length
						# causes a perpetual 32->24 diff on read.
						# See https://jira.nutanix.com/browse/ENG-877509
						ipv4{
							ip{
								value= "10.10.10.0"
								prefix_length = 32
							}
							prefix_length = 24
						}
					}
				}
				protocol_type = "ANY"
			}
			policy_action{
				action_type  = "FORWARD"
				nexthop_ip_address{
					ipv4{
						value = "10.10.10.10"
					}
				}
			}
		}
	}
`, name, desc)
}

func TestAccV2NutanixPbrResource_ProjectAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-pbr-projassoc-%d", r)
	desc := "pbr project association test"
	projectName := fmt.Sprintf("tf-pbr-pa-proj-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPbrProjectAssociationConfig(name, desc, r, projectName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("nutanix_pbr_v2.test", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_pbr_v2.test", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttr("data.nutanix_pbrs_v2.test", "routing_policies.#", "1"),
					resource.TestCheckResourceAttrPair("data.nutanix_pbrs_v2.test", "routing_policies.0.ext_id", "nutanix_pbr_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_pbrs_v2.test", "routing_policies.0.project_ext_id", "nutanix_project_v2.test", "ext_id"),
				),
			},
			{
				Config:      testPbrProjectAssociationConfig(name, desc, r, projectName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile("Update of project_ext_id is not supported"),
			},
		},
	})
}

func pbrProjectExtIDLine(override string) string {
	if override == "" {
		return `project_ext_id = nutanix_project_v2.test.ext_id`
	}
	return fmt.Sprintf(`project_ext_id = "%s"`, override)
}

func testPbrProjectAssociationConfig(name, desc string, r int, projectName, projectExtIDOverride string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		config  = (jsondecode(file("%[5]s")))
		subnets = local.config.networking.subnets
		cluster1 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_project_v2" "test" {
		name        = "%[4]s"
		project_id  = "%[4]s"
		description = "project association test"
	}

	resource "nutanix_subnet_v2" "test" {
		name              = "tf-pbr-pa-subnet-%[3]d"
		description       = "subnet for pbr project association"
		cluster_reference = local.cluster1
		subnet_type       = "VLAN"
		network_id        = local.subnets.vlan_id
		is_external       = true
		shared_with_projects = [nutanix_project_v2.test.ext_id]
		ip_config {
			ipv4 {
				ip_subnet {
					ip { value = local.subnets.network_ip }
					prefix_length = local.subnets.network_prefix
				}
				default_gateway_ip { value = local.subnets.gateway_ip }
				pool_list {
					start_ip { value = local.subnets.dhcp.start_ip }
					end_ip   { value = local.subnets.dhcp.end_ip }
				}
			}
		}
	}

	resource "nutanix_vpc_v2" "test" {
		name        = "tf-pbr-pa-vpc-%[3]d"
		description = "vpc for pbr project association"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		}
		project_ext_id       = nutanix_project_v2.test.ext_id
		shared_with_projects = [nutanix_project_v2.test.ext_id]
		depends_on           = [nutanix_subnet_v2.test]
		lifecycle {
			ignore_changes = [links, snat_ips, shared_with_projects]
		}
	}

	resource "nutanix_pbr_v2" "test" {
		name         = "%[1]s"
		description  = "%[2]s"
		priority     = 111
		vpc_ext_id   = nutanix_vpc_v2.test.id
		policies {
			policy_match {
				source {
					address_type = "ANY"
				}
				destination {
					address_type = "ANY"
				}
				protocol_type = "ANY"
			}
			policy_action {
				action_type = "PERMIT"
			}
		}
		%[6]s
		depends_on = [nutanix_project_v2.test, nutanix_vpc_v2.test]
	}

	data "nutanix_pbr_v2" "test" {
		ext_id     = nutanix_pbr_v2.test.id
		depends_on = [nutanix_pbr_v2.test]
	}

	data "nutanix_pbrs_v2" "test" {
		filter     = "name eq '%[1]s'"
		depends_on = [nutanix_pbr_v2.test]
	}
	`, name, desc, r, projectName, filepath, pbrProjectExtIDLine(projectExtIDOverride))
}

func testPbrConfigWithDefaultPriority(name, desc string) string {
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

	resource "nutanix_pbr_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		vpc_ext_id = nutanix_vpc_v2.test.ext_id
		priority = 1
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
`, name, desc)
}
