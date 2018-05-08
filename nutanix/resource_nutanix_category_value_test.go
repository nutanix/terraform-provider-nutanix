package nutanix

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixCategoryValue_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryValueDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNutanixCategoryValueConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryValueExists("nutanix_category_value.test"),
				),
			},
		},
	})
}

func testAccCheckNutanixCategoryValueExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixCategoryValueDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*NutanixClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_category_value" {
			continue
		}
		for {
			log.Println(rs.Primary.Attributes)
			_, err := conn.API.V3.GetCategoryValue("app-suppport-1", rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}

	}

	return nil
}

func testAccNutanixCategoryValueConfig() string {
	return `
resource "nutanix_category_key" "test-category-key"{
    name = "app-suppport-1"
	description = "App Support Category Key"
}


resource "nutanix_category_value" "test"{
    name = "${test-category-key.name}"
	description = "Test Category Value"
	value = "test-value"
}
`
}
