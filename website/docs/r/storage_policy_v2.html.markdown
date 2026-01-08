---
layout: "nutanix"
page_title: "NUTANIX: nutanix_storage_policy_v2"
sidebar_current: "docs-nutanix-resource-storage-policy-v2"
description: |-
  Create Storage Policies
---

# nutanix_storage_policy_v2

Provides Nutanix resource to create storage policy

-> **Note**:
Once `encryption_state` is explicitly set to `ENABLED`, it cannot be reverted back to a system-derived value.

If compression_state, encryption_state, or replication_factor are intended to be system-derived, ensure that the qos_spec block is included.



## Example

### Complete Example with All Fields

```hcl
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
  category_ext_ids = [
    "4d552748-e119-540a-b06c-3c6f0d213fa2",
    "5e663859-f220-651b-c17d-4d7f0e324fb3"
  ]
}
```

### Minimal Example with System-Derived Values

```hcl
resource "nutanix_storage_policy_v2" "minimal" {
  name = "minimal-storage-policy"

  # When using system-derived values, qos_spec must be included
  qos_spec {
    throttled_iops = 100
  }

  # Compression, encryption, and replication will use system defaults
  compression_spec {
    compression_state = "SYSTEM_DERIVED"
  }

  encryption_spec {
    encryption_state = "SYSTEM_DERIVED"
  }

  fault_tolerance_spec {
    replication_factor = "SYSTEM_DERIVED"
  }
}
```

### Example with Inline Compression

```hcl
resource "nutanix_storage_policy_v2" "inline_compression" {
  name = "inline-compression-policy"

  compression_spec {
    compression_state = "INLINE"
  }

  encryption_spec {
    encryption_state = "ENABLED"
  }

  qos_spec {
    throttled_iops = 5000
  }

  fault_tolerance_spec {
    # TWO = Original + 1 copy
    replication_factor = "TWO"
  }
}
```

### Example with Disabled Compression

```hcl
resource "nutanix_storage_policy_v2" "no_compression" {
  name = "no-compression-policy"

  compression_spec {
    compression_state = "DISABLED"
  }

  qos_spec {
    throttled_iops = 2000
  }

  fault_tolerance_spec {
    # THREE = Original + 2 copies
    replication_factor = "THREE"
  }
}
```

### Example with Categories Only

```hcl
# Rest all will be System Derived
resource "nutanix_storage_policy_v2" "categorized" {
  name = "categorized-policy"

  # Apply policy to specific categories
  category_ext_ids = [
    "4d552748-e119-540a-b06c-3c6f0d213fa2"
  ]

  qos_spec {
    throttled_iops = 100
  }
}
```

## Argument Reference

The following arguments are supported:

* `name`:- (Required) Storage Policy name. Must be unique and cannot exceed 64 characters.
* `category_ext_ids`:- (Optional) List of external identifiers for Categories to be included in the Storage Policy. Each ID must be a valid UUID format. Maximum 20 items allowed.
* `compression_spec`:- (Optional) Defines compression parameters for entities governed by the Storage Policy.
* `encryption_spec`:- (Optional) Defines encryption parameters for entities governed by the Storage Policy.
* `qos_spec`:- (Optional) Defines Storage Quality of Service (QOS) parameters for the entities.
* `fault_tolerance_spec`:- (Optional) Defines Fault Tolerance parameters for the entities.

### Compression Spec

The `compression_spec` block supports the following:

* `compression_state`:- (Required) Controls enabling or disabling compression. If no explicit preference is set, the system chooses a value.
  Valid values:
  * `"DISABLED"`:- User wants data not compressed.
  * `"POSTPROCESS"`:- User wants data compressed later.
  * `"INLINE"`:- User wants data compressed inline.
  * `"SYSTEM_DERIVED"`:- User is not interested in compression; system decides.

### Encryption Spec

The `encryption_spec` block supports the following:

* `encryption_state`:- (Required) Controls enabling encryption. Once enabled, it cannot be disabled. If no explicit preference is set, the system decides.
  Valid values:
  * `"SYSTEM_DERIVED"`:- User is not interested in encryption; system decides.
  * `"ENABLED"`:- User wants data encrypted.

  -> **Note**: Once `encryption_state` is explicitly set to `ENABLED`, it cannot be reverted back to a system-derived value.

### QoS Spec

The `qos_spec` block supports the following:

* `throttled_iops`:- (Required) Specifies throttled IOPS for governed entities. The block size for IO is 32kB. Valid range: 100 to 2147483647.

### Fault Tolerance Spec

The `fault_tolerance_spec` block supports the following:

* `replication_factor`:- (Required) Specifies the number of data copies for entities governed by the Storage Policy.
  Valid values:
  * `"SYSTEM_DERIVED"`:- User has not provided the number of copies; system decides.
  * `"TWO"`:- Two data copies (Original + 1 copy).
  * `"THREE"`:- Three data copies (Original + 2 copies).

## Attribute Reference

The following attributes are exported:

* `ext_id`:- (Computed) External identifier of the Storage Policy.
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity.
* `links`:- A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `name`:- Storage Policy name.
* `category_ext_ids`:- List of external identifiers for Categories included in the Storage Policy.
* `compression_spec`:- Compression parameters for entities governed by the Storage Policy.
  * `compression_state`:- Compression state value.
* `encryption_spec`:- Encryption parameters for entities governed by the Storage Policy.
  * `encryption_state`:- Encryption state value.
* `qos_spec`:- Storage Quality of Service (QOS) parameters for the entities.
  * `throttled_iops`:- Throttled IOPS value.
* `fault_tolerance_spec`:- Fault Tolerance parameters for the entities.
  * `replication_factor`:- Replication factor value.
* `policy_type`:- (Computed) Indicates whether the policy is user-created or system-created. Valid values: `"USER"`, `"SYSTEM"`.



See detailed information in [Nutanix Create Storage Policies v4](https://developers.nutanix.com/api-reference?namespace=datapolicies&version=v4.1#tag/StoragePolicies/operation/createStoragePolicy).