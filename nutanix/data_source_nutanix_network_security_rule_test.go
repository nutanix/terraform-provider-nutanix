package nutanix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixNetworkSecurityRuleDataSource_basic(t *testing.T) {

	//Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityRuleDataSourceConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_network_security_rule.test", "name", "RULE-1-TIERS"),
					resource.TestCheckResourceAttr(
						"data.nutanix_network_security_rule.test", "app_rule_action", "APPLY"),
				),
			},
		},
	})
}

func testAccNetworkSecurityRuleDataSourceConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "test-category-key"{
    name = "TIER-1"
	description = "TIER Category Key"
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
    name = "${nutanix_category_key.test-category-key.id}"
	  description = "ashwini Category Value"
	 value = "ashwini-1"
}


resource "nutanix_network_security_rule" "TEST-TIER" {
  name        = "RULE-1-TIERS-%d"
  description = "rule 1 tiers"

  app_rule_action = "APPLY"

  app_rule_inbound_allow_list = [
    {
      peer_specification_type = "FILTER"
      filter_type             = "CATEGORIES_MATCH_ALL"
      filter_kind_list        = ["vm"]

      filter_params = [
        {
          name   = "${nutanix_category_key.test-category-key.id}"
          values = ["${nutanix_category_value.WEB.id}"]
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
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.APP.id}"]
    },
    {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.ashwini.id}"]
    },
  ]

  app_rule_outbound_allow_list = [
    {
      peer_specification_type = "FILTER"
      filter_type             = "CATEGORIES_MATCH_ALL"
      filter_kind_list        = ["vm"]

      filter_params = [
        {
          name   = "${nutanix_category_key.test-category-key.id}"
          values = ["${nutanix_category_value.DB.id}"]
        },
      ]
    },
  ]
}

data "nutanix_network_security_rule" "test" {
	network_security_rule_id = "${nutanix_network_security_rule.TEST-TIER.id}"
}
`, r)
}

func Test_dataSourceNutanixNetworkSecurityRule(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixNetworkSecurityRule(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixNetworkSecurityRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixNetworkSecurityRuleRead(t *testing.T) {
	type args struct {
		d    *schema.ResourceData
		meta interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dataSourceNutanixNetworkSecurityRuleRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixNetworkSecurityRuleRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceNetworkSecurityRuleSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceNetworkSecurityRuleSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceNetworkSecurityRuleSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccNetworkSecurityRuleDataSourceConfig(t *testing.T) {
	type args struct {
		r int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testAccNetworkSecurityRuleDataSourceConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccNetworkSecurityRuleDataSourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
