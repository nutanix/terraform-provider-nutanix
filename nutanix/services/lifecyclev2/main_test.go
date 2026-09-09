package lifecyclev2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

// TestConfig mirrors the shared test_config_v2.json structure. The Foundation
// Central config feature is a cluster-wide singleton and needs no infrastructure
// parameters, so this struct is intentionally minimal; extend it if future
// lifecyclev2 tests require additional infra values.
type TestConfig struct {
	Lifecycle struct {
		FoundationCentral struct {
			AhvInstallationTimeoutMinutes int `json:"ahv_installation_timeout_minutes"`
			AosDownloadTimeoutMinutes     int `json:"aos_download_timeout_minutes"`
		} `json:"foundation_central"`
		// InstallerImage carries the parameters used to register an installer
		// image with Foundation Central during acceptance tests. URL points to a
		// reachable image (e.g. an AOS/AHV/ESX installer) and Version/UpdatedURL
		// are used across the create/update test steps.
		InstallerImage struct {
			Name       string `json:"name"`
			Type       string `json:"type"`
			URL        string `json:"url"`
			Version    string `json:"version"`
			UpdatedURL string `json:"updated_url"`
		} `json:"installer_image"`
	} `json:"lifecycle"`
}

var testVars TestConfig

var (
	path, _  = os.Getwd()
	filepath = path + "/../../../test_config_v2.json"
)

func loadVars(filepath string, varStuct interface{}) {
	// Read config.json from home current path
	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Got this error while reading config.json: %s", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(configData, varStuct)
	if err != nil {
		log.Printf("Got this error while unmarshalling config.json: %s", err.Error())
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	log.Println("Do some crazy stuff before tests!")
	loadVars("../../../test_config_v2.json", &testVars)
	os.Exit(m.Run())
}
