package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameAudits = "data.nutanix_audits_v2.test"

func TestAccV2NutanixAuditsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuditsDatasourceBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.0.ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.0.audit_type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.0.operation_type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.0.creation_time"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.0.service_name"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.0.links.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixAuditsDatasource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuditsDatasourceWithLimitConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.0.ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixAuditsDatasource_WithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuditsDatasourceWithFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameAudits, "audits.#"),
				),
			},
		},
	})
}

func testAuditsDatasourceBasicConfig() string {
	return `
data "nutanix_audits_v2" "test" {}
`
}

func testAuditsDatasourceWithLimitConfig() string {
	return `
data "nutanix_audits_v2" "test" {
  limit = 5
}
`
}

func testAuditsDatasourceWithFilterConfig() string {
	return `
data "nutanix_audits_v2" "test" {
  filter = "serviceName eq 'Nutanix'"
  limit  = 5
}
`
}
