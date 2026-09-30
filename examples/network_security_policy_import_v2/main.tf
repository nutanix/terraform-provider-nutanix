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

# This example demonstrates a full export -> import round trip:
#   1. create a category and a network security policy,
#   2. export it to a local file,
#   3. import that file back (recreating the policy from the data file).
#
# In a real scenario the data file usually comes from another cluster; here it
# is produced in the same run so the example is self-contained and ready to run.

resource "nutanix_category_v2" "cat1" {
  key   = "tf-example-import-cat"
  value = "tf-example-import-cat-value"
}

resource "nutanix_network_security_policy_v2" "policy1" {
  name        = "tf-example-import-nsp"
  description = "Network Security Policy example for the import data file"
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

# Produce the data file that the import below consumes.
resource "nutanix_network_security_policy_export_v2" "export" {
  policy_ext_ids = [nutanix_network_security_policy_v2.policy1.id]
  file_path      = "${path.module}/nsp_export.bin"
  depends_on     = [nutanix_network_security_policy_v2.policy1]
}

# Bulk import of the network security policies contained in the data file.
# ntnx_purge_policies = false keeps existing policies; set it to true to delete
# all existing policies before importing.
resource "nutanix_network_security_policy_import_v2" "import" {
  path                = nutanix_network_security_policy_export_v2.export.file_path
  ntnx_purge_policies = false
  depends_on          = [nutanix_network_security_policy_export_v2.export]
}

output "import_task_ext_id" {
  value = nutanix_network_security_policy_import_v2.import.task_ext_id
}

# The ext_ids (UUIDs) of the policies created by the import.
output "imported_policy_ext_ids" {
  value = nutanix_network_security_policy_import_v2.import.imported_policy_ext_ids
}
