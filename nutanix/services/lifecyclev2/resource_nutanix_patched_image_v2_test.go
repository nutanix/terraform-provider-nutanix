package lifecyclev2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_patched_image_v2 (ClaimTokens cross-resource consumer)
// ----------------------------------------------------------------------------
// A patched image is created from a base host image and customized for specific
// nodes. It REFERENCES the claim token via claim_token_ext_id — this is the
// ClaimTokens cross-resource reference under test. The integration test chains:
// create claim token -> create patched image referencing it -> read the
// claim_token_ext_id back through the resource Read AND both datasources.
//
//  1. TestAccV2NutanixPatchedImageResource_Basic — end-to-end cross-resource:
//     requires an uploaded host image ext_id and node ext_id from
//     test_config_v2.json (lifecycle.patched_image). Skipped when not configured.
//  2. TestAccV2NutanixPatchedImageResource_MissingRequired — negative: omit
//     claim_token_ext_id and assert Terraform requires it.

const (
	resourceNamePatchedImage    = "nutanix_patched_image_v2.test"
	datasourceNamePatchedImage  = "data.nutanix_patched_image_v2.test"
	datasourceNamePatchedImages = "data.nutanix_patched_images_v2.list"
)

func TestAccV2NutanixPatchedImageResource_Basic(t *testing.T) {
	imageExtID := testVars.Lifecycle.PatchedImage.HostImageExtID
	nodeExtID := testVars.Lifecycle.PatchedImage.NodeExtID
	if imageExtID == "" || nodeExtID == "" {
		t.Skip("skipping: lifecycle.patched_image.host_image_ext_id / node_ext_id not set in test_config_v2.json")
	}

	r := acctest.RandInt()
	tokenName := fmt.Sprintf("tf-claim-pi-%d", r)
	imageName := fmt.Sprintf("tf-patched-image-%d", r)
	expiry := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckPatchedImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPatchedImageConfig(tokenName, expiry, imageName, imageExtID, nodeExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePatchedImage, "name", imageName),
					resource.TestCheckResourceAttr(resourceNamePatchedImage, "host_type", "AHV"),
					resource.TestCheckResourceAttrSet(resourceNamePatchedImage, "ext_id"),
					// Cross-resource: claim_token_ext_id wired from the claim token.
					resource.TestCheckResourceAttrPair(resourceNamePatchedImage, "claim_token_ext_id", resourceNameClaimToken, "ext_id"),
					resource.TestCheckResourceAttr(resourceNamePatchedImage, "image_details.0.local_host_image_ext_id", imageExtID),
					resource.TestCheckResourceAttr(resourceNamePatchedImage, "node_configurations.0.node_ext_id", nodeExtID),
					// Cross-resource read-back via singular + plural datasources.
					resource.TestCheckResourceAttrPair(datasourceNamePatchedImage, "claim_token_ext_id", resourceNamePatchedImage, "claim_token_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNamePatchedImages, "patched_images.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixPatchedImageResource_MissingRequired(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "nutanix_patched_image_v2" "test" {
  name      = "missing-claim-token"
  host_type = "AHV"
  image_details {
    local_host_image_ext_id = "00000000-0000-0000-0000-000000000000"
  }
  node_configurations {
    node_ext_id = "00000000-0000-0000-0000-000000000001"
  }
}
`,
				ExpectError: regexp.MustCompile("The argument \"claim_token_ext_id\" is required"),
			},
		},
	})
}

func testPatchedImageConfig(tokenName, expiry, imageName, imageExtID, nodeExtID string) string {
	return fmt.Sprintf(`
resource "nutanix_claim_token_v2" "test" {
  name            = "%[1]s"
  expiry_time     = "%[2]s"
  max_usage_count = 5
}

resource "nutanix_patched_image_v2" "test" {
  claim_token_ext_id = nutanix_claim_token_v2.test.ext_id
  name               = "%[3]s"
  host_type          = "AHV"

  image_details {
    local_host_image_ext_id = "%[4]s"
  }

  node_configurations {
    node_ext_id = "%[5]s"
  }
}

data "nutanix_patched_image_v2" "test" {
  ext_id = nutanix_patched_image_v2.test.ext_id
}

data "nutanix_patched_images_v2" "list" {
  depends_on = [nutanix_patched_image_v2.test]
}
`, tokenName, expiry, imageName, imageExtID, nodeExtID)
}
