// resources/data_sources used here were introduced in nutanix/nutanix version 1.5.0-beta
terraform{
    required_providers{
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.5.0-beta"
        }
    }
}

// default foundation_port is 8000 so can be ignored
provider "nutanix" {
    // foundation_port = 8000
    foundation_endpoint = "10.xx.xx.xx"
}

/*
Description:
- Here we will image 3 bare metal nodes using ipmi. Will use kvm(ahv) hypervisor bundled in nos package
- ipmi_netmask, ipmi_gateway, ipmi_user & ipmi_password are important fields 
  for ipmi based imaging apart from required fields. These can be declare particularly
  for each node in node spec or default can be declared outside block level as well 
  (foundation treats gives priority to node level declaration).
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

    // this are defaults, you can also mention node specific creds in node spec
    ipmi_user = "<ipmi-username>"
    ipmi_password = "<ipmi-password>"

    // adding one block of nodes
    blocks{
        // adding multiple nodes
        nodes{
            hypervisor_hostname="superman-1"
            hypervisor_ip= "xx.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="xx.xx.xx.xx"
            cvm_ip= "xx.xx.xx.xx"
            node_position= "A"
            // override default password
            ipmi_password= "<node-1-password>"
        }
        nodes{
            hypervisor_hostname="superman-2"
            hypervisor_ip= "xx.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="xx.xx.xx.xx"
            cvm_ip= "xx.xx.xx.xx"
            node_position= "B"
            // override default password
            ipmi_password= "<node-3-password>"
        }
        nodes{
            hypervisor_hostname="superman-3"
            hypervisor_ip= "xx.xx.xx.xx"
            hypervisor= "kvm"
            image_now= true
            ipmi_ip="xx.xx.xx.xx"
            cvm_ip= "xx.xx.xx.xx"
            node_position= "C"
            // override default password
            ipmi_password= "<node-3-password>"
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

output "session" {
    value = resource.nutanix_foundation_image_nodes.batch1
}