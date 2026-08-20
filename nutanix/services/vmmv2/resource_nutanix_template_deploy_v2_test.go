// resource_nutanix_template_deploy_v2_test.go
//
// Acceptance tests for nutanix_deploy_templates_v2 (template deploy) and
// nutanix_vm_clone_v2 (VM clone) resources, including end-to-end Guest
// Customization (GC) profile validation via WinRM on deployed Windows VMs.
//
// Test Cases:
//
//   TestAccV2NutanixTemplateDeployResource_Basic
//     - Deploys a template with no GC profile or overrides.
//     - Validates basic template deploy lifecycle (create, read).
//
//   TestAccV2NutanixTemplateDeployResource_WithGcProfileConfigOverride
//     - Creates a GC profile with sysprep_params (use_vm_name, timezone, org,
//       locale, domain, DHCP) and deploys a template referencing it.
//     - Validates the GC profile ext_id flows into the deploy resource state.
//
//   TestAccV2NutanixTemplateDeployResource_GcProfileCase1_NoOverride
//     - GC profile: use_vm_name=true, domain (primary), admin password,
//       timezone, locale, DNS, DHCP.
//     - Override: profile-only (no field overrides).
//     - WinRM validation: timezone, computer_name=VM name, domain,
//       system_locale, ui_language, registered_org, registered_owner, dns_server.
//     - Validates both deployed VM and cloned VM.
//
//   TestAccV2NutanixTemplateDeployResource_GcProfileCase2_FullOverride
//     - GC profile: all fields set (admin password, timezone, org, owner,
//       locale, domain, DNS, IP=must_provide, computer_name=must_provide,
//       first_logon_commands, auto_logon_settings, windows_product_key).
//     - Override: every field overridden with different values, static IP
//       via dynamically found free IPs (nmap), domain switched to secondary,
//       first_logon_commands overridden, logon_count overridden.
//     - WinRM validation: all fields including domain join wait, IP address,
//       first_logon_marker file, auto_logon_count registry value.
//     - Template datasource validation: nutanix_template_v2 (Get) and
//       nutanix_templates_v2 (List with filter) verify GC profile ext_id
//       is present in the template response.
//     - Computer names, IPs are unique per run to avoid AD/network conflicts.
//
//   TestAccV2NutanixTemplateDeployResource_GcProfileCase3_WorkgroupNameOverride
//     - GC profile: use_vm_name=true, workgroup set, admin password, timezone.
//     - Override: workgroup name changed to a different dynamic value.
//     - WinRM validation: timezone, computer_name=VM name, workgroup.
//     - Workgroup names are unique per run.
//
//   TestAccV2NutanixTemplateDeployResource_GcProfileCase4_WorkgroupToDomain
//     - GC profile: workgroup (default), must_provide computer_name,
//       must_provide IP, admin password, timezone.
//     - Override: switches from workgroup to domain (secondary), sets
//       computer_name, static IP, DNS.
//     - WinRM validation: timezone, computer_name, domain (with join wait),
//       dns_server, ip_address.
//     - Computer names, IPs are unique per run.
//
//   TestAccV2NutanixTemplateDeployResource_GcProfileCase5_DiscardAll
//     - GC profile: use_vm_name=true, workgroup, admin password, timezone,
//       locale, DNS, DHCP.
//     - Override: DiscardAll=true for all fields except administrator_password
//       which is overridden with a known value ("Discard.123") to allow
//       WinRM login without needing the base image's default credentials.
//     - WinRM validation: timezone=UTC (OS default), system_locale=en-US
//       (OS default). Logs in directly with the overridden password.
//
//   TestAccV2NutanixTemplateDeployResource_GcProfileCase6_AnswerFileWorkgroup
//     - GC profile: answer_file (unattend.xml) with workgroup + admin password.
//     - Override: profile-only (no override of answer file).
//     - WinRM validation: workgroup matches the answer file value.
//     - Workgroup name is unique per run.
//
//   TestAccV2NutanixTemplateDeployResource_GcProfileCase7_AnswerFileWorkgroupOverride
//     - GC profile: answer_file with workgroup + admin password.
//     - Override: different answer_file with a new workgroup name.
//     - WinRM validation: workgroup matches the override answer file value.
//     - Both workgroup names are unique per run.
//
// Common infrastructure:
//   - All GC cases use a pre-created base Windows VM (UUID from test_config_v2.json)
//     as the template source, avoiding per-test VM creation overhead.
//   - CheckDestroy: testAccCheckGcDeployVMsDestroy explicitly deletes VMs
//     created by nutanix_deploy_templates_v2, since the resource's Delete is a no-op.
//   - WinRM: polls up to 30 min for VM readiness, with 30s retry interval.
//     Domain join is polled separately for up to 10 min after WinRM connects.

package vmmv2_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	resourceNameTemplateDeploy = "nutanix_deploy_templates_v2.test"
	resourceNameGcDeploy       = "nutanix_deploy_templates_v2.gc_deploy"
	resourceNameGcClone        = "nutanix_vm_clone_v2.gc_clone"
	resourceNameGcProfile      = "nutanix_vm_guest_customization_profile_v2.gc"
)

func TestAccV2NutanixTemplateDeployResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("test-temp-%d", r)
	templateDesc := "test temp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTemplateDeployV2Config(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "number_of_vms", "1"),
					resource.TestCheckResourceAttrSet(resourceNameTemplateDeploy, "override_vm_config_map.#"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.name", "test-tf-template-deploy"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.memory_size_bytes", "4294967296"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.num_cores_per_socket", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixTemplateDeployResource_WithGcProfileConfigOverride(t *testing.T) {
	r := acctest.RandInt()
	templateName := fmt.Sprintf("test-temp-gcp-%d", r)
	templateDesc := "template with guest customization profile"
	gcProfileName := fmt.Sprintf("test-gc-profile-%d", r)
	gcProfileDesc := "gc profile for template deploy test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTemplateDeployWithGcProfileConfig(templateName, templateDesc, gcProfileName, gcProfileDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_vm_guest_customization_profile_v2.test", "name", gcProfileName),
					resource.TestCheckResourceAttr("nutanix_vm_guest_customization_profile_v2.test", "config.#", "1"),
					resource.TestCheckResourceAttr("nutanix_vm_guest_customization_profile_v2.test",
						"config.0.sysprep_config.0.customization.0.sysprep_params.0.general_settings.0.computer_name.0.use_vm_name", "true"),
					resource.TestCheckResourceAttr("nutanix_vm_guest_customization_profile_v2.test",
						"config.0.sysprep_config.0.customization.0.sysprep_params.0.locale_settings.0.ui_language", "fr-FR"),

					resource.TestCheckResourceAttr("nutanix_template_v2.test", "template_name", templateName),
					resource.TestCheckResourceAttrSet("nutanix_template_v2.test", "template_version_spec.0.guest_customization_profile.0.ext_id"),

					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "number_of_vms", "1"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.name", "tf-gc-deploy"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.memory_size_bytes", "3221225472"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.num_cores_per_socket", "2"),
					resource.TestCheckResourceAttrSet(resourceNameTemplateDeploy,
						"override_vm_config_map.0.guest_customization_profile_config.0.profile.0.ext_id"),
				),
			},
		},
	})
}

func testTemplateDeployWithGcProfileConfig(tempName, tempDesc, gcName, gcDesc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config     = jsondecode(file("%[5]s"))
			gc_profile = local.config.vmm.gc_profile
		}

		# Source VM for the template. Creating a template *version* with a guest-
		# customization profile requires a compatible NGT (>= 4.5) installed on the
		# source VM, so use the pre-provisioned Windows VM (vmm.gc_profile.vm_name)
		# looked up by name. A bare VM has no
		# NGT and fails with "required NGT version is not installed on the source VM".
		data "nutanix_virtual_machines_v2" "gc_base" {
			filter = "name eq '${local.gc_profile.vm_name}'"
		}

		resource "nutanix_vm_guest_customization_profile_v2" "test" {
			name        = "%[3]s"
			description = "%[4]s"
			config {
				sysprep_config {
					customization {
						sysprep_params {
							general_settings {
								registered_organization = "TestOrg"
								registered_owner        = "TestOwner"
								computer_name {
									use_vm_name = true
								}
								timezone = "Indian Standard Time"
							}
							locale_settings {
								ui_language   = "fr-FR"
								system_locale = "fr-FR"
								user_locale   = "fr-FR"
							}
							workgroup_or_domain_info {
								workgroup {
									name = "WORKGROUP"
								}
							}
						}
					}
				}
			}
		}

		resource "nutanix_template_v2" "test" {
			template_name        = "%[1]s"
			template_description = "%[2]s"
			template_version_spec {
				version_source {
					template_vm_reference {
						ext_id = data.nutanix_virtual_machines_v2.gc_base.vms[0].ext_id
						guest_customization_profile {
							ext_id = nutanix_vm_guest_customization_profile_v2.test.id
						}
					}
				}
				guest_customization_profile {
					ext_id = nutanix_vm_guest_customization_profile_v2.test.id
				}
			}
			lifecycle {
				ignore_changes = [
					template_version_spec.0.version_name,
					template_version_spec.0.version_description,
					template_version_spec.0.version_source
				]
			}
			depends_on = [nutanix_vm_guest_customization_profile_v2.test]
		}

		resource "nutanix_deploy_templates_v2" "test" {
			ext_id            = nutanix_template_v2.test.id
			number_of_vms     = 1
			cluster_reference = local.cluster0
			override_vm_config_map {
				# use_vm_name=true in the GC profile makes this the Windows computer
				# name, which must satisfy sysprep rules (NetBIOS: <= 15 chars); keep
				# it short or the deploy subtask fails with "does not meet the computer
				# name requirements".
				name                 = "tf-gc-deploy"
				memory_size_bytes    = 3 * 1024 * 1024 * 1024
				num_cores_per_socket = 2
				num_sockets          = 2
				num_threads_per_core = 2
				guest_customization_profile_config {
					profile {
						ext_id = nutanix_vm_guest_customization_profile_v2.test.id
					}
					config_override_spec {
						sysprep_config {
							customization {
								sysprep_params {
									locale_settings {
										ui_language {
											value = "en-US"
										}
										system_locale {
											value = "en-US"
										}
										user_locale {
											value = "en-US"
										}
									}
								}
							}
						}
					}
				}
			}
			depends_on = [nutanix_template_v2.test]
		}
`, tempName, tempDesc, gcName, gcDesc, filepath)
}

func testTemplateDeployV2Config(name, desc, tempName, tempDesc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		    ][0]
			config = jsondecode(file("%[5]s"))
			vmm    = local.config.vmm
		}
	
		data "nutanix_subnets_v2" "subnets" { 
   			filter = "name eq '${local.vmm.subnet_name}'"
        }

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}
		resource "nutanix_template_v2" "test" {
			template_name = "%[3]s"
			template_description = "%[4]s"
			template_version_spec{
			  version_source{
				template_vm_reference{
				  ext_id = nutanix_virtual_machine_v2.test.id
				}
			  }
			}
			  lifecycle {
			  ignore_changes = [
				template_version_spec.0.version_name,
				template_version_spec.0.version_description,
				template_version_spec.0.version_source
			  ]
			} 
			depends_on = [nutanix_virtual_machine_v2.test]
		}

		resource "nutanix_deploy_templates_v2" "test" {
			ext_id = resource.nutanix_template_v2.test.id
			number_of_vms = 1
			cluster_reference = local.cluster0
			override_vm_config_map{
			  name= "test-tf-template-deploy"
			  memory_size_bytes = 4294967296
			  num_sockets=2
			  num_cores_per_socket=1
			  num_threads_per_core=1
			  nics{
				nic_backing_info{
				  virtual_ethernet_nic {
					is_connected = true
					model = "VIRTIO"
				  }
				}
				nic_network_info {
				  virtual_ethernet_nic_network_info {
					nic_type = "NORMAL_NIC"
					subnet {
					  ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
					}
					vlan_mode = "ACCESS"
					should_allow_unknown_macs = false
				  }
				}
			  }
			}
			depends_on = [
				resource.nutanix_template_v2.test
			]
		}	

`, name, desc, tempName, tempDesc, filepath)
}

// ============================================================================
// Shared infrastructure HCL for GC Profile test cases
// ============================================================================

// testGcProfileBaseInfra generates HCL for shared infra: cluster lookup, locals
// with gc_profile config, and subnet data source. The base VM is provisioned by
// testenv/terraform (nutanix_virtual_machine_v2.vm in vmm.tf) and looked up by
// name (vmm.gc_profile.vm_name from test_config_v2.json); its ext_id feeds the
// template + clone below via data.nutanix_virtual_machines_v2.gc_base.
func testGcProfileBaseInfra() string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config       = jsondecode(file("%s"))
			gc_profile   = local.config.vmm.gc_profile
			gc_subnet    = local.config.networking.gc_subnet
			primary_ad   = local.config.iam.directory_services_main.primary_ad
			secondary_ad = local.config.iam.directory_services_main.secondary_ad
		}

		data "nutanix_subnets_v2" "gc_subnet_lookup" {
			filter = "name eq '${local.gc_subnet.name}'"
		}

		data "nutanix_virtual_machines_v2" "gc_base" {
			filter = "name eq '${local.gc_profile.vm_name}'"
		}
`, filepath)
}

// testGcTemplateAndDeploy generates HCL for template creation from the base VM
// with a GC profile, template deployment, and VM clone. The deployOverrideHCL
// and cloneOverrideHCL params contain the guest_customization_profile_config
// block content for each.
func testGcTemplateAndDeploy(templateName, templateDesc, deployName, cloneName, deployOverrideHCL, cloneOverrideHCL string) string {
	return fmt.Sprintf(`
		resource "nutanix_template_v2" "gc_template" {
			template_name        = "%[1]s"
			template_description = "%[2]s"
			template_version_spec {
				version_source {
					template_vm_reference {
						ext_id = data.nutanix_virtual_machines_v2.gc_base.vms[0].ext_id
						guest_customization_profile {
							ext_id = nutanix_vm_guest_customization_profile_v2.gc.id
						}
					}
				}
				guest_customization_profile {
					ext_id = nutanix_vm_guest_customization_profile_v2.gc.id
				}
			}
			lifecycle {
				ignore_changes = [
					template_version_spec.0.version_name,
					template_version_spec.0.version_description,
					template_version_spec.0.version_source
				]
			}
			depends_on = [nutanix_vm_guest_customization_profile_v2.gc]
		}

		resource "nutanix_deploy_templates_v2" "gc_deploy" {
			ext_id            = nutanix_template_v2.gc_template.id
			number_of_vms     = 1
			cluster_reference = local.cluster0
			override_vm_config_map {
				name                 = "%[3]s"
				memory_size_bytes    = 3 * 1024 * 1024 * 1024
				num_cores_per_socket = 1
				num_sockets          = 1
				num_threads_per_core = 1
				%[5]s
			}
			depends_on = [nutanix_template_v2.gc_template]
		}

		resource "nutanix_vm_clone_v2" "gc_clone" {
			vm_ext_id          = data.nutanix_virtual_machines_v2.gc_base.vms[0].ext_id
			name               = "%[4]s"
			memory_size_bytes  = 3 * 1024 * 1024 * 1024
			num_cores_per_socket = 1
			num_sockets          = 1
			num_threads_per_core = 1
			%[6]s
			depends_on = [nutanix_vm_guest_customization_profile_v2.gc]
		}
`, templateName, templateDesc, deployName, cloneName, deployOverrideHCL, cloneOverrideHCL)
}

// ============================================================================
// Builder types and functions for GC profile HCL generation
// ============================================================================

// gcDeployProfileCfg holds all configurable fields for generating
// the nutanix_vm_guest_customization_profile_v2 "gc" resource HCL.
type gcDeployProfileCfg struct {
	Name                string
	Desc                string
	Org                 string
	Owner               string
	AdminPass           string
	Timezone            string
	UseVMName           bool
	MustProvideCompName bool
	UILang              string
	SysLocale           string
	UserLocale          string
	Workgroup           string
	DomainRef           string // "primary_ad" or "secondary_ad" -> local.<DomainRef>.*
	UseDHCP             bool
	MustProvideIP       bool
	DNSRef              string // "primary_ad" or "secondary_ad" -> local.<DNSRef>.dns
	AnswerFileXML       string // if non-empty, uses answer_file instead of sysprep_params
	FirstLogonCommands  []string
	LogonCount          int
	WindowsProductKey   string
}

func buildGcDeployProfileHCL(c gcDeployProfileCfg) string {
	if c.AnswerFileXML != "" {
		return fmt.Sprintf(`
		resource "nutanix_vm_guest_customization_profile_v2" "gc" {
			name        = "%s"
			description = "%s"
			config {
				sysprep_config {
					customization {
						answer_file {
							unattend_xml = "%s"
						}
					}
				}
			}
		}
`, c.Name, c.Desc, c.AnswerFileXML)
	}

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

	productKeyLine := ""
	if c.WindowsProductKey != "" {
		productKeyLine = fmt.Sprintf(`
								windows_product_key = "%s"`, c.WindowsProductKey)
	}

	autoLogonBlock := ""
	if c.LogonCount > 0 {
		autoLogonBlock = fmt.Sprintf(`
								auto_logon_settings {
									logon_count = %d
								}`, c.LogonCount)
	}

	firstLogonBlock := ""
	if len(c.FirstLogonCommands) > 0 {
		cmds := ""
		for _, cmd := range c.FirstLogonCommands {
			cmds += fmt.Sprintf(`"%s", `, cmd)
		}
		firstLogonBlock = fmt.Sprintf(`
							first_logon_commands = [%s]`, cmds)
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
	} else if c.DomainRef != "" {
		wgOrDomainBlock = fmt.Sprintf(`
							workgroup_or_domain_info {
								domain_settings {
									credentials {
										domain_name = local.%s.name
										password    = local.%s.password
										username    = local.%s.username
									}
								}
							}`, c.DomainRef, c.DomainRef, c.DomainRef)
	}

	networkBlock := ""
	if c.UseDHCP || c.MustProvideIP || c.DNSRef != "" {
		dnsBlock := ""
		if c.DNSRef != "" {
			dnsBlock = fmt.Sprintf(`
									dns_config {
										preferred_dns_server_address = local.%s.dns
									}`, c.DNSRef)
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

	var ignoreFields []string
	if c.AdminPass != "" {
		ignoreFields = append(ignoreFields, `config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].administrator_password`)
	}
	if c.DomainRef != "" {
		ignoreFields = append(ignoreFields, `config[0].sysprep_config[0].customization[0].sysprep_params[0].workgroup_or_domain_info[0].domain_settings[0].credentials[0].password`)
	}
	if c.WindowsProductKey != "" {
		ignoreFields = append(ignoreFields, `config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].windows_product_key`)
	}

	lifecycleBlock := ""
	if len(ignoreFields) > 0 {
		lifecycleBlock = fmt.Sprintf(`
			lifecycle {
				ignore_changes = [
					%s
				]
			}`, strings.Join(ignoreFields, ",\n\t\t\t\t\t"))
	}

	return fmt.Sprintf(`
		resource "nutanix_vm_guest_customization_profile_v2" "gc" {
			name        = "%s"
			description = "%s"
			config {
				sysprep_config {
					customization {
						sysprep_params {%s
							general_settings {%s%s%s%s%s%s%s
							}%s%s%s
						}
					}
				}
			}%s
		}
`, c.Name, c.Desc, firstLogonBlock, adminPassLine, orgLine, ownerLine, timezoneLine,
		productKeyLine, autoLogonBlock, computerNameBlock,
		localeBlock, networkBlock, wgOrDomainBlock, lifecycleBlock)
}

// gcOverrideCfg holds configurable fields for generating the
// guest_customization_profile_config override block used in deploy and clone.
type gcOverrideCfg struct {
	ProfileOnly        bool
	DiscardAll         bool
	AdminPass          string
	ComputerName       string
	Timezone           string
	Org                string
	Owner              string
	UILang             string
	SysLocale          string
	UserLocale         string
	DNSRef             string // "primary_ad" or "secondary_ad"
	IPValue            string // HCL expression for ip_address value
	WithSubnet         bool   // adds prefix_length + default_gateways from subnet
	DomainRef          string // "primary_ad" or "secondary_ad"
	Workgroup          string
	AnswerFileXML      string
	FirstLogonCommands []string
	LogonCount         int
	WindowsProductKey  string
}

func buildOverrideHCL(c gcOverrideCfg) string {
	if c.ProfileOnly {
		return `
				guest_customization_profile_config {
					profile {
						ext_id = nutanix_vm_guest_customization_profile_v2.gc.id
					}
				}
`
	}

	overrideContent := ""

	if c.DiscardAll {
		adminPassBlock := `
										administrator_password {
											discard = true
										}`
		if c.AdminPass != "" {
			adminPassBlock = fmt.Sprintf(`
										administrator_password {
											value = "%s"
										}`, c.AdminPass)
		}
		overrideContent = fmt.Sprintf(`
						sysprep_config {
							customization {
								sysprep_params {
									general_settings {%s
										computer_name {
											discard = true
										}
										timezone {
											discard = true
										}
										registered_organization {
											discard = true
										}
										registered_owner {
											discard = true
										}
									}
									locale_settings {
										ui_language {
											discard = true
										}
										system_locale {
											discard = true
										}
										user_locale {
											discard = true
										}
									}
									workgroup_or_domain_info {
										discard = true
									}
								}
							}
						}`, adminPassBlock)
	} else if c.AnswerFileXML != "" {
		overrideContent = fmt.Sprintf(`
						sysprep_config {
							customization {
								answer_file {
									unattend_xml = "%s"
								}
							}
						}`, c.AnswerFileXML)
	} else {
		generalBlock := ""
		hasGeneral := c.AdminPass != "" || c.ComputerName != "" || c.Timezone != "" || c.Org != "" || c.Owner != "" || c.WindowsProductKey != "" || c.LogonCount > 0
		if hasGeneral {
			parts := ""
			if c.AdminPass != "" {
				parts += fmt.Sprintf(`
										administrator_password {
											value = "%s"
										}`, c.AdminPass)
			}
			if c.ComputerName != "" {
				parts += fmt.Sprintf(`
										computer_name {
											value = "%s"
										}`, c.ComputerName)
			}
			if c.Timezone != "" {
				parts += fmt.Sprintf(`
										timezone {
											value = "%s"
										}`, c.Timezone)
			}
			if c.Org != "" {
				parts += fmt.Sprintf(`
										registered_organization {
											value = "%s"
										}`, c.Org)
			}
			if c.Owner != "" {
				parts += fmt.Sprintf(`
										registered_owner {
											value = "%s"
										}`, c.Owner)
			}
			if c.WindowsProductKey != "" {
				parts += fmt.Sprintf(`
										windows_product_key {
											value = "%s"
										}`, c.WindowsProductKey)
			}
			if c.LogonCount > 0 {
				parts += fmt.Sprintf(`
										auto_logon_settings {
											logon_count = %d
										}`, c.LogonCount)
			}
			generalBlock = fmt.Sprintf(`
									general_settings {%s
									}`, parts)
		}

		firstLogonBlock := ""
		if len(c.FirstLogonCommands) > 0 {
			cmds := ""
			for _, cmd := range c.FirstLogonCommands {
				cmds += fmt.Sprintf(`"%s", `, cmd)
			}
			firstLogonBlock = fmt.Sprintf(`
									first_logon_commands = [%s]`, cmds)
		}

		localeBlock := ""
		if c.UILang != "" || c.SysLocale != "" || c.UserLocale != "" {
			parts := ""
			if c.UILang != "" {
				parts += fmt.Sprintf(`
										ui_language {
											value = "%s"
										}`, c.UILang)
			}
			if c.SysLocale != "" {
				parts += fmt.Sprintf(`
										system_locale {
											value = "%s"
										}`, c.SysLocale)
			}
			if c.UserLocale != "" {
				parts += fmt.Sprintf(`
										user_locale {
											value = "%s"
										}`, c.UserLocale)
			}
			localeBlock = fmt.Sprintf(`
									locale_settings {%s
									}`, parts)
		}

		networkBlock := ""
		if c.DNSRef != "" || c.IPValue != "" {
			dnsBlock := ""
			if c.DNSRef != "" {
				dnsBlock = fmt.Sprintf(`
											dns_config {
												preferred_dns_server_address = local.%s.dns
											}`, c.DNSRef)
			}
			ipBlock := ""
			if c.IPValue != "" && c.WithSubnet {
				ipBlock = fmt.Sprintf(`
											ipv4_config {
												ip_address {
													value         = %s
													prefix_length = local.gc_subnet.prefix_length
												}
												default_gateways = [local.gc_subnet.gateway_ip]
											}`, c.IPValue)
			}
			networkBlock = fmt.Sprintf(`
									network_settings {
										nic_config_list {%s%s
										}
									}`, dnsBlock, ipBlock)
		}

		wgOrDomainBlock := ""
		if c.Workgroup != "" {
			wgOrDomainBlock = fmt.Sprintf(`
									workgroup_or_domain_info {
										workgroup {
											name = "%s"
										}
									}`, c.Workgroup)
		} else if c.DomainRef != "" {
			wgOrDomainBlock = fmt.Sprintf(`
									workgroup_or_domain_info {
										domain_settings {
											credentials {
												domain_name = local.%s.name
												password    = local.%s.password
												username    = local.%s.username
											}
										}
									}`, c.DomainRef, c.DomainRef, c.DomainRef)
		}

		overrideContent = fmt.Sprintf(`
						sysprep_config {
							customization {
								sysprep_params {%s%s%s%s%s
								}
							}
						}`, firstLogonBlock, generalBlock, localeBlock, networkBlock, wgOrDomainBlock)
	}

	return fmt.Sprintf(`
				guest_customization_profile_config {
					profile {
						ext_id = nutanix_vm_guest_customization_profile_v2.gc.id
					}
					config_override_spec {%s
					}
				}
`, overrideContent)
}

func buildAnswerFileXML(workgroup, adminPass string) string {
	return `<?xml version='1.0' encoding='utf-8'?>` +
		`<unattend xmlns='urn:schemas-microsoft-com:unattend'>` +
		`<settings pass='specialize'>` +
		`<component name='Microsoft-Windows-UnattendedJoin' processorArchitecture='amd64' ` +
		`publicKeyToken='31bf3856ad364e35' language='neutral' versionScope='nonSxS'>` +
		`<Identification><JoinWorkgroup>` + workgroup + `</JoinWorkgroup></Identification>` +
		`</component>` +
		`</settings>` +
		`<settings pass='oobeSystem'>` +
		`<component name='Microsoft-Windows-Shell-Setup' processorArchitecture='amd64' ` +
		`publicKeyToken='31bf3856ad364e35' language='neutral' versionScope='nonSxS'>` +
		`<OOBE><HideEULAPage>true</HideEULAPage><SkipMachineOOBE>true</SkipMachineOOBE></OOBE>` +
		`<UserAccounts><AdministratorPassword>` +
		`<Value>` + adminPass + `</Value><PlainText>true</PlainText>` +
		`</AdministratorPassword></UserAccounts>` +
		`</component>` +
		`</settings>` +
		`</unattend>`
}

// ============================================================================
// Case 1: GC profile deploy without override
// ============================================================================

func TestAccV2NutanixTemplateDeployResource_GcProfileCase1_NoOverride(t *testing.T) {
	r := acctest.RandInt() % 100000
	templateName := fmt.Sprintf("gc1tmpl%d", r)
	gcProfileName := fmt.Sprintf("gc1prof%d", r)
	deployName := fmt.Sprintf("gc1dep%d", r)
	cloneName := fmt.Sprintf("gc1cln%d", r)

	profileHCL := buildGcDeployProfileHCL(gcDeployProfileCfg{
		Name:       gcProfileName,
		Desc:       "GC profile case 1 - no override",
		Org:        "Nonsense Org",
		Owner:      "Nonsense Owner",
		AdminPass:  "Noice/4u",
		Timezone:   "India Standard Time",
		UseVMName:  true,
		UILang:     "en-US",
		SysLocale:  "en-US",
		UserLocale: "en-US",
		DomainRef:  "primary_ad",
		UseDHCP:    true,
		DNSRef:     "primary_ad",
	})

	overrideHCL := buildOverrideHCL(gcOverrideCfg{ProfileOnly: true})

	config := testGcProfileBaseInfra() +
		profileHCL +
		testGcTemplateAndDeploy(templateName, "GC case 1 template", deployName, cloneName, overrideHCL, overrideHCL)

	// domain + dns_server come from the primary_ad config the GC profile joins
	// with (DomainRef/DNSRef = "primary_ad"), so assert against that config rather
	// than hardcoding -- otherwise a config change (e.g. qa.nutanix.com ->
	// systest.nutanix.com) breaks the check even though the join succeeded.
	primaryAD := testVars.Iam.DirectoryServicesMain.PrimaryAD
	expectedDeploy := map[string]string{
		"timezone":         "India Standard Time",
		"computer_name":    deployName,
		"domain":           primaryAD.Name,
		"system_locale":    "en-US",
		"ui_language":      "en-US",
		"registered_org":   "Nonsense Org",
		"registered_owner": "Nonsense Owner",
		"dns_server":       primaryAD.DNS,
	}
	expectedClone := map[string]string{
		"timezone":         "India Standard Time",
		"computer_name":    cloneName,
		"domain":           primaryAD.Name,
		"system_locale":    "en-US",
		"ui_language":      "en-US",
		"registered_org":   "Nonsense Org",
		"registered_owner": "Nonsense Owner",
		"dns_server":       primaryAD.DNS,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckGcDeployVMsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameGcProfile, "name", gcProfileName),
					resource.TestCheckResourceAttr(resourceNameGcDeploy, "number_of_vms", "1"),
					resource.TestCheckResourceAttrSet(resourceNameGcDeploy,
						"override_vm_config_map.0.guest_customization_profile_config.0.profile.0.ext_id"),
					testCheckVMSettings(resourceNameGcDeploy, "Administrator", "Noice/4u", expectedDeploy),
					testCheckVMSettings(resourceNameGcClone, "Administrator", "Noice/4u", expectedClone),
				),
			},
		},
	})
}

// ============================================================================
// Case 2: Full override of all GC profile settings
// ============================================================================

func TestAccV2NutanixTemplateDeployResource_GcProfileCase2_FullOverride(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	r := acctest.RandInt() % 100000
	templateName := fmt.Sprintf("gc2tmpl%d", r)
	gcProfileName := fmt.Sprintf("gc2prof%d", r)
	deployName := fmt.Sprintf("gc2dep%d", r)
	cloneName := fmt.Sprintf("gc2cln%d", r)

	firstLogonCmd := fmt.Sprintf(`powershell -Command New-Item -Path C:\\tf_marker_%d.txt -ItemType File -Force`, r)

	profileHCL := buildGcDeployProfileHCL(gcDeployProfileCfg{
		Name:                gcProfileName,
		Desc:                "GC profile case 2 - full override",
		Org:                 "Nonsense Org",
		Owner:               "Nonsense Owner",
		AdminPass:           "Noice/4u",
		Timezone:            "Pacific Standard Time",
		MustProvideCompName: true,
		UILang:              "fr-FR",
		SysLocale:           "fr-FR",
		UserLocale:          "fr-FR",
		DomainRef:           "primary_ad",
		MustProvideIP:       true,
		DNSRef:              "primary_ad",
		FirstLogonCommands:  []string{`powershell -Command echo profile_cmd > C:\\gc_profile_cmd.txt`},
		LogonCount:          1,
		WindowsProductKey:   "VK7JG-NPHTM-C97JM-9MPGT-3V66T",
	})

	baseOverride := gcOverrideCfg{
		AdminPass:          "TEST123",
		Timezone:           "India Standard Time",
		Org:                "Secondtestcase",
		Owner:              "Secondtestcase",
		UILang:             "en-US",
		SysLocale:          "en-US",
		UserLocale:         "en-US",
		DNSRef:             "secondary_ad",
		WithSubnet:         true,
		DomainRef:          "secondary_ad",
		FirstLogonCommands: []string{firstLogonCmd},
		LogonCount:         3,
		WindowsProductKey:  "VK7JG-NPHTM-C97JM-9MPGT-3V66T",
	}

	freeIPs, err := findFreeIPs(testVars.Networking.GcSubnet.StartIP, testVars.Networking.GcSubnet.EndIP, 2)
	if err != nil {
		t.Fatalf("could not find free IPs: %v", err)
	}

	deployCompName := fmt.Sprintf("tfdep%d", r)
	cloneCompName := fmt.Sprintf("tfcln%d", r)

	deployOvr := baseOverride
	deployOvr.ComputerName = deployCompName
	deployOvr.IPValue = fmt.Sprintf(`"%s"`, freeIPs[0])

	cloneOvr := baseOverride
	cloneOvr.ComputerName = cloneCompName
	cloneOvr.IPValue = fmt.Sprintf(`"%s"`, freeIPs[1])

	config := testGcProfileBaseInfra() +
		profileHCL +
		testGcTemplateAndDeploy(templateName, "GC case 2 template", deployName, cloneName,
			buildOverrideHCL(deployOvr), buildOverrideHCL(cloneOvr))

	markerFile := fmt.Sprintf(`C:\tf_marker_%d.txt`, r)

	expectedDeploy := map[string]string{
		"timezone":           "India Standard Time",
		"computer_name":      strings.ToUpper(deployCompName),
		"domain":             "qa.nucalm.io",
		"system_locale":      "en-US",
		"ui_language":        "en-US",
		"registered_org":     "Secondtestcase",
		"registered_owner":   "Secondtestcase",
		"dns_server":         testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS,
		"ip_address":         freeIPs[0],
		"first_logon_marker": markerFile,
		"auto_logon_count":   "3",
	}
	expectedClone := map[string]string{
		"timezone":           "India Standard Time",
		"computer_name":      strings.ToUpper(cloneCompName),
		"domain":             "qa.nucalm.io",
		"system_locale":      "en-US",
		"ui_language":        "en-US",
		"registered_org":     "Secondtestcase",
		"registered_owner":   "Secondtestcase",
		"dns_server":         testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS,
		"ip_address":         freeIPs[1],
		"first_logon_marker": markerFile,
		"auto_logon_count":   "3",
	}

	templateGetDS := fmt.Sprintf(`
		data "nutanix_template_v2" "gc_get" {
			ext_id = nutanix_template_v2.gc_template.id
		}
	`)
	templateListDS := fmt.Sprintf(`
		data "nutanix_templates_v2" "gc_list" {
			filter = "templateName eq '%s'"
			depends_on = [nutanix_template_v2.gc_template]
		}
	`, templateName)

	configWithDS := config + templateGetDS + templateListDS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckGcDeployVMsDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithDS,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameGcProfile, "name", gcProfileName),
					resource.TestCheckResourceAttr(resourceNameGcDeploy, "number_of_vms", "1"),
					resource.TestCheckResourceAttrSet(resourceNameGcDeploy,
						"override_vm_config_map.0.guest_customization_profile_config.0.profile.0.ext_id"),
					testCheckVMSettingsWithIP(resourceNameGcDeploy, "Administrator", "TEST123", freeIPs[0], expectedDeploy),
					testCheckVMSettingsWithIP(resourceNameGcClone, "Administrator", "TEST123", freeIPs[1], expectedClone),

					// Template Get datasource - verify GC profile is in the response
					resource.TestCheckResourceAttrSet("data.nutanix_template_v2.gc_get", "ext_id"),
					resource.TestCheckResourceAttr("data.nutanix_template_v2.gc_get", "template_name", templateName),
					resource.TestCheckResourceAttrSet("data.nutanix_template_v2.gc_get",
						"template_version_spec.0.guest_customization_profile.0.ext_id"),

					// Template List datasource - verify GC profile is in the response
					resource.TestCheckResourceAttrSet("data.nutanix_templates_v2.gc_list", "templates.#"),
					resource.TestCheckResourceAttr("data.nutanix_templates_v2.gc_list",
						"templates.0.template_name", templateName),
					resource.TestCheckResourceAttrSet("data.nutanix_templates_v2.gc_list",
						"templates.0.template_version_spec.0.guest_customization_profile.0.ext_id"),
				),
			},
		},
	})
}

// ============================================================================
// Case 3: Workgroup name override
// ============================================================================

func TestAccV2NutanixTemplateDeployResource_GcProfileCase3_WorkgroupNameOverride(t *testing.T) {
	r := acctest.RandInt() % 100000
	templateName := fmt.Sprintf("gc3tmpl%d", r)
	gcProfileName := fmt.Sprintf("gc3prof%d", r)
	deployName := fmt.Sprintf("gc3dep%d", r)
	cloneName := fmt.Sprintf("gc3cln%d", r)

	wgProfile := fmt.Sprintf("wg%d", r)
	wgOverride := fmt.Sprintf("wgo%d", r)

	profileHCL := buildGcDeployProfileHCL(gcDeployProfileCfg{
		Name:      gcProfileName,
		Desc:      "GC profile case 3 - workgroup name override",
		AdminPass: "Noice/4u",
		Timezone:  "India Standard Time",
		UseVMName: true,
		Workgroup: wgProfile,
	})

	overrideHCL := buildOverrideHCL(gcOverrideCfg{Workgroup: wgOverride})

	config := testGcProfileBaseInfra() +
		profileHCL +
		testGcTemplateAndDeploy(templateName, "GC case 3 template", deployName, cloneName, overrideHCL, overrideHCL)

	expectedDeploy := map[string]string{
		"timezone":      "India Standard Time",
		"computer_name": deployName,
		"workgroup":     wgOverride,
	}
	expectedClone := map[string]string{
		"timezone":      "India Standard Time",
		"computer_name": cloneName,
		"workgroup":     wgOverride,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckGcDeployVMsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameGcProfile, "name", gcProfileName),
					resource.TestCheckResourceAttr(resourceNameGcDeploy, "number_of_vms", "1"),
					testCheckVMSettings(resourceNameGcDeploy, "Administrator", "Noice/4u", expectedDeploy),
					testCheckVMSettings(resourceNameGcClone, "Administrator", "Noice/4u", expectedClone),
				),
			},
		},
	})
}

// ============================================================================
// Case 4: Workgroup to domain switch (implicit via domain_settings)
// ============================================================================

func TestAccV2NutanixTemplateDeployResource_GcProfileCase4_WorkgroupToDomain(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	r := acctest.RandInt() % 100000
	templateName := fmt.Sprintf("gc4tmpl%d", r)
	gcProfileName := fmt.Sprintf("gc4prof%d", r)
	deployName := fmt.Sprintf("gc4dep%d", r)
	cloneName := fmt.Sprintf("gc4cln%d", r)

	profileHCL := buildGcDeployProfileHCL(gcDeployProfileCfg{
		Name:                gcProfileName,
		Desc:                "GC profile case 4 - workgroup to domain switch",
		AdminPass:           "Noice/4u",
		Timezone:            "India Standard Time",
		MustProvideCompName: true,
		Workgroup:           "akhil",
	})

	baseOverride := gcOverrideCfg{
		DNSRef:     "secondary_ad",
		WithSubnet: true,
		DomainRef:  "secondary_ad",
	}

	freeIPs, err := findFreeIPs(testVars.Networking.GcSubnet.StartIP, testVars.Networking.GcSubnet.EndIP, 2)
	if err != nil {
		t.Fatalf("could not find free IPs: %v", err)
	}

	deployCompName := fmt.Sprintf("c4dep%d", r)
	cloneCompName := fmt.Sprintf("c4cln%d", r)

	deployOvr := baseOverride
	deployOvr.ComputerName = deployCompName
	deployOvr.IPValue = fmt.Sprintf(`"%s"`, freeIPs[0])

	cloneOvr := baseOverride
	cloneOvr.ComputerName = cloneCompName
	cloneOvr.IPValue = fmt.Sprintf(`"%s"`, freeIPs[1])

	config := testGcProfileBaseInfra() +
		profileHCL +
		testGcTemplateAndDeploy(templateName, "GC case 4 template", deployName, cloneName,
			buildOverrideHCL(deployOvr), buildOverrideHCL(cloneOvr))

	expectedDeploy := map[string]string{
		"timezone":      "India Standard Time",
		"computer_name": strings.ToUpper(deployCompName),
		"domain":        "qa.nucalm.io",
		"dns_server":    testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS,
		"ip_address":    freeIPs[0],
	}
	expectedClone := map[string]string{
		"timezone":      "India Standard Time",
		"computer_name": strings.ToUpper(cloneCompName),
		"domain":        "qa.nucalm.io",
		"dns_server":    testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS,
		"ip_address":    freeIPs[1],
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckGcDeployVMsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameGcProfile, "name", gcProfileName),
					resource.TestCheckResourceAttr(resourceNameGcDeploy, "number_of_vms", "1"),
					testCheckVMSettingsWithIP(resourceNameGcDeploy, "Administrator", "Noice/4u", freeIPs[0], expectedDeploy),
					testCheckVMSettingsWithIP(resourceNameGcClone, "Administrator", "Noice/4u", freeIPs[1], expectedClone),
				),
			},
		},
	})
}

// ============================================================================
// Case 5: Discard all GC profile settings
// ============================================================================

func TestAccV2NutanixTemplateDeployResource_GcProfileCase5_DiscardAll(t *testing.T) {
	r := acctest.RandInt() % 100000
	templateName := fmt.Sprintf("gc5tmpl%d", r)
	gcProfileName := fmt.Sprintf("gc5prof%d", r)
	deployName := fmt.Sprintf("gc5dep%d", r)
	cloneName := fmt.Sprintf("gc5cln%d", r)

	discardAdminPass := "Discard.123"

	profileHCL := buildGcDeployProfileHCL(gcDeployProfileCfg{
		Name:       gcProfileName,
		Desc:       "GC profile case 5 - discard all",
		Org:        "Nonsense Org",
		Owner:      "Nonsense Owner",
		AdminPass:  "Noice/4u",
		Timezone:   "India Standard Time",
		UseVMName:  true,
		UILang:     "fr-FR",
		SysLocale:  "fr-FR",
		UserLocale: "fr-FR",
		Workgroup:  "akhil",
		UseDHCP:    true,
		DNSRef:     "primary_ad",
	})

	overrideHCL := buildOverrideHCL(gcOverrideCfg{DiscardAll: true, AdminPass: discardAdminPass})

	config := testGcProfileBaseInfra() +
		profileHCL +
		testGcTemplateAndDeploy(templateName, "GC case 5 template", deployName, cloneName, overrideHCL, overrideHCL)

	expectedDeploy := map[string]string{
		"timezone":      "UTC",
		"system_locale": "en-US",
	}
	expectedClone := map[string]string{
		"timezone":      "UTC",
		"system_locale": "en-US",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckGcDeployVMsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameGcProfile, "name", gcProfileName),
					resource.TestCheckResourceAttr(resourceNameGcDeploy, "number_of_vms", "1"),
					testCheckVMSettings(resourceNameGcDeploy, "Administrator", discardAdminPass, expectedDeploy),
					testCheckVMSettings(resourceNameGcClone, "Administrator", discardAdminPass, expectedClone),
				),
			},
		},
	})
}

// ============================================================================
// Case 6: Answer file with workgroup – no override
// ============================================================================

func TestAccV2NutanixTemplateDeployResource_GcProfileCase6_AnswerFileWorkgroup(t *testing.T) {
	r := acctest.RandInt() % 100000
	templateName := fmt.Sprintf("gc6tmpl%d", r)
	gcProfileName := fmt.Sprintf("gc6prof%d", r)
	deployName := fmt.Sprintf("gc6dep%d", r)
	cloneName := fmt.Sprintf("gc6cln%d", r)

	answerWG := fmt.Sprintf("awg%d", r)

	profileHCL := buildGcDeployProfileHCL(gcDeployProfileCfg{
		Name:          gcProfileName,
		Desc:          "GC profile case 6 - answer file workgroup",
		AnswerFileXML: buildAnswerFileXML(answerWG, "Noice/4u"),
	})

	overrideHCL := buildOverrideHCL(gcOverrideCfg{ProfileOnly: true})

	config := testGcProfileBaseInfra() +
		profileHCL +
		testGcTemplateAndDeploy(templateName, "GC case 6 template", deployName, cloneName, overrideHCL, overrideHCL)

	expectedDeploy := map[string]string{
		"workgroup": answerWG,
	}
	expectedClone := map[string]string{
		"workgroup": answerWG,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckGcDeployVMsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameGcProfile, "name", gcProfileName),
					resource.TestCheckResourceAttr(resourceNameGcProfile,
						"config.0.sysprep_config.0.customization.0.answer_file.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameGcProfile,
						"config.0.sysprep_config.0.customization.0.answer_file.0.unattend_xml"),
					resource.TestCheckResourceAttr(resourceNameGcDeploy, "number_of_vms", "1"),
					resource.TestCheckResourceAttrSet(resourceNameGcDeploy,
						"override_vm_config_map.0.guest_customization_profile_config.0.profile.0.ext_id"),
					testCheckVMSettings(resourceNameGcDeploy, "Administrator", "Noice/4u", expectedDeploy),
					testCheckVMSettings(resourceNameGcClone, "Administrator", "Noice/4u", expectedClone),
				),
			},
		},
	})
}

// ============================================================================
// Case 7: Answer file with workgroup – override answer file
// ============================================================================

func TestAccV2NutanixTemplateDeployResource_GcProfileCase7_AnswerFileWorkgroupOverride(t *testing.T) {
	r := acctest.RandInt() % 100000
	templateName := fmt.Sprintf("gc7tmpl%d", r)
	gcProfileName := fmt.Sprintf("gc7prof%d", r)
	deployName := fmt.Sprintf("gc7dep%d", r)
	cloneName := fmt.Sprintf("gc7cln%d", r)

	answerWG := fmt.Sprintf("awg%d", r)
	overrideWG := fmt.Sprintf("owg%d", r)

	profileHCL := buildGcDeployProfileHCL(gcDeployProfileCfg{
		Name:          gcProfileName,
		Desc:          "GC profile case 7 - answer file workgroup override",
		AnswerFileXML: buildAnswerFileXML(answerWG, "Noice/4u"),
	})

	overrideHCL := buildOverrideHCL(gcOverrideCfg{
		AnswerFileXML: buildAnswerFileXML(overrideWG, "Noice/4u"),
	})

	config := testGcProfileBaseInfra() +
		profileHCL +
		testGcTemplateAndDeploy(templateName, "GC case 7 template", deployName, cloneName, overrideHCL, overrideHCL)

	expectedDeploy := map[string]string{
		"workgroup": overrideWG,
	}
	expectedClone := map[string]string{
		"workgroup": overrideWG,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckGcDeployVMsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameGcProfile, "name", gcProfileName),
					resource.TestCheckResourceAttr(resourceNameGcProfile,
						"config.0.sysprep_config.0.customization.0.answer_file.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameGcProfile,
						"config.0.sysprep_config.0.customization.0.answer_file.0.unattend_xml"),
					resource.TestCheckResourceAttr(resourceNameGcDeploy, "number_of_vms", "1"),
					resource.TestCheckResourceAttrSet(resourceNameGcDeploy,
						"override_vm_config_map.0.guest_customization_profile_config.0.profile.0.ext_id"),
					testCheckVMSettings(resourceNameGcDeploy, "Administrator", "Noice/4u", expectedDeploy),
					testCheckVMSettings(resourceNameGcClone, "Administrator", "Noice/4u", expectedClone),
				),
			},
		},
	})
}
