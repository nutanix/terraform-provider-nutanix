package monitoringv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSdaPolicy = "data.nutanix_sda_policy_v2.test"

func TestAccV2NutanixSdaPolicyDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSdaPolicyDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSdaPolicy, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameSdaPolicy, "name"),
					resource.TestCheckResourceAttrSet(datasourceNameSdaPolicy, "title"),
				),
			},
		},
	})
}

func testSdaPolicyDatasourceConfig() string {
	return fmt.Sprintf(`
data "nutanix_sda_policies_v2" "list" {}

data "nutanix_sda_policy_v2" "test" {
  ext_id = data.nutanix_sda_policies_v2.list.sda_policies.0.ext_id
}
`)
}
