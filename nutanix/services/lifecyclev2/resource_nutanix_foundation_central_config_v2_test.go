package lifecyclev2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_foundation_central_config_v2
// --------------------------------------------------
// Foundation Central exposes a single, cluster-wide configuration object (no
// entity id, no Create/Delete). The resource is therefore a singleton: Create
// applies the desired config via an Update task, Read fetches it back, Update
// re-applies, Delete is a no-op. Coverage below:
//
//  1. TestAccV2NutanixFoundationCentralConfigResource_Basic
//     - Step 1 (Create): set both writable scalars (ahv_installation_timeout_minutes,
//       aos_download_timeout_minutes) and assert the literal values on the resource.
//       Also read the config back via the singular datasource and assert every
//       attribute matches the resource (TestCheckResourceAttrPair), including the
//       server-computed commit_id / tls_certificate_fingerprint / version.
//     - Step 2 (Update): change BOTH writable scalars to new values and assert the
//       new literals on the resource and on a freshly-read datasource.
//
//  2. TestAccV2NutanixFoundationCentralConfigDataSource_Basic
//     - Read the singleton config purely via the datasource (no resource) and
//       assert the computed attributes are present/consistent.
//
// Note: Import is NOT supported for this resource (singleton, synthetic id) — see
// comment on the ImportState omission below. Delete is a no-op; there is no
// GetById to assert against in a CheckDestroy, so (matching the LCM config
// singleton anchor) no CheckDestroy is defined.

const (
	resourceNameFCConfig   = "nutanix_foundation_central_config_v2.test"
	datasourceNameFCConfig = "data.nutanix_foundation_central_config_v2.test"
)

func TestAccV2NutanixFoundationCentralConfigResource_Basic(t *testing.T) {
	ahvTimeout := acctest.RandIntRange(30, 120)
	aosTimeout := acctest.RandIntRange(30, 120)
	ahvTimeoutUpdated := ahvTimeout + 15
	aosTimeoutUpdated := aosTimeout + 15

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create / apply the config with both writable scalars set.
			{
				Config: testFoundationCentralConfig(ahvTimeout, aosTimeout),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameFCConfig, "ahv_installation_timeout_minutes", fmt.Sprintf("%d", ahvTimeout)),
					resource.TestCheckResourceAttr(resourceNameFCConfig, "aos_download_timeout_minutes", fmt.Sprintf("%d", aosTimeout)),
					// Computed / server-populated attributes must be present after apply.
					resource.TestCheckResourceAttrSet(resourceNameFCConfig, "version"),
					// Datasource must round-trip every attribute from the resource.
					resource.TestCheckResourceAttrPair(datasourceNameFCConfig, "ahv_installation_timeout_minutes", resourceNameFCConfig, "ahv_installation_timeout_minutes"),
					resource.TestCheckResourceAttrPair(datasourceNameFCConfig, "aos_download_timeout_minutes", resourceNameFCConfig, "aos_download_timeout_minutes"),
					resource.TestCheckResourceAttrPair(datasourceNameFCConfig, "commit_id", resourceNameFCConfig, "commit_id"),
					resource.TestCheckResourceAttrPair(datasourceNameFCConfig, "tls_certificate_fingerprint", resourceNameFCConfig, "tls_certificate_fingerprint"),
					resource.TestCheckResourceAttrPair(datasourceNameFCConfig, "version", resourceNameFCConfig, "version"),
				),
			},
			// Step 2: Update both writable scalars and assert the new literal values.
			{
				Config: testFoundationCentralConfig(ahvTimeoutUpdated, aosTimeoutUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameFCConfig, "ahv_installation_timeout_minutes", fmt.Sprintf("%d", ahvTimeoutUpdated)),
					resource.TestCheckResourceAttr(resourceNameFCConfig, "aos_download_timeout_minutes", fmt.Sprintf("%d", aosTimeoutUpdated)),
					resource.TestCheckResourceAttrPair(datasourceNameFCConfig, "ahv_installation_timeout_minutes", resourceNameFCConfig, "ahv_installation_timeout_minutes"),
					resource.TestCheckResourceAttrPair(datasourceNameFCConfig, "aos_download_timeout_minutes", resourceNameFCConfig, "aos_download_timeout_minutes"),
				),
			},
			// Import not supported for this resource (singleton with a synthetic id).
		},
	})
}

func TestAccV2NutanixFoundationCentralConfigDataSource_Basic(t *testing.T) {
	datasourceName := "data.nutanix_foundation_central_config_v2.read_only"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFoundationCentralConfigDataSourceOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "version"),
					resource.TestCheckResourceAttrSet(datasourceName, "ahv_installation_timeout_minutes"),
					resource.TestCheckResourceAttrSet(datasourceName, "aos_download_timeout_minutes"),
				),
			},
		},
	})
}

func testFoundationCentralConfig(ahvTimeout, aosTimeout int) string {
	return fmt.Sprintf(`
resource "nutanix_foundation_central_config_v2" "test" {
  ahv_installation_timeout_minutes = %[1]d
  aos_download_timeout_minutes     = %[2]d
}

data "nutanix_foundation_central_config_v2" "test" {
  depends_on = [nutanix_foundation_central_config_v2.test]
}
`, ahvTimeout, aosTimeout)
}

func testFoundationCentralConfigDataSourceOnly() string {
	return `
data "nutanix_foundation_central_config_v2" "read_only" {}
`
}
