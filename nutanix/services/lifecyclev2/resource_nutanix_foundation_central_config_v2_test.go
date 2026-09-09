package lifecyclev2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_foundation_central_config_v2
// ---------------------------------------------------
// The Foundation Central configuration is a SINGLETON managed by the Foundation
// Central (FCVM) service. There is no create/delete verb — the object is mutated
// in place via an asynchronous Update that returns a task. The only user-writable
// fields are the two timeout knobs:
//   - ahv_installation_timeout_minutes
//   - aos_download_timeout_minutes
// All other fields (commit_id, tls_certificate_fingerprint, version) are computed
// and read back from the service.
//
// Scenarios covered:
//  1. Create/apply with both timeout fields set, assert literal values + that the
//     computed attributes (version) are populated. Read back via the singular
//     datasource and assert every attribute matches the resource state.
//  2. Update both timeout fields to DIFFERENT literal values and re-assert, while
//     confirming the computed fields remain populated.
//  3. Negative: supply a non-integer value for a timeout field and assert the SDK
//     type validation rejects it.
//
// Because the singleton cannot be destroyed, CheckDestroy verifies the config is
// still readable after the resource leaves state (it must NOT disappear).

const resourceNameFCC = "nutanix_foundation_central_config_v2.test"
const datasourceNameFCC = "data.nutanix_foundation_central_config_v2.test"

func TestAccV2NutanixFoundationCentralConfigResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccFoundationPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixFoundationCentralConfigStillReadable,
		Steps: []resource.TestStep{
			// Step 1 - apply initial timeout configuration.
			{
				Config: testAccFoundationCentralConfigConfig(120, 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameFCC, "ahv_installation_timeout_minutes", "120"),
					resource.TestCheckResourceAttr(resourceNameFCC, "aos_download_timeout_minutes", "60"),
					resource.TestCheckResourceAttrSet(resourceNameFCC, "version"),
					// Datasource read-back — every attribute must match the resource.
					resource.TestCheckResourceAttrPair(datasourceNameFCC, "ahv_installation_timeout_minutes", resourceNameFCC, "ahv_installation_timeout_minutes"),
					resource.TestCheckResourceAttrPair(datasourceNameFCC, "aos_download_timeout_minutes", resourceNameFCC, "aos_download_timeout_minutes"),
					resource.TestCheckResourceAttrPair(datasourceNameFCC, "commit_id", resourceNameFCC, "commit_id"),
					resource.TestCheckResourceAttrPair(datasourceNameFCC, "tls_certificate_fingerprint", resourceNameFCC, "tls_certificate_fingerprint"),
					resource.TestCheckResourceAttrPair(datasourceNameFCC, "version", resourceNameFCC, "version"),
				),
			},
			// Step 2 - update BOTH timeout fields to different literal values.
			{
				Config: testAccFoundationCentralConfigConfig(150, 90),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameFCC, "ahv_installation_timeout_minutes", "150"),
					resource.TestCheckResourceAttr(resourceNameFCC, "aos_download_timeout_minutes", "90"),
					resource.TestCheckResourceAttrSet(resourceNameFCC, "version"),
					resource.TestCheckResourceAttrPair(datasourceNameFCC, "ahv_installation_timeout_minutes", resourceNameFCC, "ahv_installation_timeout_minutes"),
					resource.TestCheckResourceAttrPair(datasourceNameFCC, "aos_download_timeout_minutes", resourceNameFCC, "aos_download_timeout_minutes"),
				),
			},
		},
	})
}

// TestAccV2NutanixFoundationCentralConfigResource_InvalidTimeout asserts the SDK
// type system rejects a non-integer timeout value.
func TestAccV2NutanixFoundationCentralConfigResource_InvalidTimeout(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccFoundationCentralConfigInvalid(),
				ExpectError: regexp.MustCompile(`(?i)(a number is required|Incorrect attribute value type|Inappropriate value)`),
			},
		},
	})
}

func testAccFoundationCentralConfigConfig(ahvTimeout, aosTimeout int) string {
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

func testAccFoundationCentralConfigInvalid() string {
	return `
resource "nutanix_foundation_central_config_v2" "test" {
  ahv_installation_timeout_minutes = "not-a-number"
}
`
}

// testAccCheckNutanixFoundationCentralConfigStillReadable verifies the Foundation
// Central configuration singleton continues to be readable after the resource is
// removed from state. Unlike a normal entity, this configuration cannot be deleted,
// so the "destroy" contract is simply that the config remains fetchable.
func testAccCheckNutanixFoundationCentralConfigStillReadable(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).LifecycleAPI
	if _, err := conn.FoundationCentralConfigAPIInstance.GetFoundationCentralConfig(); err != nil {
		return fmt.Errorf("Foundation Central config should remain readable after destroy: %v", err)
	}
	return nil
}
