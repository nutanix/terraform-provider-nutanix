terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0"
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

# Create a category for the cluster entity filter
resource "nutanix_category_v2" "example" {
  key         = "rate-limit-example"
  value       = "cluster-filter"
  description = "Category for image rate limit policy example"
}

# Create an image rate limit policy
resource "nutanix_image_rate_limit_policy_v2" "example" {
  name            = "example-rate-limit-policy"
  description     = "Example image rate limit policy"
  rate_limit_kbps = 1000

  cluster_entity_filter {
    category_ext_ids = [nutanix_category_v2.example.ext_id]
    type             = "CATEGORIES_MATCH_ALL"
    # type             = "CATEGORIES_MATCH_ALL"
    # Cluster matching with all the categories specified in the cluster_entity_filter --> category_ext_ids, will be enrolled for the rate limit policy.

    # type             = "CATEGORIES_MATCH_ANY"
    # Cluster matching with any of the categories specified in the cluster_entity_filter --> category_ext_ids, will be enrolled for the rate limit policy.
  }
}

# Get a singular image rate limit policy by ext_id
data "nutanix_image_rate_limit_policy_v2" "get_policy" {
  ext_id = nutanix_image_rate_limit_policy_v2.example.ext_id
}

# List all image rate limit policies
data "nutanix_image_rate_limit_policies_v2" "list_policies" {
  depends_on = [nutanix_image_rate_limit_policy_v2.example]
}

# List effective rate limit policies across clusters
data "nutanix_effective_image_rate_limit_policies_v2" "effective_policies" {
  depends_on = [nutanix_image_rate_limit_policy_v2.example]
}
