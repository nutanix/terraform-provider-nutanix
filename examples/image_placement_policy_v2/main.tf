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


#pull all categories data
data "nutanix_categories_v2" "categories" {}

locals {
  category_ext_id = data.nutanix_categories_v2.categories.categories.0.ext_id
}

resource "nutanix_image_placement_policy_v2" "ipp" {
  name           = "image-placement-policy"
  description    = "Image placement policy for the cluster"
  placement_type = "SOFT"
  cluster_entity_filter {
    category_ext_ids = [
      local.category_ext_id,
    ]
    type = "CATEGORIES_MATCH_ALL"
  }
  image_entity_filter {
    category_ext_ids = [
      local.category_ext_id,
    ]
    type = "CATEGORIES_MATCH_ALL"
  }

  lifecycle {
    ignore_changes = [
      cluster_entity_filter,
      image_entity_filter,
    ]
  }
}

# to suspend the image placement policy, just add action = "SUSPEND" in the resource block
resource "nutanix_image_placement_policy_v2" "ipp-suspend" {
  name           = "image-placement-policy-SUSPEND"
  description    = "Image placement policy for the cluster"
  placement_type = "SOFT"
  cluster_entity_filter {
    category_ext_ids = [
      local.category_ext_id,
    ]
    type = "CATEGORIES_MATCH_ALL"
  }
  image_entity_filter {
    category_ext_ids = [
      local.category_ext_id,
    ]
    type = "CATEGORIES_MATCH_ALL"
  }
  action = "SUSPEND"
}

# to resume the image placement policy, just add action = "RESUME" in the resource block
resource "nutanix_image_placement_policy_v2" "ipp-resume" {
  name           = "image-placement-policy-RESUME"
  description    = "Image placement policy for the cluster"
  placement_type = "SOFT"
  cluster_entity_filter {
    category_ext_ids = [
      local.category_ext_id,
    ]
    type = "CATEGORIES_MATCH_ALL"
  }
  image_entity_filter {
    category_ext_ids = [
      local.category_ext_id,
    ]
    type = "CATEGORIES_MATCH_ALL"
  }
  action = "RESUME"
}


# list all image placement policies
data "nutanix_image_placement_policies_v2" "list-ipp" {
  depends_on = [nutanix_image_placement_policy_v2.ipp, nutanix_image_placement_policy_v2.ipp-suspend, nutanix_image_placement_policy_v2.ipp-resume]
}
# list image placement policies with filter , pagination and limit
data "nutanix_image_placement_policies_v2" "filtered-ipp" {
  filter = "name eq '${nutanix_image_placement_policy_v2.ipp.name}'"
}

# list image placement policies with  pagination and limit
data "nutanix_image_placement_policies_v2" "paginated-ipp" {
  page       = 2
  limit      = 10
  depends_on = [nutanix_image_placement_policy_v2.ipp, nutanix_image_placement_policy_v2.ipp-suspend, nutanix_image_placement_policy_v2.ipp-resume]

}

# get image placement policy by id
data "nutanix_image_placement_policy_v2" "get-ipp" {
  ext_id = nutanix_image_placement_policy_v2.ipp.id
}
