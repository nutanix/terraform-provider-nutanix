// Package securityv2_test provides testing utilities for the securityv2 package.
package securityv2_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameKeyManagementServer = "data.nutanix_key_management_server_v2.test"

func TestAccV2NutanixKeyManagementServerDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-kms-%d", r)
	// Expiry time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)
	expirationTimeFormatted := expirationTime.UTC().Format("2006-01-02")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testKMSResourceConfig(name, expirationTimeFormatted) + testKMSdatasourceFetchConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServer, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "name", name),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.client_id", testVars.Security.KMS.ClientID),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.credential_expiry_date", expirationTimeFormatted),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.endpoint_url", testVars.Security.KMS.EndpointURL),
					func(s *terraform.State) error {
						kmsAttributes := s.RootModule().Resources[datasourceNameKeyManagementServer].Primary.Attributes

						keyID := kmsAttributes["access_information.0.key_id"]

						if strings.Split(keyID, ":")[0] == testVars.Security.KMS.KeyID {
							return nil
						}
						return fmt.Errorf("expected key_id to contain %q, got %q", testVars.Security.KMS.KeyID, keyID)
					},
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.tenant_id", testVars.Security.KMS.TenantID),
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServer, "access_information.0.truncated_client_secret"),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.credential_expiry_date", expirationTimeFormatted),
				),
			},
		},
	})
}

func testKMSdatasourceFetchConfig() string {
	return `
data "nutanix_key_management_server_v2" "test" {
  ext_id = nutanix_key_management_server_v2.test.id
}
`
}
