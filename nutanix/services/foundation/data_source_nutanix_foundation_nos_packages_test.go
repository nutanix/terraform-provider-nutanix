package foundation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFoundationNosPackagesDataSource(t *testing.T) {
	name := "nos_packages"
	resourcePath := "data.nutanix_foundation_nos_packages." + name
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNosPackagesDSConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourcePath, "entities.0"),
				),
			},
		},
	})
}

func testNosPackagesDSConfig(name string) string {
	return fmt.Sprintf(`data "nutanix_foundation_nos_packages" "%s" {}`, name)
}
