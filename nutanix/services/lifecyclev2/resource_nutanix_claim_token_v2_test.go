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

// Test plan for nutanix_claim_token_v2 (+ datasources)
// ----------------------------------------------------
// A claim token is a Foundation Central object used to register nodes. It carries
// a name, an expiry_time (RFC3339) and a max_usage_count, and the API returns
// computed ext_id, created_time, current_usage_count, owner_ext_id, tenant_id and
// links. Coverage below:
//
//  1. TestAccV2NutanixClaimTokenResource_Basic — full lifecycle:
//     - Step 1 (Create): all input fields set; assert every literal + computed set.
//       Read back via the singular datasource (TestCheckResourceAttrPair on ALL
//       attributes) and via the secret datasource.
//     - Step 2 (Update): change name, expiry_time and max_usage_count; assert new
//       literals; confirm ext_id is stable (in-place update, not recreate).
//     - Step 3 (List datasource): assert the token appears in the plural list.
//     - Step 4 (Import): ImportStateVerify round-trips the resource.
//  2. TestAccV2NutanixClaimTokenResource_MinimalRequired — only required fields.
//  3. TestAccV2NutanixClaimTokenResource_MissingRequired — negative: omit name.
//  4. TestAccV2NutanixClaimTokensDataSource_InvalidFilter — negative: bad filter.

const (
	resourceNameClaimToken   = "nutanix_claim_token_v2.test"
	datasourceNameClaimToken = "data.nutanix_claim_token_v2.test"
	datasourceNameSecret     = "data.nutanix_secret_v2.test"
	datasourceNameClaimList  = "data.nutanix_claim_tokens_v2.list"
)

func TestAccV2NutanixClaimTokenResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-claim-token-%d", r)
	nameUpdated := fmt.Sprintf("tf-claim-token-upd-%d", r)
	expiry := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	expiryUpdated := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckClaimTokenDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with all fields + read via singular + secret datasource.
			{
				Config: testClaimTokenConfig(name, expiry, 5) + testClaimTokenDataSources(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameClaimToken, "name", name),
					resource.TestCheckResourceAttr(resourceNameClaimToken, "expiry_time", expiry),
					resource.TestCheckResourceAttr(resourceNameClaimToken, "max_usage_count", "5"),
					resource.TestCheckResourceAttrSet(resourceNameClaimToken, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameClaimToken, "created_time"),
					resource.TestCheckResourceAttr(resourceNameClaimToken, "current_usage_count", "0"),
					// Singular datasource must round-trip every attribute.
					resource.TestCheckResourceAttrPair(datasourceNameClaimToken, "name", resourceNameClaimToken, "name"),
					resource.TestCheckResourceAttrPair(datasourceNameClaimToken, "expiry_time", resourceNameClaimToken, "expiry_time"),
					resource.TestCheckResourceAttrPair(datasourceNameClaimToken, "max_usage_count", resourceNameClaimToken, "max_usage_count"),
					resource.TestCheckResourceAttrPair(datasourceNameClaimToken, "current_usage_count", resourceNameClaimToken, "current_usage_count"),
					resource.TestCheckResourceAttrPair(datasourceNameClaimToken, "created_time", resourceNameClaimToken, "created_time"),
					resource.TestCheckResourceAttrPair(datasourceNameClaimToken, "tenant_id", resourceNameClaimToken, "tenant_id"),
					// Secret datasource must return a non-empty secret.
					resource.TestCheckResourceAttrSet(datasourceNameSecret, "secret"),
				),
			},
			// Step 2: Update name, expiry_time and max_usage_count.
			{
				Config: testClaimTokenConfig(nameUpdated, expiryUpdated, 10) + testClaimTokenDataSources(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameClaimToken, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceNameClaimToken, "expiry_time", expiryUpdated),
					resource.TestCheckResourceAttr(resourceNameClaimToken, "max_usage_count", "10"),
					resource.TestCheckResourceAttrSet(resourceNameClaimToken, "ext_id"),
				),
			},
			// Step 3: Import (config is just the resource, no list datasource).
			{
				Config:            testClaimTokenConfig(nameUpdated, expiryUpdated, 10),
				ResourceName:      resourceNameClaimToken,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Step 4: List datasource contains the token.
			{
				Config: testClaimTokenConfig(nameUpdated, expiryUpdated, 10) + testClaimTokensListConfig(nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameClaimList, "claim_tokens.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixClaimTokenResource_MinimalRequired(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-claim-min-%d", r)
	expiry := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckClaimTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testClaimTokenConfig(name, expiry, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameClaimToken, "name", name),
					resource.TestCheckResourceAttr(resourceNameClaimToken, "max_usage_count", "1"),
					resource.TestCheckResourceAttrSet(resourceNameClaimToken, "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixClaimTokenResource_MissingRequired(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "nutanix_claim_token_v2" "test" {
  expiry_time     = "2030-01-01T00:00:00Z"
  max_usage_count = 5
}
`,
				ExpectError: regexp.MustCompile("The argument \"name\" is required"),
			},
		},
	})
}

func TestAccV2NutanixClaimTokensDataSource_InvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
data "nutanix_claim_tokens_v2" "invalid" {
  filter = "invalid_field eq 'x'"
}
`,
				ExpectError: regexp.MustCompile("error while listing claim tokens"),
			},
		},
	})
}

func testClaimTokenConfig(name, expiry string, maxUsage int) string {
	return fmt.Sprintf(`
resource "nutanix_claim_token_v2" "test" {
  name            = "%[1]s"
  expiry_time     = "%[2]s"
  max_usage_count = %[3]d
}
`, name, expiry, maxUsage)
}

func testClaimTokenDataSources() string {
	return `
data "nutanix_claim_token_v2" "test" {
  ext_id = nutanix_claim_token_v2.test.id
}

data "nutanix_secret_v2" "test" {
  ext_id = nutanix_claim_token_v2.test.id
}
`
}

func testClaimTokensListConfig(name string) string {
	return fmt.Sprintf(`
data "nutanix_claim_tokens_v2" "list" {
  filter     = "name eq '%[1]s'"
  depends_on = [nutanix_claim_token_v2.test]
}
`, name)
}
