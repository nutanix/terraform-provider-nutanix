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

#creating Host category
resource "nutanix_category_v2" "host_category" {
  key   = "HostAffinityPolicy"
  value = "tf-host-affinity-host-category"
}

# Create VM category
resource "nutanix_category_v2" "vm_category" {
  key   = "HostAffinityPolicy"
  value = "tf-host-affinity-vm-category"
}

# Create policy
resource "nutanix_vm_host_affinity_policy_v2" "policy" {
  name        = "temp-host-affinity"
  description = "a description"
  host_categories  = [ nutanix_category_v2.host_category.id ]
  vm_categories  = [ nutanix_category_v2.vm_category.id ]
}

# Get all policies
data "nutanix_vm_host_affinity_policies_v2" "policies" {}

# Get specific policy
data "nutanix_vm_host_affinity_policy_v2" "policy" {
  ext_id = nutanix_vm_host_affinity_policy_v2.policy.ext_id
}