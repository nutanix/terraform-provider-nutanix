package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameIdentityProvider = "data.nutanix_saml_identity_provider_v2.test"

func TestAccV2NutanixIdentityProvidersDatasource_GetSamlIdpById(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIdentityProviderDatasourceV2Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProvider, "idp_metadata.#"),
					resource.TestCheckResourceAttr(datasourceNameIdentityProvider, "name", testVars.Iam.IdentityProviders.Name),
					resource.TestCheckResourceAttr(datasourceNameIdentityProvider, "username_attribute", testVars.Iam.IdentityProviders.UsernameAttribute),
					resource.TestCheckResourceAttr(datasourceNameIdentityProvider, "email_attribute", testVars.Iam.IdentityProviders.EmailAttribute),
					resource.TestCheckResourceAttr(datasourceNameIdentityProvider, "groups_attribute", testVars.Iam.IdentityProviders.GroupsAttribute),
					resource.TestCheckResourceAttr(datasourceNameIdentityProvider, "groups_delim", testVars.Iam.IdentityProviders.GroupsDelim),
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProvider, "custom_attributes.#"),
					resource.TestCheckResourceAttr(datasourceNameIdentityProvider, "custom_attributes.0", testVars.Iam.IdentityProviders.CustomAttributes[0]),
					resource.TestCheckResourceAttr(datasourceNameIdentityProvider, "custom_attributes.1", testVars.Iam.IdentityProviders.CustomAttributes[1]),
				),
			},
		},
	})
}

func testIdentityProviderDatasourceV2Config(filepath string) string {
	return fmt.Sprintf(`
		locals{
			config = (jsondecode(file("%s")))
			identity_providers = local.config.iam.identity_providers
		}

		resource "nutanix_saml_identity_providers_v2" "test" {
			name = local.identity_providers.name
			username_attribute = local.identity_providers.username_attr
			email_attribute = local.identity_providers.email_attr
			groups_attribute = local.identity_providers.groups_attr
			groups_delim = local.identity_providers.groups_delim
			idp_metadata_xml = file("%[2]s") # xml content
			entity_issuer = local.identity_providers.entity_issuer
			is_signed_authn_req_enabled = local.identity_providers.is_signed_authn_req_enabled
			custom_attributes = local.identity_providers.custom_attributes
		}

		data "nutanix_saml_identity_provider_v2" "test" {
			ext_id = nutanix_saml_identity_providers_v2.test.id
		}
`, filepath, xmlFilePath)
}
