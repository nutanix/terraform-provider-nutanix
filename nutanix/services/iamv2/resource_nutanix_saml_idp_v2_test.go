package iamv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameIdentityProviders = "nutanix_saml_identity_providers_v2.test"

func TestAccV2NutanixIdentityProvidersResource_CreateSamlIdp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIdentityProvidersResourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameIdentityProviders, "idp_metadata.#"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "name", testVars.Iam.IdentityProviders.Name),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "username_attribute", testVars.Iam.IdentityProviders.UsernameAttribute),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "email_attribute", testVars.Iam.IdentityProviders.EmailAttribute),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "groups_attribute", testVars.Iam.IdentityProviders.GroupsAttribute),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "groups_delim", testVars.Iam.IdentityProviders.GroupsDelim),
					resource.TestCheckResourceAttrSet(resourceNameIdentityProviders, "custom_attributes.#"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "custom_attributes.0", testVars.Iam.IdentityProviders.CustomAttributes[0]),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "custom_attributes.1", testVars.Iam.IdentityProviders.CustomAttributes[1]),
				),
			},
		},
	})
}

func TestAccV2NutanixIdentityProvidersResourceWithNoName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testIdentityProvidersResourceWithoutName(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixIdentityProvidersResourceWithNoEntityId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testIdentityProvidersResourceWithoutEntityID(filepath),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testIdentityProvidersResourceConfig(filepath string) string {
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
	}`, filepath, xmlFilePath)
}

func testIdentityProvidersResourceWithoutName(filepath string) string {
	return fmt.Sprintf(`

		locals{
			config = (jsondecode(file("%s")))
			identity_providers = local.config.iam.identity_providers
		}

		resource "nutanix_saml_identity_providers_v2" "test" {
			idp_metadata {
				entity_id = local.identity_providers.idp_metadata.entity_id
				login_url = local.identity_providers.idp_metadata.login_url
				logout_url = local.identity_providers.idp_metadata.logout_url
				certificate = local.identity_providers.idp_metadata.certificate
				name_id_policy_format = local.identity_providers.idp_metadata.name_id_policy_format
			}
			username_attribute = local.identity_providers.username_attr
			email_attribute = local.identity_providers.email_attr
			entity_issuer = local.identity_providers.entity_issuer
			is_signed_authn_req_enabled = local.identity_providers.is_signed_authn_req_enabled
		}`, filepath)
}

func testIdentityProvidersResourceWithoutEntityID(filepath string) string {
	return fmt.Sprintf(`

	locals{
		config = (jsondecode(file("%s")))
		identity_providers = local.config.iam.identity_providers
	}

	resource "nutanix_saml_identity_providers_v2" "test" {
		idp_metadata {
			login_url = local.identity_providers.idp_metadata.login_url
			logout_url = local.identity_providers.idp_metadata.logout_url
			certificate = local.identity_providers.idp_metadata.certificate
			name_id_policy_format = local.identity_providers.idp_metadata.name_id_policy_format
		}
		name = local.identity_providers.name
		username_attribute = local.identity_providers.username_attr
		email_attribute = local.identity_providers.email_attr
		entity_issuer = local.identity_providers.entity_issuer
		is_signed_authn_req_enabled = local.identity_providers.is_signed_authn_req_enabled
	}`, filepath)
}
