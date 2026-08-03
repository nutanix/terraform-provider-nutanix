terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# =============================================================================
# Shared Data Sources
# =============================================================================

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_subnets_v2" "subnet" {
  filter = "name eq '${var.subnet_name}'"
}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  subnet_ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
}

# =============================================================================
# Example 1: GC Profile — use_vm_name with Workgroup (simplest case)
# =============================================================================
# The deployed VM's computer name will match its VM name.
# Joined to a workgroup with DHCP networking.

resource "nutanix_vm_guest_customization_profile_v2" "workgroup_basic" {
  name        = "gc-workgroup-basic"
  description = "Basic workgroup profile using VM name as computer name"
  config {
    sysprep_config {
      customization {
        sysprep_params {
          general_settings {
            computer_name {
              use_vm_name = true
            }
            administrator_password = var.admin_password
            timezone               = "Pacific Standard Time"
          }
          locale_settings {
            ui_language   = "en-US"
            system_locale = "en-US"
            user_locale   = "en-US"
          }
          network_settings {
            nic_config_list {
              ipv4_config {
                use_dhcp = true
              }
            }
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
  lifecycle {
    ignore_changes = [
      config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].administrator_password
    ]
  }
}

# =============================================================================
# Example 2: GC Profile — must_provide_during_deployment for computer name
# =============================================================================
# Computer name and IP must be supplied at deployment time via overrides.
# Profile joins a domain with DNS, locale, and org/owner settings.

resource "nutanix_vm_guest_customization_profile_v2" "must_provide" {
  name        = "gc-must-provide-fields"
  description = "Profile requiring computer name and IP at deploy time"
  config {
    sysprep_config {
      customization {
        sysprep_params {
          general_settings {
            computer_name {
              must_provide_during_deployment = true
            }
            administrator_password  = var.admin_password
            timezone                = "Pacific Standard Time"
            registered_organization = "My Organization"
            registered_owner        = "IT Admin"
          }
          locale_settings {
            ui_language   = "fr-FR"
            system_locale = "fr-FR"
            user_locale   = "fr-FR"
          }
          network_settings {
            nic_config_list {
              dns_config {
                preferred_dns_server_address = var.domain_dns_server
              }
              ipv4_config {
                must_provide_during_deployment = true
              }
            }
          }
          workgroup_or_domain_info {
            domain_settings {
              credentials {
                domain_name = var.domain_name
                username    = var.domain_username
                password    = var.domain_password
              }
            }
          }
        }
      }
    }
  }
  lifecycle {
    ignore_changes = [
      config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].administrator_password,
      config[0].sysprep_config[0].customization[0].sysprep_params[0].workgroup_or_domain_info[0].domain_settings[0].credentials[0].password
    ]
  }
}

# =============================================================================
# Example 3: GC Profile — Domain join with DHCP
# =============================================================================
# VM is joined to a domain. Computer name is derived from VM name.
# Uses DHCP for IP assignment and configures a preferred DNS server.

resource "nutanix_vm_guest_customization_profile_v2" "domain_dhcp" {
  name        = "gc-domain-dhcp"
  description = "Domain-joined profile with DHCP and DNS"
  config {
    sysprep_config {
      customization {
        sysprep_params {
          general_settings {
            computer_name {
              use_vm_name = true
            }
            administrator_password  = var.admin_password
            timezone                = "India Standard Time"
            registered_organization = "Nutanix Inc"
            registered_owner        = "Cloud Admin"
          }
          locale_settings {
            ui_language   = "en-US"
            system_locale = "en-US"
            user_locale   = "en-US"
          }
          network_settings {
            nic_config_list {
              dns_config {
                preferred_dns_server_address = var.domain_dns_server
              }
              ipv4_config {
                use_dhcp = true
              }
            }
          }
          workgroup_or_domain_info {
            domain_settings {
              credentials {
                domain_name = var.domain_name
                username    = var.domain_username
                password    = var.domain_password
              }
            }
          }
        }
      }
    }
  }
  lifecycle {
    ignore_changes = [
      config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].administrator_password,
      config[0].sysprep_config[0].customization[0].sysprep_params[0].workgroup_or_domain_info[0].domain_settings[0].credentials[0].password
    ]
  }
}

# =============================================================================
# Example 4: GC Profile — Answer File (unattend.xml)
# =============================================================================
# Instead of sysprep_params, supply a complete unattend.xml answer file.
# Useful when you need full control over the sysprep XML or have a
# pre-existing answer file.

resource "nutanix_vm_guest_customization_profile_v2" "answer_file" {
  name        = "gc-answer-file"
  description = "Profile using unattend.xml answer file"
  config {
    sysprep_config {
      customization {
        answer_file {
          unattend_xml = (templatefile("${path.module}/unattend.xml", {
            admin_password = var.admin_password
            workgroup_name = "MYWORKGROUP"
          }))
        }
      }
    }
  }
}

# =============================================================================
# Example 5: GC Profile with first_logon_commands, auto_logon, product_key
# =============================================================================
# Demonstrates advanced settings: commands that run at first logon,
# auto-logon configuration, and Windows product key activation.

resource "nutanix_vm_guest_customization_profile_v2" "advanced" {
  name        = "gc-advanced-settings"
  description = "Profile with first logon commands, auto logon, and product key"
  config {
    sysprep_config {
      customization {
        sysprep_params {
          general_settings {
            computer_name {
              must_provide_during_deployment = true
            }
            administrator_password  = var.admin_password
            timezone                = "Pacific Standard Time"
            registered_organization = "Nutanix Inc"
            registered_owner        = "Cloud Admin"
            auto_logon_settings {
              logon_count = 1
            }
            windows_product_key = var.windows_product_key
          }
          first_logon_commands = [
            "powershell -Command New-Item -Path C:\\setup_complete.txt -ItemType File -Force",
            "powershell -Command Set-ExecutionPolicy RemoteSigned -Force",
          ]
          locale_settings {
            ui_language   = "en-US"
            system_locale = "en-US"
            user_locale   = "en-US"
          }
          network_settings {
            nic_config_list {
              dns_config {
                preferred_dns_server_address = var.domain_dns_server
              }
              ipv4_config {
                must_provide_during_deployment = true
              }
            }
          }
          workgroup_or_domain_info {
            domain_settings {
              credentials {
                domain_name = var.domain_name
                username    = var.domain_username
                password    = var.domain_password
              }
            }
          }
        }
      }
    }
  }
  lifecycle {
    ignore_changes = [
      config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].administrator_password,
      config[0].sysprep_config[0].customization[0].sysprep_params[0].general_settings[0].windows_product_key,
      config[0].sysprep_config[0].customization[0].sysprep_params[0].workgroup_or_domain_info[0].domain_settings[0].credentials[0].password
    ]
  }
}

# =============================================================================
# Datasource — Read a single GC profile
# =============================================================================

data "nutanix_vm_guest_customization_profile_v2" "read_one" {
  ext_id = nutanix_vm_guest_customization_profile_v2.workgroup_basic.id
}

# =============================================================================
# Datasource — List all GC profiles (with optional filter)
# =============================================================================

data "nutanix_vm_guest_customization_profiles_v2" "list_all" {
  depends_on = [nutanix_vm_guest_customization_profile_v2.workgroup_basic]
}

# =============================================================================
# Template — Create from a VM with GC profile attached
# =============================================================================

resource "nutanix_template_v2" "from_vm" {
  template_name        = "template-with-gc-profile"
  template_description = "Template created from a Windows VM with GC profile"
  template_version_spec {
    version_source {
      template_vm_reference {
        ext_id = var.source_vm_uuid
        guest_customization_profile {
          ext_id = nutanix_vm_guest_customization_profile_v2.must_provide.id
        }
      }
    }
    guest_customization_profile {
      ext_id = nutanix_vm_guest_customization_profile_v2.must_provide.id
    }
  }
  lifecycle {
    ignore_changes = [
      template_version_spec.0.version_name,
      template_version_spec.0.version_description,
      template_version_spec.0.version_source
    ]
  }
  depends_on = [nutanix_vm_guest_customization_profile_v2.must_provide]
}

# =============================================================================
# Template Deploy — Override all GC settings
# =============================================================================
# When the GC profile uses must_provide_during_deployment, you MUST supply
# the values at deploy time. You can also override any other field.

resource "nutanix_deploy_templates_v2" "full_override" {
  ext_id            = nutanix_template_v2.from_vm.id
  number_of_vms     = 1
  cluster_reference = local.cluster_ext_id
  override_vm_config_map {
    name                 = "deployed-vm-override"
    memory_size_bytes    = 4 * 1024 * 1024 * 1024
    num_cores_per_socket = 1
    num_sockets          = 2
    num_threads_per_core = 1
    guest_customization_profile_config {
      profile {
        ext_id = nutanix_vm_guest_customization_profile_v2.must_provide.id
      }
      config_override_spec {
        sysprep_config {
          customization {
            sysprep_params {
              general_settings {
                administrator_password {
                  value = var.override_admin_password
                }
                computer_name {
                  value = "DEPLOYEDVM01"
                }
                timezone {
                  value = "India Standard Time"
                }
                registered_organization {
                  value = "Override Org"
                }
                registered_owner {
                  value = "Override Owner"
                }
              }
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
              network_settings {
                nic_config_list {
                  dns_config {
                    preferred_dns_server_address = var.override_dns_server
                  }
                  ipv4_config {
                    ip_address {
                      value         = var.deploy_ip_address
                      prefix_length = var.subnet_prefix_length
                    }
                    default_gateways = [var.subnet_gateway]
                  }
                }
              }
              workgroup_or_domain_info {
                domain_settings {
                  credentials {
                    domain_name = var.override_domain_name
                    username    = var.override_domain_username
                    password    = var.override_domain_password
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  depends_on = [nutanix_template_v2.from_vm]
}

# =============================================================================
# Template Deploy — Profile only (no overrides)
# =============================================================================
# Use the GC profile as-is without overriding any field. The VM inherits
# all settings from the profile (computer name = VM name, DHCP, etc.).

resource "nutanix_deploy_templates_v2" "profile_only" {
  ext_id            = nutanix_template_v2.from_vm.id
  number_of_vms     = 1
  cluster_reference = local.cluster_ext_id
  override_vm_config_map {
    name                 = "deployed-vm-no-override"
    memory_size_bytes    = 4 * 1024 * 1024 * 1024
    num_cores_per_socket = 1
    num_sockets          = 2
    guest_customization_profile_config {
      profile {
        ext_id = nutanix_vm_guest_customization_profile_v2.domain_dhcp.id
      }
    }
  }
  depends_on = [nutanix_template_v2.from_vm]
}

# =============================================================================
# Template Deploy — Discard all GC settings
# =============================================================================
# Discards all GC profile settings so the VM boots with OS defaults.
# The administrator_password is explicitly set to allow login.

resource "nutanix_deploy_templates_v2" "discard_all" {
  ext_id            = nutanix_template_v2.from_vm.id
  number_of_vms     = 1
  cluster_reference = local.cluster_ext_id
  override_vm_config_map {
    name                 = "deployed-vm-discard"
    memory_size_bytes    = 4 * 1024 * 1024 * 1024
    num_cores_per_socket = 1
    num_sockets          = 1
    guest_customization_profile_config {
      profile {
        ext_id = nutanix_vm_guest_customization_profile_v2.must_provide.id
      }
      config_override_spec {
        sysprep_config {
          customization {
            sysprep_params {
              general_settings {
                administrator_password {
                  value = var.override_admin_password
                }
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
        }
      }
    }
  }
  depends_on = [nutanix_template_v2.from_vm]
}

# =============================================================================
# Template Deploy — Override with Answer File
# =============================================================================
# Override the GC profile's sysprep_params with an answer file instead.

resource "nutanix_deploy_templates_v2" "answer_file_override" {
  ext_id            = nutanix_template_v2.from_vm.id
  number_of_vms     = 1
  cluster_reference = local.cluster_ext_id
  override_vm_config_map {
    name                 = "deployed-vm-answerfile"
    memory_size_bytes    = 4 * 1024 * 1024 * 1024
    num_cores_per_socket = 1
    num_sockets          = 1
    guest_customization_profile_config {
      profile {
        ext_id = nutanix_vm_guest_customization_profile_v2.must_provide.id
      }
      config_override_spec {
        sysprep_config {
          customization {
            answer_file {
              unattend_xml = base64encode(templatefile("${path.module}/unattend.xml", {
                admin_password = var.override_admin_password
                workgroup_name = "OVERRIDE-WG"
              }))
            }
          }
        }
      }
    }
  }
  depends_on = [nutanix_template_v2.from_vm]
}

# =============================================================================
# VM Clone — With GC profile override (workgroup to domain switch)
# =============================================================================
# Clone an existing VM and apply a GC profile that switches from workgroup
# to domain join, overriding computer name and network settings.

resource "nutanix_vm_clone_v2" "with_gc_override" {
  vm_ext_id            = var.source_vm_uuid
  name                 = "cloned-vm-domain-override"
  memory_size_bytes    = 4 * 1024 * 1024 * 1024
  num_cores_per_socket = 1
  num_sockets          = 2
  guest_customization_profile_config {
    profile {
      ext_id = nutanix_vm_guest_customization_profile_v2.must_provide.id
    }
    config_override_spec {
      sysprep_config {
        customization {
          sysprep_params {
            general_settings {
              administrator_password {
                value = var.override_admin_password
              }
              computer_name {
                value = "CLONEDVM01"
              }
              timezone {
                value = "India Standard Time"
              }
              registered_organization {
                value = "Clone Org"
              }
              registered_owner {
                value = "Clone Owner"
              }
            }
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
            network_settings {
              nic_config_list {
                dns_config {
                  preferred_dns_server_address = var.override_dns_server
                }
                ipv4_config {
                  ip_address {
                    value         = var.clone_ip_address
                    prefix_length = var.subnet_prefix_length
                  }
                  default_gateways = [var.subnet_gateway]
                }
              }
            }
            workgroup_or_domain_info {
              domain_settings {
                credentials {
                  domain_name = var.override_domain_name
                  username    = var.override_domain_username
                  password    = var.override_domain_password
                }
              }
            }
          }
        }
      }
    }
  }
  depends_on = [nutanix_vm_guest_customization_profile_v2.must_provide]
}

# =============================================================================
# VM Clone — Profile only (no overrides)
# =============================================================================
# Clone a VM using a GC profile as-is. The cloned VM gets all settings
# from the profile without any overrides.

resource "nutanix_vm_clone_v2" "profile_only" {
  vm_ext_id            = var.source_vm_uuid
  name                 = "cloned-vm-no-override"
  memory_size_bytes    = 4 * 1024 * 1024 * 1024
  num_cores_per_socket = 1
  num_sockets          = 2
  guest_customization_profile_config {
    profile {
      ext_id = nutanix_vm_guest_customization_profile_v2.domain_dhcp.id
    }
  }
  depends_on = [nutanix_vm_guest_customization_profile_v2.domain_dhcp]
}

# =============================================================================
# VM Clone — With Answer File override
# =============================================================================

resource "nutanix_vm_clone_v2" "answer_file_override" {
  vm_ext_id            = var.source_vm_uuid
  name                 = "cloned-vm-answerfile"
  memory_size_bytes    = 4 * 1024 * 1024 * 1024
  num_cores_per_socket = 1
  num_sockets          = 1
  guest_customization_profile_config {
    profile {
      ext_id = nutanix_vm_guest_customization_profile_v2.must_provide.id
    }
    config_override_spec {
      sysprep_config {
        customization {
          answer_file {
            unattend_xml = base64encode(templatefile("${path.module}/unattend.xml", {
              admin_password = var.override_admin_password
              workgroup_name = "CLONE-WG"
            }))
          }
        }
      }
    }
  }
  depends_on = [nutanix_vm_guest_customization_profile_v2.must_provide]
}

# =============================================================================
# Template Datasources — Read template details
# =============================================================================

data "nutanix_template_v2" "get_template" {
  ext_id = nutanix_template_v2.from_vm.id
}

data "nutanix_templates_v2" "list_templates" {
  filter     = "templateName eq '${nutanix_template_v2.from_vm.template_name}'"
  depends_on = [nutanix_template_v2.from_vm]
}
