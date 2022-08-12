---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_image_nodes"
sidebar_current: "docs-nutanix-resource-foundation-image-nodes"
description: |-
Images node(s) and optionally creates clusters.
---

# nutanix_foundation_image_nodes

Images node(s) and optionally creates clusters.

## Example Usage

```hcl
resource "nutanix_foundation_image_nodes" "batch1" {
    timeouts {
        create = "65m"
    }
    nos_package = "nos_package.tar"
    cvm_netmask = "10.xx.xx.xx"
    cvm_gateway = "10.xx.xx.xx"
    hypervisor_gateway = "10.xx.xx.xx"
    hypervisor_netmask = "10.xx.xx.xx"
    ipmi_gateway = "10.xx.xx.xx"
    ipmi_netmask = "10.xx.xx.xx"
    hypervisor_iso {
        esx {
            filename = esx_image.iso
            checksum = "aasjdajkdsa8sdjnwj2902djncsdc93"
        }
    }
    blocks{
        // incase of undiscovered / manually adding node, ipmi based imaging should be done and IPMI creds are must for that.
        nodes{
            hypervisor_hostname="batman-1"
            cvm_gb_ram = 50
            hypervisor_ip= "10.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="10.xx.xx.xx"
            cvm_ip= "10.xx.xx.xx"
            node_position= "A"
            ipmi_user = "ADMIN"
            ipmi_password = "10.xx.xx.xx"
        }

        // for discovered node, ipmi creds can be ignored & set device_hint = "vm_installer"
        nodes{
            cvm_num_vcpus = 10
            cvm_gb_ram = 51
            hypervisor_hostname="batman-2"
            ipv6_address = "ffff::ffff:ffff:ffff:ffff%eth0"
            current_network_interface = "eth0"
            hypervisor_ip= "10.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="10.xx.xx.xx"
            cvm_ip= "10.xx.xx.xx"
            node_position= "B"
            device_hint = "vm_installer"

        }
        nodes{
            cvm_num_vcpus = 10
            cvm_gb_ram = 51
            hypervisor_hostname="batman-3"
            ipv6_address = "ffff::ffff:ffff:ffff:ffff%eth0"
            hypervisor_ip= "10.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="10.xx.xx.xx"
            current_network_interface = "eth0"
            cvm_ip= "10.xx.xx.xx"
            node_position= "C"
            device_hint = "vm_installer"
        }
        block_id = "999999999"
    }
    blocks {
        nodes {
            cvm_num_vcpus = 10
            cvm_gb_ram = 51
            ipv6_address = "ffff::ffff:ffff:ffff:ffff%eth0"
            current_network_interface = "eth2"
            hypervisor_hostname="superman-1"
            hypervisor_ip= "10.xx.xx.xx"
            hypervisor= "esx"
            image_now= true
            ipmi_ip="10.xx.xx.xx"
            cvm_ip= "10.xx.xx.xx"
            node_position= "D"
            device_hint = "vm_installer"
        }
        block_id = "99999999"
    }
    clusters {
        redundancy_factor = 1
        cluster_name = "superman"
        single_node_cluster = true
        cluster_init_now = true
        cluster_external_ip = "10.xx.xx.xx"
        cluster_members = ["10.xx.xx.xx"]
    }
    clusters {
        redundancy_factor = 2
        cluster_name = "batman"
        cluster_init_now = true
        cluster_external_ip = "10.xx.xx.xx"
        cluster_members = ["10.xx.xx.xx","10.xx.xx.xx","10.xx.xx.xx"]
        timezone = "Africa/Conakry"
    }
}

output "session" {
    value = resource.nutanix_foundation_image_nodes.batch1
}
```

## Argument Reference

The following arguments are supported:

* `ipmi_netmask`: - (Required incase using IPMI based imaging either here or inside node spec) default IPMI netmask
* `ipmi_gateway`: - (Required incase using IPMI based imaging either here or inside node spec) default IPMI gateway
* `ipmi_user` : - (Required incase using IPMI based imaging either here or inside node spec) IPMI username.
* `ipmi_password` : - (Required incase using IPMI based imaging either here or inside node spec) IPMI password.
* `cvm_gateway` : - (Required) CVM gateway.
* `cvm_netmask` : - (Required) CVM netmask.
* `hypervisor_gateway` : - (Required) Hypervisor gateway.
* `hypervisor_netmask` : - (Required) Hypervisor netmask.
* `nos_package` : - (Required) NOS package.
* `blocks` : - (Required) Terraform blocks of Block level parameters.
* `clusters` : - Terraform blocks of clusters config
* `xs_master_label`: - xen server master label.
* `hyperv_external_vnic` : - Hyperv External virtual network adapter name.
* `xen_config_type` : - xen config types.
* `ucsm_ip` : - UCSM IP address.
* `ucsm_password` : - UCSM password.
* `hypervisor_iso` : - Hypervisor ISO.
* `unc_path` : - UNC Path.
* `fc_settings` : - Foundation Central specific settings.
* `xs_master_password` : - xen server master password.
* `svm_rescue_args` : - Arguments to be passed to svm_rescue for AOS installation. Ensure that the arguments provided are supported by the AOS version used for imaging.
* `xs_master_ip` : - xen server master IP address.
* `hyperv_external_vswitch` : - Hyperv External vswitch name.
* `hypervisor_name_server` : - Hypervisor nameserver.
* `hyperv_sku` : - Hyperv SKU.
* `eos_metadata` : - Contains user data from Eos portal.
* `tests` : - Types of tests to be performed.
* `hyperv_product_key` : - Hyperv product key.
* `unc_username` : - UNC username.
* `install_script` : - install script.
* `hypervisor_password` : - Hypervisor password.
* `unc_password` : - UNC password.
* `xs_master_username` : - xen server master username.
* `skip_hypervisor` : - If hypervisor installation should be skipped.
* `ucsm_user` : - UCSM username.
* `layout_egg_uuid` : - Id of the custom layout which needs to be passed to imaging request.

### blocks

The following arguments are supported for each block:
* `block_id` : - Block ID.
* `nodes` :- (Required) Terraform blocks of details of nodes

### nodes

The following arguments are supported for each node:
* `ipmi_netmask` :- (Required incase using IPMI based imaging either here or outside blocks) IPMI netmask for this node
* `ipmi_gateway` :- (Required incase using IPMI based imaging either here or outside blocks) IPMI gateway for this node 
* `ipmi_password` :- (Required incase using IPMI based imaging either here or outside blocks) IPMI username
* `ipmi_user` :- (Required incase using IPMI based imaging either here or outside blocks) IPMI password
* `ipmi_ip` :- (Required) IPMI IP address.
* `hypervisor_hostname` :- (Required) Hypervisor Hostname.
* `hypervisor_ip` :- (Required) Hypervisor IP address.
* `node_position` :- (Required) Position of the node in the block.
* `image_now` :- (Optional, Default = true) If the node should be imaged now.
* `bond_mode` :- (Required if node is capable) dynamic if using LACP, static for LAG
* `rdma_passthrough` :- (Required if node is capable) passthru RDMA nic to CVM if possible, default to false
* `bond_lacp_rate` :- (Required if node is lacp configured) slow or fast if lacp if being used at the switch
* `ipv6_address` :- (Required when device_hint = "vm_installer" for imaging using node's existing cvm for imaging) IPV6 address.
* `ipv6_interface` :- (Required when device_hint = "vm_installer" for imaging using node's existing cvm for imaging) ipv6 interface.
* `image_delay` :- Imaging delay.
* `ucsm_params` :- Object of UCSM parameters.
* `ucsm_params.native_vlan` :- If the vlan is native.
* `ucsm_params.keep_ucsm_settings` :- Whether UCSM settings should be kept.
* `ucsm_params.mac_pool` :- Mac address pool.
* `ucsm_params.vlan_name` :- Name of vlan.
* `cvm_gb_ram` :- RAM capacity of CVM in GB.
* `device_hint` :- use "vm_installer" to enable CVM imaging from standalone.
* `cluster_id` :- ID of cluster.
* `ucsm_node_serial` :- UCSM node serial.
* `node_serial` :- serial number of the node.
* `ipmi_configure_now` :- Whether IPMI should be configured.
* `cvm_num_vcpus` :- Number of CVM vcpus.
* `ipmi_mac` :- IPMI MAC address.
* `rdma_mac_addr` :- mac address of nic to be used for rdma
* `bond_uplinks` :- MAC Addresses of NICs in a team/bond
* `current_network_interface` :- CVM current network interface.
* `vswitches` :- Terraform blocks of vswitch configuration. Foundation will auto-calculate this in most cases. Provide it only if you want to override foundation's defaults.
* `vswitches.lacp` :- Status of LACP.
* `vswitches.bond_mode` :- bond_mode such as balance-tcp, active-backup, etc.
* `vswitches.name` :- Name of the vswitch.
* `vswitches.uplinks` :- Terraform blocks of MAC Addresses of NICs in a team/bond.
* `vswitches.other_config` :- Terraform blocks of Auxiliary lacp configurations. Applicable only for AHV.
* `vswitches.mtu` :- MTU of the vswitch. Applicable only for AHV.
* `ucsm_managed_mode` :- UCSM Managed mode.
* `current_cvm_vlan_tag` :- Current CVM vlan tag. 0 Value with remove vlan tag. 
* `exlude_boot_serial` :- serial of boot device to be excluded (used by NX G6 platforms)
* `mitigate_low_boot_space` :- relocate bootbank files to make space for phoenix files.




### clusters

The following arguments are supported for each cluster:

* `enable_ns` : - If network segmentation should be enabled.
* `backplane_subnet` : - Backplane subnet address.
* `cluster_init_successful` : - If cluster initialization was successful.
* `backplane_netmask` : - Backplane netmask.
* `redundancy_factor` : - (Required) Cluster Redundancy Factor.
* `backplane_vlan` : - Backplane vlan.
* `cluster_name` : - (Required) Name of the cluster.
* `cluster_external_ip` : - External IP of the cluster.
* `cvm_ntp_servers` : - NTP servers of CVM.
* `single_node_cluster` : - If it is a single node cluster.
* `cluster_members` : - (Required) Members in the cluster.
* `cvm_dns_servers` : - DNS servers of CVM.
* `cluster_init_now` : - (Optional, Default = true) If cluster should be created.
* `hypervisor_ntp_servers` : - NTP servers of hypervisor.
* `timezone`: - Set timezone for every CVM

### hypervisor_iso

The following arguments are supported:

* `hyperv` : - Details of hypervisor ISO of type hyperv.
* `kvm` : - Details of hypervisor ISO of type kvm.
* `xen` : - Details of hypervisor ISO of type xen.
* `esx` : - Details of hypervisor ISO of type esx.

### hypervisor iso reference

Each of `hypervisor_iso.hyperv`, `hypervisor_iso.kvm`,
`hypervisor_iso.xen` & `hypervisor_iso.esx` would have following arguments:

* `filename` :- (Required) Checksum for ISO file.
* `checksum` :- (Required) Filename of hypervisor ISO.

### fc_settings

The following arguments are supported:

* `fc_metadata` :- Foundation Central metadata which will be transferred to the newly imaged node.
* `fc_metadata.fc_ip` :- IP address of foundation central.
* `fc_metadata.api_key` :- api_key which the node uses to register itself with foundation central.
* `foundation_central` :- If this attribute is set to True, FC workflow will be invoked.


### eos_metadata

The following arguments are supported:

* `config_id` : - Id of the Eos config uploaded in foundation GUI.
* `account_name` : - arrya of account names
* `email` : - Email address of the user who downloaded Eos config.

### tests

The following arguments are supported:

* `run_syscheck` : - Whether system checks should run.
* `run_ncc` : - Whether NCC checks should run.

## Attributes Reference

The following attributes are exported:

* `id` : - unique id of terraform resouce is set to session_id of the imaging session
* `session_id` : - session_id of the imaging session
* `cluster_urls` :- list containing cluster name and cluster urls for created clusters in current session
* `cluster_urls.#.cluster_name` :- cluster_name 
* `cluster_urls.#.cluster_url` :- url to access the cluster login

## Defaults 

The attributes like `ipmi_netmask`, `ipmi_gateway`, `ipmi_user` & `ipmi_password` can be mentioned for a node as well as for all nodes outside blocks. This attributes if mentioned in node will be used for that particular node.

## Error 

Incase of error in any individual entity i.e. node or cluster, terraform will error our after full imaging process is completed. Error will be shown for every failed node and cluster.

## lifecycle

* `Update` : - Resource will trigger new resource create call for any kind of update in resource config.
* `delete` : - Delete will be a soft delete.

See detailed information in [Nutanix Foundation Image Nodes](https://www.nutanix.dev/api_references/foundation/#/b3A6MjIyMjMzOTQ-image-a-given-set-of-nodes).
