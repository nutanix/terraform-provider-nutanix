// Package securityv2_test provides testing utilities for the securityv2 package.
package securityv2_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameKeyManagementServer = "nutanix_key_management_server_v2.test"

func TestAccV2NutanixKeyManagementServerResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-kms-%d", r)
	updatedName := fmt.Sprintf("%s-updated", name)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testKMSResourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameKeyManagementServer, "id"),
					resource.TestCheckResourceAttrSet(resourceNameKeyManagementServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "name", name),
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "access_information.0.client_id", testVars.Security.ClientID),
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "access_information.0.credential_expiry_date", testVars.Security.CredentialExpiryDate),
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "access_information.0.endpoint_url", testVars.Security.EndpointURL),
					func(s *terraform.State) error {
						kmsAttributes := s.RootModule().Resources[resourceNameKeyManagementServer].Primary.Attributes

						keyID := kmsAttributes["access_information.0.key_id"]

						if strings.Split(keyID, ":")[0] == testVars.Security.KeyID {
							return nil
						}
						return fmt.Errorf("expected key_id to contain %q, got %q", testVars.Security.KeyID, keyID)
					},
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "access_information.0.tenant_id", testVars.Security.TenantID),
					resource.TestCheckResourceAttrSet(resourceNameKeyManagementServer, "access_information.0.truncated_client_secret"),
				),
			},
			// test update
			{
				Config: testKMSResourceConfig(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameKeyManagementServer, "id"),
					resource.TestCheckResourceAttrSet(resourceNameKeyManagementServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "access_information.0.client_id", testVars.Security.ClientID),
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "access_information.0.credential_expiry_date", testVars.Security.CredentialExpiryDate),
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "access_information.0.endpoint_url", testVars.Security.EndpointURL),
					func(s *terraform.State) error {
						kmsAttributes := s.RootModule().Resources[resourceNameKeyManagementServer].Primary.Attributes

						keyID := kmsAttributes["access_information.0.key_id"]

						if strings.Split(keyID, ":")[0] == testVars.Security.KeyID {
							return nil

						}
						return fmt.Errorf("expected key_id to contain %q, got %q", testVars.Security.KeyID, keyID)
					},
					resource.TestCheckResourceAttr(resourceNameKeyManagementServer, "access_information.0.tenant_id", testVars.Security.TenantID),
					resource.TestCheckResourceAttrSet(resourceNameKeyManagementServer, "access_information.0.truncated_client_secret"),
				),
			},
		},
	})
}

func TestAccV2NutanixKeyManagementServerResource_InvalidAccessInfo(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-kms-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testKMSResourceInvalidAccessInfoConfig(name),
				ExpectError: regexp.MustCompile("error waiting for kms to be created:"),
			},
		},
	})
}

func testKMSResourceConfig(name string) string {
	return fmt.Sprintf(`
locals {
  config = jsondecode(file("%[1]s"))
  kms = local.config.security
}
resource "nutanix_key_management_server_v2" "test" {
  name = "%[2]s"
  access_information {
    endpoint_url           = local.kms.endpoint_url
    key_id                 = local.kms.key_id
    tenant_id              = local.kms.tenant_id
    client_id              = local.kms.client_id
    client_secret          = local.kms.client_secret
    credential_expiry_date = local.kms.credential_expiry_date
  }
  lifecycle {
    ignore_changes = [
      access_information[0].client_secret,
      access_information[0].key_id
    ]
  }
}
`, filepath, name)
}

func testKMSResourceInvalidAccessInfoConfig(name string) string {
	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format("2006-01-02")
	return fmt.Sprintf(`
resource "nutanix_key_management_server_v2" "test" {
  name = "%[1]s-invalid"
  access_information {
    endpoint_url           = "https://invalid-keyvault-001.vault.azure.net/"
    key_id                 = "invalid_key_id"
    tenant_id              = "ab414ed6-7d97-4f7a-b98f-fcba7cac3b8c"
    client_id              = "ae1a2b3c-5d6e-7f80-9a1b-2c3d4e5f6789"
    client_secret          = "98765432-10fe-dcba-9876-543210fedcba"
    credential_expiry_date = "%[2]s"
  }
  lifecycle {
    ignore_changes = [
      access_information[0].client_secret,
      access_information[0].key_id
    ]
  }
}
`, name, expirationTimeFormatted)
}
