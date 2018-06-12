package nutanix

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixNetworkSecurityRulesDataSource_basic(t *testing.T) {

	//Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityRulesSDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_network_security_rules.test", "entities.#", "1"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccNetworkSecurityRulesSDataSourceConfig = `
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

data "nutanix_network_security_rules" "test" {
	metadata = {
		length = 1
	}
}`

func Test_dataSourceNutanixNetworkSecurityRules(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixNetworkSecurityRules(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixNetworkSecurityRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixNetworkSecurityRulesRead(t *testing.T) {
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
			if err := dataSourceNutanixNetworkSecurityRulesRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixNetworkSecurityRulesRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceNetworkSecurityRulesSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceNetworkSecurityRulesSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceNetworkSecurityRulesSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
