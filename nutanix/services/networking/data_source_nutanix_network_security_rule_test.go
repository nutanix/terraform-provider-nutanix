package networking_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixNetworkSecurityRuleDataSource_basic(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	r := acctest.RandIntRange(0, 500)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
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
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
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

func TestAccNutanixNetworkSecurityRuleDataSource_advanced(t *testing.T) {
	// Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	r := acctest.RandIntRange(0, 500)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityRuleDataSourceAdvancedConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_network_security_rule.test", "name", fmt.Sprintf("RULE-1-TIERS-%d", r)),
					resource.TestCheckResourceAttr(
						"data.nutanix_network_security_rule.test", "app_rule_action", "MONITOR"),
				),
			},
		},
	})
}

func isGCPEnvironment() bool {
	return os.Getenv("NUTANIX_GCP") == "true"
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

func testAccNetworkSecurityRuleDataSourceAdvancedConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_category_key" "AppType" {
  name = "AppType"
}

resource "nutanix_category_value" "DB" {
  name        = data.nutanix_category_key.AppType.id
  description = "Test Category Value"
  value       = "DB-1"
}

resource "nutanix_category_key" "test-category-key" {
  name        = "TIER-1"
  description = "TIER Category Key"
}

resource "nutanix_category_value" "APP-1" {
  name        = "${nutanix_category_key.test-category-key.id}"
  description = "APP Category Value"
  value       = "APP-1"
}

resource "nutanix_category_value" "APP-2" {
  name        = "${nutanix_category_key.test-category-key.id}"
  description = "APP Category Value"
  value       = "APP-2"
}

resource "nutanix_network_security_rule" "TEST-TIER" {
  name            = "RULE-1-TIERS-%d"
  description     = "tf-test-ports"
  app_rule_action = "MONITOR"

  app_rule_inbound_allow_list {
    ip_subnet               = "0.0.0.0"
    ip_subnet_prefix_length = "0"
    peer_specification_type = "IP_SUBNET"
    protocol                = "TCP"
    tcp_port_range_list {
      end_port   = 80
      start_port = 80
    }
    tcp_port_range_list {
      end_port   = 443
      start_port = 443
    }
  }
  app_rule_inbound_allow_list {
    filter_type = "CATEGORIES_MATCH_ALL"
    filter_params {
      name = nutanix_category_key.test-category-key.id
      values = [
        nutanix_category_value.APP-1.id
      ]
    }
    filter_kind_list        = ["vm"]
    peer_specification_type = "FILTER"
    protocol                = "TCP"
    tcp_port_range_list {
      end_port   = 22
      start_port = 22
    }
  }
  app_rule_inbound_allow_list {
    filter_type = "CATEGORIES_MATCH_ALL"
    filter_params {
      name = nutanix_category_key.test-category-key.id
      values = [
        nutanix_category_value.APP-2.id
      ]
    }

    filter_kind_list        = ["vm"]
    peer_specification_type = "FILTER"
    protocol                = "ICMP"
  }

  app_rule_target_group_default_internal_policy = "ALLOW_ALL"
  app_rule_target_group_filter_kind_list = [
    "vm"
  ]
  app_rule_target_group_filter_params {
    name = nutanix_category_key.test-category-key.id
    values = [
      nutanix_category_value.APP-1.id
    ]
  }
  app_rule_target_group_filter_params {
    name = data.nutanix_category_key.AppType.id
    values = [
      nutanix_category_value.DB.id
    ]
  }
  app_rule_target_group_filter_type             = "CATEGORIES_MATCH_ALL"
  app_rule_target_group_peer_specification_type = "FILTER"

  app_rule_outbound_allow_list {
    ip_subnet               = "10.0.0.0"
    ip_subnet_prefix_length = "24"
    peer_specification_type = "IP_SUBNET"
    protocol                = "UDP"
    udp_port_range_list {
      end_port   = 53
      start_port = 53
    }
  }
}

data "nutanix_network_security_rule" "test" {
	network_security_rule_id = "${nutanix_network_security_rule.TEST-TIER.id}"
}
`, r)
}
