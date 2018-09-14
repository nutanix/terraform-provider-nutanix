package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixCategoryKey_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixCategoryKeyConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryKeyExists("nutanix_category_key.test"),
				),
			},
		},
	})
}

func TestAccNutanixCategoryKey_importBasic(t *testing.T) {
	resourceName := "nutanix_category_key.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixCategoryKeyConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryKeyExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixCategoryKeyExists(n string) resource.TestCheckFunc {
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

func testAccNutanixCategoryKeyConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "test"{
    name = "app-support-%d"
    description = "App Support CategoryKey"
}
`, r)
}
