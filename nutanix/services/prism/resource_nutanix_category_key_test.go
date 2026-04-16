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

func TestAccNutanixCategoryKey_basic(t *testing.T) {
	resourceName := "nutanix_category_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
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

func TestAccNutanixCategoryKey_update(t *testing.T) {
	resourceName := "nutanix_category_key.test_update"
	rInt := acctest.RandIntRange(0, 500)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixCategoryKeyConfigToUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("app-support-%d", rInt)),
					resource.TestCheckResourceAttr(resourceName, "description", "App Support CategoryKey"),
				),
			},
			{
				Config: testAccNutanixCategoryKeyConfigUpdated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryKeyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("app-support-%d-updated", rInt)),
					resource.TestCheckResourceAttr(resourceName, "description", "App Support CategoryKey Updated"),
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
	conn := acc.TestAccProvider.Meta().(*conns.Client)

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

func testAccNutanixCategoryKeyConfigToUpdate(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "test_update"{
    name = "app-support-%d"
    description = "App Support CategoryKey"
}
`, r)
}

func testAccNutanixCategoryKeyConfigUpdated(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "test_update"{
    name = "app-support-%d-updated"
    description = "App Support CategoryKey Updated"
}
`, r)
}
