---
layout: "nutanix"
page_title: "NUTANIX: nutanix_network_security_rule"
sidebar_current: "docs-nutanix-resource-network-security-rule"
description: |-
  Provides a Nutanix Network Security Rule resource to Create a Network Security Rule .
---

# nutanix_network_security_rule

Provides a Nutanix network security rule resource to Create a network security rule.

> NOTE: The use of network_security_rule is only applicable in AHV clusters and requires Microsegmentation to be enabled. This feature is a function of the Flow product and requires a Flow license. For more information on Flow and Microsegmentation please visit https://www.nutanix.com/products/flow

## Example Usage

### Isolation Rule Example

```hcl
resource "nutanix_network_security_rule" "isolation" {
	name        = "example-isolation-rule"
	description = "Isolation Rule Example"
	
	isolation_rule_action = "APPLY"
	
	isolation_rule_first_entity_filter_kind_list = ["vm"]
	isolation_rule_first_entity_filter_type      = "CATEGORIES_MATCH_ALL"
	isolation_rule_first_entity_filter_params {
		name   = "Environment"
		values = ["Dev"]
	}
	
	isolation_rule_second_entity_filter_kind_list = ["vm"]
	isolation_rule_second_entity_filter_type      = "CATEGORIES_MATCH_ALL"
	isolation_rule_second_entity_filter_params {
		name   = "Environment"
		values = ["Production"]
	}
}
```

### App Rule Example with associated VMs.
```hcl
data "nutanix_clusters" "clusters" {}

locals {
  cluster_uuid = [
    for cluster in data.nutanix_clusters.clusters.entities :
    cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
  ][0]
}

//Create categories.
resource "nutanix_category_key" "test-category-key" {
  name        = "TIER-1"
  description = "TIER Category Key"
}

resource "nutanix_category_key" "USER" {
  name        = "user"
  description = "user Category Key"
}


resource "nutanix_category_value" "WEB" {
  name        = "${nutanix_category_key.test-category-key.id}"
  description = "WEB Category Value"
  value       = "WEB-1"
}

resource "nutanix_category_value" "APP" {
  name        = "${nutanix_category_key.test-category-key.id}"
  description = "APP Category Value"
  value       = "APP-1"
}

resource "nutanix_category_value" "DB" {
  name        = "${nutanix_category_key.test-category-key.id}"
  description = "DB Category Value"
  value       = "DB-1"
}

resource "nutanix_category_value" "group" {
  name        = "${nutanix_category_key.USER.id}"
  description = "group Category Value"
  value       = "group-1"
}


//Create a cirros image
resource "nutanix_image" "cirros-034-disk" {
  name        = "test-image-vm-create-flow"
  source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
  description = "heres a tiny linux image, not an iso, but a real disk!"
}

//APP-1 VM.
resource "nutanix_virtual_machine" "vm-app" {
  name                 = "test-dou-vm-flow-APP-1"
  cluster_uuid         = local.cluster_uuid
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186

  nic_list {
    subnet_uuid = "c56b535c-8aff-4435-ae85-78e64a07f76d"
  }

  disk_list {
    data_source_reference = {
      kind = "image"
      uuid = nutanix_image.cirros-034-disk.id
    }

    device_properties {
      disk_address = {
        device_index = 0
        adapter_type = "SCSI"
      }
      device_type = "DISK"
    }
  }

  categories {
    name  = "Environment"
    value = "Staging"
  }

  categories {
    name  = "TIER-1"
    value = nutanix_category_value.APP.id
  }
}

#WEB-1 VM
resource "nutanix_virtual_machine" "vm-web" {
  name                 = "test-dou-vm-flow-WEB-1"
  cluster_uuid         = local.cluster_uuid
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186

  nic_list {
    subnet_uuid = "c56b535c-8aff-4435-ae85-78e64a07f76d"
  }

  disk_list {
    data_source_reference = {
      kind = "image"
      uuid = nutanix_image.cirros-034-disk.id
    }

    device_properties {
      disk_address = {
        device_index = 0
        adapter_type = "SCSI"
      }
      device_type = "DISK"
    }
  }

  categories {
    name  = "Environment"
    value = "Staging"
  }

  categories {
    name  = "TIER-1"
    value = nutanix_category_value.WEB.id 
  }
}

#DB-1 VM
resource "nutanix_virtual_machine" "vm-db" {
  name                 = "test-dou-vm-flow-DB-1"
  cluster_uuid         = local.cluster_uuid
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186

  nic_list {
    subnet_uuid = "c56b535c-8aff-4435-ae85-78e64a07f76d"
  }

  disk_list {
    data_source_reference = {
      kind = "image"
      uuid = nutanix_image.cirros-034-disk.id
    }

    device_properties {
      disk_address = {
        device_index = 0
        adapter_type = "SCSI"
      }
      device_type = "DISK"
    }
  }
  
  categories {
    name  = "Environment"
    value = "Staging"
  }

  categories {
    name  = "TIER-1"
    value = nutanix_category_value.DB.id
  }
}

//Create Application Network Policy.
resource "nutanix_network_security_rule" "TEST-TIER" {
  name        = "RULE-1-TIERS"
  description = "rule 1 tiers"

  app_rule_action = "APPLY"

  app_rule_inbound_allow_list {
    peer_specification_type = "FILTER"
    filter_type             = "CATEGORIES_MATCH_ALL"
    filter_kind_list        = ["vm"]

    filter_params {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.WEB.id}"]
    }
  }

  
  app_rule_target_group_default_internal_policy = "DENY_ALL"
  
  app_rule_target_group_peer_specification_type = "FILTER"
  
  app_rule_target_group_filter_type = "CATEGORIES_MATCH_ALL"
  
  app_rule_target_group_filter_kind_list = ["vm"]
  
  app_rule_target_group_filter_params {
    name   = "${nutanix_category_key.test-category-key.id}"
    values = ["${nutanix_category_value.APP.id}"]
  }
  app_rule_target_group_filter_params {
    name   = "${nutanix_category_key.USER.id}"
    values = ["${nutanix_category_value.group.id}"]
  }

  app_rule_target_group_filter_params {
    name   = "AppType"
    values = ["Default"]
  }

  app_rule_outbound_allow_list {
    peer_specification_type = "FILTER"
    filter_type             = "CATEGORIES_MATCH_ALL"
    filter_kind_list        = ["vm"]

    filter_params {
      name   = "${nutanix_category_key.test-category-key.id}"
      values = ["${nutanix_category_value.DB.id}"]
    }
  }

  depends_on = [nutanix_virtual_machine.vm-app, nutanix_virtual_machine.vm-web, nutanix_virtual_machine.vm-db]
}
```

### Usage with service and address groups
```hcl
resource "nutanix_service_group" "service1" {
  name = "srv-1"
  description = "test"

  service_list {
      protocol = "TCP"
      tcp_port_range_list {
        start_port = 22
        end_port = 22
      }
      tcp_port_range_list {
        start_port = 2222
        end_port = 2222
      }
  }
}

resource "nutanix_address_group" "address1" {
  name = "addr-1"
  description = "test"

  ip_address_block_list {
    ip = "10.0.0.0"
    prefix_length = 24
  }
}

resource "nutanix_category_value" "ad-group-user-1" {
	name = "AD"
	description = "group user category value"
	value = "AD"
}

resource "nutanix_network_security_rule" "VDI" {
	name           = "nsr-1"
	ad_rule_action = "APPLY"
	description    = "test"
	#   app_rule_action = "APPLY"
	ad_rule_inbound_allow_list {
		  ip_subnet               = "10.0.0.0"
		  ip_subnet_prefix_length = "8"
		  peer_specification_type = "IP_SUBNET"
		  protocol                = "ALL"

    #  peer_specification_type = "ALL"
    #  service_group_list {
    #    kind = "service_group"
    #    uuid = nutanix_service_group.service1.id
    #  }
    #  address_group_inclusion_list {
    #    kind = "address_group"
    #    uuid = nutanix_address_group.address1.id
    #  }
	}
	ad_rule_target_group_default_internal_policy = "DENY_ALL"
	ad_rule_target_group_filter_kind_list = [
	  "vm"
	]
	ad_rule_target_group_filter_params {
	  name = "AD"
	  values = [
		"AD"
	  ]
	}
	ad_rule_target_group_filter_type             = "CATEGORIES_MATCH_ALL"
	ad_rule_target_group_peer_specification_type = "FILTER"
	ad_rule_outbound_allow_list {
		peer_specification_type = "ALL"
    service_group_list {
        kind = "service_group"
        uuid = nutanix_service_group.service1.id
      }

    address_group_inclusion_list {
        kind = "address_group"
        uuid = nutanix_address_group.address1.id
      }
	}
	depends_on = [nutanix_category_value.ad-group-user-1]
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) The name for the network_security_rule.
* `categories`: - (Optional) Categories for the network_security_rule.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.
* `description`: - (Optional) A description for network_security_rule.
* `quarantine_rule_action`: - (Optional) These rules are used for quarantining suspected VMs. Target group is a required attribute. Empty inbound_allow_list will not allow anything into target group. Empty outbound_allow_list will allow everything from target group.
* `quarantine_rule_outbound_allow_list`: - (Optional)
* `quarantine_rule_target_group_default_internal_policy`: - (Optional) - Default policy for communication within target group.
* `quarantine_rule_target_group_peer_specification_type`: - (Optional) - Way to identify the object for which rule is applied.
* `quarantine_rule_target_group_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `quarantine_rule_target_group_filter_type`: - (Optional) - The type of the filter being used.
* `quarantine_rule_target_group_filter_params`: - (Optional) - A list of category key and list of values.
* `quarantine_rule_inbound_allow_list`: - (Optional)
* `app_rule_action`: - (Optional) - These rules govern what flows are allowed. Target group is a required attribute. Empty inbound_allow_list will not anything into target group. Empty outbound_allow_list will allow everything from target group.
* `app_rule_outbound_allow_list`: - (Optional)
* `app_rule_target_group_default_internal_policy`: - (Optional) - Default policy for communication within target group.
* `app_rule_target_group_peer_specification_type`: - (Optional) - Way to identify the object for which rule is applied.
* `app_rule_target_group_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `app_rule_target_group_filter_type`: - (Optional) - The type of the filter being used.
* `app_rule_target_group_filter_params`: - (Optional) - A list of category key and list of values.
* `app_rule_inbound_allow_list`: - (Optional) The set of categories that matching VMs need to have.
* `ad_rule_action`: - (Optional) - These rules govern what flows are allowed. Target group is a required attribute. Empty inbound_allow_list will not anything into target group. Empty outbound_allow_list will allow everything from target group.
* `ad_rule_outbound_allow_list`: - (Optional)
* `ad_rule_target_group_default_internal_policy`: - (Optional) - Default policy for communication within target group.
* `ad_rule_target_group_peer_specification_type`: - (Optional) - Way to identify the object for which rule is applied.
* `ad_rule_target_group_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `ad_rule_target_group_filter_type`: - (Optional) - The type of the filter being used.
* `ad_rule_target_group_filter_params`: - (Optional) - A list of category key and list of values.
* `ad_rule_inbound_allow_list`: - (Optional) The set of categories that matching VMs need to have.
* `isolation_rule_action`: - (Optional) - These rules are used for environmental isolation.
* `app_rule_outbound_allow_list`: - (Optional)
* `isolation_rule_first_entity_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `isolation_rule_first_entity_filter_type`: - (Optional) - The type of the filter being used.
* `isolation_rule_first_entity_filter_params`: - (Optional) - A list of category key and list of values.
* `isolation_rule_second_entity_filter_kind_list`: - (Optional) - List of kinds associated with this filter.
* `isolation_rule_second_entity_filter_type`: - (Optional) - The type of the filter being used.
* `isolation_rule_second_entity_filter_params`: - (Optional) - A list of category key and list of values.

## Attributes Reference

The following attributes are exported:

* `metadata`: - The network_security_rule kind metadata.
* `retrieval_uri_list`: - List of URIs where the raw network_security_rule data can be accessed.
* `size_bytes`: - The size of the network_security_rule in bytes.
* `state`: - The state of the Network Security Rule.
* `api_version` - The version of the API.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when image was last updated.
* `UUID`: - image UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when image was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - image name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, attributes supports the following:

* `kind`: - The kind name (Default value: project)(Required).
* `name`: - the name(Optional).
* `uuid`: - the UUID(Required).

See detailed information in [Nutanix Network Security Rule](https://www.nutanix.dev/api_references/prism-central-v3/#/8e7c7fb305664-create-a-new-network-security-rule).
