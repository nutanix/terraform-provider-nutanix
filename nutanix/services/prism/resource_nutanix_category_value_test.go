package prism_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixCategoryValue_basic(t *testing.T) {
	rInt := acctest.RandIntRange(0, 200)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixCategoryValueConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryValueExists("nutanix_category_value.test"),
				),
			},
		},
	})
}

func TestAccNutanixCategoryValue_update(t *testing.T) {
	rInt := acctest.RandIntRange(201, 500)
	resourceName := "nutanix_category_value.test_update"
	description := "Test Category Value"
	value := "test-value"
	descriptionUpdated := fmt.Sprintf("%s Updated", description)
	valueUpdated := fmt.Sprintf("%s-updated", value)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixCategoryValueConfigToUpdate(rInt, value, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryValueExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccNutanixCategoryValueConfigToUpdate(rInt, valueUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryValueExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "value", valueUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
				),
			},
		},
	})
}

func testAccCheckNutanixCategoryValueExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixCategoryValueDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_category_value" {
			continue
		}

		for {
			_, err := conn.API.V3.GetCategoryValue(rs.Primary.Attributes["name"], rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "CATEGORY_NAME_VALUE_MISMATCH") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}
	}

	return nil
}

func getCategoryValueResource(rInt int) string {
	return fmt.Sprintf(`
	resource "nutanix_category_key" "test-category-key"{
    name = "app-suppport-%d"
	description = "App Support Category Key"
}
`, rInt)
}

func testAccNutanixCategoryValueConfig(rInt int) string {
	return getCategoryValueResource(rInt) +
		`
resource "nutanix_category_value" "test"{
    name = nutanix_category_key.test-category-key.id
	description = "Test Category Value"
	value = "test-value"
}
`
}

func testAccNutanixCategoryValueConfigToUpdate(rInt int, value, description string) string {
	return getCategoryValueResource(rInt) +
		fmt.Sprintf(`
resource "nutanix_category_value" "test_update"{
    name = nutanix_category_key.test-category-key.id
	value = "%s"
	description = "%s"
}
`, value, description)
}
