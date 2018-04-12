package nutanix

import (
	"github.com/hashicorp/terraform/builtin/providers/template"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"strconv"
	flag "terraform-provider-nutanix/testflg"
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

func init() {
	os.Setenv("TF_ACC", "1")
}

func testAccPreCheck(t *testing.T) {
	os.Setenv("NUTANIX_USERNAME", flag.NutanixUsername)
	os.Setenv("NUTANIX_PASSWORD", flag.NutanixPassword)
	os.Setenv("NUTANIX_ENDPOINT", flag.NutanixEndpoint)
	os.Setenv("NUTANIX_PORT", flag.NutanixPort)
	os.Setenv("NUTANIX_INSECURE", strconv.FormatBool(flag.NutanixInsecure))
	if flag.NutanixUsername == "" {
		t.Fatal("username flag must be set for acceptance tests")
	}
	if flag.NutanixPassword == "" {
		t.Fatal("password must be set for acceptance tests")
	}
	if flag.NutanixEndpoint == "" {
		t.Fatal("endpoint flag must be set for acceptance tests")
	}
	if flag.NutanixInsecure == false {
		t.Fatal("insecure flag must be set true for acceptance tests")
	}
	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}
