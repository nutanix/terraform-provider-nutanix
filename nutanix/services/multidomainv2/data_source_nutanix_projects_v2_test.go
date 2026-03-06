package multidomainv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameProjectsV2 = "data.nutanix_projects_v2.test"

func TestAccV2NutanixProjectsDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-project-%d", r)
	description := "terraform test projects list datasource"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProjectV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectV2ResourceConfig(name, description) + testAccProjectsV2DatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameProjectsV2, "projects.#"),
					checkAttributeLength(dataSourceNameProjectsV2, "projects", 1),
				),
			},
		},
	})
}

func testAccProjectsV2DatasourceConfig() string {
	return `
data "nutanix_projects_v2" "test" {
  depends_on = [nutanix_project_v2.test]
}
`
}
