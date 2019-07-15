#############################################################################
# Example main.tf for Nutanix + Terraform
#
# Author: jon@nutanix.com
#
# This script is a quick demo of how to use the following provider objects:
# - providers
#     - terraform-provider-nutanix
# - resources
#     - nutanix_virtual_machine
#     - nutanix_subnet
#     - nutanix_image
# - data sources
#     - nutanix_clusters
#     -
#     -
# - script Variables
#     - clusterid's for targeting clusters within prism central
#
# Feel free to reuse, comment, and contribute, so that others may learn.
#
#############################################################################
### Define Provider Info for terraform-provider-nutanix
### This is where you define the credentials for ** Prism Central **
###
### NOTE:
###   While it may be possible to use Prism Element directly, Nutanix's
###   provider is not structured or tested for this. Using Prism Central will
###   give the broadest capabilities across the board
/*  provider "nutanix" {
  username  = "admin"
  password  = "Nutanix/1234"
  endpoint  = "10.5.80.255"
  insecure  = true
  port      = 9440
}  */

data "nutanix_clusters" "clusters" {
}

### Define Script Local Variables
### This can be used for any manner of things, but is useful for like clusterid,
###   to store a mapping of targets for provisioning
### TODO: Need to make clusters a data source object, such that consumers do
###       not need to manually provision cluster ID
locals {
  cluster1 = data.nutanix_clusters.clusters.entities[1].metadata.uuid
}

##########################
### Data Sources
##########################
### These are "lookups" to simply define an already existing object as a
### plain text name
### This is useful when managing a nutanix prism central instance from multiple
### state files, or deploying terraform into an existing / brownfield environment
### Virtual Machine Data Sources
# data "nutanix_virtual_machine" "nutanix_virtual_machine" {
#   vm_id = nutanix_virtual_machine.vm1.id
# }
### Image Data Sources
# data "nutanix_image" "test" {
#     metadata = {
#         kind = "image"
#     }
#     image_id = nutanix_image.test.id
# }
### Subnet Data Sources
# data "nutanix_subnet" "next-iac-managed" {
#     metadata = {
#         kind = "subnet"
#     }
#    image_id = nutanix_subnet.next-iac-managed.id
#}
### Cluster Data Sources
#data "nutanix_image" "test" {
#    metadata = {
#        kind = "image"
#    }
#    image_id = nutanix_image.test.id
#}
##########################
### Resources
##########################
### Image Resources (Managed by the Image Service)
## Related Product Docs:
##   https://portal.nutanix.com/#/page/docs/details?targetId=Prism-Central-Guide-Prism-v58:mul-images-manage-pc-c.html
## Related Developer Docs:
##   http://developer.nutanix.com/reference/prism_central/v3/#images
###
### Images are raw ISO, QCOW2, or VMDK files that are uploaded by a user can be
###attached to a VM.
### An ISO image is attached as a virtual CD-ROM drive, and QCOW2 and VMDK files
### are attached as SCSI disks.
### For self service portal use cases: An image has to be explicitly added to
### the self-service catalog before users can create VMs from it.

# This demo used a single dummy image from a local filer. Multiple images can be
#   presented here as separate resources, or existing images on a cluster can be
#   called in as data sources, which you can see in the data sources section
#   above.
resource "nutanix_image" "cirros-034-disk" {
  name = "cirros-034-disk"

  #source_uri  = "http://endor.dyn.nutanix.com/acro_images/DISKs/cirros-0.3.4-x86_64-disk.img"
  source_uri  = "http://download.cirros-cloud.net/0.3.4/cirros-0.3.4-x86_64-disk.img"
  description = "heres a tiny linux image, not an iso, but a real disk!"
}

### Subnet Resources (Virtual Networks within AHV)
## Related Product Docs:
##   https://portal.nutanix.com/#/page/docs/details?targetId=Prism-Central-Guide-Prism-v58:mul-network-configuration-acropolis-pc-t.html
## Related Developer Docs:
##   http://developer.nutanix.com/reference/prism_central/v3/#subnet
## Implementation Notes on Subnets
# Subnets are virtual networks (VLANs) and can either be standard L2 VLANs or
# standard L2 VLANs with additional L3 IPAM provided by Acropolis.
#
# Using Managed Networks (aka AHV's internal IPAM) allows AHV to fully control IP
# address assignment, and the DHCP DORA request does not even hit the wire. This
# can be used where AHV fully manages the DHCP environment *OR* it can be used
# with 3rd party IPAMs, where AHV does *NOT* have any DHCP pool; however, still
# has the capability to deterministically pass in layer 3 IPv4 addresses.
#
# This is especially interesting for environments that already have a 3rd party
# IPAM, as it gives you two options:
# Option A: Use the 3rd party IPAM exactly how you'd use it anywhere else, in
#           which case you'd just set up AHV's VLANs as standard L2 VLANs, and
#           then have the enterprise network setup just like usual.
# Option B: Use AHV's managed networks as a very neat way to inject IP addresses
#           into VMs, instead of using anything on the physical network side to
#           act as an IP Helper/DHCP relay.
#
# More specifically, with "Option B" you can create a "managed" network in AHV via
# terraform that doesn't have a DHCP pool assoicated, then using the assign IP
# address feature to "pass" that static IP into the guest as an acropolis
# managed IP.
#
# Let me be clear: This helps avoid physical network switch configuration as well
# as any dependency on things like cloud-init or sysprep to statically assign an
# IP address. This you could grab an IP as a resource from a 3rd party IPAM, and
# pass that IP Object into the VM via this feature, without having to interact
# with your physical network team at all.
#
# Note: an AHV managed IP address will "stick" with the VM from cradle to
# grave, even if the VM is off.

# ### Define Terraform Managed Subnets
resource "nutanix_subnet" "infra-managed-network-140" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information
  name        = "infra-managed-network-140"
  vlan_id     = 140
  subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
  subnet_ip = "172.21.32.0"

  default_gateway_ip = "172.21.32.1"
  prefix_length      = 24

  dhcp_options = {
    boot_file_name   = "bootfile"
    domain_name      = "ntnxlab"
    tftp_server_name = "172.21.32.200"
  }

  dhcp_server_address = {
    ip = "172.21.32.254"
  }

  dhcp_domain_name_server_list = ["172.21.30.223"]
  dhcp_domain_search_list      = ["ntnxlab.local"]
  #ip_config_pool_list_ranges   = ["172.21.32.3 172.21.32.253"] 
}

### Virtual Machine Resources
## Related Product Docs:
##    https://portal.nutanix.com/#/page/docs/details?targetId=Prism-Central-Guide-Prism-v58:mul-vm-create-manage-pc-c.html
## Related Developer Docs:
##    http://developer.nutanix.com/reference/prism_central/v3/#vms
## Implementation Notes on VMs
# These are VMs managed by Prism Central, which could span across AHV, ESXi, or
# Prism Self Service Portal. That said, if you're doing ESXi, it is most likely
# that would deploy them via a VMware provider against vCenter APIs.

resource "nutanix_virtual_machine" "demo-01-web" {
  # General Information
  name                 = "demo-01-web"
  description          = "demo Frontend Web Server"
  num_vcpus_per_socket = 2
  num_sockets          = 1
  memory_size_mib      = 4096

  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # What networks will this be attached to?
  nic_list {
    # subnet_reference is saying, which VLAN/network do you want to attach here?
    subnet_uuid = nutanix_subnet.infra-managed-network-140.id
    # Used to set static IP.
    # ip_endpoint_list {
    #   ip   = "172.21.32.20"
    #   type = "ASSIGNED"
    # }
  }

  # What disk/cdrom configuration will this have?
  disk_list {
    # data_source_reference in the Nutanix API refers to where the source for
    # the disk device will come from. Could be a clone of a different VM or a
    # image like we're doing here
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
  disk_list {
    # defining an additional entry in the disk_list array will create another.

    #disk_size_mib and disk_size_bytes must be set together.
    disk_size_mib   = 100000
    disk_size_bytes = 104857600000
  }
  #Using provisioners
  #Use as the following provisioner block if you know that you are geeting an reachable IP address.
  #Get ssh connection and execute commands.
  # provisioner "remote-exec" {
  #   connection {
  #     user     = "cirros"    # user from the image attached
  #     password = "cubswin:)" #password from the user 
  #     host    = "172.21.32.20" #host is now a required value for connection, you can use `self.nic_list_status[0].ip_endpoint_list[0].ip` to set the IP or if you know the IP you could set manually.
  #   }

  #   inline = [
  #     "echo \"Hello World\"",
  #   ]
  # }
}

# Show IP address
output "ip_address" {
  value = nutanix_virtual_machine.demo-01-web.nic_list_status.0.ip_endpoint_list[0]["ip"]
}

