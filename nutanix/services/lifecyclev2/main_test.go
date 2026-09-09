package lifecyclev2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

// TestConfig maps the lifecycle/Nodes section of test_config_v2.json used by the
// acceptance tests. Node tests operate on hardware that Foundation Central has
// already discovered, so infrastructure identifiers (node ext_id, serial,
// manufacturer, provider references, image ext_ids) are supplied here instead of
// being hardcoded.
type TestConfig struct {
	Lifecycle struct {
		Node struct {
			ExtID                   string `json:"ext_id"`
			NodeSerial              string `json:"node_serial"`
			Manufacturer            string `json:"manufacturer"`
			ProviderExtID           string `json:"provider_ext_id"`
			ProviderConnectionExtID string `json:"provider_connection_ext_id"`
			PatchedImageExtID       string `json:"patched_image_ext_id"`
			AosImageExtID           string `json:"aos_image_ext_id"`
		} `json:"node"`
	} `json:"lifecycle"`
}

var testVars TestConfig

var (
	path, _  = os.Getwd()
	filepath = path + "/../../../test_config_v2.json"
)

func loadVars(filepath string, varStruct interface{}) {
	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Got this error while reading test_config_v2.json: %s", err.Error())
		os.Exit(1)
	}

	if err := json.Unmarshal(configData, varStruct); err != nil {
		log.Printf("Got this error while unmarshalling test_config_v2.json: %s", err.Error())
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	log.Println("Setting up lifecycle/Nodes acceptance test configuration")
	loadVars(filepath, &testVars)
	os.Exit(m.Run())
}
