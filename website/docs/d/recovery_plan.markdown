---
layout: "nutanix"
page_title: "NUTANIX: nutanix_recovery_plan"
sidebar_current: "docs-nutanix-datasource-recovery-plan"
description: |-
  Describe a Nutanix Recovery Plan and its values (if it has them).
---

# nutanix_recovery_plan

Describe a Nutanix Recovery Plan and its values (if it has them).

## Example Usage

```hcl
resource "nutanix_recovery_plan" "recovery_plan_test" {
    name        = "%s"
    description = "%s"
    stage_list {
        stage_work{
            recover_entities{
                entity_info_list{
                    categories {
                        name = "Environment"
                        value = "Dev"
                    }
                }
            }
        }
        stage_uuid = "ab788130-0820-4d07-a1b5-b0ba4d3a42asd"
        delay_time_secs = 0
    }
    parameters{}
}
```

## Argument Reference

The following arguments are supported:

* `recovery_plan_id`: - (Required) The `id` of the Recovery Plan.

## Attributes Reference

The following attributes are exported:

* `name` The name for the Recovery Plan.
* `description` A description for Recovery Plan.

### Stage List
* `stage_list` - (Required) Input for the stages of the Recovery Plan. Each stage will perform a predefined type of task.
* `stage_list.stage_uuid` - (Optional/Computed) UUID of stage.
* `stage_list.delay_time_secs` - (Optional/Computed) Amount of time in seconds to delay the execution of next stage after execution of current stage.
* `stage_list.stage_work` - (Required) A stage specifies the work to be performed when the Recovery Plan is executed.
* `stage_list.stage_work.0.recover_entities` - (Optional/Computed) Information about entities to be recovered.
* `stage_list.stage_work.0.recover_entities.0.entity_info_list` - (Optional/Computed) Information about entities to be recovered as part of this stage. For VM, entity information will include set of scripts to be executed after recovery of VM. Only one of categories or any_entity_reference has to be provided.
* `stage_list.stage_work.0.recover_entities.0.entity_info_list.#.any_entity_reference_kind` - (Optional/Computed) Reference to a kind.
* `stage_list.stage_work.0.recover_entities.0.entity_info_list.#.any_entity_reference_uuid` - (Optional/Computed) Reference to a uuid.
* `stage_list.stage_work.0.recover_entities.0.entity_info_list.#.any_entity_reference_name` - (Optional/Computed) Reference to a name.
* `stage_list.stage_work.0.recover_entities.0.entity_info_list.#.categories` - (Optional/Computed)  Categories for filtering entities.

### Parameters
* `parameters` - (Required) Parameters for the Recovery Plan.
* `parameters.0.floating_ip_assignment_list` - (Optional/Computed) Floating IP assignment for VMs upon recovery in an Availability Zone. This is applicable only for the public cloud Availability Zones.
* `parameters.0.floating_ip_assignment_list.#.availability_zone_url` - (Required) URL of the Availability Zone.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list` - (Required) IP assignment for VMs upon recovery in the specified Availability Zone.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.test_floating_ip_config` - (Optional/Computed) Configuration for assigning floating IP to a VM on the execution of the Recovery Plan.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.test_floating_ip_config.ip` - (Optional/Computed) IP to be assigned to VM, in case of failover.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.test_floating_ip_config.should_allocate_dynamically` - (Optional/Computed) Whether to allocate the floating IPs for the VMs dynamically.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.recovery_floating_ip_config` - (Optional/Computed) Configuration for assigning floating IP to a VM on the execution of the Recovery Plan.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.recovery_floating_ip_config.ip` - (Optional/Computed) IP to be assigned to VM, in case of failover.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.recovery_floating_ip_config.should_allocate_dynamically` - (Optional/Computed) Whether to allocate the floating IPs for the VMs dynamically.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.vm_reference` - (Required) Reference to a vm.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.vm_reference.kind` - (Required) The kind name.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.vm_reference.uuid` - (Required) The uuid.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.vm_reference.name` - (Optional/Computed) The name.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.vm_nic_information` - (Required) Information about vnic to which floating IP has to be assigned.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.vm_nic_information.ip` - (Optional/Computed) IP address associated with vnic for which floating IP has to be assigned on failover.
* `parameters.0.floating_ip_assignment_list.#.vm_ip_assignment_list.#.vm_nic_information.uuid` - (Required) Uuid of the vnic of the VM to which floating IP has to be assigned.
* `parameters.0.network_mapping_list` - (Required) Network mappings to be used for the Recovery Plan. This will be represented by array of network mappings across the Availability Zones.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list` - (Required) Mapping of networks across the Availability Zones.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.availability_zone_url` - (Optional/Computed) URL of the Availability Zone.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network` - (Optional/Computed) Network configuration to be used for performing network mapping and IP preservation/mapping on Recovery Plan execution.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.virtual_network_reference` - (Optional/Computed) The reference to a virtual_network.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.virtual_network_reference.kind` - (Optional/Computed) The kind name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.virtual_network_reference.uuid` - (Optional/Computed) The uuid.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.virtual_network_reference.name` - (Optional/Computed) The name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.use_vpc_reference` - (Optional/Computed) The reference to a VPC.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.vpc_reference` - (Optional/Computed) The reference to a VPC.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.vpc_reference.kind` - (Optional/Computed) The kind name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.vpc_reference.uuid` - (Optional/Computed) The uuid.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.vpc_reference.name` - (Optional/Computed) The name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.subnet_list` - (Optional/Computed) List of subnets for the network.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.subnet_list.#.gateway_ip` - (Required) Gateway IP address for the subnet.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.subnet_list.#.external_connectivity_state` - (Optional/Computed) External connectivity state of the subnet. This is applicable only for the subnet to be created in public cloud Availability Zone.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.subnet_list.#.prefix_length` - (Required) Prefix length for the subnet.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_network.0.name` - (Required) Name of the network.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network` - (Optional/Computed) Network configuration to be used for performing network mapping and IP preservation/mapping on Recovery Plan execution.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.virtual_network_reference` - (Optional/Computed) The reference to a virtual_network.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.virtual_network_reference.kind` - (Optional/Computed) The kind name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.virtual_network_reference.uuid` - (Optional/Computed) The uuid.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.virtual_network_reference.name` - (Optional/Computed) The name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.subnet_list` - (Optional/Computed) List of subnets for the network.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.subnet_list.#.gateway_ip` - (Required) Gateway IP address for the subnet.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.subnet_list.#.external_connectivity_state` - (Optional/Computed) External connectivity state of the subnet. This is applicable only for the subnet to be created in public cloud Availability Zone.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.subnet_list.#.prefix_length` - (Required) Prefix length for the subnet.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_network.0.name` - (Required) Name of the network.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_ip_assignment_list` - (Optional/Computed) Static IP configuration for the VMs to be applied post recovery in the recovery network for migrate/ failover action on the Recovery Plan.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_ip_assignment_list.0.vm_reference` - (Optional/Computed) The reference to a vm.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_ip_assignment_list.0.vm_reference.kind` - (Optional/Computed) The kind name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_ip_assignment_list.0.vm_reference.uuid` - (Optional/Computed) The uuid.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_ip_assignment_list.0.vm_reference.name` - (Optional/Computed) The name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_ip_assignment_list.0.ip_config_list` - (Optional/Computed) List of IP configurations for a VM.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.recovery_ip_assignment_list.0.ip_config_list.#.ip_address` - (Required) IP address.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_ip_assignment_list` - (Optional/Computed) Static IP configuration for the VMs to be applied post recovery in the test network for test failover action on the Recovery Plan.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_ip_assignment_list.0.vm_reference` - (Optional/Computed) The reference to a vm.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_ip_assignment_list.0.vm_reference.kind` - (Optional/Computed) The kind name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_ip_assignment_list.0.vm_reference.uuid` - (Optional/Computed) The uuid.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_ip_assignment_list.0.vm_reference.name` - (Optional/Computed) The name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_ip_assignment_list.0.ip_config_list` - (Optional/Computed) List of IP configurations for a VM.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_ip_assignment_list.0.ip_config_list.#.ip_address` - (Required) IP address.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.cluster_reference_list` - (Optional/Computed) The clusters where the recovery and test networks reside. This is required to specify network mapping across clusters for a Recovery Plan created to handle failover within the same Availability Zone.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.cluster_reference_list.0.kind` - (Optional/Computed) The kind name.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.cluster_reference_list.0.uuid` - (Optional/Computed) The uuid.
* `parameters.0.network_mapping_list.#.availability_zone_network_mapping_list.#.test_ip_assignment_list.0.name` - (Optional/Computed) The name.


### Metadata
The metadata attribute exports the following:

* `last_update_time` - UTC date and time in RFC-3339 format when vm was last updated.
* `uuid` - vm UUID.
* `creation_time` - UTC date and time in RFC-3339 format when vm was created.
* `spec_version` - Version number of the latest spec.
* `spec_hash` - Hash of the spec. This will be returned from server.
* `name` - vm name.

### Categories
The categories attribute supports the following:

* `name` - the key name.
* `value` - value of the key.

### Reference
The `project_reference`, `owner_reference` attributes supports the following:

* `kind` - (Required) The kind name (Default value: `project`).
* `name` - (Optional) the name.
* `uuid` - (Required) the UUID.

See detailed information in [Nutanix Recovery Plan](https://www.nutanix.dev/api_references/prism-central-v3/#/c0f7aec6fa82b-get-recovery-plan).
