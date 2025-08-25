// Package securityv2_test provides testing utilities for the securityv2 package.
package securityv2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameKeyManagementServer = "data.nutanix_key_management_server_v2.test"

func TestAccV2NutanixKeyManagementServerDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-kms-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testKMSResourceConfig(name) + testKMSdatasourceFetchConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServer, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "name", name),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.client_id", testVars.Security.ClientID),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.credential_expiry_date", testVars.Security.CredentialExpiryDate),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.endpoint_url", testVars.Security.EndpointURL),
					func(s *terraform.State) error {
						kmsAttributes := s.RootModule().Resources[datasourceNameKeyManagementServer].Primary.Attributes

						keyID := kmsAttributes["access_information.0.key_id"]

						if strings.Split(keyID, ":")[0] == testVars.Security.KeyID {
							return nil

						}
						return fmt.Errorf("expected key_id to contain %q, got %q", testVars.Security.KeyID, keyID)
					},
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServer, "access_information.0.tenant_id", testVars.Security.TenantID),
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServer, "access_information.0.truncated_client_secret"),
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
