package networkingv2_test

import (
	// "encoding/json"

	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameNs = "nutanix_network_security_policy_v2.test"

func TestAccV2NutanixNetworkSecurityResource_Basic(t *testing.T) {
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

func TestAccV2NutanixNetworkSecurityResource_WithRules(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceNameNs, "rules.0.spec.0.application_rule_spec.0.src_subnet.0.value", "192.168.0.0"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.0.spec.0.application_rule_spec.0.src_subnet.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.1.spec.0.application_rule_spec.0.dest_subnet.0.value", "192.68.0.0"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.1.spec.0.application_rule_spec.0.dest_subnet.0.prefix_length", "20"),
				),
			},
		},
	})
}

func TestAccV2NutanixNetworkSecurityResource_WithMultiEnvIsolationRuleSpecRule(t *testing.T) {
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

func TestAccV2NutanixNetworkSecurityResource_InvalidExtIDReference(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nsp-%d", r)
	desc := "test nsp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testNetworkSecurityInvalidConfig(name, desc),
				ExpectError: regexp.MustCompile(`(?s)Failed validation.*firstIsolationGroup.*regex.*invalid-ext-id`),
			},
		},
	})
}

func TestAccV2NutanixNetworkSecurityResource_WithApplicationAndInfraGroupRules(t *testing.T) {
	r := acctest.RandIntRange(1, 1000)
	name := fmt.Sprintf("tf-test-nsp-%d", r)
	desc := "test nsp description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkSecurityConfigWithApplicationAndInfraGroupRules(name, desc),
				Check: resource.ComposeTestCheckFunc(

					// basic attrs
					resource.TestCheckResourceAttrSet(resourceNameNs, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNs, "name", name),
					resource.TestCheckResourceAttr(resourceNameNs, "description", desc),
					resource.TestCheckResourceAttr(resourceNameNs, "state", "ENFORCE"),
					resource.TestCheckResourceAttr(resourceNameNs, "type", "APPLICATION"),
					resource.TestCheckResourceAttr(resourceNameNs, "scope", "VPC_LIST"),
					resource.TestCheckResourceAttr(resourceNameNs, "is_hitlog_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameNs, "is_ipv6_traffic_allowed", "false"),
					resource.TestCheckResourceAttr(resourceNameNs, "vpc_reference.#", "1"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.#", "15"),

					resource.TestCheckResourceAttrSet(resourceNameNs, "rules.0.spec.0.application_rule_spec.0.secured_group_category_references.#"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.0.description", "outbound for RDP tier"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.0.spec.0.application_rule_spec.0.is_all_protocol_allowed", "true"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.0.type", "APPLICATION"),

					resource.TestCheckResourceAttrSet(resourceNameNs, "rules.7.spec.0.intra_entity_group_rule_spec.0.secured_group_category_references.#"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.7.description", "deny amongst TFAppTest tier"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.7.spec.0.intra_entity_group_rule_spec.0.secured_group_action", "DENY"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.7.type", "INTRA_GROUP"),

					resource.TestCheckResourceAttrSet(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.secured_group_category_references.#"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.description", "ALL inbound for TFAppTest tier"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.0.end_port", "22"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.0.start_port", "22"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.1.end_port", "443"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.1.start_port", "443"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.2.end_port", "2074"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.2.start_port", "2074"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.3.end_port", "3389"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.3.start_port", "3389"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.4.end_port", "5985"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.tcp_services.4.start_port", "5985"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.spec.0.application_rule_spec.0.src_allow_spec", "ALL"),
					resource.TestCheckResourceAttr(resourceNameNs, "rules.10.type", "APPLICATION"),
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
		  description = "test rule 1"
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
			  src_subnet {
				value         = "192.168.0.0"
				prefix_length = 24
			  }
			  is_all_protocol_allowed = true
			}
		  }
		}
		rules{
		  description = "test rule 2"
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
			  dest_subnet {
			  	value         = "192.68.0.0"
			  	prefix_length = 20
			  }
			  is_all_protocol_allowed = true
			}
		  }
		}
		rules{
		  description = "test rule 3"
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

func testNetworkSecurityInvalidConfig(name, desc string) string {
	return fmt.Sprintf(`
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
				"",
			  ]
			  second_isolation_group =  [
				"invalid-ext-id",
			  ]
			}
		  }
		}
		is_hitlog_enabled = true
	  }
`, name, desc)
}

func testNetworkSecurityConfigWithApplicationAndInfraGroupRules(name, desc string) string {
	return fmt.Sprintf(`

# Vpc
# list Clusters
data "nutanix_clusters_v2" "clusters" {
}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
    ][
    0
  ]
}

# Vpc
resource "nutanix_subnet_v2" "subnet" {
  name              = "tf-test-subnet-vpc"
  description       = "test subnet description"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = 765
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

resource "nutanix_vpc_v2" "vpc" {
  name        = "tf-test-vpc-nsp"
  description = "vpc description"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.subnet.id
  }
}


# categories
resource "nutanix_category_v2" "app_type" {
  key         = "NspTfTest"
  value       = "NSPMultiVmApp"
  description = "NSP"
}

resource "nutanix_category_v2" "app_tier_ssh" {
  key         = "NspTfTestTier"
  value       = "SSH"
  description = "NSP"
}
resource "nutanix_category_v2" "app_tier_rdp" {
  key         = "NspTfTestTier"
  value       = "RDP"
  description = "NSP"
}
resource "nutanix_category_v2" "app_tier_web" {
  key         = "NspTfTestTier"
  value       = "TfTestWeb"
  description = "NSP"
}
resource "nutanix_category_v2" "app_tier_app" {
  key         = "NspTfTestTier"
  value       = "TFAppTest"
  description = "NSP"
}
resource "nutanix_category_v2" "app_tier_db" {
  key         = "NspTfTestTier"
  value       = "TFDBTest"
  description = "NSP"
}
#endregion application category

#region backup categories
resource "nutanix_category_v2" "backup" {
  key         = "TF-Test-AZ-Backup-01"
  value       = "RPO24h"
  description = "NSP"
}
#endregion backup categories

#region DR categories
resource "nutanix_category_v2" "dr_gold" {
  key         = "TF-Test-AZ-DR-Gold"
  value       = "RPOZero"
  description = "NSP"
}
resource "nutanix_category_v2" "dr_silver" {
  key         = "TF-Test-AZ-DR-Silver"
  value       = "RPO15m"
  description = "NSP"
}
resource "nutanix_category_v2" "dr_bronze" {
  key         = "TF-Test-AZ-DR-Bronze"
  value       = "RPO1h"
  description = "NSP"
}


data "nutanix_categories_v2" "app_type" {
  limit      = 1
  filter     = "key eq 'NspTfTest' and value eq 'NSPMultiVmApp'"
  depends_on = [nutanix_category_v2.app_type]
}

data "nutanix_categories_v2" "app_tier_app" {
  limit      = 1
  filter     = "key eq 'NspTfTestTier' and value eq 'TFAppTest'"
  depends_on = [nutanix_category_v2.app_tier_app]
}

data "nutanix_categories_v2" "app_tier_db" {
  limit      = 1
  filter     = "key eq 'NspTfTestTier' and value eq 'TFDBTest'"
  depends_on = [nutanix_category_v2.app_tier_db]
}

data "nutanix_categories_v2" "app_tier_web" {
  limit      = 1
  filter     = "key eq 'NspTfTestTier' and value eq 'TfTestWeb'"
  depends_on = [nutanix_category_v2.app_tier_web]
}

data "nutanix_categories_v2" "app_tier_rdp" {
  limit      = 1
  filter     = "key eq 'NspTfTestTier' and value eq 'RDP'"
  depends_on = [nutanix_category_v2.app_tier_rdp]
}

data "nutanix_categories_v2" "app_tier_ssh" {
  limit  = 1
  filter = "key eq 'NspTfTestTier' and value eq 'SSH'"
  depends_on = [nutanix_category_v2.app_tier_ssh]
}

resource "nutanix_network_security_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"
  type        = "APPLICATION"
  state       = "ENFORCE"
  scope       = "VPC_LIST"

  vpc_reference = [
    nutanix_vpc_v2.vpc.id,
  ]

  lifecycle {
    create_before_destroy = true
    ignore_changes        = [rules]
  }

  #* outbound rules
  rules {
    description = "outbound for RDP tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.app_tier_rdp.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = true
        src_allow_spec          = "NONE"
      }
    }
  }

  rules {
    description = "outbound for SSH tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.app_tier_ssh.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = true
        src_allow_spec          = "NONE"
      }
    }
  }

  #* ALL inbound rules
  rules {
    description = "ALL inbound for RDP tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.app_tier_rdp.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = false
        src_allow_spec          = "ALL"
        tcp_services {
          start_port = 3389
          end_port   = 3389
        }
        tcp_services {
          start_port = 5985
          end_port   = 5985
        }
        icmp_services {
          type = 8
          code = 0
        }
      }
    }
  }

  rules {
    description = "ALL inbound for SSH tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.app_tier_ssh.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = false
        src_allow_spec          = "ALL"
        tcp_services {
          start_port = 22
          end_port   = 22
        }
        icmp_services {
          type = 8
          code = 0
        }
      }
    }
  }

  #* outbound rules
  rules {
    description = "outbound for TFAppTest tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_app.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = true
        src_allow_spec          = "NONE"
      }
    }
  }

  rules {
    description = "outbound for TFDBTest tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_db.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = true
        src_allow_spec          = "NONE"
      }
    }
  }

  rules {
    description = "outbound for TfTestWeb tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_web.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = true
        src_allow_spec          = "NONE"
      }
    }
  }

  #* preventing vms with same tier from talking with each other
  rules {
    description = "deny amongst TFAppTest tier"
    type        = "INTRA_GROUP"
    spec {
      intra_entity_group_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_app.categories[0].ext_id,
        ]
        secured_group_action = "DENY"
      }
    }
  }

  rules {
    description = "deny amongst TFDBTest tier"
    type        = "INTRA_GROUP"
    spec {
      intra_entity_group_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_db.categories[0].ext_id,
        ]
        secured_group_action = "DENY"
      }
    }
  }

  rules {
    description = "deny amongst TfTestWeb tier"
    type        = "INTRA_GROUP"
    spec {
      intra_entity_group_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_web.categories[0].ext_id,
        ]
        secured_group_action = "DENY"
      }
    }
  }

  #* ALL inbound rules
  rules {
    description = "ALL inbound for TFAppTest tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_app.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = false
        src_allow_spec          = "ALL"
        tcp_services {
          start_port = 22
          end_port   = 22
        }
        tcp_services {
          start_port = 443
          end_port   = 443
        }
        tcp_services {
          start_port = 2074
          end_port   = 2074
        }
        tcp_services {
          start_port = 3389
          end_port   = 3389
        }
        tcp_services {
          start_port = 5985
          end_port   = 5985
        }
        icmp_services {
          type = 8
          code = 0
        }
      }
    }
  }

  rules {
    description = "ALL inbound for TFDBTest tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_db.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = false
        src_allow_spec          = "ALL"
        tcp_services {
          start_port = 22
          end_port   = 22
        }
        tcp_services {
          start_port = 2074
          end_port   = 2074
        }
        tcp_services {
          start_port = 3389
          end_port   = 3389
        }
        tcp_services {
          start_port = 5985
          end_port   = 5985
        }
        icmp_services {
          type = 8
          code = 0
        }
      }
    }
  }

  rules {
    description = "ALL inbound for TfTestWeb tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_web.categories[0].ext_id,
        ]
        src_category_references = []
        is_all_protocol_allowed = false
        src_allow_spec          = "ALL"
        tcp_services {
          start_port = 22
          end_port   = 22
        }
        tcp_services {
          start_port = 4000
          end_port   = 4000
        }
        tcp_services {
          start_port = 2074
          end_port   = 2074
        }
        tcp_services {
          start_port = 3389
          end_port   = 3389
        }
        tcp_services {
          start_port = 5985
          end_port   = 5985
        }
        icmp_services {
          type = 8
          code = 0
        }
      }
    }
  }

  #* inbound between tiers
  rules {
    description = "TfTestWeb inbound to TFAppTest tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_app.categories[0].ext_id,
        ]
        src_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_web.categories[0].ext_id,
        ]
        is_all_protocol_allowed = false
        tcp_services {
          start_port = 3000
          end_port   = 3000
        }
      }
    }
  }

  rules {
    description = "TFAppTest inbound to TFDBTest tier"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_db.categories[0].ext_id,
        ]
        src_category_references = [
          try(data.nutanix_categories_v2.app_type.categories[0].ext_id, ""),
          data.nutanix_categories_v2.app_tier_app.categories[0].ext_id,
        ]
        is_all_protocol_allowed = false
        tcp_services {
          start_port = 5432
          end_port   = 5432
        }
      }
    }
  }
}


`, name, desc)
}
