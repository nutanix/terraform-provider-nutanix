package nutanix

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()

	testAccProviders = map[string]*schema.Provider{
		"nutanix": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("NUTANIX_USERNAME") == "" ||
		os.Getenv("NUTANIX_PASSWORD") == "" ||
		os.Getenv("NUTANIX_INSECURE") == "" ||
		os.Getenv("NUTANIX_PORT") == "" ||
		os.Getenv("NUTANIX_ENDPOINT") == "" ||
		os.Getenv("NUTANIX_STORAGE_CONTAINER") == "" {
		t.Fatal("`NUTANIX_USERNAME`,`NUTANIX_PASSWORD`,`NUTANIX_INSECURE`,`NUTANIX_PORT`,`NUTANIX_ENDPOINT`, `NUTANIX_STORAGE_CONTAINER` must be set for acceptance testing")
	}
}

func testAccFoundationPreCheck(t *testing.T) {
	if os.Getenv("FOUNDATION_ENDPOINT") == "" ||
		os.Getenv("FOUNDATION_PORT") == "" {
		t.Fatal("`FOUNDATION_ENDPOINT` and `FOUNDATION_PORT` must be set for foundation acceptance testing")
	}
}

func testAccEraPreCheck(t *testing.T) {
	if os.Getenv("ERA_ENDPOINT") == "" ||
		os.Getenv("ERA_USERNAME") == "" ||
		os.Getenv("ERA_PASSWORD") == "" {
		t.Fatal("`ERA_USERNAME`,`ERA_PASSWORD`,`ERA_ENDPOINT` must be set for acceptance testing")
	}
}

func randIntBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func isGCPEnvironment() bool {
	return os.Getenv("NUTANIX_GCP") == "true"
}
