terraform{
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}

resource "nutanix_image_placement_policy_v2" "example" {
  name           = "image-placement-policy"
  description    = "Image placement policy for the cluster"
  placement_type = "SOFT"

  cluster_entity_filter {
    category_ext_ids = [
      "<cluster_category_id>",
    ]
    type = "CATEGORIES_MATCH_ALL"
  }

  image_entity_filter {
    category_ext_ids = [
        "<image_category_id>",
    ]
    type = "CATEGORIES_MATCH_ALL"
  }
}

# to suspend the image placement policy, just add action = "SUSPEND" in the resource block
resource "nutanix_image_placement_policy_v2" "ipp" {
  name           = "image-placement-policy"
  description    = "Image placement policy for the cluster"
  placement_type = "SOFT"

  cluster_entity_filter {
    category_ext_ids = [
      "<cluster_category_id>",
    ]
    type = "CATEGORIES_MATCH_ALL"
  }

  image_entity_filter {
    category_ext_ids = [
      "<image_category_id>",
    ]
    type = "CATEGORIES_MATCH_ALL"
  }
  action = "SUSPEND"
}


# list image placement policies with filter , pagination and limit
data "nutanix_image_placement_policies_v2" "example"{
  page=0
  limit=10
  filter="startswith(name,'t')"
}

# get image placement policy by id
data "nutanix_image_placement_policy_v2" "example"{
  ext_id = resource.nutanix_image_placement_policy_v2.ipp.id
}