package tenancyv2_test

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

func TestMain(m *testing.M) {
	_ = testVars
	_ = filepath
	os.Exit(m.Run())
}
