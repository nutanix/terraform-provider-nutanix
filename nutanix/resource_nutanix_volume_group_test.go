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

func TestAccNutanixVolumeGroup_basic(t *testing.T) {
	// skipping as this API is not yet GA (will GA in upcoming AOS release)
	t.Skip()

	r := rand.Int31()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVolumeGroupConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVolumeGroupExists("nutanix_volume_group.test_volume"),
				),
			},
			{
				Config: testAccNutanixVolumeGroupConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVolumeGroupExists("nutanix_volume_group.test_volume"),
				),
			},
		},
	})
}

func testAccCheckNutanixVolumeGroupExists(n string) resource.TestCheckFunc {
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

func testAccCheckNutanixVolumeGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_volume_group" {
			continue
		}
		for {
			_, err := conn.API.V3.GetVolumeGroup(rs.Primary.ID)
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

func testAccNutanixVolumeGroupConfig(r int32) string {
	return fmt.Sprintf(` 

resource "nutanix_volume_group" "test_volume" {
  name        = "Test Volume Group"
  description = "Tes Volume Group Description"
}
`)
}

func testAccNutanixVolumeGroupConfigUpdate(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_volume_group" "test_volume" {
  name        = "Test Volume Group Update"
  description = "Tes Volume Group Description Update"
}
`)
}
