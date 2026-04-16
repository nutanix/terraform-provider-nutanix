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

const datasourceNameKeyManagementServers = "data.nutanix_key_management_servers_v2.test"

func TestAccV2NutanixKeyManagementServersDatasource_Basic(t *testing.T) {
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
				Config: testKMSResourceConfig(name, expirationTimeFormatted) + testKMSdatasourceListConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServers, "kms.#"),
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServers, "kms.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.azure_key_vault.0.client_id", testVars.Security.KMS.ClientID),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.azure_key_vault.0.credential_expiry_date", expirationTimeFormatted),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.azure_key_vault.0.endpoint_url", testVars.Security.KMS.EndpointURL),
					func(s *terraform.State) error {
						kmsAttributes := s.RootModule().Resources[datasourceNameKeyManagementServers].Primary.Attributes

						keyID := kmsAttributes["kms.0.access_information.0.azure_key_vault.0.key_id"]

						if strings.Split(keyID, ":")[0] == testVars.Security.KMS.KeyID {
							return nil
						}
						return fmt.Errorf("expected key_id to contain %q, got %q", testVars.Security.KMS.KeyID, keyID)
					},
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.azure_key_vault.0.tenant_id", testVars.Security.KMS.TenantID),
					resource.TestCheckResourceAttrSet(datasourceNameKeyManagementServers, "kms.0.access_information.0.azure_key_vault.0.truncated_client_secret"),
					resource.TestCheckResourceAttr(datasourceNameKeyManagementServers, "kms.0.access_information.0.azure_key_vault.0.credential_expiry_date", expirationTimeFormatted),
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
