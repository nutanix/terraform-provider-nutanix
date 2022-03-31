// resources/data_sources used here were introduced in nutanix/nutanix version >=1.4.2
terraform{
    required_providers{
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.4.2"
        }
    }
}

// default foundation_port is 8000 so can be ignored
provider "nutanix" {
    // foundation_port = 8000
    foundation_endpoint = "10.xx.xx.xx"
}

/*
Node is considered nutanix imaged node when:
- Node is already have cvm & hypervisor running and you want to reimage as per your requirement
- Node is having Discovery OS running and you want to image it with aos and appropriate hypervisor
This kind of nodes are discoverable from foundation using ipv6 address network discovery. You can check
response of nutanix_foundation_discover_nodes data source. This nodes can be imaged using ipmi (standalone foundation), 
example can be found in examples under foundation/image_bare_metal. Or this node can be imaged using cvm/discovery os running
on imaged node (example given below).

Description:
- Here we want to image nutanix imaged nodes using the node's cvm itself.
- ipmi_user & ipmi_password creds are not mandatory here incase node's cvm is used (by making device_installer=true in node spec).

Notes : Most information for a node can be directly pulled from nutanix_foundation_discover_nodes & nutanix_foundation_node_network_details
        datasources and can be substitued.
*/

// Fetch all aos packages details
data "nutanix_foundation_nos_packages" "nos" {}

resource "nutanix_foundation_image_nodes" "batch1" {

    // custom timeout, default is 60 minutes
    timeouts {
        create = "65m"
    }

    // assuming theres only 1 nos package present in foundation vm
    nos_package = data.nutanix_foundation_nos_packages.nos.entities[0]

    // cvm, hypervisor & ipmi common details
    cvm_netmask = "xx.xx.xx.xx"
    cvm_gateway = "xx.xx.xx.xx"
    hypervisor_gateway = "xx.xx.xx.xx"
    hypervisor_netmask = "xx.xx.xx.xx"
    ipmi_gateway = "xx.xx.xx.xx"
    ipmi_netmask = "xx.xx.xx.xx"

    // use this incase you want to use specific hypervisor iso for a type
    // nos package mentioned above already have ahv(kvm) package bundled in so for that no need to mention here
    # hypervisor_iso {
    #     esx {
    #         filename = "filename.iso"
    #         checksum = "xxxxxxxxxxxxxxxxxxxxx"
    #     }
    #     kvm {
    #         filename = "filename.iso"
    #         checksum = "xxxxxxxxxxxxxxxxxxxx"
    #     }
    # }

    // adding one block of nodes
    blocks{
        // adding multiple nodes
        nodes{
            ipv6_address = "xx:xx:xx:xx:xx:xx:xx:xx"
            current_network_interface = "eth0"
            hypervisor_hostname="superman-1"
            hypervisor_ip= "xx.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="xx.xx.xx.xx"
            cvm_ip= "xx.xx.xx.xx"
            node_position= "A"
            // this will make foundation consider the node as discovered node and use cvm of this node for mounting iso
            device_hint = "vm_installer"
        }
        nodes{
            ipv6_address = "xx:xx:xx:xx:xx:xx:xx:xx"
            current_network_interface = "eth0"
            hypervisor_hostname="superman-2"
            hypervisor_ip= "xx.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="xx.xx.xx.xx"
            cvm_ip= "xx.xx.xx.xx"
            node_position= "B"
             // this will make foundation consider the node as discovered node and use cvm of this node for mounting iso
            device_hint = "vm_installer"
        }
        nodes{
            ipv6_address = "xx:xx:xx:xx:xx:xx:xx:xx"
            current_network_interface = "eth0"
            hypervisor_hostname="superman-3"
            hypervisor_ip= "xx.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="xx.xx.xx.xx"
            cvm_ip= "xx.xx.xx.xx"
            node_position= "C"
             // this will make foundation consider the node as discovered node and use cvm of this node for mounting iso
            device_hint = "vm_installer"
        }
        block_id = "xxxxxx"
    }
    
    // add cluster block
    clusters {
        redundancy_factor = 2
        cluster_name = "superman"
        single_node_cluster = false // not required. make it true for single node cluster creation
        cluster_init_now = true
        cluster_external_ip = "xx.xx.xx.xx"
        cluster_members = ["xx.xx.xx.xx","xx.xx.xx.xx","xx.xx.xx.xx"]
    }
}

// output the details after imaging is finished
output "session" {
    value = resource.nutanix_foundation_image_nodes.batch1
}