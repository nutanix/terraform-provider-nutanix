package iamv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameScopeTemplate = "data.nutanix_scope_template_v2.test"

func TestAccV2NutanixScopeTemplateDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testScopeTemplateDatasourceV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameScopeTemplate, "display_name"),
					resource.TestCheckResourceAttrSet(datasourceNameScopeTemplate, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameScopeTemplate, "entities.0.entity_filter"),
					resource.TestCheckResourceAttr(datasourceNameScopeTemplate, "description", "Scope Template for Projects_2.0"),
					resource.TestCheckResourceAttr(datasourceNameScopeTemplate, "display_name", "ProjectsScopeTemplate"),
				),
			},
		},
	})
}

func testScopeTemplateDatasourceV2Config() string {
	return `
		data "nutanix_scope_templates_v2" "list" {}

		data "nutanix_scope_template_v2" "test" {
			ext_id = data.nutanix_scope_templates_v2.list.scope_templates.0.ext_id
		}
	`
}
