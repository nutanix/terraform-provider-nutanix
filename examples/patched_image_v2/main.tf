terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

# defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# A claim token is required to create a patched image. The patched image
# references it via claim_token_ext_id (ClaimTokens cross-resource reference).
resource "nutanix_claim_token_v2" "example" {
  name            = var.claim_token_name
  expiry_time     = var.claim_token_expiry_time
  max_usage_count = 5
}

# Create a patched hypervisor image customized for specific nodes.
resource "nutanix_patched_image_v2" "example" {
  claim_token_ext_id = nutanix_claim_token_v2.example.ext_id
  name               = var.patched_image_name
  host_type          = var.host_type

  image_details {
    local_host_image_ext_id = var.local_host_image_ext_id
  }

  node_configurations {
    node_ext_id = var.node_ext_id

    host_configuration {
      hostname = var.hostname

      network_details {
        management {
          ip {
            ipv4 {
              value         = var.node_ip
              prefix_length = 24
            }
          }
          gateway {
            ipv4 {
              value = var.node_gateway
            }
          }
          vlan_id = 0
        }
      }
    }
  }
}

# Read a single patched image by its external id.
data "nutanix_patched_image_v2" "get-patched-image" {
  ext_id = nutanix_patched_image_v2.example.ext_id
}

# List all patched images.
data "nutanix_patched_images_v2" "list-patched-images" {
  depends_on = [nutanix_patched_image_v2.example]
}
