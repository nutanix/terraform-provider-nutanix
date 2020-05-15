package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixNetworkSecurityRuleDataSource_basic(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	r := acctest.RandIntRange(0, 500)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityRuleDataSourceConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_network_security_rule.test", "name", fmt.Sprintf("RULE-1-TIERS-%d", r)),
					resource.TestCheckResourceAttr(
						"data.nutanix_network_security_rule.test", "app_rule_action", "APPLY"),
				),
			},
		},
	})
}

func TestAccNutanixNetworkSecurityRuleDataSource_isolation(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	r := acctest.RandIntRange(0, 500)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityRuleDataSourceConfigIsolation(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_network_security_rule.test", "name", fmt.Sprintf("test-acc-isolation-rule-%d", r)),
				),
			},
		},
	})
}

func testAccNetworkSecurityRuleDataSourceConfigIsolation(r int) string {
	return fmt.Sprintf(`
      %s

      data "nutanix_network_security_rule" "test" {
        network_security_rule_id = nutanix_network_security_rule.isolation.id
      }
  `, testAccNutanixNetworkSecurityRuleIsolationConfig(r))
}

func testAccNetworkSecurityRuleDataSourceConfig(r int) string {
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

resource "nutanix_category_value" "group"{
    name = "${nutanix_category_key.USER.id}"
	  description = "group Category Value"
	 value = "group-1"
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
    values = ["${nutanix_category_value.group.id}"]
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

data "nutanix_network_security_rule" "test" {
	network_security_rule_id = "${nutanix_network_security_rule.TEST-TIER.id}"
}
`, r)
}
