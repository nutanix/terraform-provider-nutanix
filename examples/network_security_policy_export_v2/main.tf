terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Categories referenced by the policies below.
resource "nutanix_category_v2" "cat1" {
  key   = "tf-example-cat-1"
  value = "tf-example-cat-1-value"
}

resource "nutanix_category_v2" "cat2" {
  key   = "tf-example-cat-2"
  value = "tf-example-cat-2-value"
}

resource "nutanix_category_v2" "cat3" {
  key   = "tf-example-cat-3"
  value = "tf-example-cat-3-value"
}

# A network security policy with an intra-group rule.
resource "nutanix_network_security_policy_v2" "policy1" {
  name        = "tf-example-nsp-1"
  description = "Network Security Policy example with an intra group rule"
  state       = "MONITOR"
  type        = "APPLICATION"
  scope       = "ALL_VLAN"
  rules {
    description = "intra group rule"
    type        = "INTRA_GROUP"
    spec {
      intra_entity_group_rule_spec {
        secured_group_category_references = [
          nutanix_category_v2.cat1.id,
        ]
        secured_group_action = "ALLOW"
      }
    }
  }
  is_hitlog_enabled = false
  lifecycle {
    ignore_changes = [rules]
  }
}

# A network security policy with an application rule.
resource "nutanix_network_security_policy_v2" "policy2" {
  name        = "tf-example-nsp-2"
  description = "Network Security Policy example with an application rule"
  state       = "MONITOR"
  type        = "APPLICATION"
  scope       = "ALL_VLAN"
  rules {
    description = "application rule"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          nutanix_category_v2.cat2.id,
        ]
        src_category_references = [
          nutanix_category_v2.cat3.id,
        ]
        is_all_protocol_allowed = true
      }
    }
  }
  lifecycle {
    ignore_changes = [rules]
  }
}

# project
resource "nutanix_project_v2" "example_project" {
  name = "tf-example-project"
  project_id = "tf-example-project"
  description = "Example project"
}

# -----------------------------------------------------------------------------
# EXAMPLES: EXPORT OF POLICY WITHOUT PROJECT SCOPE
# -----------------------------------------------------------------------------

# Export a snapshot of specific policies (Global scope / No project specified).
resource "nutanix_network_security_policy_export_v2" "export_selected" {
  policy_ext_ids = [
    nutanix_network_security_policy_v2.policy1.id,
    nutanix_network_security_policy_v2.policy2.id,
  ]
  file_path = "${path.module}/nsp_export_selected_global.bin"
  depends_on = [
    nutanix_network_security_policy_v2.policy1,
    nutanix_network_security_policy_v2.policy2,
  ]
}

# Export every network security policy on the cluster by omitting policy_ext_ids.
resource "nutanix_network_security_policy_export_v2" "export_all" {
  file_path = "${path.module}/nsp_export_all_global.bin"
  depends_on = [
    nutanix_network_security_policy_v2.policy1,
    nutanix_network_security_policy_v2.policy2,
  ]
}

# -----------------------------------------------------------------------------
# EXAMPLES: EXPORT OF POLICY WITH PROJECT SCOPE
# -----------------------------------------------------------------------------

# Export a snapshot of specific policies that belong to a specific project.
# This ensures that only the policies matching both the provided IDs and the
# project context are exported.
resource "nutanix_network_security_policy_export_v2" "export_selected_with_project" {
  policy_ext_ids = [
    nutanix_network_security_policy_v2.policy1.id,
  ]
  project_ext_id = nutanix_project.example_project.id
  file_path      = "${path.module}/nsp_export_selected_project.bin"
  depends_on = [
    nutanix_network_security_policy_v2.policy1
  ]
}

# Export all network security policies that belong to a specific project.
# By omitting policy_ext_ids but providing project_ext_id, it exports every
# policy scoped to that project.
resource "nutanix_network_security_policy_export_v2" "export_all_with_project" {
  project_ext_id = nutanix_project.example_project.id
  file_path      = "${path.module}/nsp_export_all_project.bin"
}

# -----------------------------------------------------------------------------
# OUTPUTS
# -----------------------------------------------------------------------------

output "export_file_path" {
  value = nutanix_network_security_policy_export_v2.export_selected.file_path
}

output "export_task_ext_id" {
  value = nutanix_network_security_policy_export_v2.export_selected.task_ext_id
}
