#############################################################################
# Demo Multi-Tier App Deployment
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
#     - nutanix_virtual_machine
#     -
#     -
# - script Variables
#     - clusterid's for targeting clusters within prism central
#
# The goal of this script is to show a multi-tier applicaiton, interacting with
# these various objects.
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
provider "nutanix" {
  username  = "admin"
  password  = "Nutanix/1234"
  endpoint  = "10.5.80.255"
  insecure  = true
  port      = 9440
}

### Define Script Local Variables
### This can be used for any manner of things, but is useful for like clusterid, to store a mapping of targets for provisioning
### TODO: Need to make clusters a data source object, such that consumers do not need to manually provision cluster ID
locals {
  cluster1   = "00054051-250f-5ccc-0000-00000000cf0d"
  ip_haproxy = "10.5.94.11"
  ip_app     = "10.5.94.12"
  ip_db      = "10.5.94.13"
}

##########################
### Data Sources
##########################
### These are "lookups" to simply define an already existing object as a plain text name
### This is useful when managing a nutanix prism central instance from multiple state files, or deploying terraform into an existing / brownfield environment
### Virtual Machine Data Sources
# data "nutanix_virtual_machine" "nutanix_virtual_machine" {
#   vm_id = "${nutanix_virtual_machine.vm1.id}"
# }
### Image Data Sources
# data "nutanix_image" "test" {
#     metadata = {
#         kind = "image"
#     }
#     image_id = "${nutanix_image.test.id}"
# }
### Subnet Data Sources
# data "nutanix_subnet" "next-iac-managed" {
#     metadata = {
#         kind = "subnet"
#     }
#    image_id = "${nutanix_subnet.next-iac-managed.id}"
#}
### Cluster Data Sources
#data "nutanix_image" "test" {
#    metadata = {
#        kind = "image"
#    }
#    image_id = "${nutanix_image.test.id}"
#}
##########################
### Resources
##########################
### Image Resources (Managed by the Image Service)
###
### Images are raw ISO, QCOW2, or VMDK files that are uploaded by a user can be attached to a VM.
### An ISO image is attached as a virtual CD-ROM drive, and QCOW2 and VMDK files are attached as SCSI disks.
### For self service portal use cases: An image has to be explicitly added to the self-service catalog before users can create VMs from it.
###
### Image Service Product Docs: https://portal.nutanix.com/#/page/docs/details?targetId=Prism-Central-Guide-Prism-v56:mul-images-manage-pc-c.html
### Image Service Developer Docs: http://developer.nutanix.com/reference/prism_central/v3/#images
# resource "nutanix_image" "centos-lamp-app" {
#   # General Information
#   name        = "CentOS-LAMP-APP.qcow2"
#   description = "CentOS LAMP - App"
#   source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-APP.qcow2"

#   metadata = {
#     kind = "image"
#   }
# }

# resource "nutanix_image" "centos-lamp-db" {
#   # General Information
#   name        = "CentOS-LAMP-DB.qcow2"
#   description = "CentOS LAMP - DB"
#   source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-DB.qcow2"

#   metadata = {
#     kind = "image"
#   }
# }

# resource "nutanix_image" "centos-lamp-haproxy" {
#   # General Information
#   name        = "CentOS-LAMP-HAPROXY.qcow2"
#   description = "CentOS LAMP - HAProxy"
#   source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-HAProxy.qcow2"

#   metadata = {
#     kind = "image"
#   }
# }

# resource "nutanix_image" "cirros-034-disk" {
#     name        = "cirros-034-disk"
#     source_uri  = "http://endor.dyn.nutanix.com/acro_images/DISKs/cirros-0.3.4-x86_64-disk.img"
#     description = "heres a tiny linux image, not an iso, but a real disk!"
#     metadata = {
#         kind = "image"
#     }
# }

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
# resource "nutanix_subnet" "infra-managed-network-140" {
#   metadata = {
#     kind = "subnet"
#   }

#   # What cluster will this VLAN live on?
#   cluster_reference = {
#     kind = "cluster"
#     uuid = "${local.cluster1}"
#   }

#   # General Information
#   name        = "next-iac-managed"
#   description = "NEXT"
#   vlan_id     = 0
#   subnet_type = "VLAN"

#   # Provision a Managed L3 Network
#   # This bit is only needed if you intend to turn on AHV's IPAM
# 	subnet_ip          = "10.250.140.0"
#   default_gateway_ip = "10.250.140.1"
#   prefix_length      = 24
#   dhcp_options {
# 		boot_file_name   = "bootfile"
# 		domain_name      = "nutanix"
# 		domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
# 		domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
# 		tftp_server_name = "10.250.140.200"
#   }
# }

### Virtual Machine Resources
## Related Product Docs:
##    https://portal.nutanix.com/#/page/docs/details?targetId=Prism-Central-Guide-Prism-v58:mul-vm-create-manage-pc-c.html
## Related Developer Docs:
##    http://developer.nutanix.com/reference/prism_central/v3/#vms
## Implementation Notes on Subnets
# These are VMs managed by Prism Central, which could span across AHV, ESXi, or
# Prism Self Service Portal. That said, if you're doing ESXi, it is most likely
# that would deploy them via a VMware provider against vCenter APIs. That said,
# the capability does exist here.

resource "nutanix_virtual_machine" "demo-01-web" {
  # WRT virtual machine, metadata section may only be useful when setting the "projects" and "categories" constructs
  metadata {
    kind = "vm"
  }

  # General Information
  name                 = "demo-01-web"
  description          = "demo Frontend Web Server"
  num_vcpus_per_socket = 2
  num_sockets          = 1
  memory_size_mib      = 4096
  power_state          = "ON"

  # What cluster will this VLAN live on?
  cluster_reference = {
    kind = "cluster"
    uuid = "${local.cluster1}"
  }

  # What networks will this be attached to?
  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.next-iac-managed.id}"
    }

    ip_endpoint_list = {
      ip   = "${local.ip_haproxy}"
      type = "ASSIGNED"
    }
  }]

  # What disk/cdrom configuration will this have?
  disk_list = [{
    data_source_reference = [{
      kind = "image"
      name = "Centos7"
      uuid = "${nutanix_image.centos-lamp-haproxy.id}"
    }]

    device_properties = [{
      device_type = "DISK"
    }]

    disk_size_mib = 5000
  }]
}

resource "nutanix_virtual_machine" "demo-01-app" {
  # Can I hard code image to be kind image?
  # We're going to make this implicit in future API releases, so hard coding it is safe on the plugin side
  # WRT virtual machine, metadata section may only be useful when setting the "projects" and "categories" constructs
  metadata {
    kind = "vm"
  }

  # General Information
  name                 = "demo-01-app"
  description          = "Demo Java middleware App server"
  num_vcpus_per_socket = 2
  num_sockets          = 1
  memory_size_mib      = 8192
  power_state          = "ON"

  # What cluster will this VLAN live on?
  cluster_reference = {
    kind = "cluster"
    uuid = "${local.cluster1}"
  }

  # What networks will this be attached to?
  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.next-iac-managed.id}"
    }

    ip_endpoint_list = {
      ip   = "${local.ip_app}"
      type = "ASSIGNED"
    }
  }]

  #What disk/cdrom configuration will this have?
  disk_list = [{
    data_source_reference = [{
      kind = "image"
      name = "Centos7"
      uuid = "${nutanix_image.centos-lamp-app.id}"
    }]

    device_properties = [{
      device_type = "DISK"
    }]

    disk_size_mib = 5000
  }]
}

resource "nutanix_virtual_machine" "demo-01-db" {
  # Can I hard code image to be kind image?
  # We're going to make this implicit in future API releases, so hard coding it is safe on the plugin side
  # WRT virtual machine, metadata section may only be useful when setting the "projects" and "categories" constructs
  metadata {
    kind = "vm"
  }

  # General Information
  name                 = "demo-01-db"
  description          = "demo MySQL Database Server"
  num_vcpus_per_socket = 4
  num_sockets          = 1
  memory_size_mib      = 16384
  power_state          = "ON"

  # What cluster will this VLAN live on?
  cluster_reference = {
    kind = "cluster"
    uuid = "${local.cluster1}"
  }

  #What networks will this be attached to?
  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.next-iac-managed.id}"
    }

    ip_endpoint_list = {
      ip   = "${local.ip_db}"
      type = "ASSIGNED"
    }
  }]

  # What disk/cdrom configuration will this have?
  disk_list = [{
    data_source_reference = [{
      kind = "image"
      name = "Centos7"
      uuid = "${nutanix_image.centos-lamp-db.id}"
    }]

    device_properties = [{
      device_type = "DISK"
    }]

    disk_size_mib = 5000
  }]
}
