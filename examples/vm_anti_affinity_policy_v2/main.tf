terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
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

#creating category
resource "nutanix_category_v2" "category" {
  key   = "AntiAffinityPolicy"
  value = "tf-anti-affinity-policy"
}

# Create policy
resource "nutanix_vm_anti_affinity_policy_v2" "policy" {
  name        = "temp-anti-affinity"
  description = "a description"
  categories  = [ nutanix_category_v2.category.id ]
}

# Get all policies
data "nutanix_vm_anti_affinity_policies_v2" "policies" {}

# Get specific policy
data "nutanix_vm_anti_affinity_policy_v2" "policy" {
  ext_id = nutanix_vm_anti_affinity_policy_v2.policy.ext_id
}