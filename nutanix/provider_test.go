package nutanix

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"nutanix": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("NUTANIX_USERNAME") == "" ||
		os.Getenv("NUTANIX_PASSWORD") == "" ||
		os.Getenv("NUTANIX_INSECURE") == "" ||
		os.Getenv("NUTANIX_PORT") == "" ||
		os.Getenv("NUTANIX_ENDPOINT") == "" {
		t.Fatal("`NUTANIX_USERNAME`,`NUTANIX_PASSWORD`,`NUTANIX_INSECURE`,`NUTANIX_PORT`,`NUTANIX_ENDPOINT` must be set for acceptance testing")
	}
}

func randIntBetween(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func isGCPEnvironment() bool {
	return os.Getenv("NUTANIX_GCP") == "true"
}
