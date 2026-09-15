package datapoliciesv2_test

import (
	"os"
	"testing"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// testVars holds the shared test fixtures loaded from test_config_v2.json.
// filepath is the absolute path to that file, injected into Terraform configs
// that read the fixtures directly via jsondecode(file(...)).
var (
	testVars = acc.MustConfig()
	filepath = acc.ConfigPath()
)

// TestMain forces basic-auth for this package's acceptance tests.
//
// The data-policies APIs (protection policy) do not support API-key
// authentication. Rather than declaring a basic-auth `provider "nutanix"` block
// in every test config, we swap the API key for the basic-auth credentials from
// test_config_v2.json here, before the provider is configured. The tests use
// ProtoV5ProviderFactories, so the default "nutanix" provider is configured from
// these environment variables.
//
// In a basic-auth environment (NUTANIX_API_KEY unset) this is a no-op.
func TestMain(m *testing.M) {
	if os.Getenv("NUTANIX_API_KEY") != "" {
		if testVars.UsernameForTest != "" && testVars.PasswordForTest != "" {
			os.Setenv("NUTANIX_USERNAME", testVars.UsernameForTest)
			os.Setenv("NUTANIX_PASSWORD", testVars.PasswordForTest)
		}
		os.Unsetenv("NUTANIX_API_KEY")
	}
	os.Exit(m.Run())
}
