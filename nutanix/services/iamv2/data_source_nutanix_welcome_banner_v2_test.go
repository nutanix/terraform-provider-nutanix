package iamv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_welcome_banner_v2 (singleton datasource):
//
// GetWelcomeBanner fetches the cluster-wide configured welcome banner. It is a
// read-only singleton config endpoint (no ext_id input, no Create/Update/Delete
// resource in this code-gen scope). Coverage:
//
//  1. Basic read: query the singleton datasource with no arguments and verify
//     the datasource is populated (id + is_enabled are always present in state).
//     content / created_time / last_updated_time are optional on the wire (only
//     returned once a banner is configured), so they are asserted with
//     TestCheckResourceAttrSet only when the banner is enabled — here we assert
//     the always-present attributes so the test is stable across clusters.
//  2. No-argument config: verify the singleton datasource accepts zero inputs
//     (no ext_id) and still resolves — this validates the singleton wiring.

const datasourceNameWelcomeBanner = "data.nutanix_welcome_banner_v2.test"

func TestAccV2NutanixWelcomeBannerDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testWelcomeBannerDatasourceV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameWelcomeBanner, "id"),
					resource.TestCheckResourceAttrSet(datasourceNameWelcomeBanner, "is_enabled"),
				),
			},
		},
	})
}

func testWelcomeBannerDatasourceV2Config() string {
	return `
		data "nutanix_welcome_banner_v2" "test" {}
	`
}
