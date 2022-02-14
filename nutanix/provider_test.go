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

func randIntBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func isGCPEnvironment() bool {
	return os.Getenv("NUTANIX_GCP") == "true"
}
