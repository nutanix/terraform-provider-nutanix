package nutanix

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func TestAccNutanixNetworkSecurityRule_basic(t *testing.T) {
	//Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

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

func testAccNutanixNetworkSecurityRuleConfig(r int32) string {
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
      name   = "${nutanix_category_key.USER.id}"
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
`)
}

func testAccNutanixNetworkSecurityRuleConfigUpdate(r int32) string {
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
  name        = "RULE-1-TIERS Updated"
  description = "rule 1 tiers Updated"

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

  app_rule_target_group_filter_params = {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.APP.id}"]
  }

  app_rule_target_group_filter_params = {
      name   = "${nutanix_category_key.USER.id}"
      values = ["${nutanix_category_value.ashwini.id}"]
  }


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
`)
}

func Test_resourceNutanixNetworkSecurityRule(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceNutanixNetworkSecurityRule(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceNutanixNetworkSecurityRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceNutanixNetworkSecurityRuleCreate(t *testing.T) {
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
			if err := resourceNutanixNetworkSecurityRuleCreate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixNetworkSecurityRuleCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixNetworkSecurityRuleRead(t *testing.T) {
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
			if err := resourceNutanixNetworkSecurityRuleRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixNetworkSecurityRuleRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixNetworkSecurityRuleUpdate(t *testing.T) {
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
			if err := resourceNutanixNetworkSecurityRuleUpdate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixNetworkSecurityRuleUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixNetworkSecurityRuleDelete(t *testing.T) {
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
			if err := resourceNutanixNetworkSecurityRuleDelete(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixNetworkSecurityRuleDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixNetworkSecurityRuleExists(t *testing.T) {
	type args struct {
		conn *v3.Client
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resourceNutanixNetworkSecurityRuleExists(tt.args.conn, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixNetworkSecurityRuleExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resourceNutanixNetworkSecurityRuleExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNetworkSecurityRuleResources(t *testing.T) {
	type args struct {
		d                   *schema.ResourceData
		networkSecurityRule *v3.NetworkSecurityRuleResources
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
			if err := getNetworkSecurityRuleResources(tt.args.d, tt.args.networkSecurityRule); (err != nil) != tt.wantErr {
				t.Errorf("getNetworkSecurityRuleResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_networkSecurityRuleStateRefreshFunc(t *testing.T) {
	type args struct {
		client *v3.Client
		uuid   string
	}
	tests := []struct {
		name string
		args args
		want resource.StateRefreshFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := networkSecurityRuleStateRefreshFunc(tt.args.client, tt.args.uuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("networkSecurityRuleStateRefreshFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expandFilterParams(t *testing.T) {
	type args struct {
		fp map[string][]string
	}
	tests := []struct {
		name string
		args args
		want []map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expandFilterParams(tt.args.fp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expandFilterParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNetworkSecurityRuleSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNetworkSecurityRuleSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNetworkSecurityRuleSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixNetworkSecurityRuleExists(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name string
		args args
		want resource.TestCheckFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testAccCheckNutanixNetworkSecurityRuleExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixNetworkSecurityRuleExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixNetworkSecurityRuleDestroy(t *testing.T) {
	type args struct {
		s *terraform.State
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
			if err := testAccCheckNutanixNetworkSecurityRuleDestroy(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("testAccCheckNutanixNetworkSecurityRuleDestroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_testAccNutanixNetworkSecurityRuleConfig(t *testing.T) {
	type args struct {
		r int32
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
			if got := testAccNutanixNetworkSecurityRuleConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixNetworkSecurityRuleConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccNutanixNetworkSecurityRuleConfigUpdate(t *testing.T) {
	type args struct {
		r int32
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
			if got := testAccNutanixNetworkSecurityRuleConfigUpdate(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixNetworkSecurityRuleConfigUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
