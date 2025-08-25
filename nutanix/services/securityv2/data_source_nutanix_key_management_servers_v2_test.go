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

const datasourceNameKeyManagementServers = "data.nutanix_key_management_servers_v2.test"

func TestAccV2NutanixKeyManagementServersDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-kms-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testKMSResourceConfig(name) + testKMSdatasourceListConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServers, "kms.#"),
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServers, "kms.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.client_id", testVars.Security.ClientID),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.credential_expiry_date", testVars.Security.CredentialExpiryDate),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.endpoint_url", testVars.Security.EndpointURL),
					func(s *terraform.State) error {
						kmsAttributes := s.RootModule().Resources[datasourceNameKeyManagementServers].Primary.Attributes

						keyID := kmsAttributes["kms.0.access_information.0.key_id"]

						if strings.Split(keyID, ":")[0] == testVars.Security.KeyID {
							return nil

						}
						return fmt.Errorf("expected key_id to contain %q, got %q", testVars.Security.KeyID, keyID)
					},
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.tenant_id", testVars.Security.TenantID),
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServers, "kms.0.access_information.0.truncated_client_secret"),
				),
			},
		},
	})
}

func testKMSdatasourceListConfig() string {
	return `
data "nutanix_key_management_servers_v2" "test" {
  depends_on = [nutanix_key_management_server_v2.test]
}
`
}
