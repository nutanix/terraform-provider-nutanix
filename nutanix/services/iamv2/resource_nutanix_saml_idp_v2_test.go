package iamv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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

const datasourceNameSamlProject = "data.nutanix_saml_identity_provider_v2.project_test"
const datasourceNameSamlProjectList = "data.nutanix_saml_identity_providers_v2.list_test"

func TestAccV2NutanixIdentityProvidersResource_DefaultProjectAndSharing(t *testing.T) {
	r := acctest.RandInt()
	projectName := fmt.Sprintf("tf-saml-share-proj-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with no sharing
			{
				Config: testSamlIdpProjectConfig(projectName, ",", "none"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "name", testVars.Iam.IdentityProviders.Name),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "shared_with_projects.#", "0"),
				),
			},
			// Step 2: Share with specific project
			{
				Config: testSamlIdpProjectConfig(projectName, ";", "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "groups_delim", ";"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameIdentityProviders, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProjectList, "identity_providers.0.ext_id", resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProjectList, "identity_providers.0.shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
				),
			},
			// Step 3: Unshare specific project
			{
				Config: testSamlIdpProjectConfig(projectName, ",", "empty"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "groups_delim", ","),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.#", "0"),
				),
			},
			// Step 4: Re-share with specific project
			{
				Config: testSamlIdpProjectConfig(projectName, ";", "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "groups_delim", ";"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameIdentityProviders, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProjectList, "identity_providers.0.ext_id", resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.shared_with_projects.#", "1"),
				),
			},
			// Step 5: Unshare specific project + Share with all in one go
			{
				Config: testSamlIdpProjectConfig(projectName, ",", "share_all"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "groups_delim", ","),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "true"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "true"),
				),
			},
		},
	})
}

func TestAccV2NutanixIdentityProvidersResource_NonDefaultProjectSharingFails(t *testing.T) {
	r := acctest.RandInt()
	project1Name := fmt.Sprintf("tf-saml-proj1-%d", r)
	project2Name := fmt.Sprintf("tf-saml-proj2-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSamlIdpNonDefaultProjectConfig(project1Name, project2Name, "none"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameIdentityProviders, "project_ext_id", "nutanix_project_v2.project1", "ext_id"),
				),
			},
			{
				Config:      testSamlIdpNonDefaultProjectConfig(project1Name, project2Name, "share"),
				ExpectError: regexp.MustCompile("SAML IdP belonging to non-default project is not allowed to be shared/unshared"),
			},
		},
	})
}

func testSamlIdpProjectConfig(projectName, groupsDelim, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `
		share_with_all_projects = false
		shared_with_projects    = [nutanix_project_v2.share_project.ext_id]`
	case "empty":
		shareBlock = `
		share_with_all_projects = false
		shared_with_projects    = []`
	case "share_all":
		shareBlock = `share_with_all_projects = true`
	case "unshare_all":
		shareBlock = `share_with_all_projects = false`
	}
	listDataSource := ""
	if shareState == "share" {
		listDataSource = `
	data "nutanix_saml_identity_providers_v2" "list_test" {
		filter     = "sharedWithProjects/any(p:p eq '${nutanix_project_v2.share_project.id}')"
		depends_on = [nutanix_saml_identity_providers_v2.test]
	}`
	} else if shareState == "empty" {
		listDataSource = `
	data "nutanix_saml_identity_providers_v2" "list_test" {
		filter     = "sharedWithProjects/any(p:p eq '${nutanix_project_v2.share_project.id}')"
		depends_on = [nutanix_saml_identity_providers_v2.test]
	}`
	}
	return fmt.Sprintf(`
	locals {
		config             = jsondecode(file("%[1]s"))
		identity_providers = local.config.iam.identity_providers
	}

	resource "nutanix_project_v2" "share_project" {
		name       = "%[2]s"
		project_id = "%[2]s"
		description = "project for saml idp sharing test"
	}

	resource "nutanix_saml_identity_providers_v2" "test" {
		name = local.identity_providers.name
		username_attribute = local.identity_providers.username_attr
		email_attribute = local.identity_providers.email_attr
		groups_attribute = local.identity_providers.groups_attr
		groups_delim               = "%[3]s"
		idp_metadata_xml           = file("%[5]s")
		idp_metadata {
			entity_id = local.identity_providers.idp_metadata.entity_id
			login_url = local.identity_providers.idp_metadata.login_url
			logout_url = local.identity_providers.idp_metadata.logout_url
			certificate = local.identity_providers.idp_metadata.certificate
			name_id_policy_format = local.identity_providers.idp_metadata.name_id_policy_format
		}
		entity_issuer = local.identity_providers.entity_issuer
		%[4]s
		custom_attributes = local.identity_providers.custom_attributes
		depends_on = [nutanix_project_v2.share_project]
		lifecycle {
			ignore_changes = [
			  idp_metadata,
		  ]
		}
	}

	data "nutanix_saml_identity_provider_v2" "project_test" {
		ext_id = nutanix_saml_identity_providers_v2.test.id
	}
	%[6]s
`, filepath, projectName, groupsDelim, shareBlock, xmlFilePath, listDataSource)
}

func testSamlIdpNonDefaultProjectConfig(proj1, proj2, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `shared_with_projects = [nutanix_project_v2.project2.ext_id]`
	case "empty":
		shareBlock = `shared_with_projects = []`
	}
	return fmt.Sprintf(`
	locals {
		config             = jsondecode(file("%[1]s"))
		identity_providers = local.config.iam.identity_providers
	}

	resource "nutanix_project_v2" "project1" {
		name       = "%[2]s"
		project_id = "%[2]s"
		description = "first project for saml idp test"
	}

	resource "nutanix_project_v2" "project2" {
		name       = "%[3]s"
		project_id = "%[3]s"
		description = "second project for saml idp test"
	}

	resource "nutanix_saml_identity_providers_v2" "test" {
		name                       = local.identity_providers.name
		username_attribute         = local.identity_providers.username_attr
		email_attribute            = local.identity_providers.email_attr
		groups_attribute           = local.identity_providers.groups_attr
		groups_delim               = local.identity_providers.groups_delim
		idp_metadata_xml           = file("%[5]s")
		idp_metadata {
			entity_id = local.identity_providers.idp_metadata.entity_id
			login_url = local.identity_providers.idp_metadata.login_url
			logout_url = local.identity_providers.idp_metadata.logout_url
			certificate = local.identity_providers.idp_metadata.certificate
			name_id_policy_format = local.identity_providers.idp_metadata.name_id_policy_format
    }
		entity_issuer              = local.identity_providers.entity_issuer
		custom_attributes          = local.identity_providers.custom_attributes
		project_ext_id             = nutanix_project_v2.project1.ext_id
		%[4]s
		depends_on = [nutanix_project_v2.project1, nutanix_project_v2.project2]
		lifecycle {
			ignore_changes = [
			  idp_metadata,
		  ]
		}
	}
`, filepath, proj1, proj2, shareBlock, xmlFilePath)
}

func TestAccV2NutanixIdentityProvidersResource_ShareWithAllProjects(t *testing.T) {
	r := acctest.RandInt()
	projectName := fmt.Sprintf("tf-saml-shareall-proj-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create with no sharing
			{
				Config: testSamlIdpShareAllConfig(projectName, ",", "none"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "name", testVars.Iam.IdentityProviders.Name),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProjectList, "identity_providers.0.ext_id", resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.shared_with_projects.#", "0"),
				),
			},
			// Step 2: Share with specific project
			{
				Config: testSamlIdpShareAllConfig(projectName, ";", "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "groups_delim", ";"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameIdentityProviders, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProjectList, "identity_providers.0.ext_id", resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.shared_with_projects.#", "1"),
				),
			},
			// Step 3: Share with all projects
			{
				Config: testSamlIdpShareAllConfig(projectName, ",", "share_all"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "groups_delim", ","),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "true"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "true"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProjectList, "identity_providers.0.ext_id", resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.share_with_all_projects", "true"),
				),
			},
			// Step 4: Unshare all + Share with specific project in one go
			{
				Config: testSamlIdpShareAllConfig(projectName, ",", "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(resourceNameIdentityProviders, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameIdentityProviders, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameSamlProjectList, "identity_providers.0.ext_id", resourceNameIdentityProviders, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.share_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameSamlProjectList, "identity_providers.0.shared_with_projects.#", "1"),
				),
			},
		},
	})
}

func testSamlIdpShareAllConfig(projectName, groupsDelim, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `
		share_with_all_projects = false
		shared_with_projects    = [nutanix_project_v2.share_project.ext_id]`
	case "share_all":
		shareBlock = `share_with_all_projects = true`
	case "unshare_all":
		shareBlock = `share_with_all_projects = false`
	}
	return fmt.Sprintf(`
	locals {
		config             = jsondecode(file("%[1]s"))
		identity_providers = local.config.iam.identity_providers
	}

	resource "nutanix_project_v2" "share_project" {
		name       = "%[2]s"
		project_id = "%[2]s"
		description = "project for saml idp share-all test"
	}

	resource "nutanix_saml_identity_providers_v2" "test" {
		name                       = local.identity_providers.name
		username_attribute         = local.identity_providers.username_attr
		email_attribute            = local.identity_providers.email_attr
		groups_attribute           = local.identity_providers.groups_attr
		groups_delim               = "%[3]s"
		idp_metadata_xml           = file("%[5]s")
		idp_metadata {
			entity_id = local.identity_providers.idp_metadata.entity_id
			login_url = local.identity_providers.idp_metadata.login_url
			logout_url = local.identity_providers.idp_metadata.logout_url
			certificate = local.identity_providers.idp_metadata.certificate
			name_id_policy_format = local.identity_providers.idp_metadata.name_id_policy_format
		}
		entity_issuer              = local.identity_providers.entity_issuer
		custom_attributes          = local.identity_providers.custom_attributes
		%[4]s
		depends_on = [nutanix_project_v2.share_project]
		lifecycle {
			ignore_changes = [
			  idp_metadata,
		  ]
		}
	}

	data "nutanix_saml_identity_provider_v2" "project_test" {
		ext_id = nutanix_saml_identity_providers_v2.test.id
	}

	data "nutanix_saml_identity_providers_v2" "list_test" {
		filter     = "name eq '${nutanix_saml_identity_providers_v2.test.name}'"
		depends_on = [nutanix_saml_identity_providers_v2.test]
	}
`, filepath, projectName, groupsDelim, shareBlock, xmlFilePath)
}
