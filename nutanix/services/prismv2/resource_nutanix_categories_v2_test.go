package prismv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameCategory = "nutanix_category_v2.test"

func TestAccV2NutanixCategoryResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	value := fmt.Sprintf("test category value-%d", r)
	desc := "test category description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryV2Config(r, value, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "key", fmt.Sprintf("test-cat-%d", r)),
					resource.TestCheckResourceAttr(resourceNameCategory, "value", value),
					resource.TestCheckResourceAttr(resourceNameCategory, "description", desc),
					resource.TestCheckResourceAttr(resourceNameCategory, "type", "USER"),
				),
			},
		},
	})
}

func TestAccV2NutanixCategoryResource_Update(t *testing.T) {
	r := acctest.RandInt()
	value := fmt.Sprintf("test category value-%d", r)
	desc := "test category description"
	updatedValue := fmt.Sprintf("test category value updated-%d", r)
	updateDesc := "test category description updated"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryV2Config(r, value, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "key", fmt.Sprintf("test-cat-%d", r)),
					resource.TestCheckResourceAttr(resourceNameCategory, "value", value),
					resource.TestCheckResourceAttr(resourceNameCategory, "description", desc),
					resource.TestCheckResourceAttr(resourceNameCategory, "type", "USER"),
				),
			},
			{
				Config: testAccCategoryV2Config(r, updatedValue, updateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "key", fmt.Sprintf("test-cat-%d", r)),
					resource.TestCheckResourceAttr(resourceNameCategory, "value", updatedValue),
					resource.TestCheckResourceAttr(resourceNameCategory, "description", updateDesc),
					resource.TestCheckResourceAttr(resourceNameCategory, "type", "USER"),
				),
			},
		},
	})
}

func TestAccV2NutanixCategoryResource_WithNoKey(t *testing.T) {
	r := acctest.RandInt()
	value := fmt.Sprintf("test category value-%d", r)
	desc := "test category description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCategoryV2ConfigWithNoKey(r, value, desc),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testAccCategoryV2Config(r int, val, desc string) string {
	return fmt.Sprintf(`
	resource "nutanix_category_v2" "test" {
		key = "test-cat-%d"
		value = "%[2]s"
		description = "%[3]s"
	}
`, r, val, desc)
}

func testAccCategoryV2ConfigWithNoKey(r int, val, desc string) string {
	return fmt.Sprintf(`
	resource "nutanix_category_v2" "test" {
		value = "%[2]s"
		description = "%[3]s"
	  }
`, r, val, desc)
}
