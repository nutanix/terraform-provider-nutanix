package iamv2_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// testVars holds the shared test fixtures loaded from test_config_v2.json.
// filepath is the absolute path to that file, injected into Terraform configs
// that read the fixtures directly via jsondecode(file(...)).
// xmlFilePath points to the IdP metadata file downloaded in TestMain.
var (
	testVars    = acc.MustConfig()
	filepath    = acc.ConfigPath()
	xmlFilePath = func() string {
		wd, _ := os.Getwd()
		return wd + "/../../../test_idp_metadata.txt"
	}()
)

func TestMain(m *testing.M) {
	log.Println("Do some crazy stuff before tests!")
	// Best-effort download of the SAML IdP metadata used by identity-provider
	// tests; errors are non-fatal (matches the previous behavior).
	_ = downloadFile(testVars.Iam.IdentityProviders.IdpMetadataXML, xmlFilePath)
	os.Exit(m.Run())
}

// downloadFile downloads a file from a given URL and saves it to the specified path.
func downloadFile(url, destinationFilePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}
