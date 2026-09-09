package lifecyclev2_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// testAccCheckClaimTokenDestroy verifies that every claim token created during the
// test has been removed. It calls GetClaimTokenById for each resource and expects
// an error (the token no longer exists).
func testAccCheckClaimTokenDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).LifecycleAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_claim_token_v2" {
			continue
		}
		if _, err := conn.ClaimTokensAPIInstance.GetClaimTokenById(utils.StringPtr(rs.Primary.ID)); err == nil {
			return fmt.Errorf("claim token %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// testAccCheckPatchedImageDestroy verifies that every patched image created during
// the test has been removed. Any claim tokens created for the test are also checked.
func testAccCheckPatchedImageDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).LifecycleAPI

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "nutanix_patched_image_v2":
			if _, err := conn.PatchedImagesAPIInstance.GetPatchedImageById(utils.StringPtr(rs.Primary.ID)); err == nil {
				return fmt.Errorf("patched image %s still exists", rs.Primary.ID)
			}
		case "nutanix_claim_token_v2":
			if _, err := conn.ClaimTokensAPIInstance.GetClaimTokenById(utils.StringPtr(rs.Primary.ID)); err == nil {
				return fmt.Errorf("claim token %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}
