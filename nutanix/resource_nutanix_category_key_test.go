package nutanix

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixCategoryKey_basic(t *testing.T) {
	r := rand.Int31()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNutanixCategoryKeyConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryKeyExists("nutanix_category_key.test"),
				),
			},
		},
	})
}

func testAccCheckNutanixCategoryKeyExists(n string) resource.TestCheckFunc {
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

func testAccCheckNutanixCategoryKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_category_key" {
			continue
		}
		for {
			_, err := conn.API.V3.GetCategoryKey(rs.Primary.ID)
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

func testAccNutanixCategoryKeyConfig(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "test"{
    name = "app-support-%d"
    description = "App Support CategoryKey"
}
`, r)
}
