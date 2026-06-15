terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
    }
  }
}

provider "nutanix" {
  ndb_endpoint = var.ndb_endpoint
  ndb_username = var.ndb_username
  ndb_password = var.ndb_password
  insecure     = true
}

variable "ndb_endpoint" {
  type = string
}

variable "ndb_username" {
  type = string
}

variable "ndb_password" {
  type      = string
  sensitive = true
}

variable "pe_name" {
  type = string
}

variable "pe_cluster_ip" {
  type = string
}

variable "pe_username" {
  type = string
}

variable "pe_password" {
  type      = string
  sensitive = true
}

variable "pc_ip" {
  type    = string
  default = ""
}

variable "pc_username" {
  type    = string
  default = ""
}

variable "pc_password" {
  type      = string
  sensitive = true
  default   = ""
}

variable "storage_container" {
  type = string
  default = ""
}

variable "dns_servers" {
  type = list(string)
}

variable "ntp_servers" {
  type = list(string)
}

resource "nutanix_ndb_onboarding" "wizard" {
  # Full wizard mode: Step 1 (optional) -> Step 6.
  enable_full_onboarding = true
  # auto: use user-provided values when set; otherwise pick discovered options.
  # strict: fail if user-provided selection is not found in discovered options.
  selection_mode = "auto"

  # Step 1 (optional): provide Prism Central details when needed.
  dynamic "prism_central_info" {
    for_each = var.pc_ip != "" ? [1] : []
    content {
      ip_address = var.pc_ip
      username   = var.pc_username
      password   = var.pc_password
    }
  }

  # Step 2 (required): Prism Element.
  prism_element_info {
    name       = var.pe_name
    cluster_ip = var.pe_cluster_ip
    username   = var.pe_username
    password   = var.pe_password
  }

  # Step 3: optional DNS/NTP override.
  # If omitted, provider uses existing NDB DNS/NTP values.
  # Avoid hardcoded environment-specific defaults in shared examples.
  ndb_config {
    dns_servers = var.dns_servers
    ntp_servers = var.ntp_servers
    timezone    = "UTC"
  }

  # Step 4: optional storage override.
  # If empty/missing in auto mode, provider selects first discovered container.
  storage {
    container_name = var.storage_container
  }

  # Step 5: network selection (choose existing network name or skip).
  network_details {
    skip                  = false
    existing_network_name = "vlan0"
  }

  # Step 6: setup trigger and wait timeout.
  setup {
    trigger         = true
    timeout_minutes = 120
  }
}

output "onboarding_effective_values" {
  value = {
    cluster_id                  = nutanix_ndb_onboarding.wizard.cluster_id
    effective_storage_container = nutanix_ndb_onboarding.wizard.effective_storage_container
    effective_dns_servers       = nutanix_ndb_onboarding.wizard.effective_dns_servers
    effective_ntp_servers       = nutanix_ndb_onboarding.wizard.effective_ntp_servers
    effective_network_name      = nutanix_ndb_onboarding.wizard.effective_network_name
  }
}

output "onboarding_available_options" {
  value = {
    storage_containers = nutanix_ndb_onboarding.wizard.available_storage_containers
    dns_servers        = nutanix_ndb_onboarding.wizard.available_dns_servers
    ntp_servers        = nutanix_ndb_onboarding.wizard.available_ntp_servers
    network_names      = nutanix_ndb_onboarding.wizard.available_network_names
  }
}
