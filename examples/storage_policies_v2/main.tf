terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.4.0"
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

resource "nutanix_storage_policy_v2" "example" {
  # Required: Storage Policy name (max 64 characters, must be unique)
  name = "my-storage-policy"

  # Optional: Compression specification
  compression_spec {
    # Required: Compression state
    # Valid values: "DISABLED", "POSTPROCESS", "INLINE", "SYSTEM_DERIVED"
    compression_state = "POSTPROCESS"
  }

  # Optional: Encryption specification
  encryption_spec {
    # Required: Encryption state
    # Valid values: "SYSTEM_DERIVED", "ENABLED"
    # Note: Once set to "ENABLED", it cannot be reverted
    encryption_state = "ENABLED"
  }

  # Optional: Quality of Service specification
  qos_spec {
    # Required: Throttled IOPS (range: 100 to 2147483647)
    throttled_iops = 1000
  }

  # Optional: Fault Tolerance specification
  fault_tolerance_spec {
    # Required: Replication factor
    # Valid values: "SYSTEM_DERIVED", "TWO", "THREE"
    # TWO = Original + 1 copy, THREE = Original + 2 copies
    replication_factor = "THREE"
  }

  # Optional: List of category external IDs (0-20 items), 
  # Apply policy to specific categories
  # Each ID must be a valid UUID format
  # Category external IDs can be fetched from the data source "nutanix_categories_v2"
  # Example:
  # data "nutanix_categories_v2" "category-list" {
  #   filter = "key eq 'category_key'"
  # }
  # category_ext_ids = [
  #   data.nutanix_categories_v2.category-list.categories.0.ext_id,
  #   data.nutanix_categories_v2.category-list.categories.1.ext_id
  # ]
  category_ext_ids = [
    "4d552748-e119-540a-b06c-3c6f0d213fa2",
    "5e663859-f220-651b-c17d-4d7f0e324fb3"
  ]
}


# get storage policy by ext id
data "nutanix_storage_policy_v2" "fetch" {
  ext_id = nutanix_storage_policy_v2.example.id
}

# list of storage policies
data "nutanix_storage_policies_v2" "storage-policies"{ }