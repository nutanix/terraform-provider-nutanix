package objectstoresv2_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/objects/v4/config"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	import1 "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/config"
)

// testVars holds the shared test fixtures loaded from test_config_v2.json.
// filepath is the absolute path to that file, injected into Terraform configs
// that read the fixtures directly via jsondecode(file(...)).
// certificateJSONFile is generated in TestMain from the fixtures.
var (
	testVars            = acc.MustConfig()
	filepath            = acc.ConfigPath()
	certificateJSONFile = func() string {
		wd, _ := os.Getwd()
		return wd + "/../../../object_store_cert.json"
	}()
)

func TestMain(m *testing.M) {
	log.Println("Do some crazy stuff before tests!")
	if err := createCertificateJSONFile(); err != nil {
		log.Printf("warning: failed to create certificate JSON file: %s", err)
	}
	os.Exit(m.Run())
}

func createCertificateJSONFile() error {
	alternateIps := testVars.ObjectStore.PublicNetworkIPs

	certificate := config.NewCertificate()
	if len(alternateIps) > 0 {
		certificate.AlternateIps = make([]import1.IPAddress, 1)
		certificate.AlternateIps[0] = import1.IPAddress{
			Ipv4: &import1.IPv4Address{
				Value: utils.StringPtr(alternateIps[0]),
			},
		}
	}

	certificateJSON, err := json.MarshalIndent(certificate, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal certificate data: %w", err)
	}
	err = os.WriteFile(certificateJSONFile, certificateJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}
	log.Printf("Certificate JSON file created at: %s", certificateJSONFile)
	return nil
}
