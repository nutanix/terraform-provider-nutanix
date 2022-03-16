package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNutanixNetworkSecurityRule_basic(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	rInt := acctest.RandInt()
	resourceName := "nutanix_network_security_rule.TEST-TIER"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixNetworkSecurityRuleConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists(resourceName),
				),
			},
			{
				Config: testAccNutanixNetworkSecurityRuleConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccNutanixNetworkSecurityRule_isolation(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	rInt := acctest.RandInt()
	resourceName := "nutanix_network_security_rule.isolation"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixNetworkSecurityRuleIsolationConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixNetworkSecurityRule_adrule(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	rInt := acctest.RandInt()
	resourceName := "nutanix_network_security_rule.VDI"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixNetworkSecurityRuleConfigAdRule(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists(resourceName),
				),
			},
			{
				Config: testAccNutanixNetworkSecurityRuleConfigAdRuleUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixNetworkSecurityRuleWithServiceAndAddressGroupsInInbound(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	rInt := acctest.RandInt()
	sgName := fmt.Sprintf("tf-service-group-%d", rInt)
	agName := fmt.Sprintf("tf-addr-group-%d", rInt)
	securityPolicyName := fmt.Sprintf("tf-sec-policy-%d", rInt)

	resourceName := "nutanix_network_security_rule.VDI"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testNetworkSecurityRuleConfigWithServiceAndAddressGroupsInInbound(sgName, agName, securityPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists(resourceName),
				),
			},
			// {
			// 	Config: testNetworkSecurityRuleConfigWithServiceAndAddressGroupsInOutbound(sgName, agName, securityPolicyName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckNutanixNetworkSecurityRuleExists(resourceName),
			// 	),
			// },
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixNetworkSecurityRuleWithServiceAndAddressGroupsInOutbound(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	rInt := acctest.RandInt()
	sgName := fmt.Sprintf("tf-service-group-%d", rInt)
	agName := fmt.Sprintf("tf-addr-group-%d", rInt)
	securityPolicyName := fmt.Sprintf("tf-sec-policy-%d", rInt)

	resourceName := "nutanix_network_security_rule.VDI"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testNetworkSecurityRuleConfigWithServiceAndAddressGroupsInOutbound(sgName, agName, securityPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixNetworkSecurityRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixNetworkSecurityRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_network_security_rule" {
			continue
		}
		for {
			_, err := conn.API.V3.GetNetworkSecurityRule(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}
	}

	return nil
}

func testAccNutanixNetworkSecurityRuleIsolationConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_network_security_rule" "isolation" {
	name        = "test-acc-isolation-rule-%d"
	description = "Isolation Test Acc"

	isolation_rule_action = "APPLY"

	isolation_rule_first_entity_filter_kind_list = ["vm"]
	isolation_rule_first_entity_filter_type      = "CATEGORIES_MATCH_ALL"
	isolation_rule_first_entity_filter_params {
		name   = "Environment"
		values = ["Dev"]
	}

	isolation_rule_second_entity_filter_kind_list = ["vm"]
	isolation_rule_second_entity_filter_type      = "CATEGORIES_MATCH_ALL"
	isolation_rule_second_entity_filter_params {
		name   = "Environment"
		values = ["Production"]
	}
}
`, r)
}

func testAccNutanixNetworkSecurityRuleConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "test-category-key"{
    name = "TIER-1"
	  description = "TIER Category Key"
}

resource "nutanix_category_key" "USER"{
    name = "user"
	  description = "user Category Key"
}


resource "nutanix_category_value" "WEB"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "WEB Category Value"
	  value = "WEB-1"
}

resource "nutanix_category_value" "APP"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "APP Category Value"
	 value = "APP-1"
}

resource "nutanix_category_value" "DB"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "DB Category Value"
	 value = "DB-1"
}

resource "nutanix_category_value" "ashwini"{
    name = "${nutanix_category_key.USER.id}"
	  description = "ashwini Category Value"
	 value = "ashwini-1"
}


resource "nutanix_network_security_rule" "TEST-TIER" {
  name        = "RULE-1-TIERS-%d"
  description = "rule 1 tiers"

  app_rule_action = "APPLY"

  app_rule_inbound_allow_list {
    peer_specification_type = "FILTER"
    filter_type             = "CATEGORIES_MATCH_ALL"
    filter_kind_list        = ["vm"]

    filter_params {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.WEB.id}"]
	}

	icmp_type_code_list {
		code = 1
		type = 1
	}

	tcp_port_range_list {
		end_port = 22
		start_port = 80
	}

	udp_port_range_list {
		end_port = 82
		start_port = 8080
	}
  }

  app_rule_target_group_default_internal_policy = "DENY_ALL"

  app_rule_target_group_peer_specification_type = "FILTER"

  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"

  app_rule_target_group_filter_kind_list = ["vm"]

  app_rule_target_group_filter_params {
    name   = "${nutanix_category_key.test-category-key.id}"
    values = ["${nutanix_category_value.APP.id}"]
  }

  app_rule_target_group_filter_params {
    name   = "${nutanix_category_key.USER.id}"
    values = ["${nutanix_category_value.ashwini.id}"]
  }

  app_rule_outbound_allow_list {
    peer_specification_type = "FILTER"
    filter_type             = "CATEGORIES_MATCH_ALL"
    filter_kind_list        = ["vm"]

    filter_params {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.DB.id}"]
    }
  }
}
`, r)
}

func testAccNutanixNetworkSecurityRuleConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "test-category-key"{
    name = "TIER-1"
	  description = "TIER Category Key"
}

resource "nutanix_category_key" "USER"{
    name = "user"
	  description = "user Category Key"
}

resource "nutanix_category_value" "WEB"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "WEB Category Value"
	 value = "WEB-1"
}

resource "nutanix_category_value" "APP"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "APP Category Value"
	 value = "APP-1"
}

resource "nutanix_category_value" "DB"{
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "DB Category Value"
	 value = "DB-1"
}

resource "nutanix_category_value" "ashwini"{
    name = "${nutanix_category_key.USER.id}"
	  description = "ashwini Category Value"
	 value = "ashwini-1"
}


resource "nutanix_network_security_rule" "TEST-TIER" {
  name        = "RULE-1-TIERS-%d"
  description = "rule 1 tiers Updated"

  app_rule_action = "APPLY"

  app_rule_inbound_allow_list {
    peer_specification_type = "FILTER"
    filter_type             = "CATEGORIES_MATCH_ALL"
    filter_kind_list        = ["vm"]

    filter_params {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.WEB.id}"]
	}

	icmp_type_code_list {
		code = 1
		type = 1
	}

	tcp_port_range_list {
		end_port = 22
		start_port = 80
	}

	udp_port_range_list {
		end_port = 82
		start_port = 8080
	}
  }

  app_rule_target_group_default_internal_policy = "DENY_ALL"

  app_rule_target_group_peer_specification_type = "FILTER"

  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"

  app_rule_target_group_filter_kind_list = ["vm"]

  app_rule_target_group_filter_params {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.APP.id}"]
  }

  app_rule_target_group_filter_params {
      name   = "${nutanix_category_key.USER.id}"
      values = ["${nutanix_category_value.ashwini.id}"]
  }


  app_rule_outbound_allow_list {
    peer_specification_type = "FILTER"
    filter_type             = "CATEGORIES_MATCH_ALL"
    filter_kind_list        = ["vm"]

    filter_params {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.DB.id}"]
    }
  }
}
`, r)
}

func testAccNutanixNetworkSecurityRuleConfigAdRule(r int) string {
	return fmt.Sprintf(`
	resource "nutanix_category_value" "ad-group-user-1" {
		name = "ADGroup"
		description = "group user category value"
		value = "%s"
	}
	resource "nutanix_network_security_rule" "VDI" {
		name           = "tf-%d"
		ad_rule_action = "APPLY"
		description    = "test"
		#   app_rule_action = "APPLY"
		ad_rule_inbound_allow_list {
		  ip_subnet               = "10.0.0.0"
		  ip_subnet_prefix_length = "8"
		  peer_specification_type = "IP_SUBNET"
		  protocol                = "ALL"
		}
		ad_rule_target_group_default_internal_policy = "DENY_ALL"
		ad_rule_target_group_filter_kind_list = [
		  "vm"
		]
		ad_rule_target_group_filter_params {
		  name = "ADGroup"
		  values = [
			"%s"
		  ]
		}
		ad_rule_target_group_filter_type             = "CATEGORIES_MATCH_ALL"
		ad_rule_target_group_peer_specification_type = "FILTER"
		ad_rule_outbound_allow_list {
		  ip_subnet               = "10.0.0.0"
		  ip_subnet_prefix_length = "8"
		  peer_specification_type = "IP_SUBNET"
		  protocol                = "ALL"
		}
		depends_on = [nutanix_category_value.ad-group-user-1]
	  }
`, testVars.AdRuleTarget.Values, r, testVars.AdRuleTarget.Values)
}

func testAccNutanixNetworkSecurityRuleConfigAdRuleUpdate(r int) string {
	return fmt.Sprintf(`
	resource "nutanix_category_value" "ad-group-user-1" {
		name = "ADGroup"
		description = "group user category value"
		value = "%s"
	}
	resource "nutanix_network_security_rule" "VDI" {
		name           = "tf-%d"
		ad_rule_action = "APPLY"
		description    = "test update"
		#   app_rule_action = "APPLY"
		ad_rule_inbound_allow_list {
		  ip_subnet               = "10.0.0.0"
		  ip_subnet_prefix_length = "8"
		  peer_specification_type = "IP_SUBNET"
		  protocol                = "ALL"
		}
		ad_rule_target_group_default_internal_policy = "DENY_ALL"
		ad_rule_target_group_filter_kind_list = [
		  "vm"
		]
		ad_rule_target_group_filter_params {
		  name = "ADGroup"
		  values = [
			"%s"
		  ]
		}
		ad_rule_target_group_filter_type             = "CATEGORIES_MATCH_ALL"
		ad_rule_target_group_peer_specification_type = "FILTER"
		ad_rule_outbound_allow_list {
		  ip_subnet               = "10.0.0.0"
		  ip_subnet_prefix_length = "8"
		  peer_specification_type = "IP_SUBNET"
		  protocol                = "ALL"
		}
		depends_on = [nutanix_category_value.ad-group-user-1]
	  }
`, testVars.AdRuleTarget.Values, r, testVars.AdRuleTarget.Values)
}

func testNetworkSecurityRuleConfigWithServiceAndAddressGroupsInInbound(sgName, agName, securityPolicyName string) string {
	return fmt.Sprintf(`
	resource "nutanix_category_value" "ad-group-user-1" {
		name = "ADGroup"
		description = "group user category value"
		value = "%s"
	}
	resource "nutanix_network_security_rule" "VDI" {
		name           = "%s-in"
		ad_rule_action = "APPLY"
		description    = "test"
		#   app_rule_action = "APPLY"
		ad_rule_inbound_allow_list {
			peer_specification_type = "ALL"  
			service_group_list {
			  kind = "service_group"
			  uuid = nutanix_service_group.service1.id
			}
			address_group_inclusion_list {
			  kind = "address_group"
			  uuid = nutanix_address_group.address1.id
			}
		}
		ad_rule_target_group_default_internal_policy = "DENY_ALL"
		ad_rule_target_group_filter_kind_list = [
		  "vm"
		]
		ad_rule_target_group_filter_params {
		  name = "ADGroup"
		  values = [
			"%s"
		  ]
		}
		ad_rule_target_group_filter_type             = "CATEGORIES_MATCH_ALL"
		ad_rule_target_group_peer_specification_type = "FILTER"
		ad_rule_outbound_allow_list {
		  ip_subnet               = "10.0.0.0"
		  ip_subnet_prefix_length = "8"
		  peer_specification_type = "IP_SUBNET"
		  protocol                = "ALL"
		}
		depends_on = [nutanix_category_value.ad-group-user-1]
	  }
	  resource "nutanix_service_group" "service1" {
		name = "%s-in"
		description = "test"
	  
		service_list {
			protocol = "TCP"
			tcp_port_range_list {
			  start_port = 22
			  end_port = 22
			}
			tcp_port_range_list {
			  start_port = 2222
			  end_port = 2222
			}
		}
	  }
	  resource "nutanix_address_group" "address1" {
		name = "%s-in"
		description = "test"
	  
		ip_address_block_list {
		  ip = "10.0.0.0"
		  prefix_length = 24
		}
	  }
`, testVars.AdRuleTarget.Values, securityPolicyName, testVars.AdRuleTarget.Values, sgName, agName)
}

func testNetworkSecurityRuleConfigWithServiceAndAddressGroupsInOutbound(sgName, agName, securityPolicyName string) string {
	return fmt.Sprintf(`
	resource "nutanix_category_value" "ad-group-user-1" {
		name = "ADGroup"
		description = "group user category value"
		value = "%s"
	}
	resource "nutanix_network_security_rule" "VDI" {
		name           = "%s-out"
		ad_rule_action = "APPLY"
		description    = "test"
		#   app_rule_action = "APPLY"
		ad_rule_inbound_allow_list {
			ip_subnet               = "10.0.0.0"
			ip_subnet_prefix_length = "8"
			peer_specification_type = "IP_SUBNET"
			protocol                = "ALL"
		}
		ad_rule_target_group_default_internal_policy = "DENY_ALL"
		ad_rule_target_group_filter_kind_list = [
		  "vm"
		]
		ad_rule_target_group_filter_params {
		  name = "ADGroup"
		  values = [
			"%s"
		  ]
		}
		ad_rule_target_group_filter_type             = "CATEGORIES_MATCH_ALL"
		ad_rule_target_group_peer_specification_type = "FILTER"
		ad_rule_outbound_allow_list {
			peer_specification_type = "ALL"
			service_group_list {
				kind = "service_group"
				uuid = nutanix_service_group.service1.id
			}
		
			address_group_inclusion_list {
				kind = "address_group"
				uuid = nutanix_address_group.address1.id
			}
		}
		depends_on = [nutanix_category_value.ad-group-user-1]
	  }
	  resource "nutanix_service_group" "service1" {
		name = "%s-out"
		description = "test"
	  
		service_list {
			protocol = "TCP"
			tcp_port_range_list {
			  start_port = 22
			  end_port = 22
			}
			tcp_port_range_list {
			  start_port = 2222
			  end_port = 2222
			}
		}
	  }
	  resource "nutanix_address_group" "address1" {
		name = "%s-out"
		description = "test"
	  
		ip_address_block_list {
		  ip = "10.0.0.0"
		  prefix_length = 24
		}
	  }
`, testVars.AdRuleTarget.Values, securityPolicyName, testVars.AdRuleTarget.Values, sgName, agName)
}
