package acctest

import (
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/provider"
)

var TestAccProviders map[string]*schema.Provider

var TestAccProvider *schema.Provider
var TestAccProvider2 *schema.Provider

// TestAccProtoV5ProviderFactories is an opt-in alternative to TestAccProviders
// for TestCases that declare their own provider block(s) in the Terraform
// config (e.g. protection-policy tests that force basic auth, or remote-PC
// tests configuring "nutanix-2"). The legacy Providers map makes the
// terraform-plugin-sdk/v2 test harness inject an empty `provider "nutanix" {}`
// (and `provider "nutanix-2" {}`) block unconditionally, which collides with a
// test-supplied provider block ("Duplicate provider configuration"). Using
// provider factories makes the harness skip that injection, so a test's own
// provider block is the only one present.
//
// The factories wrap the same TestAccProvider/TestAccProvider2 instances used
// by TestAccProviders, so CheckDestroy helpers that read TestAccProvider.Meta()
// continue to work unchanged.
var TestAccProtoV5ProviderFactories map[string]func() (tfprotov5.ProviderServer, error)

func init() {
	TestAccProvider = provider.Provider()
	TestAccProvider2 = provider.Provider()

	TestAccProviders = map[string]*schema.Provider{
		"nutanix":   TestAccProvider,
		"nutanix-2": TestAccProvider2,
	}

	TestAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"nutanix": func() (tfprotov5.ProviderServer, error) {
			return schema.NewGRPCProviderServer(TestAccProvider), nil
		},
		"nutanix-2": func() (tfprotov5.ProviderServer, error) {
			return schema.NewGRPCProviderServer(TestAccProvider2), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ *schema.Provider = provider.Provider()
}

// TestAccPreCheck validates the environment required to run acceptance tests.
//
// # Endpoint configuration
//
// The provider can target either a Prism Central directly or a Nutanix Central
// (NC) deployment. Which "mode" is active is decided by the base endpoint that
// is set, so at least one of NUTANIX_ENDPOINT or NUTANIX_CENTRAL_ENDPOINT must
// be present. The three supported scenarios are:
//
//  1. PC only (traditional Prism Central tests)
//     - Set:   NUTANIX_ENDPOINT  (Prism Central IP/FQDN)
//     - Unset: NUTANIX_CENTRAL_ENDPOINT
//     - Every API call goes directly to that Prism Central.
//
//  2. PC tests routed through Nutanix Central
//     - Set:   NUTANIX_CENTRAL_ENDPOINT  (NC base endpoint)
//     NUTANIX_DOMAIN_NAME       (registered PC domain name)
//     - Unset: NUTANIX_ENDPOINT
//     - "Other PC" calls (clusters, storage containers, VPCs, subnets, ...)
//     are routed to <NUTANIX_DOMAIN_NAME>.domains.<NUTANIX_CENTRAL_ENDPOINT>,
//     so NUTANIX_DOMAIN_NAME is mandatory to exercise PC resources via NC.
//
//  3. Service Provider (SP) tenancy tests
//     - Set:   NUTANIX_CENTRAL_ENDPOINT  (NC base endpoint)
//     NUTANIX_USERNAME + NUTANIX_PASSWORD  (SP admin credentials)
//     - Unset: NUTANIX_ENDPOINT
//     - Service-Provider calls (tenants, domain shortlisting, domain
//     allocations) are routed to sp.services.<NUTANIX_CENTRAL_ENDPOINT>.
//     NUTANIX_DOMAIN_NAME is only needed if the test also touches PC
//     resources (as in scenario 2).
//
// Note: NUTANIX_ENDPOINT and NUTANIX_CENTRAL_ENDPOINT are mutually exclusive for
// a single provider mode (NC mode wins when both happen to be set). This check
// therefore only enforces that at least one is present; some NC-workflow
// pre-checks (e.g. TestAccNCPreCheck) intentionally set both to stand up an
// additional direct-PC provider.
//
// Always required (both modes): NUTANIX_INSECURE, NUTANIX_PORT, and credentials
// (NUTANIX_USERNAME + NUTANIX_PASSWORD, or NUTANIX_API_KEY).
func TestAccPreCheck(t *testing.T) {
	// NUTANIX_INSECURE and NUTANIX_PORT are required in every mode.
	if os.Getenv("NUTANIX_INSECURE") == "" || os.Getenv("NUTANIX_PORT") == "" {
		t.Fatal("`NUTANIX_INSECURE` and `NUTANIX_PORT` must be set for acceptance testing")
	}

	// Exactly one base endpoint drives the provider mode: NUTANIX_ENDPOINT for
	// direct Prism Central access, or NUTANIX_CENTRAL_ENDPOINT for Nutanix
	// Central. At least one must be set.
	if os.Getenv("NUTANIX_ENDPOINT") == "" && os.Getenv("NUTANIX_CENTRAL_ENDPOINT") == "" {
		t.Fatal("one of `NUTANIX_ENDPOINT` (direct Prism Central) or `NUTANIX_CENTRAL_ENDPOINT` (Nutanix Central) must be set for acceptance testing")
	}

	// Check authentication - either username/password OR api_key must be set
	hasBasicAuth := os.Getenv("NUTANIX_USERNAME") != "" && os.Getenv("NUTANIX_PASSWORD") != ""
	hasAPIKey := os.Getenv("NUTANIX_API_KEY") != ""

	if !hasBasicAuth && !hasAPIKey {
		t.Fatal("Either `NUTANIX_USERNAME` and `NUTANIX_PASSWORD`, or `NUTANIX_API_KEY` must be set for acceptance testing")
	}
}

// TestAccNCPreCheck verifies the environment required by acceptance tests that
// exercise an NC-managed domain workflow.
func TestAccNCPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("NUTANIX_CENTRAL_ENDPOINT") == "" || os.Getenv("NUTANIX_DOMAIN_NAME") == "" {
		t.Skip("skipping NC workflow test: set NUTANIX_CENTRAL_ENDPOINT and NUTANIX_DOMAIN_NAME")
	}
	if os.Getenv("NUTANIX_ENDPOINT") == "" {
		t.Skip("skipping NC workflow test: set NUTANIX_ENDPOINT for the direct PC local-domain provider")
	}
	TestAccPreCheck(t)
}

func TestAccNCPreCheckPEWorkflow(t *testing.T) {
	t.Helper()
	if os.Getenv("NUTANIX_CENTRAL_ENDPOINT") != "" || os.Getenv("NUTANIX_DOMAIN_NAME") != "" {
		t.Skip("PC skipping PE workflow tests as it is not supported in NC mode")
	}
	if os.Getenv("NUTANIX_ENDPOINT") == "" {
		t.Skip("NUTANIX_ENDPOINT should be set for PE workflow tests")
	}
	TestAccPreCheck(t)
}

// TestAccTenantAdminPreCheck verifies the environment required by the Tenant
// Admin acceptance tests that route entirely through Nutanix Central (NC).
//
// These tests are mutually exclusive with the regular PC acceptance tests: they
// run ONLY when NUTANIX_TENANT_NAME is set. A standard PC test run does not set
// NUTANIX_TENANT_NAME, so every tenant admin case is skipped during PC runs (and
// vice-versa). This is the mechanism that keeps the two suites from executing
// together. It skips (never fatals) so a mixed run is reported as SKIP.
func TestAccTenantAdminPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("NUTANIX_TENANT_NAME") == "" {
		t.Skip("skipping tenant admin test: NUTANIX_TENANT_NAME is not set (PC acceptance run)")
	}
	if os.Getenv("NUTANIX_CENTRAL_ENDPOINT") == "" || os.Getenv("NUTANIX_DOMAIN_NAME") == "" {
		t.Skip("skipping tenant admin test: set NUTANIX_CENTRAL_ENDPOINT and NUTANIX_DOMAIN_NAME")
	}
	// Tenant admin requires basic auth OR api_key routed through NC.
	hasBasicAuth := os.Getenv("NUTANIX_USERNAME") != "" && os.Getenv("NUTANIX_PASSWORD") != ""
	hasAPIKey := os.Getenv("NUTANIX_API_KEY") != ""
	if !hasBasicAuth && !hasAPIKey {
		t.Skip("skipping tenant admin test: set NUTANIX_USERNAME+NUTANIX_PASSWORD or NUTANIX_API_KEY")
	}
}

// TestAccPreCheckStorageContainer checks for storage container requirement
// Use this in addition to TestAccPreCheck for tests that create VMs with disks
func TestAccPreCheckStorageContainer(t *testing.T) {
	TestAccPreCheck(t)
	if os.Getenv("NUTANIX_STORAGE_CONTAINER") == "" {
		t.Fatal("`NUTANIX_STORAGE_CONTAINER` must be set for VM creation tests")
	}
}

func TestAccFoundationPreCheck(t *testing.T) {
	if os.Getenv("FOUNDATION_ENDPOINT") == "" ||
		os.Getenv("FOUNDATION_PORT") == "" {
		t.Fatal("`FOUNDATION_ENDPOINT` and `FOUNDATION_PORT` must be set for foundation acceptance testing")
	}
}

func TestAccEraPreCheck(t *testing.T) {
	if os.Getenv("NDB_ENDPOINT") == "" ||
		os.Getenv("NDB_USERNAME") == "" ||
		os.Getenv("NDB_PASSWORD") == "" {
		t.Fatal("`NDB_USERNAME`,`NDB_PASSWORD`,`NDB_ENDPOINT` must be set for acceptance testing")
	}
}

// v3FlowNextGenSkipMarkers are substrings of the PC errors returned when the
// legacy v3 microsegmentation APIs (address groups, service groups, network
// security rules) are removed on the cluster. This happens either because Flow
// network security next-gen is enabled, or because the v3 API itself has been
// deprecated in favour of the v4 network security policy APIs. On such clusters
// those v3-only resources return HTTP 410 GONE and the tests exercising them
// cannot pass, so we skip rather than fail them.
var v3FlowNextGenSkipMarkers = []string{
	"Flow network security next-gen is enabled",
	"The network security rule APIs are no longer supported",
}

// SkipIfV3FlowNextGen returns a resource.TestCase ErrorCheck that converts the
// "v3 API gone / Flow next-gen enabled" apply error into a skipped test, so it
// is reported as SKIP instead of FAIL in the test summary. Any other error is
// returned unchanged so genuine failures still fail.
func SkipIfV3FlowNextGen(t *testing.T) resource.ErrorCheckFunc {
	return func(err error) error {
		if err != nil {
			for _, marker := range v3FlowNextGenSkipMarkers {
				if strings.Contains(err.Error(), marker) {
					t.Skipf("skipping: legacy v3 microsegmentation API unavailable on this PC (Flow next-gen enabled or v3 API deprecated): %v", err)
				}
			}
		}
		return err
	}
}

// PCEndpointHost returns the hostname used for Prism Central API calls in
// acceptance-test helpers (e.g. local-exec curl polls). Matches provider
// routing in sdkconfig.ResolveEndpoint for PC service calls:
//
//   - NC mode (NUTANIX_CENTRAL_ENDPOINT set): <NUTANIX_DOMAIN_NAME>.domains.<NUTANIX_CENTRAL_ENDPOINT>
//   - Direct PC (traditional): NUTANIX_ENDPOINT
//
// NC wins when both endpoints are set, same as the provider.
func PCEndpointHost() string {
	if central := os.Getenv("NUTANIX_CENTRAL_ENDPOINT"); central != "" {
		if domain := os.Getenv("NUTANIX_DOMAIN_NAME"); domain != "" {
			return domain + ".domains." + central
		}
		return central
	}
	return os.Getenv("NUTANIX_ENDPOINT")
}

// PCEndpointPort returns NUTANIX_PORT, defaulting to 9440 when unset.
func PCEndpointPort() string {
	if p := os.Getenv("NUTANIX_PORT"); p != "" {
		return p
	}
	return "9440"
}

func RandIntBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func isGCPEnvironment() bool {
	return os.Getenv("NUTANIX_GCP") == "true"
}
