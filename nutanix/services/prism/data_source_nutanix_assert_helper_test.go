package prism_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixAssertHelperDS(t *testing.T) {
	name := "checks"
	errorMsg := "Error message for nutanix assert helper"
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNutanixAssertHelperDS(name, "false", errorMsg),
				ExpectError: regexp.MustCompile(errorMsg),
			},
			{
				Config: testAccNutanixAssertHelperDS(name, "true", errorMsg),
			},
		},
	})
}

func testAccNutanixAssertHelperDS(name, condition, errMsg string) string {
	return fmt.Sprintf(`
	data "nutanix_assert_helper" "%s" {
		checks {
			condition = %s
			error_message = "%s"
		}
	}
	`, name, condition, errMsg)
}
