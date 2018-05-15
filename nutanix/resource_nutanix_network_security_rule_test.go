package nutanix

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixNetworkSecurityRule_basic(t *testing.T) {
	r := rand.Int31()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixNetworkSecurityRuleConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists("nutanix_network_security_rule.TEST-TIER"),
				),
			},
			{
				Config: testAccNutanixNetworkSecurityRuleConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixNetworkSecurityRuleExists("nutanix_network_security_rule.TEST-TIER"),
				),
			},
		},
	})
}

func testAccCheckNutanixNetworkSecurityRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixNetworkSecurityRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*NutanixClient)

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

func testAccNutanixNetworkSecurityRuleConfig(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_network_security_rule" "TEST-TIER" {
  name        = "RULE-1-TIERS"
  description = "rule 1 tiers"

  app_rule_action = "APPLY"

  app_rule_inbound_allow_list = [
    {
      peer_specification_type = "FILTER"
      filter_type             = "CATEGORIES_MATCH_ALL"
      filter_kind_list        = ["vm"]

      filter_params = [
        {
          name   = "TIER"
          values = ["WEB"]
        },
      ]
    },
  ]

  app_rule_target_group_default_internal_policy = "DENY_ALL"

  app_rule_target_group_peer_specification_type = "FILTER"

  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"

  app_rule_target_group_filter_kind_list = ["vm"]

  app_rule_target_group_filter_params = [
    {
      name   = "TIER"
      values = ["APP"]
    },
    {
      name   = "TIER"
      values = ["ashwini"]
    },
  ]

  app_rule_outbound_allow_list = [
    {
      peer_specification_type = "FILTER"
      filter_type             = "CATEGORIES_MATCH_ALL"
      filter_kind_list        = ["vm"]

      filter_params = [
        {
          name   = "TIER"
          values = ["DB"]
        },
      ]
    },
  ]
}
`)
}

func testAccNutanixNetworkSecurityRuleConfigUpdate(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_network_security_rule" "TEST-TIER" {
  name        = "RULE-1-TIERS-1"
  description = "rule 1 tiers 1"

  app_rule_action = "APPLY"

  app_rule_inbound_allow_list = [
    {
      peer_specification_type = "FILTER"
      filter_type             = "CATEGORIES_MATCH_ALL"
      filter_kind_list        = ["vm"]

      filter_params = [
        {
          name   = "TIER"
          values = ["WEB"]
        },
      ]
    },
  ]

  app_rule_target_group_default_internal_policy = "DENY_ALL"

  app_rule_target_group_peer_specification_type = "FILTER"

  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"

  app_rule_target_group_filter_kind_list = ["vm"]

  app_rule_target_group_filter_params = [
    {
      name   = "TIER"
      values = ["APP"]
    },
    {
      name   = "TIER"
      values = ["ashwini"]
    },
  ]

  app_rule_outbound_allow_list = [
    {
      peer_specification_type = "FILTER"
      filter_type             = "CATEGORIES_MATCH_ALL"
      filter_kind_list        = ["vm"]

      filter_params = [
        {
          name   = "TIER"
          values = ["DB"]
        },
      ]
    },
  ]
}
`)
}
