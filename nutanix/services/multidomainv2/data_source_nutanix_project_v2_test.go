package multidomainv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameProjectV2 = "data.nutanix_project_v2.test"

func TestAccV2NutanixProjectDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-project-%d", r)
	description := "terraform test project datasource"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProjectV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectV2ResourceConfig(name, description) + testAccProjectV2DatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameProjectV2, "ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameProjectV2, "name", name),
					resource.TestCheckResourceAttr(dataSourceNameProjectV2, "description", description),
				),
			},
		},
	})
}

func testAccProjectV2DatasourceConfig() string {
	return `
data "nutanix_project_v2" "test" {
  ext_id = nutanix_project_v2.test.id
  depends_on = [nutanix_project_v2.test]
}
`
}
