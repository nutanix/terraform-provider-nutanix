// resource_nutanix_vm_guest_customization_profile_v2_test.go
//
// Acceptance tests for the nutanix_vm_guest_customization_profile_v2 resource
// and its associated datasources (singular and plural).
//
// Test Cases:
//
//   TestAccV2NutanixVmGuestCustomizationProfileResource_SysprepParamsWithUpdates
//     - Comprehensive 10-step create + update + datasource test covering the
//       full lifecycle of a sysprep_params-based GC profile.
//     - Step 1 (Create): Full initial config — admin password, timezone,
//       registered org/owner, computer_name=must_provide, locale (en-US),
//       workgroup, DHCP, DNS server. Validates all attributes including
//       ext_id, create_time, update_time, links, config structure.
//     - Step 2 (Update name/desc/org/owner): Changes name, description,
//       registered_organization, registered_owner. Validates updated values.
//     - Step 3 (Update locale): Switches locale from en-US to fr-FR for
//       ui_language, system_locale, user_locale.
//     - Step 4 (Update computer_name): Switches from must_provide to
//       use_vm_name=true.
//     - Step 5 (Update workgroup): Changes workgroup name from "WORKGROUP"
//       to "UpdatedWG".
//     - Step 6 (Workgroup → Domain): Removes workgroup, adds domain_settings
//       with credentials (domain_primary from test config). Validates domain
//       settings are present and workgroup is removed.
//     - Step 7 (Update timezone): Changes timezone from "Pacific Standard
//       Time" to "India Standard Time".
//     - Step 8 (Update IPv4 config): Switches from use_dhcp=true to
//       must_provide_during_deployment=true.
//     - Step 9 (Singular datasource): Reads the profile via
//       nutanix_vm_guest_customization_profile_v2 datasource and validates
//       all attributes match the resource state (name, desc, timezone, org,
//       owner, computer_name, locale, domain, network, DNS).
//     - Step 10 (Plural datasource): Lists profiles via
//       nutanix_vm_guest_customization_profiles_v2 datasource with name
//       filter. Validates profiles.# = 1 and all attributes match including
//       nested sysprep_params fields.
//     - CheckDestroy: Verifies the profile is deleted from the API after
//       test completion.
//
//   TestAccV2NutanixVmGuestCustomizationProfileResource_AnswerFile
//     - Creates a GC profile using answer_file (unattend.xml) instead of
//       sysprep_params.
//     - Validates the answer_file block is present with unattend_xml set.
//     - CheckDestroy: Verifies profile deletion.

package vmmv2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gcpRequest "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vmguestcustomizationprofiles"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	resourceNameVmGcProfile    = "nutanix_vm_guest_customization_profile_v2.test"
	datasourceNameVmGcProfile  = "data.nutanix_vm_guest_customization_profile_v2.test"
	datasourceNameVmGcProfiles = "data.nutanix_vm_guest_customization_profiles_v2.test"
)

// gcProfileTestCfg holds all configurable fields for generating GC profile HCL.
type gcProfileTestCfg struct {
	Name                string
	Desc                string
	Org                 string
	Owner               string
	AdminPass           string
	Timezone            string
	UILang              string
	SysLocale           string
	UserLocale          string
	UseVMName           bool
	MustProvideCompName bool
	Workgroup           string
	DomainName          string
	DomainUser          string
	DomainPass          string
	UseDHCP             bool
	MustProvideIP       bool
	DNSServer           string
}

func buildGcProfileConfig(c gcProfileTestCfg) string {
	computerNameBlock := ""
	if c.UseVMName {
		computerNameBlock = `
								computer_name {
									use_vm_name = true
								}`
	} else if c.MustProvideCompName {
		computerNameBlock = `
								computer_name {
									must_provide_during_deployment = true
								}`
	}

	adminPassLine := ""
	if c.AdminPass != "" {
		adminPassLine = fmt.Sprintf(`
								administrator_password  = "%s"`, c.AdminPass)
	}

	orgLine := ""
	if c.Org != "" {
		orgLine = fmt.Sprintf(`
								registered_organization = "%s"`, c.Org)
	}

	ownerLine := ""
	if c.Owner != "" {
		ownerLine = fmt.Sprintf(`
								registered_owner        = "%s"`, c.Owner)
	}

	timezoneLine := ""
	if c.Timezone != "" {
		timezoneLine = fmt.Sprintf(`
								timezone = "%s"`, c.Timezone)
	}

	localeBlock := ""
	if c.UILang != "" || c.SysLocale != "" || c.UserLocale != "" {
		localeBlock = fmt.Sprintf(`
							locale_settings {
								ui_language   = "%s"
								system_locale = "%s"
								user_locale   = "%s"
							}`, c.UILang, c.SysLocale, c.UserLocale)
	}

	wgOrDomainBlock := ""
	if c.Workgroup != "" {
		wgOrDomainBlock = fmt.Sprintf(`
							workgroup_or_domain_info {
								workgroup {
									name = "%s"
								}
							}`, c.Workgroup)
	} else if c.DomainName != "" {
		wgOrDomainBlock = fmt.Sprintf(`
							workgroup_or_domain_info {
								domain_settings {
									credentials {
										domain_name = "%s"
										password    = "%s"
										username    = "%s"
									}
								}
							}`, c.DomainName, c.DomainPass, c.DomainUser)
	}

	networkBlock := ""
	if c.UseDHCP || c.MustProvideIP {
		dnsBlock := ""
		if c.DNSServer != "" {
			dnsBlock = fmt.Sprintf(`
									dns_config {
										preferred_dns_server_address = "%s"
									}`, c.DNSServer)
		}
		ipv4Block := ""
		if c.UseDHCP {
			ipv4Block = `
									ipv4_config {
										use_dhcp = true
									}`
		} else if c.MustProvideIP {
			ipv4Block = `
									ipv4_config {
										must_provide_during_deployment = true
									}`
		}
		networkBlock = fmt.Sprintf(`
							network_settings {
								nic_config_list {%s%s
								}
							}`, dnsBlock, ipv4Block)
	}

	lifecycleBlock := ""
	if c.AdminPass != "" && c.DomainPass != "" {
		lifecycleBlock = `
			lifecycle {
				ignore_changes = [
					config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].administrator_password,
					config[0].sysprep_config[0].customization[0].sysprep_params[0].workgroup_or_domain_info[0].domain_settings[0].credentials[0].password
				]
			}`
	} else if c.AdminPass != "" {
		lifecycleBlock = `
			lifecycle {
				ignore_changes = [
					config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].administrator_password
				]
			}`
	}

	return fmt.Sprintf(`
		resource "nutanix_vm_guest_customization_profile_v2" "test" {
			name        = "%s"
			description = "%s"
			config {
				sysprep_config {
					customization {
						sysprep_params {
							general_settings {%s%s%s%s%s
							}%s%s%s
						}
					}
				}
			}%s
		}
`, c.Name, c.Desc, adminPassLine, orgLine, ownerLine, timezoneLine, computerNameBlock,
		localeBlock, networkBlock, wgOrDomainBlock, lifecycleBlock)
}

// ============================================================================
// Comprehensive create + multi-attribute update test
// ============================================================================

func TestAccV2NutanixVmGuestCustomizationProfileResource_SysprepParamsWithUpdates(t *testing.T) {
	r := acctest.RandInt()
	baseName := fmt.Sprintf("test-gc-profile-%d", r)

	sp := "config.0.sysprep_config.0.customization.0.sysprep_params.0."
	lsp := "profiles.0.config.0.sysprep_config.0.customization.0.sysprep_params.0."

	// Step 1: Create with full initial config
	step1 := gcProfileTestCfg{
		Name:                baseName,
		Desc:                "gc profile create and update test",
		Org:                 "TestOrg",
		Owner:               "TestOwner",
		AdminPass:           "Noice/4u",
		Timezone:            "Pacific Standard Time",
		UILang:              "en-US",
		SysLocale:           "en-US",
		UserLocale:          "en-US",
		MustProvideCompName: true,
		Workgroup:           "WORKGROUP",
		UseDHCP:             true,
		DNSServer:           testVars.Iam.DirectoryServicesMain.PrimaryAD.DNS,
	}

	// Step 2: Update name, description, organization, owner
	step2 := step1
	step2.Name = fmt.Sprintf("test-gc-profile-%d-upd", r)
	step2.Desc = "updated description"
	step2.Org = "UpdatedOrg"
	step2.Owner = "UpdatedOwner"

	// Step 3: Update locale settings (en-US → fr-FR)
	step3 := step2
	step3.UILang = "fr-FR"
	step3.SysLocale = "fr-FR"
	step3.UserLocale = "fr-FR"

	// Step 4: Update computer_name from must_provide to use_vm_name
	step4 := step3
	step4.MustProvideCompName = false
	step4.UseVMName = true

	// Step 5: Update workgroup name only
	step5 := step4
	step5.Workgroup = "UpdatedWG"

	// Step 6: Switch from workgroup to domain
	step6 := step5
	step6.Workgroup = ""
	step6.DomainName = testVars.Iam.DirectoryServicesMain.PrimaryAD.Name
	step6.DomainUser = testVars.Iam.DirectoryServicesMain.PrimaryAD.Username
	step6.DomainPass = testVars.Iam.DirectoryServicesMain.PrimaryAD.Password

	// Step 7: Update timezone
	step7 := step6
	step7.Timezone = "India Standard Time"

	// Step 8: Update ipv4_config from use_dhcp to must_provide_during_deployment
	step8 := step7
	step8.UseDHCP = false
	step8.MustProvideIP = true

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmGcProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildGcProfileConfig(step1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", step1.Name),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", step1.Desc),
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"general_settings.0.registered_organization", "TestOrg"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"general_settings.0.registered_owner", "TestOwner"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"general_settings.0.timezone", "Pacific Standard Time"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"general_settings.0.computer_name.0.must_provide_during_deployment", "true"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"locale_settings.0.ui_language", "en-US"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"locale_settings.0.system_locale", "en-US"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"locale_settings.0.user_locale", "en-US"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"workgroup_or_domain_info.0.workgroup.0.name", "WORKGROUP"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"network_settings.0.nic_config_list.0.ipv4_config.0.use_dhcp", "true"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"network_settings.0.nic_config_list.0.dns_config.0.preferred_dns_server_address", testVars.Iam.DirectoryServicesMain.PrimaryAD.DNS),
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "links.#"),
				),
			},
			{
				Config: buildGcProfileConfig(step2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", step2.Name),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", step2.Desc),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"general_settings.0.registered_organization", "UpdatedOrg"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"general_settings.0.registered_owner", "UpdatedOwner"),
				),
			},
			{
				Config: buildGcProfileConfig(step3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"locale_settings.0.ui_language", "fr-FR"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"locale_settings.0.system_locale", "fr-FR"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"locale_settings.0.user_locale", "fr-FR"),
				),
			},
			{
				Config: buildGcProfileConfig(step4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"general_settings.0.computer_name.0.use_vm_name", "true"),
				),
			},
			{
				Config: buildGcProfileConfig(step5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"workgroup_or_domain_info.0.workgroup.0.name", "UpdatedWG"),
				),
			},
			{
				Config: buildGcProfileConfig(step6),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"workgroup_or_domain_info.0.domain_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"workgroup_or_domain_info.0.domain_settings.0.credentials.0.domain_name", testVars.Iam.DirectoryServicesMain.PrimaryAD.Name),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"workgroup_or_domain_info.0.domain_settings.0.credentials.0.username", testVars.Iam.DirectoryServicesMain.PrimaryAD.Username),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"workgroup_or_domain_info.0.workgroup.#", "0"),
				),
			},
			{
				Config: buildGcProfileConfig(step7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"general_settings.0.timezone", "India Standard Time"),
				),
			},
			{
				Config: buildGcProfileConfig(step8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, sp+"network_settings.0.nic_config_list.0.ipv4_config.0.must_provide_during_deployment", "true"),
				),
			},
			// Step 9: Validate singular datasource read with full attribute checks
			{
				Config: buildGcProfileConfig(step8) + `
					data "nutanix_vm_guest_customization_profile_v2" "test" {
						ext_id = nutanix_vm_guest_customization_profile_v2.test.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfile, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, "name", step2.Name),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, "description", step2.Desc),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfile, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfile, "update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfile, "links.#"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, "config.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, "config.0.sysprep_config.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"general_settings.0.timezone", "India Standard Time"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"general_settings.0.registered_organization", "UpdatedOrg"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"general_settings.0.registered_owner", "UpdatedOwner"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"general_settings.0.computer_name.0.use_vm_name", "true"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"locale_settings.0.ui_language", "fr-FR"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"locale_settings.0.system_locale", "fr-FR"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"locale_settings.0.user_locale", "fr-FR"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"workgroup_or_domain_info.0.domain_settings.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"workgroup_or_domain_info.0.domain_settings.0.credentials.0.domain_name", testVars.Iam.DirectoryServicesMain.PrimaryAD.Name),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"network_settings.0.nic_config_list.0.ipv4_config.0.must_provide_during_deployment", "true"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, sp+"network_settings.0.nic_config_list.0.dns_config.0.preferred_dns_server_address", testVars.Iam.DirectoryServicesMain.PrimaryAD.DNS),
				),
			},
			// Step 10: Validate plural datasource list with filter and full attribute checks
			{
				Config: buildGcProfileConfig(step8) + fmt.Sprintf(`
					data "nutanix_vm_guest_customization_profiles_v2" "test" {
						filter     = "name eq '%s'"
						depends_on = [nutanix_vm_guest_customization_profile_v2.test]
					}
				`, step2.Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, "profiles.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfiles, "profiles.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, "profiles.0.name", step2.Name),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, "profiles.0.description", step2.Desc),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfiles, "profiles.0.create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfiles, "profiles.0.update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfiles, "profiles.0.links.#"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, "profiles.0.config.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, "profiles.0.config.0.sysprep_config.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"general_settings.0.timezone", "India Standard Time"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"general_settings.0.registered_organization", "UpdatedOrg"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"general_settings.0.registered_owner", "UpdatedOwner"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"general_settings.0.computer_name.0.use_vm_name", "true"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"locale_settings.0.ui_language", "fr-FR"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"locale_settings.0.system_locale", "fr-FR"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"locale_settings.0.user_locale", "fr-FR"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"workgroup_or_domain_info.0.domain_settings.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"workgroup_or_domain_info.0.domain_settings.0.credentials.0.domain_name", testVars.Iam.DirectoryServicesMain.PrimaryAD.Name),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"network_settings.0.nic_config_list.0.ipv4_config.0.must_provide_during_deployment", "true"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfiles, lsp+"network_settings.0.nic_config_list.0.dns_config.0.preferred_dns_server_address", testVars.Iam.DirectoryServicesMain.PrimaryAD.DNS),
				),
			},
		},
	})
}

// ============================================================================
// Answer file create test
// ============================================================================

func TestAccV2NutanixVmGuestCustomizationProfileResource_AnswerFile(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-af-%d", r)
	desc := "answer file profile"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmGcProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfileAnswerFileConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.answer_file.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.answer_file.0.unattend_xml"),
				),
			},
		},
	})
}

func testAccCheckNutanixVmGcProfileDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	client := conn.VmmAPI.VmGuestCustomizationProfilesAPIInstance

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_vm_guest_customization_profile_v2" {
			continue
		}
		req := gcpRequest.GetVmGuestCustomizationProfileByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := client.GetVmGuestCustomizationProfileById(context.Background(), &req)
		if err == nil {
			return fmt.Errorf("VM Guest Customization Profile still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testVmGcProfileAnswerFileConfig(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_vm_guest_customization_profile_v2" "test" {
			name        = "%s"
			description = "%s"
			config {
				sysprep_config {
					customization {
						answer_file {
							unattend_xml = "<unattend xmlns='urn:schemas-microsoft-com:unattend'></unattend>"
						}
					}
				}
			}
		}
`, name, desc)
}
