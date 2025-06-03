package iamv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameIdentityProviders = "data.nutanix_saml_identity_providers_v2.test"

func TestAccV2NutanixIdentityProvidersDatasource_ListAllIDPS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIdentityProvidersDatasourceConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.#"),
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.0.username_attribute"),
				),
			},
		},
	})
}

func TestAccV2NutanixIdentityProvidersDatasource_WithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIdentityProvidersDatasourceWithFilterConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.0.name", testVars.Iam.IdentityProviders.Name),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.0.username_attribute", testVars.Iam.IdentityProviders.UsernameAttribute),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.0.email_attribute", testVars.Iam.IdentityProviders.EmailAttribute),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.0.groups_attribute", testVars.Iam.IdentityProviders.GroupsAttribute),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.0.groups_delim", testVars.Iam.IdentityProviders.GroupsDelim),
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.0.custom_attributes.#"),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.0.custom_attributes.0", testVars.Iam.IdentityProviders.CustomAttributes[0]),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.0.custom_attributes.1", testVars.Iam.IdentityProviders.CustomAttributes[1]),
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.0.idp_metadata.0.certificate"),
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.0.idp_metadata.0.entity_id"),
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.0.idp_metadata.0.login_url"),
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.0.idp_metadata.0.name_id_policy_format"),
				),
			},
		},
	})
}

func TestAccV2NutanixIdentityProvidersDatasource_WithLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIdentityProvidersDatasourceWithLimitConfig(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.#"),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixIdentityProvidersDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIdentityProvidersDatasourceWithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameIdentityProviders, "identity_providers.#"),
					resource.TestCheckResourceAttr(datasourceNameIdentityProviders, "identity_providers.#", "0"),
				),
			},
		},
	})
}

func testIdentityProvidersDatasourceConfig(filepath string) string {
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
	data "nutanix_saml_identity_providers_v2" "test"{
		depends_on = [ resource.nutanix_saml_identity_providers_v2.test ]
	}
	`, filepath, xmlFilePath)
}

func testIdentityProvidersDatasourceWithFilterConfig(filepath string) string {
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

	data "nutanix_saml_identity_providers_v2" "test" {
		filter = "extId eq '${resource.nutanix_saml_identity_providers_v2.test.id}'"
		depends_on = [ resource.nutanix_saml_identity_providers_v2.test ]
	}

	`, filepath, xmlFilePath)
}

func testIdentityProvidersDatasourceWithLimitConfig(filepath string) string {
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

		data "nutanix_saml_identity_providers_v2" "test" {
			limit      = 1
			depends_on = [ resource.nutanix_saml_identity_providers_v2.test ]
		}
	`, filepath, xmlFilePath)
}

func testIdentityProvidersDatasourceWithInvalidFilterConfig() string {
	return `
	data "nutanix_saml_identity_providers_v2" "test" {
		filter = "extId eq 'invalid_filter'"
	}
	`
}
