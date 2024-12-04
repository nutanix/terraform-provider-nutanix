package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameNs = "nutanix_network_security_policy_v2.test"

func TestAccNutanixNetworkSecurityV2Resource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nsp-%d", r)
	desc := "test nsp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkSecurityConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNs, "name", name),
					resource.TestCheckResourceAttr(resourceNameNs, "description", desc),
					resource.TestCheckResourceAttr(resourceNameNs, "state", "SAVE"),
					resource.TestCheckResourceAttrSet(resourceNameNs, "links.#"),
					resource.TestCheckResourceAttr(resourceNameNs, "type", "ISOLATION"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.0.type", "TWO_ENV_ISOLATION"),
				),
			},
		},
	})
}

func TestAccNutanixNetworkSecurityV2Resource_WithRules(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nsp-%d", r)
	desc := "test nsp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkSecurityConfigWithRules(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNs, "name", name),
					resource.TestCheckResourceAttr(resourceNameNs, "description", desc),
					resource.TestCheckResourceAttr(resourceNameNs, "state", "SAVE"),
					resource.TestCheckResourceAttrSet(resourceNameNs, "links.#"),
					resource.TestCheckResourceAttr(resourceNameNs, "type", "APPLICATION"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.#", "3"),
					resource.TestCheckResourceAttrSet(resourceNameNs, "vpc_reference.#"),
					resource.TestCheckResourceAttr(resourceNameNs, "is_hitlog_enabled", "false"),
				),
			},
		},
	})
}

func TestAccNutanixNetworkSecurityV2Resource_WithMultiEnvIsolationRuleSpecRule(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nsp-%d", r)
	desc := "test nsp description"

	isolationGroup1Key := "rules.0.spec.0.multi_env_isolation_rule_spec.0.spec.0.all_to_all_isolation_group.0.isolation_group.0.group_category_references.#"
	isolationGroup2Key := "rules.0.spec.0.multi_env_isolation_rule_spec.0.spec.0.all_to_all_isolation_group.0.isolation_group.1.group_category_references.#"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkSecurityConfigWithMultiEnvIsolationRuleSpecRule(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNs, "name", name),
					resource.TestCheckResourceAttr(resourceNameNs, "description", desc),
					resource.TestCheckResourceAttr(resourceNameNs, "state", "SAVE"),
					resource.TestCheckResourceAttrSet(resourceNameNs, "links.#"),
					resource.TestCheckResourceAttr(resourceNameNs, "type", "ISOLATION"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.#", "1"),
					resource.TestCheckResourceAttr(resourceNameNs, isolationGroup1Key, "2"),
					resource.TestCheckResourceAttr(resourceNameNs, isolationGroup2Key, "2"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.0.type", "MULTI_ENV_ISOLATION"),
					resource.TestCheckResourceAttrSet(resourceNameNs, "vpc_reference.#"),
					resource.TestCheckResourceAttr(resourceNameNs, "is_hitlog_enabled", "false"),
				),
			},
		},
	})
}

func testNetworkSecurityConfig(name, desc string) string {
	return fmt.Sprintf(`

    data "nutanix_categories_v2" "test" {}

	resource "nutanix_network_security_policy_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		state = "SAVE"
		type = "ISOLATION"
		rules{
		  type = "TWO_ENV_ISOLATION"
		  spec{
			two_env_isolation_rule_spec{
			  first_isolation_group = [
				data.nutanix_categories_v2.test.categories.0.ext_id,
			  ]
			  second_isolation_group =  [
				data.nutanix_categories_v2.test.categories.1.ext_id,
			  ]
			}
		  }
		}
		is_hitlog_enabled = true
	  }
`, name, desc)
}

func testNetworkSecurityConfigWithRules(name, desc string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
	}
	
	resource "nutanix_subnet_v2" "test" {
		name = "tf-test-subnet-vpc-%[1]s"
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
		name =  "tf-test-vpc-%[1]s"
		description = "vpc %[2]s"
		external_subnets{
		  subnet_reference = nutanix_subnet_v2.test.id
		}
		depends_on = [nutanix_subnet_v2.test]
	}

	data "nutanix_categories_v2" "test" {}

	resource "nutanix_network_security_policy_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		type = "APPLICATION"
		state = "SAVE"
		rules{
		  description = "test"
		  type  = "APPLICATION"
		  spec{
			application_rule_spec{
			  secured_group_category_references = [
				data.nutanix_categories_v2.test.categories.0.ext_id,
				data.nutanix_categories_v2.test.categories.1.ext_id
			  ]
			  src_category_references = [
				data.nutanix_categories_v2.test.categories.2.ext_id
			  ]
			  is_all_protocol_allowed = true
			}
		  }
		}
		rules{
		  description = "test22"
		  type  = "APPLICATION"
		  spec{
			application_rule_spec{
			  secured_group_category_references = [
				data.nutanix_categories_v2.test.categories.3.ext_id,
				data.nutanix_categories_v2.test.categories.4.ext_id
			  ]
			  dest_category_references = [
				data.nutanix_categories_v2.test.categories.5.ext_id
			  ]
			  is_all_protocol_allowed = true
			}
		  }
		}
		rules{
		  type = "INTRA_GROUP"
		  spec{
			intra_entity_group_rule_spec{
			  secured_group_category_references = [
				data.nutanix_categories_v2.test.categories.6.ext_id,
				data.nutanix_categories_v2.test.categories.7.ext_id
			  ]
			  secured_group_action = "ALLOW"
			}
		  }
		}
		  
		vpc_reference = [
		  nutanix_vpc_v2.test.id
		]
		is_hitlog_enabled = false
		depends_on = [nutanix_vpc_v2.test]
	  }
	  `, name, desc)
}

func testNetworkSecurityConfigWithMultiEnvIsolationRuleSpecRule(name, desc string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
	}
	
	resource "nutanix_subnet_v2" "test" {
		name = "tf-test-subnet-vpc-%[1]s"
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
		name =  "tf-test-vpc-%[1]s"
		description = "vpc %[2]s"
		external_subnets{
		  subnet_reference = nutanix_subnet_v2.test.id
		}
		depends_on = [nutanix_subnet_v2.test]
	}

	data "nutanix_categories_v2" "test" {}

	resource "nutanix_network_security_policy_v2" "test" {
		name = "%[1]s"
		description = "%[2]s"
		type = "ISOLATION"
		state = "SAVE"
		rules{
		  description = "test"
		  type  = "MULTI_ENV_ISOLATION"
		  spec{
			multi_env_isolation_rule_spec{
			    spec{
					all_to_all_isolation_group{
						isolation_group{
							group_category_references = [
								data.nutanix_categories_v2.test.categories.0.ext_id,
								data.nutanix_categories_v2.test.categories.1.ext_id
							]
						}
						isolation_group{
							group_category_references = [	
								data.nutanix_categories_v2.test.categories.2.ext_id,
								data.nutanix_categories_v2.test.categories.3.ext_id	
							]
						}
					}
				}
			}
		  }
		}
		
		  
		vpc_reference = [
		  nutanix_vpc_v2.test.id
		]
		is_hitlog_enabled = false
		depends_on = [nutanix_vpc_v2.test]
	  }
	  `, name, desc)
}
