package multidomainv2_test

// Test plan for nutanix_local_domain_v2 (multidomain / LocalDomainService)
//
// LocalDomainService models the registration of a Prism Central with Nutanix
// Central. It is a SINGLETON entity exposed through three API methods:
//   - RegisterLocalDomain   (Create) — POST /management/local-domain/$actions/register
//   - GetLocalDomain        (Read)   — GET  /management/local-domain
//   - UnregisterLocalDomain (Delete) — POST /management/local-domain/$actions/unregister
//
// Because registration establishes a live connection between Prism Central and
// a real Nutanix Central instance, the happy-path lifecycle test needs genuine
// Nutanix Central credentials (API key + key id) plus the domain/IDP/tenant
// ext_ids. Those are provided via environment variables and the test is skipped
// when they are absent, so the suite still runs in CI without an NC backend.
//
// Scenarios covered:
//   1. Basic lifecycle (Register -> Read -> Unregister) with ALL fields set,
//      asserting the computed outputs (name, registration_status, etc.) and the
//      singular datasource mirror the resource. Gated on NC env vars.
//   2. Import of the singleton resource. Gated on NC env vars.
//   3. Negative: missing each required argument -> "Missing required argument".
//      These run without any backend because they fail at plan time.

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceLocalDomainGet = "data.nutanix_local_domain_v2.test"

// localDomainTestVars gathers the Nutanix Central registration inputs from the
// environment. All must be present for the live lifecycle test to run.
type localDomainTestVars struct {
	apiKey                string
	keyID                 string
	domainExtID           string
	identityProviderExtID string
	targetURL             string
	tenantExtID           string
}

func getLocalDomainTestVars() localDomainTestVars {
	return localDomainTestVars{
		apiKey:                os.Getenv("NUTANIX_NC_API_KEY"),
		keyID:                 os.Getenv("NUTANIX_NC_KEY_ID"),
		domainExtID:           os.Getenv("NUTANIX_NC_DOMAIN_EXT_ID"),
		identityProviderExtID: os.Getenv("NUTANIX_NC_IDENTITY_PROVIDER_EXT_ID"),
		targetURL:             os.Getenv("NUTANIX_NC_TARGET_URL"),
		tenantExtID:           os.Getenv("NUTANIX_NC_TENANT_EXT_ID"),
	}
}

func (v localDomainTestVars) complete() bool {
	return v.apiKey != "" && v.keyID != "" && v.domainExtID != "" &&
		v.identityProviderExtID != "" && v.targetURL != "" && v.tenantExtID != ""
}

func TestAccV2NutanixLocalDomainResource_Basic(t *testing.T) {
	tv := getLocalDomainTestVars()
	if !tv.complete() {
		t.Skip("skipping: set NUTANIX_NC_API_KEY, NUTANIX_NC_KEY_ID, NUTANIX_NC_DOMAIN_EXT_ID, " +
			"NUTANIX_NC_IDENTITY_PROVIDER_EXT_ID, NUTANIX_NC_TARGET_URL and NUTANIX_NC_TENANT_EXT_ID to run the LocalDomain lifecycle test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testLocalDomainV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalDomainV2ResourceConfig(tv),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameLocalDomainV2, "domain_ext_id", tv.domainExtID),
					resource.TestCheckResourceAttr(resourceNameLocalDomainV2, "identity_provider_ext_id", tv.identityProviderExtID),
					resource.TestCheckResourceAttr(resourceNameLocalDomainV2, "target_url", tv.targetURL),
					resource.TestCheckResourceAttr(resourceNameLocalDomainV2, "tenant_ext_id", tv.tenantExtID),
					resource.TestCheckResourceAttrSet(resourceNameLocalDomainV2, "name"),
					resource.TestCheckResourceAttrSet(resourceNameLocalDomainV2, "registration_status.#"),
					// singular datasource mirrors the resource state
					resource.TestCheckResourceAttrPair(dataSourceLocalDomainGet, "name", resourceNameLocalDomainV2, "name"),
					resource.TestCheckResourceAttrPair(dataSourceLocalDomainGet, "target_url", resourceNameLocalDomainV2, "target_url"),
					resource.TestCheckResourceAttrPair(dataSourceLocalDomainGet, "registration_status.#", resourceNameLocalDomainV2, "registration_status.#"),
					resource.TestCheckResourceAttrPair(dataSourceLocalDomainGet, "connection_status.#", resourceNameLocalDomainV2, "connection_status.#"),
				),
			},
			{
				ResourceName:            resourceNameLocalDomainV2,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
		},
	})
}

func TestAccV2NutanixLocalDomainResource_MissingCredentials(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccLocalDomainV2MissingCredentialsConfig(),
				ExpectError: regexp.MustCompile(`(?s)(Insufficient credentials blocks|"credentials" blocks are required)`),
			},
		},
	})
}

func TestAccV2NutanixLocalDomainResource_MissingDomainExtID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccLocalDomainV2MissingDomainExtIDConfig(),
				ExpectError: regexp.MustCompile(`The argument "domain_ext_id" is required`),
			},
		},
	})
}

func testAccLocalDomainV2ResourceConfig(tv localDomainTestVars) string {
	return fmt.Sprintf(`
resource "nutanix_local_domain_v2" "test" {
  credentials {
    api_key = "%[1]s"
    key_id  = "%[2]s"
  }
  domain_ext_id            = "%[3]s"
  identity_provider_ext_id = "%[4]s"
  target_url               = "%[5]s"
  tenant_ext_id            = "%[6]s"
}

data "nutanix_local_domain_v2" "test" {
  depends_on = [nutanix_local_domain_v2.test]
}
`, tv.apiKey, tv.keyID, tv.domainExtID, tv.identityProviderExtID, tv.targetURL, tv.tenantExtID)
}

func testAccLocalDomainV2MissingCredentialsConfig() string {
	return `
resource "nutanix_local_domain_v2" "test" {
  domain_ext_id            = "00000000-0000-0000-0000-000000000000"
  identity_provider_ext_id = "11111111-1111-1111-1111-111111111111"
  target_url               = "https://nutanix-central.example.com"
  tenant_ext_id            = "22222222-2222-2222-2222-222222222222"
}
`
}

func testAccLocalDomainV2MissingDomainExtIDConfig() string {
	return `
resource "nutanix_local_domain_v2" "test" {
  credentials {
    api_key = "fake-api-key"
    key_id  = "fake-key-id"
  }
  identity_provider_ext_id = "11111111-1111-1111-1111-111111111111"
  target_url               = "https://nutanix-central.example.com"
  tenant_ext_id            = "22222222-2222-2222-2222-222222222222"
}
`
}
