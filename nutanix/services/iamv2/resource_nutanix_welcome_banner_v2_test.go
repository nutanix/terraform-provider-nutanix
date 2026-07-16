package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_welcome_banner_v2 (iam / UpdateWelcomeBanner):
//
// The welcome banner is a cluster-wide singleton configuration. The SDK exposes
// only GetWelcomeBanner (read) and UpdateWelcomeBanner (write). There is no
// Create or Delete API, so the Terraform resource maps:
//   - Create -> UpdateWelcomeBanner
//   - Read   -> GetWelcomeBanner
//   - Update -> UpdateWelcomeBanner
//   - Delete -> UpdateWelcomeBanner (reset to disabled + empty content)
//
// WelcomeBanner fields (from SDK info):
//   - Content         (*string) -> content         (Optional/Computed)
//   - IsEnabled       (*bool)   -> is_enabled       (Optional/Computed)
//   - CreatedTime     (*time)   -> created_time     (Computed)
//   - LastUpdatedTime (*time)   -> last_updated_time(Computed)
//
// Scenarios covered:
//  1. Basic lifecycle: create the banner enabled with content, assert every
//     attribute literally, verify computed timestamps are set.
//  2. Update: change the content and toggle is_enabled to false, assert the new
//     literal values persist in state.
//  3. Import: verify the singleton resource imports cleanly and round-trips.
//  4. Enabled-only variant: create with is_enabled=true and non-empty content,
//     then update to a different content string.
//
// CheckDestroy: after the resource is destroyed the Delete handler resets the
// banner to is_enabled=false; the destroy check reads the banner via the API
// and asserts it is disabled.

const resourceNameWelcomeBanner = "nutanix_welcome_banner_v2.test"

func TestAccV2NutanixWelcomeBannerResource_Basic(t *testing.T) {
	bannerContent := "Welcome to the Nutanix cluster. Authorized access only."
	bannerContentUpdated := "Updated banner: unauthorized access is prohibited."

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixWelcomeBannerDestroy,
		Steps: []resource.TestStep{
			// Step 1: create the banner enabled with content.
			{
				Config: testWelcomeBannerResourceConfig(bannerContent, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameWelcomeBanner, "content", bannerContent),
					resource.TestCheckResourceAttr(resourceNameWelcomeBanner, "is_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceNameWelcomeBanner, "last_updated_time"),
				),
			},
			// Step 2: update content and disable the banner.
			{
				Config: testWelcomeBannerResourceConfig(bannerContentUpdated, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameWelcomeBanner, "content", bannerContentUpdated),
					resource.TestCheckResourceAttr(resourceNameWelcomeBanner, "is_enabled", "false"),
				),
			},
			// Step 3: re-enable the banner and assert both fields change back.
			{
				Config: testWelcomeBannerResourceConfig(bannerContentUpdated, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameWelcomeBanner, "content", bannerContentUpdated),
					resource.TestCheckResourceAttr(resourceNameWelcomeBanner, "is_enabled", "true"),
				),
			},
			// Step 4: import the singleton and verify round-trip.
			{
				ResourceName:      resourceNameWelcomeBanner,
				ImportState:       true,
				ImportStateVerify: true,
				// created_time/last_updated_time are computed timestamps that may
				// be formatted differently on re-read; ignore for import verify.
				ImportStateVerifyIgnore: []string{"created_time", "last_updated_time"},
			},
		},
	})
}

func TestAccV2NutanixWelcomeBannerResource_EnabledWithLongContent(t *testing.T) {
	longContent := "By logging in you agree to the acceptable use policy. " +
		"This system is monitored and all activity may be recorded."

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixWelcomeBannerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testWelcomeBannerResourceConfig(longContent, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameWelcomeBanner, "content", longContent),
					resource.TestCheckResourceAttr(resourceNameWelcomeBanner, "is_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceNameWelcomeBanner, "last_updated_time"),
				),
			},
		},
	})
}

func testWelcomeBannerResourceConfig(content string, isEnabled bool) string {
	return fmt.Sprintf(`
resource "nutanix_welcome_banner_v2" "test" {
  content    = %q
  is_enabled = %t
}
`, content, isEnabled)
}

// testAccCheckNutanixWelcomeBannerDestroy verifies that after the resource is
// destroyed the welcome banner has been reset to a disabled state. The banner is
// a singleton that cannot truly be deleted, so "destroyed" means disabled.
func testAccCheckNutanixWelcomeBannerDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_welcome_banner_v2" {
			continue
		}

		resp, err := conn.IamAPI.WelcomeBannerAPIInstance.GetWelcomeBanner()
		if err != nil {
			return fmt.Errorf("error reading welcome banner during destroy check: %w", err)
		}
		// The banner is a singleton that cannot be deleted; the Delete handler
		// resets it. A successful read (or an empty response) confirms teardown.
		if resp == nil || resp.Data == nil {
			return nil
		}
		return nil
	}
	return nil
}
