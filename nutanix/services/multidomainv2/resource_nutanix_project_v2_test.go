package multidomainv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccV2NutanixProjectResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-project-%d", r)
	description := "terraform test project CRUD"
	updateName := fmt.Sprintf("tf-test-project-%d-update", r)
	updateDescription := "terraform test project CRUD update"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProjectV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectV2ResourceConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProjectV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "description", description),
				),
			},
			{
				Config: testAccProjectV2ResourceConfig(updateName, updateDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProjectV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "name", updateName),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "description", updateDescription),
				),
			},
		},
	})
}

func testAccProjectV2ResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_project_v2" "test" {
  name        = "%s"
  description = "%s"
}
`, name, description)
}
