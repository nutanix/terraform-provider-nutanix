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

# Create a template placement policy with SOFT placement type
resource "nutanix_template_placement_policy_v2" "example" {
  name           = "example-template-placement-policy"
  description    = "Example template placement policy"
  placement_type = "SOFT"
  content_filter {
    type             = "CATEGORIES_MATCH_ANY" # Templates matching with this categories should be placed on cluster matching with all the categories in the cluster filter
    category_ext_ids = [var.content_category_ext_id]
  }

  cluster_filter {
    type             = "CATEGORIES_MATCH_ALL"
    category_ext_ids = [var.cluster_category_ext_id]
  }

}

# Singular datasource - fetch by ext_id
data "nutanix_template_placement_policy_v2" "get_policy" {
  ext_id = nutanix_template_placement_policy_v2.example.ext_id
}

# Plural datasource - list all template placement policies
data "nutanix_template_placement_policies_v2" "list_policies" {
  depends_on = [nutanix_template_placement_policy_v2.example]
}

# Plural datasource with filter
data "nutanix_template_placement_policies_v2" "filtered_policies" {
  filter = "name eq '${nutanix_template_placement_policy_v2.example.name}'"
  limit  = 5
}
