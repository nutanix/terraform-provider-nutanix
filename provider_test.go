package nutanix

import (
	"flag"
	"github.com/hashicorp/terraform/builtin/providers/template"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccTemplateProvider *schema.Provider
var terraformState string
var NutanixUsername string
var NutanixPassword string
var NutanixEndpoint string
var NutanixInsecure bool
var NutanixPort string

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
	NutanixUsername = *flag.String("username", "", "username for api call")
	NutanixPassword = *flag.String("password", "", "password for api call")
	NutanixEndpoint = *flag.String("endpoint", "", "endpoint must be set")
	NutanixInsecure = *flag.Bool("insecure", false, "insecure flag")
	NutanixPort = *flag.String("port", "9440", "port for api call")
}

func testAccPreCheck(t *testing.T) {
	if NutanixUsername == "" {
		t.Fatal("username flag must be set for acceptance tests")
	}
	if NutanixPassword == "" {
		t.Fatal("password must be set for acceptance tests")
	}
	if NutanixEndpoint == "" {
		t.Fatal("endpoint flag must be set for acceptance tests")
	}
	if NutanixInsecure == false {
		t.Fatal("insecure flag must be set true for acceptance tests")
	}
	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}
