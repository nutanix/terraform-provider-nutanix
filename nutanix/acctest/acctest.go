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

func TestAccPreCheck(t *testing.T) {
	// Check common required variables
	if os.Getenv("NUTANIX_INSECURE") == "" ||
		os.Getenv("NUTANIX_PORT") == "" ||
		os.Getenv("NUTANIX_ENDPOINT") == "" ||
		os.Getenv("NUTANIX_STORAGE_CONTAINER") == "" {
		t.Fatal("`NUTANIX_INSECURE`,`NUTANIX_PORT`,`NUTANIX_ENDPOINT`,`NUTANIX_STORAGE_CONTAINER` must be set for acceptance testing")
	}

	// Check authentication - either username/password OR api_key must be set
	hasBasicAuth := os.Getenv("NUTANIX_USERNAME") != "" && os.Getenv("NUTANIX_PASSWORD") != ""
	hasAPIKey := os.Getenv("NUTANIX_API_KEY") != ""

	if !hasBasicAuth && !hasAPIKey {
		t.Fatal("Either `NUTANIX_USERNAME` and `NUTANIX_PASSWORD`, or `NUTANIX_API_KEY` must be set for acceptance testing")
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

func RandIntBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func isGCPEnvironment() bool {
	return os.Getenv("NUTANIX_GCP") == "true"
}
