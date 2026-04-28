package monitoringv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameAudit = "data.nutanix_audit_v2.test"

func TestAccV2NutanixAuditDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuditDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameAudit, "ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudit, "audit_type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudit, "operation_type"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudit, "status"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudit, "creation_time"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudit, "service_name"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudit, "source_entity.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameAudit, "user_reference.#"),
				),
			},
		},
	})
}

func testAuditDatasourceConfig() string {
	return fmt.Sprintf(`
data "nutanix_audits_v2" "audits" {
  limit = 1
}

data "nutanix_audit_v2" "test" {
  ext_id = data.nutanix_audits_v2.audits.audits.0.ext_id
  depends_on = [data.nutanix_audits_v2.audits]
}
`)
}
