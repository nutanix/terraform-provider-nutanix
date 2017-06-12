package nutanix

import (
	"github.com/hashicorp/terraform/builtin/providers/template"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccTemplateProvider *schema.Provider
var terraformState string

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccTemplateProvider = template.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"nutanix":  testAccProvider,
		"template": testAccTemplateProvider,
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
	if v := os.Getenv("NUTANIX_USERNAME"); v == "" {
		t.Fatal("NUTANIX_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("NUTANIX_PASSWORD"); v == "" {
		t.Fatal("NUTANIX_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("NUTANIX_ENDPOINT"); v == "" {
		t.Fatal("NUTANIX_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("NUTANIX_INSECURE"); v == "" {
		t.Fatal("NUTANIX_INSECURE must be set for acceptance tests")
	}
	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}
