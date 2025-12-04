package acctest

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/provider"
)

var TestAccProviders map[string]*schema.Provider

var TestAccProvider *schema.Provider
var TestAccProvider2 *schema.Provider

func init() {
	TestAccProvider = provider.Provider()
	TestAccProvider2 = provider.Provider()

	TestAccProviders = map[string]*schema.Provider{
		"nutanix":   TestAccProvider,
		"nutanix-2": TestAccProvider2,
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
	if os.Getenv("NUTANIX_USERNAME") == "" ||
		os.Getenv("NUTANIX_PASSWORD") == "" ||
		os.Getenv("NUTANIX_INSECURE") == "" ||
		os.Getenv("NUTANIX_PORT") == "" ||
		os.Getenv("NUTANIX_ENDPOINT") == "" ||
		os.Getenv("NUTANIX_STORAGE_CONTAINER") == "" {
		t.Fatal("`NUTANIX_USERNAME`,`NUTANIX_PASSWORD`,`NUTANIX_INSECURE`,`NUTANIX_PORT`,`NUTANIX_ENDPOINT`, `NUTANIX_STORAGE_CONTAINER` must be set for acceptance testing")
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

func RandIntBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func isGCPEnvironment() bool {
	return os.Getenv("NUTANIX_GCP") == "true"
}
