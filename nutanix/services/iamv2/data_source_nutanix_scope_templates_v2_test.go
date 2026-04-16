package iamv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameScopeTemplates = "data.nutanix_scope_templates_v2.test"

func TestAccV2NutanixScopeTemplatesDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testScopeTemplatesDatasourceV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameScopeTemplates, "scope_templates.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixScopeTemplatesDatasource_WithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testScopeTemplatesDatasourceV2WithFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameScopeTemplates, "scope_templates.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixScopeTemplatesDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testScopeTemplatesDatasourceV2WithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameScopeTemplates, "scope_templates.#"),
					resource.TestCheckResourceAttr(datasourceNameScopeTemplates, "scope_templates.#", "0"),
				),
			},
		},
	})
}

func testScopeTemplatesDatasourceV2Config() string {
	return `
		data "nutanix_scope_templates_v2" "test" {}
	`
}

func testScopeTemplatesDatasourceV2WithFilterConfig() string {
	return `
		data "nutanix_scope_templates_v2" "test" {
			filter = "displayName ne 'nonexistent_filter_value_12345'"
		}
	`
}

func testScopeTemplatesDatasourceV2WithInvalidFilterConfig() string {
	return `
		data "nutanix_scope_templates_v2" "test" {
			filter = "displayName eq 'invalid_filter_no_match_12345'"
		}
	`
}
